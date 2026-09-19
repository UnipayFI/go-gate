package announcement

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*AnnouncementClient)(nil)

// AnnouncementClient is the REST client for Gate's announcement endpoints under
// /api/v4/ann/*.
type AnnouncementClient struct {
	*client.Client
}

// NewAnnouncementClient constructs an announcement REST client.
func NewAnnouncementClient(options ...client.Options) *AnnouncementClient {
	return &AnnouncementClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *AnnouncementClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
