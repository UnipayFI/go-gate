package subaccount

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*SubAccountClient)(nil)

// SubAccountClient is the REST client for Gate's sub-account endpoints under /api/v4/sub_accounts/*.
type SubAccountClient struct {
	*client.Client
}

// NewSubAccountClient constructs a sub-account REST client.
func NewSubAccountClient(options ...client.Options) *SubAccountClient {
	return &SubAccountClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *SubAccountClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
