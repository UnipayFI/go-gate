package tradfi

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListCategoriesService -- GET /api/v4/tradfi/symbols/categories (private)
//
// Lists the trading-symbol categories.
type ListCategoriesService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewListCategoriesService() *ListCategoriesService {
	return &ListCategoriesService{c: c}
}

func (s *ListCategoriesService) Do(ctx context.Context) (*TradfiCategoriesResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/symbols/categories").WithSign()
	return request.Do[TradfiCategoriesResponse](req)
}

// TradfiCategoriesResponse is the envelope of the symbol-categories query.
type TradfiCategoriesResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiCategory `json:"list"`
	} `json:"data"`
}

// TradfiCategory is a single trading-symbol category.
type TradfiCategory struct {
	CategoryID   int    `json:"category_id"`
	IsFavorite   bool   `json:"is_favorite"`
	CategoryName string `json:"category_name"`
}

// ListCommissionsService -- GET /api/v4/tradfi/symbols/commissions (private)
//
// Lists per-lot commission rates for symbols and/or categories. At least one of
// symbols or categoryCode must be provided.
type ListCommissionsService struct {
	c      *TradfiClient
	params map[string]string
}

func (c *TradfiClient) NewListCommissionsService() *ListCommissionsService {
	return &ListCommissionsService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of symbol codes.
func (s *ListCommissionsService) SetSymbols(symbols string) *ListCommissionsService {
	s.params["symbols"] = symbols
	return s
}

// SetCategoryCode narrows the result to a comma-separated list of category codes.
func (s *ListCommissionsService) SetCategoryCode(categoryCode string) *ListCommissionsService {
	s.params["category_code"] = categoryCode
	return s
}

func (s *ListCommissionsService) Do(ctx context.Context) (*TradfiCommissionsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/symbols/commissions", s.params).WithSign()
	return request.Do[TradfiCommissionsResponse](req)
}

// TradfiCommissionsResponse is the envelope of the symbol-commissions query.
type TradfiCommissionsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiCommission `json:"list"`
	} `json:"data"`
}

// TradfiCommission is a single symbol's per-lot commission rate.
type TradfiCommission struct {
	CategoryCode string          `json:"category_code"`
	Symbol       string          `json:"symbol"`
	FeePerLot    decimal.Decimal `json:"fee_per_lot"`
}

// ListSymbolsService -- GET /api/v4/tradfi/symbols (private)
//
// Lists the tradable symbols with their status and trading windows.
type ListSymbolsService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewListSymbolsService() *ListSymbolsService {
	return &ListSymbolsService{c: c}
}

func (s *ListSymbolsService) Do(ctx context.Context) (*TradfiSymbolsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/symbols").WithSign()
	return request.Do[TradfiSymbolsResponse](req)
}

// TradfiSymbolsResponse is the envelope of the symbol-list query.
type TradfiSymbolsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiSymbol `json:"list"`
	} `json:"data"`
}

// TradfiSymbol is a single tradable symbol. close_time / open_time /
// next_open_time are integer-second Unix timestamps (next_open_time 0 = none).
type TradfiSymbol struct {
	Symbol                   string            `json:"symbol"`
	SymbolDesc               string            `json:"symbol_desc"`
	CategoryID               int               `json:"category_id"`
	Status                   string            `json:"status"`
	TradeMode                string            `json:"trade_mode"`
	IconLink                 string            `json:"icon_link"`
	CloseTime                time.Time         `json:"close_time,format:unix"`
	OpenTime                 time.Time         `json:"open_time,format:unix"`
	NextOpenTime             time.Time         `json:"next_open_time,format:unix"`
	SettlementCurrency       string            `json:"settlement_currency"`
	SettlementCurrencySymbol string            `json:"settlement_currency_symbol"`
	PricePrecision           int               `json:"price_precision"`
	IsBase                   bool              `json:"is_base"`
	Leverages                []decimal.Decimal `json:"leverages"`
	SymbolDescs              []string          `json:"symbol_descs"`
}

// ListSymbolDetailsService -- GET /api/v4/tradfi/symbols/detail (private)
//
// Returns detailed contract specifications for up to 10 comma-separated symbols.
type ListSymbolDetailsService struct {
	c      *TradfiClient
	params map[string]string
}

func (c *TradfiClient) NewListSymbolDetailsService(symbols string) *ListSymbolDetailsService {
	return &ListSymbolDetailsService{c: c, params: map[string]string{
		"symbols": symbols,
	}}
}

func (s *ListSymbolDetailsService) Do(ctx context.Context) (*TradfiSymbolDetailsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/symbols/detail", s.params).WithSign()
	return request.Do[TradfiSymbolDetailsResponse](req)
}

// TradfiSymbolDetailsResponse is the envelope of the symbol-details query.
type TradfiSymbolDetailsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiSymbolDetail `json:"list"`
	} `json:"data"`
}

// TradfiSymbolDetail is a single symbol's contract specification. leverage is
// returned as a free-form string (e.g. a ratio), so it is kept as a string.
type TradfiSymbolDetail struct {
	Symbol             string          `json:"symbol"`
	SymbolDesc         string          `json:"symbol_desc"`
	CategoryName       string          `json:"category_name"`
	ContractVolume     decimal.Decimal `json:"contract_volume"`
	SettlementCurrency string          `json:"settlement_currency"`
	MaxOrderVolume     decimal.Decimal `json:"max_order_volume"`
	MinOrderVolume     decimal.Decimal `json:"min_order_volume"`
	Leverage           string          `json:"leverage"`
	PricePrecision     int             `json:"price_precision"`
	PriceSLLevel       decimal.Decimal `json:"price_sl_level"`
	SwapCostType       string          `json:"swap_cost_type"`
	BuySwapCostRate    decimal.Decimal `json:"buy_swap_cost_rate"`
	SellSwapCostRate   decimal.Decimal `json:"sell_swap_cost_rate"`
	SwapCost3Day       decimal.Decimal `json:"swap_cost_3day"`
	TradeTimezone      string          `json:"trade_timezone"`
	TradeMode          string          `json:"trade_mode"`
	IconLink           string          `json:"icon_link"`
}

// ListKlinesService -- GET /api/v4/tradfi/symbols/{symbol}/klines (private)
//
// Returns the candlestick series for a symbol at the requested period.
type ListKlinesService struct {
	c      *TradfiClient
	symbol string
	params map[string]string
}

func (c *TradfiClient) NewListKlinesService(symbol, klineType string) *ListKlinesService {
	return &ListKlinesService{c: c, symbol: symbol, params: map[string]string{
		"kline_type": klineType,
	}}
}

// SetBeginTime bounds the series to candles at or after this time.
func (s *ListKlinesService) SetBeginTime(beginTime time.Time) *ListKlinesService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the series to candles at or before this time.
func (s *ListKlinesService) SetEndTime(endTime time.Time) *ListKlinesService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetLimit caps the number of candles returned (max 500).
func (s *ListKlinesService) SetLimit(limit int) *ListKlinesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListKlinesService) Do(ctx context.Context) (*TradfiKlinesResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/symbols/"+s.symbol+"/klines", s.params).WithSign()
	return request.Do[TradfiKlinesResponse](req)
}

// TradfiKlinesResponse is the envelope of the kline query.
type TradfiKlinesResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiKline `json:"list"`
	} `json:"data"`
}

// TradfiKline is a single candlestick. t is the integer-second Unix timestamp.
type TradfiKline struct {
	Open  decimal.Decimal `json:"o"`
	Close decimal.Decimal `json:"c"`
	High  decimal.Decimal `json:"h"`
	Low   decimal.Decimal `json:"l"`
	Time  time.Time       `json:"t,format:unix"`
}

// GetTickerService -- GET /api/v4/tradfi/symbols/{symbol}/tickers (private)
//
// Returns the latest ticker snapshot for a symbol.
type GetTickerService struct {
	c      *TradfiClient
	symbol string
}

func (c *TradfiClient) NewGetTickerService(symbol string) *GetTickerService {
	return &GetTickerService{c: c, symbol: symbol}
}

func (s *GetTickerService) Do(ctx context.Context) (*TradfiTickerResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/symbols/"+s.symbol+"/tickers").WithSign()
	return request.Do[TradfiTickerResponse](req)
}

// TradfiTickerResponse is the envelope of the ticker query.
type TradfiTickerResponse struct {
	Label     string       `json:"label"`
	Message   string       `json:"message"`
	Timestamp time.Time    `json:"timestamp,format:unixmilli"`
	Data      TradfiTicker `json:"data"`
}

// TradfiTicker is a symbol's latest ticker. price_change is the percentage
// change multiplied by 100. close_time / open_time / next_open_time are
// integer-second Unix timestamps (next_open_time 0 = none).
type TradfiTicker struct {
	HighestPrice        decimal.Decimal `json:"highest_price"`
	LowestPrice         decimal.Decimal `json:"lowest_price"`
	PriceChange         decimal.Decimal `json:"price_change"`
	PriceChangeAmount   decimal.Decimal `json:"price_change_amount"`
	TodayOpenPrice      decimal.Decimal `json:"today_open_price"`
	LastTodayClosePrice decimal.Decimal `json:"last_today_close_price"`
	LastPrice           decimal.Decimal `json:"last_price"`
	BidPrice            decimal.Decimal `json:"bid_price"`
	AskPrice            decimal.Decimal `json:"ask_price"`
	Favorite            bool            `json:"favorite"`
	Status              string          `json:"status"`
	CloseTime           time.Time       `json:"close_time,format:unix"`
	OpenTime            time.Time       `json:"open_time,format:unix"`
	NextOpenTime        time.Time       `json:"next_open_time,format:unix"`
	TradeMode           string          `json:"trade_mode"`
	CategoryName        string          `json:"category_name"`
}
