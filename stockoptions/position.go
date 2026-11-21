package stockoptions

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListPositionsService -- GET /api/v4/stock_options/positions (private)
//
// Lists the current user's option positions, market values and PnL, sorted by
// open_time descending then symbol ascending.
type ListPositionsService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewListPositionsService() *ListPositionsService {
	return &ListPositionsService{c: c, params: map[string]string{}}
}

// SetUnderlying narrows the result to a single option underlying (e.g. "AAPL").
func (s *ListPositionsService) SetUnderlying(underlying string) *ListPositionsService {
	s.params["underlying"] = underlying
	return s
}

// SetSymbol narrows the result to a full option contract symbol.
func (s *ListPositionsService) SetSymbol(symbol string) *ListPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetPriceType selects the valuation price for market value, unrealized PnL and
// daily PnL: "latest" (default, last traded price) or "mark" (mark price).
func (s *ListPositionsService) SetPriceType(priceType string) *ListPositionsService {
	s.params["price_type"] = priceType
	return s
}

func (s *ListPositionsService) Do(ctx context.Context) (*StockOptionPositionsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/positions", s.params).WithSign()
	return request.Do[StockOptionPositionsResponse](req)
}

// StockOptionPositionsResponse is the envelope of the active-positions query.
// Data.TodayPnL is the aggregate daily PnL across current positions and
// contracts fully closed today.
type StockOptionPositionsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		TodayPnL decimal.Decimal       `json:"today_pnl"`
		List     []StockOptionPosition `json:"list"`
	} `json:"data"`
}

// StockOptionPosition is a single active stock-option position. side is "long"
// or "short"; open_time is Unix seconds (0 when unavailable); trade_status is
// "open" or "closed"; trade_mode is 0 (disabled) or 4 (buy/sell allowed).
type StockOptionPosition struct {
	Symbol          string          `json:"symbol"`
	SymbolDisplay   string          `json:"symbol_display"`
	Underlying      string          `json:"underlying"`
	Side            string          `json:"side"`
	OpenTime        time.Time       `json:"open_time,format:unix"`
	Quantity        decimal.Decimal `json:"quantity"`
	Available       decimal.Decimal `json:"available"`
	Frozen          decimal.Decimal `json:"frozen"`
	AvgEntryPrice   decimal.Decimal `json:"avg_entry_price"`
	LastPrice       decimal.Decimal `json:"last_price"`
	MarkPrice       decimal.Decimal `json:"mark_price"`
	PreviousClose   decimal.Decimal `json:"previous_close"`
	MarketValue     decimal.Decimal `json:"market_value"`
	UnrealizedPnL   decimal.Decimal `json:"unrealized_pnl"`
	TodayPnL        decimal.Decimal `json:"today_pnl"`
	TodayBuyQty     decimal.Decimal `json:"today_buy_qty"`
	TodaySellQty    decimal.Decimal `json:"today_sell_qty"`
	TodayBuyValue   decimal.Decimal `json:"today_buy_value"`
	TodaySellValue  decimal.Decimal `json:"today_sell_value"`
	TradeStatus     string          `json:"trade_status"`
	TradeMode       int             `json:"trade_mode"`
	Multiplier      decimal.Decimal `json:"multiplier"`
	PricePrecision  int             `json:"price_precision"`
	MinOrderVolume  decimal.Decimal `json:"min_order_volume"`
	MaxOrderVolume  decimal.Decimal `json:"max_order_volume"`
	StepOrderVolume decimal.Decimal `json:"step_order_volume"`
}

// ClosePositionService -- POST /api/v4/stock_options/positions/close (private)
//
// Submits a closing order for part or all of an option position. Omit
// order_type for a market closing order; set order_type to "limit" with a
// positive price for a limit closing order. Time in force is fixed to "day".
type ClosePositionService struct {
	c    *StockOptionsClient
	body map[string]any
}

func (c *StockOptionsClient) NewClosePositionService(symbol string, quantity decimal.Decimal) *ClosePositionService {
	return &ClosePositionService{c: c, body: map[string]any{
		"symbol":   symbol,
		"quantity": quantity.String(),
	}}
}

// SetOrderType sets the closing order type ("market" or "limit"). Defaults to
// "market" when omitted.
func (s *ClosePositionService) SetOrderType(orderType string) *ClosePositionService {
	s.body["order_type"] = orderType
	return s
}

// SetPrice sets the limit closing price (required for limit orders, ignored for
// market orders).
func (s *ClosePositionService) SetPrice(price decimal.Decimal) *ClosePositionService {
	s.body["price"] = price.String()
	return s
}

// SetClientOrderID sets a custom closing order identifier.
func (s *ClosePositionService) SetClientOrderID(clientOrderID string) *ClosePositionService {
	s.body["client_order_id"] = clientOrderID
	return s
}

func (s *ClosePositionService) Do(ctx context.Context) (*StockOptionClosePositionResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/stock_options/positions/close", s.body).WithSign()
	return request.Do[StockOptionClosePositionResponse](req)
}

// StockOptionClosePositionResponse is the envelope of the close-position request.
type StockOptionClosePositionResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		OrderID string `json:"order_id"`
	} `json:"data"`
}
