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
// whether the caller has activated the service; margin_enabled whether margin
// financing is on. loan_amount is the absolute value of the available balance
// when it is negative, otherwise 0; leave_liquidity / leave_liquidity_rate are
// the excess liquidity and its ratio; aee is the available excess equity;
// long_borrowing_rate is the annual margin financing rate, sent as "" (and
// decoded as zero) while no valid rate is available.
type StockUserAssets struct {
	Equity                    decimal.Decimal `json:"equity"`
	Balance                   decimal.Decimal `json:"balance"`
	Available                 decimal.Decimal `json:"available"`
	LoanAmount                decimal.Decimal `json:"loan_amount"`
	BuyPower                  decimal.Decimal `json:"buy_power"`
	Leverage                  decimal.Decimal `json:"leverage"`
	LeaveLiquidity            decimal.Decimal `json:"leave_liquidity"`
	LeaveLiquidityRate        decimal.Decimal `json:"leave_liquidity_rate"`
	AEE                       decimal.Decimal `json:"aee"`
	OptionsMargin             decimal.Decimal `json:"options_margin"`
	PositionMarketValue       decimal.Decimal `json:"position_market_value"`
	PositionPnL               decimal.Decimal `json:"position_pnl"`
	TodayPnL                  decimal.Decimal `json:"today_pnl"`
	OptionPositionMarketValue decimal.Decimal `json:"option_position_market_value"`
	OptionPositionPnL         decimal.Decimal `json:"option_position_pnl"`
	OptionTodayPnL            decimal.Decimal `json:"option_today_pnl"`
	UserExists                bool            `json:"user_exists"`
	MarginEnabled             bool            `json:"margin_enabled"`
	LongBorrowingRate         decimal.Decimal `json:"long_borrowing_rate"`
}

// ListBalanceRecordsService -- GET /api/v4/stock/users/balance-records (private)
//
// Lists the caller's raw balance records (cash and stock quantity changes),
// always ordered by record id: id desc by default, or id asc returning only
// records after start_id when SetStartID is used without SetSortDirection.
type ListBalanceRecordsService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListBalanceRecordsService() *ListBalanceRecordsService {
	return &ListBalanceRecordsService{c: c, params: map[string]string{}}
}

// SetType filters by record type, comma-separated (e.g. "deposit", "withdraw",
// "fee", "dividend", "sell", "buy", "award", "stock_transfer_in",
// "stock_transfer_out").
func (s *ListBalanceRecordsService) SetType(recordType string) *ListBalanceRecordsService {
	s.params["type"] = recordType
	return s
}

// SetBusiness filters by business identifier, comma-separated.
func (s *ListBalanceRecordsService) SetBusiness(business string) *ListBalanceRecordsService {
	s.params["business"] = business
	return s
}

// SetSymbol narrows the result to a single symbol.
func (s *ListBalanceRecordsService) SetSymbol(symbol string) *ListBalanceRecordsService {
	s.params["symbol"] = symbol
	return s
}

// SetBeginTime bounds the result to records created at or after this time
// (sent in seconds; when both bounds are set, end_time must be after
// begin_time).
func (s *ListBalanceRecordsService) SetBeginTime(beginTime time.Time) *ListBalanceRecordsService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to records created at or before this time
// (sent in seconds).
func (s *ListBalanceRecordsService) SetEndTime(endTime time.Time) *ListBalanceRecordsService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetStartID returns only records whose id is greater than startID, for
// incremental retrieval.
func (s *ListBalanceRecordsService) SetStartID(startID int64) *ListBalanceRecordsService {
	s.params["start_id"] = strconv.FormatInt(startID, 10)
	return s
}

// SetSortDirection orders the result by record id ("asc" or "desc").
func (s *ListBalanceRecordsService) SetSortDirection(sortDirection string) *ListBalanceRecordsService {
	s.params["sort_direction"] = sortDirection
	return s
}

// SetIncludeDetail returns the raw detail JSON string on each record when true.
func (s *ListBalanceRecordsService) SetIncludeDetail(includeDetail bool) *ListBalanceRecordsService {
	s.params["include_detail"] = strconv.FormatBool(includeDetail)
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListBalanceRecordsService) SetPage(page int) *ListBalanceRecordsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (defaults to 10; the server
// caps values above 100 at 100).
func (s *ListBalanceRecordsService) SetPageSize(pageSize int) *ListBalanceRecordsService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListBalanceRecordsService) Do(ctx context.Context) (*StockBalanceRecordsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/users/balance-records", s.params).WithSign()
	return request.Do[StockBalanceRecordsResponse](req)
}

// StockBalanceRecordsResponse is the envelope of the balance-records query.
type StockBalanceRecordsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64                `json:"total"`
		TotalPage int64                `json:"total_page"`
		List      []StockBalanceRecord `json:"list"`
	} `json:"data"`
}

// StockBalanceRecord is a single raw balance record. change is the cash change
// (positive credit, negative debit) and balance the cash balance snapshot after
// it, not a stock position balance; stock_change is the stock quantity change
// parsed from detail.amount ("" and so zero for other records); detail is the
// raw business detail as a JSON string, present only with include_detail=true;
// time_create is the business event time and create_time the record creation
// time, both integer-second Unix timestamps.
type StockBalanceRecord struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Change      decimal.Decimal `json:"change"`
	Balance     decimal.Decimal `json:"balance"`
	StockChange decimal.Decimal `json:"stock_change"`
	Type        string          `json:"type"`
	Business    string          `json:"business"`
	BusinessID  int64           `json:"business_id"`
	Symbol      string          `json:"symbol"`
	Detail      string          `json:"detail"`
	Version     int             `json:"version"`
	TimeCreate  time.Time       `json:"time_create,format:unix"`
	CreateTime  time.Time       `json:"create_time,format:unix"`
}

// SetMarginFeatureService -- PUT /api/v4/stock/users/margin-feature (private)
//
// Enables (true) or disables (false) margin financing for the caller. Gate
// returns an empty data object on success.
type SetMarginFeatureService struct {
	c    *StockClient
	body map[string]any
}

func (c *StockClient) NewSetMarginFeatureService(marginFeature bool) *SetMarginFeatureService {
	return &SetMarginFeatureService{c: c, body: map[string]any{
		"margin_feature": marginFeature,
	}}
}

func (s *SetMarginFeatureService) Do(ctx context.Context) (*StockSetMarginFeatureResponse, error) {
	req := request.Put(ctx, s.c, "/api/v4/stock/users/margin-feature", s.body).WithSign()
	return request.Do[StockSetMarginFeatureResponse](req)
}

// StockSetMarginFeatureResponse is the envelope of the margin-financing toggle.
// data is an empty object on success.
type StockSetMarginFeatureResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct{}  `json:"data"`
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
