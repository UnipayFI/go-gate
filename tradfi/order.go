package tradfi

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListOrdersService -- GET /api/v4/tradfi/orders (private)
//
// Lists the caller's active (unfinished) orders.
type ListOrdersService struct {
	c *TradfiClient
}

func (c *TradfiClient) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c}
}

func (s *ListOrdersService) Do(ctx context.Context) (*TradfiOrdersResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/orders").WithSign()
	return request.Do[TradfiOrdersResponse](req)
}

// TradfiOrdersResponse is the envelope of the active-order query.
type TradfiOrdersResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiOrder `json:"list"`
	} `json:"data"`
}

// TradfiOrder is a single active order. side is 1=sell, 2=buy; time_setup is an
// integer-second Unix timestamp.
type TradfiOrder struct {
	OrderID    int64           `json:"order_id"`
	Symbol     string          `json:"symbol"`
	SymbolDesc string          `json:"symbol_desc"`
	PriceType  string          `json:"price_type"`
	State      int             `json:"state"`
	StateDesc  string          `json:"state_desc"`
	Finished   int             `json:"finished"`
	Side       int             `json:"side"`
	Volume     decimal.Decimal `json:"volume"`
	Price      decimal.Decimal `json:"price"`
	PriceTP    decimal.Decimal `json:"price_tp"`
	PriceSL    decimal.Decimal `json:"price_sl"`
	TimeSetup  time.Time       `json:"time_setup,format:unix"`
}

// CreateOrderService -- POST /api/v4/tradfi/orders (private)
//
// Places a single order. side is 1=sell, 2=buy; priceType is "trigger" or
// "market".
type CreateOrderService struct {
	c    *TradfiClient
	body map[string]any
}

func (c *TradfiClient) NewCreateOrderService(symbol string, side int, priceType string, price, volume decimal.Decimal) *CreateOrderService {
	return &CreateOrderService{c: c, body: map[string]any{
		"symbol":     symbol,
		"side":       side,
		"price_type": priceType,
		"price":      price.String(),
		"volume":     volume.String(),
	}}
}

// SetPriceTP sets the take-profit price.
func (s *CreateOrderService) SetPriceTP(priceTP decimal.Decimal) *CreateOrderService {
	s.body["price_tp"] = priceTP.String()
	return s
}

// SetPriceSL sets the stop-loss price.
func (s *CreateOrderService) SetPriceSL(priceSL decimal.Decimal) *CreateOrderService {
	s.body["price_sl"] = priceSL.String()
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*TradfiCreateOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/tradfi/orders", s.body).WithSign()
	return request.Do[TradfiCreateOrderResponse](req)
}

// TradfiCreateOrderResponse is the envelope of the create-order request.
type TradfiCreateOrderResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		ID string `json:"id"`
	} `json:"data"`
}

// ModifyOrderService -- PUT /api/v4/tradfi/orders/{order_id} (private)
//
// Modifies an active order's price and optional take-profit/stop-loss prices.
type ModifyOrderService struct {
	c       *TradfiClient
	orderID int64
	body    map[string]any
}

func (c *TradfiClient) NewModifyOrderService(orderID int64, price decimal.Decimal) *ModifyOrderService {
	return &ModifyOrderService{c: c, orderID: orderID, body: map[string]any{
		"price": price.String(),
	}}
}

// SetPriceTP sets the take-profit price. Passing "0" clears the existing one.
func (s *ModifyOrderService) SetPriceTP(priceTP decimal.Decimal) *ModifyOrderService {
	s.body["price_tp"] = priceTP.String()
	return s
}

// SetPriceSL sets the stop-loss price. Passing "0" clears the existing one.
func (s *ModifyOrderService) SetPriceSL(priceSL decimal.Decimal) *ModifyOrderService {
	s.body["price_sl"] = priceSL.String()
	return s
}

func (s *ModifyOrderService) Do(ctx context.Context) (*TradfiModifyOrderResponse, error) {
	req := request.Put(ctx, s.c, "/api/v4/tradfi/orders/"+strconv.FormatInt(s.orderID, 10), s.body).WithSign()
	return request.Do[TradfiModifyOrderResponse](req)
}

// TradfiModifyOrderResponse is the envelope of the modify-order request.
type TradfiModifyOrderResponse struct {
	Timestamp time.Time               `json:"timestamp,format:unixmilli"`
	Data      TradfiModifyOrderResult `json:"data"`
}

// TradfiModifyOrderResult is the order state after modification. On this
// endpoint state is returned as a string.
type TradfiModifyOrderResult struct {
	OrderID int64           `json:"order_id"`
	Symbol  string          `json:"symbol"`
	State   string          `json:"state"`
	Volume  decimal.Decimal `json:"volume"`
	Price   decimal.Decimal `json:"price"`
	PriceTP decimal.Decimal `json:"price_tp"`
	PriceSL decimal.Decimal `json:"price_sl"`
}

// CancelOrderService -- DELETE /api/v4/tradfi/orders/{order_id} (private)
//
// Cancels an active order. Gate returns no meaningful body on success.
type CancelOrderService struct {
	c       *TradfiClient
	orderID int64
}

func (c *TradfiClient) NewCancelOrderService(orderID int64) *CancelOrderService {
	return &CancelOrderService{c: c, orderID: orderID}
}

func (s *CancelOrderService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v4/tradfi/orders/"+strconv.FormatInt(s.orderID, 10)).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListOrderHistoryService -- GET /api/v4/tradfi/orders/history (private)
//
// Lists the caller's historical orders (earliest queryable one month ago).
type ListOrderHistoryService struct {
	c      *TradfiClient
	params map[string]string
}

func (c *TradfiClient) NewListOrderHistoryService() *ListOrderHistoryService {
	return &ListOrderHistoryService{c: c, params: map[string]string{}}
}

// SetBeginTime bounds the result to orders at or after this time.
func (s *ListOrderHistoryService) SetBeginTime(beginTime time.Time) *ListOrderHistoryService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to orders at or before this time.
func (s *ListOrderHistoryService) SetEndTime(endTime time.Time) *ListOrderHistoryService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetSymbol narrows the result to a single symbol.
func (s *ListOrderHistoryService) SetSymbol(symbol string) *ListOrderHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetSide narrows the result to an order side (1=sell, 2=buy).
func (s *ListOrderHistoryService) SetSide(side int) *ListOrderHistoryService {
	s.params["side"] = strconv.Itoa(side)
	return s
}

func (s *ListOrderHistoryService) Do(ctx context.Context) (*TradfiOrderHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/orders/history", s.params).WithSign()
	return request.Do[TradfiOrderHistoryResponse](req)
}

// TradfiOrderHistoryResponse is the envelope of the historical-order query.
type TradfiOrderHistoryResponse struct {
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []TradfiOrderHistory `json:"list"`
	} `json:"data"`
}

// TradfiOrderHistory is a single historical order. time_setup / time_done are
// integer-second Unix timestamps; side is 1=sell, 2=buy.
type TradfiOrderHistory struct {
	OrderID      int64           `json:"order_id"`
	Symbol       string          `json:"symbol"`
	SymbolDesc   string          `json:"symbol_desc"`
	PriceType    string          `json:"price_type"`
	OrderOptType int             `json:"order_opt_type"`
	State        int             `json:"state"`
	StateDesc    string          `json:"state_desc"`
	Side         int             `json:"side"`
	Volume       decimal.Decimal `json:"volume"`
	FillVolume   decimal.Decimal `json:"fill_volume"`
	ClosePnL     decimal.Decimal `json:"close_pnl"`
	Price        decimal.Decimal `json:"price"`
	TriggerPrice decimal.Decimal `json:"trigger_price"`
	PriceTP      decimal.Decimal `json:"price_tp"`
	PriceSL      decimal.Decimal `json:"price_sl"`
	TimeSetup    time.Time       `json:"time_setup,format:unix"`
	TimeDone     time.Time       `json:"time_done,format:unix"`
}

// GetOrderLogService -- GET /api/v4/tradfi/orders/log/{log_id} (private)
//
// Returns an order's details by the log id returned from order placement.
type GetOrderLogService struct {
	c     *TradfiClient
	logID int64
}

func (c *TradfiClient) NewGetOrderLogService(logID int64) *GetOrderLogService {
	return &GetOrderLogService{c: c, logID: logID}
}

func (s *GetOrderLogService) Do(ctx context.Context) (*TradfiOrderLogResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/tradfi/orders/log/"+strconv.FormatInt(s.logID, 10)).WithSign()
	return request.Do[TradfiOrderLogResponse](req)
}

// TradfiOrderLogResponse is the envelope of the order-by-log query.
type TradfiOrderLogResponse struct {
	Timestamp time.Time      `json:"timestamp,format:unixmilli"`
	Data      TradfiOrderLog `json:"data"`
}

// TradfiOrderLog is an order's details resolved by its log id. side is 1=sell,
// 2=buy.
type TradfiOrderLog struct {
	OrderID   int64           `json:"order_id"`
	LogID     int64           `json:"log_id"`
	Symbol    string          `json:"symbol"`
	PriceType string          `json:"price_type"`
	State     int             `json:"state"`
	Side      int             `json:"side"`
	Volume    decimal.Decimal `json:"volume"`
	Price     decimal.Decimal `json:"price"`
}
