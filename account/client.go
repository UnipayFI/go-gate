package account

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*AccountClient)(nil)

// AccountClient is the REST client for Gate's account endpoints under /api/v4/account/*.
type AccountClient struct {
	*client.Client
}

// NewAccountClient constructs a account REST client.
func NewAccountClient(options ...client.Options) *AccountClient {
	return &AccountClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *AccountClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
