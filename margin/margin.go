package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListMarginAccountsService -- GET /api/v4/margin/accounts (private)
//
// Returns the authenticated user's isolated margin accounts, optionally filtered
// to a single trading pair.
type ListMarginAccountsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListMarginAccountsService() *ListMarginAccountsService {
	return &ListMarginAccountsService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair (e.g. BTC_USDT).
func (s *ListMarginAccountsService) SetCurrencyPair(currencyPair string) *ListMarginAccountsService {
	s.params["currency_pair"] = currencyPair
	return s
}

func (s *ListMarginAccountsService) Do(ctx context.Context) ([]MarginAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/accounts", s.params).WithSign()
	resp, err := request.Do[[]MarginAccount](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginAccount is a single isolated margin account (one trading pair). base is
// the base-currency sub-account, quote the quote-currency sub-account.
type MarginAccount struct {
	CurrencyPair string                `json:"currency_pair"`
	AccountType  string                `json:"account_type"`
	Leverage     decimal.Decimal       `json:"leverage"`
	Locked       bool                  `json:"locked"`
	Risk         decimal.Decimal       `json:"risk"`
	MMR          decimal.Decimal       `json:"mmr"`
	Base         MarginAccountCurrency `json:"base"`
	Quote        MarginAccountCurrency `json:"quote"`
}

// MarginAccountCurrency is one currency's balance inside an isolated margin
// account.
type MarginAccountCurrency struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
	Borrowed  decimal.Decimal `json:"borrowed"`
	Interest  decimal.Decimal `json:"interest"`
}

// ListMarginAccountBookService -- GET /api/v4/margin/account_book (private)
//
// Returns the isolated margin account balance change history. The query time
// range cannot exceed 30 days.
type ListMarginAccountBookService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListMarginAccountBookService() *ListMarginAccountBookService {
	return &ListMarginAccountBookService{c: c, params: map[string]string{}}
}

// SetCurrency filters by currency. When set, currency_pair must also be set.
func (s *ListMarginAccountBookService) SetCurrency(currency string) *ListMarginAccountBookService {
	s.params["currency"] = currency
	return s
}

// SetCurrencyPair filters by margin trading pair (used together with currency).
func (s *ListMarginAccountBookService) SetCurrencyPair(currencyPair string) *ListMarginAccountBookService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetType filters by account book change type.
func (s *ListMarginAccountBookService) SetType(changeType string) *ListMarginAccountBookService {
	s.params["type"] = changeType
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListMarginAccountBookService) SetFrom(from time.Time) *ListMarginAccountBookService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (defaults to now).
func (s *ListMarginAccountBookService) SetTo(to time.Time) *ListMarginAccountBookService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetPage sets the page number.
func (s *ListMarginAccountBookService) SetPage(page int) *ListMarginAccountBookService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned.
func (s *ListMarginAccountBookService) SetLimit(limit int) *ListMarginAccountBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListMarginAccountBookService) Do(ctx context.Context) ([]MarginAccountBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/account_book", s.params).WithSign()
	resp, err := request.Do[[]MarginAccountBook](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginAccountBook is a single isolated margin balance change record.
type MarginAccountBook struct {
	ID           string          `json:"id"`
	Time         time.Time       `json:"time,string,format:unix"`
	TimeMs       time.Time       `json:"time_ms,format:unixmilli"`
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	Change       decimal.Decimal `json:"change"`
	Balance      decimal.Decimal `json:"balance"`
	Type         string          `json:"type"`
}

// ListFundingAccountsService -- GET /api/v4/margin/funding_accounts (private)
//
// Returns the authenticated user's funding (lending) accounts, optionally
// filtered to a single currency.
type ListFundingAccountsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListFundingAccountsService() *ListFundingAccountsService {
	return &ListFundingAccountsService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListFundingAccountsService) SetCurrency(currency string) *ListFundingAccountsService {
	s.params["currency"] = currency
	return s
}

func (s *ListFundingAccountsService) Do(ctx context.Context) ([]FundingAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/funding_accounts", s.params).WithSign()
	resp, err := request.Do[[]FundingAccount](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FundingAccount is a single currency's lending (funding) balance.
type FundingAccount struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
	Lent      decimal.Decimal `json:"lent"`
	TotalLent decimal.Decimal `json:"total_lent"`
}

// GetAutoRepayStatusService -- GET /api/v4/margin/auto_repay (private)
//
// Returns the user's cross-margin auto-repayment setting.
type GetAutoRepayStatusService struct {
	c *MarginClient
}

func (c *MarginClient) NewGetAutoRepayStatusService() *GetAutoRepayStatusService {
	return &GetAutoRepayStatusService{c: c}
}

func (s *GetAutoRepayStatusService) Do(ctx context.Context) (*AutoRepayStatus, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/auto_repay").WithSign()
	return request.Do[AutoRepayStatus](req)
}

// AutoRepayStatus is the cross-margin auto-repayment toggle: "on" or "off".
type AutoRepayStatus struct {
	Status string `json:"status"`
}

// SetAutoRepayService -- POST /api/v4/margin/auto_repay (private)
//
// Enables or disables cross-margin auto repayment. status is "on" or "off".
type SetAutoRepayService struct {
	c      *MarginClient
	status string
}

func (c *MarginClient) NewSetAutoRepayService(status string) *SetAutoRepayService {
	return &SetAutoRepayService{c: c, status: status}
}

func (s *SetAutoRepayService) Do(ctx context.Context) (*AutoRepayStatus, error) {
	req := request.Post(ctx, s.c, "/api/v4/margin/auto_repay").WithSign().SetQuery("status", s.status)
	return request.Do[AutoRepayStatus](req)
}

// GetMarginTransferableService -- GET /api/v4/margin/transferable (private)
//
// Returns the maximum amount of a currency that can be transferred out of an
// isolated margin account.
type GetMarginTransferableService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetMarginTransferableService(currency string) *GetMarginTransferableService {
	return &GetMarginTransferableService{c: c, params: map[string]string{"currency": currency}}
}

// SetCurrencyPair scopes the query to a single trading pair.
func (s *GetMarginTransferableService) SetCurrencyPair(currencyPair string) *GetMarginTransferableService {
	s.params["currency_pair"] = currencyPair
	return s
}

func (s *GetMarginTransferableService) Do(ctx context.Context) (*MarginTransferable, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/transferable", s.params).WithSign()
	return request.Do[MarginTransferable](req)
}

// MarginTransferable is the maximum transferable amount of a currency.
type MarginTransferable struct {
	Currency     string          `json:"currency"`
	CurrencyPair string          `json:"currency_pair"`
	Amount       decimal.Decimal `json:"amount"`
}

// GetUserMarginTierService -- GET /api/v4/margin/user/loan_margin_tiers (private)
//
// Returns the user's own leverage lending tiers for a market.
type GetUserMarginTierService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetUserMarginTierService(currencyPair string) *GetUserMarginTierService {
	return &GetUserMarginTierService{c: c, params: map[string]string{"currency_pair": currencyPair}}
}

func (s *GetUserMarginTierService) Do(ctx context.Context) ([]MarginTier, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/user/loan_margin_tiers", s.params).WithSign()
	resp, err := request.Do[[]MarginTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetMarketMarginTierService -- GET /api/v4/margin/loan_margin_tiers
//
// Returns the current market leverage lending tiers for a trading pair.
type GetMarketMarginTierService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetMarketMarginTierService(currencyPair string) *GetMarketMarginTierService {
	return &GetMarketMarginTierService{c: c, params: map[string]string{"currency_pair": currencyPair}}
}

func (s *GetMarketMarginTierService) Do(ctx context.Context) ([]MarginTier, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/loan_margin_tiers", s.params)
	resp, err := request.Do[[]MarginTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginTier is one gradient of the market's leverage lending schedule.
type MarginTier struct {
	UpperLimit decimal.Decimal `json:"upper_limit"`
	MMR        decimal.Decimal `json:"mmr"`
	Leverage   decimal.Decimal `json:"leverage"`
}

// SetUserMarketLeverageService -- POST /api/v4/margin/leverage/user_market_setting (private)
//
// Sets the user's leverage multiplier for a margin market. The endpoint returns
// no content.
type SetUserMarketLeverageService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewSetUserMarketLeverageService(currencyPair, leverage string) *SetUserMarketLeverageService {
	return &SetUserMarketLeverageService{c: c, body: map[string]any{
		"currency_pair": currencyPair,
		"leverage":      leverage,
	}}
}

func (s *SetUserMarketLeverageService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/margin/leverage/user_market_setting", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListMarginUserAccountService -- GET /api/v4/margin/user/account (private)
//
// Returns the user's isolated margin accounts, supporting both risk-ratio and
// margin-ratio isolated accounts.
type ListMarginUserAccountService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListMarginUserAccountService() *ListMarginUserAccountService {
	return &ListMarginUserAccountService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *ListMarginUserAccountService) SetCurrencyPair(currencyPair string) *ListMarginUserAccountService {
	s.params["currency_pair"] = currencyPair
	return s
}

func (s *ListMarginUserAccountService) Do(ctx context.Context) ([]MarginUserAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/user/account", s.params).WithSign()
	resp, err := request.Do[[]MarginUserAccount](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginUserAccount is a single isolated margin account as returned by
// /margin/user/account. base is the base-currency sub-account, quote the
// quote-currency sub-account.
type MarginUserAccount struct {
	CurrencyPair string                `json:"currency_pair"`
	AccountType  string                `json:"account_type"`
	Leverage     decimal.Decimal       `json:"leverage"`
	Locked       bool                  `json:"locked"`
	Risk         decimal.Decimal       `json:"risk"`
	MMR          decimal.Decimal       `json:"mmr"`
	Base         MarginAccountCurrency `json:"base"`
	Quote        MarginAccountCurrency `json:"quote"`
}

// ListCrossMarginLoansService -- GET /api/v4/margin/cross/loans (private)
//
// Returns the cross-margin borrow history, newest first by default. Deprecated
// upstream; status filters by the (now-fixed) loan status.
type ListCrossMarginLoansService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListCrossMarginLoansService(status int) *ListCrossMarginLoansService {
	return &ListCrossMarginLoansService{c: c, params: map[string]string{
		"status": strconv.Itoa(status),
	}}
}

// SetCurrency filters by currency (all currencies when unset).
func (s *ListCrossMarginLoansService) SetCurrency(currency string) *ListCrossMarginLoansService {
	s.params["currency"] = currency
	return s
}

// SetLimit caps the number of records returned.
func (s *ListCrossMarginLoansService) SetLimit(limit int) *ListCrossMarginLoansService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListCrossMarginLoansService) SetOffset(offset int) *ListCrossMarginLoansService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetReverse toggles sort order; the default (true) is descending by creation
// time, set false for ascending.
func (s *ListCrossMarginLoansService) SetReverse(reverse bool) *ListCrossMarginLoansService {
	s.params["reverse"] = strconv.FormatBool(reverse)
	return s
}

func (s *ListCrossMarginLoansService) Do(ctx context.Context) ([]CrossMarginLoan, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/cross/loans", s.params).WithSign()
	resp, err := request.Do[[]CrossMarginLoan](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossMarginLoan is a single cross-margin borrow record.
type CrossMarginLoan struct {
	ID             string          `json:"id"`
	CreateTime     time.Time       `json:"create_time,format:unixmilli"`
	UpdateTime     time.Time       `json:"update_time,format:unixmilli"`
	Currency       string          `json:"currency"`
	Amount         decimal.Decimal `json:"amount"`
	Text           string          `json:"text"`
	Status         int             `json:"status"`
	Repaid         decimal.Decimal `json:"repaid"`
	RepaidInterest decimal.Decimal `json:"repaid_interest"`
	UnpaidInterest decimal.Decimal `json:"unpaid_interest"`
}

// ListCrossMarginRepaymentsService -- GET /api/v4/margin/cross/repayments (private)
//
// Returns cross-margin repayment records, newest first by default. Deprecated
// upstream.
type ListCrossMarginRepaymentsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewListCrossMarginRepaymentsService() *ListCrossMarginRepaymentsService {
	return &ListCrossMarginRepaymentsService{c: c, params: map[string]string{}}
}

// SetCurrency filters by currency.
func (s *ListCrossMarginRepaymentsService) SetCurrency(currency string) *ListCrossMarginRepaymentsService {
	s.params["currency"] = currency
	return s
}

// SetLoanID filters by the originating loan record ID.
func (s *ListCrossMarginRepaymentsService) SetLoanID(loanID string) *ListCrossMarginRepaymentsService {
	s.params["loan_id"] = loanID
	return s
}

// SetLimit caps the number of records returned.
func (s *ListCrossMarginRepaymentsService) SetLimit(limit int) *ListCrossMarginRepaymentsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListCrossMarginRepaymentsService) SetOffset(offset int) *ListCrossMarginRepaymentsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetReverse toggles sort order; the default (true) is descending by creation
// time, set false for ascending.
func (s *ListCrossMarginRepaymentsService) SetReverse(reverse bool) *ListCrossMarginRepaymentsService {
	s.params["reverse"] = strconv.FormatBool(reverse)
	return s
}

func (s *ListCrossMarginRepaymentsService) Do(ctx context.Context) ([]CrossMarginRepayment, error) {
	req := request.Get(ctx, s.c, "/api/v4/margin/cross/repayments", s.params).WithSign()
	resp, err := request.Do[[]CrossMarginRepayment](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossMarginRepayment is a single cross-margin repayment record.
type CrossMarginRepayment struct {
	ID            string          `json:"id"`
	CreateTime    time.Time       `json:"create_time,format:unix"`
	LoanID        string          `json:"loan_id"`
	Currency      string          `json:"currency"`
	Principal     decimal.Decimal `json:"principal"`
	Interest      decimal.Decimal `json:"interest"`
	RepaymentType string          `json:"repayment_type"`
}
