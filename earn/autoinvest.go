package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// AutoInvestPlanItem is one currency allocation of an auto-invest portfolio; the
// ratios of all items in a request must sum to 100.
type AutoInvestPlanItem struct {
	Asset string `json:"asset"`
	Ratio string `json:"ratio"`
}

// CreateAutoInvestPlanService -- POST /api/v4/earn/autoinvest/plans/create (private)
//
// Creates an auto-invest plan that periodically buys a portfolio of currencies.
type CreateAutoInvestPlanService struct {
	c    *EarnClient
	body map[string]any
}

// NewCreateAutoInvestPlanService builds a plan-creation request. planMoney is the
// pricing currency (USDT or BTC), planAmount the per-period amount, planPeriodType
// one of daily/weekly/biweekly/monthly/hourly/4-hourly, planPeriodDay the cycle
// day, planPeriodHour the execution hour (0-23) and items the portfolio.
func (c *EarnClient) NewCreateAutoInvestPlanService(planMoney string, planAmount decimal.Decimal, planPeriodType string, planPeriodDay, planPeriodHour int64, items []AutoInvestPlanItem) *CreateAutoInvestPlanService {
	return &CreateAutoInvestPlanService{c: c, body: map[string]any{
		"plan_money":       planMoney,
		"plan_amount":      planAmount.String(),
		"plan_period_type": planPeriodType,
		"plan_period_day":  planPeriodDay,
		"plan_period_hour": planPeriodHour,
		"items":            items,
	}}
}

// SetPlanName sets the plan name (0-50 characters).
func (s *CreateAutoInvestPlanService) SetPlanName(planName string) *CreateAutoInvestPlanService {
	s.body["plan_name"] = planName
	return s
}

// SetPlanDes sets the plan description.
func (s *CreateAutoInvestPlanService) SetPlanDes(planDes string) *CreateAutoInvestPlanService {
	s.body["plan_des"] = planDes
	return s
}

// SetFundSource sets the fund source: "spot" or "earn" (default "spot").
func (s *CreateAutoInvestPlanService) SetFundSource(fundSource string) *CreateAutoInvestPlanService {
	s.body["fund_source"] = fundSource
	return s
}

// SetFundFlow sets the fund flow direction: "auto_invest" or "earn" (default
// "auto_invest").
func (s *CreateAutoInvestPlanService) SetFundFlow(fundFlow string) *CreateAutoInvestPlanService {
	s.body["fund_flow"] = fundFlow
	return s
}

// SetType sets the creation type: 0 for normal, 1 for quick investment.
func (s *CreateAutoInvestPlanService) SetType(planType int64) *CreateAutoInvestPlanService {
	s.body["type"] = planType
	return s
}

func (s *CreateAutoInvestPlanService) Do(ctx context.Context) (*AutoInvestPlanCreateResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/autoinvest/plans/create", s.body).WithSign()
	return request.Do[AutoInvestPlanCreateResult](req)
}

// AutoInvestPlanCreateResult is the created auto-invest plan summary.
type AutoInvestPlanCreateResult struct {
	ID         int64           `json:"id"`
	Amount     decimal.Decimal `json:"amount"`
	Money      string          `json:"money"`
	NextTime   time.Time       `json:"next_time,format:unix"`
	PeriodType string          `json:"period_type"`
	PeriodDay  int64           `json:"period_day"`
	PeriodHour int64           `json:"period_hour"`
	FundFlow   string          `json:"fund_flow"`
	FundSource string          `json:"fund_source"`
}

// UpdateAutoInvestPlanService -- POST /api/v4/earn/autoinvest/plans/update (private)
//
// Updates the fund source/flow of an existing auto-invest plan. Gate returns an
// empty body on success.
type UpdateAutoInvestPlanService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewUpdateAutoInvestPlanService(planID int64) *UpdateAutoInvestPlanService {
	return &UpdateAutoInvestPlanService{c: c, body: map[string]any{
		"plan_id": planID,
	}}
}

// SetFundSource sets the fund source: "spot" or "earn" (default "spot").
func (s *UpdateAutoInvestPlanService) SetFundSource(fundSource string) *UpdateAutoInvestPlanService {
	s.body["fund_source"] = fundSource
	return s
}

// SetFundFlow sets the fund flow direction: "auto_invest" or "earn" (default
// "auto_invest").
func (s *UpdateAutoInvestPlanService) SetFundFlow(fundFlow string) *UpdateAutoInvestPlanService {
	s.body["fund_flow"] = fundFlow
	return s
}

func (s *UpdateAutoInvestPlanService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/autoinvest/plans/update", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// StopAutoInvestPlanService -- POST /api/v4/earn/autoinvest/plans/stop (private)
//
// Stops an auto-invest plan. Gate returns an empty body on success.
type StopAutoInvestPlanService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewStopAutoInvestPlanService(planID int64) *StopAutoInvestPlanService {
	return &StopAutoInvestPlanService{c: c, body: map[string]any{
		"plan_id": planID,
	}}
}

func (s *StopAutoInvestPlanService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/autoinvest/plans/stop", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// AddAutoInvestPositionService -- POST /api/v4/earn/autoinvest/plans/add_position (private)
//
// Immediately adds a position to an auto-invest plan. Gate returns an empty body
// on success.
type AddAutoInvestPositionService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewAddAutoInvestPositionService(planID int64, amount decimal.Decimal) *AddAutoInvestPositionService {
	return &AddAutoInvestPositionService{c: c, body: map[string]any{
		"plan_id": planID,
		"amount":  amount.String(),
	}}
}

func (s *AddAutoInvestPositionService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/autoinvest/plans/add_position", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListAutoInvestCoinsService -- GET /api/v4/earn/autoinvest/coins (private)
//
// Returns the currencies that support auto-invest for a given pricing currency.
type ListAutoInvestCoinsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListAutoInvestCoinsService() *ListAutoInvestCoinsService {
	return &ListAutoInvestCoinsService{c: c, params: map[string]string{}}
}

// SetPlanMoney sets the pricing currency: USDT or BTC (default USDT).
func (s *ListAutoInvestCoinsService) SetPlanMoney(planMoney string) *ListAutoInvestCoinsService {
	s.params["plan_money"] = planMoney
	return s
}

func (s *ListAutoInvestCoinsService) Do(ctx context.Context) ([]AutoInvestCoinsItem, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/autoinvest/coins", s.params).WithSign()
	resp, err := request.Do[[]AutoInvestCoinsItem](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AutoInvestCoinsItem is a currency that supports auto-invest.
type AutoInvestCoinsItem struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	AssetIconURL string `json:"asset_icon_url"`
	Sort         int64  `json:"sort"`
}

// GetAutoInvestMinAmountService -- POST /api/v4/earn/autoinvest/min_invest_amount (private)
//
// Computes the minimum investment amount for a given pricing currency and
// portfolio.
type GetAutoInvestMinAmountService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewGetAutoInvestMinAmountService(money string, items []AutoInvestPlanItem) *GetAutoInvestMinAmountService {
	return &GetAutoInvestMinAmountService{c: c, body: map[string]any{
		"money": money,
		"items": items,
	}}
}

func (s *GetAutoInvestMinAmountService) Do(ctx context.Context) (*AutoInvestMinAmount, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/autoinvest/min_invest_amount", s.body).WithSign()
	return request.Do[AutoInvestMinAmount](req)
}

// AutoInvestMinAmount is the minimum investment amount for a portfolio.
type AutoInvestMinAmount struct {
	MinAmount decimal.Decimal `json:"min_amount"`
}

// ListAutoInvestRecordsService -- GET /api/v4/earn/autoinvest/plans/records (private)
//
// Returns the execution records of an auto-invest plan, paginated.
type ListAutoInvestRecordsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListAutoInvestRecordsService(planID int64) *ListAutoInvestRecordsService {
	return &ListAutoInvestRecordsService{c: c, params: map[string]string{
		"plan_id": strconv.FormatInt(planID, 10),
	}}
}

// SetPage selects the result page (1-based).
func (s *ListAutoInvestRecordsService) SetPage(page int64) *ListAutoInvestRecordsService {
	s.params["page"] = strconv.FormatInt(page, 10)
	return s
}

// SetPageSize caps the number of records per page (maximum 100).
func (s *ListAutoInvestRecordsService) SetPageSize(pageSize int64) *ListAutoInvestRecordsService {
	s.params["page_size"] = strconv.FormatInt(pageSize, 10)
	return s
}

func (s *ListAutoInvestRecordsService) Do(ctx context.Context) (*AutoInvestRecordList, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/autoinvest/plans/records", s.params).WithSign()
	return request.Do[AutoInvestRecordList](req)
}

// AutoInvestRecordList is one page of auto-invest plan execution records.
type AutoInvestRecordList struct {
	Page      int64              `json:"page"`
	PageSize  int64              `json:"page_size"`
	TotalPage int64              `json:"total_page"`
	Total     int64              `json:"total"`
	List      []AutoInvestRecord `json:"list"`
}

// AutoInvestRecord is a single auto-invest plan execution record.
type AutoInvestRecord struct {
	ID            int64           `json:"id"`
	Type          string          `json:"type"`
	Money         string          `json:"money"`
	UserID        int64           `json:"user_id"`
	PlanID        int64           `json:"plan_id"`
	PlanVersion   int64           `json:"plan_version"`
	Amount        decimal.Decimal `json:"amount"`
	CreateTime    time.Time       `json:"create_time,format:unix"`
	UpdateTime    time.Time       `json:"update_time,format:unix"`
	Status        string          `json:"status"`
	StatusType    int64           `json:"status_type"`
	Side          int64           `json:"side"`
	StatusMessage string          `json:"status_message"`
	Detail        string          `json:"detail"`
	Asset         string          `json:"asset"`
}

// ListAutoInvestOrdersService -- GET /api/v4/earn/autoinvest/orders (private)
//
// Returns the order-level details of a single auto-invest plan execution record.
type ListAutoInvestOrdersService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListAutoInvestOrdersService(planID, recordID int64) *ListAutoInvestOrdersService {
	return &ListAutoInvestOrdersService{c: c, params: map[string]string{
		"plan_id":   strconv.FormatInt(planID, 10),
		"record_id": strconv.FormatInt(recordID, 10),
	}}
}

func (s *ListAutoInvestOrdersService) Do(ctx context.Context) ([]AutoInvestOrderItem, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/autoinvest/orders", s.params).WithSign()
	resp, err := request.Do[[]AutoInvestOrderItem](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AutoInvestOrderItem is a single auto-invest execution order.
type AutoInvestOrderItem struct {
	ID         int64           `json:"id"`
	Type       string          `json:"type"`
	Amount     decimal.Decimal `json:"amount"`
	PlanID     int64           `json:"plan_id"`
	Side       int64           `json:"side"`
	Asset      string          `json:"asset"`
	RecordID   int64           `json:"record_id"`
	TotalMoney decimal.Decimal `json:"total_money"`
	Market     string          `json:"market"`
	Price      decimal.Decimal `json:"price"`
	CreateTime time.Time       `json:"create_time,format:unix"`
	Total      decimal.Decimal `json:"total"`
	FundFlow   string          `json:"fund_flow"`
	ErrorCode  int64           `json:"error_code"`
	ErrorMsg   string          `json:"error_msg"`
	Status     int64           `json:"status"`
}

// ListAutoInvestConfigService -- GET /api/v4/earn/autoinvest/config (private)
//
// Returns the per-currency auto-invest configuration (investment limits).
type ListAutoInvestConfigService struct {
	c *EarnClient
}

func (c *EarnClient) NewListAutoInvestConfigService() *ListAutoInvestConfigService {
	return &ListAutoInvestConfigService{c: c}
}

func (s *ListAutoInvestConfigService) Do(ctx context.Context) ([]AutoInvestConfigItem, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/autoinvest/config").WithSign()
	resp, err := request.Do[[]AutoInvestConfigItem](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AutoInvestConfigItem is a single currency's auto-invest configuration.
type AutoInvestConfigItem struct {
	Coin     string          `json:"coin"`
	MaxLimit decimal.Decimal `json:"max_limit"`
}

// GetAutoInvestPlanDetailService -- GET /api/v4/earn/autoinvest/plans/detail (private)
//
// Returns the details of a single auto-invest plan.
type GetAutoInvestPlanDetailService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetAutoInvestPlanDetailService(planID int64) *GetAutoInvestPlanDetailService {
	return &GetAutoInvestPlanDetailService{c: c, params: map[string]string{
		"plan_id": strconv.FormatInt(planID, 10),
	}}
}

func (s *GetAutoInvestPlanDetailService) Do(ctx context.Context) (*AutoInvestPlanDetail, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/autoinvest/plans/detail", s.params).WithSign()
	return request.Do[AutoInvestPlanDetail](req)
}

// AutoInvestPlanDetail is the full configuration and state of an auto-invest
// plan.
type AutoInvestPlanDetail struct {
	ID         int64                     `json:"id"`
	Version    int64                     `json:"version"`
	Name       string                    `json:"name"`
	CreateTime time.Time                 `json:"create_time,format:unix"`
	UpdateTime time.Time                 `json:"update_time,format:unix"`
	UserID     int64                     `json:"user_id"`
	Money      string                    `json:"money"`
	Amount     decimal.Decimal           `json:"amount"`
	PeriodType string                    `json:"period_type"`
	PeriodDay  int64                     `json:"period_day"`
	PeriodHour int64                     `json:"period_hour"`
	AssetType  int64                     `json:"asset_type"`
	Portfolio  []AutoInvestPortfolioItem `json:"portfolio"`
	NextTime   time.Time                 `json:"next_time,format:unix"`
	Period     int64                     `json:"period"`
	FundSource string                    `json:"fund_source"`
	FundFlow   string                    `json:"fund_flow"`
}

// AutoInvestPortfolioItem is one currency allocation of an auto-invest plan with
// its cumulative statistics.
type AutoInvestPortfolioItem struct {
	Asset        string          `json:"asset"`
	Ratio        decimal.Decimal `json:"ratio"`
	CumInvest    decimal.Decimal `json:"cum_invest"`
	CumHold      decimal.Decimal `json:"cum_hold"`
	CumRedeem    decimal.Decimal `json:"cum_redeem"`
	AvgPrice     decimal.Decimal `json:"avg_price"`
	RedeemStatus int64           `json:"redeem_status"`
	LendAmount   decimal.Decimal `json:"lend_amount"`
}

// ListAutoInvestPlansService -- GET /api/v4/earn/autoinvest/plans/list_info (private)
//
// Returns the authenticated user's auto-invest plans by status, paginated.
type ListAutoInvestPlansService struct {
	c      *EarnClient
	params map[string]string
}

// NewListAutoInvestPlansService lists plans by status: "active" or "history".
func (c *EarnClient) NewListAutoInvestPlansService(status string) *ListAutoInvestPlansService {
	return &ListAutoInvestPlansService{c: c, params: map[string]string{
		"status": status,
	}}
}

// SetPage selects the result page (1-based).
func (s *ListAutoInvestPlansService) SetPage(page int64) *ListAutoInvestPlansService {
	s.params["page"] = strconv.FormatInt(page, 10)
	return s
}

// SetPageSize caps the number of records per page (maximum 100).
func (s *ListAutoInvestPlansService) SetPageSize(pageSize int64) *ListAutoInvestPlansService {
	s.params["page_size"] = strconv.FormatInt(pageSize, 10)
	return s
}

func (s *ListAutoInvestPlansService) Do(ctx context.Context) (*AutoInvestPlanList, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/autoinvest/plans/list_info", s.params).WithSign()
	return request.Do[AutoInvestPlanList](req)
}

// AutoInvestPlanList is one page of auto-invest plans; each entry has the same
// shape as AutoInvestPlanDetail.
type AutoInvestPlanList struct {
	Page       int64                  `json:"page"`
	PageSize   int64                  `json:"page_size"`
	PageCount  int64                  `json:"page_count"`
	TotalCount int64                  `json:"total_count"`
	List       []AutoInvestPlanDetail `json:"list"`
}
