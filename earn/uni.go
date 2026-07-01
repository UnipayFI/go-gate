package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListUniCurrenciesService -- GET /api/v4/earn/uni/currencies
//
// Returns the list of currencies available for Uni-lending and their limits.
type ListUniCurrenciesService struct {
	c *EarnClient
}

func (c *EarnClient) NewListUniCurrenciesService() *ListUniCurrenciesService {
	return &ListUniCurrenciesService{c: c}
}

func (s *ListUniCurrenciesService) Do(ctx context.Context) ([]UniLendingCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/currencies")
	resp, err := request.Do[[]UniLendingCurrency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetUniCurrencyService -- GET /api/v4/earn/uni/currencies/{currency}
//
// Returns the Uni-lending limits for a single currency.
type GetUniCurrencyService struct {
	c        *EarnClient
	currency string
}

func (c *EarnClient) NewGetUniCurrencyService(currency string) *GetUniCurrencyService {
	return &GetUniCurrencyService{c: c, currency: currency}
}

func (s *GetUniCurrencyService) Do(ctx context.Context) (*UniLendingCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/currencies/"+s.currency)
	return request.Do[UniLendingCurrency](req)
}

// UniLendingCurrency is a currency that can be lent through Uni-lending and its
// per-order/rate limits.
type UniLendingCurrency struct {
	Currency      string          `json:"currency"`
	MinLendAmount decimal.Decimal `json:"min_lend_amount"`
	MaxLendAmount decimal.Decimal `json:"max_lend_amount"`
	MaxRate       decimal.Decimal `json:"max_rate"`
	MinRate       decimal.Decimal `json:"min_rate"`
}

// ListUserUniLendsService -- GET /api/v4/earn/uni/lends (private)
//
// Returns the authenticated user's Uni-lending orders.
type ListUserUniLendsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListUserUniLendsService() *ListUserUniLendsService {
	return &ListUserUniLendsService{c: c, params: map[string]string{}}
}

// SetCurrency filters the result to a single currency.
func (s *ListUserUniLendsService) SetCurrency(currency string) *ListUserUniLendsService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the result page (1-based).
func (s *ListUserUniLendsService) SetPage(page int) *ListUserUniLendsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUserUniLendsService) SetLimit(limit int) *ListUserUniLendsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListUserUniLendsService) Do(ctx context.Context) ([]UserUniLend, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/lends", s.params).WithSign()
	resp, err := request.Do[[]UserUniLend](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UserUniLend is a single Uni-lending order of the authenticated user.
type UserUniLend struct {
	Currency           string          `json:"currency"`
	CurrentAmount      decimal.Decimal `json:"current_amount"`
	Amount             decimal.Decimal `json:"amount"`
	LentAmount         decimal.Decimal `json:"lent_amount"`
	FrozenAmount       decimal.Decimal `json:"frozen_amount"`
	MinRate            decimal.Decimal `json:"min_rate"`
	InterestStatus     string          `json:"interest_status"`
	ReinvestLeftAmount decimal.Decimal `json:"reinvest_left_amount"`
	CreateTime         time.Time       `json:"create_time,format:unix"`
	UpdateTime         time.Time       `json:"update_time,format:unix"`
}

// CreateUniLendService -- POST /api/v4/earn/uni/lends (private)
//
// Lends funds into the Uni-lending pool or redeems previously lent funds.
type CreateUniLendService struct {
	c    *EarnClient
	body map[string]any
}

// NewCreateUniLendService creates a lend or redeem request. lendType is "lend"
// or "redeem".
func (c *EarnClient) NewCreateUniLendService(currency string, amount decimal.Decimal, lendType string) *CreateUniLendService {
	return &CreateUniLendService{c: c, body: map[string]any{
		"currency": currency,
		"amount":   amount.String(),
		"type":     lendType,
	}}
}

// SetMinRate sets the minimum acceptable hourly lending rate. Required when
// lending; if set too high, lending may fail and no interest is earned.
func (s *CreateUniLendService) SetMinRate(minRate decimal.Decimal) *CreateUniLendService {
	s.body["min_rate"] = minRate.String()
	return s
}

func (s *CreateUniLendService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/uni/lends", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ChangeUniLendService -- PATCH /api/v4/earn/uni/lends (private)
//
// Amends the minimum hourly interest rate of an existing Uni-lending order.
type ChangeUniLendService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewChangeUniLendService(currency string, minRate decimal.Decimal) *ChangeUniLendService {
	return &ChangeUniLendService{c: c, body: map[string]any{
		"currency": currency,
		"min_rate": minRate.String(),
	}}
}

func (s *ChangeUniLendService) Do(ctx context.Context) error {
	req := request.Patch(ctx, s.c, "/api/v4/earn/uni/lends", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListUniLendRecordsService -- GET /api/v4/earn/uni/lend_records (private)
//
// Returns the authenticated user's lend/redeem transaction records.
type ListUniLendRecordsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListUniLendRecordsService() *ListUniLendRecordsService {
	return &ListUniLendRecordsService{c: c, params: map[string]string{}}
}

// SetCurrency filters the result to a single currency.
func (s *ListUniLendRecordsService) SetCurrency(currency string) *ListUniLendRecordsService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the result page (1-based).
func (s *ListUniLendRecordsService) SetPage(page int) *ListUniLendRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUniLendRecordsService) SetLimit(limit int) *ListUniLendRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom bounds the result to records at or after this time.
func (s *ListUniLendRecordsService) SetFrom(from time.Time) *ListUniLendRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo bounds the result to records at or before this time.
func (s *ListUniLendRecordsService) SetTo(to time.Time) *ListUniLendRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetType filters by operation type: "lend" or "redeem".
func (s *ListUniLendRecordsService) SetType(lendType string) *ListUniLendRecordsService {
	s.params["type"] = lendType
	return s
}

func (s *ListUniLendRecordsService) Do(ctx context.Context) ([]UniLendRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/lend_records", s.params).WithSign()
	resp, err := request.Do[[]UniLendRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniLendRecord is a single lend/redeem transaction record.
type UniLendRecord struct {
	Currency         string          `json:"currency"`
	Amount           decimal.Decimal `json:"amount"`
	LastWalletAmount decimal.Decimal `json:"last_wallet_amount"`
	LastLentAmount   decimal.Decimal `json:"last_lent_amount"`
	LastFrozenAmount decimal.Decimal `json:"last_frozen_amount"`
	Type             string          `json:"type"`
	CreateTime       time.Time       `json:"create_time,format:unix"`
}

// GetUniInterestService -- GET /api/v4/earn/uni/interests/{currency} (private)
//
// Returns the authenticated user's total interest income for a currency.
type GetUniInterestService struct {
	c        *EarnClient
	currency string
}

func (c *EarnClient) NewGetUniInterestService(currency string) *GetUniInterestService {
	return &GetUniInterestService{c: c, currency: currency}
}

func (s *GetUniInterestService) Do(ctx context.Context) (*UniInterest, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/interests/"+s.currency).WithSign()
	return request.Do[UniInterest](req)
}

// UniInterest is a currency's accumulated Uni-lending interest income.
type UniInterest struct {
	Currency string          `json:"currency"`
	Interest decimal.Decimal `json:"interest"`
}

// ListUniInterestRecordsService -- GET /api/v4/earn/uni/interest_records (private)
//
// Returns the authenticated user's hourly interest (dividend) records.
type ListUniInterestRecordsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListUniInterestRecordsService() *ListUniInterestRecordsService {
	return &ListUniInterestRecordsService{c: c, params: map[string]string{}}
}

// SetCurrency filters the result to a single currency.
func (s *ListUniInterestRecordsService) SetCurrency(currency string) *ListUniInterestRecordsService {
	s.params["currency"] = currency
	return s
}

// SetPage selects the result page (1-based).
func (s *ListUniInterestRecordsService) SetPage(page int) *ListUniInterestRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *ListUniInterestRecordsService) SetLimit(limit int) *ListUniInterestRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom bounds the result to records at or after this time.
func (s *ListUniInterestRecordsService) SetFrom(from time.Time) *ListUniInterestRecordsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo bounds the result to records at or before this time.
func (s *ListUniInterestRecordsService) SetTo(to time.Time) *ListUniInterestRecordsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListUniInterestRecordsService) Do(ctx context.Context) ([]UniInterestRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/interest_records", s.params).WithSign()
	resp, err := request.Do[[]UniInterestRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniInterestRecord is a single hourly interest settlement record. Status is
// 0 (failed) or 1 (success).
type UniInterestRecord struct {
	Status         int             `json:"status"`
	Currency       string          `json:"currency"`
	ActualRate     decimal.Decimal `json:"actual_rate"`
	Interest       decimal.Decimal `json:"interest"`
	InterestStatus string          `json:"interest_status"`
	CreateTime     time.Time       `json:"create_time,format:unix"`
}

// GetUniInterestStatusService -- GET /api/v4/earn/uni/interest_status/{currency} (private)
//
// Returns the interest-compounding status of a currency for the user.
type GetUniInterestStatusService struct {
	c        *EarnClient
	currency string
}

func (c *EarnClient) NewGetUniInterestStatusService(currency string) *GetUniInterestStatusService {
	return &GetUniInterestStatusService{c: c, currency: currency}
}

func (s *GetUniInterestStatusService) Do(ctx context.Context) (*UniInterestStatus, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/interest_status/"+s.currency).WithSign()
	return request.Do[UniInterestStatus](req)
}

// UniInterestStatus is a currency's interest-compounding status:
// interest_dividend (normal dividend) or interest_reinvest (reinvestment).
type UniInterestStatus struct {
	Currency       string `json:"currency"`
	InterestStatus string `json:"interest_status"`
}

// ListUniChartService -- GET /api/v4/earn/uni/chart
//
// Returns the annualized-rate trend chart for a currency over a time window
// (maximum span 30 days).
type ListUniChartService struct {
	c      *EarnClient
	params map[string]string
}

// NewListUniChartService builds the chart request for the given currency
// (asset) over [from, to].
func (c *EarnClient) NewListUniChartService(asset string, from, to time.Time) *ListUniChartService {
	return &ListUniChartService{c: c, params: map[string]string{
		"asset": asset,
		"from":  strconv.FormatInt(from.Unix(), 10),
		"to":    strconv.FormatInt(to.Unix(), 10),
	}}
}

func (s *ListUniChartService) Do(ctx context.Context) ([]UniChartPoint, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/chart", s.params)
	resp, err := request.Do[[]UniChartPoint](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniChartPoint is one sample of the annualized-rate trend chart.
type UniChartPoint struct {
	Time  time.Time       `json:"time,format:unix"`
	Value decimal.Decimal `json:"value"`
}

// ListUniRateService -- GET /api/v4/earn/uni/rate
//
// Returns the estimated annualized interest rate of every Uni-lending currency.
type ListUniRateService struct {
	c *EarnClient
}

func (c *EarnClient) NewListUniRateService() *ListUniRateService {
	return &ListUniRateService{c: c}
}

func (s *ListUniRateService) Do(ctx context.Context) ([]UniRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/uni/rate")
	resp, err := request.Do[[]UniRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UniRate is a currency's estimated annualized Uni-lending interest rate.
type UniRate struct {
	Currency string          `json:"currency"`
	EstRate  decimal.Decimal `json:"est_rate"`
}
