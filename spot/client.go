package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*SpotClient)(nil)

// SpotClient is the REST client for Gate's spot, margin and unified-account
// endpoints under /api/v4/spot/*. It embeds the shared core client, so every
// product client reuses the same signing/transport layer.
type SpotClient struct {
	*client.Client
}

// NewSpotClient constructs a spot REST client.
func NewSpotClient(options ...client.Options) *SpotClient {
	return &SpotClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset and stores it so that
// signed requests carry a timestamp Gate accepts. Gate rejects requests whose
// Timestamp drifts more than a minute from its own clock, so call this once at
// startup (and periodically for long-lived processes).
func (c *SpotClient) SyncServerTime(ctx context.Context) error {
	localBefore := time.Now().UnixMilli()
	resp, err := c.NewGetServerTimeService().Do(ctx)
	if err != nil {
		return err
	}
	localAfter := time.Now().UnixMilli()
	local := (localBefore + localAfter) / 2
	c.SetTimeOffset(local - resp.ServerTime.UnixMilli())
	c.GetLogger().Infof("Time sync: local=%d, server=%d, offset=%dms",
		local, resp.ServerTime.UnixMilli(), c.GetTimeOffsetMs())
	return nil
}
