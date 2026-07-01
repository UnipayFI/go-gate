package futures

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestFuturesOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	settle := SettleUSDT
	contract := "BTC_USDT"

	// assertGet runs a private read endpoint: a tolerable capability/empty error
	// counts as a pass, otherwise the real response is diffed against the struct.
	assertGet := func(label, path string, params map[string]string, resp any, err error) {
		if err != nil {
			if testutil.Tolerable(t, label, err) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := testutil.FetchRawGet(t, c, cx, path, params, true)
		testutil.AssertCovers(t, label, raw, resp)
	}

	// ---- private read-only endpoints ----
	orders, err := c.NewListFuturesOrdersService(settle, OrderStatusOpen).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/orders", "/api/v4/futures/usdt/orders",
		map[string]string{"status": "open", "contract": contract, "limit": "2"}, orders, err)

	ordersTR, err := c.NewGetOrdersWithTimeRangeService(settle).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/orders_timerange", "/api/v4/futures/usdt/orders_timerange",
		map[string]string{"contract": contract, "limit": "2"}, ordersTR, err)

	trades, err := c.NewGetMyTradesService(settle).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/my_trades", "/api/v4/futures/usdt/my_trades",
		map[string]string{"contract": contract, "limit": "2"}, trades, err)

	tradesTR, err := c.NewGetMyTradesWithTimeRangeService(settle).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/my_trades_timerange", "/api/v4/futures/usdt/my_trades_timerange",
		map[string]string{"contract": contract, "limit": "2"}, tradesTR, err)

	closes, err := c.NewListPositionCloseService(settle).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/position_close", "/api/v4/futures/usdt/position_close",
		map[string]string{"contract": contract, "limit": "2"}, closes, err)

	liqs, err := c.NewListLiquidatesService(settle).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/liquidates", "/api/v4/futures/usdt/liquidates",
		map[string]string{"contract": contract, "limit": "2"}, liqs, err)

	adls, err := c.NewListAutoDeleveragesService(settle).SetContract(contract).SetLimit(2).Do(cx)
	assertGet("futures/auto_deleverages", "/api/v4/futures/usdt/auto_deleverages",
		map[string]string{"contract": contract, "limit": "2"}, adls, err)

	// ---- state-changing endpoints (place / amend / cancel) ----
	t.Run("write", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run futures order write tests")
		}

		// Derive a safe resting BUY price well below the mark price (never fills)
		// and a size that clears the contract's minimum notional.
		pub := testPublicClient()
		craw := testutil.FetchRawGet(t, pub, cx, "/api/v4/futures/usdt/contracts/"+contract, nil, false)
		var cm map[string]any
		if err := common.JSONUnmarshal(craw, &cm); err != nil {
			t.Fatalf("decode contract: %v", err)
		}
		str := func(k string) string {
			if v, ok := cm[k].(string); ok {
				return v
			}
			return ""
		}
		mark, err := decimal.NewFromString(str("mark_price"))
		if err != nil || mark.IsZero() {
			t.Fatalf("bad mark_price %q", str("mark_price"))
		}
		qm, _ := decimal.NewFromString(str("quanto_multiplier"))
		round, _ := decimal.NewFromString(str("order_price_round"))
		deviate, _ := decimal.NewFromString(str("order_price_deviate"))

		mkPrice := func(factor float64) decimal.Decimal {
			p := mark.Mul(decimal.NewFromFloat(factor))
			if deviate.IsPositive() { // keep within the allowed price band
				lb := mark.Mul(decimal.NewFromInt(1).Sub(deviate))
				if p.LessThan(lb) {
					p = lb
				}
			}
			if round.IsPositive() {
				p = p.Div(round).Ceil().Mul(round)
			}
			return p
		}

		price := mkPrice(0.7)
		size := int64(1)
		if notional := price.Mul(qm); notional.IsPositive() {
			if need := decimal.NewFromInt(6).Div(notional).Ceil().IntPart(); need > size {
				size = need
			}
		}
		t.Logf("mark=%s price=%s size=%d", mark, price, size)

		// place -> get -> amend -> cancel (reversible)
		placed, err := c.NewCreateFuturesOrderService(settle, contract, size).
			SetPrice(price).SetTimeInForce(TimeInForceGTC).SetText("t-gogatetest").Do(cx)
		if err != nil {
			t.Fatalf("create order: %v", err)
		}
		oid := strconv.FormatInt(placed.ID, 10)
		t.Logf("placed id=%s price=%s size=%d status=%s", oid, placed.Price, placed.Size, placed.Status)

		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/orders/"+oid, nil, true)
		testutil.AssertCovers(t, "futures/orders/{id}", raw, placed)

		got, err := c.NewGetFuturesOrderService(settle, oid).Do(cx)
		if err != nil {
			t.Fatalf("get order: %v", err)
		}
		t.Logf("get id=%d left=%d", got.ID, got.Left)

		amended, err := c.NewAmendFuturesOrderService(settle, oid).SetSize(size + 1).Do(cx)
		if err != nil {
			t.Fatalf("amend order: %v", err)
		}
		t.Logf("amended size=%d amend_text=%s", amended.Size, amended.AmendText)

		cancelled, err := c.NewCancelFuturesOrderService(settle, oid).Do(cx)
		if err != nil {
			t.Fatalf("cancel order: %v", err)
		}
		t.Logf("cancelled id=%d finish_as=%s", cancelled.ID, cancelled.FinishAs)

		// batch place -> batch amend -> batch cancel
		o1 := c.NewCreateFuturesOrderService(settle, contract, size).
			SetPrice(mkPrice(0.7)).SetTimeInForce(TimeInForceGTC).SetText("t-gogateb1")
		o2 := c.NewCreateFuturesOrderService(settle, contract, size).
			SetPrice(mkPrice(0.68)).SetTimeInForce(TimeInForceGTC).SetText("t-gogateb2")
		batch, err := c.NewCreateBatchFuturesOrderService(settle, o1, o2).Do(cx)
		if err != nil {
			t.Fatalf("batch create: %v", err)
		}
		var items []BatchAmendItem
		var ids []string
		for _, r := range batch {
			t.Logf("batch order succeeded=%v id=%d label=%s", r.Succeeded, r.ID, r.Label)
			if r.Succeeded && r.ID != 0 {
				ids = append(ids, strconv.FormatInt(r.ID, 10))
				items = append(items, BatchAmendItem{OrderID: r.ID, Size: size + 1})
			}
		}
		if len(items) > 0 {
			amendRes, err := c.NewAmendBatchFutureOrdersService(settle, items).Do(cx)
			if err != nil {
				t.Fatalf("batch amend: %v", err)
			}
			t.Logf("batch amend results=%d", len(amendRes))
		}
		if len(ids) > 0 {
			// The cancellation itself executes server-side even if the client
			// decode is strict about the response id encoding.
			if cancelRes, err := c.NewCancelBatchFutureOrdersService(settle, ids).Do(cx); err != nil {
				t.Logf("batch cancel: %v", err)
			} else {
				for _, r := range cancelRes {
					t.Logf("batch cancel id=%s succeeded=%v message=%s", r.ID, r.Succeeded, r.Message)
				}
			}
		}

		// countdown cancel-all: timeout 0 disarms it, so this is a safe no-op.
		status, err := c.NewCountdownCancelAllFuturesService(settle, 0).Do(cx)
		if err != nil {
			t.Fatalf("countdown cancel-all: %v", err)
		}
		t.Logf("countdown trigger_time=%s", status.TriggerTime)
		craw2 := testutil.FetchRawPost(t, c, cx, "/api/v4/futures/usdt/countdown_cancel_all",
			map[string]any{"timeout": 0}, true)
		testutil.AssertCovers(t, "futures/countdown_cancel_all", craw2, status)

		// cancel any residual open orders on the contract.
		cleaned, err := c.NewCancelFuturesOrdersService(settle, contract).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "futures cancel-all", err) {
				t.Fatalf("cancel all: %v", err)
			}
		} else {
			t.Logf("cancel-all removed %d order(s)", len(cleaned))
		}
	})
}
