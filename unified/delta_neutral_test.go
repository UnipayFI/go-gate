package unified

import (
	"context"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestUnifiedExtra(t *testing.T) {
	// authed builds a signed, time-synced unified client, skipping the subtest
	// when credentials are unset.
	authed := func(t *testing.T, cx context.Context) *UnifiedClient {
		t.Helper()
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		return c
	}

	// --- Private reads ---

	t.Run("DeltaNeutral", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		resp, err := c.NewGetUnifiedDeltaNeutralService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/delta_neutral", err) {
				return
			}
			t.Fatalf("delta_neutral: %v", err)
		}
		t.Logf("delta neutral: enabled=%v", resp.Enabled)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/delta_neutral", nil, true)
		testutil.AssertCovers(t, "unified/delta_neutral", raw, resp)
	})

	t.Run("EstimatedQuickRepayment", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		resp, err := c.NewGetEstimatedQuickRepaymentService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/estimated_quick_repayment", err) {
				return
			}
			t.Fatalf("estimated_quick_repayment: %v", err)
		}
		t.Logf("estimated quick repayment: debts=%d", len(resp.DebtCurrencies))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/unified/estimated_quick_repayment", nil, true)
		testutil.AssertCovers(t, "unified/estimated_quick_repayment", raw, resp)
	})

	// --- State-changing endpoints (guarded) ---

	t.Run("SetDeltaNeutral", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run state-changing SetUnifiedDeltaNeutral")
		}
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		cur, err := c.NewGetUnifiedDeltaNeutralService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/delta_neutral set(read)", err) {
				return
			}
			t.Fatalf("read delta_neutral: %v", err)
		}
		// Re-apply the current setting so the write is a no-op / reversible.
		resp, err := c.NewSetUnifiedDeltaNeutralService(cur.Enabled).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "unified/delta_neutral set", err) {
				return
			}
			t.Fatalf("set delta_neutral: %v", err)
		}
		t.Logf("re-applied delta neutral enabled=%v", resp.Enabled)
	})

	t.Run("QuickRepayment", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run state-changing CreateQuickRepayment")
		}
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		// Repay USDT liability using BTC; with no matching liability Gate rejects
		// the request, so nothing is actually repaid.
		_, err := c.NewCreateQuickRepaymentService([]string{"USDT"}, []string{"BTC"}).Do(cx)
		if err != nil {
			t.Logf("quick repayment: %v (tolerable)", err)
			return
		}
		t.Log("quick repayment accepted")
	})
}
