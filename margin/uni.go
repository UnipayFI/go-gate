package margin

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListUniCurrencyPairsService -- GET /api/v4/margin/uni/currency_pairs
//
// Returns every unified (isolated) margin lending market and its borrow rules.
type ListUniCurrencyPairsService struct {
	c *MarginClient
}

func (c *MarginClient) NewListUniCurrencyPairsService() *ListUniCurrencyPairsService {
	return &ListUniCurrencyPairsService{c: c}
}

func (s *ListUniCurrencyPairsService) Do(ctx context.Context) ([]UniCurrencyPair, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/currency_pairs")
	resp, err := request.Do[[]UniCurrencyPair](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetUniCurrencyPairService -- GET /api/v4/margin/uni/currency_pairs/{currency_pair}
//
// Returns the borrow rules for a single unified-margin lending market.
type GetUniCurrencyPairService struct {
	c            *MarginClient
	currencyPair string
}

func (c *MarginClient) NewGetUniCurrencyPairService(currencyPair string) *GetUniCurrencyPairService {
	return &GetUniCurrencyPairService{c: c, currencyPair: currencyPair}
}

func (s *GetUniCurrencyPairService) Do(ctx context.Context) (*UniCurrencyPair, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/currency_pairs/"+s.currencyPair)
	return request.Do[UniCurrencyPair](req)
}

// UniCurrencyPair is a unified (isolated) margin lending market and its borrow
// rules.
type UniCurrencyPair struct {
	CurrencyPair         string          `json:"currency_pair"`
	BaseMinBorrowAmount  decimal.Decimal `json:"base_min_borrow_amount"`
	QuoteMinBorrowAmount decimal.Decimal `json:"quote_min_borrow_amount"`
	Leverage             decimal.Decimal `json:"leverage"`
}

// GetMarginUniEstimateRateService -- GET /api/v4/margin/uni/estimate_rate (private)
//
// Estimates the current hourly interest rate for the requested currencies (max
// 10). Rates change hourly with lending depth, so the result is an estimate.
type GetMarginUniEstimateRateService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetMarginUniEstimateRateService(currencies []string) *GetMarginUniEstimateRateService {
	return &GetMarginUniEstimateRateService{c: c, params: map[string]string{
		"currencies": strings.Join(currencies, ","),
	}}
}

func (s *GetMarginUniEstimateRateService) Do(ctx context.Context) (MarginUniEstimateRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/estimate_rate", s.params).WithSign()
	resp, err := request.Do[MarginUniEstimateRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginUniEstimateRate maps each queried currency to its estimated hourly
// interest rate.
type MarginUniEstimateRate = map[string]decimal.Decimal

// ListUniLoansService -- GET /api/v4/margin/uni/loans (private)
//
// Returns the account's outstanding unified-margin borrowings.
type ListUniLoansService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListUniLoansService() *ListUniLoansService {
	return &ListUniLoansService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *ListUniLoansService) SetCurrencyPair(currencyPair string) *ListUniLoansService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetCurrency narrows the result to a single currency.
func (s *ListUniLoansService) SetCurrency(currency string) *ListUniLoansService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the page number (1-based).
func (s *ListUniLoansService) SetPage(page int) *ListUniLoansService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUniLoansService) SetLimit(limit int) *ListUniLoansService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListUniLoansService) Do(ctx context.Context) ([]UniLoan, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/loans", s.params).WithSign()
	resp, err := request.Do[[]UniLoan](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniLoan is a single outstanding unified-margin borrowing.
type UniLoan struct {
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	Amount       decimal.Decimal `json:"amount"`
	Type         string          `json:"type"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	UpdateTime   time.Time       `json:"update_time,format:unix"`
}

// CreateUniLoanService -- POST /api/v4/margin/uni/loans (private)
//
// Borrows or repays under unified (isolated) margin. type is "borrow" or
// "repay". Returns no content.
type CreateUniLoanService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewCreateUniLoanService(currencyPair, currency string, amount decimal.Decimal, loanType string) *CreateUniLoanService {
	return &CreateUniLoanService{c: c, body: map[string]any{
		"currency_pair": currencyPair,
		"currency":      currency,
		"amount":        amount.String(),
		"type":          loanType,
	}}
}

// SetRepaidAll repays the full outstanding amount, overriding amount. For repay
// operations only.
func (s *CreateUniLoanService) SetRepaidAll(repaidAll bool) *CreateUniLoanService {
	s.body["repaid_all"] = repaidAll
	return s
}

func (s *CreateUniLoanService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/margin/uni/loans", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListUniLoanRecordsService -- GET /api/v4/margin/uni/loan_records (private)
//
// Returns the account's unified-margin borrow/repay history.
type ListUniLoanRecordsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListUniLoanRecordsService() *ListUniLoanRecordsService {
	return &ListUniLoanRecordsService{c: c, params: map[string]string{}}
}

// SetType filters by record type ("borrow" or "repay").
func (s *ListUniLoanRecordsService) SetType(recordType string) *ListUniLoanRecordsService {
	s.params["type"] = recordType
	return s
}

// SetCurrency narrows the result to a single currency.
func (s *ListUniLoanRecordsService) SetCurrency(currency string) *ListUniLoanRecordsService {
	s.params["currency"] = currency
	return s
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *ListUniLoanRecordsService) SetCurrencyPair(currencyPair string) *ListUniLoanRecordsService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetPage selects the page number (1-based).
func (s *ListUniLoanRecordsService) SetPage(page int) *ListUniLoanRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUniLoanRecordsService) SetLimit(limit int) *ListUniLoanRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListUniLoanRecordsService) Do(ctx context.Context) ([]UniLoanRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/loan_records", s.params).WithSign()
	resp, err := request.Do[[]UniLoanRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniLoanRecord is a single unified-margin borrow or repay event.
type UniLoanRecord struct {
	Type         string          `json:"type"`
	CurrencyPair string          `json:"currency_pair"`
	Currency     string          `json:"currency"`
	Amount       decimal.Decimal `json:"amount"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
}

// ListUniLoanInterestRecordsService -- GET /api/v4/margin/uni/interest_records (private)
//
// Returns the account's unified-margin interest deduction history.
type ListUniLoanInterestRecordsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListUniLoanInterestRecordsService() *ListUniLoanInterestRecordsService {
	return &ListUniLoanInterestRecordsService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *ListUniLoanInterestRecordsService) SetCurrencyPair(currencyPair string) *ListUniLoanInterestRecordsService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetCurrency narrows the result to a single currency.
func (s *ListUniLoanInterestRecordsService) SetCurrency(currency string) *ListUniLoanInterestRecordsService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the page number (1-based).
func (s *ListUniLoanInterestRecordsService) SetPage(page int) *ListUniLoanInterestRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned.
func (s *ListUniLoanInterestRecordsService) SetLimit(limit int) *ListUniLoanInterestRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start time of the query window (unix seconds).
func (s *ListUniLoanInterestRecordsService) SetFrom(from time.Time) *ListUniLoanInterestRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window (unix seconds).
func (s *ListUniLoanInterestRecordsService) SetTo(to time.Time) *ListUniLoanInterestRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListUniLoanInterestRecordsService) Do(ctx context.Context) ([]UniLoanInterestRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/interest_records", s.params).WithSign()
	resp, err := request.Do[[]UniLoanInterestRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniLoanInterestRecord is a single unified-margin interest deduction.
type UniLoanInterestRecord struct {
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	ActualRate   decimal.Decimal `json:"actual_rate"`
	Interest     decimal.Decimal `json:"interest"`
	Status       int             `json:"status"`
	Type         string          `json:"type"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
}

// GetUniBorrowableService -- GET /api/v4/margin/uni/borrowable (private)
//
// Returns the maximum amount of a currency the account can still borrow within
// the given unified-margin pair.
type GetUniBorrowableService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetUniBorrowableService(currency, currencyPair string) *GetUniBorrowableService {
	return &GetUniBorrowableService{c: c, params: map[string]string{
		"currency":      currency,
		"currency_pair": currencyPair,
	}}
}

func (s *GetUniBorrowableService) Do(ctx context.Context) (*UniBorrowable, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/uni/borrowable", s.params).WithSign()
	return request.Do[UniBorrowable](req)
}

// UniBorrowable is the maximum borrowable amount of a currency within a
// unified-margin pair.
type UniBorrowable struct {
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	Borrowable   decimal.Decimal `json:"borrowable"`
}
