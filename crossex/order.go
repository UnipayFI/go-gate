package crossex

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateOrderService -- POST /api/v4/crossex/orders (private)
//
// Places a single cross-exchange order. symbol is the venue-qualified pair
// (e.g. BINANCE_SPOT_BTC_USDT) and side is BUY or SELL.
type CreateOrderService struct {
	c    *CrossexClient
	body map[string]any
}

func (c *CrossexClient) NewCreateOrderService(symbol, side string) *CreateOrderService {
	return &CreateOrderService{c: c, body: map[string]any{
		"symbol": symbol,
		"side":   side,
	}}
}

// SetText attaches a client-defined order id (letters, digits, '-' and '_').
func (s *CreateOrderService) SetText(text string) *CreateOrderService {
	s.body["text"] = text
	return s
}

// SetType sets the order type (LIMIT default, or MARKET).
func (s *CreateOrderService) SetType(orderType string) *CreateOrderService {
	s.body["type"] = orderType
	return s
}

// SetTimeInForce sets the time in force (GTC default, IOC, FOK, POC, RPI).
func (s *CreateOrderService) SetTimeInForce(tif string) *CreateOrderService {
	s.body["time_in_force"] = tif
	return s
}

// SetQty sets the order quantity in the base currency.
func (s *CreateOrderService) SetQty(qty decimal.Decimal) *CreateOrderService {
	s.body["qty"] = qty.String()
	return s
}

// SetPrice sets the limit order price (required for limit orders).
func (s *CreateOrderService) SetPrice(price decimal.Decimal) *CreateOrderService {
	s.body["price"] = price.String()
	return s
}

// SetQuoteQty sets the order quote quantity (required for spot/margin market buys).
func (s *CreateOrderService) SetQuoteQty(quoteQty decimal.Decimal) *CreateOrderService {
	s.body["quote_qty"] = quoteQty.String()
	return s
}

// SetReduceOnly marks the order reduce-only ("true" or "false").
func (s *CreateOrderService) SetReduceOnly(reduceOnly string) *CreateOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetPositionSide sets the position side (NONE, LONG, SHORT).
func (s *CreateOrderService) SetPositionSide(positionSide string) *CreateOrderService {
	s.body["position_side"] = positionSide
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*CrossexOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/orders", s.body).WithSign()
	return request.Do[CrossexOrderResult](req)
}

// CrossexOrderResult is the minimal acknowledgement returned when creating,
// amending, cancelling or closing an order.
type CrossexOrderResult struct {
	OrderID string `json:"order_id"`
	Text    string `json:"text"`
}

// GetOrderService -- GET /api/v4/crossex/orders/{order_id} (private)
//
// Returns a single cross-exchange order by its id or custom text id.
type GetOrderService struct {
	c       *CrossexClient
	orderID string
}

func (c *CrossexClient) NewGetOrderService(orderID string) *GetOrderService {
	return &GetOrderService{c: c, orderID: orderID}
}

func (s *GetOrderService) Do(ctx context.Context) (*CrossexOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/orders/"+s.orderID).WithSign()
	return request.Do[CrossexOrder](req)
}

// AmendOrderService -- PUT /api/v4/crossex/orders/{order_id} (private)
//
// Modifies the quantity and/or price of an open cross-exchange order.
type AmendOrderService struct {
	c       *CrossexClient
	orderID string
	body    map[string]any
}

func (c *CrossexClient) NewAmendOrderService(orderID string) *AmendOrderService {
	return &AmendOrderService{c: c, orderID: orderID, body: map[string]any{}}
}

// SetQty sets the new order quantity.
func (s *AmendOrderService) SetQty(qty decimal.Decimal) *AmendOrderService {
	s.body["qty"] = qty.String()
	return s
}

// SetPrice sets the new order price.
func (s *AmendOrderService) SetPrice(price decimal.Decimal) *AmendOrderService {
	s.body["price"] = price.String()
	return s
}

func (s *AmendOrderService) Do(ctx context.Context) (*CrossexOrderResult, error) {
	req := request.Put(ctx, s.c, "/api/v4/crossex/orders/"+s.orderID, s.body).WithSign()
	return request.Do[CrossexOrderResult](req)
}

// CancelOrderService -- DELETE /api/v4/crossex/orders/{order_id} (private)
//
// Cancels a single cross-exchange order by its id or custom text id.
type CancelOrderService struct {
	c       *CrossexClient
	orderID string
}

func (c *CrossexClient) NewCancelOrderService(orderID string) *CancelOrderService {
	return &CancelOrderService{c: c, orderID: orderID}
}

func (s *CancelOrderService) Do(ctx context.Context) (*CrossexOrderResult, error) {
	req := request.Delete(ctx, s.c, "/api/v4/crossex/orders/"+s.orderID).WithSign()
	return request.Do[CrossexOrderResult](req)
}

// ListOpenOrdersService -- GET /api/v4/crossex/open_orders (private)
//
// Returns all of the account's currently open cross-exchange orders.
type ListOpenOrdersService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListOpenOrdersService() *ListOpenOrdersService {
	return &ListOpenOrdersService{c: c, params: map[string]string{}}
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListOpenOrdersService) SetSymbol(symbol string) *ListOpenOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetExchangeType narrows the result to a single venue.
func (s *ListOpenOrdersService) SetExchangeType(exchangeType string) *ListOpenOrdersService {
	s.params["exchange_type"] = exchangeType
	return s
}

// SetBusinessType narrows the result to a single business type (SPOT/FUTURE/MARGIN).
func (s *ListOpenOrdersService) SetBusinessType(businessType string) *ListOpenOrdersService {
	s.params["business_type"] = businessType
	return s
}

func (s *ListOpenOrdersService) Do(ctx context.Context) ([]CrossexOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/open_orders", s.params).WithSign()
	resp, err := request.Do[[]CrossexOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListHistoryOrdersService -- GET /api/v4/crossex/history_orders (private)
//
// Returns the account's historical cross-exchange orders.
type ListHistoryOrdersService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListHistoryOrdersService() *ListHistoryOrdersService {
	return &ListHistoryOrdersService{c: c, params: map[string]string{}}
}

// SetPage selects the result page.
func (s *ListHistoryOrdersService) SetPage(page int) *ListHistoryOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListHistoryOrdersService) SetLimit(limit int) *ListHistoryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListHistoryOrdersService) SetSymbol(symbol string) *ListHistoryOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListHistoryOrdersService) SetFrom(from time.Time) *ListHistoryOrdersService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision).
func (s *ListHistoryOrdersService) SetTo(to time.Time) *ListHistoryOrdersService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

// SetAttributes filters by order attribute (COMMON, LIQ, REDUCE).
func (s *ListHistoryOrdersService) SetAttributes(attributes string) *ListHistoryOrdersService {
	s.params["attributes"] = attributes
	return s
}

func (s *ListHistoryOrdersService) Do(ctx context.Context) ([]CrossexOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/history_orders", s.params).WithSign()
	resp, err := request.Do[[]CrossexOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexOrder is a cross-exchange order and its live state. create_time and
// update_time are millisecond Unix timestamps. reduce_only is a string ("true"
// or "false"), matching the Gate wire format.
type CrossexOrder struct {
	UserID             string          `json:"user_id"`
	OrderID            string          `json:"order_id"`
	Text               string          `json:"text"`
	State              string          `json:"state"`
	Symbol             string          `json:"symbol"`
	Side               string          `json:"side"`
	Type               string          `json:"type"`
	Attribute          string          `json:"attribute"`
	ExchangeType       string          `json:"exchange_type"`
	BusinessType       string          `json:"business_type"`
	Qty                decimal.Decimal `json:"qty"`
	QuoteQty           decimal.Decimal `json:"quote_qty"`
	Price              decimal.Decimal `json:"price"`
	TimeInForce        string          `json:"time_in_force"`
	ExecutedQty        decimal.Decimal `json:"executed_qty"`
	ExecutedAmount     decimal.Decimal `json:"executed_amount"`
	ExecutedAvgPrice   decimal.Decimal `json:"executed_avg_price"`
	FeeCoin            string          `json:"fee_coin"`
	Fee                decimal.Decimal `json:"fee"`
	ReduceOnly         string          `json:"reduce_only"`
	Leverage           decimal.Decimal `json:"leverage"`
	Reason             string          `json:"reason"`
	LastExecutedQty    decimal.Decimal `json:"last_executed_qty"`
	LastExecutedPrice  decimal.Decimal `json:"last_executed_price"`
	LastExecutedAmount decimal.Decimal `json:"last_executed_amount"`
	PositionSide       string          `json:"position_side"`
	CreateTime         time.Time       `json:"create_time,string,format:unixmilli"`
	UpdateTime         time.Time       `json:"update_time,string,format:unixmilli"`
}

// ListHistoryTradesService -- GET /api/v4/crossex/history_trades (private)
//
// Returns the account's cross-exchange trade-fill history.
type ListHistoryTradesService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewListHistoryTradesService() *ListHistoryTradesService {
	return &ListHistoryTradesService{c: c, params: map[string]string{}}
}

// SetPage selects the result page.
func (s *ListHistoryTradesService) SetPage(page int) *ListHistoryTradesService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single list (max 1000).
func (s *ListHistoryTradesService) SetLimit(limit int) *ListHistoryTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetSymbol narrows the result to a single trading pair.
func (s *ListHistoryTradesService) SetSymbol(symbol string) *ListHistoryTradesService {
	s.params["symbol"] = symbol
	return s
}

// SetFrom sets the start time (millisecond precision).
func (s *ListHistoryTradesService) SetFrom(from time.Time) *ListHistoryTradesService {
	s.params["from"] = strconv.FormatInt(from.UnixMilli(), 10)
	return s
}

// SetTo sets the end time (millisecond precision).
func (s *ListHistoryTradesService) SetTo(to time.Time) *ListHistoryTradesService {
	s.params["to"] = strconv.FormatInt(to.UnixMilli(), 10)
	return s
}

func (s *ListHistoryTradesService) Do(ctx context.Context) ([]CrossexTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/history_trades", s.params).WithSign()
	resp, err := request.Do[[]CrossexTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CancelOrderReq identifies a single order to cancel in a batch request. An
// order may be targeted by order_id or by custom text; when both are set,
// order_id takes precedence.
type CancelOrderReq struct {
	OrderID string `json:"order_id,omitempty"`
	Text    string `json:"text,omitempty"`
}

// CancelBatchOrdersService -- POST /api/v4/crossex/batch_cancel_orders (private)
//
// Cancels multiple cross-exchange orders in one request. Each item targets an
// order by order_id or by custom text; when both are set on an item, order_id
// takes precedence. Each result reports whether the cancel was accepted along
// with any error label and message.
type CancelBatchOrdersService struct {
	c    *CrossexClient
	reqs []CancelOrderReq
}

func (c *CrossexClient) NewCancelBatchOrdersService(reqs ...CancelOrderReq) *CancelBatchOrdersService {
	return &CancelBatchOrdersService{c: c, reqs: reqs}
}

func (s *CancelBatchOrdersService) Do(ctx context.Context) ([]CrossexCancelResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/crossex/batch_cancel_orders").WithSign().SetBody(s.reqs)
	resp, err := request.Do[[]CrossexCancelResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexCancelResult is one element of a batch_cancel_orders response. Accepted
// is a string ("true" or "false"), matching the Gate wire format; on failure
// Label and Message describe the reason.
type CrossexCancelResult struct {
	OrderID  string `json:"order_id"`
	Text     string `json:"text"`
	Accepted string `json:"accepted"`
	Label    string `json:"label"`
	Message  string `json:"message"`
}

// CrossexTrade is one cross-exchange trade fill. create_time is a millisecond
// Unix timestamp.
type CrossexTrade struct {
	UserID        string          `json:"user_id"`
	TransactionID string          `json:"transaction_id"`
	OrderID       string          `json:"order_id"`
	Text          string          `json:"text"`
	Symbol        string          `json:"symbol"`
	ExchangeType  string          `json:"exchange_type"`
	BusinessType  string          `json:"business_type"`
	Side          string          `json:"side"`
	Qty           decimal.Decimal `json:"qty"`
	Price         decimal.Decimal `json:"price"`
	Fee           decimal.Decimal `json:"fee"`
	FeeCoin       string          `json:"fee_coin"`
	FeeRate       decimal.Decimal `json:"fee_rate"`
	MatchRole     string          `json:"match_role"`
	RPnL          decimal.Decimal `json:"rpnl"`
	PositionMode  string          `json:"position_mode"`
	PositionSide  string          `json:"position_side"`
	CreateTime    time.Time       `json:"create_time,string,format:unixmilli"`
}
