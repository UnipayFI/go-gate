package unified

import (
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestUnifiedLoan(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	t.Run("ListUnifiedLoans", func(t *testing.T) {
		list, err := c.NewListUnifiedLoansService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/loans", err) {
				return
			}
			t.Fatalf("list unified loans: %v", err)
		}
		t.Logf("loans=%d", len(list))
		for _, l := range list {
			t.Logf("  %s amount=%s type=%s", l.Currency, l.Amount, l.Type)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/loans",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "unified/loans", raw, list)
	})

	t.Run("ListUnifiedLoanRecords", func(t *testing.T) {
		list, err := c.NewListUnifiedLoanRecordsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/loan_records", err) {
				return
			}
			t.Fatalf("list unified loan records: %v", err)
		}
		t.Logf("loanRecords=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/loan_records",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "unified/loan_records", raw, list)
	})

	t.Run("ListUnifiedLoanInterestRecords", func(t *testing.T) {
		list, err := c.NewListUnifiedLoanInterestRecordsService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/interest_records", err) {
				return
			}
			t.Fatalf("list unified interest records: %v", err)
		}
		t.Logf("interestRecords=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/interest_records",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "unified/interest_records", raw, list)
	})

	t.Run("CreateUnifiedLoan", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write tests disabled; set GATE_TEST_WRITE=1")
		}
		// Borrow a tiny amount, then repay it in full so the operation reverses.
		res, err := c.NewCreateUnifiedLoanService("USDT", decimal.NewFromInt(1), "borrow").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/loans borrow", err) {
				return
			}
			t.Fatalf("borrow unified loan: %v", err)
		}
		t.Logf("borrow tran_id=%d", res.TranID)

		repay, err := c.NewCreateUnifiedLoanService("USDT", decimal.NewFromInt(1), "repay").
			SetRepaidAll(true).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/loans repay", err) {
				return
			}
			t.Fatalf("repay unified loan: %v", err)
		}
		t.Logf("repay tran_id=%d", repay.TranID)
	})
}
