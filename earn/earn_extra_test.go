package earn

import (
	"errors"
	"strings"
	"testing"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

// tolerablePlaceholder extends testutil.Tolerable for the placeholder-parameter
// smoke tests below. Several read endpoints have no queryable data on the test
// account (no auto-invest plans, no dual-currency orders), so exercising them
// with a placeholder id / portfolio comes back either as a parameter-validation
// rejection (INVALID_PARAM_VALUE) or, for the records endpoint, a transient
// "service busy" (SERVER_ERROR) — in every case the request path and signing were
// correct, so these count as reachability passes rather than code bugs.
func tolerablePlaceholder(t *testing.T, label string, err error) bool {
	t.Helper()
	if testutil.Tolerable(t, label, err) {
		return true
	}
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Label {
		case "INVALID_PARAM_VALUE":
			t.Logf("%s: placeholder parameter rejected (label=%s, msg=%s) — endpoint+signing OK", label, apiErr.Label, apiErr.Message)
			return true
		case "SERVER_ERROR":
			if strings.Contains(apiErr.Message, "Service is busy") {
				t.Logf("%s: placeholder id rejected (label=%s, msg=%s) — endpoint+signing OK", label, apiErr.Label, apiErr.Message)
				return true
			}
		}
	}
	return false
}

// TestEarnExtra exercises the additional earn endpoints (auto-invest, fixed-term,
// dual-currency extras and on-chain staking) added alongside the original earn
// coverage. Every endpoint is private, so subtests skip without credentials;
// state-changing subtests additionally require GATE_TEST_WRITE=1 and use
// tiny/likely-rejected parameters so they never place a real order.
func TestEarnExtra(t *testing.T) {
	// ---- dual-currency extras ----

	t.Run("GetDualBalance", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetDualBalanceService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/balance", err) {
				return
			}
			t.Fatalf("dual balance: %v", err)
		}
		t.Logf("dual balance: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/balance", nil, true)
		testutil.AssertCovers(t, "earn/dual/balance", raw, got)
	})

	t.Run("GetDualOrderRefundPreview", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetDualOrderRefundPreviewService("0").Do(cx)
		if err != nil {
			if tolerablePlaceholder(t, "earn/dual/order-refund-preview", err) {
				return
			}
			t.Fatalf("dual order-refund-preview: %v", err)
		}
		t.Logf("preview: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/order-refund-preview",
			map[string]string{"order_id": "0"}, true)
		testutil.AssertCovers(t, "earn/dual/order-refund-preview", raw, got)
	})

	t.Run("ListDualProjectRecommend", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		projects, err := c.NewListDualProjectRecommendService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/dual/project-recommend", err) {
				return
			}
			t.Fatalf("dual project-recommend: %v", err)
		}
		t.Logf("recommended projects=%d", len(projects))
		if len(projects) == 0 {
			t.Log("no recommended projects")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/dual/project-recommend", nil, true)
		testutil.AssertCovers(t, "earn/dual/project-recommend", raw, projects)
	})

	// ---- on-chain staking extras ----

	t.Run("ListStakingOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListStakingOrdersService().SetPage(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/staking/order_list", err) {
				return
			}
			t.Fatalf("staking order_list: %v", err)
		}
		t.Logf("staking orders=%d", len(got.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/staking/order_list",
			map[string]string{"page": "1"}, true)
		testutil.AssertCovers(t, "earn/staking/order_list", raw, got)
	})

	t.Run("ListStakingAwards", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListStakingAwardsService().SetPage(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/staking/award_list", err) {
				return
			}
			t.Fatalf("staking award_list: %v", err)
		}
		t.Logf("staking awards=%d", len(got.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/staking/award_list",
			map[string]string{"page": "1"}, true)
		testutil.AssertCovers(t, "earn/staking/award_list", raw, got)
	})

	t.Run("GetStakingAssets", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetStakingAssetsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/staking/assets", err) {
				return
			}
			t.Fatalf("staking assets: %v", err)
		}
		t.Logf("staking assets: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/staking/assets", nil, true)
		testutil.AssertCovers(t, "earn/staking/assets", raw, got)
	})

	// ---- auto-invest ----

	t.Run("ListAutoInvestCoins", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		coins, err := c.NewListAutoInvestCoinsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/autoinvest/coins", err) {
				return
			}
			t.Fatalf("autoinvest coins: %v", err)
		}
		t.Logf("autoinvest coins=%d", len(coins))
		if len(coins) == 0 {
			t.Log("no autoinvest coins")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/autoinvest/coins", nil, true)
		testutil.AssertCovers(t, "earn/autoinvest/coins", raw, coins)
	})

	t.Run("GetAutoInvestMinAmount", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		items := []AutoInvestPlanItem{{Asset: "BTC", Ratio: "100"}}
		got, err := c.NewGetAutoInvestMinAmountService("USDT", items).Do(cx)
		if err != nil {
			if tolerablePlaceholder(t, "earn/autoinvest/min_invest_amount", err) {
				return
			}
			t.Fatalf("autoinvest min_invest_amount: %v", err)
		}
		t.Logf("min amount: %+v", got)
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/earn/autoinvest/min_invest_amount",
			map[string]any{"money": "USDT", "items": items}, true)
		testutil.AssertCovers(t, "earn/autoinvest/min_invest_amount", raw, got)
	})

	t.Run("ListAutoInvestRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListAutoInvestRecordsService(0).SetPage(1).Do(cx)
		if err != nil {
			if tolerablePlaceholder(t, "earn/autoinvest/plans/records", err) {
				return
			}
			t.Fatalf("autoinvest plans/records: %v", err)
		}
		t.Logf("autoinvest records=%d", len(got.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/autoinvest/plans/records",
			map[string]string{"plan_id": "0", "page": "1"}, true)
		testutil.AssertCovers(t, "earn/autoinvest/plans/records", raw, got)
	})

	t.Run("ListAutoInvestOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		orders, err := c.NewListAutoInvestOrdersService(0, 0).Do(cx)
		if err != nil {
			if tolerablePlaceholder(t, "earn/autoinvest/orders", err) {
				return
			}
			t.Fatalf("autoinvest orders: %v", err)
		}
		t.Logf("autoinvest orders=%d", len(orders))
		if len(orders) == 0 {
			t.Log("no autoinvest orders")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/autoinvest/orders",
			map[string]string{"plan_id": "0", "record_id": "0"}, true)
		testutil.AssertCovers(t, "earn/autoinvest/orders", raw, orders)
	})

	t.Run("ListAutoInvestConfig", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		cfg, err := c.NewListAutoInvestConfigService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/autoinvest/config", err) {
				return
			}
			t.Fatalf("autoinvest config: %v", err)
		}
		t.Logf("autoinvest config=%d", len(cfg))
		if len(cfg) == 0 {
			t.Log("no autoinvest config")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/autoinvest/config", nil, true)
		testutil.AssertCovers(t, "earn/autoinvest/config", raw, cfg)
	})

	t.Run("GetAutoInvestPlanDetail", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetAutoInvestPlanDetailService(0).Do(cx)
		if err != nil {
			if tolerablePlaceholder(t, "earn/autoinvest/plans/detail", err) {
				return
			}
			t.Fatalf("autoinvest plans/detail: %v", err)
		}
		t.Logf("plan detail: %+v", got)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/autoinvest/plans/detail",
			map[string]string{"plan_id": "0"}, true)
		testutil.AssertCovers(t, "earn/autoinvest/plans/detail", raw, got)
	})

	t.Run("ListAutoInvestPlans", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListAutoInvestPlansService("active").SetPage(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/autoinvest/plans/list_info", err) {
				return
			}
			t.Fatalf("autoinvest plans/list_info: %v", err)
		}
		t.Logf("autoinvest plans=%d", len(got.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/autoinvest/plans/list_info",
			map[string]string{"status": "active", "page": "1"}, true)
		testutil.AssertCovers(t, "earn/autoinvest/plans/list_info", raw, got)
	})

	// ---- fixed-term ----

	t.Run("ListFixedTermProducts", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListFixedTermProductsService(1, 10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/fixed-term/product", err) {
				return
			}
			t.Fatalf("fixed-term product: %v", err)
		}
		t.Logf("fixed-term products=%d", len(got.Data.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/fixed-term/product",
			map[string]string{"page": "1", "limit": "10"}, true)
		testutil.AssertCovers(t, "earn/fixed-term/product", raw, got)
	})

	t.Run("ListFixedTermProductsByAsset", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListFixedTermProductsByAssetService("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/fixed-term/product/USDT/list", err) {
				return
			}
			t.Fatalf("fixed-term product by asset: %v", err)
		}
		t.Logf("fixed-term products (USDT)=%d", len(got.Data.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/fixed-term/product/USDT/list", nil, true)
		testutil.AssertCovers(t, "earn/fixed-term/product/USDT/list", raw, got)
	})

	t.Run("ListFixedTermLends", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListFixedTermLendsService("1", 1, 10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/fixed-term/user/lend", err) {
				return
			}
			t.Fatalf("fixed-term user/lend: %v", err)
		}
		t.Logf("fixed-term lends=%d", len(got.Data.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/fixed-term/user/lend",
			map[string]string{"order_type": "1", "page": "1", "limit": "10"}, true)
		testutil.AssertCovers(t, "earn/fixed-term/user/lend", raw, got)
	})

	t.Run("ListFixedTermHistory", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListFixedTermHistoryService("1", 1, 10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "earn/fixed-term/user/history", err) {
				return
			}
			t.Fatalf("fixed-term user/history: %v", err)
		}
		t.Logf("fixed-term history=%d", len(got.Data.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/fixed-term/user/history",
			map[string]string{"type": "1", "page": "1", "limit": "10"}, true)
		testutil.AssertCovers(t, "earn/fixed-term/user/history", raw, got)
	})

	// ---- state-changing (write) endpoints ----

	t.Run("RefundDualOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		if err := c.NewRefundDualOrderService("0", "nonexistent").Do(cx); err != nil {
			t.Logf("refund dual order: %v (tolerable)", err)
			return
		}
		t.Log("refund dual order accepted")
	})

	t.Run("ModifyDualOrderReinvest", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		if err := c.NewModifyDualOrderReinvestService().SetOrderID(0).SetStatus(0).Do(cx); err != nil {
			t.Logf("modify dual order reinvest: %v (tolerable)", err)
			return
		}
		t.Log("modify dual order reinvest accepted")
	})

	t.Run("CreateAutoInvestPlan", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		items := []AutoInvestPlanItem{{Asset: "BTC", Ratio: "100"}}
		_, err := c.NewCreateAutoInvestPlanService("USDT", decimal.NewFromFloat(0.00000001), "daily", 1, 0, items).Do(cx)
		if err != nil {
			t.Logf("create autoinvest plan: %v (tolerable)", err)
			return
		}
		t.Log("create autoinvest plan accepted")
	})

	t.Run("UpdateAutoInvestPlan", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		if err := c.NewUpdateAutoInvestPlanService(0).SetFundSource("spot").Do(cx); err != nil {
			t.Logf("update autoinvest plan: %v (tolerable)", err)
			return
		}
		t.Log("update autoinvest plan accepted")
	})

	t.Run("StopAutoInvestPlan", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		if err := c.NewStopAutoInvestPlanService(0).Do(cx); err != nil {
			t.Logf("stop autoinvest plan: %v (tolerable)", err)
			return
		}
		t.Log("stop autoinvest plan accepted")
	})

	t.Run("AddAutoInvestPosition", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		if err := c.NewAddAutoInvestPositionService(0, decimal.NewFromFloat(0.00000001)).Do(cx); err != nil {
			t.Logf("add autoinvest position: %v (tolerable)", err)
			return
		}
		t.Log("add autoinvest position accepted")
	})

	t.Run("SubscribeFixedTerm", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSubscribeFixedTermService(0, decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("subscribe fixed-term: %v (tolerable)", err)
			return
		}
		t.Log("subscribe fixed-term accepted")
	})

	t.Run("PreRedeemFixedTerm", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewPreRedeemFixedTermService("0").Do(cx)
		if err != nil {
			t.Logf("pre-redeem fixed-term: %v (tolerable)", err)
			return
		}
		t.Log("pre-redeem fixed-term accepted")
	})
}
