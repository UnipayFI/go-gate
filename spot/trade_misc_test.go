package spot

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestSpotTradeMisc(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	t.Run("GetFee", func(t *testing.T) {
		fee, err := c.NewGetFeeService().SetCurrencyPair("BTC_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "spot/fee", err) {
				return
			}
			t.Fatalf("get fee: %v", err)
		}
		t.Logf("fee: %+v", fee)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/fee",
			map[string]string{"currency_pair": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "spot/fee", raw, fee)
	})

	t.Run("GetBatchSpotFee", func(t *testing.T) {
		fees, err := c.NewGetBatchSpotFeeService("BTC_USDT,ETH_USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "spot/batch_fee", err) {
				return
			}
			t.Fatalf("batch fee: %v", err)
		}
		t.Logf("batch_fee pairs=%d", len(fees))
		for pair, f := range fees {
			t.Logf("  %s taker=%s maker=%s", pair, f.TakerFee, f.MakerFee)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/batch_fee",
			map[string]string{"currency_pairs": "BTC_USDT,ETH_USDT"}, true)
		testutil.AssertCovers(t, "spot/batch_fee", raw, fees)
	})

	t.Run("ListSpotAccountBook", func(t *testing.T) {
		list, err := c.NewListSpotAccountBookService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "spot/account_book", err) {
				return
			}
			t.Fatalf("account book: %v", err)
		}
		t.Logf("account_book records=%d", len(list))
		for _, b := range list {
			t.Logf("  %s %s change=%s balance=%s type=%s", b.Time, b.Currency, b.Change, b.Balance, b.Type)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/account_book",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "spot/account_book", raw, list)
	})

	t.Run("ListMyTrades", func(t *testing.T) {
		list, err := c.NewListMyTradesService().SetCurrencyPair("BTC_USDT").SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "spot/my_trades", err) {
				return
			}
			t.Fatalf("my trades: %v", err)
		}
		t.Logf("my_trades=%d", len(list))
		for _, tr := range list {
			t.Logf("  %s %s %s amount=%s price=%s role=%s", tr.ID, tr.CurrencyPair, tr.Side, tr.Amount, tr.Price, tr.Role)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/my_trades",
			map[string]string{"currency_pair": "BTC_USDT", "limit": "2"}, true)
		testutil.AssertCovers(t, "spot/my_trades", raw, list)
	})
}
