package rebate

import (
	"context"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*RebateClient)(nil)

// RebateClient is the REST client for Gate's rebate endpoints under /api/v4/rebate/*.
type RebateClient struct {
	*client.Client
}

// NewRebateClient constructs a rebate REST client.
func NewRebateClient(options ...client.Options) *RebateClient {
	return &RebateClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *RebateClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
