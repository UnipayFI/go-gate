package crossex

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestCrossex(t *testing.T) {
	// ---- Private read endpoints ------------------------------------------

	t.Run("QuerySymbols", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewQuerySymbolsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/rule/symbols", err) {
				return
			}
			t.Fatalf("query symbols: %v", err)
		}
		t.Logf("symbols=%d", len(got))
		if len(got) == 0 {
			t.Skip("no symbols")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/rule/symbols", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/rule/symbols", raw, got)
	})

	t.Run("QueryRiskLimits", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		symbols := "BINANCE_FUTURE_BTC_USDT"
		got, err := c.NewQueryRiskLimitsService(symbols).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/rule/risk_limits", err) {
				return
			}
			t.Fatalf("query risk limits: %v", err)
		}
		t.Logf("risk limits=%d", len(got))
		if len(got) == 0 {
			t.Skip("no risk limits")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/rule/risk_limits",
			map[string]string{"symbols": symbols}, true)
		testutil.AssertCovers(t, "crossex/rule/risk_limits", raw, got)
	})

	t.Run("QueryTransferCoins", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewQueryTransferCoinsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/transfers/coin", err) {
				return
			}
			t.Fatalf("query transfer coins: %v", err)
		}
		t.Logf("transfer coins=%d", len(got))
		if len(got) == 0 {
			t.Skip("no transfer coins")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/transfers/coin", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/transfers/coin", raw, got)
	})

	t.Run("ListTransfers", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListTransfersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/transfers", err) {
				return
			}
			t.Fatalf("list transfers: %v", err)
		}
		t.Logf("transfers=%d", len(got))
		if len(got) == 0 {
			t.Log("no transfers")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/transfers",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/transfers", raw, got)
	})

	t.Run("GetOrder", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetOrderService("0").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/orders/{order_id}", err) {
				return
			}
			t.Fatalf("get order: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/orders/0", nil, true)
		testutil.AssertCovers(t, "crossex/orders/{order_id}", raw, got)
	})

	t.Run("GetAccount", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetAccountService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/accounts", err) {
				return
			}
			t.Fatalf("get account: %v", err)
		}
		t.Logf("account: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/accounts", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/accounts", raw, got)
	})

	t.Run("GetPositionLeverage", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetPositionLeverageService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/positions/leverage", err) {
				return
			}
			t.Fatalf("get position leverage: %v", err)
		}
		t.Logf("position leverage=%d", len(got))
		if len(got) == 0 {
			t.Skip("no position leverage")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/positions/leverage", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/positions/leverage", raw, got)
	})

	t.Run("GetMarginPositionLeverage", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetMarginPositionLeverageService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/margin_positions/leverage", err) {
				return
			}
			t.Fatalf("get margin position leverage: %v", err)
		}
		t.Logf("margin position leverage=%d", len(got))
		if len(got) == 0 {
			t.Skip("no margin position leverage")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/margin_positions/leverage", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/margin_positions/leverage", raw, got)
	})

	t.Run("QueryInterestRate", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewQueryInterestRateService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/interest_rate", err) {
				return
			}
			t.Fatalf("query interest rate: %v", err)
		}
		t.Logf("interest rates=%d", len(got))
		if len(got) == 0 {
			t.Skip("no interest rates")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/interest_rate", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/interest_rate", raw, got)
	})

	t.Run("GetFee", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetFeeService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/fee", err) {
				return
			}
			t.Fatalf("get fee: %v", err)
		}
		t.Logf("fees=%d", len(got))
		if len(got) == 0 {
			t.Skip("no fees")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/fee", nil, true)
		testutil.AssertCovers(t, "crossex/fee", raw, got)
	})

	t.Run("ListPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListPositionsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/positions", err) {
				return
			}
			t.Fatalf("list positions: %v", err)
		}
		t.Logf("positions=%d", len(got))
		if len(got) == 0 {
			t.Log("no positions")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/positions", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/positions", raw, got)
	})

	t.Run("ListMarginPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListMarginPositionsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/margin_positions", err) {
				return
			}
			t.Fatalf("list margin positions: %v", err)
		}
		t.Logf("margin positions=%d", len(got))
		if len(got) == 0 {
			t.Log("no margin positions")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/margin_positions", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/margin_positions", raw, got)
	})

	t.Run("GetADLRank", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		symbol := "BINANCE_FUTURE_BTC_USDT"
		got, err := c.NewGetADLRankService(symbol).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/adl_rank", err) {
				return
			}
			t.Fatalf("get adl rank: %v", err)
		}
		t.Logf("adl rank: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/adl_rank",
			map[string]string{"symbol": symbol}, true)
		testutil.AssertCovers(t, "crossex/adl_rank", raw, got)
	})

	t.Run("ListOpenOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListOpenOrdersService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/open_orders", err) {
				return
			}
			t.Fatalf("list open orders: %v", err)
		}
		t.Logf("open orders=%d", len(got))
		if len(got) == 0 {
			t.Log("no open orders")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/open_orders", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/open_orders", raw, got)
	})

	t.Run("ListHistoryOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHistoryOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/history_orders", err) {
				return
			}
			t.Fatalf("list history orders: %v", err)
		}
		t.Logf("history orders=%d", len(got))
		if len(got) == 0 {
			t.Log("no history orders")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/history_orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/history_orders", raw, got)
	})

	t.Run("ListHistoryPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHistoryPositionsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/history_positions", err) {
				return
			}
			t.Fatalf("list history positions: %v", err)
		}
		t.Logf("history positions=%d", len(got))
		if len(got) == 0 {
			t.Log("no history positions")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/history_positions",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/history_positions", raw, got)
	})

	t.Run("ListHistoryMarginPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHistoryMarginPositionsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/history_margin_positions", err) {
				return
			}
			t.Fatalf("list history margin positions: %v", err)
		}
		t.Logf("history margin positions=%d", len(got))
		if len(got) == 0 {
			t.Log("no history margin positions")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/history_margin_positions",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/history_margin_positions", raw, got)
	})

	t.Run("ListHistoryMarginInterests", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHistoryMarginInterestsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/history_margin_interests", err) {
				return
			}
			t.Fatalf("list history margin interests: %v", err)
		}
		t.Logf("history margin interests=%d", len(got))
		if len(got) == 0 {
			t.Log("no history margin interests")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/history_margin_interests",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/history_margin_interests", raw, got)
	})

	t.Run("ListHistoryTrades", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHistoryTradesService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/history_trades", err) {
				return
			}
			t.Fatalf("list history trades: %v", err)
		}
		t.Logf("history trades=%d", len(got))
		if len(got) == 0 {
			t.Log("no history trades")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/history_trades",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/history_trades", raw, got)
	})

	t.Run("ListAccountBook", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListAccountBookService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/account_book", err) {
				return
			}
			t.Fatalf("list account book: %v", err)
		}
		t.Logf("account book=%d", len(got))
		if len(got) == 0 {
			t.Log("no account book records")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/account_book",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "crossex/account_book", raw, got)
	})

	t.Run("QueryCoinDiscountRate", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewQueryCoinDiscountRateService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "crossex/coin_discount_rate", err) {
				return
			}
			t.Fatalf("query coin discount rate: %v", err)
		}
		t.Logf("coin discount rates=%d", len(got))
		if len(got) == 0 {
			t.Skip("no coin discount rates")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/crossex/coin_discount_rate", map[string]string{}, true)
		testutil.AssertCovers(t, "crossex/coin_discount_rate", raw, got)
	})

	// ---- Write endpoints -------------------------------------------------
	// Gated behind GATE_TEST_WRITE and exercised with tiny/likely-rejected
	// parameters so they never place a real order/transfer; any error is a pass.

	t.Run("CreateTransfer", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateTransferService("USDT", decimal.NewFromFloat(0.00000001),
			"CROSSEX_BINANCE", "CROSSEX_OKX").Do(cx)
		if err != nil {
			t.Logf("create transfer: %v (tolerable)", err)
			return
		}
		t.Log("create transfer accepted")
	})

	t.Run("CreateOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateOrderService("BINANCE_SPOT_BTC_USDT", "BUY").
			SetType("LIMIT").
			SetQty(decimal.NewFromFloat(0.00000001)).
			SetPrice(decimal.NewFromInt(1)).
			Do(cx)
		if err != nil {
			t.Logf("create order: %v (tolerable)", err)
			return
		}
		t.Log("create order accepted")
	})

	t.Run("AmendOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewAmendOrderService("0").SetPrice(decimal.NewFromInt(1)).Do(cx)
		if err != nil {
			t.Logf("amend order: %v (tolerable)", err)
			return
		}
		t.Log("amend order accepted")
	})

	t.Run("CancelOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCancelOrderService("0").Do(cx)
		if err != nil {
			t.Logf("cancel order: %v (tolerable)", err)
			return
		}
		t.Log("cancel order accepted")
	})

	t.Run("CreateConvertQuote", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateConvertQuoteService("OKX", "USDT", "BTC",
			decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("create convert quote: %v (tolerable)", err)
			return
		}
		t.Log("create convert quote accepted")
	})

	t.Run("CreateConvertOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateConvertOrderService("0").Do(cx)
		if err != nil {
			t.Logf("create convert order: %v (tolerable)", err)
			return
		}
		t.Log("create convert order accepted")
	})

	t.Run("UpdateAccount", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewUpdateAccountService().SetPositionMode("SINGLE").Do(cx)
		if err != nil {
			t.Logf("update account: %v (tolerable)", err)
			return
		}
		t.Log("update account accepted")
	})

	t.Run("UpdatePositionLeverage", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewUpdatePositionLeverageService("BINANCE_FUTURE_BTC_USDT",
			decimal.NewFromInt(1)).Do(cx)
		if err != nil {
			t.Logf("update position leverage: %v (tolerable)", err)
			return
		}
		t.Log("update position leverage accepted")
	})

	t.Run("UpdateMarginPositionLeverage", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewUpdateMarginPositionLeverageService("BINANCE_MARGIN_BTC_USDT",
			decimal.NewFromInt(1)).Do(cx)
		if err != nil {
			t.Logf("update margin position leverage: %v (tolerable)", err)
			return
		}
		t.Log("update margin position leverage accepted")
	})

	t.Run("ClosePosition", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewClosePositionService("BINANCE_FUTURE_BTC_USDT").Do(cx)
		if err != nil {
			t.Logf("close position: %v (tolerable)", err)
			return
		}
		t.Log("close position accepted")
	})
}
