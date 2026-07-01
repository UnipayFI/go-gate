package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// SwapETH2Service -- POST /api/v4/earn/staking/eth2/swap (private)
//
// Swaps between ETH and ETH2 staking tokens. side is "1" for a forward swap
// (ETH -> ETH2) and "2" for a reverse swap (ETH2 -> ETH).
type SwapETH2Service struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewSwapETH2Service(side string, amount decimal.Decimal) *SwapETH2Service {
	return &SwapETH2Service{c: c, body: map[string]any{
		"side":   side,
		"amount": amount.String(),
	}}
}

func (s *SwapETH2Service) Do(ctx context.Context) (*ETH2SwapResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/staking/eth2/swap", s.body).WithSign()
	return request.Do[ETH2SwapResult](req)
}

// ETH2SwapResult is the response of an ETH2 swap. Gate returns an empty body on
// success — the 2xx status is the confirmation — so this struct carries no
// fields.
type ETH2SwapResult struct{}

// RateListETH2Service -- GET /api/v4/earn/staking/eth2/rate_records (private)
//
// Returns the ETH2 staking return-rate records for the last 31 days.
type RateListETH2Service struct {
	c *EarnClient
}

func (c *EarnClient) NewRateListETH2Service() *RateListETH2Service {
	return &RateListETH2Service{c: c}
}

func (s *RateListETH2Service) Do(ctx context.Context) (*ETH2RateList, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/staking/eth2/rate_records").WithSign()
	return request.Do[ETH2RateList](req)
}

// ETH2RateList wraps the daily ETH2 staking rate history (the response is a
// {"rates":[...]} object, not a bare array).
type ETH2RateList struct {
	Rates []ETH2RateRecord `json:"rates"`
}

// ETH2RateRecord is one day's ETH2 staking return rate.
type ETH2RateRecord struct {
	DateTime time.Time       `json:"date_time,format:unix"`
	Date     string          `json:"date"`
	Rate     decimal.Decimal `json:"rate"`
}

// ListDualInvestmentPlansService -- GET /api/v4/earn/dual/investment_plan
//
// Returns the available Dual Investment products, optionally narrowed to a
// single plan.
type ListDualInvestmentPlansService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListDualInvestmentPlansService() *ListDualInvestmentPlansService {
	return &ListDualInvestmentPlansService{c: c, params: map[string]string{}}
}

// SetPlanID narrows the result to a single financial project.
func (s *ListDualInvestmentPlansService) SetPlanID(planID int64) *ListDualInvestmentPlansService {
	s.params["plan_id"] = strconv.FormatInt(planID, 10)
	return s
}

func (s *ListDualInvestmentPlansService) Do(ctx context.Context) ([]DualInvestmentPlan, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/investment_plan", s.params)
	resp, err := request.Do[[]DualInvestmentPlan](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DualInvestmentPlan is a single Dual Investment product and its terms.
type DualInvestmentPlan struct {
	ID               int64           `json:"id"`
	InstrumentName   string          `json:"instrument_name"`
	Type             string          `json:"type"`
	InvestCurrency   string          `json:"invest_currency"`
	ExerciseCurrency string          `json:"exercise_currency"`
	ExercisePrice    decimal.Decimal `json:"exercise_price"`
	DeliveryTime     time.Time       `json:"delivery_time,format:unix"`
	MinCopies        int             `json:"min_copies"`
	MaxCopies        int             `json:"max_copies"`
	PerValue         decimal.Decimal `json:"per_value"`
	MinAmount        decimal.Decimal `json:"min_amount"`
	APYDisplay       decimal.Decimal `json:"apy_display"`
	StartTime        time.Time       `json:"start_time,format:unix"`
	EndTime          time.Time       `json:"end_time,format:unix"`
	Status           string          `json:"status"`
}

// ListDualOrdersService -- GET /api/v4/earn/dual/orders (private)
//
// Returns the authenticated user's Dual Investment orders.
type ListDualOrdersService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListDualOrdersService() *ListDualOrdersService {
	return &ListDualOrdersService{c: c, params: map[string]string{}}
}

// SetFrom bounds the result to orders settled at or after this time.
func (s *ListDualOrdersService) SetFrom(from time.Time) *ListDualOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo bounds the result to orders settled at or before this time.
func (s *ListDualOrdersService) SetTo(to time.Time) *ListDualOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetPage selects the result page (1-based).
func (s *ListDualOrdersService) SetPage(page int) *ListDualOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListDualOrdersService) SetLimit(limit int) *ListDualOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListDualOrdersService) Do(ctx context.Context) ([]DualOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/orders", s.params).WithSign()
	resp, err := request.Do[[]DualOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DualOrder is a single Dual Investment order of the authenticated user.
type DualOrder struct {
	ID                 int64           `json:"id"`
	PlanID             int64           `json:"plan_id"`
	Copies             decimal.Decimal `json:"copies"`
	InvestAmount       decimal.Decimal `json:"invest_amount"`
	SettlementAmount   decimal.Decimal `json:"settlement_amount"`
	CreateTime         time.Time       `json:"create_time,format:unix"`
	CompleteTime       time.Time       `json:"complete_time,format:unix"`
	Status             string          `json:"status"`
	InvestCurrency     string          `json:"invest_currency"`
	ExerciseCurrency   string          `json:"exercise_currency"`
	ExercisePrice      decimal.Decimal `json:"exercise_price"`
	SettlementPrice    decimal.Decimal `json:"settlement_price"`
	SettlementCurrency string          `json:"settlement_currency"`
	APYDisplay         decimal.Decimal `json:"apy_display"`
	APYSettlement      decimal.Decimal `json:"apy_settlement"`
	DeliveryTime       time.Time       `json:"delivery_time,format:unix"`
	Text               string          `json:"text"`
}

// PlaceDualOrderService -- POST /api/v4/earn/dual/orders (private)
//
// Subscribes to a Dual Investment product. amount is the subscription amount
// (mutually exclusive with the copies field on the wire).
type PlaceDualOrderService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewPlaceDualOrderService(planID string, amount decimal.Decimal) *PlaceDualOrderService {
	return &PlaceDualOrderService{c: c, body: map[string]any{
		"plan_id": planID,
		"amount":  amount.String(),
	}}
}

// SetText attaches a custom order id. It must start with "t-", be at most 28
// bytes after the prefix, and contain only letters, digits, "_", "-" or ".".
func (s *PlaceDualOrderService) SetText(text string) *PlaceDualOrderService {
	s.body["text"] = text
	return s
}

func (s *PlaceDualOrderService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/dual/orders", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListStructuredProductsService -- GET /api/v4/earn/structured/products
//
// Returns the Structured Product list for a given status ("in_process",
// "will_begin", "wait_settlement" or "done").
type ListStructuredProductsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListStructuredProductsService(status string) *ListStructuredProductsService {
	return &ListStructuredProductsService{c: c, params: map[string]string{
		"status": status,
	}}
}

// SetType filters by product type (e.g. "SharkFin2.0", "BullishSharkFin",
// "BearishSharkFin", "DoubleNoTouch", "RangeAccrual", "SnowBall").
func (s *ListStructuredProductsService) SetType(productType string) *ListStructuredProductsService {
	s.params["type"] = productType
	return s
}

// SetPage selects the result page (1-based).
func (s *ListStructuredProductsService) SetPage(page int) *ListStructuredProductsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListStructuredProductsService) SetLimit(limit int) *ListStructuredProductsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListStructuredProductsService) Do(ctx context.Context) ([]StructuredProduct, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/structured/products", s.params)
	resp, err := request.Do[[]StructuredProduct](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// StructuredProduct is a single Structured Product and its terms.
type StructuredProduct struct {
	ID               int64           `json:"id"`
	Type             string          `json:"type"`
	NameEn           string          `json:"name_en"`
	RiskLabel        string          `json:"risk_label"`
	InvestmentCoin   string          `json:"investment_coin"`
	InvestmentPeriod string          `json:"investment_period"`
	MinAnnualRate    decimal.Decimal `json:"min_annual_rate"`
	MidAnnualRate    decimal.Decimal `json:"mid_annual_rate"`
	MaxAnnualRate    decimal.Decimal `json:"max_annual_rate"`
	WatchMarket      string          `json:"watch_market"`
	StartTime        time.Time       `json:"start_time,format:unix"`
	EndTime          time.Time       `json:"end_time,format:unix"`
	Status           string          `json:"status"`
}

// ListStructuredOrdersService -- GET /api/v4/earn/structured/orders (private)
//
// Returns the authenticated user's Structured Product orders.
type ListStructuredOrdersService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListStructuredOrdersService() *ListStructuredOrdersService {
	return &ListStructuredOrdersService{c: c, params: map[string]string{}}
}

// SetFrom bounds the result to orders at or after this time.
func (s *ListStructuredOrdersService) SetFrom(from time.Time) *ListStructuredOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo bounds the result to orders at or before this time.
func (s *ListStructuredOrdersService) SetTo(to time.Time) *ListStructuredOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetPage selects the result page (1-based).
func (s *ListStructuredOrdersService) SetPage(page int) *ListStructuredOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListStructuredOrdersService) SetLimit(limit int) *ListStructuredOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListStructuredOrdersService) Do(ctx context.Context) ([]StructuredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/structured/orders", s.params).WithSign()
	resp, err := request.Do[[]StructuredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// StructuredOrder is a single Structured Product order of the authenticated user.
type StructuredOrder struct {
	ID         int64           `json:"id"`
	Pid        string          `json:"pid"`
	LockCoin   string          `json:"lock_coin"`
	Amount     decimal.Decimal `json:"amount"`
	Status     string          `json:"status"`
	Income     decimal.Decimal `json:"income"`
	CreateTime time.Time       `json:"create_time,format:unix"`
}

// PlaceStructuredOrderService -- POST /api/v4/earn/structured/orders (private)
//
// Subscribes to a Structured Product. pid is the product id and amount is the
// buy quantity.
type PlaceStructuredOrderService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewPlaceStructuredOrderService(pid string, amount decimal.Decimal) *PlaceStructuredOrderService {
	return &PlaceStructuredOrderService{c: c, body: map[string]any{
		"pid":    pid,
		"amount": amount.String(),
	}}
}

func (s *PlaceStructuredOrderService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/structured/orders", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// FindCoinService -- GET /api/v4/earn/staking/coins
//
// Returns the coins available for staking and their reward terms.
type FindCoinService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewFindCoinService() *FindCoinService {
	return &FindCoinService{c: c, params: map[string]string{}}
}

// SetCoinType filters by currency type: "swap" (voucher), "lock" (locked
// position) or "debt" (US Treasury bond).
func (s *FindCoinService) SetCoinType(coinType string) *FindCoinService {
	s.params["cointype"] = coinType
	return s
}

func (s *FindCoinService) Do(ctx context.Context) ([]StakingCoin, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/staking/coins", s.params)
	resp, err := request.Do[[]StakingCoin](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// StakingCoin is a coin available for staking and its reward terms.
type StakingCoin struct {
	PID                 int64                   `json:"pid"`
	ProductType         int                     `json:"productType"`
	IsDefi              int                     `json:"isDefi"`
	Currency            string                  `json:"currency"`
	EstimateAPR         decimal.Decimal         `json:"estimateApr"`
	MinStakeAmount      decimal.Decimal         `json:"minStakeAmount"`
	MaxStakeAmount      decimal.Decimal         `json:"maxStakeAmount"`
	ProtocolName        string                  `json:"protocolName"`
	RedeemPeriod        int                     `json:"redeemPeriod"`
	ExchangeRate        decimal.Decimal         `json:"exchangeRate"`
	ExchangeRateReserve decimal.Decimal         `json:"exchangeRateReserve"`
	ExtraInterest       []StakingExtraInterest  `json:"extraInterest"`
	CurrencyRewards     []StakingCurrencyReward `json:"currencyRewards"`
}

// StakingExtraInterest is a bonus-interest campaign attached to a staking coin.
type StakingExtraInterest struct {
	StartTime       time.Time                `json:"start_time,format:unix"`
	EndTime         time.Time                `json:"end_time,format:unix"`
	RewardCoin      string                   `json:"reward_coin"`
	SegmentInterest []StakingSegmentInterest `json:"segment_interest"`
}

// StakingSegmentInterest is one tiered bonus-rate bracket keyed by staked amount.
type StakingSegmentInterest struct {
	MoneyMin  decimal.Decimal `json:"money_min"`
	MoneyMax  decimal.Decimal `json:"money_max"`
	MoneyRate decimal.Decimal `json:"money_rate"`
}

// StakingCurrencyReward is the reward currency and its annual rate for a staking
// coin.
type StakingCurrencyReward struct {
	APR               decimal.Decimal `json:"apr"`
	RewardCoin        string          `json:"reward_coin"`
	RewardDelayDays   int             `json:"reward_delay_days"`
	InterestDelayDays int             `json:"interest_delay_days"`
}

// SwapStakingCoinService -- POST /api/v4/earn/staking/swap (private)
//
// Stakes or redeems an on-chain staking coin. side is "0" to stake and "1" to
// redeem.
type SwapStakingCoinService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewSwapStakingCoinService(coin, side string, amount decimal.Decimal) *SwapStakingCoinService {
	return &SwapStakingCoinService{c: c, body: map[string]any{
		"coin":   coin,
		"side":   side,
		"amount": amount.String(),
	}}
}

// SetPID sets the DeFi-type mining protocol identifier.
func (s *SwapStakingCoinService) SetPID(pid int) *SwapStakingCoinService {
	s.body["pid"] = pid
	return s
}

func (s *SwapStakingCoinService) Do(ctx context.Context) (*StakingSwapResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/staking/swap", s.body).WithSign()
	return request.Do[StakingSwapResult](req)
}

// StakingSwapResult is the result of a staking stake/redeem operation.
type StakingSwapResult struct {
	ID             int64           `json:"id"`
	PID            int64           `json:"pid"`
	UID            int64           `json:"uid"`
	Coin           string          `json:"coin"`
	Type           int             `json:"type"`
	Subtype        string          `json:"subtype"`
	Amount         decimal.Decimal `json:"amount"`
	ExchangeRate   decimal.Decimal `json:"exchange_rate"`
	ExchangeAmount decimal.Decimal `json:"exchange_amount"`
	UpdateStamp    time.Time       `json:"updateStamp,format:unix"`
	CreateStamp    time.Time       `json:"createStamp,format:unix"`
	Status         int             `json:"status"`
	ProtocolType   int             `json:"protocol_type"`
	ClientOrderID  string          `json:"client_order_id"`
	Source         string          `json:"source"`
}
