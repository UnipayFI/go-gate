package stockoptions

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetFeeRateService -- GET /api/v4/stock_options/fee-rate (public)
//
// Returns the fee configuration for stock options trading: rates and
// applicable minimum/maximum fees for SEC, TAF, OCC, ORF and CAT charges.
type GetFeeRateService struct {
	c *StockOptionsClient
}

func (c *StockOptionsClient) NewGetFeeRateService() *GetFeeRateService {
	return &GetFeeRateService{c: c}
}

func (s *GetFeeRateService) Do(ctx context.Context) (*StockOptionFeeRateResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/fee-rate")
	return request.Do[StockOptionFeeRateResponse](req)
}

// StockOptionFeeRateResponse is the envelope of the fee-rate query.
type StockOptionFeeRateResponse struct {
	Label     string             `json:"label"`
	Timestamp time.Time          `json:"timestamp,format:unixmilli"`
	Data      StockOptionFeeRate `json:"data"`
}

// StockOptionFeeRate holds the stock-options fee rates and min/max bounds.
// MaxSEC is sent as "" (decoded as zero) when no maximum limit applies.
type StockOptionFeeRate struct {
	SEC    decimal.Decimal `json:"sec"`
	MinSEC decimal.Decimal `json:"min_sec"`
	MaxSEC decimal.Decimal `json:"max_sec"`
	TAF    decimal.Decimal `json:"taf"`
	MinTAF decimal.Decimal `json:"min_taf"`
	MaxTAF decimal.Decimal `json:"max_taf"`
	OCC    decimal.Decimal `json:"occ"`
	ORF    decimal.Decimal `json:"orf"`
	CAT    decimal.Decimal `json:"cat"`
}

// ListUnderlyingsService -- GET /api/v4/stock_options/underlyings (public)
//
// Lists available stock-option underlyings with their descriptions, icons,
// price precision and order quantity limits.
type ListUnderlyingsService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewListUnderlyingsService() *ListUnderlyingsService {
	return &ListUnderlyingsService{c: c, params: map[string]string{}}
}

// SetUnderlyings narrows the result to a comma-separated list of underlying
// symbols (e.g. "AAPL,TSLA"). Omit to query all available underlyings.
func (s *ListUnderlyingsService) SetUnderlyings(underlyings string) *ListUnderlyingsService {
	s.params["underlyings"] = underlyings
	return s
}

func (s *ListUnderlyingsService) Do(ctx context.Context) (*StockOptionUnderlyingsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/underlyings", s.params)
	return request.Do[StockOptionUnderlyingsResponse](req)
}

// StockOptionUnderlyingsResponse is the envelope of the underlyings query.
type StockOptionUnderlyingsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []StockOptionUnderlying `json:"list"`
	} `json:"data"`
}

// StockOptionUnderlying is a single available option underlying.
type StockOptionUnderlying struct {
	Underlying      string            `json:"underlying"`
	Description     string            `json:"description"`
	Descs           []StockOptionDesc `json:"descs"`
	IconLink        string            `json:"icon_link"`
	PricePrecision  int               `json:"price_precision"`
	MinOrderVolume  decimal.Decimal   `json:"min_order_volume"`
	MaxOrderVolume  decimal.Decimal   `json:"max_order_volume"`
	StepOrderVolume decimal.Decimal   `json:"step_order_volume"`
}

// StockOptionDesc is one localized underlying description.
type StockOptionDesc struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// ListContractsService -- GET /api/v4/stock_options/contracts (public)
//
// Lists option contracts for an underlying, paginated in ascending contract
// symbol order. Default and maximum page_size are both 1000.
type ListContractsService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewListContractsService(underlying string) *ListContractsService {
	return &ListContractsService{c: c, params: map[string]string{
		"underlying": underlying,
	}}
}

// SetExchange narrows the result to a trading market (currently only "us").
func (s *ListContractsService) SetExchange(exchange string) *ListContractsService {
	s.params["exchange"] = exchange
	return s
}

// SetExpirationDate filters by expiration date in YYYY-MM-DD format.
func (s *ListContractsService) SetExpirationDate(expirationDate string) *ListContractsService {
	s.params["expiration_date"] = expirationDate
	return s
}

// SetOptionType filters by option type ("call" or "put"). Omit for both.
func (s *ListContractsService) SetOptionType(optionType string) *ListContractsService {
	s.params["option_type"] = optionType
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListContractsService) SetPage(page int) *ListContractsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the contracts per page (defaults to and maxes at 1000).
func (s *ListContractsService) SetPageSize(pageSize int) *ListContractsService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListContractsService) Do(ctx context.Context) (*StockOptionContractsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/contracts", s.params)
	return request.Do[StockOptionContractsResponse](req)
}

// StockOptionContractsResponse is the envelope of the contracts query.
type StockOptionContractsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64                 `json:"total"`
		TotalPage int                   `json:"total_page"`
		List      []StockOptionContract `json:"list"`
	} `json:"data"`
}

// StockOptionContract is a single stock-option contract specification.
// option_type is "call" or "put"; trade_status is "open" or "closed";
// trade_mode is 0 (trading disabled) or 4 (buying and selling allowed; regular
// accounts still require sufficient available long positions to sell).
type StockOptionContract struct {
	Symbol          string          `json:"symbol"`
	SymbolDisplay   string          `json:"symbol_display"`
	Underlying      string          `json:"underlying"`
	Exchange        string          `json:"exchange"`
	ExpirationDate  string          `json:"expiration_date"`
	SettlementDate  string          `json:"settlement_date"`
	StrikePrice     decimal.Decimal `json:"strike_price"`
	OptionType      string          `json:"option_type"`
	Multiplier      decimal.Decimal `json:"multiplier"`
	OpenInterest    decimal.Decimal `json:"open_interest"`
	QuoteCurrency   string          `json:"quote_currency"`
	SettlementType  string          `json:"settlement_type"`
	TradeStatus     string          `json:"trade_status"`
	TradeMode       int             `json:"trade_mode"`
	PricePrecision  int             `json:"price_precision"`
	MinOrderVolume  decimal.Decimal `json:"min_order_volume"`
	MaxOrderVolume  decimal.Decimal `json:"max_order_volume"`
	StepOrderVolume decimal.Decimal `json:"step_order_volume"`
}

// ListExpirationsService -- GET /api/v4/stock_options/market/expirations (public)
//
// Lists available expiration dates and their Unix timestamps for an underlying.
type ListExpirationsService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewListExpirationsService(underlying string) *ListExpirationsService {
	return &ListExpirationsService{c: c, params: map[string]string{
		"underlying": underlying,
	}}
}

func (s *ListExpirationsService) Do(ctx context.Context) (*StockOptionExpirationsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/market/expirations", s.params)
	return request.Do[StockOptionExpirationsResponse](req)
}

// StockOptionExpirationsResponse is the envelope of the expiration-dates query.
type StockOptionExpirationsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Underlying string                  `json:"underlying"`
		List       []StockOptionExpiration `json:"list"`
	} `json:"data"`
}

// StockOptionExpiration is a single expiration date and its Unix timestamp.
type StockOptionExpiration struct {
	ExpirationDate   string    `json:"expiration_date"`
	ExpirationDateTS time.Time `json:"expiration_date_ts,format:unix"`
}

// GetOptionChainService -- GET /api/v4/stock_options/market/chain (public)
//
// Queries option quotes (bid/ask, mark price, IV, Greeks) across strike prices
// for an underlying and expiration date (YYYY-MM-DD).
type GetOptionChainService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewGetOptionChainService(underlying, expirationDate string) *GetOptionChainService {
	return &GetOptionChainService{c: c, params: map[string]string{
		"underlying":      underlying,
		"expiration_date": expirationDate,
	}}
}

// SetOptionType filters by option type ("call", "put" or "all"). Omit for both.
func (s *GetOptionChainService) SetOptionType(optionType string) *GetOptionChainService {
	s.params["option_type"] = optionType
	return s
}

func (s *GetOptionChainService) Do(ctx context.Context) (*StockOptionChainResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/market/chain", s.params)
	return request.Do[StockOptionChainResponse](req)
}

// StockOptionChainResponse is the envelope of the option-chain query.
type StockOptionChainResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Underlying       string                 `json:"underlying"`
		ExpirationDate   string                 `json:"expiration_date"`
		ExpirationDateTS time.Time              `json:"expiration_date_ts,format:unix"`
		IndexPrice       decimal.Decimal        `json:"index_price"`
		List             []StockOptionChainItem `json:"list"`
	} `json:"data"`
}

// StockOptionChainItem is one contract's quote and Greeks inside an option
// chain. Unavailable numeric values (e.g. bid/ask/rho/breakeven when no quote
// exists) are sent as "" and decoded as zero.
type StockOptionChainItem struct {
	Symbol         string          `json:"symbol"`
	SymbolDisplay  string          `json:"symbol_display"`
	StrikePrice    decimal.Decimal `json:"strike_price"`
	OptionType     string          `json:"option_type"`
	Bid            decimal.Decimal `json:"bid"`
	Ask            decimal.Decimal `json:"ask"`
	MarkPrice      decimal.Decimal `json:"mark_price"`
	MarkIV         decimal.Decimal `json:"mark_iv"`
	Last           decimal.Decimal `json:"last"`
	DisplayPrice   decimal.Decimal `json:"display_price"`
	BreakevenPrice decimal.Decimal `json:"breakeven_price"`
	ToBreakevenPct decimal.Decimal `json:"to_breakeven_pct"`
	TodayChangePct decimal.Decimal `json:"today_change_pct"`
	OpenInterest   decimal.Decimal `json:"open_interest"`
	Position       decimal.Decimal `json:"position"`
	Delta          decimal.Decimal `json:"delta"`
	Gamma          decimal.Decimal `json:"gamma"`
	Theta          decimal.Decimal `json:"theta"`
	Vega           decimal.Decimal `json:"vega"`
	Rho            decimal.Decimal `json:"rho"`
	TradeStatus    string          `json:"trade_status"`
	TradeMode      int             `json:"trade_mode"`
}

// GetOrderBookService -- GET /api/v4/stock_options/market/{symbol}/orderbook (public)
//
// Queries the best bid and ask for an option contract symbol. Each side holds
// at most one price level (without quantities) or is empty when unavailable.
type GetOrderBookService struct {
	c      *StockOptionsClient
	symbol string
}

func (c *StockOptionsClient) NewGetOrderBookService(symbol string) *GetOrderBookService {
	return &GetOrderBookService{c: c, symbol: symbol}
}

func (s *GetOrderBookService) Do(ctx context.Context) (*StockOptionOrderBookResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/market/"+s.symbol+"/orderbook")
	return request.Do[StockOptionOrderBookResponse](req)
}

// StockOptionOrderBookResponse is the envelope of the option order-book query.
type StockOptionOrderBookResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Symbol        string                      `json:"symbol"`
		SymbolDisplay string                      `json:"symbol_display"`
		Bids          []StockOptionOrderBookLevel `json:"bids"`
		Asks          []StockOptionOrderBookLevel `json:"asks"`
	} `json:"data"`
}

// StockOptionOrderBookLevel is a single price level in an option order book.
type StockOptionOrderBookLevel struct {
	Price decimal.Decimal `json:"price"`
}
