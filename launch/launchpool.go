package launch

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListLaunchPoolProjectsService -- GET /api/v4/launch/project-list (public)
//
// Lists LaunchPool projects with their reward-pool configuration. The query
// parameters are camelCase (mortgageCoin, searchCoin, limitRule, sortType,
// pageSize); snake_case spellings are silently ignored by Gate.
type ListLaunchPoolProjectsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListLaunchPoolProjectsService() *ListLaunchPoolProjectsService {
	return &ListLaunchPoolProjectsService{c: c, params: map[string]string{}}
}

// SetStatus filters by project status: 0 all, 1 ongoing, 2 warming up, 3 ended,
// 4 ongoing and warming up.
func (s *ListLaunchPoolProjectsService) SetStatus(status int) *ListLaunchPoolProjectsService {
	s.params["status"] = strconv.Itoa(status)
	return s
}

// SetMortgageCoin filters by staking currency (exact match).
func (s *ListLaunchPoolProjectsService) SetMortgageCoin(mortgageCoin string) *ListLaunchPoolProjectsService {
	s.params["mortgageCoin"] = mortgageCoin
	return s
}

// SetSearchCoin filters by reward currency or project name (fuzzy match).
func (s *ListLaunchPoolProjectsService) SetSearchCoin(searchCoin string) *ListLaunchPoolProjectsService {
	s.params["searchCoin"] = searchCoin
	return s
}

// SetLimitRule filters by pool rule: 0 regular pool, 1 beginner pool.
func (s *ListLaunchPoolProjectsService) SetLimitRule(limitRule int) *ListLaunchPoolProjectsService {
	s.params["limitRule"] = strconv.Itoa(limitRule)
	return s
}

// SetSortType selects the order: 1 max APR descending, 2 max APR ascending.
func (s *ListLaunchPoolProjectsService) SetSortType(sortType int) *ListLaunchPoolProjectsService {
	s.params["sortType"] = strconv.Itoa(sortType)
	return s
}

// SetPage selects the page number (default 1).
func (s *ListLaunchPoolProjectsService) SetPage(page int) *ListLaunchPoolProjectsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of projects per page (default 10, max 30).
func (s *ListLaunchPoolProjectsService) SetPageSize(pageSize int) *ListLaunchPoolProjectsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListLaunchPoolProjectsService) Do(ctx context.Context) ([]LaunchPoolProject, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/project-list", s.params)
	resp, err := request.Do[[]LaunchPoolProject](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateLaunchPoolOrderService -- POST /api/v4/launch/create-order (private)
//
// Stakes amount into reward pool rid of LaunchPool project pid to start
// staking mining.
type CreateLaunchPoolOrderService struct {
	c    *LaunchClient
	body map[string]any
}

func (c *LaunchClient) NewCreateLaunchPoolOrderService(pid, rid int64, amount decimal.Decimal) *CreateLaunchPoolOrderService {
	return &CreateLaunchPoolOrderService{c: c, body: map[string]any{
		"pid":    pid,
		"rid":    rid,
		"amount": amount.String(),
	}}
}

func (s *CreateLaunchPoolOrderService) Do(ctx context.Context) (*LaunchPoolOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/launch/create-order", s.body).WithSign()
	return request.Do[LaunchPoolOrderResult](req)
}

// RedeemLaunchPoolService -- POST /api/v4/launch/redeem (private)
//
// Redeems amount from reward pool rid of LaunchPool project pid, ending
// staking mining for the redeemed amount.
type RedeemLaunchPoolService struct {
	c    *LaunchClient
	body map[string]any
}

func (c *LaunchClient) NewRedeemLaunchPoolService(pid, rid int64, amount decimal.Decimal) *RedeemLaunchPoolService {
	return &RedeemLaunchPoolService{c: c, body: map[string]any{
		"pid":    pid,
		"rid":    rid,
		"amount": amount.String(),
	}}
}

func (s *RedeemLaunchPoolService) Do(ctx context.Context) (*LaunchPoolRedeemResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/launch/redeem", s.body).WithSign()
	return request.Do[LaunchPoolRedeemResult](req)
}

// ListLaunchPoolPledgeRecordsService -- GET /api/v4/launch/user-pledge-records (private)
//
// Lists the user's LaunchPool staking and redemption records. Like the project
// list, its query parameters are camelCase (pageSize, startTime, endTime).
type ListLaunchPoolPledgeRecordsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListLaunchPoolPledgeRecordsService() *ListLaunchPoolPledgeRecordsService {
	return &ListLaunchPoolPledgeRecordsService{c: c, params: map[string]string{}}
}

// SetPage selects the page number (default 1).
func (s *ListLaunchPoolPledgeRecordsService) SetPage(page int) *ListLaunchPoolPledgeRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (max 30).
func (s *ListLaunchPoolPledgeRecordsService) SetPageSize(pageSize int) *ListLaunchPoolPledgeRecordsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetType filters by record type: 1 pledge, 2 redemption.
func (s *ListLaunchPoolPledgeRecordsService) SetType(recordType int) *ListLaunchPoolPledgeRecordsService {
	s.params["type"] = strconv.Itoa(recordType)
	return s
}

// SetStartTime sets the start of the query range, formatted
// "YYYY-MM-DD HH:MM:SS".
func (s *ListLaunchPoolPledgeRecordsService) SetStartTime(startTime string) *ListLaunchPoolPledgeRecordsService {
	s.params["startTime"] = startTime
	return s
}

// SetEndTime sets the end of the query range, formatted "YYYY-MM-DD HH:MM:SS".
func (s *ListLaunchPoolPledgeRecordsService) SetEndTime(endTime string) *ListLaunchPoolPledgeRecordsService {
	s.params["endTime"] = endTime
	return s
}

// SetCoin filters by staking currency.
func (s *ListLaunchPoolPledgeRecordsService) SetCoin(coin string) *ListLaunchPoolPledgeRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *ListLaunchPoolPledgeRecordsService) Do(ctx context.Context) ([]LaunchPoolPledgeRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/user-pledge-records", s.params).WithSign()
	resp, err := request.Do[[]LaunchPoolPledgeRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListLaunchPoolRewardRecordsService -- GET /api/v4/launch/get-user-reward-records (private)
//
// Lists the user's LaunchPool staking reward records. Paging is camelCase
// (pageSize) while the time range is snake_case (start_time, end_time).
type ListLaunchPoolRewardRecordsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListLaunchPoolRewardRecordsService() *ListLaunchPoolRewardRecordsService {
	return &ListLaunchPoolRewardRecordsService{c: c, params: map[string]string{}}
}

// SetPage selects the page number (default 1).
func (s *ListLaunchPoolRewardRecordsService) SetPage(page int) *ListLaunchPoolRewardRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (max 30).
func (s *ListLaunchPoolRewardRecordsService) SetPageSize(pageSize int) *ListLaunchPoolRewardRecordsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetStartTime sets the start of the query range.
func (s *ListLaunchPoolRewardRecordsService) SetStartTime(startTime time.Time) *ListLaunchPoolRewardRecordsService {
	s.params["start_time"] = strconv.FormatInt(startTime.Unix(), 10)
	return s
}

// SetEndTime sets the end of the query range.
func (s *ListLaunchPoolRewardRecordsService) SetEndTime(endTime time.Time) *ListLaunchPoolRewardRecordsService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetCoin filters by reward currency.
func (s *ListLaunchPoolRewardRecordsService) SetCoin(coin string) *ListLaunchPoolRewardRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *ListLaunchPoolRewardRecordsService) Do(ctx context.Context) ([]LaunchPoolRewardRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/get-user-reward-records", s.params).WithSign()
	resp, err := request.Do[[]LaunchPoolRewardRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LaunchPoolProject is a LaunchPool project and its reward pools. Times are
// bare Unix seconds; warmup_timest is 0 (decoded as the Unix epoch) when the
// project has no warm-up phase. project_state is 1 ongoing, 2 warming up,
// 3 ended. warmup_timest is returned live but absent from Gate's API spec.
type LaunchPoolProject struct {
	PID          int64                  `json:"pid"`
	Name         string                 `json:"name"`
	TotalAmount  decimal.Decimal        `json:"total_amount"`
	StartTime    time.Time              `json:"start_timest,format:unix"`
	WarmupTime   time.Time              `json:"warmup_timest,format:unix"`
	EndTime      time.Time              `json:"end_timest,format:unix"`
	Days         int                    `json:"days"`
	ProjectState int                    `json:"project_state"`
	RewardPools  []LaunchPoolRewardPool `json:"reward_pools"`
}

// LaunchPoolRewardPool is one staking pool of a LaunchPool project. coin is the
// staking currency and rate_year the estimated APR in percent ("6.24" =
// 6.24%). personal_max_amount / personal_min_amount are "0" when unbounded.
// limit_rule is 0 regular pool, 1 beginner pool (3 is also seen live,
// undocumented). limit_rule, net_deposit_cnt_sum (always "" so far) and
// transaction_type are returned live but absent from Gate's API spec.
type LaunchPoolRewardPool struct {
	RID                   int64                         `json:"rid"`
	Coin                  string                        `json:"coin"`
	RateYear              decimal.Decimal               `json:"rate_year"`
	AlreadyBuyTotalAmount decimal.Decimal               `json:"already_buy_total_amount"`
	PersonalMaxAmount     decimal.Decimal               `json:"personal_max_amount"`
	PersonalMinAmount     decimal.Decimal               `json:"personal_min_amount"`
	LimitRule             int                           `json:"limit_rule"`
	NetDepositCountSum    string                        `json:"net_deposit_cnt_sum"`
	TransactionType       int                           `json:"transaction_type"`
	TransactionConfig     []LaunchPoolTransactionConfig `json:"transaction_config"`
}

// LaunchPoolTransactionConfig is one participation tier of a reward pool: a
// trading-volume requirement and the staking cap it unlocks. Both arrive as
// bare JSON numbers.
type LaunchPoolTransactionConfig struct {
	TransactionAmount decimal.Decimal `json:"transaction_amount"`
	MortgageLimit     decimal.Decimal `json:"mortgage_limit"`
}

// LaunchPoolOrderResult is the result of a LaunchPool staking order.
type LaunchPoolOrderResult struct {
	FlowID int64 `json:"flow_id"`
}

// LaunchPoolRedeemResult is the result of a LaunchPool redemption.
type LaunchPoolRedeemResult struct {
	Success bool `json:"success"`
}

// LaunchPoolPledgeRecord is one LaunchPool staking (type 1) or redemption
// (type 2) record. coin is the staking currency. create_timest is the creation
// date as a string (layout unconfirmed live).
type LaunchPoolPledgeRecord struct {
	ID         int64           `json:"id"`
	CreateTime string          `json:"create_timest"`
	RewardCoin string          `json:"reward_coin"`
	Coin       string          `json:"coin"`
	Type       int             `json:"type"`
	Amount     decimal.Decimal `json:"amount"`
}

// LaunchPoolRewardRecord is one LaunchPool reward distribution. amount_base is
// the base reward and amount_ext the bonus reward. reward_timest is the reward
// date as a string (layout unconfirmed live).
type LaunchPoolRewardRecord struct {
	RewardTime          string          `json:"reward_timest"`
	Coin                string          `json:"coin"`
	ValidMortgageAmount decimal.Decimal `json:"valid_mortgage_amount"`
	AmountBase          decimal.Decimal `json:"amount_base"`
	AmountExt           decimal.Decimal `json:"amount_ext"`
}
