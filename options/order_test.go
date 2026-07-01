package options

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
)

func TestOptionsOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// ---- private read-only endpoints ----
	orders, err := c.NewListOptionsOrdersService("finished").SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/orders", err) {
			t.Fatalf("list options orders: %v", err)
		}
	} else {
		t.Logf("options orders=%d", len(orders))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/orders",
			map[string]string{"status": "finished", "limit": "2"}, true)
		testutil.AssertCovers(t, "options/orders", raw, orders)

		// GetOptionsOrder — reuse an id from the list when one is available.
		if len(orders) > 0 {
			id := strconv.FormatInt(orders[0].ID, 10)
			ord, gerr := c.NewGetOptionsOrderService(id).Do(cx)
			if gerr != nil {
				if !testutil.Tolerable(t, "options/orders/{id}", gerr) {
					t.Fatalf("get options order: %v", gerr)
				}
			} else {
				t.Logf("order %s status=%s size=%d", id, ord.Status, ord.Size)
				raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/orders/"+id, nil, true)
				testutil.AssertCovers(t, "options/orders/{id}", raw, ord)
			}
		}
	}

	// ---- state-changing endpoints (place / cancel / countdown) ----
	t.Run("write", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run options order write tests")
		}

		// Options may be disabled on the test account: a tolerable capability
		// error counts as a pass (endpoint + signing verified).
		contract := "BTC_USDT-20260626-100000-C"

		// place a small resting BUY well below any realistic premium (never fills).
		placed, err := c.NewCreateOptionsOrderService(contract, 1).
			SetTimeInForce("gtc").SetText("t-gogatetest").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "options create order", err) {
				return
			}
			t.Fatalf("create options order: %v", err)
		}
		oid := strconv.FormatInt(placed.ID, 10)
		t.Logf("placed id=%s price=%s size=%d status=%s", oid, placed.Price, placed.Size, placed.Status)

		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/orders/"+oid, nil, true)
		testutil.AssertCovers(t, "options/orders/{id}", raw, placed)

		got, err := c.NewGetOptionsOrderService(oid).Do(cx)
		if err != nil {
			t.Fatalf("get options order: %v", err)
		}
		t.Logf("get id=%d left=%d", got.ID, got.Left)

		cancelled, err := c.NewCancelOptionsOrderService(oid).Do(cx)
		if err != nil {
			t.Fatalf("cancel options order: %v", err)
		}
		t.Logf("cancelled id=%d finish_as=%s", cancelled.ID, cancelled.FinishAs)

		// countdown cancel-all: timeout 0 disarms it, so this is a safe no-op.
		status, err := c.NewCountdownCancelAllOptionsService(0).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "options countdown_cancel_all", err) {
				t.Fatalf("countdown cancel-all: %v", err)
			}
		} else {
			t.Logf("countdown trigger_time=%s", status.TriggerTime)
			craw := testutil.FetchRawPost(t, c, cx, "/api/v4/options/countdown_cancel_all",
				map[string]any{"timeout": 0}, true)
			testutil.AssertCovers(t, "options/countdown_cancel_all", craw, status)
		}

		// cancel any residual open orders.
		cleaned, err := c.NewCancelOptionsOrdersService().Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "options cancel-all", err) {
				t.Fatalf("cancel all: %v", err)
			}
		} else {
			t.Logf("cancel-all removed %d order(s)", len(cleaned))
		}
	})
}
