package stock

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*StockClient)(nil)

// StockClient is the REST client for Gate's Stock endpoints under
// /api/v4/stock/*. It embeds the shared core client, so every product client
// reuses the same signing/transport layer. The Stock business provides
// traditional-finance stock spot trading; like the other envelope-wrapped
// products it returns an APIV4-compatible {data, timestamp} body on success and
// {label, message, data, timestamp} on failure, so each Service returns that
// envelope struct.
type StockClient struct {
	*client.Client
}

// NewStockClient constructs a Stock REST client.
func NewStockClient(options ...client.Options) *StockClient {
	return &StockClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset and stores it so that
// signed requests carry a timestamp Gate accepts. Gate rejects requests whose
// Timestamp drifts more than a minute from its own clock, so call this once at
// startup (and periodically for long-lived processes). Stock has no public time
// endpoint, so it aligns against the shared /api/v4/spot/time clock.
func (c *StockClient) SyncServerTime(ctx context.Context) error {
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
