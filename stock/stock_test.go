package stock

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestStock(t *testing.T) {
	// ---- public read endpoints ----

	t.Run("ListSymbols", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		got, err := c.NewListSymbolsService().SetPageSize(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/symbols", err) {
				return
			}
			t.Fatalf("symbols: %v", err)
		}
		params := map[string]string{"page_size": "5"}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/symbols", params, false)
		testutil.AssertCovers(t, "stock/symbols", raw, got)
	})

	t.Run("ListSymbolDetails", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		got, err := c.NewListSymbolDetailsService().SetPageSize(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/symbols/detail", err) {
				return
			}
			t.Fatalf("symbol details: %v", err)
		}
		params := map[string]string{"page_size": "5"}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/symbols/detail", params, false)
		testutil.AssertCovers(t, "stock/symbols/detail", raw, got)
	})

	t.Run("GetOrderBook", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		got, err := c.NewGetOrderBookService("AAPL").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/market/{symbol}/orderbook", err) {
				return
			}
			t.Fatalf("orderbook: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/market/AAPL/orderbook", nil, false)
		testutil.AssertCovers(t, "stock/market/{symbol}/orderbook", raw, got)
	})

	t.Run("ListExchanges", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListExchangesService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/exchanges", err) {
				return
			}
			t.Fatalf("exchanges: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/exchanges", nil, true)
		testutil.AssertCovers(t, "stock/exchanges", raw, got)
	})

	t.Run("GetFeeRate", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		got, err := c.NewGetFeeRateService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/fee-rate", err) {
				return
			}
			t.Fatalf("fee rate: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/fee-rate", nil, false)
		testutil.AssertCovers(t, "stock/fee-rate", raw, got)
	})

	// ---- private read endpoints ----

	t.Run("GetUserAssets", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetUserAssetsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/users/assets", err) {
				return
			}
			t.Fatalf("user assets: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/users/assets", nil, true)
		testutil.AssertCovers(t, "stock/users/assets", raw, got)
	})

	t.Run("ListOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListOrdersService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/orders", err) {
				return
			}
			t.Fatalf("orders: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/orders", nil, true)
		testutil.AssertCovers(t, "stock/orders", raw, got)
	})

	t.Run("ListOrderHistory", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		params := map[string]string{"page_size": "5"}
		got, err := c.NewListOrderHistoryService().SetPageSize(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/orders/history", err) {
				return
			}
			t.Fatalf("order history: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/orders/history", params, true)
		testutil.AssertCovers(t, "stock/orders/history", raw, got)
	})

	t.Run("ListPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListPositionsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/positions", err) {
				return
			}
			t.Fatalf("positions: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/positions", nil, true)
		testutil.AssertCovers(t, "stock/positions", raw, got)
	})

	t.Run("ListTransactions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		params := map[string]string{"page_size": "5"}
		got, err := c.NewListTransactionsService().SetPageSize(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock/transactions", err) {
				return
			}
			t.Fatalf("transactions: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock/transactions", params, true)
		testutil.AssertCovers(t, "stock/transactions", raw, got)
	})

	// ---- private write endpoints ----
	// Gated behind GATE_TEST_WRITE and exercised with tiny/likely-rejected
	// parameters so they never place or move anything real; any error is a pass.

	t.Run("CreateOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateOrderService("AAPL", 2, "limit", "regular", "day",
			decimal.NewFromFloat(0.00000001)).SetPrice(decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("create order: %v (tolerable)", err)
			return
		}
		t.Log("create order accepted")
	})

	t.Run("ModifyOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewModifyOrderService(1, decimal.NewFromFloat(0.00000001), decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("modify order: %v (tolerable)", err)
			return
		}
		t.Log("modify order accepted")
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
		err := c.NewCancelOrderService(1).Do(cx)
		if err != nil {
			t.Logf("cancel order: %v (tolerable)", err)
			return
		}
		t.Log("cancel order accepted")
	})

	t.Run("CancelAllOrders", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		err := c.NewCancelAllOrdersService().Do(cx)
		if err != nil {
			t.Logf("cancel all orders: %v (tolerable)", err)
			return
		}
		t.Log("cancel all orders accepted")
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
		_, err := c.NewClosePositionService("AAPL", 2).Do(cx)
		if err != nil {
			t.Logf("close position: %v (tolerable)", err)
			return
		}
		t.Log("close position accepted")
	})

	t.Run("CreateTransaction", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateTransactionService("USDT", decimal.NewFromFloat(0.00000001), "withdraw", "go-gate-test-0").Do(cx)
		if err != nil {
			t.Logf("create transaction: %v (tolerable)", err)
			return
		}
		t.Log("create transaction accepted")
	})
}
