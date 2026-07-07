package account

import (
	"context"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestAccountExtra(t *testing.T) {
	// authed builds a signed, time-synced account client, skipping the subtest
	// when credentials are unset.
	authed := func(t *testing.T, cx context.Context) *AccountClient {
		t.Helper()
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		return c
	}

	t.Run("MainKeys", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		list, err := c.NewGetMainKeysService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "account/main_keys", err) {
				return
			}
			t.Fatalf("main_keys: %v", err)
		}
		t.Logf("main keys: %d", len(list))
		if len(list) == 0 {
			t.Skip("no main keys")
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/account/main_keys", nil, true)
		testutil.AssertCovers(t, "account/main_keys", raw, list)
	})
}
