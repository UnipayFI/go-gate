package delivery

import (
	"context"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/request"
)

var _ request.Client = (*DeliveryClient)(nil)

// DeliveryClient is the REST client for Gate's delivery (dated-futures)
// endpoints under /api/v4/delivery/{settle}/*. Delivery settles in USDT.
type DeliveryClient struct {
	*client.Client
}

// NewDeliveryClient constructs a delivery REST client.
func NewDeliveryClient(options ...client.Options) *DeliveryClient {
	return &DeliveryClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *DeliveryClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}
