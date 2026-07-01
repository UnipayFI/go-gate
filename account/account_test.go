package account

import (
	"fmt"
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestAccount(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	t.Run("detail", func(t *testing.T) {
		detail, err := c.NewGetAccountDetailService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "account/detail", err) {
				return
			}
			t.Fatalf("account detail: %v", err)
		}
		t.Logf("detail: user_id=%d tier=%d vip_tier=%d role=%d mode=%d",
			detail.UserID, detail.Tier, detail.VIPTier, detail.CopyTradingRole, detail.Key.Mode)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/account/detail", nil, true)
		testutil.AssertCovers(t, "account/detail", raw, detail)
	})

	t.Run("rate_limit", func(t *testing.T) {
		limits, err := c.NewGetAccountRateLimitService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "account/rate_limit", err) {
				return
			}
			t.Fatalf("account rate limit: %v", err)
		}
		t.Logf("rate limit tiers=%d", len(limits))
		for _, l := range limits {
			t.Logf("  tier=%s ratio=%s main_ratio=%s updated_at=%s", l.Tier, l.Ratio, l.MainRatio, l.UpdatedAt)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/account/rate_limit", nil, true)
		testutil.AssertCovers(t, "account/rate_limit", raw, limits)
	})

	t.Run("stp_groups", func(t *testing.T) {
		groups, err := c.NewListSTPGroupsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "account/stp_groups", err) {
				return
			}
			t.Fatalf("list stp groups: %v", err)
		}
		t.Logf("stp groups=%d", len(groups))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/account/stp_groups", nil, true)
		testutil.AssertCovers(t, "account/stp_groups", raw, groups)

		// Follow through into the group's user list when one exists.
		if len(groups) == 0 {
			t.Log("no stp groups; skipping stp_groups/{id}/users read")
			return
		}
		id := groups[0].ID
		users, err := c.NewListSTPGroupsUsersService(id).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "account/stp_groups/users", err) {
				return
			}
			t.Fatalf("list stp group users: %v", err)
		}
		t.Logf("stp group %d users=%d", id, len(users))
		path := fmt.Sprintf("/api/v4/account/stp_groups/%d/users", id)
		raw = testutil.FetchRawGet(t, c, cx, path, nil, true)
		testutil.AssertCovers(t, "account/stp_groups/users", raw, users)
	})

	t.Run("debit_fee", func(t *testing.T) {
		fee, err := c.NewGetDebitFeeService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "account/debit_fee", err) {
				return
			}
			t.Fatalf("get debit fee: %v", err)
		}
		t.Logf("debit_fee=%d", fee.DebitFee)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/account/debit_fee", nil, true)
		testutil.AssertCovers(t, "account/debit_fee", raw, fee)
	})

	t.Run("stp_group_lifecycle", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run STP group mutations")
		}
		name := fmt.Sprintf("gogate-%d", time.Now().Unix())
		grp, err := c.NewCreateSTPGroupService(name).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "create stp group", err) {
				return
			}
			t.Fatalf("create stp group: %v", err)
		}
		t.Logf("created stp group id=%d name=%s", grp.ID, grp.Name)

		// Add the main account itself, then remove it so the change is reversible.
		detail, err := c.NewGetAccountDetailService().Do(cx)
		if err != nil {
			t.Logf("detail for self user id: %v", err)
			return
		}
		added, err := c.NewAddSTPGroupUsersService(grp.ID, []int64{detail.UserID}).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "add stp user", err) {
				t.Errorf("add stp user: %v", err)
			}
			return
		}
		t.Logf("added, users=%d", len(added))

		removed, err := c.NewDeleteSTPGroupUsersService(grp.ID, detail.UserID).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "delete stp user", err) {
				t.Errorf("delete stp user: %v", err)
			}
			return
		}
		t.Logf("removed, remaining users=%d", len(removed))
	})

	t.Run("set_debit_fee", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to run debit-fee mutation")
		}
		cur, err := c.NewGetDebitFeeService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "get debit fee", err) {
				return
			}
			t.Fatalf("get debit fee: %v", err)
		}
		// Re-set the current value so the mutation is a no-op / reversible.
		if err := c.NewSetDebitFeeService(cur.DebitFee).Do(cx); err != nil {
			if !testutil.Tolerable(t, "set debit fee", err) {
				t.Errorf("set debit fee: %v", err)
			}
			return
		}
		t.Logf("re-applied debit_fee=%d", cur.DebitFee)
	})
}
