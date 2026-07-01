package futures

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestFuturesMarket(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// ListFuturesContracts
	contracts, err := c.NewListFuturesContractsService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("list contracts: %v", err)
	}
	if len(contracts) == 0 {
		t.Fatal("no contracts returned")
	}
	t.Logf("contracts=%d first=%+v", len(contracts), contracts[0])
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/contracts",
		map[string]string{"limit": "2"}, false)
	testutil.AssertCovers(t, "futures/contracts", raw, contracts)

	// GetFuturesContract
	contract, err := c.NewGetFuturesContractService(SettleUSDT, "BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("get contract: %v", err)
	}
	t.Logf("contract: %+v", contract)
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/contracts/BTC_USDT", nil, false)
	testutil.AssertCovers(t, "futures/contracts/BTC_USDT", raw, contract)

	// ListFuturesOrderBook
	ob, err := c.NewListFuturesOrderBookService(SettleUSDT, "BTC_USDT").SetLimit(2).SetWithID(true).Do(cx)
	if err != nil {
		t.Fatalf("order book: %v", err)
	}
	t.Logf("order book: id=%d asks=%d bids=%d", ob.ID, len(ob.Asks), len(ob.Bids))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/order_book",
		map[string]string{"contract": "BTC_USDT", "limit": "2", "with_id": "true"}, false)
	testutil.AssertCovers(t, "futures/order_book", raw, ob)

	// ListFuturesTrades
	trades, err := c.NewListFuturesTradesService(SettleUSDT, "BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("trades: %v", err)
	}
	if len(trades) == 0 {
		t.Fatal("no trades returned")
	}
	t.Logf("trades=%d first=%+v", len(trades), trades[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/trades",
		map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
	testutil.AssertCovers(t, "futures/trades", raw, trades)

	// ListFuturesCandlesticks
	candles, err := c.NewListFuturesCandlesticksService(SettleUSDT, "BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("candlesticks: %v", err)
	}
	if len(candles) == 0 {
		t.Fatal("no candlesticks returned")
	}
	t.Logf("candles=%d first=%+v", len(candles), candles[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/candlesticks",
		map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
	testutil.AssertCovers(t, "futures/candlesticks", raw, candles)

	// ListFuturesPremiumIndex
	premium, err := c.NewListFuturesPremiumIndexService(SettleUSDT, "BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("premium index: %v", err)
	}
	if len(premium) == 0 {
		t.Fatal("no premium index returned")
	}
	t.Logf("premium=%d first=%+v", len(premium), premium[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/premium_index",
		map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
	testutil.AssertCovers(t, "futures/premium_index", raw, premium)

	// ListFuturesTickers
	tickers, err := c.NewListFuturesTickersService(SettleUSDT).SetContract("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("tickers: %v", err)
	}
	if len(tickers) == 0 {
		t.Fatal("no tickers returned")
	}
	t.Logf("ticker: %+v", tickers[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/tickers",
		map[string]string{"contract": "BTC_USDT"}, false)
	testutil.AssertCovers(t, "futures/tickers", raw, tickers)

	// ListFuturesFundingRateHistory
	funding, err := c.NewListFuturesFundingRateHistoryService(SettleUSDT, "BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("funding rate: %v", err)
	}
	if len(funding) == 0 {
		t.Fatal("no funding rate returned")
	}
	t.Logf("funding=%d first=%+v", len(funding), funding[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/funding_rate",
		map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
	testutil.AssertCovers(t, "futures/funding_rate", raw, funding)
}
