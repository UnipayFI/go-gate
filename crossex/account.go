package crossex

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetAccountService -- GET /api/v4/crossex/accounts (private)
//
// Returns the cross-exchange account assets, margins and per-venue balances.
type GetAccountService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c, params: map[string]string{}}
}

// SetExchangeType selects the venue in isolated-per-venue mode (omit in
// cross-exchange mode).
func (s *GetAccountService) SetExchangeType(exchangeType string) *GetAccountService {
	s.params["exchange_type"] = exchangeType
	return s
}

func (s *GetAccountService) Do(ctx context.Context) (*CrossexAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/accounts", s.params).WithSign()
	return request.Do[CrossexAccount](req)
}

// CrossexAccount is the cross-exchange account summary and its per-venue asset
// breakdown. create_time and update_time are millisecond Unix timestamps.
type CrossexAccount struct {
	UserID                string                `json:"user_id"`
	AvailableMargin       decimal.Decimal       `json:"available_margin"`
	MarginBalance         decimal.Decimal       `json:"margin_balance"`
	InitialMargin         decimal.Decimal       `json:"initial_margin"`
	MaintenanceMargin     decimal.Decimal       `json:"maintenance_margin"`
	InitialMarginRate     decimal.Decimal       `json:"initial_margin_rate"`
	MaintenanceMarginRate decimal.Decimal       `json:"maintenance_margin_rate"`
	PositionMode          string                `json:"position_mode"`
	AccountLimit          string                `json:"account_limit"`
	CreateTime            time.Time             `json:"create_time,string,format:unixmilli"`
	UpdateTime            time.Time             `json:"update_time,string,format:unixmilli"`
	AccountMode           string                `json:"account_mode"`
	ExchangeType          string                `json:"exchange_type"`
	Assets                []CrossexAccountAsset `json:"assets"`
}

// CrossexAccountAsset is a per-venue, per-currency balance within the account.
type CrossexAccountAsset struct {
	UserID                     string          `json:"user_id"`
	Coin                       string          `json:"coin"`
	ExchangeType               string          `json:"exchange_type"`
	Balance                    decimal.Decimal `json:"balance"`
	UPnL                       decimal.Decimal `json:"upnl"`
	Equity                     decimal.Decimal `json:"equity"`
	FuturesInitialMargin       decimal.Decimal `json:"futures_initial_margin"`
	FuturesMaintenanceMargin   decimal.Decimal `json:"futures_maintenance_margin"`
	BorrowingInitialMargin     decimal.Decimal `json:"borrowing_initial_margin"`
	BorrowingMaintenanceMargin decimal.Decimal `json:"borrowing_maintenance_margin"`
	AvailableBalance           decimal.Decimal `json:"available_balance"`
	Liability                  decimal.Decimal `json:"liability"`
}

// UpdateAccountService -- PUT /api/v4/crossex/accounts (private)
//
// Modifies the account's futures position mode and/or account mode.
type UpdateAccountService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewUpdateAccountService() *UpdateAccountService {
	return &UpdateAccountService{c: c, body: map[string]any{}}
}

// SetPositionMode sets the futures position mode (SINGLE or DUAL).
func (s *UpdateAccountService) SetPositionMode(positionMode string) *UpdateAccountService {
	s.body["position_mode"] = positionMode
	return s
}

// SetAccountMode sets the account mode (CROSS_EXCHANGE or ISOLATED_EXCHANGE).
func (s *UpdateAccountService) SetAccountMode(accountMode string) *UpdateAccountService {
	s.body["account_mode"] = accountMode
	return s
}

// SetExchangeType sets the venue the change targets (required in isolated mode).
func (s *UpdateAccountService) SetExchangeType(exchangeType string) *UpdateAccountService {
	s.body["exchange_type"] = exchangeType
	return s
}

func (s *UpdateAccountService) Do(ctx context.Context) (*CrossexAccountUpdate, error) {
	req := request.Put(ctx, s.c, "/api/v4/crossex/accounts", s.body).WithSign()
	return request.Do[CrossexAccountUpdate](req)
}

// CrossexAccountUpdate echoes the account-mode change that was requested.
type CrossexAccountUpdate struct {
	PositionMode string `json:"position_mode"`
	AccountMode  string `json:"account_mode"`
	ExchangeType string `json:"exchange_type"`
}

// ListAccountBookService -- GET /api/v4/crossex/account_book (private)
//
// Returns the account's asset-change (bill) history.
type ListAccountBookService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListAccountBookService() *ListAccountBookService {
	return &ListAccountBookService{c: c, params: map[string]string{}}
}

// SetPage selects the result page.
func (s *ListAccountBookService) SetPage(page int) *ListAccountBookService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list (max 1000).
func (s *ListAccountBookService) SetLimit(limit int) *ListAccountBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetCoin narrows the result to a single currency.
func (s *ListAccountBookService) SetCoin(coin string) *ListAccountBookService {
	s.params["coin"] = coin
	return s
}

// SetStatementType filters by bill entry type (e.g. TRANSACTION, TRADING_FEE).
func (s *ListAccountBookService) SetStatementType(statementType string) *ListAccountBookService {
	s.params["statement_type"] = statementType
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListAccountBookService) SetFrom(from time.Time) *ListAccountBookService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision).
func (s *ListAccountBookService) SetTo(to time.Time) *ListAccountBookService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

func (s *ListAccountBookService) Do(ctx context.Context) ([]CrossexAccountBookRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/account_book", s.params).WithSign()
	resp, err := request.Do[[]CrossexAccountBookRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexAccountBookRecord is one asset-change (bill) entry. create_time is a
// millisecond Unix timestamp; change is positive for inflows, negative for
// outflows.
type CrossexAccountBookRecord struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	BusinessID    string          `json:"business_id"`
	StatementType string          `json:"statement_type"`
	ExchangeType  string          `json:"exchange_type"`
	Coin          string          `json:"coin"`
	Change        decimal.Decimal `json:"change"`
	Balance       decimal.Decimal `json:"balance"`
	CreateTime    time.Time       `json:"create_time,string,format:unixmilli"`
}
