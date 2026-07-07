package tradfi

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestTradfi(t *testing.T) {
	// ---- private read endpoints ----

	t.Run("GetMT5Account", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetMT5AccountService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/users/mt5-account", err) {
				return
			}
			t.Fatalf("mt5 account: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/users/mt5-account", nil, true)
		testutil.AssertCovers(t, "tradfi/users/mt5-account", raw, got)
	})

	t.Run("ListCategories", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListCategoriesService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/symbols/categories", err) {
				return
			}
			t.Fatalf("categories: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/symbols/categories", nil, true)
		testutil.AssertCovers(t, "tradfi/symbols/categories", raw, got)
	})

	t.Run("ListCommissions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		params := map[string]string{"symbols": "EURUSD"}
		got, err := c.NewListCommissionsService().SetSymbols("EURUSD").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/symbols/commissions", err) {
				return
			}
			t.Fatalf("commissions: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/symbols/commissions", params, true)
		testutil.AssertCovers(t, "tradfi/symbols/commissions", raw, got)
	})

	t.Run("ListSymbols", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListSymbolsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/symbols", err) {
				return
			}
			t.Fatalf("symbols: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/symbols", nil, true)
		testutil.AssertCovers(t, "tradfi/symbols", raw, got)
	})

	t.Run("ListSymbolDetails", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		params := map[string]string{"symbols": "EURUSD"}
		got, err := c.NewListSymbolDetailsService("EURUSD").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/symbols/detail", err) {
				return
			}
			t.Fatalf("symbol details: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/symbols/detail", params, true)
		testutil.AssertCovers(t, "tradfi/symbols/detail", raw, got)
	})

	t.Run("ListKlines", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		params := map[string]string{"kline_type": "1m", "limit": "5"}
		got, err := c.NewListKlinesService("EURUSD", "1m").SetLimit(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/symbols/{symbol}/klines", err) {
				return
			}
			t.Fatalf("klines: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/symbols/EURUSD/klines", params, true)
		testutil.AssertCovers(t, "tradfi/symbols/{symbol}/klines", raw, got)
	})

	t.Run("GetTicker", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetTickerService("EURUSD").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/symbols/{symbol}/tickers", err) {
				return
			}
			t.Fatalf("ticker: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/symbols/EURUSD/tickers", nil, true)
		testutil.AssertCovers(t, "tradfi/symbols/{symbol}/tickers", raw, got)
	})

	t.Run("GetAssets", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetAssetsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/users/assets", err) {
				return
			}
			t.Fatalf("assets: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/users/assets", nil, true)
		testutil.AssertCovers(t, "tradfi/users/assets", raw, got)
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
			if testutil.Tolerable(t, "tradfi/transactions", err) {
				return
			}
			t.Fatalf("transactions: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/transactions", params, true)
		testutil.AssertCovers(t, "tradfi/transactions", raw, got)
	})

	t.Run("ListOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListOrdersService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/orders", err) {
				return
			}
			t.Fatalf("orders: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/orders", nil, true)
		testutil.AssertCovers(t, "tradfi/orders", raw, got)
	})

	t.Run("ListOrderHistory", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListOrderHistoryService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/orders/history", err) {
				return
			}
			t.Fatalf("order history: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/orders/history", nil, true)
		testutil.AssertCovers(t, "tradfi/orders/history", raw, got)
	})

	t.Run("GetOrderLog", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetOrderLogService(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/orders/log/{log_id}", err) {
				return
			}
			t.Fatalf("order log: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/orders/log/1", nil, true)
		testutil.AssertCovers(t, "tradfi/orders/log/{log_id}", raw, got)
	})

	t.Run("ListPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListPositionsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/positions", err) {
				return
			}
			t.Fatalf("positions: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/positions", nil, true)
		testutil.AssertCovers(t, "tradfi/positions", raw, got)
	})

	t.Run("ListPositionHistory", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		params := map[string]string{"page_size": "5"}
		got, err := c.NewListPositionHistoryService().SetPageSize(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "tradfi/positions/history", err) {
				return
			}
			t.Fatalf("position history: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/tradfi/positions/history", params, true)
		testutil.AssertCovers(t, "tradfi/positions/history", raw, got)
	})

	// ---- private write endpoints ----
	// Gated behind GATE_TEST_WRITE and exercised with tiny/likely-rejected
	// parameters so they never place or move anything real; any error is a pass.

	t.Run("CreateUser", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateUserService().Do(cx)
		if err != nil {
			t.Logf("create user: %v (tolerable)", err)
			return
		}
		t.Log("create user accepted")
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
		_, err := c.NewCreateTransactionService("USDT", decimal.NewFromFloat(0.00000001), "withdraw").Do(cx)
		if err != nil {
			t.Logf("create transaction: %v (tolerable)", err)
			return
		}
		t.Log("create transaction accepted")
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
		_, err := c.NewCreateOrderService("EURUSD", 2, "market",
			decimal.NewFromFloat(0.00000001), decimal.NewFromFloat(0.00000001)).Do(cx)
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
		_, err := c.NewModifyOrderService(1, decimal.NewFromFloat(0.00000001)).Do(cx)
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

	t.Run("ModifyPosition", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewModifyPositionService(1).SetPriceTP(decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("modify position: %v (tolerable)", err)
			return
		}
		t.Log("modify position accepted")
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
		_, err := c.NewClosePositionService(1, 2).Do(cx)
		if err != nil {
			t.Logf("close position: %v (tolerable)", err)
			return
		}
		t.Log("close position accepted")
	})
}
