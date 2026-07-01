package earn

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestEarnUni(t *testing.T) {
	pub := testPublicClient()
	cx := testutil.Ctx(t)

	// ---- public endpoints ----

	t.Run("ListUniCurrencies", func(t *testing.T) {
		list, err := pub.NewListUniCurrenciesService().Do(cx)
		if err != nil {
			t.Fatalf("list uni currencies: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no uni currencies returned")
		}
		t.Logf("currencies=%d first=%+v", len(list), list[0])
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/earn/uni/currencies", nil, false)
		testutil.AssertCovers(t, "earn/uni/currencies", raw, list)
	})

	t.Run("GetUniCurrency", func(t *testing.T) {
		cur, err := pub.NewGetUniCurrencyService("USDT").Do(cx)
		if err != nil {
			t.Fatalf("get uni currency: %v", err)
		}
		t.Logf("USDT: %+v", cur)
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/earn/uni/currencies/USDT", nil, false)
		testutil.AssertCovers(t, "earn/uni/currencies/USDT", raw, cur)
	})

	t.Run("ListUniChart", func(t *testing.T) {
		to := time.Now()
		from := to.Add(-24 * time.Hour)
		pts, err := pub.NewListUniChartService("USDT", from, to).Do(cx)
		if err != nil {
			t.Fatalf("list uni chart: %v", err)
		}
		t.Logf("chart points=%d", len(pts))
		if len(pts) > 0 {
			t.Logf("first point: %+v", pts[0])
		}
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/earn/uni/chart", map[string]string{
			"asset": "USDT",
			"from":  strconv.FormatInt(from.Unix(), 10),
			"to":    strconv.FormatInt(to.Unix(), 10),
		}, false)
		testutil.AssertCovers(t, "earn/uni/chart", raw, pts)
	})

	t.Run("ListUniRate", func(t *testing.T) {
		rates, err := pub.NewListUniRateService().Do(cx)
		if err != nil {
			t.Fatalf("list uni rate: %v", err)
		}
		if len(rates) == 0 {
			t.Fatal("no uni rates returned")
		}
		t.Logf("rates=%d first=%+v", len(rates), rates[0])
		raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/earn/uni/rate", nil, false)
		testutil.AssertCovers(t, "earn/uni/rate", raw, rates)
	})

	// ---- private endpoints ----

	t.Run("private", func(t *testing.T) {
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}

		t.Run("ListUserUniLends", func(t *testing.T) {
			list, err := c.NewListUserUniLendsService().SetCurrency("USDT").SetLimit(10).Do(cx)
			if err != nil {
				if testutil.Tolerable(t, "earn/uni/lends", err) {
					return
				}
				t.Fatalf("list user uni lends: %v", err)
			}
			t.Logf("lends=%d", len(list))
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/uni/lends",
				map[string]string{"currency": "USDT", "limit": "10"}, true)
			testutil.AssertCovers(t, "earn/uni/lends", raw, list)
		})

		t.Run("ListUniLendRecords", func(t *testing.T) {
			list, err := c.NewListUniLendRecordsService().SetCurrency("USDT").SetLimit(10).Do(cx)
			if err != nil {
				if testutil.Tolerable(t, "earn/uni/lend_records", err) {
					return
				}
				t.Fatalf("list uni lend records: %v", err)
			}
			t.Logf("lend records=%d", len(list))
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/uni/lend_records",
				map[string]string{"currency": "USDT", "limit": "10"}, true)
			testutil.AssertCovers(t, "earn/uni/lend_records", raw, list)
		})

		t.Run("GetUniInterest", func(t *testing.T) {
			interest, err := c.NewGetUniInterestService("USDT").Do(cx)
			if err != nil {
				if testutil.Tolerable(t, "earn/uni/interests/USDT", err) {
					return
				}
				t.Fatalf("get uni interest: %v", err)
			}
			t.Logf("interest: %+v", interest)
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/uni/interests/USDT", nil, true)
			testutil.AssertCovers(t, "earn/uni/interests/USDT", raw, interest)
		})

		t.Run("ListUniInterestRecords", func(t *testing.T) {
			list, err := c.NewListUniInterestRecordsService().SetCurrency("USDT").SetLimit(10).Do(cx)
			if err != nil {
				if testutil.Tolerable(t, "earn/uni/interest_records", err) {
					return
				}
				t.Fatalf("list uni interest records: %v", err)
			}
			t.Logf("interest records=%d", len(list))
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/uni/interest_records",
				map[string]string{"currency": "USDT", "limit": "10"}, true)
			testutil.AssertCovers(t, "earn/uni/interest_records", raw, list)
		})

		t.Run("GetUniInterestStatus", func(t *testing.T) {
			status, err := c.NewGetUniInterestStatusService("USDT").Do(cx)
			if err != nil {
				if testutil.Tolerable(t, "earn/uni/interest_status/USDT", err) {
					return
				}
				t.Fatalf("get uni interest status: %v", err)
			}
			t.Logf("interest status: %+v", status)
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/earn/uni/interest_status/USDT", nil, true)
			testutil.AssertCovers(t, "earn/uni/interest_status/USDT", raw, status)
		})

		t.Run("CreateAndChangeUniLend", func(t *testing.T) {
			if !testutil.WriteEnabled() {
				t.Skip("set GATE_TEST_WRITE=1 to exercise uni lend create/amend/redeem")
			}
			// Lend the minimum (1 USDT) at the pool's minimum hourly rate.
			minRate := decimal.RequireFromString("0.00000011")
			if err := c.NewCreateUniLendService("USDT", decimal.RequireFromString("1"), "lend").
				SetMinRate(minRate).Do(cx); err != nil {
				if testutil.Tolerable(t, "earn/uni/lends lend", err) {
					return
				}
				t.Fatalf("create uni lend: %v", err)
			}
			// Amend the minimum rate on the open order.
			if err := c.NewChangeUniLendService("USDT", decimal.RequireFromString("0.00000012")).Do(cx); err != nil {
				if !testutil.Tolerable(t, "earn/uni/lends amend", err) {
					t.Errorf("change uni lend: %v", err)
				}
			}
			// Redeem to reverse the lend.
			if err := c.NewCreateUniLendService("USDT", decimal.RequireFromString("1"), "redeem").Do(cx); err != nil {
				if !testutil.Tolerable(t, "earn/uni/lends redeem", err) {
					t.Errorf("redeem uni lend: %v", err)
				}
			}
		})
	})
}
