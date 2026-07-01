package margin

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestMarginUni(t *testing.T) {
	pub := testPublicClient()
	cx := testutil.Ctx(t)

	// ListUniCurrencyPairs (public).
	pairs, err := pub.NewListUniCurrencyPairsService().Do(cx)
	if err != nil {
		t.Fatalf("list uni currency pairs: %v", err)
	}
	if len(pairs) == 0 {
		t.Fatal("no uni currency pairs returned")
	}
	t.Logf("uni currency pairs=%d first=%+v", len(pairs), pairs[0])
	raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/margin/uni/currency_pairs", nil, false)
	testutil.AssertCovers(t, "margin/uni/currency_pairs", raw, pairs)

	// GetUniCurrencyPair (public).
	pair, err := pub.NewGetUniCurrencyPairService("BTC_USDT").Do(cx)
	if err != nil {
		t.Fatalf("get uni currency pair: %v", err)
	}
	t.Logf("uni currency pair: %+v", pair)
	raw = testutil.FetchRawGet(t, pub, cx, "/api/v4/margin/uni/currency_pairs/BTC_USDT", nil, false)
	testutil.AssertCovers(t, "margin/uni/currency_pairs/BTC_USDT", raw, pair)

	// Remaining endpoints are signed. estimate_rate requires a Timestamp header
	// on the live API (despite the docs saying otherwise), so it is signed too.
	c := testClient(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// GetMarginUniEstimateRate (signed).
	rates, err := c.NewGetMarginUniEstimateRateService([]string{"BTC", "USDT"}).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "margin/uni/estimate_rate", err) {
			t.Fatalf("estimate rate: %v", err)
		}
	} else {
		t.Logf("estimate rates: %+v", rates)
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/margin/uni/estimate_rate",
			map[string]string{"currencies": "BTC,USDT"}, true)
		testutil.AssertCovers(t, "margin/uni/estimate_rate", raw, rates)
	}

	// ListUniLoans (signed read).
	loans, err := c.NewListUniLoansService().SetCurrencyPair("BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "margin/uni/loans", err) {
			t.Fatalf("list uni loans: %v", err)
		}
	} else {
		t.Logf("uni loans=%d", len(loans))
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/margin/uni/loans",
			map[string]string{"currency_pair": "BTC_USDT", "limit": "2"}, true)
		testutil.AssertCovers(t, "margin/uni/loans", raw, loans)
	}

	// ListUniLoanRecords (signed read).
	records, err := c.NewListUniLoanRecordsService().SetCurrencyPair("BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "margin/uni/loan_records", err) {
			t.Fatalf("list uni loan records: %v", err)
		}
	} else {
		t.Logf("uni loan records=%d", len(records))
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/margin/uni/loan_records",
			map[string]string{"currency_pair": "BTC_USDT", "limit": "2"}, true)
		testutil.AssertCovers(t, "margin/uni/loan_records", raw, records)
	}

	// ListUniLoanInterestRecords (signed read).
	interest, err := c.NewListUniLoanInterestRecordsService().SetCurrencyPair("BTC_USDT").SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "margin/uni/interest_records", err) {
			t.Fatalf("list uni interest records: %v", err)
		}
	} else {
		t.Logf("uni interest records=%d", len(interest))
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/margin/uni/interest_records",
			map[string]string{"currency_pair": "BTC_USDT", "limit": "2"}, true)
		testutil.AssertCovers(t, "margin/uni/interest_records", raw, interest)
	}

	// GetUniBorrowable (signed read).
	borrowable, err := c.NewGetUniBorrowableService("USDT", "BTC_USDT").Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "margin/uni/borrowable", err) {
			t.Fatalf("get uni borrowable: %v", err)
		}
	} else {
		t.Logf("uni borrowable: %+v", borrowable)
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/margin/uni/borrowable",
			map[string]string{"currency": "USDT", "currency_pair": "BTC_USDT"}, true)
		testutil.AssertCovers(t, "margin/uni/borrowable", raw, borrowable)
	}

	// CreateUniLoan (state-changing): a tiny reversible repay, opt-in only.
	t.Run("CreateUniLoan", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("GATE_TEST_WRITE not set; skipping state-changing uni loan test")
		}
		err := c.NewCreateUniLoanService("BTC_USDT", "USDT", decimal.NewFromInt(1), "repay").Do(cx)
		if err != nil && !testutil.Tolerable(t, "margin/uni/loans create", err) {
			t.Fatalf("create uni loan: %v", err)
		}
	})
}
