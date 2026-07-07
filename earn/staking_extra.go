package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListStakingOrdersService -- GET /api/v4/earn/staking/order_list (private)
//
// Returns the authenticated user's on-chain staking (stake/redeem) orders,
// paginated.
type ListStakingOrdersService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListStakingOrdersService() *ListStakingOrdersService {
	return &ListStakingOrdersService{c: c, params: map[string]string{}}
}

// SetPID filters by product id.
func (s *ListStakingOrdersService) SetPID(pid int) *ListStakingOrdersService {
	s.params["pid"] = strconv.Itoa(pid)
	return s
}

// SetCoin filters by currency name.
func (s *ListStakingOrdersService) SetCoin(coin string) *ListStakingOrdersService {
	s.params["coin"] = coin
	return s
}

// SetType filters by order type: 0 for staking, 1 for redemption.
func (s *ListStakingOrdersService) SetType(orderType int) *ListStakingOrdersService {
	s.params["type"] = strconv.Itoa(orderType)
	return s
}

// SetPage selects the result page (1-based).
func (s *ListStakingOrdersService) SetPage(page int) *ListStakingOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *ListStakingOrdersService) Do(ctx context.Context) (*StakingOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/staking/order_list", s.params).WithSign()
	return request.Do[StakingOrderList](req)
}

// StakingOrderList is one page of on-chain staking orders.
type StakingOrderList struct {
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	PageCount  int                `json:"pageCount"`
	TotalCount int                `json:"totalCount"`
	List       []StakingOrderItem `json:"list"`
}

// StakingOrderItem is a single on-chain staking (stake/redeem) order.
type StakingOrderItem struct {
	PID            int64           `json:"pid"`
	Coin           string          `json:"coin"`
	Amount         decimal.Decimal `json:"amount"`
	Type           int             `json:"type"`
	Status         int             `json:"status"`
	RedeemStamp    time.Time       `json:"redeem_stamp,format:unix"`
	CreateStamp    time.Time       `json:"createStamp,format:unix"`
	ExchangeAmount decimal.Decimal `json:"exchange_amount"`
	Fee            decimal.Decimal `json:"fee"`
}

// ListStakingAwardsService -- GET /api/v4/earn/staking/award_list (private)
//
// Returns the authenticated user's on-chain staking dividend (reward) records,
// paginated.
type ListStakingAwardsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListStakingAwardsService() *ListStakingAwardsService {
	return &ListStakingAwardsService{c: c, params: map[string]string{}}
}

// SetPID filters by product id.
func (s *ListStakingAwardsService) SetPID(pid int) *ListStakingAwardsService {
	s.params["pid"] = strconv.Itoa(pid)
	return s
}

// SetCoin filters by currency name.
func (s *ListStakingAwardsService) SetCoin(coin string) *ListStakingAwardsService {
	s.params["coin"] = coin
	return s
}

// SetPage selects the result page (1-based).
func (s *ListStakingAwardsService) SetPage(page int) *ListStakingAwardsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *ListStakingAwardsService) Do(ctx context.Context) (*StakingAwardList, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/staking/award_list", s.params).WithSign()
	return request.Do[StakingAwardList](req)
}

// StakingAwardList is one page of on-chain staking dividend records.
type StakingAwardList struct {
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	PageCount  int                `json:"pageCount"`
	TotalCount int                `json:"totalCount"`
	List       []StakingAwardItem `json:"list"`
}

// StakingAwardItem is a single on-chain staking dividend record.
type StakingAwardItem struct {
	PID              int64           `json:"pid"`
	MortgageCoin     string          `json:"mortgage_coin"`
	Amount           decimal.Decimal `json:"amount"`
	RewardCoin       string          `json:"reward_coin"`
	Interest         decimal.Decimal `json:"interest"`
	Fee              decimal.Decimal `json:"fee"`
	Status           int             `json:"status"`
	BonusDate        string          `json:"bonus_date"`
	ShouldBonusStamp time.Time       `json:"should_bonus_stamp,format:unix"`
}

// GetStakingAssetsService -- GET /api/v4/earn/staking/assets (private)
//
// Returns the authenticated user's aggregated on-chain staking assets (one entry
// per staked currency), optionally filtered by currency. The response is a JSON
// array.
type GetStakingAssetsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetStakingAssetsService() *GetStakingAssetsService {
	return &GetStakingAssetsService{c: c, params: map[string]string{}}
}

// SetCoin filters by currency name.
func (s *GetStakingAssetsService) SetCoin(coin string) *GetStakingAssetsService {
	s.params["coin"] = coin
	return s
}

func (s *GetStakingAssetsService) Do(ctx context.Context) ([]StakingAsset, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/staking/assets", s.params).WithSign()
	resp, err := request.Do[[]StakingAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// StakingAsset is the user's aggregated on-chain staking position and earnings.
type StakingAsset struct {
	PID                  int64                    `json:"pid"`
	MortgageCoin         string                   `json:"mortgage_coin"`
	MortgageAmount       decimal.Decimal          `json:"mortgage_amount"`
	CreateStamp          time.Time                `json:"createStamp,format:unix"`
	ExtraIncome          decimal.Decimal          `json:"extra_income"`
	FreezeAmount         decimal.Decimal          `json:"freeze_amount"`
	MoveIncome           decimal.Decimal          `json:"move_income"`
	Type                 int                      `json:"type"`
	Status               int                      `json:"status"`
	IncomeTotal          decimal.Decimal          `json:"income_total"`
	YesterdayIncomeMulti []any                    `json:"yesterday_income_multi"`
	RewardCoins          []StakingAssetRewardCoin `json:"reward_coins"`
	DefiIncome           StakingAssetDefiIncome   `json:"defi_income"`
}

// StakingAssetRewardCoin is a per-currency reward configuration of a staking
// asset.
type StakingAssetRewardCoin struct {
	RewardCoin        string `json:"reward_coin"`
	InterestDelayDays int    `json:"interest_delay_days"`
	RewardDelayDays   int    `json:"reward_delay_days"`
}

// StakingAssetDefiIncome is the DeFi earnings breakdown of a staking asset.
type StakingAssetDefiIncome struct {
	Total []StakingAssetDefiIncomeItem `json:"total"`
}

// StakingAssetDefiIncomeItem is a single currency's DeFi earning amount.
type StakingAssetDefiIncomeItem struct {
	Coin   string          `json:"coin"`
	Amount decimal.Decimal `json:"amount"`
}
