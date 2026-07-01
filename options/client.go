package options

import (
	"context"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*OptionsClient)(nil)

// OptionsClient is the REST client for Gate's options endpoints under /api/v4/options/*.
type OptionsClient struct {
	*client.Client
}

// NewOptionsClient constructs a options REST client.
func NewOptionsClient(options ...client.Options) *OptionsClient {
	return &OptionsClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *OptionsClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
