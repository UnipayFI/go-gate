package wallet

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestWalletTransfer(t *testing.T) {
	// ListCurrencyChains is public and always runs.
	t.Run("ListCurrencyChains", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)
		chains, err := c.NewListCurrencyChainsService("USDT").Do(cx)
		if err != nil {
			t.Fatalf("currency chains: %v", err)
		}
		if len(chains) == 0 {
			t.Fatal("no currency chains returned")
		}
		t.Logf("chains=%d first=%+v", len(chains), chains[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/currency_chains",
			map[string]string{"currency": "USDT"}, false)
		testutil.AssertCovers(t, "wallet/currency_chains", raw, chains)
	})

	// Everything below is private; testClient skips the rest when creds are unset.
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	t.Run("GetTotalBalance", func(t *testing.T) {
		bal, err := c.NewGetTotalBalanceService().SetCurrency("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/total_balance", err) {
				return
			}
			t.Fatalf("total balance: %v", err)
		}
		t.Logf("total=%+v details=%d", bal.Total, len(bal.Details))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/total_balance",
			map[string]string{"currency": "USDT"}, true)
		testutil.AssertCovers(t, "wallet/total_balance", raw, bal)
	})

	t.Run("GetTradeFee", func(t *testing.T) {
		fee, err := c.NewGetTradeFeeService().SetCurrencyPair("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/fee", err) {
				return
			}
			t.Fatalf("trade fee: %v", err)
		}
		t.Logf("fee=%+v", fee)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/fee",
			map[string]string{"currency_pair": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "wallet/fee", raw, fee)
	})

	t.Run("GetDepositAddress", func(t *testing.T) {
		addr, err := c.NewGetDepositAddressService("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/deposit_address", err) {
				return
			}
			t.Fatalf("deposit address: %v", err)
		}
		t.Logf("addr=%+v", addr)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/deposit_address",
			map[string]string{"currency": "USDT"}, true)
		testutil.AssertCovers(t, "wallet/deposit_address", raw, addr)
	})

	t.Run("ListSubAccountTransfers", func(t *testing.T) {
		list, err := c.NewListSubAccountTransfersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_transfers", err) {
				return
			}
			t.Fatalf("sub account transfers: %v", err)
		}
		t.Logf("sub transfers=%d", len(list))
		if len(list) == 0 {
			t.Log("no sub-account transfer records; skipping cover")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/sub_account_transfers",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "wallet/sub_account_transfers", raw, list)
	})

	t.Run("ListSavedAddress", func(t *testing.T) {
		list, err := c.NewListSavedAddressService("USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/saved_address", err) {
				return
			}
			t.Fatalf("saved address: %v", err)
		}
		t.Logf("saved addresses=%d", len(list))
		if len(list) == 0 {
			t.Log("no saved addresses; skipping cover")
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/saved_address",
			map[string]string{"currency": "USDT"}, true)
		testutil.AssertCovers(t, "wallet/saved_address", raw, list)
	})

	t.Run("GetTransferOrderStatus", func(t *testing.T) {
		// Without a real tx_id the lookup is expected to fail; treat any error as
		// non-fatal so the endpoint path/signing are still exercised.
		st, err := c.NewGetTransferOrderStatusService("1").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/order_status", err) {
				return
			}
			t.Logf("order status (expected without a real tx_id): %v", err)
			return
		}
		t.Logf("order status=%+v", st)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/order_status",
			map[string]string{"tx_id": "1"}, true)
		testutil.AssertCovers(t, "wallet/order_status", raw, st)
	})

	t.Run("Transfer", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write ops disabled; set GATE_TEST_WRITE=1")
		}
		// Tiny spot->futures move, reversed immediately to stay balance-neutral.
		amount := decimal.NewFromFloat(0.01)
		out, err := c.NewTransferService("USDT", "spot", "futures", amount).SetSettle("usdt").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/transfers", err) {
				return
			}
			t.Fatalf("transfer: %v", err)
		}
		t.Logf("transfer tx_id=%d", out.TxID)
		back, err := c.NewTransferService("USDT", "futures", "spot", amount).SetSettle("usdt").Do(cx)
		if err != nil {
			t.Fatalf("reverse transfer: %v", err)
		}
		t.Logf("reverse tx_id=%d", back.TxID)
	})

	t.Run("TransferWithSubAccount", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write ops disabled; set GATE_TEST_WRITE=1")
		}
		// Requires a real sub-account; without one this is expected to fail.
		out, err := c.NewTransferWithSubAccountService("10000000", "USDT", decimal.NewFromFloat(0.01), "to").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_transfers(post)", err) {
				return
			}
			t.Logf("sub-account transfer (needs a real sub-account): %v", err)
			return
		}
		t.Logf("sub-account transfer tx_id=%d", out.TxID)
	})

	t.Run("SubAccountToSubAccount", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write ops disabled; set GATE_TEST_WRITE=1")
		}
		// Requires two real sub-accounts; without them this is expected to fail.
		out, err := c.NewSubAccountToSubAccountService(
			"USDT", "10000000", "spot", "10000001", "spot", decimal.NewFromFloat(0.01)).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/sub_account_to_sub_account", err) {
				return
			}
			t.Logf("sub-to-sub transfer (needs real sub-accounts): %v", err)
			return
		}
		t.Logf("sub-to-sub transfer tx_id=%d", out.TxID)
	})
}
