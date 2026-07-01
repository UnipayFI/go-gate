package spot

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

// testPublicClient builds an unauthenticated spot client for public endpoints.
func testPublicClient() *SpotClient {
	opts := []client.Options{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewSpotClient(opts...)
}

// testClient builds an authenticated spot client, skipping when creds are unset.
func testClient(t *testing.T) *SpotClient {
	t.Helper()
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.Options{client.WithAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewSpotClient(opts...)
}

func TestSpotServerTime(t *testing.T) {
	c := testPublicClient()
	resp, err := c.NewGetServerTimeService().Do(testutil.Ctx(t))
	if err != nil {
		t.Fatalf("server time: %v", err)
	}
	t.Logf("serverTime=%s (%d)", resp.ServerTime, resp.ServerTime.UnixMilli())
	if resp.ServerTime.IsZero() {
		t.Fatal("server time is zero")
	}
}

func TestSpotCurrencyPair(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	pair, err := c.NewGetCurrencyPairService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("currency pair: %v", err)
	}
	t.Logf("pair: %+v", pair)
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/currency_pairs/BTC_USDT", nil, false)
	testutil.AssertCovers(t, "spot/currency_pairs/BTC_USDT", raw, pair)
}

func TestSpotTickers(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	list, err := c.NewGetTickersService().SetCurrencyPair("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("tickers: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("no tickers returned")
	}
	t.Logf("ticker: %+v", list[0])
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/tickers",
		map[string]string{"currency_pair": "BTC_USDT"}, false)
	testutil.AssertCovers(t, "spot/tickers", raw, list)
}

func TestSpotAccounts(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	list, err := c.NewListSpotAccountsService().Do(cx)
	if err != nil {
		t.Fatalf("spot accounts: %v", err)
	}
	t.Logf("accounts=%d", len(list))
	for _, a := range list {
		t.Logf("  %s available=%s locked=%s", a.Currency, a.Available, a.Locked)
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/accounts", nil, true)
	testutil.AssertCovers(t, "spot/accounts", raw, list)
}
