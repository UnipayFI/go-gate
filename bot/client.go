package bot

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*BotClient)(nil)

// BotClient is the REST client for Gate's AIHub trading-bot endpoints under
// /api/v4/bot/*. It embeds the shared core client, so it reuses the same
// signing/transport layer as every other product client.
type BotClient struct {
	*client.Client
}

// NewBotClient constructs a bot REST client.
func NewBotClient(options ...client.Options) *BotClient {
	return &BotClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *BotClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
