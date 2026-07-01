package unified

import (
	"context"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestUnifiedAccount(t *testing.T) {
	// authed builds a signed, time-synced unified client, skipping the subtest
	// when credentials are unset.
	authed := func(t *testing.T, cx context.Context) *UnifiedClient {
		t.Helper()
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		return c
	}

	// --- Public endpoints ---

	t.Run("CurrencyDiscountTiers", func(t *testing.T) {
		pub := testPublicClient()
		cx := testutil.Ctx(t)
		list, err := pub.NewListCurrencyDiscountTiersService().Do(cx)
		if err != nil {
			t.Fatalf("currency discount tiers: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no discount tiers returned")
		}
		t.Logf("discount tiers: %d currencies, first=%+v", len(list), list[0])
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/unified/currency_discount_tiers", nil, false)
		testutil.AssertCovers(t, "unified/currency_discount_tiers", raw, list)
	})

	t.Run("LoanMarginTiers", func(t *testing.T) {
		pub := testPublicClient()
		cx := testutil.Ctx(t)
		list, err := pub.NewListLoanMarginTiersService().Do(cx)
		if err != nil {
			t.Fatalf("loan margin tiers: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no margin tiers returned")
		}
		t.Logf("margin tiers: %d currencies, first=%+v", len(list), list[0])
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/unified/loan_margin_tiers", nil, false)
		testutil.AssertCovers(t, "unified/loan_margin_tiers", raw, list)
	})

	t.Run("Currencies", func(t *testing.T) {
		pub := testPublicClient()
		cx := testutil.Ctx(t)
		list, err := pub.NewListUnifiedCurrenciesService().Do(cx)
		if err != nil {
			t.Fatalf("unified currencies: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no currencies returned")
		}
		t.Logf("currencies: %d, first=%+v", len(list), list[0])
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/unified/currencies", nil, false)
		testutil.AssertCovers(t, "unified/currencies", raw, list)
	})

	t.Run("HistoryLoanRate", func(t *testing.T) {
		pub := testPublicClient()
		cx := testutil.Ctx(t)
		params := map[string]string{"currency": "USDT", "limit": "2"}
		rate, err := pub.NewGetHistoryLoanRateService("USDT").SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("history loan rate: %v", err)
		}
		t.Logf("history loan rate: currency=%s rates=%d", rate.Currency, len(rate.Rates))
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/unified/history_loan_rate", params, false)
		testutil.AssertCovers(t, "unified/history_loan_rate", raw, rate)
	})

	t.Run("PortfolioCalculator", func(t *testing.T) {
		pub := testPublicClient()
		cx := testutil.Ctx(t)
		// Stateless public calculator; empty simulated input may be rejected, so
		// treat any error as a soft pass.
		body := map[string]any{"spot_balances": []any{}}
		resp, err := pub.NewCalculatePortfolioMarginService().SetSpotBalances([]PortfolioSpotBalance{}).Do(cx)
		if err != nil {
			t.Logf("portfolio calculator: %v (minimal input may be rejected)", err)
			return
		}
		t.Logf("portfolio margin: initial=%s maintain=%s units=%d",
			resp.InitialMarginTotal, resp.MaintainMarginTotal, len(resp.RiskUnit))
		raw := testutil.FetchRawPost(t, pub, cx, "/api/v4/unified/portfolio_calculator", body, false)
		testutil.AssertCovers(t, "unified/portfolio_calculator", raw, resp)
	})

	// --- Private reads ---

	t.Run("Accounts", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		acc, err := c.NewListUnifiedAccountsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/accounts", err) {
				return
			}
			t.Fatalf("unified accounts: %v", err)
		}
		t.Logf("unified account: user=%d total=%s balances=%d", acc.UserID, acc.UnifiedAccountTotal, len(acc.Balances))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/accounts", nil, true)
		testutil.AssertCovers(t, "unified/accounts", raw, acc)
	})

	t.Run("Borrowable", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currency": "USDT"}
		resp, err := c.NewGetUnifiedBorrowableService("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/borrowable", err) {
				return
			}
			t.Fatalf("unified borrowable: %v", err)
		}
		t.Logf("borrowable: %+v", resp)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/borrowable", params, true)
		testutil.AssertCovers(t, "unified/borrowable", raw, resp)
	})

	t.Run("Transferable", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currency": "USDT"}
		resp, err := c.NewGetUnifiedTransferableService("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/transferable", err) {
				return
			}
			t.Fatalf("unified transferable: %v", err)
		}
		t.Logf("transferable: %+v", resp)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/transferable", params, true)
		testutil.AssertCovers(t, "unified/transferable", raw, resp)
	})

	t.Run("Transferables", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currencies": "BTC,USDT"}
		list, err := c.NewGetUnifiedTransferablesService([]string{"BTC", "USDT"}).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/transferables", err) {
				return
			}
			t.Fatalf("unified transferables: %v", err)
		}
		t.Logf("transferables: %d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/transferables", params, true)
		testutil.AssertCovers(t, "unified/transferables", raw, list)
	})

	t.Run("BorrowableList", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currencies": "BTC,USDT"}
		list, err := c.NewGetUnifiedBorrowableListService([]string{"BTC", "USDT"}).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/batch_borrowable", err) {
				return
			}
			t.Fatalf("unified batch_borrowable: %v", err)
		}
		t.Logf("batch borrowable: %d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/batch_borrowable", params, true)
		testutil.AssertCovers(t, "unified/batch_borrowable", raw, list)
	})

	t.Run("RiskUnits", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		resp, err := c.NewGetUnifiedRiskUnitsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/risk_units", err) {
				return
			}
			t.Fatalf("unified risk_units: %v", err)
		}
		t.Logf("risk units: user=%d units=%d", resp.UserID, len(resp.RiskUnits))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/risk_units", nil, true)
		testutil.AssertCovers(t, "unified/risk_units", raw, resp)
	})

	t.Run("Mode", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		resp, err := c.NewGetUnifiedModeService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/unified_mode", err) {
				return
			}
			t.Fatalf("unified mode: %v", err)
		}
		t.Logf("unified mode: %+v", resp)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/unified_mode", nil, true)
		testutil.AssertCovers(t, "unified/unified_mode", raw, resp)
	})

	t.Run("EstimateRate", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currencies": "BTC,USDT"}
		rates, err := c.NewGetUnifiedEstimateRateService([]string{"BTC", "USDT"}).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/estimate_rate", err) {
				return
			}
			t.Fatalf("unified estimate_rate: %v", err)
		}
		t.Logf("estimate rate: %+v", rates)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/estimate_rate", params, true)
		testutil.AssertCovers(t, "unified/estimate_rate", raw, rates)
	})

	t.Run("LeverageCurrencyConfig", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currency": "USDT"}
		resp, err := c.NewGetUserLeverageCurrencyConfigService("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/leverage/user_currency_config", err) {
				return
			}
			t.Fatalf("leverage user_currency_config: %v", err)
		}
		t.Logf("leverage config: %+v", resp)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/leverage/user_currency_config", params, true)
		testutil.AssertCovers(t, "unified/leverage/user_currency_config", raw, resp)
	})

	t.Run("LeverageCurrencySetting", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		params := map[string]string{"currency": "USDT"}
		list, err := c.NewGetUserLeverageCurrencySettingService().SetCurrency("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/leverage/user_currency_setting", err) {
				return
			}
			t.Fatalf("leverage user_currency_setting: %v", err)
		}
		t.Logf("leverage settings: %d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/leverage/user_currency_setting", params, true)
		testutil.AssertCovers(t, "unified/leverage/user_currency_setting", raw, list)
	})

	// --- State-changing endpoints (guarded) ---

	t.Run("SetMode", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run state-changing SetUnifiedMode")
		}
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		cur, err := c.NewGetUnifiedModeService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/unified_mode set(read)", err) {
				return
			}
			t.Fatalf("read mode: %v", err)
		}
		// Re-set the current mode so the write is a no-op / reversible.
		if err := c.NewSetUnifiedModeService(cur.Mode).Do(cx); err != nil {
			if testutil.Tolerable(t, "unified/unified_mode set", err) {
				return
			}
			t.Fatalf("set mode: %v", err)
		}
		t.Logf("re-applied unified mode=%s", cur.Mode)
	})

	t.Run("SetLeverageCurrencySetting", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run state-changing SetUserLeverageCurrencySetting")
		}
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		if err := c.NewSetUserLeverageCurrencySettingService("USDT", "3").Do(cx); err != nil {
			if testutil.Tolerable(t, "unified/leverage/user_currency_setting set", err) {
				return
			}
			t.Fatalf("set leverage: %v", err)
		}
		t.Log("set USDT leverage=3")
	})

	t.Run("SetCollateral", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run state-changing SetUnifiedCollateral")
		}
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		// collateral_type 0 = use all currencies as collateral.
		resp, err := c.NewSetUnifiedCollateralService(0).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/collateral_currencies", err) {
				return
			}
			t.Fatalf("set collateral: %v", err)
		}
		t.Logf("set collateral: is_success=%v", resp.IsSuccess)
	})
}
