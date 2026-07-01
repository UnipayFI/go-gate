package request

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/pkg/log"
	"github.com/gorilla/websocket"
)

// WsClient is what the subscribe framework needs from a *client.WebSocketClient.
type WsClient interface {
	GetURL() string
	GetAPIKey() string
	GetAPISecret() string
	GetSignFn() SignFn
	GetLogger() log.Logger
	GetDialer() *websocket.Dialer
	// TimestampMs returns the current subscription timestamp in milliseconds,
	// adjusted by the configured client/server clock offset. Gate signs the
	// subscribe frame's whole-second time field, so calibrating this keeps
	// private subscriptions inside Gate's acceptance window on skewed clocks.
	TimestampMs() int64
}

// wsAuth is the per-request credential object attached to private subscriptions.
type wsAuth struct {
	Method string `json:"method"`
	Key    string `json:"KEY"`
	Sign   string `json:"SIGN"`
}

// wsRequest is a Gate v4 stream control frame (subscribe / unsubscribe).
type wsRequest struct {
	Time    int64   `json:"time"`
	ID      *int64  `json:"id,omitempty"`
	Channel string  `json:"channel"`
	Event   string  `json:"event"`
	Payload any     `json:"payload,omitempty"`
	Auth    *wsAuth `json:"auth,omitempty"`
}

// WsPush is the envelope Gate pushes for a stream event. Time is unix seconds,
// TimeMs unix milliseconds; Result carries the channel-specific typed payload.
type WsPush[T any] struct {
	Time    time.Time `json:"time,format:unix"`
	TimeMs  time.Time `json:"time_ms,format:unixmilli"`
	ID      *int64    `json:"id,omitempty"`
	Channel string    `json:"channel"`
	Event   string    `json:"event"`
	Result  T         `json:"result"`
}

// wsHeader is a lightweight view used to classify an inbound frame (data push vs
// subscribe/unsubscribe ack vs error) before committing to a typed decode.
type wsHeader struct {
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Error   *WsError `json:"error"`
}

// WsError is a Gate WebSocket error object.
type WsError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *WsError) Error() string {
	return fmt.Sprintf("<WsError> code=%d, msg=%s", e.Code, e.Message)
}

// wsSign signs a private subscription: hex(HMAC-SHA512(secret, "channel=X&event=Y&time=Z")).
func wsSign(client WsClient, channel, event string, ts int64) (string, error) {
	prehash := fmt.Sprintf("channel=%s&event=%s&time=%d", channel, event, ts)
	if fn := client.GetSignFn(); fn != nil {
		return fn(client.GetAPISecret(), prehash)
	}
	mac := hmac.New(sha512.New, common.StringToBytes(client.GetAPISecret()))
	mac.Write(common.StringToBytes(prehash))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// Subscribe opens a dedicated connection to the client's gateway, subscribes to
// channel with payload (signing when private), and invokes cb for every data
// push, decoding the push's Result into *T. It returns a done channel (close to
// unsubscribe + disconnect) and a stop channel (closed when the reader exits).
func Subscribe[T any](ctx context.Context, client WsClient, channel string, payload any, private bool, cb func(*WsPush[T], error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeRaw(ctx, client, channel, payload, private, func(message []byte, e error) {
		if e != nil {
			cb(nil, e)
			return
		}
		var push WsPush[T]
		if uerr := common.JSONUnmarshal(message, &push); uerr != nil {
			cb(nil, uerr)
			return
		}
		cb(&push, nil)
	})
}

// SubscribeRaw is like Subscribe but delivers each data frame's raw bytes.
func SubscribeRaw(ctx context.Context, client WsClient, channel string, payload any, private bool, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeRaw(ctx, client, channel, payload, private, cb)
}

func subscribeRaw(ctx context.Context, client WsClient, channel string, payload any, private bool, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	conn, _, err := client.GetDialer().DialContext(ctx, client.GetURL(), nil)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadLimit(10 << 20)

	if err := writeSubscribe(client, conn, channel, "subscribe", payload, private); err != nil {
		conn.Close()
		return nil, nil, err
	}

	doneC := make(chan struct{})
	stopC := make(chan struct{})
	// silent suppresses callback delivery once the caller closes doneC to stop:
	// the ReadMessage error from the watcher closing the conn must not reach cb.
	var silent atomic.Bool

	go keepAlive(conn, common.DEFAULT_KEEP_ALIVE_INTERVAL)
	go func() {
		select {
		case <-stopC:
			silent.Store(true)
		case <-doneC:
			silent.Store(true)
		}
		// Best-effort unsubscribe before closing.
		_ = writeSubscribe(client, conn, channel, "unsubscribe", payload, private)
		conn.Close()
	}()
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !silent.Load() {
					cb(nil, err)
				}
				close(stopC)
				return
			}
			client.GetLogger().Debugf("ws recv: %s", common.BytesToString(message))

			var hdr wsHeader
			if err := common.JSONUnmarshal(message, &hdr); err != nil {
				cb(nil, err)
				continue
			}
			if hdr.Error != nil {
				if !silent.Load() {
					cb(nil, hdr.Error)
				}
				// A subscribe rejection is fatal for this connection.
				if hdr.Event == "subscribe" || hdr.Event == "unsubscribe" {
					close(stopC)
					return
				}
				continue
			}
			switch hdr.Event {
			case "subscribe", "unsubscribe", "":
				// control acks (result: {status:"success"}); ignore.
			default:
				cb(message, nil)
			}
		}
	}()
	return doneC, stopC, nil
}

func writeSubscribe(client WsClient, conn *websocket.Conn, channel, event string, payload any, private bool) error {
	// Whole-second timestamp derived from the client's clock offset (0 when
	// unset) so private subscribe signatures stay inside Gate's window.
	ts := client.TimestampMs() / 1000
	req := wsRequest{Time: ts, Channel: channel, Event: event, Payload: payload}
	if private {
		sign, err := wsSign(client, channel, event, ts)
		if err != nil {
			return err
		}
		req.Auth = &wsAuth{Method: "api_key", Key: client.GetAPIKey(), Sign: sign}
	}
	data, err := common.JSONMarshal(req)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}

// keepAlive sends a periodic protocol-level ping; Gate replies with pong,
// keeping quiet channels' connections alive.
func keepAlive(conn *websocket.Conn, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
			return
		}
	}
}
