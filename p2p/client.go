package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/request"
)

var _ request.Client = (*P2PClient)(nil)

// P2PClient is the REST client for Gate's P2P merchant endpoints under
// /api/v4/p2p/merchant/*. It embeds the shared core client, so it reuses the
// same signing/transport layer as every other product client.
type P2PClient struct {
	*client.Client
}

// NewP2PClient constructs a P2P REST client.
func NewP2PClient(options ...client.Options) *P2PClient {
	return &P2PClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset (via GET /api/v4/spot/time,
// Gate's single server-time source) and stores it so signed requests carry a
// Timestamp Gate accepts.
func (c *P2PClient) SyncServerTime(ctx context.Context) error {
	offset, server, err := request.FetchServerTimeOffset(ctx, c)
	if err != nil {
		return err
	}
	c.SetTimeOffset(offset)
	c.GetLogger().Infof("Time sync: server=%d, offset=%dms", server, c.GetTimeOffsetMs())
	return nil
}

// P2PResponse is the common envelope every /api/v4/p2p/merchant/* endpoint wraps
// its payload in. Unlike Gate's core v4 API (which returns the payload directly),
// the P2P subsystem always answers with this envelope: Code 0 means success and
// the typed Data holds the endpoint-specific payload.
type P2PResponse[T any] struct {
	Timestamp time.Time `json:"timestamp,format:unix"`
	Method    string    `json:"method"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Data      T         `json:"data"`
	Version   string    `json:"version"`
}

// P2PEmptyData is the empty object Gate returns in the Data field of P2P actions
// (confirm payment/receipt, cancel) that carry no payload beyond the envelope.
type P2PEmptyData struct{}

// doP2P executes a P2P request and decodes its business envelope. Unlike Gate's
// core v4 API, the P2P subsystem answers with an HTTP 200 even on business
// errors, signalling failure via a non-zero envelope Code (and returning an
// empty {} Data object rather than the success shape). doP2P inspects Code first
// and surfaces a non-zero Code as a *client.APIError (Label = Message, e.g.
// "NO_ACCESS") so callers get a typed error instead of a Data-decode failure.
func doP2P[T any](req *request.Request) (*P2PResponse[T], error) {
	raw, err := request.DoRaw(req)
	if err != nil {
		return nil, err
	}
	var probe struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := common.JSONUnmarshal(raw, &probe); err != nil {
		return nil, fmt.Errorf("p2p: decode envelope: %w (body: %s)", err, common.BytesToString(raw))
	}
	if probe.Code != 0 {
		return nil, &client.APIError{Label: probe.Message, Message: probe.Message}
	}
	var out P2PResponse[T]
	if err := common.JSONUnmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("p2p: decode response: %w (body: %s)", err, common.BytesToString(raw))
	}
	return &out, nil
}
