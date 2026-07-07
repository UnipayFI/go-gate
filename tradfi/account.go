package tradfi

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetMT5AccountService -- GET /api/v4/tradfi/users/mt5-account (private)
//
// Queries the caller's MT5 account information (leverage, status, liquidation
// margin ratio).
type GetMT5AccountService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewGetMT5AccountService() *GetMT5AccountService {
	return &GetMT5AccountService{c: c}
}

func (s *GetMT5AccountService) Do(ctx context.Context) (*TradfiMT5AccountResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/users/mt5-account").WithSign()
	return request.Do[TradfiMT5AccountResponse](req)
}

// TradfiMT5AccountResponse is the envelope of the MT5-account query.
type TradfiMT5AccountResponse struct {
	Code      int              `json:"code"`
	Message   string           `json:"message"`
	Timestamp time.Time        `json:"timestamp,format:unixmilli"`
	Data      TradfiMT5Account `json:"data"`
}

// TradfiMT5Account is the MT5 account information block.
type TradfiMT5Account struct {
	MT5UID       int64           `json:"mt5_uid"`
	Leverage     int             `json:"leverage"`
	StopOutLevel decimal.Decimal `json:"stop_out_level"`
	Status       int             `json:"status"`
}

// CreateUserService -- POST /api/v4/tradfi/users (private)
//
// Creates (opens) a TradFi user for the caller and returns its initial status.
type CreateUserService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewCreateUserService() *CreateUserService {
	return &CreateUserService{c: c}
}

func (s *CreateUserService) Do(ctx context.Context) (*TradfiCreateUserResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/tradfi/users").WithSign()
	return request.Do[TradfiCreateUserResponse](req)
}

// TradfiCreateUserResponse is the envelope of the create-user request.
type TradfiCreateUserResponse struct {
	Timestamp time.Time            `json:"timestamp,format:unixmilli"`
	Data      TradfiCreateUserData `json:"data"`
}

// TradfiCreateUserData is the newly created TradFi user's status.
type TradfiCreateUserData struct {
	Status   int    `json:"status"`
	Leverage int    `json:"leverage"`
	MT5UID   string `json:"mt5_uid"`
}

// GetAssetsService -- GET /api/v4/tradfi/users/assets (private)
//
// Queries the caller's TradFi account assets (equity, balance, margin, PnL).
type GetAssetsService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewGetAssetsService() *GetAssetsService {
	return &GetAssetsService{c: c}
}

func (s *GetAssetsService) Do(ctx context.Context) (*TradfiAssetsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/users/assets").WithSign()
	return request.Do[TradfiAssetsResponse](req)
}

// TradfiAssetsResponse is the envelope of the account-assets query.
type TradfiAssetsResponse struct {
	Timestamp time.Time    `json:"timestamp,format:unixmilli"`
	Data      TradfiAssets `json:"data"`
}

// TradfiAssets is the TradFi account asset snapshot.
type TradfiAssets struct {
	Equity        decimal.Decimal `json:"equity"`
	MarginLevel   decimal.Decimal `json:"margin_level"`
	Balance       decimal.Decimal `json:"balance"`
	Margin        decimal.Decimal `json:"margin"`
	MarginFree    decimal.Decimal `json:"margin_free"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	MT5UID        string          `json:"mt5_uid"`
}

// ListTransactionsService -- GET /api/v4/tradfi/transactions (private)
//
// Lists the caller's fund transfer-in/out (and dividend) records, paginated.
type ListTransactionsService struct {
	c      *TradfiClient
	params map[string]string
}

func (c *TradfiClient) NewListTransactionsService() *ListTransactionsService {
	return &ListTransactionsService{c: c, params: map[string]string{}}
}

// SetBeginTime bounds the result to records at or after this time.
func (s *ListTransactionsService) SetBeginTime(beginTime time.Time) *ListTransactionsService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to records at or before this time.
func (s *ListTransactionsService) SetEndTime(endTime time.Time) *ListTransactionsService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetType filters by transaction type ("deposit", "withdraw" or "dividend").
func (s *ListTransactionsService) SetType(transactionType string) *ListTransactionsService {
	s.params["type"] = transactionType
	return s
}

// SetPage selects the result page (1-based).
func (s *ListTransactionsService) SetPage(page int) *ListTransactionsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (default 10, max 50).
func (s *ListTransactionsService) SetPageSize(pageSize int) *ListTransactionsService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListTransactionsService) Do(ctx context.Context) (*TradfiTransactionsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/transactions", s.params).WithSign()
	return request.Do[TradfiTransactionsResponse](req)
}

// TradfiTransactionsResponse is the envelope of the transaction-records query.
// The server timestamp is nested inside data on this endpoint.
type TradfiTransactionsResponse struct {
	Data struct {
		Total     int                 `json:"total"`
		TotalPage int                 `json:"total_page"`
		List      []TradfiTransaction `json:"list"`
		Timestamp time.Time           `json:"timestamp,format:unixmilli"`
	} `json:"data"`
}

// TradfiTransaction is a single fund transfer-in/out (or dividend) record.
type TradfiTransaction struct {
	Asset    string          `json:"asset"`
	Type     string          `json:"type"`
	TypeDesc string          `json:"type_desc"`
	Change   decimal.Decimal `json:"change"`
	Balance  decimal.Decimal `json:"balance"`
	Time     time.Time       `json:"time,format:unix"`
}

// CreateTransactionService -- POST /api/v4/tradfi/transactions (private)
//
// Transfers funds into or out of the TradFi account. transactionType is
// "deposit" (transfer in) or "withdraw" (transfer out); change supports up to
// two decimal places.
type CreateTransactionService struct {
	c    *TradfiClient
	body map[string]any
}

func (c *TradfiClient) NewCreateTransactionService(asset string, change decimal.Decimal, transactionType string) *CreateTransactionService {
	return &CreateTransactionService{c: c, body: map[string]any{
		"asset":  asset,
		"change": change.String(),
		"type":   transactionType,
	}}
}

func (s *CreateTransactionService) Do(ctx context.Context) (*TradfiCreateTransactionResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/tradfi/transactions", s.body).WithSign()
	return request.Do[TradfiCreateTransactionResponse](req)
}

// TradfiCreateTransactionResponse is the envelope of the fund transfer request.
// data is an empty object on success.
type TradfiCreateTransactionResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct{}  `json:"data"`
}
