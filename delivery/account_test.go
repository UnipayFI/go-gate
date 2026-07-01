package delivery

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestDeliveryAccount(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// ListDeliveryAccounts -- private read; tolerable when the account has no
	// delivery balance/wallet.
	acct, err := c.NewListDeliveryAccountsService(SettleUSDT).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "delivery/accounts", err) {
			return
		}
		t.Fatalf("delivery accounts: %v", err)
	}
	t.Logf("account: currency=%s total=%s available=%s update_time=%s",
		acct.Currency, acct.Total, acct.Available, acct.UpdateTime)
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/accounts", nil, true)
	testutil.AssertCovers(t, "delivery/accounts", raw, acct)

	// ListDeliveryAccountBook -- private read; tolerable on empty history.
	params := map[string]string{"limit": "2"}
	book, err := c.NewListDeliveryAccountBookService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "delivery/account_book", err) {
			return
		}
		t.Fatalf("delivery account_book: %v", err)
	}
	t.Logf("account_book=%d", len(book))
	for _, b := range book {
		t.Logf("  time=%s type=%s change=%s balance=%s", b.Time, b.Type, b.Change, b.Balance)
	}
	if len(book) == 0 {
		t.Log("delivery account_book empty; nothing to diff")
		return
	}
	rawBook := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/account_book", params, true)
	testutil.AssertCovers(t, "delivery/account_book", rawBook, book)
}
