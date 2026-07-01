package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// CreateOrderService -- POST /api/v4/spot/orders (private)
//
// Places a spot, margin, cross-margin or unified-account order. Only
// currency_pair, side and amount are required; a limit order also needs a price.
type CreateOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCreateOrderService(currencyPair string, side Side, amount decimal.Decimal) *CreateOrderService {
	return &CreateOrderService{c: c, body: map[string]any{
		"currency_pair": currencyPair,
		"side":          string(side),
		"amount":        amount.String(),
	}}
}

// SetType selects the order type (limit or market). Defaults to limit server-side.
func (s *CreateOrderService) SetType(orderType OrderType) *CreateOrderService {
	s.body["type"] = string(orderType)
	return s
}

// SetPrice sets the limit price. Required when type is limit.
func (s *CreateOrderService) SetPrice(price decimal.Decimal) *CreateOrderService {
	s.body["price"] = price.String()
	return s
}

// SetAccount selects which balance the order draws on.
func (s *CreateOrderService) SetAccount(account Account) *CreateOrderService {
	s.body["account"] = string(account)
	return s
}

// SetTimeInForce sets how long the order stays active.
func (s *CreateOrderService) SetTimeInForce(tif TimeInForce) *CreateOrderService {
	s.body["time_in_force"] = string(tif)
	return s
}

// SetIceberg sets the displayed quantity for an iceberg order.
func (s *CreateOrderService) SetIceberg(iceberg decimal.Decimal) *CreateOrderService {
	s.body["iceberg"] = iceberg.String()
	return s
}

// SetAutoBorrow enables automatic borrowing of any insufficient amount (margin).
func (s *CreateOrderService) SetAutoBorrow(autoBorrow bool) *CreateOrderService {
	s.body["auto_borrow"] = autoBorrow
	return s
}

// SetAutoRepay enables automatic repayment of the cross-margin loan when the
// order ends.
func (s *CreateOrderService) SetAutoRepay(autoRepay bool) *CreateOrderService {
	s.body["auto_repay"] = autoRepay
	return s
}

// SetStpAct sets the self-trade-prevention strategy.
func (s *CreateOrderService) SetStpAct(stpAct StpAct) *CreateOrderService {
	s.body["stp_act"] = string(stpAct)
	return s
}

// SetText sets a user-defined order label. Must be prefixed with "t-".
func (s *CreateOrderService) SetText(text string) *CreateOrderService {
	s.body["text"] = text
	return s
}

// SetActionMode selects the response detail level (ACK, RESULT or FULL).
func (s *CreateOrderService) SetActionMode(actionMode string) *CreateOrderService {
	s.body["action_mode"] = actionMode
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/orders", s.body).WithSign()
	return request.Do[Order](req)
}

// CreateBatchOrdersService -- POST /api/v4/spot/batch_orders (private)
//
// Places multiple orders in a single request (max 10). Each order is described
// with the same builder used by NewCreateOrderService.
type CreateBatchOrdersService struct {
	c      *SpotClient
	orders []*CreateOrderService
}

func (c *SpotClient) NewCreateBatchOrdersService(orders ...*CreateOrderService) *CreateBatchOrdersService {
	return &CreateBatchOrdersService{c: c, orders: orders}
}

func (s *CreateBatchOrdersService) Do(ctx context.Context) ([]BatchOrderResult, error) {
	bodies := make([]map[string]any, len(s.orders))
	for i, o := range s.orders {
		bodies[i] = o.body
	}
	req := request.Post(ctx, s.c, "/api/v4/spot/batch_orders").WithSign().SetBody(bodies)
	resp, err := request.Do[[]BatchOrderResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListAllOpenOrdersService -- GET /api/v4/spot/open_orders (private)
//
// Returns the open orders of every trading pair that has any, grouped by pair.
type ListAllOpenOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListAllOpenOrdersService() *ListAllOpenOrdersService {
	return &ListAllOpenOrdersService{c: c, params: map[string]string{}}
}

// SetPage sets the page number.
func (s *ListAllOpenOrdersService) SetPage(page int) *ListAllOpenOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of open orders returned per currency pair.
func (s *ListAllOpenOrdersService) SetLimit(limit int) *ListAllOpenOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetAccount restricts the query to a single account type.
func (s *ListAllOpenOrdersService) SetAccount(account Account) *ListAllOpenOrdersService {
	s.params["account"] = string(account)
	return s
}

func (s *ListAllOpenOrdersService) Do(ctx context.Context) ([]OpenOrders, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/open_orders", s.params).WithSign()
	resp, err := request.Do[[]OpenOrders](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListOrdersService -- GET /api/v4/spot/orders (private)
//
// Lists orders for a currency pair filtered by status ("open" or "finished").
type ListOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListOrdersService(currencyPair, status string) *ListOrdersService {
	return &ListOrdersService{c: c, params: map[string]string{
		"currency_pair": currencyPair,
		"status":        status,
	}}
}

// SetPage sets the page number.
func (s *ListOrdersService) SetPage(page int) *ListOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (max 100 when status is open).
func (s *ListOrdersService) SetLimit(limit int) *ListOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetAccount restricts the query to a single account type.
func (s *ListOrdersService) SetAccount(account Account) *ListOrdersService {
	s.params["account"] = string(account)
	return s
}

// SetFrom sets the start time of the query window (finished orders only).
func (s *ListOrdersService) SetFrom(from time.Time) *ListOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window (finished orders only).
func (s *ListOrdersService) SetTo(to time.Time) *ListOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetSide filters to one side (finished orders only).
func (s *ListOrdersService) SetSide(side Side) *ListOrdersService {
	s.params["side"] = string(side)
	return s
}

func (s *ListOrdersService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/orders", s.params).WithSign()
	resp, err := request.Do[[]Order](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CancelOrdersService -- DELETE /api/v4/spot/orders (private)
//
// Cancels all open orders, optionally scoped to a currency pair, side or account.
type CancelOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelOrdersService() *CancelOrdersService {
	return &CancelOrdersService{c: c, params: map[string]string{}}
}

// SetCurrencyPair limits the cancellation to a single trading pair.
func (s *CancelOrdersService) SetCurrencyPair(currencyPair string) *CancelOrdersService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetSide limits the cancellation to one side.
func (s *CancelOrdersService) SetSide(side Side) *CancelOrdersService {
	s.params["side"] = string(side)
	return s
}

// SetAccount limits the cancellation to a single account type.
func (s *CancelOrdersService) SetAccount(account Account) *CancelOrdersService {
	s.params["account"] = string(account)
	return s
}

// SetActionMode selects the response detail level (ACK, RESULT or FULL).
func (s *CancelOrdersService) SetActionMode(actionMode string) *CancelOrdersService {
	s.params["action_mode"] = actionMode
	return s
}

func (s *CancelOrdersService) Do(ctx context.Context) ([]Order, error) {
	req := request.Delete(ctx, s.c, "/api/v4/spot/orders", s.params).WithSign()
	resp, err := request.Do[[]Order](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CancelOrderReq identifies a single order to cancel in a batch request.
type CancelOrderReq struct {
	CurrencyPair string  `json:"currency_pair"`
	ID           string  `json:"id"`
	Account      Account `json:"account,omitempty"`
	ActionMode   string  `json:"action_mode,omitempty"`
}

// CancelBatchOrdersService -- POST /api/v4/spot/cancel_batch_orders (private)
//
// Cancels a list of orders by ID (max 20), spanning multiple currency pairs.
type CancelBatchOrdersService struct {
	c    *SpotClient
	reqs []CancelOrderReq
}

func (c *SpotClient) NewCancelBatchOrdersService(reqs ...CancelOrderReq) *CancelBatchOrdersService {
	return &CancelBatchOrdersService{c: c, reqs: reqs}
}

func (s *CancelBatchOrdersService) Do(ctx context.Context) ([]CancelOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/cancel_batch_orders").WithSign().SetBody(s.reqs)
	resp, err := request.Do[[]CancelOrderResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOrderService -- GET /api/v4/spot/orders/{order_id} (private)
//
// Returns the details of a single order. currency_pair is required for pending
// orders and optional for finished ones.
type GetOrderService struct {
	c       *SpotClient
	orderID string
	params  map[string]string
}

func (c *SpotClient) NewGetOrderService(orderID, currencyPair string) *GetOrderService {
	return &GetOrderService{c: c, orderID: orderID, params: map[string]string{
		"currency_pair": currencyPair,
	}}
}

// SetAccount specifies the account to query.
func (s *GetOrderService) SetAccount(account Account) *GetOrderService {
	s.params["account"] = string(account)
	return s
}

func (s *GetOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/orders/"+s.orderID, s.params).WithSign()
	return request.Do[Order](req)
}

// AmendOrderService -- PATCH /api/v4/spot/orders/{order_id} (private)
//
// Amends a live order's price and/or amount. currency_pair must be supplied
// either here or on the order.
type AmendOrderService struct {
	c       *SpotClient
	orderID string
	body    map[string]any
}

func (c *SpotClient) NewAmendOrderService(orderID string) *AmendOrderService {
	return &AmendOrderService{c: c, orderID: orderID, body: map[string]any{}}
}

// SetCurrencyPair sets the trading pair of the order being amended.
func (s *AmendOrderService) SetCurrencyPair(currencyPair string) *AmendOrderService {
	s.body["currency_pair"] = currencyPair
	return s
}

// SetAmount sets the new trading quantity. Only one of amount or price is allowed.
func (s *AmendOrderService) SetAmount(amount decimal.Decimal) *AmendOrderService {
	s.body["amount"] = amount.String()
	return s
}

// SetPrice sets the new trading price. Only one of amount or price is allowed.
func (s *AmendOrderService) SetPrice(price decimal.Decimal) *AmendOrderService {
	s.body["price"] = price.String()
	return s
}

// SetAmendText attaches custom info to the amendment.
func (s *AmendOrderService) SetAmendText(amendText string) *AmendOrderService {
	s.body["amend_text"] = amendText
	return s
}

// SetAccount specifies the account of the order being amended.
func (s *AmendOrderService) SetAccount(account Account) *AmendOrderService {
	s.body["account"] = string(account)
	return s
}

// SetActionMode selects the response detail level (ACK, RESULT or FULL).
func (s *AmendOrderService) SetActionMode(actionMode string) *AmendOrderService {
	s.body["action_mode"] = actionMode
	return s
}

func (s *AmendOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Patch(ctx, s.c, "/api/v4/spot/orders/"+s.orderID, s.body).WithSign()
	return request.Do[Order](req)
}

// CancelOrderService -- DELETE /api/v4/spot/orders/{order_id} (private)
//
// Cancels a single order by ID.
type CancelOrderService struct {
	c       *SpotClient
	orderID string
	params  map[string]string
}

func (c *SpotClient) NewCancelOrderService(orderID, currencyPair string) *CancelOrderService {
	return &CancelOrderService{c: c, orderID: orderID, params: map[string]string{
		"currency_pair": currencyPair,
	}}
}

// SetAccount specifies the account of the order being cancelled.
func (s *CancelOrderService) SetAccount(account Account) *CancelOrderService {
	s.params["account"] = string(account)
	return s
}

// SetActionMode selects the response detail level (ACK, RESULT or FULL).
func (s *CancelOrderService) SetActionMode(actionMode string) *CancelOrderService {
	s.params["action_mode"] = actionMode
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Delete(ctx, s.c, "/api/v4/spot/orders/"+s.orderID, s.params).WithSign()
	return request.Do[Order](req)
}

// CountdownCancelAllSpotService -- POST /api/v4/spot/countdown_cancel_all (private)
//
// Arms a dead-man's-switch: if not renewed within timeout seconds, all pending
// spot orders (optionally scoped to a pair) are cancelled. A timeout of 0
// disarms it.
type CountdownCancelAllSpotService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCountdownCancelAllSpotService(timeout int) *CountdownCancelAllSpotService {
	return &CountdownCancelAllSpotService{c: c, body: map[string]any{
		"timeout": timeout,
	}}
}

// SetCurrencyPair scopes the countdown to a single trading pair.
func (s *CountdownCancelAllSpotService) SetCurrencyPair(currencyPair string) *CountdownCancelAllSpotService {
	s.body["currency_pair"] = currencyPair
	return s
}

func (s *CountdownCancelAllSpotService) Do(ctx context.Context) (*CountdownStatus, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/countdown_cancel_all", s.body).WithSign()
	return request.Do[CountdownStatus](req)
}

// BatchAmendItem describes one order amendment inside a batch request. Build it
// with NewBatchAmendItem and chain setters for the fields to change.
type BatchAmendItem struct {
	body map[string]any
}

// NewBatchAmendItem starts a batch amendment for one order.
func NewBatchAmendItem(orderID, currencyPair string) *BatchAmendItem {
	return &BatchAmendItem{body: map[string]any{
		"order_id":      orderID,
		"currency_pair": currencyPair,
	}}
}

// SetAmount sets the new trading quantity. Only one of amount or price is allowed.
func (i *BatchAmendItem) SetAmount(amount decimal.Decimal) *BatchAmendItem {
	i.body["amount"] = amount.String()
	return i
}

// SetPrice sets the new trading price. Only one of amount or price is allowed.
func (i *BatchAmendItem) SetPrice(price decimal.Decimal) *BatchAmendItem {
	i.body["price"] = price.String()
	return i
}

// SetAccount specifies the account of the order being amended.
func (i *BatchAmendItem) SetAccount(account Account) *BatchAmendItem {
	i.body["account"] = string(account)
	return i
}

// SetAmendText attaches custom info to the amendment.
func (i *BatchAmendItem) SetAmendText(amendText string) *BatchAmendItem {
	i.body["amend_text"] = amendText
	return i
}

// SetActionMode selects the response detail level (ACK, RESULT or FULL).
func (i *BatchAmendItem) SetActionMode(actionMode string) *BatchAmendItem {
	i.body["action_mode"] = actionMode
	return i
}

// AmendBatchOrdersService -- POST /api/v4/spot/amend_batch_orders (private)
//
// Amends up to 5 orders in a single request. Results are returned in the same
// order as the request items.
type AmendBatchOrdersService struct {
	c     *SpotClient
	items []*BatchAmendItem
}

func (c *SpotClient) NewAmendBatchOrdersService(items ...*BatchAmendItem) *AmendBatchOrdersService {
	return &AmendBatchOrdersService{c: c, items: items}
}

func (s *AmendBatchOrdersService) Do(ctx context.Context) ([]BatchOrderResult, error) {
	bodies := make([]map[string]any, len(s.items))
	for i, item := range s.items {
		bodies[i] = item.body
	}
	req := request.Post(ctx, s.c, "/api/v4/spot/amend_batch_orders").WithSign().SetBody(bodies)
	resp, err := request.Do[[]BatchOrderResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateCrossLiquidateOrderService -- POST /api/v4/spot/cross_liquidate_orders (private)
//
// Places a buy order to close a position in a currency that has been disabled
// for cross-margin trading. currency_pair, amount and price are required.
type CreateCrossLiquidateOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCreateCrossLiquidateOrderService(currencyPair string, amount, price decimal.Decimal) *CreateCrossLiquidateOrderService {
	return &CreateCrossLiquidateOrderService{c: c, body: map[string]any{
		"currency_pair": currencyPair,
		"amount":        amount.String(),
		"price":         price.String(),
	}}
}

// SetText sets a user-defined order label. Must be prefixed with "t-".
func (s *CreateCrossLiquidateOrderService) SetText(text string) *CreateCrossLiquidateOrderService {
	s.body["text"] = text
	return s
}

// SetActionMode selects the response detail level (ACK, RESULT or FULL).
func (s *CreateCrossLiquidateOrderService) SetActionMode(actionMode string) *CreateCrossLiquidateOrderService {
	s.body["action_mode"] = actionMode
	return s
}

func (s *CreateCrossLiquidateOrderService) Do(ctx context.Context) (*CrossLiquidateOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/cross_liquidate_orders", s.body).WithSign()
	return request.Do[CrossLiquidateOrder](req)
}

// Order is a spot/margin order and its lifecycle state.
type Order struct {
	ID                 string          `json:"id"`
	Text               string          `json:"text"`
	AmendText          string          `json:"amend_text"`
	CreateTime         time.Time       `json:"create_time,string,format:unix"`
	UpdateTime         time.Time       `json:"update_time,string,format:unix"`
	CreateTimeMs       time.Time       `json:"create_time_ms,format:unixmilli"`
	UpdateTimeMs       time.Time       `json:"update_time_ms,format:unixmilli"`
	Status             OrderStatus     `json:"status"`
	CurrencyPair       string          `json:"currency_pair"`
	Type               OrderType       `json:"type"`
	Account            Account         `json:"account"`
	Side               Side            `json:"side"`
	Amount             decimal.Decimal `json:"amount"`
	Price              decimal.Decimal `json:"price"`
	TimeInForce        TimeInForce     `json:"time_in_force"`
	Iceberg            decimal.Decimal `json:"iceberg"`
	AutoBorrow         bool            `json:"auto_borrow"`
	AutoRepay          bool            `json:"auto_repay"`
	Left               decimal.Decimal `json:"left"`
	FilledAmount       decimal.Decimal `json:"filled_amount"`
	FillPrice          decimal.Decimal `json:"fill_price"`
	FilledTotal        decimal.Decimal `json:"filled_total"`
	AvgDealPrice       decimal.Decimal `json:"avg_deal_price"`
	Fee                decimal.Decimal `json:"fee"`
	FeeCurrency        string          `json:"fee_currency"`
	PointFee           decimal.Decimal `json:"point_fee"`
	GTFee              decimal.Decimal `json:"gt_fee"`
	GTMakerFee         decimal.Decimal `json:"gt_maker_fee"`
	GTTakerFee         decimal.Decimal `json:"gt_taker_fee"`
	GTDiscount         bool            `json:"gt_discount"`
	RebatedFee         decimal.Decimal `json:"rebated_fee"`
	RebatedFeeCurrency string          `json:"rebated_fee_currency"`
	StpID              int             `json:"stp_id"`
	StpAct             StpAct          `json:"stp_act"`
	FinishAs           string          `json:"finish_as"`
}

// OpenOrders groups the open orders of a single trading pair.
type OpenOrders struct {
	CurrencyPair string  `json:"currency_pair"`
	Total        int     `json:"total"`
	Orders       []Order `json:"orders"`
}

// CancelOrderResult is the per-order outcome of a batch cancellation.
type CancelOrderResult struct {
	CurrencyPair string  `json:"currency_pair"`
	ID           string  `json:"id"`
	Text         string  `json:"text"`
	Succeeded    bool    `json:"succeeded"`
	Label        string  `json:"label"`
	Message      string  `json:"message"`
	Account      Account `json:"account"`
}

// BatchOrderResult is one order's result inside a batch create/amend response.
// It carries the order fields plus the per-order success flag and error detail.
type BatchOrderResult struct {
	OrderID            string          `json:"order_id"`
	AmendText          string          `json:"amend_text"`
	Text               string          `json:"text"`
	Succeeded          bool            `json:"succeeded"`
	Label              string          `json:"label"`
	Message            string          `json:"message"`
	ID                 string          `json:"id"`
	CreateTime         time.Time       `json:"create_time,string,format:unix"`
	UpdateTime         time.Time       `json:"update_time,string,format:unix"`
	CreateTimeMs       time.Time       `json:"create_time_ms,format:unixmilli"`
	UpdateTimeMs       time.Time       `json:"update_time_ms,format:unixmilli"`
	Status             OrderStatus     `json:"status"`
	CurrencyPair       string          `json:"currency_pair"`
	Type               OrderType       `json:"type"`
	Account            Account         `json:"account"`
	Side               Side            `json:"side"`
	Amount             decimal.Decimal `json:"amount"`
	Price              decimal.Decimal `json:"price"`
	TimeInForce        TimeInForce     `json:"time_in_force"`
	Iceberg            decimal.Decimal `json:"iceberg"`
	AutoBorrow         bool            `json:"auto_borrow"`
	AutoRepay          bool            `json:"auto_repay"`
	Left               decimal.Decimal `json:"left"`
	FilledAmount       decimal.Decimal `json:"filled_amount"`
	FillPrice          decimal.Decimal `json:"fill_price"`
	FilledTotal        decimal.Decimal `json:"filled_total"`
	AvgDealPrice       decimal.Decimal `json:"avg_deal_price"`
	Fee                decimal.Decimal `json:"fee"`
	FeeCurrency        string          `json:"fee_currency"`
	PointFee           decimal.Decimal `json:"point_fee"`
	GTFee              decimal.Decimal `json:"gt_fee"`
	GTDiscount         bool            `json:"gt_discount"`
	RebatedFee         decimal.Decimal `json:"rebated_fee"`
	RebatedFeeCurrency string          `json:"rebated_fee_currency"`
	StpID              int             `json:"stp_id"`
	StpAct             StpAct          `json:"stp_act"`
	FinishAs           string          `json:"finish_as"`
}

// CountdownStatus reports when the armed countdown will fire. triggerTime is a
// millisecond epoch.
type CountdownStatus struct {
	TriggerTime time.Time `json:"triggerTime,format:unixmilli"`
}

// CrossLiquidateOrder is the spot order produced by a cross-margin liquidation.
type CrossLiquidateOrder struct {
	ID                 string          `json:"id"`
	Text               string          `json:"text"`
	AmendText          string          `json:"amend_text"`
	CreateTime         time.Time       `json:"create_time,string,format:unix"`
	UpdateTime         time.Time       `json:"update_time,string,format:unix"`
	CreateTimeMs       time.Time       `json:"create_time_ms,format:unixmilli"`
	UpdateTimeMs       time.Time       `json:"update_time_ms,format:unixmilli"`
	Status             OrderStatus     `json:"status"`
	CurrencyPair       string          `json:"currency_pair"`
	Type               OrderType       `json:"type"`
	Account            Account         `json:"account"`
	Side               Side            `json:"side"`
	Amount             decimal.Decimal `json:"amount"`
	Price              decimal.Decimal `json:"price"`
	TimeInForce        TimeInForce     `json:"time_in_force"`
	Iceberg            decimal.Decimal `json:"iceberg"`
	AutoBorrow         bool            `json:"auto_borrow"`
	AutoRepay          bool            `json:"auto_repay"`
	Left               decimal.Decimal `json:"left"`
	FilledAmount       decimal.Decimal `json:"filled_amount"`
	FillPrice          decimal.Decimal `json:"fill_price"`
	FilledTotal        decimal.Decimal `json:"filled_total"`
	AvgDealPrice       decimal.Decimal `json:"avg_deal_price"`
	Fee                decimal.Decimal `json:"fee"`
	FeeCurrency        string          `json:"fee_currency"`
	PointFee           decimal.Decimal `json:"point_fee"`
	GTFee              decimal.Decimal `json:"gt_fee"`
	GTMakerFee         decimal.Decimal `json:"gt_maker_fee"`
	GTTakerFee         decimal.Decimal `json:"gt_taker_fee"`
	GTDiscount         bool            `json:"gt_discount"`
	RebatedFee         decimal.Decimal `json:"rebated_fee"`
	RebatedFeeCurrency string          `json:"rebated_fee_currency"`
	StpID              int             `json:"stp_id"`
	StpAct             StpAct          `json:"stp_act"`
	FinishAs           string          `json:"finish_as"`
}
