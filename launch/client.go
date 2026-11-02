package launch

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*LaunchClient)(nil)

// LaunchClient is the REST client for Gate's Launch endpoints under
// /api/v4/launch/* (LaunchPool staking mining, HODLer Airdrop and CandyDrop).
// It embeds the shared core client, so it reuses the same signing/transport
// layer as every other product client.
type LaunchClient struct {
	*client.Client
}

// NewLaunchClient constructs a Launch REST client.
func NewLaunchClient(options ...client.Options) *LaunchClient {
	return &LaunchClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *LaunchClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
