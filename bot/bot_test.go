package bot

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestBot(t *testing.T) {
	// tiny is a dust-sized amount used by the write-gated tests so they are
	// rejected by Gate long before a real strategy is ever created.
	tiny := decimal.NewFromFloat(0.00000001)

	// ---- Private read endpoints ----

	t.Run("GetStrategyRecommend", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewGetStrategyRecommendService().SetMarket("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "bot/strategy/recommend", err) {
				return
			}
			t.Fatalf("strategy recommend: %v", err)
		}
		t.Logf("recommendations=%d", len(got.Data.Recommendations))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/bot/strategy/recommend",
			map[string]string{"market": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "bot/strategy/recommend", raw, got)
	})

	t.Run("ListPortfolioRunning", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListPortfolioRunningService().SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "bot/portfolio/running", err) {
				return
			}
			t.Fatalf("portfolio running: %v", err)
		}
		t.Logf("running items=%d total=%d", len(got.Data.Items), got.Data.Total)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/bot/portfolio/running",
			map[string]string{"page_size": "2"}, true)
		testutil.AssertCovers(t, "bot/portfolio/running", raw, got)
	})

	t.Run("GetPortfolioDetail", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		// Detail needs a real strategy id/type, so discover one from the running
		// list first; skip when the account has none.
		running, err := c.NewListPortfolioRunningService().SetPageSize(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "bot/portfolio/detail", err) {
				return
			}
			t.Fatalf("portfolio running (for detail): %v", err)
		}
		if len(running.Data.Items) == 0 {
			t.Skip("no running strategy to fetch detail for")
		}
		item := running.Data.Items[0]
		got, err := c.NewGetPortfolioDetailService(item.StrategyID, string(item.StrategyType)).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "bot/portfolio/detail", err) {
				return
			}
			t.Fatalf("portfolio detail: %v", err)
		}
		t.Logf("detail status=%s", got.Data.Status)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/bot/portfolio/detail",
			map[string]string{"strategy_id": item.StrategyID, "strategy_type": string(item.StrategyType)}, true)
		testutil.AssertCovers(t, "bot/portfolio/detail", raw, got)
	})

	// ---- Write endpoints ----
	// All strategy-creation and stop calls are gated behind GATE_TEST_WRITE and
	// use dust-sized / likely-rejected parameters so they never open a real
	// strategy; any error is treated as a pass.

	t.Run("CreateSpotGrid", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateSpotGridService("spot_grid", "BTC_USDT", tiny, tiny, tiny.Mul(decimal.NewFromInt(2)), 2, 0).Do(cx)
		if err != nil {
			t.Logf("create spot grid: %v (tolerable)", err)
			return
		}
		t.Log("create spot grid accepted")
	})

	t.Run("CreateMarginGrid", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateMarginGridService("margin_grid", "BTC_USDT", tiny, tiny, tiny.Mul(decimal.NewFromInt(2)), 2, 0, decimal.NewFromInt(1)).Do(cx)
		if err != nil {
			t.Logf("create margin grid: %v (tolerable)", err)
			return
		}
		t.Log("create margin grid accepted")
	})

	t.Run("CreateInfiniteGrid", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateInfiniteGridService("infinite_grid", "BTC_USDT", tiny, tiny, tiny).Do(cx)
		if err != nil {
			t.Logf("create infinite grid: %v (tolerable)", err)
			return
		}
		t.Log("create infinite grid accepted")
	})

	t.Run("CreateFuturesGrid", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateFuturesGridService("futures_grid", "BTC_USDT", tiny, tiny, tiny.Mul(decimal.NewFromInt(2)), 2, 0, decimal.NewFromInt(1)).Do(cx)
		if err != nil {
			t.Logf("create futures grid: %v (tolerable)", err)
			return
		}
		t.Log("create futures grid accepted")
	})

	t.Run("CreateSpotMartingale", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateSpotMartingaleService("spot_martingale", "BTC_USDT", tiny, decimal.NewFromFloat(0.02), 2, decimal.NewFromFloat(0.01)).Do(cx)
		if err != nil {
			t.Logf("create spot martingale: %v (tolerable)", err)
			return
		}
		t.Log("create spot martingale accepted")
	})

	t.Run("CreateContractMartingale", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateContractMartingaleService("contract_martingale", "BTC_USDT", tiny, decimal.NewFromFloat(0.02), 2, decimal.NewFromFloat(0.01), ContractMartingaleDirection("long"), decimal.NewFromInt(1)).Do(cx)
		if err != nil {
			t.Logf("create contract martingale: %v (tolerable)", err)
			return
		}
		t.Log("create contract martingale accepted")
	})

	t.Run("StopPortfolio", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewStopPortfolioService("0", "spot_grid").Do(cx)
		if err != nil {
			t.Logf("stop portfolio: %v (tolerable)", err)
			return
		}
		t.Log("stop portfolio accepted")
	})
}
