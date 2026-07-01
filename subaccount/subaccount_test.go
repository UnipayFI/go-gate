package subaccount

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestSubAccount(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// ListSubAccounts -- always runs; an account with no sub-accounts returns [].
	subs, err := c.NewListSubAccountsService().Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "sub_accounts", err) {
			return
		}
		t.Fatalf("list sub-accounts: %v", err)
	}
	t.Logf("sub-accounts=%d", len(subs))
	for _, s := range subs {
		t.Logf("  user_id=%d login=%s state=%d type=%d", s.UserID, s.LoginName, s.State, s.Type)
	}
	if len(subs) > 0 {
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/sub_accounts", nil, true)
		testutil.AssertCovers(t, "sub_accounts", raw, subs)
	}

	// ListUnifiedMode -- always runs.
	modes, err := c.NewListUnifiedModeService().Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "sub_accounts/unified_mode", err) {
			t.Errorf("list unified mode: %v", err)
		}
	} else {
		t.Logf("unified modes=%d", len(modes))
		if len(modes) > 0 {
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/sub_accounts/unified_mode", nil, true)
			testutil.AssertCovers(t, "sub_accounts/unified_mode", raw, modes)
		}
	}

	// Per-sub-account reads need a real sub user_id.
	if len(subs) == 0 {
		t.Log("no sub-accounts; skipping GetSubAccount / keys reads")
	} else {
		userID := subs[0].UserID
		idStr := strconv.FormatInt(userID, 10)

		sub, err := c.NewGetSubAccountService(userID).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "sub_accounts/{user_id}", err) {
				t.Errorf("get sub-account: %v", err)
			}
		} else {
			t.Logf("sub-account: %+v", sub)
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/sub_accounts/"+idStr, nil, true)
			testutil.AssertCovers(t, "sub_accounts/{user_id}", raw, sub)
		}

		keys, err := c.NewListSubAccountKeysService(userID).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "sub_accounts/{user_id}/keys", err) {
				t.Errorf("list sub-account keys: %v", err)
			}
		} else {
			t.Logf("sub-account keys=%d", len(keys))
			if len(keys) > 0 {
				raw := testutil.FetchRawGet(t, c, cx, "/api/v4/sub_accounts/"+idStr+"/keys", nil, true)
				testutil.AssertCovers(t, "sub_accounts/{user_id}/keys", raw, keys)

				key := keys[0].Key
				gk, err := c.NewGetSubAccountKeyService(userID, key).Do(cx)
				if err != nil {
					if !testutil.Tolerable(t, "sub_accounts/{user_id}/keys/{key}", err) {
						t.Errorf("get sub-account key: %v", err)
					}
				} else {
					t.Logf("sub-account key: %+v", gk)
					raw := testutil.FetchRawGet(t, c, cx, "/api/v4/sub_accounts/"+idStr+"/keys/"+key, nil, true)
					testutil.AssertCovers(t, "sub_accounts/{user_id}/keys/{key}", raw, gk)
				}
			}
		}
	}

	// State-changing endpoints: create/update/delete/lock/unlock. Opt-in only and
	// tolerant of capability errors, since sub-account management is provisioned.
	t.Run("write", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("GATE_TEST_WRITE!=1; skipping sub-account create/update/delete/lock/unlock")
		}

		// CreateSubAccounts -- make a fresh sub-account to exercise the full key
		// lifecycle without disturbing existing ones.
		loginName := fmt.Sprintf("gotest%d", time.Now().Unix())
		created, err := c.NewCreateSubAccountsService(loginName).
			SetRemark("go-gate sdk test").
			Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "create sub_accounts", err) {
				return
			}
			t.Fatalf("create sub-account: %v", err)
		}
		userID := created.UserID
		idStr := strconv.FormatInt(userID, 10)
		t.Logf("created sub-account user_id=%d login=%s", userID, created.LoginName)

		// CreateSubAccountKeys
		key, err := c.NewCreateSubAccountKeysService(userID).
			SetName("gotest").
			SetPerms([]SubAccountKeyPerm{{Name: "spot", ReadOnly: true}}).
			Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "create sub_account key", err) {
				t.Errorf("create sub-account key: %v", err)
			}
		} else {
			t.Logf("created key=%s", key.Key)

			if _, gerr := c.NewGetSubAccountKeyService(userID, key.Key).Do(cx); gerr != nil {
				if !testutil.Tolerable(t, "get sub_account key", gerr) {
					t.Errorf("get sub-account key: %v", gerr)
				}
			}

			// UpdateSubAccountKeys -- broaden the permission scope.
			if uerr := c.NewUpdateSubAccountKeysService(userID, key.Key).
				SetPerms([]SubAccountKeyPerm{{Name: "wallet", ReadOnly: true}}).
				Do(cx); uerr != nil {
				if !testutil.Tolerable(t, "update sub_account key", uerr) {
					t.Errorf("update sub-account key: %v", uerr)
				}
			}

			// DeleteSubAccountKeys -- clean up the key we created.
			if derr := c.NewDeleteSubAccountKeysService(userID, key.Key).Do(cx); derr != nil {
				if !testutil.Tolerable(t, "delete sub_account key", derr) {
					t.Errorf("delete sub-account key: %v", derr)
				}
			}
		}

		// LockSubAccount then UnlockSubAccount -- reversible pair.
		if lerr := c.NewLockSubAccountService(userID).Do(cx); lerr != nil {
			if !testutil.Tolerable(t, "lock sub_account", lerr) {
				t.Errorf("lock sub-account: %v", lerr)
			}
		} else {
			t.Logf("locked sub-account %s", idStr)
		}
		if uerr := c.NewUnlockSubAccountService(userID).Do(cx); uerr != nil {
			if !testutil.Tolerable(t, "unlock sub_account", uerr) {
				t.Errorf("unlock sub-account: %v", uerr)
			}
		} else {
			t.Logf("unlocked sub-account %s", idStr)
		}
	})
}
