package futures

import (
	"context"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*FuturesClient)(nil)

// FuturesClient is the REST client for Gate's perpetual-futures endpoints under
// /api/v4/futures/{settle}/*. The settle currency (usdt / btc) is a per-request
// path argument, so one client serves every settlement.
type FuturesClient struct {
	*client.Client
}

// NewFuturesClient constructs a perpetual-futures REST client.
func NewFuturesClient(options ...client.Options) *FuturesClient {
	return &FuturesClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *FuturesClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
