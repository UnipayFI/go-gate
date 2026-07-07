package crossex

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*CrossexClient)(nil)

// CrossexClient is the REST client for Gate's cross-exchange (CrossEx) endpoints
// under /api/v4/crossex/*. It embeds the shared core client, so every product
// client reuses the same signing/transport layer.
type CrossexClient struct {
	*client.Client
}

// NewCrossexClient constructs a crossex REST client.
func NewCrossexClient(options ...client.Options) *CrossexClient {
	return &CrossexClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *CrossexClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
