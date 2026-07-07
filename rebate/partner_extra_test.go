package rebate

import (
	"context"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestRebateExtra(t *testing.T) {
	// authed builds a signed, time-synced rebate client, skipping the subtest
	// when credentials are unset.
	authed := func(t *testing.T, cx context.Context) *RebateClient {
		t.Helper()
		c := testClient(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		return c
	}

	t.Run("PartnerRecentApplications", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		raw, err := c.NewPartnerRecentApplicationsService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "rebate/partner/applications/recent", err) {
				return
			}
			t.Fatalf("partner applications/recent: %v", err)
		}
		t.Logf("partner applications/recent: %d bytes", len(raw))
	})

	t.Run("PartnerEligibility", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		raw, err := c.NewPartnerEligibilityService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "rebate/partner/eligibility", err) {
				return
			}
			t.Fatalf("partner eligibility: %v", err)
		}
		t.Logf("partner eligibility: %d bytes", len(raw))
	})

	t.Run("PartnerAggregatedData", func(t *testing.T) {
		cx := testutil.Ctx(t)
		c := authed(t, cx)
		raw, err := c.NewPartnerAggregatedDataService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "rebate/partner/data/aggregated", err) {
				return
			}
			t.Fatalf("partner data/aggregated: %v", err)
		}
		t.Logf("partner data/aggregated: %d bytes", len(raw))
	})
}
