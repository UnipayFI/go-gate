package rebate

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
)

// PartnerRecentApplicationsService -- GET /api/v4/rebate/partner/applications/recent (private)
//
// Returns the partner's recent application records. Gate does not document the
// response schema, so the raw JSON payload is returned for the caller to decode.
type PartnerRecentApplicationsService struct {
	c *RebateClient
}

func (c *RebateClient) NewPartnerRecentApplicationsService() *PartnerRecentApplicationsService {
	return &PartnerRecentApplicationsService{c: c}
}

func (s *PartnerRecentApplicationsService) Do(ctx context.Context) (json.RawMessage, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/partner/applications/recent").WithSign()
	raw, err := request.DoRaw(req)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

// PartnerEligibilityService -- GET /api/v4/rebate/partner/eligibility (private)
//
// Checks whether the account is eligible to apply as a partner. Gate does not
// document the response schema, so the raw JSON payload is returned.
type PartnerEligibilityService struct {
	c *RebateClient
}

func (c *RebateClient) NewPartnerEligibilityService() *PartnerEligibilityService {
	return &PartnerEligibilityService{c: c}
}

func (s *PartnerEligibilityService) Do(ctx context.Context) (json.RawMessage, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/partner/eligibility").WithSign()
	raw, err := request.DoRaw(req)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

// PartnerAggregatedDataService -- GET /api/v4/rebate/partner/data/aggregated (private)
//
// Returns aggregated partner agent statistics. Gate does not document the
// response schema, so the raw JSON payload is returned.
type PartnerAggregatedDataService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewPartnerAggregatedDataService() *PartnerAggregatedDataService {
	return &PartnerAggregatedDataService{c: c, params: map[string]string{}}
}

// SetStartDate sets the query start time, formatted as "yyyy-mm-dd hh:ii:ss"
// (UTC+8). Defaults server-side to the start of the current period.
func (s *PartnerAggregatedDataService) SetStartDate(startDate string) *PartnerAggregatedDataService {
	s.params["start_date"] = startDate
	return s
}

// SetEndDate sets the query end time, formatted as "yyyy-mm-dd hh:ii:ss"
// (UTC+8). Defaults server-side to the end of the current period.
func (s *PartnerAggregatedDataService) SetEndDate(endDate string) *PartnerAggregatedDataService {
	s.params["end_date"] = endDate
	return s
}

// SetBusinessType filters by business type: 0 - All (default), 1 - Spot,
// 2 - Futures, 3 - Alpha, 4 - Web3, 5 - and so on.
func (s *PartnerAggregatedDataService) SetBusinessType(businessType int) *PartnerAggregatedDataService {
	s.params["business_type"] = strconv.Itoa(businessType)
	return s
}

func (s *PartnerAggregatedDataService) Do(ctx context.Context) (json.RawMessage, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/partner/data/aggregated", s.params).WithSign()
	raw, err := request.DoRaw(req)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}
