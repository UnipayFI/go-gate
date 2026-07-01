package flashswap

import (
	"context"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*FlashSwapClient)(nil)

// FlashSwapClient is the REST client for Gate's flash-swap endpoints under /api/v4/flash_swap/*.
type FlashSwapClient struct {
	*client.Client
}

// NewFlashSwapClient constructs a flash-swap REST client.
func NewFlashSwapClient(options ...client.Options) *FlashSwapClient {
	return &FlashSwapClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *FlashSwapClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
