package spot

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestSpotOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// --- private reads (always run) ---

	// ListOrders (open) on a single pair.
	orders, err := c.NewListOrdersService("BTC_USDT", "open").SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "spot/orders", err) {
			t.Fatalf("list orders: %v", err)
		}
	} else {
		t.Logf("open orders on BTC_USDT: %d", len(orders))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/orders",
			map[string]string{"currency_pair": "BTC_USDT", "status": "open", "limit": "2"}, true)
		testutil.AssertCovers(t, "spot/orders", raw, orders)
	}

	// ListAllOpenOrders across every pair with pending orders.
	openOrders, err := c.NewListAllOpenOrdersService().SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "spot/open_orders", err) {
			t.Fatalf("list all open orders: %v", err)
		}
	} else {
		t.Logf("pairs with open orders: %d", len(openOrders))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/open_orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "spot/open_orders", raw, openOrders)
	}

	// --- state-changing (opt-in) ---
	if !testutil.WriteEnabled() {
		t.Skip("set GATE_TEST_WRITE=1 to run spot order place/amend/cancel tests")
	}

	// Derive a safe limit price ~30% below market and an amount whose notional
	// clears the pair's minimum, using the pair's own precisions.
	pair, err := c.NewGetCurrencyPairService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("currency pair: %v", err)
	}
	tickers, err := c.NewGetTickersService().SetCurrencyPair("BTC_USDT").Do(cx)
	if err != nil || len(tickers) == 0 {
		t.Fatalf("ticker: %v", err)
	}
	last := tickers[0].Last
	price := last.Mul(decimal.NewFromFloat(0.7)).Round(int32(pair.Precision))
	notional := decimal.Max(pair.MinQuoteAmount, decimal.NewFromInt(6)).Mul(decimal.NewFromFloat(1.3))
	amount := notional.Div(price).Round(int32(pair.AmountPrecision))
	if !price.IsPositive() || !amount.IsPositive() {
		t.Fatalf("computed non-positive price=%s amount=%s", price, amount)
	}
	t.Logf("using price=%s amount=%s (last=%s)", price, amount, last)

	// CreateOrder: limit BUY well below market so it rests open, then GetOrder,
	// amend the price, and cancel — fully reversible.
	placed, err := c.NewCreateOrderService("BTC_USDT", SideBuy, amount).
		SetType(OrderTypeLimit).
		SetPrice(price).
		SetTimeInForce(TimeInForceGTC).
		SetText("t-gogate-order").
		Do(cx)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	t.Logf("placed id=%s status=%s left=%s", placed.ID, placed.Status, placed.Left)

	// GetOrder + coverage check on the real order shape.
	got, err := c.NewGetOrderService(placed.ID, "BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("get order: %v", err)
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/orders/"+placed.ID,
		map[string]string{"currency_pair": "BTC_USDT"}, true)
	testutil.AssertCovers(t, "spot/orders/{order_id}", raw, got)

	// AmendOrder: nudge the price lower (still below market).
	amendPrice := price.Mul(decimal.NewFromFloat(0.99)).Round(int32(pair.Precision))
	amended, err := c.NewAmendOrderService(placed.ID).
		SetCurrencyPair("BTC_USDT").
		SetPrice(amendPrice).
		SetAmendText("t-gogate-amend").
		Do(cx)
	if err != nil {
		t.Fatalf("amend order: %v", err)
	}
	t.Logf("amended id=%s price=%s", amended.ID, amended.Price)

	// CancelOrder.
	cancelled, err := c.NewCancelOrderService(placed.ID, "BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("cancel order: %v", err)
	}
	t.Logf("cancelled id=%s status=%s", cancelled.ID, cancelled.Status)

	// CreateBatchOrders: two resting limit buys.
	o1 := c.NewCreateOrderService("BTC_USDT", SideBuy, amount).
		SetType(OrderTypeLimit).SetPrice(price).SetText("t-gogate-batch1")
	o2 := c.NewCreateOrderService("BTC_USDT", SideBuy, amount).
		SetType(OrderTypeLimit).
		SetPrice(price.Mul(decimal.NewFromFloat(0.98)).Round(int32(pair.Precision))).
		SetText("t-gogate-batch2")
	batch, err := c.NewCreateBatchOrdersService(o1, o2).Do(cx)
	if err != nil {
		t.Fatalf("create batch orders: %v", err)
	}
	var toCancel []CancelOrderReq
	for _, b := range batch {
		t.Logf("batch order id=%s succeeded=%v label=%s message=%s", b.ID, b.Succeeded, b.Label, b.Message)
		if b.Succeeded && b.ID != "" {
			toCancel = append(toCancel, CancelOrderReq{CurrencyPair: "BTC_USDT", ID: b.ID})
		}
	}

	// AmendBatchOrders on the first batch order (if it landed).
	if len(toCancel) > 0 {
		item := NewBatchAmendItem(toCancel[0].ID, "BTC_USDT").
			SetPrice(price.Mul(decimal.NewFromFloat(0.97)).Round(int32(pair.Precision)))
		amendBatch, err := c.NewAmendBatchOrdersService(item).Do(cx)
		if err != nil {
			t.Fatalf("amend batch orders: %v", err)
		}
		for _, ab := range amendBatch {
			t.Logf("amend-batch id=%s succeeded=%v label=%s", ab.ID, ab.Succeeded, ab.Label)
		}
	}

	// CancelBatchOrders by ID.
	if len(toCancel) > 0 {
		results, err := c.NewCancelBatchOrdersService(toCancel...).Do(cx)
		if err != nil {
			t.Fatalf("cancel batch orders: %v", err)
		}
		for _, r := range results {
			t.Logf("cancel-batch id=%s succeeded=%v label=%s", r.ID, r.Succeeded, r.Label)
		}
	}

	// CancelOrders: sweep any leftover open BTC_USDT orders from this test.
	swept, err := c.NewCancelOrdersService().SetCurrencyPair("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("cancel all orders: %v", err)
	}
	t.Logf("cancel-all swept %d order(s)", len(swept))

	// CountdownCancelAllSpot: disarm (timeout 0) — safe no-op that still exercises
	// the endpoint and returns the trigger time.
	status, err := c.NewCountdownCancelAllSpotService(0).SetCurrencyPair("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("countdown cancel all: %v", err)
	}
	t.Logf("countdown triggerTime=%s", status.TriggerTime)

	// CreateCrossLiquidateOrder: only valid for cross-margin accounts holding a
	// disabled currency, so it is expected to error for normal accounts. Exercise
	// the path and cancel it defensively if it somehow rests.
	liq, err := c.NewCreateCrossLiquidateOrderService("BTC_USDT", amount, price).
		SetText("t-gogate-liq").
		Do(cx)
	if err != nil {
		t.Logf("cross liquidate order (expected to fail without a disabled cross-margin currency): %v", err)
	} else {
		t.Logf("cross liquidate order id=%s status=%s", liq.ID, liq.Status)
		if liq.ID != "" {
			if _, cerr := c.NewCancelOrderService(liq.ID, "BTC_USDT").SetAccount(AccountCrossMargin).Do(cx); cerr != nil {
				t.Logf("cleanup cancel cross-liquidate order: %v", cerr)
			}
		}
	}
}
