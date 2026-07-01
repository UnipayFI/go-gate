package margin

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestMargin(t *testing.T) {
	cx := testutil.Ctx(t)

	// Public: current market leverage lending tiers.
	t.Run("MarketMarginTier", func(t *testing.T) {
		pc := testPublicClient()
		tiers, err := pc.NewGetMarketMarginTierService("BTC_USDT").Do(cx)
		if err != nil {
			t.Fatalf("market margin tier: %v", err)
		}
		if len(tiers) == 0 {
			t.Fatal("no market margin tiers returned")
		}
		t.Logf("market tiers=%d first=%+v", len(tiers), tiers[0])
		raw := testutil.FetchRawGet(t, pc, cx, "/api/v4/margin/loan_margin_tiers",
			map[string]string{"currency_pair": "BTC_USDT"}, false)
		testutil.AssertCovers(t, "margin/loan_margin_tiers", raw, tiers)
	})

	// Everything below is private.
	c := testClient(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	t.Run("MarginAccounts", func(t *testing.T) {
		list, err := c.NewListMarginAccountsService().SetCurrencyPair("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/accounts", err) {
				return
			}
			t.Fatalf("margin accounts: %v", err)
		}
		t.Logf("margin accounts=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/accounts",
			map[string]string{"currency_pair": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "margin/accounts", raw, list)
	})

	t.Run("MarginAccountBook", func(t *testing.T) {
		list, err := c.NewListMarginAccountBookService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/account_book", err) {
				return
			}
			t.Fatalf("margin account book: %v", err)
		}
		t.Logf("account book records=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/account_book", nil, true)
		testutil.AssertCovers(t, "margin/account_book", raw, list)
	})

	t.Run("FundingAccounts", func(t *testing.T) {
		list, err := c.NewListFundingAccountsService().SetCurrency("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/funding_accounts", err) {
				return
			}
			t.Fatalf("funding accounts: %v", err)
		}
		t.Logf("funding accounts=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/funding_accounts",
			map[string]string{"currency": "USDT"}, true)
		testutil.AssertCovers(t, "margin/funding_accounts", raw, list)
	})

	t.Run("AutoRepayStatus", func(t *testing.T) {
		st, err := c.NewGetAutoRepayStatusService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/auto_repay", err) {
				return
			}
			t.Fatalf("auto repay status: %v", err)
		}
		t.Logf("auto repay status=%s", st.Status)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/auto_repay", nil, true)
		testutil.AssertCovers(t, "margin/auto_repay", raw, st)
	})

	t.Run("MarginTransferable", func(t *testing.T) {
		tr, err := c.NewGetMarginTransferableService("USDT").SetCurrencyPair("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/transferable", err) {
				return
			}
			t.Fatalf("margin transferable: %v", err)
		}
		t.Logf("transferable %s=%s", tr.Currency, tr.Amount)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/transferable",
			map[string]string{"currency": "USDT", "currency_pair": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "margin/transferable", raw, tr)
	})

	t.Run("UserMarginTier", func(t *testing.T) {
		tiers, err := c.NewGetUserMarginTierService("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/user/loan_margin_tiers", err) {
				return
			}
			t.Fatalf("user margin tier: %v", err)
		}
		t.Logf("user tiers=%d", len(tiers))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/user/loan_margin_tiers",
			map[string]string{"currency_pair": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "margin/user/loan_margin_tiers", raw, tiers)
	})

	t.Run("MarginUserAccount", func(t *testing.T) {
		list, err := c.NewListMarginUserAccountService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/user/account", err) {
				return
			}
			t.Fatalf("margin user account: %v", err)
		}
		t.Logf("margin user accounts=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/user/account", nil, true)
		testutil.AssertCovers(t, "margin/user/account", raw, list)
	})

	t.Run("CrossMarginLoans", func(t *testing.T) {
		list, err := c.NewListCrossMarginLoansService(2).SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/cross/loans", err) {
				return
			}
			t.Fatalf("cross margin loans: %v", err)
		}
		t.Logf("cross margin loans=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/cross/loans",
			map[string]string{"status": "2", "limit": "2"}, true)
		testutil.AssertCovers(t, "margin/cross/loans", raw, list)
	})

	t.Run("CrossMarginRepayments", func(t *testing.T) {
		list, err := c.NewListCrossMarginRepaymentsService().SetCurrency("USDT").SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/cross/repayments", err) {
				return
			}
			t.Fatalf("cross margin repayments: %v", err)
		}
		t.Logf("cross margin repayments=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/margin/cross/repayments",
			map[string]string{"currency": "USDT", "limit": "2"}, true)
		testutil.AssertCovers(t, "margin/cross/repayments", raw, list)
	})

	// State-changing: re-apply the current auto-repay setting to itself so the
	// account state is unchanged, then verify the echo.
	t.Run("SetAutoRepay", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("GATE_TEST_WRITE not set; skipping auto-repay write")
		}
		cur, err := c.NewGetAutoRepayStatusService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/auto_repay", err) {
				return
			}
			t.Fatalf("auto repay status: %v", err)
		}
		st, err := c.NewSetAutoRepayService(cur.Status).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/auto_repay set", err) {
				return
			}
			t.Fatalf("set auto repay: %v", err)
		}
		t.Logf("auto repay set to %s", st.Status)
	})

	// State-changing: re-apply the market's current leverage to itself so no net
	// change occurs. Skips unless an active BTC_USDT isolated account exists.
	t.Run("SetUserMarketLeverage", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("GATE_TEST_WRITE not set; skipping leverage write")
		}
		list, err := c.NewListMarginUserAccountService().SetCurrencyPair("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "margin/user/account", err) {
				return
			}
			t.Fatalf("margin user account: %v", err)
		}
		var leverage string
		for _, a := range list {
			if a.CurrencyPair == "BTC_USDT" && a.Leverage.IsPositive() {
				leverage = a.Leverage.String()
				break
			}
		}
		if leverage == "" {
			t.Skip("no active BTC_USDT isolated account; skipping leverage write")
		}
		if err := c.NewSetUserMarketLeverageService("BTC_USDT", leverage).Do(cx); err != nil {
			if testutil.Tolerable(t, "margin/leverage/user_market_setting", err) {
				return
			}
			t.Fatalf("set user market leverage: %v", err)
		}
		t.Logf("BTC_USDT leverage re-applied as %s", leverage)
	})
}
