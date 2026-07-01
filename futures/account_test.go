package futures

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestFuturesAccount(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// GET /api/v4/futures/usdt/accounts
	acct, err := c.NewListFuturesAccountsService(SettleUSDT).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "futures/accounts", err) {
			t.Fatalf("futures accounts: %v", err)
		}
	} else {
		t.Logf("account: total=%s available=%s currency=%s marginMode=%d",
			acct.Total, acct.Available, acct.Currency, acct.MarginMode)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/accounts", nil, true)
		testutil.AssertCovers(t, "futures/accounts", raw, acct)
	}

	// GET /api/v4/futures/usdt/account_book
	book, err := c.NewListFuturesAccountBookService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "futures/account_book", err) {
			t.Fatalf("futures account_book: %v", err)
		}
	} else if len(book) == 0 {
		t.Logf("futures account_book: empty")
	} else {
		t.Logf("account_book[0]: type=%s change=%s balance=%s",
			book[0].Type, book[0].Change, book[0].Balance)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/account_book",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "futures/account_book", raw, book)
	}

	// GET /api/v4/futures/usdt/fee
	fees, err := c.NewGetFuturesFeeService(SettleUSDT).SetContract("BTC_USDT").Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "futures/fee", err) {
			t.Fatalf("futures fee: %v", err)
		}
	} else {
		for k, v := range fees {
			t.Logf("fee[%s]: taker=%s maker=%s", k, v.TakerFee, v.MakerFee)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/fee",
			map[string]string{"contract": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "futures/fee", raw, fees)
	}
}
