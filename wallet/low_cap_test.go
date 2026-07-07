package wallet

import (
	"context"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestWalletExtra(t *testing.T) {
	// authed builds a signed, time-synced wallet client, skipping the subtest
	// when credentials are unset.
	authed := func(t *testing.T, cx context.Context) *WalletClient {
		t.Helper()
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		return c
	}

	t.Run("LowCapExchangeList", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		list, err := c.NewListLowCapExchangeService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "wallet/getLowCapExchangeList", err) {
				return
			}
			t.Fatalf("getLowCapExchangeList: %v", err)
		}
		t.Logf("low cap tokens: %d", len(list))
		if len(list) == 0 {
			t.Skip("no low cap tokens")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/wallet/getLowCapExchangeList", nil, true)
		testutil.AssertCovers(t, "wallet/getLowCapExchangeList", raw, list)
	})
}
