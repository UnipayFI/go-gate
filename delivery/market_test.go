package delivery

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestDeliveryMarket(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// Contracts: also the source of a live, expiry-suffixed contract name for the
	// contract-specific endpoints below.
	contracts, err := c.NewListDeliveryContractsService(SettleUSDT).Do(cx)
	if err != nil {
		t.Fatalf("list contracts: %v", err)
	}
	if len(contracts) == 0 {
		t.Skip("no delivery contracts listed; skipping market test")
	}
	t.Logf("contracts=%d first=%s", len(contracts), contracts[0].Name)
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/contracts", nil, false)
	testutil.AssertCovers(t, "delivery/contracts", raw, contracts)

	contract := contracts[0].Name

	// Single contract.
	one, err := c.NewGetDeliveryContractService(SettleUSDT, contract).Do(cx)
	if err != nil {
		t.Fatalf("get contract: %v", err)
	}
	t.Logf("contract %s: mark=%s expire=%s", one.Name, one.MarkPrice, one.ExpireTime)
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/contracts/"+contract, nil, false)
	testutil.AssertCovers(t, "delivery/contracts/{contract}", raw, one)

	// Order book (with_id to populate the id field).
	ob, err := c.NewListDeliveryOrderBookService(SettleUSDT, contract).SetLimit(2).SetWithID(true).Do(cx)
	if err != nil {
		t.Fatalf("order book: %v", err)
	}
	t.Logf("order book id=%d asks=%d bids=%d", ob.ID, len(ob.Asks), len(ob.Bids))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/order_book",
		map[string]string{"contract": contract, "limit": "2", "with_id": "true"}, false)
	testutil.AssertCovers(t, "delivery/order_book", raw, ob)

	// Trades.
	trades, err := c.NewListDeliveryTradesService(SettleUSDT, contract).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("trades: %v", err)
	}
	t.Logf("trades=%d", len(trades))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/trades",
		map[string]string{"contract": contract, "limit": "2"}, false)
	testutil.AssertCovers(t, "delivery/trades", raw, trades)

	// Candlesticks.
	candles, err := c.NewListDeliveryCandlesticksService(SettleUSDT, contract).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("candlesticks: %v", err)
	}
	t.Logf("candles=%d", len(candles))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/candlesticks",
		map[string]string{"contract": contract, "limit": "2"}, false)
	testutil.AssertCovers(t, "delivery/candlesticks", raw, candles)

	// Tickers (narrowed to the one contract).
	tickers, err := c.NewListDeliveryTickersService(SettleUSDT).SetContract(contract).Do(cx)
	if err != nil {
		t.Fatalf("tickers: %v", err)
	}
	if len(tickers) == 0 {
		t.Fatal("no tickers returned")
	}
	t.Logf("ticker: %+v", tickers[0])
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/tickers",
		map[string]string{"contract": contract}, false)
	testutil.AssertCovers(t, "delivery/tickers", raw, tickers)

	// Insurance ledger.
	insurance, err := c.NewListDeliveryInsuranceLedgerService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("insurance: %v", err)
	}
	t.Logf("insurance=%d", len(insurance))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/insurance",
		map[string]string{"limit": "2"}, false)
	testutil.AssertCovers(t, "delivery/insurance", raw, insurance)
}
