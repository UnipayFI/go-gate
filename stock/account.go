package stock

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetUserAssetsService -- GET /api/v4/stock/users/assets (private)
//
// Queries the caller's stock account assets (equity, balance, position market
// value and PnL). pnlCalcType selects the cost basis (1=average cost,
// 2=diluted cost); pnlCalcPrice selects the price basis (1=intraday, 2=latest
// extended-hours).
type GetUserAssetsService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewGetUserAssetsService() *GetUserAssetsService {
	return &GetUserAssetsService{c: c, params: map[string]string{}}
}

// SetPnLCalcType selects the PnL cost basis (1=average cost, 2=diluted cost).
func (s *GetUserAssetsService) SetPnLCalcType(pnlCalcType int) *GetUserAssetsService {
	s.params["pnl_calc_type"] = strconv.Itoa(pnlCalcType)
	return s
}

// SetPnLCalcPrice selects the PnL price basis (1=intraday, 2=latest
// extended-hours).
func (s *GetUserAssetsService) SetPnLCalcPrice(pnlCalcPrice int) *GetUserAssetsService {
	s.params["pnl_calc_price"] = strconv.Itoa(pnlCalcPrice)
	return s
}

func (s *GetUserAssetsService) Do(ctx context.Context) (*StockUserAssetsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/users/assets", s.params).WithSign()
	return request.Do[StockUserAssetsResponse](req)
}

// StockUserAssetsResponse is the envelope of the user-assets query.
type StockUserAssetsResponse struct {
	Timestamp time.Time       `json:"timestamp,format:unixmilli"`
	Data      StockUserAssets `json:"data"`
}

// StockUserAssets is the stock account asset snapshot. user_exists reports
// whether the caller has activated the service.
type StockUserAssets struct {
	Equity              decimal.Decimal `json:"equity"`
	Balance             decimal.Decimal `json:"balance"`
	Available           decimal.Decimal `json:"available"`
	PositionMarketValue decimal.Decimal `json:"position_market_value"`
	PositionPnL         decimal.Decimal `json:"position_pnl"`
	TodayPnL            decimal.Decimal `json:"today_pnl"`
	UserExists          bool            `json:"user_exists"`
}

// ListTransactionsService -- GET /api/v4/stock/transactions (private)
//
// Lists the caller's transaction records (deposits, withdrawals, fees,
// dividends, trades, awards and stock transfers), paginated. When refID is set
// the server queries by that idempotent id and ignores the other filters.
type ListTransactionsService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListTransactionsService() *ListTransactionsService {
	return &ListTransactionsService{c: c, params: map[string]string{}}
}

// SetBeginTime bounds the result to records at or after this time (query range
// must not exceed 3 months).
func (s *ListTransactionsService) SetBeginTime(beginTime time.Time) *ListTransactionsService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to records at or before this time (query range
// must not exceed 3 months).
func (s *ListTransactionsService) SetEndTime(endTime time.Time) *ListTransactionsService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetRefID queries by business idempotent id, ignoring the other filters.
func (s *ListTransactionsService) SetRefID(refID string) *ListTransactionsService {
	s.params["ref_id"] = refID
	return s
}

// SetType filters by transaction type (e.g. "deposit", "withdraw", "fee",
// "dividend", "sell", "buy", "award", "stock_transfer_in",
// "stock_transfer_out").
func (s *ListTransactionsService) SetType(transactionType string) *ListTransactionsService {
	s.params["type"] = transactionType
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListTransactionsService) SetPage(page int) *ListTransactionsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (defaults to 10, max 500).
func (s *ListTransactionsService) SetPageSize(pageSize int) *ListTransactionsService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListTransactionsService) Do(ctx context.Context) (*StockTransactionsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/transactions", s.params).WithSign()
	return request.Do[StockTransactionsResponse](req)
}

// StockTransactionsResponse is the envelope of the transaction-records query.
type StockTransactionsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64              `json:"total"`
		TotalPage int64              `json:"total_page"`
		List      []StockTransaction `json:"list"`
	} `json:"data"`
}

// StockTransaction is a single transaction record. time is an integer-second
// Unix timestamp; detail carries free-form business details.
type StockTransaction struct {
	Asset         string          `json:"asset"`
	Symbol        string          `json:"symbol"`
	SymbolDisplay string          `json:"symbol_display"`
	Type          string          `json:"type"`
	TypeDesc      string          `json:"type_desc"`
	Change        decimal.Decimal `json:"change"`
	Balance       decimal.Decimal `json:"balance"`
	RefID         string          `json:"ref_id"`
	Time          time.Time       `json:"time,format:unix"`
	UnitText      string          `json:"unit_text"`
	Detail        map[string]any  `json:"detail"`
}

// CreateTransactionService -- POST /api/v4/stock/transactions (private)
//
// Transfers funds into or out of the stock account. asset is "USDT" only;
// transactionType is "deposit" (transfer in) or "withdraw" (transfer out);
// refID is a required business idempotent id.
type CreateTransactionService struct {
	c    *StockClient
	body map[string]any
}

func (c *StockClient) NewCreateTransactionService(asset string, change decimal.Decimal, transactionType, refID string) *CreateTransactionService {
	return &CreateTransactionService{c: c, body: map[string]any{
		"asset":  asset,
		"change": change.String(),
		"type":   transactionType,
		"ref_id": refID,
	}}
}

func (s *CreateTransactionService) Do(ctx context.Context) (*StockCreateTransactionResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/stock/transactions", s.body).WithSign()
	return request.Do[StockCreateTransactionResponse](req)
}

// StockCreateTransactionResponse is the envelope of the fund-transfer request.
// data is an empty object on success.
type StockCreateTransactionResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct{}  `json:"data"`
}
