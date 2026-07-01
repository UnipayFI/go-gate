package margin

import (
	"context"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*MarginClient)(nil)

// MarginClient is the REST client for Gate's margin & unified-margin endpoints under /api/v4/margin/* and /api/v4/margin/uni/*.
type MarginClient struct {
	*client.Client
}

// NewMarginClient constructs a margin & unified-margin REST client.
func NewMarginClient(options ...client.Options) *MarginClient {
	return &MarginClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *MarginClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
