package loan

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// MultiCollateralInput is one collateral currency and amount, used as an item of
// the collateral list when placing or adjusting a multi-currency collateral order.
type MultiCollateralInput struct {
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

// MultiRepayItemInput is one repayment instruction. Set RepaidAll to repay the
// full outstanding of the currency, otherwise Amount is used as a partial repay.
type MultiRepayItemInput struct {
	Currency  string          `json:"currency"`
	Amount    decimal.Decimal `json:"amount"`
	RepaidAll bool            `json:"repaid_all"`
}

// ListMultiCollateralOrdersService -- GET /api/v4/loan/multi_collateral/orders (private)
//
// Returns the authenticated account's multi-currency collateral loan orders.
type ListMultiCollateralOrdersService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListMultiCollateralOrdersService() *ListMultiCollateralOrdersService {
	return &ListMultiCollateralOrdersService{c: c, params: map[string]string{}}
}

// SetPage sets the page number.
func (s *ListMultiCollateralOrdersService) SetPage(page int) *ListMultiCollateralOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListMultiCollateralOrdersService) SetLimit(limit int) *ListMultiCollateralOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetSort selects the sort order: time_desc (default), ltv_asc or ltv_desc.
func (s *ListMultiCollateralOrdersService) SetSort(sort string) *ListMultiCollateralOrdersService {
	s.params["sort"] = sort
	return s
}

// SetOrderType filters by order type: current or fixed (defaults to current).
func (s *ListMultiCollateralOrdersService) SetOrderType(orderType string) *ListMultiCollateralOrdersService {
	s.params["order_type"] = orderType
	return s
}

func (s *ListMultiCollateralOrdersService) Do(ctx context.Context) ([]MultiCollateralOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/orders", s.params).WithSign()
	resp, err := request.Do[[]MultiCollateralOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateMultiCollateralService -- POST /api/v4/loan/multi_collateral/orders (private)
//
// Places a multi-currency collateral loan order borrowing borrowAmount of
// borrowCurrency against the supplied collateral currencies.
type CreateMultiCollateralService struct {
	c    *LoanClient
	body map[string]any
}

func (c *LoanClient) NewCreateMultiCollateralService(borrowCurrency string, borrowAmount decimal.Decimal, collaterals []MultiCollateralInput) *CreateMultiCollateralService {
	return &CreateMultiCollateralService{c: c, body: map[string]any{
		"borrow_currency":       borrowCurrency,
		"borrow_amount":         borrowAmount,
		"collateral_currencies": collaterals,
	}}
}

// SetOrderID sets a client-supplied order ID.
func (s *CreateMultiCollateralService) SetOrderID(orderID string) *CreateMultiCollateralService {
	s.body["order_id"] = orderID
	return s
}

// SetOrderType selects current (default) or fixed rate.
func (s *CreateMultiCollateralService) SetOrderType(orderType string) *CreateMultiCollateralService {
	s.body["order_type"] = orderType
	return s
}

// SetFixedType selects the fixed-rate lending period: 7d or 30d (fixed rate only).
func (s *CreateMultiCollateralService) SetFixedType(fixedType string) *CreateMultiCollateralService {
	s.body["fixed_type"] = fixedType
	return s
}

// SetFixedRate sets the fixed interest rate (fixed rate only).
func (s *CreateMultiCollateralService) SetFixedRate(fixedRate decimal.Decimal) *CreateMultiCollateralService {
	s.body["fixed_rate"] = fixedRate
	return s
}

// SetAutoRenew enables fixed-rate auto-renewal.
func (s *CreateMultiCollateralService) SetAutoRenew(autoRenew bool) *CreateMultiCollateralService {
	s.body["auto_renew"] = autoRenew
	return s
}

// SetAutoRepay enables fixed-rate auto-repayment.
func (s *CreateMultiCollateralService) SetAutoRepay(autoRepay bool) *CreateMultiCollateralService {
	s.body["auto_repay"] = autoRepay
	return s
}

func (s *CreateMultiCollateralService) Do(ctx context.Context) (*MultiCreateOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/loan/multi_collateral/orders", s.body).WithSign()
	return request.Do[MultiCreateOrderResult](req)
}

// GetMultiCollateralOrderDetailService -- GET /api/v4/loan/multi_collateral/orders/{order_id} (private)
//
// Returns the details of a single multi-currency collateral order.
type GetMultiCollateralOrderDetailService struct {
	c       *LoanClient
	orderID string
}

func (c *LoanClient) NewGetMultiCollateralOrderDetailService(orderID string) *GetMultiCollateralOrderDetailService {
	return &GetMultiCollateralOrderDetailService{c: c, orderID: orderID}
}

func (s *GetMultiCollateralOrderDetailService) Do(ctx context.Context) (*MultiCollateralOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/orders/"+s.orderID).WithSign()
	return request.Do[MultiCollateralOrder](req)
}

// ListMultiRepayRecordsService -- GET /api/v4/loan/multi_collateral/repay (private)
//
// Returns multi-currency collateral repayment records for the given operation
// type (repay - regular repayment, liquidate - liquidation).
type ListMultiRepayRecordsService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListMultiRepayRecordsService(operationType string) *ListMultiRepayRecordsService {
	return &ListMultiRepayRecordsService{c: c, params: map[string]string{"type": operationType}}
}

// SetBorrowCurrency filters by borrowed currency.
func (s *ListMultiRepayRecordsService) SetBorrowCurrency(borrowCurrency string) *ListMultiRepayRecordsService {
	s.params["borrow_currency"] = borrowCurrency
	return s
}

// SetPage sets the page number.
func (s *ListMultiRepayRecordsService) SetPage(page int) *ListMultiRepayRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListMultiRepayRecordsService) SetLimit(limit int) *ListMultiRepayRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListMultiRepayRecordsService) SetFrom(from time.Time) *ListMultiRepayRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (defaults to now).
func (s *ListMultiRepayRecordsService) SetTo(to time.Time) *ListMultiRepayRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListMultiRepayRecordsService) Do(ctx context.Context) ([]MultiRepayRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/repay", s.params).WithSign()
	resp, err := request.Do[[]MultiRepayRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// RepayMultiCollateralLoanService -- POST /api/v4/loan/multi_collateral/repay (private)
//
// Repays a multi-currency collateral loan with the supplied repayment items.
type RepayMultiCollateralLoanService struct {
	c    *LoanClient
	body map[string]any
}

func (c *LoanClient) NewRepayMultiCollateralLoanService(orderID int64, repayItems []MultiRepayItemInput) *RepayMultiCollateralLoanService {
	return &RepayMultiCollateralLoanService{c: c, body: map[string]any{
		"order_id":    orderID,
		"repay_items": repayItems,
	}}
}

func (s *RepayMultiCollateralLoanService) Do(ctx context.Context) (*MultiRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/loan/multi_collateral/repay", s.body).WithSign()
	return request.Do[MultiRepayResult](req)
}

// ListMultiCollateralRecordsService -- GET /api/v4/loan/multi_collateral/mortgage (private)
//
// Returns the account's collateral adjustment (append / redeem) records.
type ListMultiCollateralRecordsService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListMultiCollateralRecordsService() *ListMultiCollateralRecordsService {
	return &ListMultiCollateralRecordsService{c: c, params: map[string]string{}}
}

// SetPage sets the page number.
func (s *ListMultiCollateralRecordsService) SetPage(page int) *ListMultiCollateralRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListMultiCollateralRecordsService) SetLimit(limit int) *ListMultiCollateralRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListMultiCollateralRecordsService) SetFrom(from time.Time) *ListMultiCollateralRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (defaults to now).
func (s *ListMultiCollateralRecordsService) SetTo(to time.Time) *ListMultiCollateralRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetCollateralCurrency filters by collateral currency.
func (s *ListMultiCollateralRecordsService) SetCollateralCurrency(collateralCurrency string) *ListMultiCollateralRecordsService {
	s.params["collateral_currency"] = collateralCurrency
	return s
}

func (s *ListMultiCollateralRecordsService) Do(ctx context.Context) ([]MultiCollateralRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/mortgage", s.params).WithSign()
	resp, err := request.Do[[]MultiCollateralRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OperateMultiCollateralService -- POST /api/v4/loan/multi_collateral/mortgage (private)
//
// Adds (append) or withdraws (redeem) collateral on a multi-currency collateral
// order.
type OperateMultiCollateralService struct {
	c    *LoanClient
	body map[string]any
}

func (c *LoanClient) NewOperateMultiCollateralService(orderID int64, operationType string) *OperateMultiCollateralService {
	return &OperateMultiCollateralService{c: c, body: map[string]any{
		"order_id": orderID,
		"type":     operationType,
	}}
}

// SetCollaterals sets the collateral currency list to adjust.
func (s *OperateMultiCollateralService) SetCollaterals(collaterals []MultiCollateralInput) *OperateMultiCollateralService {
	s.body["collaterals"] = collaterals
	return s
}

func (s *OperateMultiCollateralService) Do(ctx context.Context) (*MultiCollateralAdjustResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/loan/multi_collateral/mortgage", s.body).WithSign()
	return request.Do[MultiCollateralAdjustResult](req)
}

// ListUserCurrencyQuotaService -- GET /api/v4/loan/multi_collateral/currency_quota (private)
//
// Returns the user's collateral / borrowing currency quota information.
// quotaType is collateral or borrow; currency is comma-separated for collateral
// and a single currency for borrow.
type ListUserCurrencyQuotaService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListUserCurrencyQuotaService(quotaType, currency string) *ListUserCurrencyQuotaService {
	return &ListUserCurrencyQuotaService{c: c, params: map[string]string{
		"type":     quotaType,
		"currency": currency,
	}}
}

func (s *ListUserCurrencyQuotaService) Do(ctx context.Context) ([]MultiUserCurrencyQuota, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/currency_quota", s.params).WithSign()
	resp, err := request.Do[[]MultiUserCurrencyQuota](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListMultiCollateralCurrenciesService -- GET /api/v4/loan/multi_collateral/currencies
//
// Returns the borrowing and collateral currencies supported for multi-currency
// collateral loans.
type ListMultiCollateralCurrenciesService struct {
	c *LoanClient
}

func (c *LoanClient) NewListMultiCollateralCurrenciesService() *ListMultiCollateralCurrenciesService {
	return &ListMultiCollateralCurrenciesService{c: c}
}

func (s *ListMultiCollateralCurrenciesService) Do(ctx context.Context) (*MultiCollateralCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/currencies")
	return request.Do[MultiCollateralCurrency](req)
}

// GetMultiCollateralLtvService -- GET /api/v4/loan/multi_collateral/ltv
//
// Returns the multi-currency collateral ratio thresholds (fixed, independent of
// currency).
type GetMultiCollateralLtvService struct {
	c *LoanClient
}

func (c *LoanClient) NewGetMultiCollateralLtvService() *GetMultiCollateralLtvService {
	return &GetMultiCollateralLtvService{c: c}
}

func (s *GetMultiCollateralLtvService) Do(ctx context.Context) (*MultiCollateralLtv, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/ltv")
	return request.Do[MultiCollateralLtv](req)
}

// GetMultiCollateralFixRateService -- GET /api/v4/loan/multi_collateral/fixed_rate
//
// Returns each currency's 7-day and 30-day fixed interest rates.
type GetMultiCollateralFixRateService struct {
	c *LoanClient
}

func (c *LoanClient) NewGetMultiCollateralFixRateService() *GetMultiCollateralFixRateService {
	return &GetMultiCollateralFixRateService{c: c}
}

func (s *GetMultiCollateralFixRateService) Do(ctx context.Context) ([]MultiCollateralFixRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/fixed_rate")
	resp, err := request.Do[[]MultiCollateralFixRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetMultiCollateralCurrentRateService -- GET /api/v4/loan/multi_collateral/current_rate
//
// Returns the current (previous-hour) interest rate for the given currencies.
type GetMultiCollateralCurrentRateService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewGetMultiCollateralCurrentRateService(currencies []string) *GetMultiCollateralCurrentRateService {
	return &GetMultiCollateralCurrentRateService{c: c, params: map[string]string{
		"currencies": strings.Join(currencies, ","),
	}}
}

// SetVipLevel filters by VIP level (defaults to 0).
func (s *GetMultiCollateralCurrentRateService) SetVipLevel(vipLevel string) *GetMultiCollateralCurrentRateService {
	s.params["vip_level"] = vipLevel
	return s
}

func (s *GetMultiCollateralCurrentRateService) Do(ctx context.Context) ([]MultiCollateralCurrentRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/multi_collateral/current_rate", s.params)
	resp, err := request.Do[[]MultiCollateralCurrentRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MultiCreateOrderResult is the response of placing a multi-currency collateral order.
type MultiCreateOrderResult struct {
	OrderID int64 `json:"order_id"`
}

// MultiRepayResult is the response of a multi-currency collateral repayment.
type MultiRepayResult struct {
	OrderID          int64                      `json:"order_id"`
	RepaidCurrencies []MultiRepayCurrencyResult `json:"repaid_currencies"`
}

// MultiRepayCurrencyResult is the per-currency outcome of a repayment.
type MultiRepayCurrencyResult struct {
	Succeeded       bool            `json:"succeeded"`
	Label           string          `json:"label"`
	Message         string          `json:"message"`
	Currency        string          `json:"currency"`
	RepaidPrincipal decimal.Decimal `json:"repaid_principal"`
	RepaidInterest  decimal.Decimal `json:"repaid_interest"`
}

// MultiCollateralOrder is a multi-currency collateral loan order.
type MultiCollateralOrder struct {
	OrderID                 string                        `json:"order_id"`
	OrderType               string                        `json:"order_type"`
	FixedType               string                        `json:"fixed_type"`
	FixedRate               decimal.Decimal               `json:"fixed_rate"`
	ExpireTime              time.Time                     `json:"expire_time,format:unix"`
	AutoRenew               bool                          `json:"auto_renew"`
	AutoRepay               bool                          `json:"auto_repay"`
	CurrentLTV              decimal.Decimal               `json:"current_ltv"`
	Status                  string                        `json:"status"`
	BorrowTime              time.Time                     `json:"borrow_time,format:unix"`
	TotalLeftRepayUSDT      decimal.Decimal               `json:"total_left_repay_usdt"`
	TotalLeftCollateralUSDT decimal.Decimal               `json:"total_left_collateral_usdt"`
	BorrowCurrencies        []MultiBorrowCurrencyInfo     `json:"borrow_currencies"`
	CollateralCurrencies    []MultiCollateralCurrencyInfo `json:"collateral_currencies"`
}

// MultiBorrowCurrencyInfo is one borrowed currency of a multi-collateral order.
type MultiBorrowCurrencyInfo struct {
	Currency           string          `json:"currency"`
	IndexPrice         decimal.Decimal `json:"index_price"`
	LeftRepayPrincipal decimal.Decimal `json:"left_repay_principal"`
	LeftRepayInterest  decimal.Decimal `json:"left_repay_interest"`
	LeftRepayUSDT      decimal.Decimal `json:"left_repay_usdt"`
}

// MultiCollateralCurrencyInfo is one collateral currency of a multi-collateral order.
type MultiCollateralCurrencyInfo struct {
	Currency           string          `json:"currency"`
	IndexPrice         decimal.Decimal `json:"index_price"`
	LeftCollateral     decimal.Decimal `json:"left_collateral"`
	LeftCollateralUSDT decimal.Decimal `json:"left_collateral_usdt"`
}

// MultiRepayRecord is a multi-currency collateral repayment record.
type MultiRepayRecord struct {
	OrderID               int64                            `json:"order_id"`
	RecordID              int64                            `json:"record_id"`
	InitLTV               decimal.Decimal                  `json:"init_ltv"`
	BeforeLTV             decimal.Decimal                  `json:"before_ltv"`
	AfterLTV              decimal.Decimal                  `json:"after_ltv"`
	BorrowTime            time.Time                        `json:"borrow_time,format:unix"`
	RepayTime             time.Time                        `json:"repay_time,format:unix"`
	BorrowCurrencies      []MultiRepayRecordCurrency       `json:"borrow_currencies"`
	CollateralCurrencies  []MultiRepayRecordCurrency       `json:"collateral_currencies"`
	RepaidCurrencies      []MultiRepayRecordRepaidCurrency `json:"repaid_currencies"`
	TotalInterestList     []MultiRepayRecordInterest       `json:"total_interest_list"`
	LeftRepayInterestList []MultiRepayRecordLeftInterest   `json:"left_repay_interest_list"`
}

// MultiRepayRecordCurrency is a borrow/collateral currency snapshot in a repay record.
type MultiRepayRecordCurrency struct {
	Currency         string          `json:"currency"`
	IndexPrice       decimal.Decimal `json:"index_price"`
	BeforeAmount     decimal.Decimal `json:"before_amount"`
	BeforeAmountUSDT decimal.Decimal `json:"before_amount_usdt"`
	AfterAmount      decimal.Decimal `json:"after_amount"`
	AfterAmountUSDT  decimal.Decimal `json:"after_amount_usdt"`
}

// MultiRepayRecordRepaidCurrency is one repaid currency of a repay record.
type MultiRepayRecordRepaidCurrency struct {
	Currency         string          `json:"currency"`
	IndexPrice       decimal.Decimal `json:"index_price"`
	RepaidAmount     decimal.Decimal `json:"repaid_amount"`
	RepaidPrincipal  decimal.Decimal `json:"repaid_principal"`
	RepaidInterest   decimal.Decimal `json:"repaid_interest"`
	RepaidAmountUSDT decimal.Decimal `json:"repaid_amount_usdt"`
}

// MultiRepayRecordInterest is one currency's total interest in a repay record.
type MultiRepayRecordInterest struct {
	Currency   string          `json:"currency"`
	IndexPrice decimal.Decimal `json:"index_price"`
	Amount     decimal.Decimal `json:"amount"`
	AmountUSDT decimal.Decimal `json:"amount_usdt"`
}

// MultiRepayRecordLeftInterest is one currency's remaining interest in a repay record.
type MultiRepayRecordLeftInterest struct {
	Currency         string          `json:"currency"`
	IndexPrice       decimal.Decimal `json:"index_price"`
	BeforeAmount     decimal.Decimal `json:"before_amount"`
	BeforeAmountUSDT decimal.Decimal `json:"before_amount_usdt"`
	AfterAmount      decimal.Decimal `json:"after_amount"`
	AfterAmountUSDT  decimal.Decimal `json:"after_amount_usdt"`
}

// MultiCollateralRecord is a collateral adjustment record.
type MultiCollateralRecord struct {
	OrderID              int64                           `json:"order_id"`
	RecordID             int64                           `json:"record_id"`
	BeforeLTV            decimal.Decimal                 `json:"before_ltv"`
	AfterLTV             decimal.Decimal                 `json:"after_ltv"`
	OperateTime          time.Time                       `json:"operate_time,format:unix"`
	BorrowCurrencies     []MultiCollateralRecordCurrency `json:"borrow_currencies"`
	CollateralCurrencies []MultiCollateralRecordCurrency `json:"collateral_currencies"`
}

// MultiCollateralRecordCurrency is a currency snapshot in a collateral adjustment record.
type MultiCollateralRecordCurrency struct {
	Currency         string          `json:"currency"`
	IndexPrice       decimal.Decimal `json:"index_price"`
	BeforeAmount     decimal.Decimal `json:"before_amount"`
	BeforeAmountUSDT decimal.Decimal `json:"before_amount_usdt"`
	AfterAmount      decimal.Decimal `json:"after_amount"`
	AfterAmountUSDT  decimal.Decimal `json:"after_amount_usdt"`
}

// MultiCollateralAdjustResult is the response of adding or withdrawing collateral.
type MultiCollateralAdjustResult struct {
	OrderID              int64                              `json:"order_id"`
	CollateralCurrencies []MultiCollateralAdjustCurrencyRes `json:"collateral_currencies"`
}

// MultiCollateralAdjustCurrencyRes is the per-currency outcome of a collateral adjustment.
type MultiCollateralAdjustCurrencyRes struct {
	Succeeded bool            `json:"succeeded"`
	Label     string          `json:"label"`
	Message   string          `json:"message"`
	Currency  string          `json:"currency"`
	Amount    decimal.Decimal `json:"amount"`
}

// MultiUserCurrencyQuota is a user's collateral / borrowing quota for a currency.
type MultiUserCurrencyQuota struct {
	Currency      string          `json:"currency"`
	IndexPrice    decimal.Decimal `json:"index_price"`
	MinQuota      decimal.Decimal `json:"min_quota"`
	LeftQuota     decimal.Decimal `json:"left_quota"`
	LeftQuoteUSDT decimal.Decimal `json:"left_quote_usdt"`
}

// MultiCollateralCurrency lists the supported borrowing and collateral currencies.
type MultiCollateralCurrency struct {
	LoanCurrencies       []MultiLoanItem       `json:"loan_currencies"`
	CollateralCurrencies []MultiCollateralItem `json:"collateral_currencies"`
}

// MultiLoanItem is one supported borrowing currency and its latest price.
type MultiLoanItem struct {
	Currency string          `json:"currency"`
	Price    decimal.Decimal `json:"price"`
}

// MultiCollateralItem is one supported collateral currency, its index price and discount.
type MultiCollateralItem struct {
	Currency   string          `json:"currency"`
	IndexPrice decimal.Decimal `json:"index_price"`
	Discount   decimal.Decimal `json:"discount"`
}

// MultiCollateralLtv is the multi-currency collateral ratio thresholds.
type MultiCollateralLtv struct {
	InitLTV      decimal.Decimal `json:"init_ltv"`
	AlertLTV     decimal.Decimal `json:"alert_ltv"`
	LiquidateLTV decimal.Decimal `json:"liquidate_ltv"`
}

// MultiCollateralFixRate is a currency's 7-day and 30-day fixed interest rates.
type MultiCollateralFixRate struct {
	Currency   string          `json:"currency"`
	Rate7D     decimal.Decimal `json:"rate_7d"`
	Rate30D    decimal.Decimal `json:"rate_30d"`
	UpdateTime time.Time       `json:"update_time,format:unix"`
}

// MultiCollateralCurrentRate is a currency's current interest rate.
type MultiCollateralCurrentRate struct {
	Currency    string          `json:"currency"`
	CurrentRate decimal.Decimal `json:"current_rate"`
}
