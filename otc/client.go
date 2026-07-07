package otc

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*OTCClient)(nil)

// OTCClient is the REST client for Gate's OTC (fiat and stablecoin) endpoints
// under /api/v4/otc/*. It embeds the shared core client, so it reuses the same
// signing/transport layer as every other product client.
type OTCClient struct {
	*client.Client
}

// NewOTCClient constructs an OTC REST client.
func NewOTCClient(options ...client.Options) *OTCClient {
	return &OTCClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts. Gate rejects requests whose Timestamp drifts more than
// a minute from its own clock, so call this once at startup (and periodically for
// long-lived processes).
func (c *OTCClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}

// OTCAckResponse is the plain acknowledgement envelope Gate returns for OTC
// write operations that carry no data payload (order create/cancel/paid, bank
// delete/set-default and the bank-supplement submissions). code is the business
// status code, message the human-readable result and timestamp the server time
// the response was produced.
type OTCAckResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
