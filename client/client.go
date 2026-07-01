package client

import (
	"time"

	"github.com/UnipayFI/go-gate/v4/pkg/log"
	"github.com/go-resty/resty/v2"
)

// Client is the shared, product-agnostic REST core. Every Gate business line
// (spot, margin, unified, perpetual futures, delivery, options, earn, ...) is
// just a set of /api/v4/* request paths layered on top of this same signing +
// transport machinery, so the core carries no product-specific state.
type Client struct {
	client *resty.Client

	apiKey       string
	apiSecret    string
	logger       log.Logger
	signFn       SignFn
	timeOffsetMs int64
}

func NewClient(options ...Options) *Client {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}

	baseURL := opt.network.RestBaseURL()
	if opt.baseURL != "" {
		baseURL = opt.baseURL
	}
	opt.client.SetBaseURL(baseURL)

	return &Client{
		client:       opt.client,
		apiKey:       opt.apiKey,
		apiSecret:    opt.apiSecret,
		logger:       opt.logger,
		signFn:       opt.signFn,
		timeOffsetMs: opt.timeOffsetMs,
	}
}

func (c *Client) GetHttpClient() *resty.Client { return c.client }

func (c *Client) GetAPIKey() string { return c.apiKey }

func (c *Client) GetAPISecret() string { return c.apiSecret }

func (c *Client) GetLogger() log.Logger { return c.logger }

func (c *Client) GetSignFn() SignFn { return c.signFn }

func (c *Client) GetTimeOffsetMs() int64 { return c.timeOffsetMs }

func (c *Client) SetTimeOffset(offsetMs int64) { c.timeOffsetMs = offsetMs }

// TimestampMs returns the current request timestamp in milliseconds, adjusted
// by the configured client/server clock offset. Gate's signature carries the
// timestamp in whole seconds (derived from this value); syncing the offset keeps
// signed requests inside Gate's acceptance window.
func (c *Client) TimestampMs() int64 {
	return time.Now().UnixMilli() - c.timeOffsetMs
}
