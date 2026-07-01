package options

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestOptionsMarket(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// underlyings
	underlyings, err := c.NewListOptionsUnderlyingsService().Do(cx)
	if err != nil {
		t.Fatalf("underlyings: %v", err)
	}
	if len(underlyings) == 0 {
		t.Fatal("no underlyings returned")
	}
	t.Logf("underlyings=%d first=%+v", len(underlyings), underlyings[0])
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/underlyings", nil, false)
	testutil.AssertCovers(t, "options/underlyings", raw, underlyings)

	// expirations
	exps, err := c.NewListOptionsExpirationsService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("expirations: %v", err)
	}
	t.Logf("expirations=%v", exps)

	// contracts
	contracts, err := c.NewListOptionsContractsService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("contracts: %v", err)
	}
	if len(contracts) == 0 {
		t.Skip("no BTC_USDT option contracts listed; skipping contract-specific checks")
	}
	t.Logf("contracts=%d first=%+v", len(contracts), contracts[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/contracts",
		map[string]string{"underlying": "BTC_USDT"}, false)
	testutil.AssertCovers(t, "options/contracts", raw, contracts)

	contract := contracts[0].Name

	// single contract
	one, err := c.NewGetOptionsContractService(contract).Do(cx)
	if err != nil {
		t.Fatalf("contract %s: %v", contract, err)
	}
	t.Logf("contract: %+v", one)
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/contracts/"+contract, nil, false)
	testutil.AssertCovers(t, "options/contracts/{contract}", raw, one)

	// settlements
	settlements, err := c.NewListOptionsSettlementsService("BTC_USDT").SetLimit(3).Do(cx)
	if err != nil {
		t.Fatalf("settlements: %v", err)
	}
	t.Logf("settlements=%d", len(settlements))
	if len(settlements) > 0 {
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/settlements",
			map[string]string{"underlying": "BTC_USDT", "limit": "3"}, false)
		testutil.AssertCovers(t, "options/settlements", raw, settlements)

		// single settlement, keyed by the record's own contract + time
		st := settlements[0]
		sOne, err := c.NewGetOptionsSettlementService(st.Contract, "BTC_USDT", st.Time).Do(cx)
		if err != nil {
			t.Fatalf("settlement %s: %v", st.Contract, err)
		}
		t.Logf("settlement: %+v", sOne)
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/settlements/"+st.Contract,
			map[string]string{"underlying": "BTC_USDT", "at": strconv.FormatInt(st.Time.Unix(), 10)}, false)
		testutil.AssertCovers(t, "options/settlements/{contract}", raw, sOne)
	}

	// order book
	ob, err := c.NewListOptionsOrderBookService(contract).SetLimit(2).SetWithID(true).Do(cx)
	if err != nil {
		t.Fatalf("order book: %v", err)
	}
	t.Logf("order book: id=%d asks=%d bids=%d", ob.ID, len(ob.Asks), len(ob.Bids))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/order_book",
		map[string]string{"contract": contract, "limit": "2", "with_id": "true"}, false)
	testutil.AssertCovers(t, "options/order_book", raw, ob)

	// tickers
	tickers, err := c.NewListOptionsTickersService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("tickers: %v", err)
	}
	if len(tickers) == 0 {
		t.Fatal("no tickers returned")
	}
	t.Logf("tickers=%d first=%+v", len(tickers), tickers[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/tickers",
		map[string]string{"underlying": "BTC_USDT"}, false)
	testutil.AssertCovers(t, "options/tickers", raw, tickers)

	// underlying ticker
	ut, err := c.NewListOptionsUnderlyingTickersService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("underlying ticker: %v", err)
	}
	t.Logf("underlying ticker: %+v", ut)
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/underlying/tickers/BTC_USDT", nil, false)
	testutil.AssertCovers(t, "options/underlying/tickers", raw, ut)

	// contract candlesticks
	candles, err := c.NewListOptionsCandlesticksService(contract).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("candlesticks: %v", err)
	}
	t.Logf("candlesticks=%d", len(candles))
	if len(candles) > 0 {
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/candlesticks",
			map[string]string{"contract": contract, "limit": "2"}, false)
		testutil.AssertCovers(t, "options/candlesticks", raw, candles)
	}

	// underlying candlesticks
	ucandles, err := c.NewListOptionsUnderlyingCandlesticksService("BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("underlying candlesticks: %v", err)
	}
	if len(ucandles) == 0 {
		t.Fatal("no underlying candlesticks returned")
	}
	t.Logf("underlying candlesticks=%d first=%+v", len(ucandles), ucandles[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/underlying/candlesticks",
		map[string]string{"underlying": "BTC_USDT", "limit": "2"}, false)
	testutil.AssertCovers(t, "options/underlying/candlesticks", raw, ucandles)

	// trades (filter by option type; a single contract may have no fills)
	trades, err := c.NewListOptionsTradesService().SetType("C").SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("trades: %v", err)
	}
	t.Logf("trades=%d", len(trades))
	if len(trades) > 0 {
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/options/trades",
			map[string]string{"type": "C", "limit": "2"}, false)
		testutil.AssertCovers(t, "options/trades", raw, trades)
	}
}
