package launch

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestLaunch(t *testing.T) {
	// ---- Public endpoints ----

	t.Run("ListLaunchPoolProjects", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		got, err := c.NewListLaunchPoolProjectsService().SetStatus(0).SetPageSize(30).Do(cx)
		if err != nil {
			t.Fatalf("launchpool projects: %v", err)
		}
		if len(got) == 0 {
			t.Fatal("no launchpool projects returned")
		}
		// start_timest is bare Unix seconds; a wrong unit would land far from now.
		if y := got[0].StartTime.Year(); y < 2020 || y > 2100 {
			t.Fatalf("start_timest decoded to implausible time %v", got[0].StartTime)
		}
		t.Logf("projects=%d first=%+v", len(got), got[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/project-list",
			map[string]string{"status": "0", "pageSize": "30"}, false)
		testutil.AssertCovers(t, "launch/project-list", raw, got)
	})

	t.Run("ListLaunchPoolProjectsFiltered", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		// An ended project whose last reward pool carries transaction_config
		// tiers, so the nested tier keys are diffed too.
		got, err := c.NewListLaunchPoolProjectsService().SetStatus(3).SetSearchCoin("LAPTOP").Do(cx)
		if err != nil {
			t.Fatalf("launchpool projects (filtered): %v", err)
		}
		t.Logf("projects=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/project-list",
			map[string]string{"status": "3", "searchCoin": "LAPTOP"}, false)
		testutil.AssertCovers(t, "launch/project-list(searchCoin)", raw, got)
	})

	// ---- Endpoints Gate's spec marks public but the live gateway requires a
	// signature for (unsigned calls get MISSING_REQUIRED_HEADER) ----

	t.Run("ListHodlerAirdropProjects", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHodlerAirdropProjectsService().SetSize(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/hodler-airdrop/project-list", err) {
				return
			}
			t.Fatalf("hodler airdrop projects: %v", err)
		}
		t.Logf("projects=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/hodler-airdrop/project-list",
			map[string]string{"size": "10"}, true)
		testutil.AssertCovers(t, "launch/hodler-airdrop/project-list", raw, got)
	})

	t.Run("ListCandyDropActivities", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListCandyDropActivitiesService().SetLimit(30).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/activity-list", err) {
				return
			}
			t.Fatalf("candydrop activities: %v", err)
		}
		t.Logf("activities=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/candydrop/activity-list",
			map[string]string{"limit": "30"}, true)
		testutil.AssertCovers(t, "launch/candydrop/activity-list", raw, got)
	})

	t.Run("GetCandyDropActivityRules", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		// The rules need a real activity, so discover one from the list first.
		activities, err := c.NewListCandyDropActivitiesService().SetLimit(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/activity-rules", err) {
				return
			}
			t.Fatalf("candydrop activities (for rules): %v", err)
		}
		if len(activities) == 0 {
			t.Skip("no candydrop activity to fetch rules for")
		}
		currency := activities[0].Currency
		got, err := c.NewGetCandyDropActivityRulesService().SetCurrency(currency).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/activity-rules", err) {
				return
			}
			t.Fatalf("candydrop activity rules: %v", err)
		}
		t.Logf("currency=%s prize_pools=%d", got.Currency, len(got.PrizePools))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/candydrop/activity-rules",
			map[string]string{"currency": currency}, true)
		testutil.AssertCovers(t, "launch/candydrop/activity-rules", raw, got)
	})

	// ---- Private read endpoints ----

	t.Run("ListLaunchPoolPledgeRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListLaunchPoolPledgeRecordsService().SetPageSize(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/user-pledge-records", err) {
				return
			}
			t.Fatalf("pledge records: %v", err)
		}
		t.Logf("records=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/user-pledge-records",
			map[string]string{"pageSize": "10"}, true)
		testutil.AssertCovers(t, "launch/user-pledge-records", raw, got)
	})

	t.Run("ListLaunchPoolRewardRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListLaunchPoolRewardRecordsService().SetPageSize(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/get-user-reward-records", err) {
				return
			}
			t.Fatalf("reward records: %v", err)
		}
		t.Logf("records=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/get-user-reward-records",
			map[string]string{"pageSize": "10"}, true)
		testutil.AssertCovers(t, "launch/get-user-reward-records", raw, got)
	})

	t.Run("ListHodlerAirdropOrderRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHodlerAirdropOrderRecordsService().SetSize(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/hodler-airdrop/user-order-records", err) {
				return
			}
			t.Fatalf("hodler airdrop order records: %v", err)
		}
		t.Logf("records=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/hodler-airdrop/user-order-records",
			map[string]string{"size": "10"}, true)
		testutil.AssertCovers(t, "launch/hodler-airdrop/user-order-records", raw, got)
	})

	t.Run("ListHodlerAirdropRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListHodlerAirdropRecordsService().SetSize(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/hodler-airdrop/user-airdrop-records", err) {
				return
			}
			t.Fatalf("hodler airdrop records: %v", err)
		}
		t.Logf("records=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/hodler-airdrop/user-airdrop-records",
			map[string]string{"size": "10"}, true)
		testutil.AssertCovers(t, "launch/hodler-airdrop/user-airdrop-records", raw, got)
	})

	t.Run("GetCandyDropTaskProgress", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		// Progress only exists for an ongoing activity the account registered
		// for; skip when there is none.
		activities, err := c.NewListCandyDropActivitiesService().
			SetStatus("ongoing").SetRegisterStatus("registered").SetLimit(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/task-progress", err) {
				return
			}
			t.Fatalf("candydrop activities (for progress): %v", err)
		}
		if len(activities) == 0 {
			t.Skip("no registered ongoing candydrop activity")
		}
		currency := activities[0].Currency
		got, err := c.NewGetCandyDropTaskProgressService().SetCurrency(currency).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/task-progress", err) {
				return
			}
			t.Fatalf("candydrop task progress: %v", err)
		}
		t.Logf("currency=%s tasks=%d", got.Currency, len(got.Tasks))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/candydrop/task-progress",
			map[string]string{"currency": currency}, true)
		testutil.AssertCovers(t, "launch/candydrop/task-progress", raw, got)
	})

	t.Run("ListCandyDropParticipationRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListCandyDropParticipationRecordsService().SetLimit(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/participation-records", err) {
				return
			}
			t.Fatalf("candydrop participation records: %v", err)
		}
		t.Logf("records=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/candydrop/participation-records",
			map[string]string{"limit": "10"}, true)
		testutil.AssertCovers(t, "launch/candydrop/participation-records", raw, got)
	})

	t.Run("ListCandyDropAirdropRecords", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListCandyDropAirdropRecordsService().SetLimit(10).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "launch/candydrop/airdrop-records", err) {
				return
			}
			t.Fatalf("candydrop airdrop records: %v", err)
		}
		t.Logf("records=%d", len(got))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/launch/candydrop/airdrop-records",
			map[string]string{"limit": "10"}, true)
		testutil.AssertCovers(t, "launch/candydrop/airdrop-records", raw, got)
	})
}
