package futures

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

// TestFuturesExtra covers the futures endpoints added on top of the original
// surface: contracts_all, batch funding_rates, positions_timerange, the
// split-mode leverage getters/setters, set_position_mode, BBO orders and the
// price-triggered-order amend.
func TestFuturesExtra(t *testing.T) {
	t.Run("ListAllContracts", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		contracts, err := c.NewListAllFuturesContractsService(SettleUSDT).SetLimit(50).Do(cx)
		if err != nil {
			t.Fatalf("contracts_all: %v", err)
		}
		t.Logf("contracts_all=%d", len(contracts))
		if len(contracts) == 0 {
			t.Skip("no contracts")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/contracts_all",
			map[string]string{"limit": "50"}, false)
		testutil.AssertCovers(t, "futures/contracts_all", raw, contracts)
	})

	t.Run("BatchFundingRates", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		rates, err := c.NewBatchFundingRatesService(SettleUSDT, []string{"BTC_USDT", "ETH_USDT"}).Do(cx)
		if err != nil {
			t.Fatalf("funding_rates: %v", err)
		}
		t.Logf("funding_rates groups=%d", len(rates))
		if len(rates) == 0 {
			t.Skip("no funding rates")
		}
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/futures/usdt/funding_rates",
			map[string]any{"contracts": []string{"BTC_USDT", "ETH_USDT"}}, false)
		testutil.AssertCovers(t, "futures/funding_rates", raw, rates)
	})

	t.Run("ListPositionsWithTimeRange", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		positions, err := c.NewListPositionsWithTimeRangeService(SettleUSDT, "BTC_USDT").SetLimit(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/positions_timerange", err) {
				return
			}
			t.Fatalf("positions_timerange: %v", err)
		}
		t.Logf("positions_timerange=%d", len(positions))
		if len(positions) == 0 {
			t.Log("no historical positions")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/positions_timerange",
			map[string]string{"contract": "BTC_USDT", "limit": "10"}, true)
		testutil.AssertCovers(t, "futures/positions_timerange", raw, positions)
	})

	t.Run("GetPositionLeverage", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		info, err := c.NewGetPositionLeverageService(SettleUSDT, "BTC_USDT", "cross", "dual_long").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/get_leverage", err) {
				return
			}
			t.Fatalf("get_leverage: %v", err)
		}
		t.Logf("get_leverage=%s", info.Leverage)
	})

	// State-changing endpoints: gated behind GATE_TEST_WRITE, exercised with
	// values likely rejected by the current account/mode so nothing is placed.
	t.Run("SetPositionLeverage", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSetPositionLeverageService(SettleUSDT, "BTC_USDT", "5", "isolated").Do(cx)
		if err != nil {
			t.Logf("set_leverage: %v (tolerable)", err)
			return
		}
		t.Log("set_leverage accepted")
	})

	t.Run("SetPositionMode", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSetPositionModeService(SettleUSDT, "single").Do(cx)
		if err != nil {
			t.Logf("set_position_mode: %v (tolerable)", err)
			return
		}
		t.Log("set_position_mode accepted")
	})

	t.Run("CreateBBOOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateBBOOrderService(SettleUSDT, "BTC_USDT", 1, "buy", 1).
			SetTimeInForce(TimeInForceIOC).Do(cx)
		if err != nil {
			t.Logf("bbo_orders: %v (tolerable)", err)
			return
		}
		t.Log("bbo_orders accepted")
	})

	t.Run("AmendPriceTriggeredOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewAmendPriceTriggeredOrderService(SettleUSDT, 0).
			SetPrice(decimal.RequireFromString("10000")).Do(cx)
		if err != nil {
			t.Logf("price_orders/amend: %v (tolerable)", err)
			return
		}
		t.Log("price_orders/amend accepted")
	})
}
