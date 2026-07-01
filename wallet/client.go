package wallet

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*WalletClient)(nil)

// WalletClient is the REST client for Gate's wallet & withdrawal endpoints under /api/v4/wallet/* and /api/v4/withdrawals/*.
type WalletClient struct {
	*client.Client
}

// NewWalletClient constructs a wallet & withdrawal REST client.
func NewWalletClient(options ...client.Options) *WalletClient {
	return &WalletClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *WalletClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
