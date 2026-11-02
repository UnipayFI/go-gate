package launch

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListCandyDropActivitiesService -- GET /api/v4/launch/candydrop/activity-list (private)
//
// Lists CandyDrop activities with multi-dimensional filters. Gate's API spec
// marks it public, but the live gateway rejects unsigned requests
// (MISSING_REQUIRED_HEADER for Timestamp / KEY), so it is signed.
type ListCandyDropActivitiesService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListCandyDropActivitiesService() *ListCandyDropActivitiesService {
	return &ListCandyDropActivitiesService{c: c, params: map[string]string{}}
}

// SetStatus filters by activity status: ongoing, upcoming or ended. All when
// unset.
func (s *ListCandyDropActivitiesService) SetStatus(status string) *ListCandyDropActivitiesService {
	s.params["status"] = status
	return s
}

// SetRuleName filters by task type: spot, futures, deposit, invite,
// trading_bot, simple_earn, first_deposit, alpha, flash_swap, tradfi or etf.
func (s *ListCandyDropActivitiesService) SetRuleName(ruleName string) *ListCandyDropActivitiesService {
	s.params["rule_name"] = ruleName
	return s
}

// SetRegisterStatus filters by participation: registered or unregistered. All
// when unset.
func (s *ListCandyDropActivitiesService) SetRegisterStatus(registerStatus string) *ListCandyDropActivitiesService {
	s.params["register_status"] = registerStatus
	return s
}

// SetCurrency filters by currency name.
func (s *ListCandyDropActivitiesService) SetCurrency(currency string) *ListCandyDropActivitiesService {
	s.params["currency"] = currency
	return s
}

// SetLimit caps the number of activities returned (default 10, max 30).
func (s *ListCandyDropActivitiesService) SetLimit(limit int) *ListCandyDropActivitiesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset (default 0).
func (s *ListCandyDropActivitiesService) SetOffset(offset int) *ListCandyDropActivitiesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListCandyDropActivitiesService) Do(ctx context.Context) ([]CandyDropActivity, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/candydrop/activity-list", s.params).WithSign()
	resp, err := request.Do[[]CandyDropActivity](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetCandyDropActivityRulesService -- GET /api/v4/launch/candydrop/activity-rules (private)
//
// Returns the rules of one CandyDrop activity: its prize pools and their
// tasks. Set at least one of activity_id or currency. Gate's API spec marks it
// public, but the live gateway rejects unsigned requests
// (MISSING_REQUIRED_HEADER for Timestamp / KEY), so it is signed.
type GetCandyDropActivityRulesService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewGetCandyDropActivityRulesService() *GetCandyDropActivityRulesService {
	return &GetCandyDropActivityRulesService{c: c, params: map[string]string{}}
}

// SetActivityID selects the activity by ID (alternative to SetCurrency).
func (s *GetCandyDropActivityRulesService) SetActivityID(activityID int64) *GetCandyDropActivityRulesService {
	s.params["activity_id"] = strconv.FormatInt(activityID, 10)
	return s
}

// SetCurrency selects the activity by project/currency name (alternative to
// SetActivityID).
func (s *GetCandyDropActivityRulesService) SetCurrency(currency string) *GetCandyDropActivityRulesService {
	s.params["currency"] = currency
	return s
}

func (s *GetCandyDropActivityRulesService) Do(ctx context.Context) (*CandyDropActivityRules, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/candydrop/activity-rules", s.params).WithSign()
	return request.Do[CandyDropActivityRules](req)
}

// RegisterCandyDropService -- POST /api/v4/launch/candydrop/register (private)
//
// Registers for the CandyDrop activity of currency.
type RegisterCandyDropService struct {
	c    *LaunchClient
	body map[string]any
}

func (c *LaunchClient) NewRegisterCandyDropService(currency string) *RegisterCandyDropService {
	return &RegisterCandyDropService{c: c, body: map[string]any{
		"currency": currency,
	}}
}

// SetActivityID pins the activity ID, used together with currency.
func (s *RegisterCandyDropService) SetActivityID(activityID int64) *RegisterCandyDropService {
	s.body["activity_id"] = activityID
	return s
}

func (s *RegisterCandyDropService) Do(ctx context.Context) (*CandyDropRegisterResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/launch/candydrop/register", s.body).WithSign()
	return request.Do[CandyDropRegisterResult](req)
}

// GetCandyDropTaskProgressService -- GET /api/v4/launch/candydrop/task-progress (private)
//
// Returns the user's progress on the tasks of an ongoing CandyDrop activity
// they registered for. Set at least one of activity_id or currency.
type GetCandyDropTaskProgressService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewGetCandyDropTaskProgressService() *GetCandyDropTaskProgressService {
	return &GetCandyDropTaskProgressService{c: c, params: map[string]string{}}
}

// SetActivityID selects the activity by ID (alternative to SetCurrency).
func (s *GetCandyDropTaskProgressService) SetActivityID(activityID int64) *GetCandyDropTaskProgressService {
	s.params["activity_id"] = strconv.FormatInt(activityID, 10)
	return s
}

// SetCurrency selects the activity by project/currency name (alternative to
// SetActivityID).
func (s *GetCandyDropTaskProgressService) SetCurrency(currency string) *GetCandyDropTaskProgressService {
	s.params["currency"] = currency
	return s
}

func (s *GetCandyDropTaskProgressService) Do(ctx context.Context) (*CandyDropTaskProgress, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/candydrop/task-progress", s.params).WithSign()
	return request.Do[CandyDropTaskProgress](req)
}

// ListCandyDropParticipationRecordsService -- GET /api/v4/launch/candydrop/participation-records (private)
//
// Lists the user's CandyDrop participation records.
type ListCandyDropParticipationRecordsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListCandyDropParticipationRecordsService() *ListCandyDropParticipationRecordsService {
	return &ListCandyDropParticipationRecordsService{c: c, params: map[string]string{}}
}

// SetCurrency filters by currency name.
func (s *ListCandyDropParticipationRecordsService) SetCurrency(currency string) *ListCandyDropParticipationRecordsService {
	s.params["currency"] = currency
	return s
}

// SetStatus filters by status: ongoing, awaiting_draw, won or not_win.
func (s *ListCandyDropParticipationRecordsService) SetStatus(status string) *ListCandyDropParticipationRecordsService {
	s.params["status"] = status
	return s
}

// SetStartTime sets the start of the query range.
func (s *ListCandyDropParticipationRecordsService) SetStartTime(startTime time.Time) *ListCandyDropParticipationRecordsService {
	s.params["start_time"] = strconv.FormatInt(startTime.Unix(), 10)
	return s
}

// SetEndTime sets the end of the query range.
func (s *ListCandyDropParticipationRecordsService) SetEndTime(endTime time.Time) *ListCandyDropParticipationRecordsService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetPage selects the page number (default 1).
func (s *ListCandyDropParticipationRecordsService) SetPage(page int) *ListCandyDropParticipationRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records per page (default 10, max 30).
func (s *ListCandyDropParticipationRecordsService) SetLimit(limit int) *ListCandyDropParticipationRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListCandyDropParticipationRecordsService) Do(ctx context.Context) ([]CandyDropParticipationRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/candydrop/participation-records", s.params).WithSign()
	resp, err := request.Do[[]CandyDropParticipationRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListCandyDropAirdropRecordsService -- GET /api/v4/launch/candydrop/airdrop-records (private)
//
// Lists the CandyDrop airdrops the user has received.
type ListCandyDropAirdropRecordsService struct {
	c      *LaunchClient
	params map[string]string
}

func (c *LaunchClient) NewListCandyDropAirdropRecordsService() *ListCandyDropAirdropRecordsService {
	return &ListCandyDropAirdropRecordsService{c: c, params: map[string]string{}}
}

// SetCurrency filters by currency name.
func (s *ListCandyDropAirdropRecordsService) SetCurrency(currency string) *ListCandyDropAirdropRecordsService {
	s.params["currency"] = currency
	return s
}

// SetStartTime sets the start of the query range.
func (s *ListCandyDropAirdropRecordsService) SetStartTime(startTime time.Time) *ListCandyDropAirdropRecordsService {
	s.params["start_time"] = strconv.FormatInt(startTime.Unix(), 10)
	return s
}

// SetEndTime sets the end of the query range.
func (s *ListCandyDropAirdropRecordsService) SetEndTime(endTime time.Time) *ListCandyDropAirdropRecordsService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetPage selects the page number (default 1).
func (s *ListCandyDropAirdropRecordsService) SetPage(page int) *ListCandyDropAirdropRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records per page (default 10, max 30).
func (s *ListCandyDropAirdropRecordsService) SetLimit(limit int) *ListCandyDropAirdropRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListCandyDropAirdropRecordsService) Do(ctx context.Context) ([]CandyDropAirdropRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/launch/candydrop/airdrop-records", s.params).WithSign()
	resp, err := request.Do[[]CandyDropAirdropRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CandyDropActivity is one CandyDrop activity. start_time / end_time are
// UTC-formatted date strings (layout unconfirmed live); rule_name lists the
// task types; user_max_rewards is the most one user can win.
type CandyDropActivity struct {
	ID             int64           `json:"id"`
	Currency       string          `json:"currency"`
	TotalRewards   decimal.Decimal `json:"total_rewards"`
	StartTime      string          `json:"start_time"`
	EndTime        string          `json:"end_time"`
	RuleName       []string        `json:"rule_name"`
	Participants   int64           `json:"participants"`
	UserMaxRewards decimal.Decimal `json:"user_max_rewards"`
}

// CandyDropActivityRules is the rule set of one CandyDrop activity.
// start_time / end_time are UTC-formatted date strings (layout unconfirmed
// live).
type CandyDropActivityRules struct {
	Currency     string               `json:"currency"`
	TotalRewards decimal.Decimal      `json:"total_rewards"`
	StartTime    string               `json:"start_time"`
	EndTime      string               `json:"end_time"`
	PrizePools   []CandyDropPrizePool `json:"prize_pools"`
}

// CandyDropPrizePool is one prize pool of a CandyDrop activity: its total
// amount, the per-user cap and the tasks that earn a share.
type CandyDropPrizePool struct {
	PrizePoolType int             `json:"prize_pool_type"`
	PrizeAll      decimal.Decimal `json:"prize_all"`
	PrizeLimit    decimal.Decimal `json:"prize_limit"`
	Tasks         []CandyDropTask `json:"tasks"`
}

// CandyDropTask is one task of a CandyDrop prize pool. task_desc and
// exclusive_label are "-" when not configured.
type CandyDropTask struct {
	TaskName       string `json:"task_name"`
	TaskDesc       string `json:"task_desc"`
	ExclusiveLabel string `json:"exclusive_label"`
}

// CandyDropRegisterResult is the result of a CandyDrop registration.
type CandyDropRegisterResult struct {
	Success bool `json:"success"`
}

// CandyDropTaskProgress is the user's task progress in one CandyDrop activity.
// start_time / end_time are UTC-formatted date strings (layout unconfirmed
// live).
type CandyDropTaskProgress struct {
	Currency     string                      `json:"currency"`
	TotalRewards decimal.Decimal             `json:"total_rewards"`
	StartTime    string                      `json:"start_time"`
	EndTime      string                      `json:"end_time"`
	Tasks        []CandyDropTaskProgressItem `json:"tasks"`
}

// CandyDropTaskProgressItem is the progress of one registered task. value is
// the task's progress figure (trading volume, deposit amount, ...), kept as
// the raw string.
type CandyDropTaskProgressItem struct {
	TaskName string `json:"task_name"`
	TaskDesc string `json:"task_desc"`
	Value    string `json:"value"`
}

// CandyDropParticipationRecord is one CandyDrop participation. status is
// ongoing, awaiting_draw, won or not_win; register_time is a UTC-formatted
// date string (layout unconfirmed live).
type CandyDropParticipationRecord struct {
	ID           int64  `json:"id"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
	RegisterTime string `json:"register_time"`
}

// CandyDropAirdropRecord is one CandyDrop airdrop received. convert_amount is
// the auto-exchanged amount; airdrop_time is a UTC-formatted date string
// (layout unconfirmed live).
type CandyDropAirdropRecord struct {
	Currency      string          `json:"currency"`
	AirdropTime   string          `json:"airdrop_time"`
	Rewards       decimal.Decimal `json:"rewards"`
	ConvertAmount decimal.Decimal `json:"convert_amount"`
}
