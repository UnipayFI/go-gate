package stock

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListPositionsService -- GET /api/v4/stock/positions (private)
//
// Lists the caller's active positions. pnlCalcType selects the cost basis
// (1=average cost, 2=diluted cost); pnlCalcPrice selects the price basis
// (1=intraday, 2=latest extended-hours).
type ListPositionsService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListPositionsService() *ListPositionsService {
	return &ListPositionsService{c: c, params: map[string]string{}}
}

// SetPnLCalcType selects the PnL cost basis (1=average cost, 2=diluted cost).
func (s *ListPositionsService) SetPnLCalcType(pnlCalcType int) *ListPositionsService {
	s.params["pnl_calc_type"] = strconv.Itoa(pnlCalcType)
	return s
}

// SetPnLCalcPrice selects the PnL price basis (1=intraday, 2=latest
// extended-hours).
func (s *ListPositionsService) SetPnLCalcPrice(pnlCalcPrice int) *ListPositionsService {
	s.params["pnl_calc_price"] = strconv.Itoa(pnlCalcPrice)
	return s
}

// SetSymbol narrows the result to a single symbol.
func (s *ListPositionsService) SetSymbol(symbol string) *ListPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetExchange narrows the result to an exchange ("us", "hk" or "kr").
func (s *ListPositionsService) SetExchange(exchange string) *ListPositionsService {
	s.params["exchange"] = exchange
	return s
}

func (s *ListPositionsService) Do(ctx context.Context) (*StockPositionsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/positions", s.params).WithSign()
	return request.Do[StockPositionsResponse](req)
}

// StockPositionsResponse is the envelope of the active-position query.
type StockPositionsResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []StockPosition `json:"list"`
	} `json:"data"`
}

// StockPosition is a single active position. extended_last_price is null
// outside extended-hours trading.
type StockPosition struct {
	Symbol                 string          `json:"symbol"`
	Exchange               string          `json:"exchange"`
	QuoteCurrency          string          `json:"quote_currency"`
	QuoteCurrencyPrecision int             `json:"quote_currency_precision"`
	FXRate                 decimal.Decimal `json:"fx_rate"`
	TradeStatus            string          `json:"trade_status"`
	SymbolDesc             string          `json:"symbol_desc"`
	PositionPnL            decimal.Decimal `json:"position_pnl"`
	TodayPnL               decimal.Decimal `json:"today_pnl"`
	PnLRate                decimal.Decimal `json:"pnl_rate"`
	TodaySellAmount        decimal.Decimal `json:"today_sell_amount"`
	TodayBuyAmount         decimal.Decimal `json:"today_buy_amount"`
	TodaySellVolume        decimal.Decimal `json:"today_sell_volume"`
	TodayBuyVolume         decimal.Decimal `json:"today_buy_volume"`
	YesterdayVolume        decimal.Decimal `json:"yesterday_volume"`
	Volume                 decimal.Decimal `json:"volume"`
	Available              decimal.Decimal `json:"available"`
	TransferOutPendingQty  decimal.Decimal `json:"transfer_out_pending_qty"`
	AvgCostPrice           decimal.Decimal `json:"avg_cost_price"`
	DilutedCostPrice       decimal.Decimal `json:"diluted_cost_price"`
	LastPrice              decimal.Decimal `json:"last_price"`
	ExtendedLastPrice      decimal.Decimal `json:"extended_last_price"`
	MaxOrderVolume         decimal.Decimal `json:"max_order_volume"`
	StepOrderVolume        decimal.Decimal `json:"step_order_volume"`
	MinOrderVolume         decimal.Decimal `json:"min_order_volume"`
	PricePrecision         int             `json:"price_precision"`
	PriceProtection        decimal.Decimal `json:"price_protection"`
	SellPriceProtection    decimal.Decimal `json:"sell_price_protection"`
	BuyPriceProtection     decimal.Decimal `json:"buy_price_protection"`
	CommissionRate         decimal.Decimal `json:"commission_rate"`
	SlippageRate           decimal.Decimal `json:"slippage_rate"`
}

// ClosePositionService -- POST /api/v4/stock/positions/close (private)
//
// Closes a position. closeType is 1 for a partial close (closeVolume required)
// or 2 for a full close.
type ClosePositionService struct {
	c    *StockClient
	body map[string]any
}

func (c *StockClient) NewClosePositionService(symbol string, closeType int) *ClosePositionService {
	return &ClosePositionService{c: c, body: map[string]any{
		"symbol":     symbol,
		"close_type": closeType,
	}}
}

// SetCloseVolume sets the volume to close (required when closeType is 1).
func (s *ClosePositionService) SetCloseVolume(closeVolume decimal.Decimal) *ClosePositionService {
	s.body["close_volume"] = closeVolume.String()
	return s
}

func (s *ClosePositionService) Do(ctx context.Context) (*StockClosePositionResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/stock/positions/close", s.body).WithSign()
	return request.Do[StockClosePositionResponse](req)
}

// StockClosePositionResponse is the envelope of the close-position request.
type StockClosePositionResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		OrderID int64 `json:"order_id"`
	} `json:"data"`
}
