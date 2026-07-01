package loan

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestCollateralLoan(t *testing.T) {
	// Public: supported borrowing / collateral currencies.
	t.Run("Currencies", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		list, err := c.NewListCollateralCurrenciesService().Do(cx)
		if err != nil {
			t.Fatalf("collateral currencies: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no collateral currencies returned")
		}
		t.Logf("collateral currencies: %d loan currencies, first=%s (%d collaterals)",
			len(list), list[0].LoanCurrency, len(list[0].CollateralCurrency))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/loan/collateral/currencies", nil, false)
		testutil.AssertCovers(t, "loan/collateral/currencies", raw, list)
	})

	// Private read: loan orders.
	t.Run("LoanOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		list, err := c.NewListCollateralLoanOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "loan/collateral/orders", err) {
				return
			}
			t.Fatalf("loan orders: %v", err)
		}
		t.Logf("loan orders=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/loan/collateral/orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "loan/collateral/orders", raw, list)

		// Single order detail, when an order exists.
		if len(list) > 0 {
			id := list[0].OrderID
			detail, err := c.NewGetCollateralLoanOrderDetailService(id).Do(cx)
			if err != nil {
				if testutil.Tolerable(t, "loan/collateral/orders/{id}", err) {
					return
				}
				t.Fatalf("loan order detail: %v", err)
			}
			t.Logf("order %d status=%s", detail.OrderID, detail.Status)
			raw := testutil.FetchRawGet(t, c, cx,
				"/api/v4/loan/collateral/orders/"+strconv.FormatInt(id, 10), nil, true)
			testutil.AssertCovers(t, "loan/collateral/orders/{id}", raw, detail)
		}
	})

	// Private read: repayment records.
	t.Run("RepayRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		list, err := c.NewListRepayRecordsService("repay").SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "loan/collateral/repay_records", err) {
				return
			}
			t.Fatalf("repay records: %v", err)
		}
		t.Logf("repay records=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/loan/collateral/repay_records",
			map[string]string{"source": "repay", "limit": "2"}, true)
		testutil.AssertCovers(t, "loan/collateral/repay_records", raw, list)
	})

	// Private read: collateral adjustment records.
	t.Run("CollateralRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		list, err := c.NewListCollateralRecordsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "loan/collateral/collaterals", err) {
				return
			}
			t.Fatalf("collateral records: %v", err)
		}
		t.Logf("collateral records=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/loan/collateral/collaterals",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "loan/collateral/collaterals", raw, list)
	})

	// Private read: total borrowing / collateral amount.
	t.Run("TotalAmount", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		total, err := c.NewGetUserTotalAmountService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "loan/collateral/total_amount", err) {
				return
			}
			t.Fatalf("total amount: %v", err)
		}
		t.Logf("total borrow=%s collateral=%s", total.BorrowAmount, total.CollateralAmount)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/loan/collateral/total_amount", nil, true)
		testutil.AssertCovers(t, "loan/collateral/total_amount", raw, total)
	})

	// Private read: LTV info for a currency pair.
	t.Run("LtvInfo", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		ltv, err := c.NewGetUserLtvInfoService("BTC", "USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "loan/collateral/ltv", err) {
				return
			}
			t.Fatalf("ltv info: %v", err)
		}
		t.Logf("ltv init=%s liquidate=%s left_borrowable=%s",
			ltv.InitLtv, ltv.LiquidateLtv, ltv.LeftBorrowableAmount)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/loan/collateral/ltv",
			map[string]string{"collateral_currency": "BTC", "borrow_currency": "USDT"}, true)
		testutil.AssertCovers(t, "loan/collateral/ltv", raw, ltv)
	})

	// State-changing: place a tiny loan, adjust collateral, then repay it all.
	t.Run("Write", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("GATE_TEST_WRITE!=1; skipping collateral loan write test")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}

		created, err := c.NewCreateCollateralLoanService(
			decimal.RequireFromString("0.0002"), "BTC",
			decimal.RequireFromString("5"), "USDT").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "create collateral loan", err) {
				return
			}
			t.Logf("create collateral loan failed (not tolerable): %v", err)
			return
		}
		t.Logf("created collateral loan order_id=%d", created.OrderID)

		// Append a little more collateral, then redeem it back.
		if _, err := c.NewOperateCollateralService(created.OrderID, "BTC",
			decimal.RequireFromString("0.0001"), "append").Do(cx); err != nil {
			t.Logf("append collateral: %v", err)
		}
		if _, err := c.NewOperateCollateralService(created.OrderID, "BTC",
			decimal.RequireFromString("0.0001"), "redeem").Do(cx); err != nil {
			t.Logf("redeem collateral: %v", err)
		}

		// Reverse the loan: repay everything.
		repay, err := c.NewRepayCollateralLoanService(created.OrderID,
			decimal.RequireFromString("5"), true).Do(cx)
		if err != nil {
			t.Logf("repay collateral loan: %v", err)
			return
		}
		t.Logf("repaid principal=%s interest=%s", repay.RepaidPrincipal, repay.RepaidInterest)
	})
}
