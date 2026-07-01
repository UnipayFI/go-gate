package spot

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestSpotDepth(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// Order book: request with_id so id/current/update are populated.
	ob, err := c.NewListOrderBookService("BTC_USDT").SetLimit(5).SetWithID(true).Do(cx)
	if err != nil {
		t.Fatalf("order book: %v", err)
	}
	t.Logf("order book id=%d current=%s asks=%d bids=%d", ob.ID, ob.Current, len(ob.Asks), len(ob.Bids))
	if len(ob.Asks) == 0 || len(ob.Bids) == 0 {
		t.Fatal("empty order book")
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/order_book",
		map[string]string{"currency_pair": "BTC_USDT", "limit": "5", "with_id": "true"}, false)
	testutil.AssertCovers(t, "spot/order_book", raw, ob)

	// Recent market trades.
	trades, err := c.NewListTradesService("BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("trades: %v", err)
	}
	if len(trades) == 0 {
		t.Fatal("no trades returned")
	}
	t.Logf("trade: %+v", trades[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/spot/trades",
		map[string]string{"currency_pair": "BTC_USDT", "limit": "2"}, false)
	testutil.AssertCovers(t, "spot/trades", raw, trades)

	// Candlesticks decode from array-of-arrays; just assert len>0 and log.
	candles, err := c.NewListCandlesticksService("BTC_USDT").SetInterval("1m").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("candlesticks: %v", err)
	}
	if len(candles) == 0 {
		t.Fatal("no candlesticks returned")
	}
	t.Logf("candle: ts=%s open=%s high=%s low=%s close=%s baseVol=%s quoteVol=%s closed=%v",
		candles[0].Timestamp, candles[0].Open, candles[0].High, candles[0].Low,
		candles[0].Close, candles[0].BaseVolume, candles[0].QuoteVolume, candles[0].WindowClosed)
}
