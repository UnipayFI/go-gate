package futures

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestFuturesPriceOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// List open price-triggered orders (read-only).
	t.Run("list", func(t *testing.T) {
		list, err := c.NewListPriceTriggeredOrdersService(SettleUSDT, OrderStatusOpen).
			SetContract("BTC_USDT").SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/price_orders", err) {
				return
			}
			t.Fatalf("list price orders: %v", err)
		}
		t.Logf("openPriceOrders=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/price_orders",
			map[string]string{"status": "open", "contract": "BTC_USDT", "limit": "2"}, true)
		testutil.AssertCovers(t, "futures/price_orders", raw, list)
	})

	// Create -> query -> cancel, plus a cancel-all pass, all reversible and gated
	// behind GATE_TEST_WRITE. The triggers are set far from the market so nothing
	// actually fires; the orders rest until we cancel them.
	t.Run("createGetCancel", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise create/get/cancel price-triggered orders")
		}

		// A buy-on-breakout auto order: rule 1 fires when price >= 1_000_000
		// (well above the last price), so it never triggers here.
		created, err := c.NewCreatePriceTriggeredOrderService(SettleUSDT, "BTC_USDT",
			decimal.NewFromInt(1000), decimal.NewFromInt(1000000), 1).
			SetSize(1).
			SetTif(TimeInForceGTC).
			Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/price_orders create", err) {
				return
			}
			t.Fatalf("create price order: %v", err)
		}
		orderID := strconv.FormatInt(created.ID, 10)
		t.Logf("created price order id=%s", orderID)

		got, err := c.NewGetPriceTriggeredOrderService(SettleUSDT, orderID).Do(cx)
		if err != nil {
			t.Fatalf("get price order %s: %v", orderID, err)
		}
		t.Logf("price order: status=%s contract=%s", got.Status, got.Initial.Contract)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/price_orders/"+orderID, nil, true)
		testutil.AssertCovers(t, "futures/price_orders/{order_id}", raw, got)

		cancelled, err := c.NewCancelPriceTriggeredOrderService(SettleUSDT, orderID).Do(cx)
		if err != nil {
			t.Fatalf("cancel price order %s: %v", orderID, err)
		}
		t.Logf("cancelled price order id=%d status=%s", cancelled.ID, cancelled.Status)

		// Place a second order and clear it via cancel-all for the contract.
		created2, err := c.NewCreatePriceTriggeredOrderService(SettleUSDT, "BTC_USDT",
			decimal.NewFromInt(1000), decimal.NewFromInt(1000000), 1).
			SetSize(1).
			SetTif(TimeInForceGTC).
			Do(cx)
		if err != nil {
			t.Fatalf("create second price order: %v", err)
		}
		t.Logf("created second price order id=%d", created2.ID)

		list, err := c.NewCancelPriceTriggeredOrderListService(SettleUSDT).
			SetContract("BTC_USDT").Do(cx)
		if err != nil {
			t.Fatalf("cancel price order list: %v", err)
		}
		t.Logf("cancel-all removed %d price order(s)", len(list))
	})
}
