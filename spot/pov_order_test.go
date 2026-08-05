package spot

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestSpotPOVOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// List running POV orders (private read).
	list, err := c.NewListSpotPOVOrdersService("open").SetLimit(2).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "spot/pov_orders list", err) {
			return
		}
		t.Fatalf("list pov orders: %v", err)
	}
	t.Logf("pov orders=%d", len(list))
	for _, o := range list {
		t.Logf("  id=%s pair=%s status=%s side=%s rate=%d ttl=%s",
			o.ID, o.CurrencyPair, o.Status, o.Side, o.ParticipationRate, o.TTL)
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/pov_orders",
		map[string]string{"status": "open", "limit": "2"}, true)
	testutil.AssertCovers(t, "spot/pov_orders", raw, list)

	// Query a single POV order if one already exists.
	if len(list) == 0 {
		t.Log("no open pov orders; skipping detail cover")
		return
	}
	id := list[0].ID
	one, err := c.NewGetSpotPOVOrderService(id).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "spot/pov_orders/{id}", err) {
			t.Fatalf("get pov order: %v", err)
		}
		return
	}
	t.Logf("order %s status=%s", id, one.Status)
	oneRaw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/pov_orders/"+id, nil, true)
	testutil.AssertCovers(t, "spot/pov_orders/{id}", oneRaw, one)
}
