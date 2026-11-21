package stockoptions

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestStockOptions(t *testing.T) {
	// ---- public read endpoints ----

	t.Run("GetFeeRate", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		got, err := c.NewGetFeeRateService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/fee-rate", err) {
				return
			}
			t.Fatalf("fee-rate: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/fee-rate", nil, false)
		testutil.AssertCovers(t, "stock_options/fee-rate", raw, got)
	})

	t.Run("ListUnderlyings", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		got, err := c.NewListUnderlyingsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/underlyings", err) {
				return
			}
			t.Fatalf("underlyings: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/underlyings", nil, false)
		testutil.AssertCovers(t, "stock_options/underlyings", raw, got)
	})

	t.Run("ListContracts", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		underlyings, err := c.NewListUnderlyingsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/underlyings", err) {
				return
			}
			t.Fatalf("underlyings: %v", err)
		}
		if len(underlyings.Data.List) == 0 {
			t.Skip("no stock option underlyings returned")
		}
		u := underlyings.Data.List[0].Underlying
		got, err := c.NewListContractsService(u).SetPageSize(5).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/contracts", err) {
				return
			}
			t.Fatalf("contracts: %v", err)
		}
		params := map[string]string{"underlying": u, "page_size": "5"}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/contracts", params, false)
		testutil.AssertCovers(t, "stock_options/contracts", raw, got)
	})

	t.Run("ListExpirations", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		underlyings, err := c.NewListUnderlyingsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/underlyings", err) {
				return
			}
			t.Fatalf("underlyings: %v", err)
		}
		if len(underlyings.Data.List) == 0 {
			t.Skip("no stock option underlyings returned")
		}
		u := underlyings.Data.List[0].Underlying
		got, err := c.NewListExpirationsService(u).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/market/expirations", err) {
				return
			}
			t.Fatalf("expirations: %v", err)
		}
		params := map[string]string{"underlying": u}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/market/expirations", params, false)
		testutil.AssertCovers(t, "stock_options/market/expirations", raw, got)
	})

	t.Run("GetOptionChainAndOrderBook", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		underlyings, err := c.NewListUnderlyingsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/underlyings", err) {
				return
			}
			t.Fatalf("underlyings: %v", err)
		}
		if len(underlyings.Data.List) == 0 {
			t.Skip("no stock option underlyings returned")
		}
		u := underlyings.Data.List[0].Underlying
		exps, err := c.NewListExpirationsService(u).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/market/expirations", err) {
				return
			}
			t.Fatalf("expirations: %v", err)
		}
		if len(exps.Data.List) == 0 {
			t.Skip("no expiration dates for " + u)
		}
		expDate := exps.Data.List[0].ExpirationDate
		chain, err := c.NewGetOptionChainService(u, expDate).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/market/chain", err) {
				return
			}
			t.Fatalf("chain: %v", err)
		}
		params := map[string]string{"underlying": u, "expiration_date": expDate}
		rawChain := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/market/chain", params, false)
		testutil.AssertCovers(t, "stock_options/market/chain", rawChain, chain)

		if len(chain.Data.List) == 0 {
			t.Skip("no contracts in option chain")
		}
		sym := chain.Data.List[0].Symbol
		ob, err := c.NewGetOrderBookService(sym).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/market/{symbol}/orderbook", err) {
				return
			}
			t.Fatalf("orderbook: %v", err)
		}
		rawOB := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/market/"+sym+"/orderbook", nil, false)
		testutil.AssertCovers(t, "stock_options/market/{symbol}/orderbook", rawOB, ob)
	})

	// ---- private read endpoints ----

	t.Run("ListPositions", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListPositionsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/positions", err) {
				return
			}
			t.Fatalf("positions: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/positions", nil, true)
		testutil.AssertCovers(t, "stock_options/positions", raw, got)
	})

	t.Run("ListOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListOrdersService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/orders", err) {
				return
			}
			t.Fatalf("orders: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/orders", nil, true)
		testutil.AssertCovers(t, "stock_options/orders", raw, got)
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
			if testutil.Tolerable(t, "stock_options/orders/history", err) {
				return
			}
			t.Fatalf("order history: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/orders/history", params, true)
		testutil.AssertCovers(t, "stock_options/orders/history", raw, got)
	})

	t.Run("GetOrder", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		hist, err := c.NewListOrderHistoryService().SetPageSize(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/orders/history", err) {
				return
			}
			t.Fatalf("order history: %v", err)
		}
		if len(hist.Data.List) == 0 {
			t.Skip("no historical stock option order to query")
		}
		orderID := hist.Data.List[0].OrderID
		got, err := c.NewGetOrderService(orderID).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "stock_options/orders/{order_id}", err) {
				return
			}
			t.Fatalf("get order: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/stock_options/orders/"+orderID, nil, true)
		testutil.AssertCovers(t, "stock_options/orders/{order_id}", raw, got)
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
		_, err := c.NewCreateOrderService("AAPL260814P00130000", 2, "limit", "day",
			decimal.NewFromFloat(0.00000001)).SetPrice(decimal.NewFromFloat(0.00000001)).Do(cx)
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
		_, err := c.NewAmendOrderService("1", decimal.NewFromFloat(0.00000001), decimal.NewFromFloat(0.00000001)).Do(cx)
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
		_, err := c.NewCancelOrderService("1").Do(cx)
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
		_, err := c.NewCancelAllOrdersService().SetSymbol("AAPL260814P00130000").Do(cx)
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
		_, err := c.NewClosePositionService("AAPL260814P00130000", decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("close position: %v (tolerable)", err)
			return
		}
		t.Log("close position accepted")
	})
}
