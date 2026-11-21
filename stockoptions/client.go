package stockoptions

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*StockOptionsClient)(nil)

// StockOptionsClient is the REST client for Gate's Stock Options endpoints under
// /api/v4/stock_options/*. It embeds the shared core client, so it reuses the
// same signing/transport layer as every other product client. Like stock and
// tradfi, every endpoint wraps its payload in a {label, data, timestamp}
// envelope (timestamp in Unix milliseconds).
type StockOptionsClient struct {
	*client.Client
}

// NewStockOptionsClient constructs a Stock Options REST client.
func NewStockOptionsClient(options ...client.Options) *StockOptionsClient {
	return &StockOptionsClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *StockOptionsClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
