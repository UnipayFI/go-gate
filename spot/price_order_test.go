package spot

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestSpotPriceOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// List running auto orders (private read).
	list, err := c.NewListSpotPriceTriggeredOrdersService("open").SetLimit(2).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "spot/price_orders list", err) {
			return
		}
		t.Fatalf("list price orders: %v", err)
	}
	t.Logf("price orders=%d", len(list))
	for _, o := range list {
		t.Logf("  id=%d market=%s status=%s side=%s", o.ID, o.Market, o.Status, o.Put.Side)
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/price_orders",
		map[string]string{"status": "open", "limit": "2"}, true)
	testutil.AssertCovers(t, "spot/price_orders", raw, list)

	// Query a single auto order if one already exists.
	if len(list) > 0 {
		id := strconv.FormatInt(list[0].ID, 10)
		one, err := c.NewGetSpotPriceTriggeredOrderService(id).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "spot/price_orders/{id}", err) {
				t.Fatalf("get price order: %v", err)
			}
		} else {
			t.Logf("order %s status=%s", id, one.Status)
			oneRaw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/price_orders/"+id, nil, true)
			testutil.AssertCovers(t, "spot/price_orders/{id}", oneRaw, one)
		}
	}

	// State-changing: place -> query -> cancel a tiny, far-from-market auto order.
	if !testutil.WriteEnabled() {
		t.Skip("set GATE_TEST_WRITE=1 to exercise create/cancel price orders")
	}

	// Trigger to buy only if BTC falls to 1000 (far below market), so the order
	// stays queued and never fires; the queued limit buy is likewise far off.
	created, err := c.NewCreateSpotPriceTriggeredOrderService(
		"BTC_USDT",
		decimal.NewFromInt(1000), "<=", 3600,
		SideBuy, decimal.NewFromInt(1000), decimal.NewFromFloat(0.006), Account("normal"),
	).SetPutType(OrderTypeLimit).SetPutTimeInForce(TimeInForceGTC).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "spot/price_orders create", err) {
			return
		}
		t.Fatalf("create price order: %v", err)
	}
	t.Logf("created price order id=%d", created.ID)

	id := strconv.FormatInt(created.ID, 10)
	got, err := c.NewGetSpotPriceTriggeredOrderService(id).Do(cx)
	if err != nil {
		t.Fatalf("get created price order: %v", err)
	}
	t.Logf("queried order status=%s", got.Status)

	cancelled, err := c.NewCancelSpotPriceTriggeredOrderService(id).Do(cx)
	if err != nil {
		t.Fatalf("cancel price order: %v", err)
	}
	t.Logf("cancelled order status=%s", cancelled.Status)
}
