package client

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/UnipayFI/go-gate/pkg/log"
	"github.com/gorilla/websocket"
)

// WebSocketClient holds the configuration for a single Gate v4 stream gateway.
// Gate splits streams by product across different hosts (spot on
// api.gateio.ws/ws/v4/, perpetual futures on fx-ws.gateio.ws per settle
// currency, delivery and options on their own paths), so each product package
// constructs this core with its gateway URL. Credentials are only needed for
// private channels, whose subscribe frames carry a per-request auth object.
type WebSocketClient struct {
	url          string
	apiKey       string
	apiSecret    string
	signFn       SignFn
	logger       log.Logger
	dialer       *websocket.Dialer
	timeOffsetMs int64
}

type WebSocketOption struct {
	url          string
	apiKey       string
	apiSecret    string
	signFn       SignFn
	logger       log.Logger
	dialer       *websocket.Dialer
	timeOffsetMs int64
}

type WebSocketOptions func(*WebSocketOption)

func defaultWebSocketOption() *WebSocketOption {
	return &WebSocketOption{
		logger: log.GetDefaultLogger(),
		dialer: defaultDialer(),
	}
}

func defaultDialer() *websocket.Dialer {
	return &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true,
	}
}

// NewWebSocketClient constructs a stream client for the given gateway URL. The
// url is normally supplied by a product package via WithWebSocketURL.
func NewWebSocketClient(defaultURL string, options ...WebSocketOptions) *WebSocketClient {
	opt := defaultWebSocketOption()
	for _, option := range options {
		option(opt)
	}
	u := defaultURL
	if opt.url != "" {
		u = opt.url
	}
	return &WebSocketClient{
		url:          u,
		apiKey:       opt.apiKey,
		apiSecret:    opt.apiSecret,
		signFn:       opt.signFn,
		logger:       opt.logger,
		dialer:       opt.dialer,
		timeOffsetMs: opt.timeOffsetMs,
	}
}

func (c *WebSocketClient) GetURL() string               { return c.url }
func (c *WebSocketClient) GetAPIKey() string            { return c.apiKey }
func (c *WebSocketClient) GetAPISecret() string         { return c.apiSecret }
func (c *WebSocketClient) GetSignFn() SignFn            { return c.signFn }
func (c *WebSocketClient) GetLogger() log.Logger        { return c.logger }
func (c *WebSocketClient) GetDialer() *websocket.Dialer { return c.dialer }

func (c *WebSocketClient) GetTimeOffsetMs() int64 { return c.timeOffsetMs }

// SetTimeOffset stores a client/server clock offset in milliseconds so that
// private subscription signatures carry a timestamp Gate accepts. Set it once at
// startup (before subscribing) — typically by copying a REST client's offset
// obtained from SyncServerTime — since Gate signs the WebSocket subscribe frame's
// time field and rejects skewed private subscriptions.
func (c *WebSocketClient) SetTimeOffset(offsetMs int64) { c.timeOffsetMs = offsetMs }

// TimestampMs returns the current subscription timestamp in milliseconds,
// adjusted by the configured client/server clock offset. The subscribe frame's
// whole-second time field is derived from this value.
func (c *WebSocketClient) TimestampMs() int64 {
	return time.Now().UnixMilli() - c.timeOffsetMs
}

// WithWebSocketAuth sets the credentials used to sign private channel
// subscriptions. Gate WebSocket has no passphrase.
func WithWebSocketAuth(apiKey, apiSecret string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		opt.apiKey = apiKey
		opt.apiSecret = apiSecret
	}
}

// WithWebSocketURL overrides the gateway URL. Empty is ignored.
func WithWebSocketURL(u string) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.url = u }
}

// WithWebSocketTimeOffset sets a fixed client/server clock offset in milliseconds
// used when signing private subscriptions (subscribe timestamp = localMillis -
// offsetMs). Usually copied from a REST client's SyncServerTime result. The
// offset can also be set after construction via SetTimeOffset.
func WithWebSocketTimeOffset(offsetMs int64) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.timeOffsetMs = offsetMs }
}

func WithWebSocketLogger(logger log.Logger) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.logger = logger }
}

// WithWebSocketSignFn overrides the default HMAC-SHA512 subscription signer.
func WithWebSocketSignFn(signFn SignFn) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.signFn = signFn }
}

// WithWebSocketProxy routes the stream dial through the given proxy (http,
// https, socks5, socks5h). Invalid URLs are logged and skipped.
func WithWebSocketProxy(proxyURL string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		if proxyURL == "" {
			return
		}
		u, err := url.Parse(proxyURL)
		if err != nil {
			opt.logger.Errorf("WithWebSocketProxy: invalid proxy URL %q: %v", proxyURL, err)
			return
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			opt.dialer.Proxy = http.ProxyURL(u)
			opt.dialer.NetDialContext = nil
		case "socks5", "socks5h":
			dialCtx, err := socks5DialContext(u)
			if err != nil {
				opt.logger.Errorf("WithWebSocketProxy: socks5 setup failed: %v", err)
				return
			}
			opt.dialer.Proxy = nil
			opt.dialer.NetDialContext = dialCtx
		default:
			opt.logger.Errorf("WithWebSocketProxy: unsupported scheme %q", u.Scheme)
		}
	}
}
