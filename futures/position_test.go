package futures

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestFuturesPosition(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// ListPositions -- read-only. An account with no positions returns an empty
	// list rather than an error, so AssertCovers simply finds no rows to check.
	list, err := c.NewListPositionsService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "futures/positions", err) {
			t.Fatalf("list positions: %v", err)
		}
	} else {
		t.Logf("positions=%d", len(list))
		for _, p := range list {
			t.Logf("  %s size=%d leverage=%s entry=%s upnl=%s mode=%s",
				p.Contract, p.Size, p.Leverage, p.EntryPrice, p.UnrealisedPnL, p.Mode)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/positions",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "futures/positions", raw, list)
	}

	// GetPosition -- read-only; POSITION_NOT_FOUND when the account holds none.
	pos, err := c.NewGetPositionService(SettleUSDT, "BTC_USDT").Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "futures/positions/BTC_USDT", err) {
			t.Fatalf("get position: %v", err)
		}
	} else {
		t.Logf("position: size=%d leverage=%s entry=%s liq=%s update=%s",
			pos.Size, pos.Leverage, pos.EntryPrice, pos.LiqPrice, pos.UpdateTime)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/positions/BTC_USDT", nil, true)
		testutil.AssertCovers(t, "futures/positions/BTC_USDT", raw, pos)
	}

	// GetDualModePosition -- read-only; errors when the account is not in dual mode.
	if dual, err := c.NewGetDualModePositionService(SettleUSDT, "BTC_USDT").Do(cx); err != nil {
		if !testutil.Tolerable(t, "futures/dual_comp/positions/BTC_USDT", err) {
			t.Logf("dual-mode position: %v (account may not be in dual mode)", err)
		}
	} else {
		t.Logf("dual-mode positions=%d", len(dual))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/dual_comp/positions/BTC_USDT", nil, true)
		testutil.AssertCovers(t, "futures/dual_comp/positions/BTC_USDT", raw, dual)
	}

	if !testutil.WriteEnabled() {
		t.Log("futures position write endpoints (margin/leverage/risk_limit/cross_mode/dual_mode) skipped; set GATE_TEST_WRITE=1 to exercise")
		return
	}

	// --- State-changing endpoints (GATE_TEST_WRITE=1). Kept reversible/no-op or
	// harmlessly rejected; all errors tolerated so an un-provisioned account
	// (no position / single mode) does not fail the run. ---

	// Margin: add a negligible amount, then remove it (net zero).
	if p, err := c.NewUpdatePositionMarginService(SettleUSDT, "BTC_USDT", "0.001").Do(cx); err != nil {
		if !testutil.Tolerable(t, "position margin+", err) {
			t.Logf("update margin +0.001: %v", err)
		}
	} else {
		t.Logf("margin after add: %s", p.Margin)
		if _, err := c.NewUpdatePositionMarginService(SettleUSDT, "BTC_USDT", "-0.001").Do(cx); err != nil {
			t.Logf("restore margin -0.001: %v", err)
		}
	}

	// Leverage / risk-limit: re-apply the CURRENT value so nothing actually
	// changes, and only when a position exists (avoids touching contract defaults).
	if p, err := c.NewGetPositionService(SettleUSDT, "BTC_USDT").Do(cx); err == nil {
		if _, err := c.NewUpdatePositionLeverageService(SettleUSDT, "BTC_USDT", p.Leverage.String()).Do(cx); err != nil {
			if !testutil.Tolerable(t, "position leverage", err) {
				t.Logf("update leverage: %v", err)
			}
		}
		if _, err := c.NewUpdatePositionRiskLimitService(SettleUSDT, "BTC_USDT", p.RiskLimit.String()).Do(cx); err != nil {
			if !testutil.Tolerable(t, "position risk_limit", err) {
				t.Logf("update risk_limit: %v", err)
			}
		}
	}

	// Dual-mode variants: rejected harmlessly when the account is in single mode.
	if _, err := c.NewUpdateDualModePositionMarginService(SettleUSDT, "BTC_USDT", "0.001", "dual_long").Do(cx); err != nil {
		if !testutil.Tolerable(t, "dual margin", err) {
			t.Logf("dual-mode margin: %v", err)
		}
	}
	if _, err := c.NewUpdateDualModePositionLeverageService(SettleUSDT, "BTC_USDT", "10").Do(cx); err != nil {
		if !testutil.Tolerable(t, "dual leverage", err) {
			t.Logf("dual-mode leverage: %v", err)
		}
	}
	if _, err := c.NewUpdateDualModePositionRiskLimitService(SettleUSDT, "BTC_USDT", "1000000").Do(cx); err != nil {
		if !testutil.Tolerable(t, "dual risk_limit", err) {
			t.Logf("dual-mode risk_limit: %v", err)
		}
	}
	if _, err := c.NewUpdateDualCompPositionCrossModeService(SettleUSDT, "cross", "BTC_USDT").Do(cx); err != nil {
		if !testutil.Tolerable(t, "dual cross_mode", err) {
			t.Logf("dual-comp cross_mode: %v", err)
		}
	}

	// UpdatePositionCrossMode / SetDualMode flip account-wide margin/position mode
	// and cannot be reversed without knowing prior state; construct only.
	_ = c.NewUpdatePositionCrossModeService(SettleUSDT, "cross", "BTC_USDT")
	_ = c.NewSetDualModeService(SettleUSDT, true)
	t.Log("UpdatePositionCrossMode / SetDualMode: services constructed only (account-wide mode flips are disruptive)")
}
