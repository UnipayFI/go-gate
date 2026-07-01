package request

import (
	"context"
	"time"
)

// ServerTime is Gate's exchange clock (GET /api/v4/spot/time), a millisecond
// epoch number. It is the single server-time source for every product; the
// futures/delivery/options gateways have no separate time endpoint.
type ServerTime struct {
	ServerTime time.Time `json:"server_time,format:unixmilli"`
}

// FetchServerTimeOffset measures the client/server clock offset in milliseconds
// via GET /api/v4/spot/time, midpointing the local clock around the round trip.
// Each product client's SyncServerTime wraps this and stores the offset so signed
// requests carry a Timestamp Gate accepts.
func FetchServerTimeOffset(ctx context.Context, c Client) (offsetMs, serverMs int64, err error) {
	localBefore := time.Now().UnixMilli()
	resp, err := Do[ServerTime](Get(ctx, c, "/api/v4/spot/time"))
	if err != nil {
		return 0, 0, err
	}
	localAfter := time.Now().UnixMilli()
	local := (localBefore + localAfter) / 2
	serverMs = resp.ServerTime.UnixMilli()
	return local - serverMs, serverMs, nil
}
