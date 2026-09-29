package rebate

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

// TestRebate live-tests every rebate endpoint. All of them are private and
// capability-gated (the account is likely not an agency/partner/broker), so a
// Gate "not allowed / not found" response is treated as a pass: the request path
// and signing were exercised correctly. Real data is diffed with AssertCovers.
func TestRebate(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// check runs the standard capability-gated read assertion: tolerate the
	// "account lacks this capability" errors, fail on any other error (such as
	// a response that does not decode), and diff the raw response otherwise.
	check := func(label, path string, params map[string]string, err error, resp any) {
		t.Helper()
		if err != nil {
			if testutil.Tolerable(t, label, err) {
				return
			}
			t.Errorf("%s: %v", label, err)
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, path, params, true)
		testutil.AssertCovers(t, label, raw, resp)
	}

	agencyTx, err := c.NewAgencyTransactionHistoryService().SetLimit(2).Do(cx)
	if agencyTx != nil {
		t.Logf("agency transaction total=%d records=%d", agencyTx.Total, len(agencyTx.List))
	}
	check("rebate/agency/transaction_history", "/api/v4/rebate/agency/transaction_history",
		map[string]string{"limit": "2"}, err, agencyTx)

	agencyComm, err := c.NewAgencyCommissionsHistoryService().SetLimit(2).Do(cx)
	if agencyComm != nil {
		t.Logf("agency commission total=%d records=%d", agencyComm.Total, len(agencyComm.List))
	}
	check("rebate/agency/commission_history", "/api/v4/rebate/agency/commission_history",
		map[string]string{"limit": "2"}, err, agencyComm)

	partnerTx, err := c.NewPartnerTransactionHistoryService().SetLimit(2).Do(cx)
	if partnerTx != nil {
		t.Logf("partner transaction total=%d records=%d", partnerTx.Total, len(partnerTx.List))
	}
	check("rebate/partner/transaction_history", "/api/v4/rebate/partner/transaction_history",
		map[string]string{"limit": "2"}, err, partnerTx)

	partnerComm, err := c.NewPartnerCommissionsHistoryService().SetLimit(2).Do(cx)
	if partnerComm != nil {
		t.Logf("partner commission total=%d records=%d", partnerComm.Total, len(partnerComm.List))
	}
	check("rebate/partner/commission_history", "/api/v4/rebate/partner/commission_history",
		map[string]string{"limit": "2"}, err, partnerComm)

	partnerSub, err := c.NewPartnerSubListService().SetLimit(2).Do(cx)
	if partnerSub != nil {
		t.Logf("partner sub total=%d records=%d", partnerSub.Total, len(partnerSub.List))
	}
	check("rebate/partner/sub_list", "/api/v4/rebate/partner/sub_list",
		map[string]string{"limit": "2"}, err, partnerSub)

	brokerComm, err := c.NewRebateBrokerCommissionHistoryService().SetLimit(2).Do(cx)
	if brokerComm != nil {
		t.Logf("broker commission total=%d records=%d", brokerComm.Total, len(brokerComm.List))
	}
	check("rebate/broker/commission_history", "/api/v4/rebate/broker/commission_history",
		map[string]string{"limit": "2"}, err, brokerComm)

	brokerTx, err := c.NewRebateBrokerTransactionHistoryService().SetLimit(2).Do(cx)
	if brokerTx != nil {
		t.Logf("broker transaction total=%d records=%d", brokerTx.Total, len(brokerTx.List))
	}
	check("rebate/broker/transaction_history", "/api/v4/rebate/broker/transaction_history",
		map[string]string{"limit": "2"}, err, brokerTx)

	userInfo, err := c.NewRebateUserInfoService().Do(cx)
	if userInfo != nil {
		t.Logf("rebate user info invite_uid=%d", userInfo.InviteUID)
	}
	check("rebate/user/info", "/api/v4/rebate/user/info", nil, err, userInfo)

	subRelation, err := c.NewUserSubRelationService("10000").Do(cx)
	if subRelation != nil {
		t.Logf("user sub relation records=%d", len(subRelation.List))
	}
	check("rebate/user/sub_relation", "/api/v4/rebate/user/sub_relation",
		map[string]string{"user_id_list": "10000"}, err, subRelation)
}
