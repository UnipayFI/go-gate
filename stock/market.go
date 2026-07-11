package stock

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListSymbolsService -- GET /api/v4/stock/symbols (public)
//
// Lists the tradable stock symbols with their quote currency, trading status
// and precision, paginated.
type ListSymbolsService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListSymbolsService() *ListSymbolsService {
	return &ListSymbolsService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of symbols.
func (s *ListSymbolsService) SetSymbols(symbols string) *ListSymbolsService {
	s.params["symbols"] = symbols
	return s
}

// SetExchange narrows the result to an exchange ("us", "hk" or "kr").
func (s *ListSymbolsService) SetExchange(exchange string) *ListSymbolsService {
	s.params["exchange"] = exchange
	return s
}

// SetWithDescI18n returns the multilingual symbol description when set to true.
func (s *ListSymbolsService) SetWithDescI18n(withDescI18n bool) *ListSymbolsService {
	s.params["with_desc_i18n"] = strconv.FormatBool(withDescI18n)
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListSymbolsService) SetPage(page int) *ListSymbolsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (defaults to 10, max 500).
func (s *ListSymbolsService) SetPageSize(pageSize int) *ListSymbolsService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListSymbolsService) Do(ctx context.Context) (*StockSymbolsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/symbols", s.params)
	return request.Do[StockSymbolsResponse](req)
}

// StockSymbolsResponse is the envelope of the symbol-list query.
type StockSymbolsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64         `json:"total"`
		TotalPage int           `json:"total_page"`
		List      []StockSymbol `json:"list"`
	} `json:"data"`
}

// StockSymbol is a single tradable symbol. trade_mode is the current session
// trading mode (0=disabled, 1=buy only, 2=sell only, 4=buy and sell);
// order_fill_timing is 1=immediate, 2=after pre-market opens, 3=after regular
// session opens.
type StockSymbol struct {
	Symbol                 string            `json:"symbol"`
	Exchange               string            `json:"exchange"`
	ExchangeDesc           string            `json:"exchange_desc"`
	QuoteCurrency          string            `json:"quote_currency"`
	QuoteCurrencyPrecision int               `json:"quote_currency_precision"`
	FXRate                 decimal.Decimal   `json:"fx_rate"`
	SymbolDesc             string            `json:"symbol_desc"`
	Category               string            `json:"category"`
	TradeStatus            string            `json:"trade_status"`
	TradeMode              int               `json:"trade_mode"`
	OrderFillTiming        int               `json:"order_fill_timing"`
	IconLink               string            `json:"icon_link"`
	QuoteCurrencySymbol    string            `json:"quote_currency_symbol"`
	PricePrecision         int               `json:"price_precision"`
	VolumePrecision        int               `json:"volume_precision"`
	IsIPO                  bool              `json:"is_ipo"`
	IPOPrice               decimal.Decimal   `json:"ipo_price"`
	SellPriceProtection    decimal.Decimal   `json:"sell_price_protection"`
	BuyPriceProtection     decimal.Decimal   `json:"buy_price_protection"`
	SymbolDescs            []StockSymbolDesc `json:"symbol_descs"`
}

// StockSymbolDesc is one localized symbol description.
type StockSymbolDesc struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// ListSymbolDetailsService -- GET /api/v4/stock/symbols/detail (public)
//
// Returns detailed contract specifications (order-volume bounds, protection and
// fee rates) for stock symbols, paginated.
type ListSymbolDetailsService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListSymbolDetailsService() *ListSymbolDetailsService {
	return &ListSymbolDetailsService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of symbols.
func (s *ListSymbolDetailsService) SetSymbols(symbols string) *ListSymbolDetailsService {
	s.params["symbols"] = symbols
	return s
}

// SetExchange narrows the result to an exchange ("us", "hk" or "kr").
func (s *ListSymbolDetailsService) SetExchange(exchange string) *ListSymbolDetailsService {
	s.params["exchange"] = exchange
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListSymbolDetailsService) SetPage(page int) *ListSymbolDetailsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (defaults to 10, max 500).
func (s *ListSymbolDetailsService) SetPageSize(pageSize int) *ListSymbolDetailsService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListSymbolDetailsService) Do(ctx context.Context) (*StockSymbolDetailsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/symbols/detail", s.params)
	return request.Do[StockSymbolDetailsResponse](req)
}

// StockSymbolDetailsResponse is the envelope of the symbol-details query.
type StockSymbolDetailsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64               `json:"total"`
		TotalPage int                 `json:"total_page"`
		List      []StockSymbolDetail `json:"list"`
	} `json:"data"`
}

// StockSymbolDetail is a single symbol's contract specification.
type StockSymbolDetail struct {
	Symbol                 string            `json:"symbol"`
	Exchange               string            `json:"exchange"`
	ExchangeDesc           string            `json:"exchange_desc"`
	QuoteCurrency          string            `json:"quote_currency"`
	QuoteCurrencyPrecision int               `json:"quote_currency_precision"`
	FXRate                 decimal.Decimal   `json:"fx_rate"`
	SymbolDesc             string            `json:"symbol_desc"`
	Category               string            `json:"category"`
	SettlementCurrency     string            `json:"settlement_currency"`
	MaxOrderVolume         decimal.Decimal   `json:"max_order_volume"`
	StepOrderVolume        decimal.Decimal   `json:"step_order_volume"`
	MinOrderVolume         decimal.Decimal   `json:"min_order_volume"`
	PricePrecision         int               `json:"price_precision"`
	VolumePrecision        int               `json:"volume_precision"`
	IsIPO                  bool              `json:"is_ipo"`
	IPOPrice               decimal.Decimal   `json:"ipo_price"`
	PriceProtection        decimal.Decimal   `json:"price_protection"`
	SellPriceProtection    decimal.Decimal   `json:"sell_price_protection"`
	BuyPriceProtection     decimal.Decimal   `json:"buy_price_protection"`
	SlippageRate           decimal.Decimal   `json:"slippage_rate"`
	CommissionRate         decimal.Decimal   `json:"commission_rate"`
	TradeStatus            string            `json:"trade_status"`
	Status                 string            `json:"status"`
	TradeMode              int               `json:"trade_mode"`
	OrderFillTiming        int               `json:"order_fill_timing"`
	SymbolDescs            []StockSymbolDesc `json:"symbol_descs"`
	IconLink               string            `json:"icon_link"`
}

// GetOrderBookService -- GET /api/v4/stock/market/{symbol}/orderbook (public)
//
// Returns the market order book (bids and asks) for a symbol; user_order marks
// whether a level belongs to the caller's own order.
type GetOrderBookService struct {
	c      *StockClient
	symbol string
}

func (c *StockClient) NewGetOrderBookService(symbol string) *GetOrderBookService {
	return &GetOrderBookService{c: c, symbol: symbol}
}

func (s *GetOrderBookService) Do(ctx context.Context) (*StockOrderBookResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/market/"+s.symbol+"/orderbook")
	return request.Do[StockOrderBookResponse](req)
}

// StockOrderBookResponse is the envelope of the order-book query.
type StockOrderBookResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Symbol string                `json:"symbol"`
		Bids   []StockOrderBookLevel `json:"bids"`
		Asks   []StockOrderBookLevel `json:"asks"`
	} `json:"data"`
}

// StockOrderBookLevel is a single price level in the order book.
type StockOrderBookLevel struct {
	Price     decimal.Decimal `json:"p"`
	UserOrder bool            `json:"user_order"`
}

// ListExchangesService -- GET /api/v4/stock/exchanges (private)
//
// Lists the supported exchanges ("us", "hk", "kr") and whether each supports
// stock transfer. Despite being documented as public, the live endpoint
// requires a signed request.
type ListExchangesService struct {
	c *StockClient
}

func (c *StockClient) NewListExchangesService() *ListExchangesService {
	return &ListExchangesService{c: c}
}

func (s *ListExchangesService) Do(ctx context.Context) (*StockExchangesResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/exchanges").WithSign()
	return request.Do[StockExchangesResponse](req)
}

// StockExchangesResponse is the envelope of the supported-exchanges query.
type StockExchangesResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []StockExchange `json:"list"`
	} `json:"data"`
}

// StockExchange is a single supported exchange.
type StockExchange struct {
	Exchange        string `json:"exchange"`
	ExchangeDesc    string `json:"exchange_desc"`
	IconLink        string `json:"icon_link"`
	SupportTransfer bool   `json:"support_transfer"`
}

// GetFeeRateService -- GET /api/v4/stock/fee-rate (public)
//
// Returns the maker/taker fee rates per VIP level.
type GetFeeRateService struct {
	c *StockClient
}

func (c *StockClient) NewGetFeeRateService() *GetFeeRateService {
	return &GetFeeRateService{c: c}
}

func (s *GetFeeRateService) Do(ctx context.Context) (*StockFeeRateResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/fee-rate")
	return request.Do[StockFeeRateResponse](req)
}

// StockFeeRateResponse is the envelope of the fee-rate query.
type StockFeeRateResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []StockFeeRate `json:"list"`
	} `json:"data"`
}

// StockFeeRate is a single VIP level's maker/taker fee rate.
type StockFeeRate struct {
	VIPLevel int             `json:"vip_level"`
	MakerFee decimal.Decimal `json:"maker_fee"`
	TakerFee decimal.Decimal `json:"taker_fee"`
}
