package options

import (
	"context"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestOptionsExtra(t *testing.T) {
	// authed builds a signed, time-synced options client, skipping the subtest
	// when credentials are unset.
	authed := func(t *testing.T, cx context.Context) *OptionsClient {
		t.Helper()
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		return c
	}

	// --- State-changing endpoints (guarded) ---

	t.Run("AmendOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run state-changing AmendOptionsOrder")
		}
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		// Amend a non-existent order id: Gate rejects it, so we never touch a real
		// order.
		_, err := c.NewAmendOptionsOrderService(1, "BTC_USDT-20260101-50000-C", decimal.RequireFromString("0.1"), 1).Do(cx)
		if err != nil {
			t.Logf("amend options order: %v (tolerable)", err)
			return
		}
		t.Log("amend options order accepted")
	})
}
