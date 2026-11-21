package stockoptions

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListOrdersService -- GET /api/v4/stock_options/orders (private)
//
// Lists the current user's unfinished option orders (status 1=new, 2=partially
// filled), sorted by time_setup descending.
type ListOrdersService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c, params: map[string]string{}}
}

// SetUnderlying narrows the result to a single option underlying (e.g. "AAPL").
func (s *ListOrdersService) SetUnderlying(underlying string) *ListOrdersService {
	s.params["underlying"] = underlying
	return s
}

// SetSymbol narrows the result to a full option contract symbol.
func (s *ListOrdersService) SetSymbol(symbol string) *ListOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetSide narrows the result to one order side (1=sell, 2=buy).
func (s *ListOrdersService) SetSide(side int) *ListOrdersService {
	s.params["side"] = strconv.Itoa(side)
	return s
}

func (s *ListOrdersService) Do(ctx context.Context) (*StockOptionOrdersResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/orders", s.params).WithSign()
	return request.Do[StockOptionOrdersResponse](req)
}

// StockOptionOrdersResponse is the envelope of the open-orders query.
type StockOptionOrdersResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []StockOptionOrder `json:"list"`
	} `json:"data"`
}

// StockOptionOrder is a single active (unfinished) stock-option order.
// side is 1=sell, 2=buy; order_type is "market" or "limit"; status is 1=new,
// 2=partially filled; time_in_force is "day" or "gtc"; time_setup is Unix
// seconds.
type StockOptionOrder struct {
	OrderID        string          `json:"order_id"`
	Symbol         string          `json:"symbol"`
	SymbolDisplay  string          `json:"symbol_display"`
	Underlying     string          `json:"underlying"`
	Side           int             `json:"side"`
	OrderType      string          `json:"order_type"`
	Status         int             `json:"status"`
	StatusDesc     string          `json:"status_desc"`
	TimeInForce    string          `json:"time_in_force"`
	Quantity       decimal.Decimal `json:"quantity"`
	FilledQuantity decimal.Decimal `json:"filled_quantity"`
	Price          decimal.Decimal `json:"price"`
	TimeSetup      time.Time       `json:"time_setup,format:unix"`
}

// CreateOrderService -- POST /api/v4/stock_options/orders (private)
//
// Places a single option order. side is 1=sell, 2=buy; orderType is "market" or
// "limit"; timeInForce is "day" or "gtc" (market orders only support "day").
// Limit orders require SetPrice.
type CreateOrderService struct {
	c    *StockOptionsClient
	body map[string]any
}

func (c *StockOptionsClient) NewCreateOrderService(symbol string, side int, orderType, timeInForce string, quantity decimal.Decimal) *CreateOrderService {
	return &CreateOrderService{c: c, body: map[string]any{
		"symbol":        symbol,
		"side":          side,
		"order_type":    orderType,
		"time_in_force": timeInForce,
		"quantity":      quantity.String(),
	}}
}

// SetPrice sets the limit order price (required for limit orders, ignored for
// market orders).
func (s *CreateOrderService) SetPrice(price decimal.Decimal) *CreateOrderService {
	s.body["price"] = price.String()
	return s
}

// SetClientOrderID sets a custom order identifier (up to 30 UTF-8 bytes).
func (s *CreateOrderService) SetClientOrderID(clientOrderID string) *CreateOrderService {
	s.body["client_order_id"] = clientOrderID
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*StockOptionCreateOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/stock_options/orders", s.body).WithSign()
	return request.Do[StockOptionCreateOrderResponse](req)
}

// StockOptionCreateOrderResponse is the envelope of the create-order request.
type StockOptionCreateOrderResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CancelAllOrdersService -- DELETE /api/v4/stock_options/orders (private)
//
// Cancels all or a filtered subset of the current user's option limit orders.
type CancelAllOrdersService struct {
	c    *StockOptionsClient
	body map[string]any
}

func (c *StockOptionsClient) NewCancelAllOrdersService() *CancelAllOrdersService {
	return &CancelAllOrdersService{c: c, body: map[string]any{}}
}

// SetUnderlying narrows cancellation to a single option underlying.
func (s *CancelAllOrdersService) SetUnderlying(underlying string) *CancelAllOrdersService {
	s.body["underlying"] = underlying
	return s
}

// SetSymbol narrows cancellation to a single option contract symbol.
func (s *CancelAllOrdersService) SetSymbol(symbol string) *CancelAllOrdersService {
	s.body["symbol"] = symbol
	return s
}

// SetSide narrows cancellation to one order side (1=sell, 2=buy).
func (s *CancelAllOrdersService) SetSide(side int) *CancelAllOrdersService {
	s.body["side"] = side
	return s
}

func (s *CancelAllOrdersService) Do(ctx context.Context) (*StockOptionCancelAllOrdersResponse, error) {
	req := request.Delete(ctx, s.c, "/api/v4/stock_options/orders").WithSign()
	if len(s.body) > 0 {
		req.SetBody(s.body)
	}
	return request.Do[StockOptionCancelAllOrdersResponse](req)
}

// StockOptionCancelAllOrdersResponse is the envelope of the bulk-cancel request.
type StockOptionCancelAllOrdersResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Success bool `json:"success"`
	} `json:"data"`
}

// GetOrderService -- GET /api/v4/stock_options/orders/{order_id} (private)
//
// Queries a single option order (active or terminal) by order ID.
type GetOrderService struct {
	c       *StockOptionsClient
	orderID string
}

func (c *StockOptionsClient) NewGetOrderService(orderID string) *GetOrderService {
	return &GetOrderService{c: c, orderID: orderID}
}

func (s *GetOrderService) Do(ctx context.Context) (*StockOptionOrderDetailResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/orders/"+s.orderID).WithSign()
	return request.Do[StockOptionOrderDetailResponse](req)
}

// StockOptionOrderDetailResponse is the envelope of the single-order query.
type StockOptionOrderDetailResponse struct {
	Label     string                 `json:"label"`
	Timestamp time.Time              `json:"timestamp,format:unixmilli"`
	Data      StockOptionOrderDetail `json:"data"`
}

// StockOptionOrderDetail is the full representation of an option order, used by
// both GetOrderService and ListOrderHistoryService. side is 1=sell, 2=buy;
// order_type is "market" or "limit"; status is 1=new, 2=partially filled,
// 3=filled, 4=cancelled, 5=rejected; time_in_force is "day" or "gtc";
// time_setup and time_done are Unix seconds.
type StockOptionOrderDetail struct {
	OrderID        string          `json:"order_id"`
	Symbol         string          `json:"symbol"`
	SymbolDisplay  string          `json:"symbol_display"`
	Underlying     string          `json:"underlying"`
	Side           int             `json:"side"`
	OrderType      string          `json:"order_type"`
	Status         int             `json:"status"`
	StatusDesc     string          `json:"status_desc"`
	TimeInForce    string          `json:"time_in_force"`
	Quantity       decimal.Decimal `json:"quantity"`
	FilledQuantity decimal.Decimal `json:"filled_quantity"`
	Price          decimal.Decimal `json:"price"`
	AvgFillPrice   decimal.Decimal `json:"avg_fill_price"`
	Fee            decimal.Decimal `json:"fee"`
	TimeSetup      time.Time       `json:"time_setup,format:unix"`
	TimeDone       time.Time       `json:"time_done,format:unix"`
}

// AmendOrderService -- PUT /api/v4/stock_options/orders/{order_id} (private)
//
// Updates the price and new total quantity of an active option limit order.
// quantity is the new total order quantity (not remaining unfilled) and must be
// at least filled_quantity.
type AmendOrderService struct {
	c       *StockOptionsClient
	orderID string
	body    map[string]any
}

func (c *StockOptionsClient) NewAmendOrderService(orderID string, price, quantity decimal.Decimal) *AmendOrderService {
	return &AmendOrderService{c: c, orderID: orderID, body: map[string]any{
		"price":    price.String(),
		"quantity": quantity.String(),
	}}
}

func (s *AmendOrderService) Do(ctx context.Context) (*StockOptionAmendOrderResponse, error) {
	req := request.Put(ctx, s.c, "/api/v4/stock_options/orders/"+s.orderID, s.body).WithSign()
	return request.Do[StockOptionAmendOrderResponse](req)
}

// StockOptionAmendOrderResponse is the envelope of the amend-order request.
type StockOptionAmendOrderResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CancelOrderService -- DELETE /api/v4/stock_options/orders/{order_id} (private)
//
// Cancels a single option limit order by ID.
type CancelOrderService struct {
	c       *StockOptionsClient
	orderID string
}

func (c *StockOptionsClient) NewCancelOrderService(orderID string) *CancelOrderService {
	return &CancelOrderService{c: c, orderID: orderID}
}

func (s *CancelOrderService) Do(ctx context.Context) (*StockOptionCancelOrderResponse, error) {
	req := request.Delete(ctx, s.c, "/api/v4/stock_options/orders/"+s.orderID).WithSign()
	return request.Do[StockOptionCancelOrderResponse](req)
}

// StockOptionCancelOrderResponse is the envelope of the cancel-order request.
type StockOptionCancelOrderResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		ID string `json:"id"`
	} `json:"data"`
}

// ListOrderHistoryService -- GET /api/v4/stock_options/orders/history (private)
//
// Lists the current user's finished option orders (status 3=filled, 4=cancelled,
// 5=rejected), paginated. Default and maximum page_size are both 100.
type ListOrderHistoryService struct {
	c      *StockOptionsClient
	params map[string]string
}

func (c *StockOptionsClient) NewListOrderHistoryService() *ListOrderHistoryService {
	return &ListOrderHistoryService{c: c, params: map[string]string{}}
}

// SetUnderlying narrows the result to a single option underlying.
func (s *ListOrderHistoryService) SetUnderlying(underlying string) *ListOrderHistoryService {
	s.params["underlying"] = underlying
	return s
}

// SetSymbol narrows the result to a full option contract symbol.
func (s *ListOrderHistoryService) SetSymbol(symbol string) *ListOrderHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetBeginTime bounds the result to orders created at or after beginTime.
func (s *ListOrderHistoryService) SetBeginTime(beginTime time.Time) *ListOrderHistoryService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to orders created at or before endTime.
func (s *ListOrderHistoryService) SetEndTime(endTime time.Time) *ListOrderHistoryService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetSide narrows the result to one order side (1=sell, 2=buy).
func (s *ListOrderHistoryService) SetSide(side int) *ListOrderHistoryService {
	s.params["side"] = strconv.Itoa(side)
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListOrderHistoryService) SetPage(page int) *ListOrderHistoryService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the items per page (defaults to and maxes at 100).
func (s *ListOrderHistoryService) SetPageSize(pageSize int) *ListOrderHistoryService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListOrderHistoryService) Do(ctx context.Context) (*StockOptionOrderHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock_options/orders/history", s.params).WithSign()
	return request.Do[StockOptionOrderHistoryResponse](req)
}

// StockOptionOrderHistoryResponse is the envelope of the historical-orders query.
type StockOptionOrderHistoryResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64                    `json:"total"`
		TotalPage int                      `json:"total_page"`
		List      []StockOptionOrderDetail `json:"list"`
	} `json:"data"`
}
