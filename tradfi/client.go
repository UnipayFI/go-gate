package tradfi

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*TradfiClient)(nil)

// TradfiClient is the REST client for Gate's TradFi (traditional-finance / MT5)
// endpoints under /api/v4/tradfi/*. It embeds the shared core client, so every
// product client reuses the same signing/transport layer. Unlike the core Gate
// v4 endpoints, TradFi wraps every payload in a {timestamp, data, ...} business
// envelope, so each Service returns that envelope struct.
type TradfiClient struct {
	*client.Client
}

// NewTradfiClient constructs a TradFi REST client.
func NewTradfiClient(options ...client.Options) *TradfiClient {
	return &TradfiClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset and stores it so that
// signed requests carry a timestamp Gate accepts. Gate rejects requests whose
// Timestamp drifts more than a minute from its own clock, so call this once at
// startup (and periodically for long-lived processes). TradFi has no public
// time endpoint, so it aligns against the shared /api/v4/spot/time clock.
func (c *TradfiClient) SyncServerTime(ctx context.Context) error {
	localBefore := time.Now().UnixMilli()
	raw, err := request.DoRaw(request.Get(ctx, c, "/api/v4/spot/time"))
	if err != nil {
		return err
	}
	localAfter := time.Now().UnixMilli()
	var out struct {
		ServerTime int64 `json:"server_time"`
	}
	if err := common.JSONUnmarshal(raw, &out); err != nil {
		return err
	}
	c.SetTimeOffset((localBefore+localAfter)/2 - out.ServerTime)
	return nil
}
