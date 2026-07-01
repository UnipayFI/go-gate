package spot

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestSpotCurrency(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// ListCurrencies
	currencies, err := c.NewListCurrenciesService().Do(cx)
	if err != nil {
		t.Fatalf("list currencies: %v", err)
	}
	if len(currencies) == 0 {
		t.Fatal("no currencies returned")
	}
	t.Logf("currencies=%d first=%+v", len(currencies), currencies[0])
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/spot/currencies", nil, false)
	testutil.AssertCovers(t, "spot/currencies", raw, currencies)

	// GetCurrency
	cur, err := c.NewGetCurrencyService("USDT").Do(cx)
	if err != nil {
		t.Fatalf("get currency: %v", err)
	}
	t.Logf("currency: %+v", cur)
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/spot/currencies/USDT", nil, false)
	testutil.AssertCovers(t, "spot/currencies/USDT", raw, cur)

	// GetSpotInsuranceHistory
	from := time.Now().Add(-30 * 24 * time.Hour)
	to := time.Now()
	insurance, err := c.NewGetSpotInsuranceHistoryService("margin", "USDT", from, to).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("insurance history: %v", err)
	}
	t.Logf("insurance records=%d", len(insurance))
	if len(insurance) > 0 {
		t.Logf("insurance[0]: %+v", insurance[0])
		params := map[string]string{
			"business": "margin",
			"currency": "USDT",
			"from":     strconv.FormatInt(from.Unix(), 10),
			"to":       strconv.FormatInt(to.Unix(), 10),
			"limit":    "2",
		}
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/spot/insurance_history", params, false)
		testutil.AssertCovers(t, "spot/insurance_history", raw, insurance)
	}
}
