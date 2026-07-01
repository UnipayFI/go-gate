package loan

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestMultiCollateralLoan(t *testing.T) {
	cx := testutil.Ctx(t)

	// ----- PUBLIC endpoints -----
	pub := testPublicClient()

	// GET /api/v4/loan/multi_collateral/currencies
	currencies, err := pub.NewListMultiCollateralCurrenciesService().Do(cx)
	if err != nil {
		t.Fatalf("multi currencies: %v", err)
	}
	t.Logf("loan_currencies=%d collateral_currencies=%d", len(currencies.LoanCurrencies), len(currencies.CollateralCurrencies))
	raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/loan/multi_collateral/currencies", nil, false)
	testutil.AssertCovers(t, "loan/multi_collateral/currencies", raw, currencies)

	// GET /api/v4/loan/multi_collateral/ltv
	ltv, err := pub.NewGetMultiCollateralLtvService().Do(cx)
	if err != nil {
		t.Fatalf("multi ltv: %v", err)
	}
	t.Logf("ltv init=%s alert=%s liquidate=%s", ltv.InitLTV, ltv.AlertLTV, ltv.LiquidateLTV)
	raw = testutil.FetchRawGet(t, pub, cx, "/api/v4/loan/multi_collateral/ltv", nil, false)
	testutil.AssertCovers(t, "loan/multi_collateral/ltv", raw, ltv)

	// GET /api/v4/loan/multi_collateral/fixed_rate
	fixRates, err := pub.NewGetMultiCollateralFixRateService().Do(cx)
	if err != nil {
		t.Fatalf("multi fixed_rate: %v", err)
	}
	if len(fixRates) == 0 {
		t.Fatal("no fixed rates returned")
	}
	t.Logf("fixed_rate[0]: %+v", fixRates[0])
	raw = testutil.FetchRawGet(t, pub, cx, "/api/v4/loan/multi_collateral/fixed_rate", nil, false)
	testutil.AssertCovers(t, "loan/multi_collateral/fixed_rate", raw, fixRates)

	// GET /api/v4/loan/multi_collateral/current_rate
	currentRates, err := pub.NewGetMultiCollateralCurrentRateService([]string{"USDT"}).Do(cx)
	if err != nil {
		t.Fatalf("multi current_rate: %v", err)
	}
	if len(currentRates) == 0 {
		t.Fatal("no current rates returned")
	}
	t.Logf("current_rate[0]: %+v", currentRates[0])
	raw = testutil.FetchRawGet(t, pub, cx, "/api/v4/loan/multi_collateral/current_rate",
		map[string]string{"currencies": "USDT"}, false)
	testutil.AssertCovers(t, "loan/multi_collateral/current_rate", raw, currentRates)

	// ----- PRIVATE reads -----
	c := testClient(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// GET /api/v4/loan/multi_collateral/orders
	orders, err := c.NewListMultiCollateralOrdersService().SetLimit(2).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "loan/multi_collateral/orders", err) {
			return
		}
		t.Fatalf("multi orders: %v", err)
	}
	t.Logf("orders=%d", len(orders))
	raw = testutil.FetchRawGet(t, c, cx, "/api/v4/loan/multi_collateral/orders",
		map[string]string{"limit": "2"}, true)
	testutil.AssertCovers(t, "loan/multi_collateral/orders", raw, orders)

	if len(orders) > 0 {
		// GET /api/v4/loan/multi_collateral/orders/{order_id}
		id := orders[0].OrderID
		detail, err := c.NewGetMultiCollateralOrderDetailService(id).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "loan/multi_collateral/orders/{id}", err) {
				t.Fatalf("multi order detail: %v", err)
			}
		} else {
			t.Logf("order detail: %+v", detail)
			raw = testutil.FetchRawGet(t, c, cx, "/api/v4/loan/multi_collateral/orders/"+id, nil, true)
			testutil.AssertCovers(t, "loan/multi_collateral/orders/{id}", raw, detail)
		}
	}

	// GET /api/v4/loan/multi_collateral/repay
	repayRecords, err := c.NewListMultiRepayRecordsService("repay").SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "loan/multi_collateral/repay", err) {
			t.Fatalf("multi repay records: %v", err)
		}
	} else {
		t.Logf("repay records=%d", len(repayRecords))
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/loan/multi_collateral/repay",
			map[string]string{"type": "repay", "limit": "2"}, true)
		testutil.AssertCovers(t, "loan/multi_collateral/repay", raw, repayRecords)
	}

	// GET /api/v4/loan/multi_collateral/mortgage
	mortgageRecords, err := c.NewListMultiCollateralRecordsService().SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "loan/multi_collateral/mortgage", err) {
			t.Fatalf("multi mortgage records: %v", err)
		}
	} else {
		t.Logf("mortgage records=%d", len(mortgageRecords))
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/loan/multi_collateral/mortgage",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "loan/multi_collateral/mortgage", raw, mortgageRecords)
	}

	// GET /api/v4/loan/multi_collateral/currency_quota
	quotas, err := c.NewListUserCurrencyQuotaService("collateral", "BTC").Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "loan/multi_collateral/currency_quota", err) {
			t.Fatalf("multi currency quota: %v", err)
		}
	} else {
		t.Logf("quotas=%d", len(quotas))
		raw = testutil.FetchRawGet(t, c, cx, "/api/v4/loan/multi_collateral/currency_quota",
			map[string]string{"type": "collateral", "currency": "BTC"}, true)
		testutil.AssertCovers(t, "loan/multi_collateral/currency_quota", raw, quotas)
	}

	// ----- STATE-CHANGING endpoints (guarded) -----
	if !testutil.WriteEnabled() {
		t.Skip("write tests disabled (set GATE_TEST_WRITE=1); create/repay/operate skipped")
	}

	// POST /api/v4/loan/multi_collateral/orders — borrow a tiny amount of USDT
	// against BTC collateral, then unwind via query + repay.
	created, err := c.NewCreateMultiCollateralService(
		"USDT",
		decimal.RequireFromString("5"),
		[]MultiCollateralInput{{Currency: "BTC", Amount: decimal.RequireFromString("0.0002")}},
	).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "loan/multi_collateral/orders create", err) {
			t.Fatalf("create multi collateral: %v", err)
		}
		return
	}
	t.Logf("created order_id=%d", created.OrderID)

	// POST /api/v4/loan/multi_collateral/mortgage — append a little more collateral.
	adjusted, err := c.NewOperateMultiCollateralService(created.OrderID, "append").
		SetCollaterals([]MultiCollateralInput{{Currency: "BTC", Amount: decimal.RequireFromString("0.0001")}}).
		Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "loan/multi_collateral/mortgage operate", err) {
			t.Fatalf("operate multi collateral: %v", err)
		}
	} else {
		t.Logf("adjust result: %+v", adjusted)
	}

	// POST /api/v4/loan/multi_collateral/repay — repay the borrow in full.
	repaid, err := c.NewRepayMultiCollateralLoanService(created.OrderID,
		[]MultiRepayItemInput{{Currency: "USDT", RepaidAll: true}}).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "loan/multi_collateral/repay post", err) {
			t.Fatalf("repay multi collateral: %v", err)
		}
	} else {
		t.Logf("repaid: %+v", repaid)
	}
}
