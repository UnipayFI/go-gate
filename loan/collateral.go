package loan

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListCollateralLoanOrdersService -- GET /api/v4/loan/collateral/orders (private)
//
// Returns the authenticated account's collateral loan orders, optionally
// filtered by collateral / borrowed currency.
type ListCollateralLoanOrdersService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListCollateralLoanOrdersService() *ListCollateralLoanOrdersService {
	return &ListCollateralLoanOrdersService{c: c, params: map[string]string{}}
}

// SetCollateralCurrency narrows the result to a single collateral currency.
func (s *ListCollateralLoanOrdersService) SetCollateralCurrency(collateralCurrency string) *ListCollateralLoanOrdersService {
	s.params["collateral_currency"] = collateralCurrency
	return s
}

// SetBorrowCurrency narrows the result to a single borrowed currency.
func (s *ListCollateralLoanOrdersService) SetBorrowCurrency(borrowCurrency string) *ListCollateralLoanOrdersService {
	s.params["borrow_currency"] = borrowCurrency
	return s
}

// SetPage sets the page number.
func (s *ListCollateralLoanOrdersService) SetPage(page int) *ListCollateralLoanOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListCollateralLoanOrdersService) SetLimit(limit int) *ListCollateralLoanOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListCollateralLoanOrdersService) Do(ctx context.Context) ([]CollateralLoanOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/orders", s.params).WithSign()
	resp, err := request.Do[[]CollateralLoanOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateCollateralLoanService -- POST /api/v4/loan/collateral/orders (private)
//
// Places a collateral loan order: locks collateral_amount of collateral_currency
// and borrows borrow_amount of borrow_currency.
type CreateCollateralLoanService struct {
	c    *LoanClient
	body map[string]any
}

func (c *LoanClient) NewCreateCollateralLoanService(collateralAmount decimal.Decimal, collateralCurrency string, borrowAmount decimal.Decimal, borrowCurrency string) *CreateCollateralLoanService {
	return &CreateCollateralLoanService{c: c, body: map[string]any{
		"collateral_amount":   collateralAmount.String(),
		"collateral_currency": collateralCurrency,
		"borrow_amount":       borrowAmount.String(),
		"borrow_currency":     borrowCurrency,
	}}
}

func (s *CreateCollateralLoanService) Do(ctx context.Context) (*CollateralOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/loan/collateral/orders", s.body).WithSign()
	return request.Do[CollateralOrderResult](req)
}

// GetCollateralLoanOrderDetailService -- GET /api/v4/loan/collateral/orders/{order_id} (private)
//
// Returns the details of a single collateral loan order.
type GetCollateralLoanOrderDetailService struct {
	c       *LoanClient
	orderID int64
}

func (c *LoanClient) NewGetCollateralLoanOrderDetailService(orderID int64) *GetCollateralLoanOrderDetailService {
	return &GetCollateralLoanOrderDetailService{c: c, orderID: orderID}
}

func (s *GetCollateralLoanOrderDetailService) Do(ctx context.Context) (*CollateralLoanOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/orders/"+strconv.FormatInt(s.orderID, 10)).WithSign()
	return request.Do[CollateralLoanOrder](req)
}

// RepayCollateralLoanService -- POST /api/v4/loan/collateral/repay (private)
//
// Repays a collateral loan order. Set repaidAll for full repayment; otherwise
// repayAmount is the partial amount to repay.
type RepayCollateralLoanService struct {
	c    *LoanClient
	body map[string]any
}

func (c *LoanClient) NewRepayCollateralLoanService(orderID int64, repayAmount decimal.Decimal, repaidAll bool) *RepayCollateralLoanService {
	return &RepayCollateralLoanService{c: c, body: map[string]any{
		"order_id":     orderID,
		"repay_amount": repayAmount.String(),
		"repaid_all":   repaidAll,
	}}
}

func (s *RepayCollateralLoanService) Do(ctx context.Context) (*CollateralRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/loan/collateral/repay", s.body).WithSign()
	return request.Do[CollateralRepayResult](req)
}

// ListRepayRecordsService -- GET /api/v4/loan/collateral/repay_records (private)
//
// Returns collateral loan repayment records for the given source ("repay" for a
// regular repayment, "liquidate" for a liquidation).
type ListRepayRecordsService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListRepayRecordsService(source string) *ListRepayRecordsService {
	return &ListRepayRecordsService{c: c, params: map[string]string{"source": source}}
}

// SetBorrowCurrency narrows the result to a single borrowed currency.
func (s *ListRepayRecordsService) SetBorrowCurrency(borrowCurrency string) *ListRepayRecordsService {
	s.params["borrow_currency"] = borrowCurrency
	return s
}

// SetCollateralCurrency narrows the result to a single collateral currency.
func (s *ListRepayRecordsService) SetCollateralCurrency(collateralCurrency string) *ListRepayRecordsService {
	s.params["collateral_currency"] = collateralCurrency
	return s
}

// SetPage sets the page number.
func (s *ListRepayRecordsService) SetPage(page int) *ListRepayRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListRepayRecordsService) SetLimit(limit int) *ListRepayRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListRepayRecordsService) SetFrom(from time.Time) *ListRepayRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (defaults to now server-side).
func (s *ListRepayRecordsService) SetTo(to time.Time) *ListRepayRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListRepayRecordsService) Do(ctx context.Context) ([]CollateralRepayRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/repay_records", s.params).WithSign()
	resp, err := request.Do[[]CollateralRepayRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListCollateralRecordsService -- GET /api/v4/loan/collateral/collaterals (private)
//
// Returns the authenticated account's collateral adjustment records.
type ListCollateralRecordsService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListCollateralRecordsService() *ListCollateralRecordsService {
	return &ListCollateralRecordsService{c: c, params: map[string]string{}}
}

// SetPage sets the page number.
func (s *ListCollateralRecordsService) SetPage(page int) *ListCollateralRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListCollateralRecordsService) SetLimit(limit int) *ListCollateralRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListCollateralRecordsService) SetFrom(from time.Time) *ListCollateralRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (defaults to now server-side).
func (s *ListCollateralRecordsService) SetTo(to time.Time) *ListCollateralRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetBorrowCurrency narrows the result to a single borrowed currency.
func (s *ListCollateralRecordsService) SetBorrowCurrency(borrowCurrency string) *ListCollateralRecordsService {
	s.params["borrow_currency"] = borrowCurrency
	return s
}

// SetCollateralCurrency narrows the result to a single collateral currency.
func (s *ListCollateralRecordsService) SetCollateralCurrency(collateralCurrency string) *ListCollateralRecordsService {
	s.params["collateral_currency"] = collateralCurrency
	return s
}

func (s *ListCollateralRecordsService) Do(ctx context.Context) ([]CollateralRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/collaterals", s.params).WithSign()
	resp, err := request.Do[[]CollateralRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OperateCollateralService -- POST /api/v4/loan/collateral/collaterals (private)
//
// Adjusts an order's collateral: type "append" adds collateral, "redeem"
// withdraws it.
type OperateCollateralService struct {
	c    *LoanClient
	body map[string]any
}

func (c *LoanClient) NewOperateCollateralService(orderID int64, collateralCurrency string, collateralAmount decimal.Decimal, operationType string) *OperateCollateralService {
	return &OperateCollateralService{c: c, body: map[string]any{
		"order_id":            orderID,
		"collateral_currency": collateralCurrency,
		"collateral_amount":   collateralAmount.String(),
		"type":                operationType,
	}}
}

func (s *OperateCollateralService) Do(ctx context.Context) (*CollateralAdjustResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/loan/collateral/collaterals", s.body).WithSign()
	return request.Do[CollateralAdjustResult](req)
}

// GetUserTotalAmountService -- GET /api/v4/loan/collateral/total_amount (private)
//
// Returns the authenticated account's total borrowing and collateral amount
// (both denominated in USDT).
type GetUserTotalAmountService struct {
	c *LoanClient
}

func (c *LoanClient) NewGetUserTotalAmountService() *GetUserTotalAmountService {
	return &GetUserTotalAmountService{c: c}
}

func (s *GetUserTotalAmountService) Do(ctx context.Context) (*CollateralTotalAmount, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/total_amount").WithSign()
	return request.Do[CollateralTotalAmount](req)
}

// GetUserLtvInfoService -- GET /api/v4/loan/collateral/ltv (private)
//
// Returns the authenticated account's collateralization ratios and remaining
// borrowable amount for a collateral / borrowed currency pair.
type GetUserLtvInfoService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewGetUserLtvInfoService(collateralCurrency, borrowCurrency string) *GetUserLtvInfoService {
	return &GetUserLtvInfoService{c: c, params: map[string]string{
		"collateral_currency": collateralCurrency,
		"borrow_currency":     borrowCurrency,
	}}
}

func (s *GetUserLtvInfoService) Do(ctx context.Context) (*CollateralLtv, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/ltv", s.params).WithSign()
	return request.Do[CollateralLtv](req)
}

// ListCollateralCurrenciesService -- GET /api/v4/loan/collateral/currencies
//
// Returns the supported borrowing currencies and their eligible collateral
// currencies. When loan_currency is set, only that currency's entry is returned.
type ListCollateralCurrenciesService struct {
	c      *LoanClient
	params map[string]string
}

func (c *LoanClient) NewListCollateralCurrenciesService() *ListCollateralCurrenciesService {
	return &ListCollateralCurrenciesService{c: c, params: map[string]string{}}
}

// SetLoanCurrency narrows the result to a single borrowing currency.
func (s *ListCollateralCurrenciesService) SetLoanCurrency(loanCurrency string) *ListCollateralCurrenciesService {
	s.params["loan_currency"] = loanCurrency
	return s
}

func (s *ListCollateralCurrenciesService) Do(ctx context.Context) ([]CollateralCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v4/loan/collateral/currencies", s.params)
	resp, err := request.Do[[]CollateralCurrency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CollateralLoanOrder is a single collateral loan order.
type CollateralLoanOrder struct {
	OrderID            int64           `json:"order_id"`
	CollateralCurrency string          `json:"collateral_currency"`
	CollateralAmount   decimal.Decimal `json:"collateral_amount"`
	BorrowCurrency     string          `json:"borrow_currency"`
	BorrowAmount       decimal.Decimal `json:"borrow_amount"`
	RepaidAmount       decimal.Decimal `json:"repaid_amount"`
	RepaidPrincipal    decimal.Decimal `json:"repaid_principal"`
	RepaidInterest     decimal.Decimal `json:"repaid_interest"`
	InitLtv            decimal.Decimal `json:"init_ltv"`
	CurrentLtv         decimal.Decimal `json:"current_ltv"`
	LiquidateLtv       decimal.Decimal `json:"liquidate_ltv"`
	Status             string          `json:"status"`
	BorrowTime         time.Time       `json:"borrow_time,format:unix"`
	LeftRepayTotal     decimal.Decimal `json:"left_repay_total"`
	LeftRepayPrincipal decimal.Decimal `json:"left_repay_principal"`
	LeftRepayInterest  decimal.Decimal `json:"left_repay_interest"`
}

// CollateralOrderResult is the identifier returned when placing a collateral loan order.
type CollateralOrderResult struct {
	OrderID int64 `json:"order_id"`
}

// CollateralRepayResult is the principal / interest actually settled by a repayment.
type CollateralRepayResult struct {
	RepaidPrincipal decimal.Decimal `json:"repaid_principal"`
	RepaidInterest  decimal.Decimal `json:"repaid_interest"`
}

// CollateralRepayRecord is a single collateral loan repayment record.
type CollateralRepayRecord struct {
	OrderID              int64           `json:"order_id"`
	RecordID             int64           `json:"record_id"`
	RepaidAmount         decimal.Decimal `json:"repaid_amount"`
	BorrowCurrency       string          `json:"borrow_currency"`
	CollateralCurrency   string          `json:"collateral_currency"`
	InitLtv              decimal.Decimal `json:"init_ltv"`
	BorrowTime           time.Time       `json:"borrow_time,format:unix"`
	RepayTime            time.Time       `json:"repay_time,format:unix"`
	TotalInterest        decimal.Decimal `json:"total_interest"`
	BeforeLeftPrincipal  decimal.Decimal `json:"before_left_principal"`
	AfterLeftPrincipal   decimal.Decimal `json:"after_left_principal"`
	BeforeLeftCollateral decimal.Decimal `json:"before_left_collateral"`
	AfterLeftCollateral  decimal.Decimal `json:"after_left_collateral"`
}

// CollateralRecord is a single collateral adjustment record.
type CollateralRecord struct {
	OrderID            int64           `json:"order_id"`
	RecordID           int64           `json:"record_id"`
	BorrowCurrency     string          `json:"borrow_currency"`
	BorrowAmount       decimal.Decimal `json:"borrow_amount"`
	CollateralCurrency string          `json:"collateral_currency"`
	BeforeCollateral   decimal.Decimal `json:"before_collateral"`
	AfterCollateral    decimal.Decimal `json:"after_collateral"`
	BeforeLtv          decimal.Decimal `json:"before_ltv"`
	AfterLtv           decimal.Decimal `json:"after_ltv"`
	OperateTime        time.Time       `json:"operate_time,format:unix"`
}

// CollateralAdjustResult is the (empty) body returned by a collateral append / redeem.
type CollateralAdjustResult struct{}

// CollateralTotalAmount is the account's total borrowing and collateral amount,
// both denominated in USDT.
type CollateralTotalAmount struct {
	BorrowAmount     decimal.Decimal `json:"borrow_amount"`
	CollateralAmount decimal.Decimal `json:"collateral_amount"`
}

// CollateralLtv is an account's collateralization ratios and remaining
// borrowable amount for a collateral / borrowed currency pair.
type CollateralLtv struct {
	CollateralCurrency   string          `json:"collateral_currency"`
	BorrowCurrency       string          `json:"borrow_currency"`
	InitLtv              decimal.Decimal `json:"init_ltv"`
	AlertLtv             decimal.Decimal `json:"alert_ltv"`
	LiquidateLtv         decimal.Decimal `json:"liquidate_ltv"`
	MinBorrowAmount      decimal.Decimal `json:"min_borrow_amount"`
	LeftBorrowableAmount decimal.Decimal `json:"left_borrowable_amount"`
}

// CollateralCurrency is a borrowing currency and the collateral currencies it accepts.
type CollateralCurrency struct {
	LoanCurrency       string   `json:"loan_currency"`
	CollateralCurrency []string `json:"collateral_currency"`
}
