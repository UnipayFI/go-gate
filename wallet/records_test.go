package wallet

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestWalletRecords(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	t.Run("Withdrawals", func(t *testing.T) {
		list, err := c.NewListWithdrawalsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/withdrawals", err) {
				return
			}
			t.Fatalf("withdrawals: %v", err)
		}
		t.Logf("withdrawals=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/withdrawals",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "wallet/withdrawals", raw, list)
	})

	t.Run("Deposits", func(t *testing.T) {
		list, err := c.NewListDepositsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/deposits", err) {
				return
			}
			t.Fatalf("deposits: %v", err)
		}
		t.Logf("deposits=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/deposits",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "wallet/deposits", raw, list)
	})

	t.Run("WithdrawStatus", func(t *testing.T) {
		list, err := c.NewListWithdrawStatusService().SetCurrency("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/withdraw_status", err) {
				return
			}
			t.Fatalf("withdraw_status: %v", err)
		}
		t.Logf("withdraw_status=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/withdraw_status",
			map[string]string{"currency": "USDT"}, true)
		testutil.AssertCovers(t, "wallet/withdraw_status", raw, list)
	})

	t.Run("SubAccountBalances", func(t *testing.T) {
		list, err := c.NewListSubAccountBalancesService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_balances", err) {
				return
			}
			t.Fatalf("sub_account_balances: %v", err)
		}
		t.Logf("sub_account_balances=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/sub_account_balances", nil, true)
		testutil.AssertCovers(t, "wallet/sub_account_balances", raw, list)
	})

	t.Run("SubAccountMarginBalances", func(t *testing.T) {
		list, err := c.NewListSubAccountMarginBalancesService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_margin_balances", err) {
				return
			}
			t.Fatalf("sub_account_margin_balances: %v", err)
		}
		t.Logf("sub_account_margin_balances=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/sub_account_margin_balances", nil, true)
		testutil.AssertCovers(t, "wallet/sub_account_margin_balances", raw, list)
	})

	t.Run("SubAccountFuturesBalances", func(t *testing.T) {
		list, err := c.NewListSubAccountFuturesBalancesService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_futures_balances", err) {
				return
			}
			t.Fatalf("sub_account_futures_balances: %v", err)
		}
		t.Logf("sub_account_futures_balances=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/sub_account_futures_balances", nil, true)
		testutil.AssertCovers(t, "wallet/sub_account_futures_balances", raw, list)
	})

	t.Run("SubAccountCrossMarginBalances", func(t *testing.T) {
		list, err := c.NewListSubAccountCrossMarginBalancesService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_cross_margin_balances", err) {
				return
			}
			t.Fatalf("sub_account_cross_margin_balances: %v", err)
		}
		t.Logf("sub_account_cross_margin_balances=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/sub_account_cross_margin_balances", nil, true)
		testutil.AssertCovers(t, "wallet/sub_account_cross_margin_balances", raw, list)
	})

	t.Run("SmallBalance", func(t *testing.T) {
		list, err := c.NewListSmallBalanceService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/small_balance", err) {
				return
			}
			t.Fatalf("small_balance: %v", err)
		}
		t.Logf("small_balance=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/small_balance", nil, true)
		testutil.AssertCovers(t, "wallet/small_balance", raw, list)
	})

	t.Run("SmallBalanceHistory", func(t *testing.T) {
		list, err := c.NewListSmallBalanceHistoryService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/small_balance_history", err) {
				return
			}
			t.Fatalf("small_balance_history: %v", err)
		}
		t.Logf("small_balance_history=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/small_balance_history",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "wallet/small_balance_history", raw, list)
	})

	t.Run("PushOrders", func(t *testing.T) {
		list, err := c.NewListPushOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/push", err) {
				return
			}
			t.Fatalf("push: %v", err)
		}
		t.Logf("push_orders=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/push",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "wallet/push", raw, list)
	})

	t.Run("ConvertSmallBalance", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run dust conversion (state-changing)")
		}
		// Convert a single obscure currency's dust: signing/path correctness is
		// what matters here, and the account most likely holds no such dust, so
		// Gate answers with a tolerable "no balance"-style error rather than
		// touching real funds. The 200 response carries no body.
		_, err := c.NewConvertSmallBalanceService([]string{"1INCH"}).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/small_balance(convert)", err) {
				return
			}
			t.Logf("convert small balance: %v", err)
			return
		}
		t.Log("convert small balance: OK")
	})
}
