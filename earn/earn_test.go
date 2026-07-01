package earn

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestEarn(t *testing.T) {
	t.Run("ListDualInvestmentPlans", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		plans, err := c.NewListDualInvestmentPlansService().Do(cx)
		if err != nil {
			t.Fatalf("dual investment plans: %v", err)
		}
		t.Logf("dual plans=%d", len(plans))
		if len(plans) == 0 {
			t.Skip("no dual investment plans available")
		}
		t.Logf("plan[0]: %+v", plans[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/investment_plan", nil, false)
		testutil.AssertCovers(t, "earn/dual/investment_plan", raw, plans)
	})

	t.Run("ListStructuredProducts", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		params := map[string]string{"status": "in_process", "limit": "2"}
		products, err := c.NewListStructuredProductsService("in_process").SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("structured products: %v", err)
		}
		t.Logf("structured products=%d", len(products))
		if len(products) == 0 {
			t.Skip("no in_process structured products available")
		}
		t.Logf("product[0]: %+v", products[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/structured/products", params, false)
		testutil.AssertCovers(t, "earn/structured/products", raw, products)
	})

	t.Run("FindCoin", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		coins, err := c.NewFindCoinService().Do(cx)
		if err != nil {
			t.Fatalf("staking coins: %v", err)
		}
		t.Logf("staking coins=%d", len(coins))
		if len(coins) == 0 {
			t.Skip("no staking coins available")
		}
		t.Logf("coin[0]: %+v", coins[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/staking/coins", nil, false)
		testutil.AssertCovers(t, "earn/staking/coins", raw, coins)
	})

	t.Run("RateListETH2", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		records, err := c.NewRateListETH2Service().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/staking/eth2/rate_records", err) {
				return
			}
			t.Fatalf("eth2 rate records: %v", err)
		}
		t.Logf("eth2 rate records=%d", len(records.Rates))
		if len(records.Rates) == 0 {
			t.Log("no eth2 rate records")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/staking/eth2/rate_records", nil, true)
		testutil.AssertCovers(t, "earn/staking/eth2/rate_records", raw, records)
	})

	t.Run("ListDualOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		orders, err := c.NewListDualOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/orders", err) {
				return
			}
			t.Fatalf("dual orders: %v", err)
		}
		t.Logf("dual orders=%d", len(orders))
		if len(orders) == 0 {
			t.Log("no dual orders")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "earn/dual/orders", raw, orders)
	})

	t.Run("ListStructuredOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		orders, err := c.NewListStructuredOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/structured/orders", err) {
				return
			}
			t.Fatalf("structured orders: %v", err)
		}
		t.Logf("structured orders=%d", len(orders))
		if len(orders) == 0 {
			t.Log("no structured orders")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/structured/orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "earn/structured/orders", raw, orders)
	})

	// State-changing subscriptions/swaps. Gated behind GATE_TEST_WRITE and
	// exercised with tiny/likely-rejected parameters so they never place a real
	// investment; any error is treated as a pass.
	t.Run("SwapETH2", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSwapETH2Service("1", decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("eth2 swap: %v (tolerable)", err)
			return
		}
		t.Log("eth2 swap accepted")
	})

	t.Run("PlaceDualOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		err := c.NewPlaceDualOrderService("0", decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("place dual order: %v (tolerable)", err)
			return
		}
		t.Log("place dual order accepted")
	})

	t.Run("PlaceStructuredOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		err := c.NewPlaceStructuredOrderService("0", decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("place structured order: %v (tolerable)", err)
			return
		}
		t.Log("place structured order accepted")
	})

	t.Run("SwapStakingCoin", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSwapStakingCoinService("GT", "0", decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("swap staking coin: %v (tolerable)", err)
			return
		}
		t.Log("swap staking coin accepted")
	})
}
