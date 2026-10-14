package earn

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

// TestDualStock exercises the stock Dual Investment read endpoints. They all
// require signing, so every subtest skips without credentials.
func TestDualStock(t *testing.T) {
	t.Run("ListDualStockInvestmentPlans", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		plans, err := c.NewListDualStockInvestmentPlansService().SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/stock/investment_plan", err) {
				return
			}
			t.Fatalf("stock dual plans: %v", err)
		}
		t.Logf("stock dual plans=%d", len(plans))
		if len(plans) == 0 {
			t.Log("no stock dual plans")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/stock/investment_plan",
			map[string]string{"page_size": "2"}, true)
		testutil.AssertCovers(t, "earn/dual/stock/investment_plan", raw, plans)
	})

	t.Run("ListDualStockProjectRecommend", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		projects, err := c.NewListDualStockProjectRecommendService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/stock/project-recommend", err) {
				return
			}
			t.Fatalf("stock dual project-recommend: %v", err)
		}
		t.Logf("recommended stock projects=%d", len(projects))
		if len(projects) == 0 {
			t.Log("no recommended stock projects")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/stock/project-recommend", nil, true)
		testutil.AssertCovers(t, "earn/dual/stock/project-recommend", raw, projects)
	})

	t.Run("ListDualStockOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		orders, err := c.NewListDualStockOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/stock/orders", err) {
				return
			}
			t.Fatalf("stock dual orders: %v", err)
		}
		t.Logf("stock dual orders=%d", len(orders))
		if len(orders) == 0 {
			t.Log("no stock dual orders")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/stock/orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "earn/dual/stock/orders", raw, orders)
	})

	t.Run("GetDualStockBalance", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetDualStockBalanceService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/stock/balance", err) {
				return
			}
			t.Fatalf("stock dual balance: %v", err)
		}
		t.Logf("stock dual balance: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/stock/balance", nil, true)
		testutil.AssertCovers(t, "earn/dual/stock/balance", raw, got)
	})

	t.Run("GetDualStockOrderRefundPreview", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetDualStockOrderRefundPreviewService("0").Do(cx)
		if err != nil {
			if tolerablePlaceholder(t, "earn/dual/stock/order-refund-preview", err) {
				return
			}
			t.Fatalf("stock dual order-refund-preview: %v", err)
		}
		t.Logf("preview: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/stock/order-refund-preview",
			map[string]string{"order_id": "0"}, true)
		testutil.AssertCovers(t, "earn/dual/stock/order-refund-preview", raw, got)
	})
}
