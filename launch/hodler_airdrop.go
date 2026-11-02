package launch

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListHodlerAirdropProjectsService -- GET /api/v4/launch/hodler-airdrop/project-list (private)
//
// Lists HODLer Airdrop activities, filtered by status, currency/project name
// and participation. Gate's API spec marks it public, but the live gateway
// rejects unsigned requests (MISSING_REQUIRED_HEADER for Timestamp / KEY), so
// it is signed; signed callers also get their own participation data.
type ListHodlerAirdropProjectsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListHodlerAirdropProjectsService() *ListHodlerAirdropProjectsService {
	return &ListHodlerAirdropProjectsService{c: c, params: map[string]string{}}
}

// SetStatus filters by activity status: ACTIVE (ongoing + warming up),
// UNDERWAY (ongoing), PREHEAT (warming up) or FINISH (ended). All when unset.
func (s *ListHodlerAirdropProjectsService) SetStatus(status string) *ListHodlerAirdropProjectsService {
	s.params["status"] = status
	return s
}

// SetKeyword filters by currency or project name (fuzzy match).
func (s *ListHodlerAirdropProjectsService) SetKeyword(keyword string) *ListHodlerAirdropProjectsService {
	s.params["keyword"] = keyword
	return s
}

// SetJoin filters by participation: 0 all (default), 1 only participated.
func (s *ListHodlerAirdropProjectsService) SetJoin(join int) *ListHodlerAirdropProjectsService {
	s.params["join"] = strconv.Itoa(join)
	return s
}

// SetPage selects the page number (default 1).
func (s *ListHodlerAirdropProjectsService) SetPage(page int) *ListHodlerAirdropProjectsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetSize caps the number of activities per page (default 10).
func (s *ListHodlerAirdropProjectsService) SetSize(size int) *ListHodlerAirdropProjectsService {
	s.params["size"] = strconv.Itoa(size)
	return s
}

func (s *ListHodlerAirdropProjectsService) Do(ctx context.Context) ([]HodlerAirdropProject, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/hodler-airdrop/project-list", s.params).WithSign()
	resp, err := request.Do[[]HodlerAirdropProject](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateHodlerAirdropOrderService -- POST /api/v4/launch/hodler-airdrop/order (private)
//
// Joins HODLer Airdrop activity hodlerID. Requires holding GT and passing KYC;
// sub-accounts and institutional accounts are not supported.
type CreateHodlerAirdropOrderService struct {
	c    *LaunchClient
	body map[string]any
}

func (c *LaunchClient) NewCreateHodlerAirdropOrderService(hodlerID int64) *CreateHodlerAirdropOrderService {
	return &CreateHodlerAirdropOrderService{c: c, body: map[string]any{
		"hodler_id": hodlerID,
	}}
}

func (s *CreateHodlerAirdropOrderService) Do(ctx context.Context) (*HodlerAirdropOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/launch/hodler-airdrop/order", s.body).WithSign()
	return request.Do[HodlerAirdropOrderResult](req)
}

// ListHodlerAirdropOrderRecordsService -- GET /api/v4/launch/hodler-airdrop/user-order-records (private)
//
// Lists the user's HODLer Airdrop participation records with the effective GT
// holding and airdrop amounts of each activity.
type ListHodlerAirdropOrderRecordsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListHodlerAirdropOrderRecordsService() *ListHodlerAirdropOrderRecordsService {
	return &ListHodlerAirdropOrderRecordsService{c: c, params: map[string]string{}}
}

// SetKeyword filters by currency name.
func (s *ListHodlerAirdropOrderRecordsService) SetKeyword(keyword string) *ListHodlerAirdropOrderRecordsService {
	s.params["keyword"] = keyword
	return s
}

// SetStartTime sets the start of the query range.
func (s *ListHodlerAirdropOrderRecordsService) SetStartTime(startTime time.Time) *ListHodlerAirdropOrderRecordsService {
	s.params["start_timest"] = strconv.FormatInt(startTime.Unix(), 10)
	return s
}

// SetEndTime sets the end of the query range.
func (s *ListHodlerAirdropOrderRecordsService) SetEndTime(endTime time.Time) *ListHodlerAirdropOrderRecordsService {
	s.params["end_timest"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetPage selects the page number (default 1).
func (s *ListHodlerAirdropOrderRecordsService) SetPage(page int) *ListHodlerAirdropOrderRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetSize caps the number of records per page (default 10).
func (s *ListHodlerAirdropOrderRecordsService) SetSize(size int) *ListHodlerAirdropOrderRecordsService {
	s.params["size"] = strconv.Itoa(size)
	return s
}

func (s *ListHodlerAirdropOrderRecordsService) Do(ctx context.Context) ([]HodlerAirdropOrderRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/hodler-airdrop/user-order-records", s.params).WithSign()
	resp, err := request.Do[[]HodlerAirdropOrderRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListHodlerAirdropRecordsService -- GET /api/v4/launch/hodler-airdrop/user-airdrop-records (private)
//
// Lists the HODLer Airdrop distributions the user has received: base and bonus
// airdrops plus their auto-exchange status.
type ListHodlerAirdropRecordsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListHodlerAirdropRecordsService() *ListHodlerAirdropRecordsService {
	return &ListHodlerAirdropRecordsService{c: c, params: map[string]string{}}
}

// SetKeyword filters by currency name.
func (s *ListHodlerAirdropRecordsService) SetKeyword(keyword string) *ListHodlerAirdropRecordsService {
	s.params["keyword"] = keyword
	return s
}

// SetStartTime sets the start of the query range.
func (s *ListHodlerAirdropRecordsService) SetStartTime(startTime time.Time) *ListHodlerAirdropRecordsService {
	s.params["start_timest"] = strconv.FormatInt(startTime.Unix(), 10)
	return s
}

// SetEndTime sets the end of the query range.
func (s *ListHodlerAirdropRecordsService) SetEndTime(endTime time.Time) *ListHodlerAirdropRecordsService {
	s.params["end_timest"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetPage selects the page number (default 1).
func (s *ListHodlerAirdropRecordsService) SetPage(page int) *ListHodlerAirdropRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetSize caps the number of records per page (default 10).
func (s *ListHodlerAirdropRecordsService) SetSize(size int) *ListHodlerAirdropRecordsService {
	s.params["size"] = strconv.Itoa(size)
	return s
}

func (s *ListHodlerAirdropRecordsService) Do(ctx context.Context) ([]HodlerAirdropRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/hodler-airdrop/user-airdrop-records", s.params).WithSign()
	resp, err := request.Do[[]HodlerAirdropRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// HodlerAirdropProject is a HODLer Airdrop activity. asset is the airdropped
// currency; open_timest / close_timest are UTC "Y-m-d H:i:s" strings.
// per_gt_reward_token is the airdrop per GT held ("" — decoded as zero — while
// it is still being calculated) and max_queue_amount the personal GT cap.
type HodlerAirdropProject struct {
	HodlerID         int64           `json:"hodler_id,string"`
	Name             string          `json:"name"`
	Asset            string          `json:"asset"`
	Status           string          `json:"status"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	OpenTime         string          `json:"open_timest"`
	CloseTime        string          `json:"close_timest"`
	PerGTRewardToken decimal.Decimal `json:"per_gt_reward_token"`
	UserCount        int64           `json:"user_count,string"`
	MaxQueueAmount   decimal.Decimal `json:"max_queue_amount"`
}

// HodlerAirdropOrderResult is the result of joining a HODLer Airdrop activity:
// the currencies that will be airdropped.
type HodlerAirdropOrderResult struct {
	InterestCoin []string `json:"interest_coin"`
}

// HodlerAirdropOrderRecord is one HODLer Airdrop participation.
// user_gt_avg_amount is the effective (average) GT holding; the receive
// amounts are "" (decoded as zero) until the activity ends, and
// receive_amount_ext is the Simple Earn bonus. close_timest is a UTC
// "Y-m-d H:i:s" string.
type HodlerAirdropOrderRecord struct {
	HodlerID          int64           `json:"hodler_id,string"`
	Asset             string          `json:"asset"`
	UserGTAvgAmount   decimal.Decimal `json:"user_gt_avg_amount"`
	ReceiveAmountBase decimal.Decimal `json:"receive_amount_base"`
	ReceiveAmountExt  decimal.Decimal `json:"receive_amount_ext"`
	CloseTime         string          `json:"close_timest"`
}

// HodlerAirdropRecord is one HODLer Airdrop distribution. receive_amount_ext
// is the Simple Earn bonus. delivery_timest is a UTC "Y-m-d H:i:s" string, ""
// until distributed. exchange is "0" not exchanged, "1" auto-exchanged to GT.
type HodlerAirdropRecord struct {
	HodlerID          int64           `json:"hodler_id,string"`
	Asset             string          `json:"asset"`
	ReceiveAmountBase decimal.Decimal `json:"receive_amount_base"`
	ReceiveAmountExt  decimal.Decimal `json:"receive_amount_ext"`
	DeliveryTime      string          `json:"delivery_timest"`
	Exchange          string          `json:"exchange"`
}
