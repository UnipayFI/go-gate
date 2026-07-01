package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListPriceTriggeredOrdersService -- GET /api/v4/futures/{settle}/price_orders (private)
//
// Lists the account's price-triggered (auto) futures orders in a given status.
type ListPriceTriggeredOrdersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListPriceTriggeredOrdersService(settle Settle, status OrderStatus) *ListPriceTriggeredOrdersService {
	return &ListPriceTriggeredOrdersService{
		c:      c,
		settle: settle,
		params: map[string]string{"status": string(status)},
	}
}

// SetContract narrows the result to a single futures contract.
func (s *ListPriceTriggeredOrdersService) SetContract(contract string) *ListPriceTriggeredOrdersService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned.
func (s *ListPriceTriggeredOrdersService) SetLimit(limit int) *ListPriceTriggeredOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the given number of records (for pagination).
func (s *ListPriceTriggeredOrdersService) SetOffset(offset int) *ListPriceTriggeredOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListPriceTriggeredOrdersService) Do(ctx context.Context) ([]FuturesPriceTriggeredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/price_orders", s.params).WithSign()
	resp, err := request.Do[[]FuturesPriceTriggeredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreatePriceTriggeredOrderService -- POST /api/v4/futures/{settle}/price_orders (private)
//
// Places a price-triggered (auto) futures order that submits the initial order
// once the trigger condition is met.
type CreatePriceTriggeredOrderService struct {
	c       *FuturesClient
	settle  Settle
	initial map[string]any
	trigger map[string]any
	body    map[string]any
}

// NewCreatePriceTriggeredOrderService takes the required order/trigger fields.
// price is the initial order price (0 for market); triggerPrice and rule define
// the trigger condition (rule 1: fire when price >= triggerPrice, which must be
// above the last price; rule 2: fire when price <= triggerPrice, below the last
// price).
func (c *FuturesClient) NewCreatePriceTriggeredOrderService(settle Settle, contract string, price, triggerPrice decimal.Decimal, rule int) *CreatePriceTriggeredOrderService {
	return &CreatePriceTriggeredOrderService{
		c:      c,
		settle: settle,
		initial: map[string]any{
			"contract": contract,
			"price":    price,
		},
		trigger: map[string]any{
			"price": triggerPrice,
			"rule":  rule,
		},
		body: map[string]any{},
	}
}

// SetSize sets the signed contract size of the initial order (positive long,
// negative short); 0 together with close/auto_size closes a position.
func (s *CreatePriceTriggeredOrderService) SetSize(size int64) *CreatePriceTriggeredOrderService {
	s.initial["size"] = size
	return s
}

// SetTif sets the time-in-force strategy of the initial order.
func (s *CreatePriceTriggeredOrderService) SetTif(tif TimeInForce) *CreatePriceTriggeredOrderService {
	s.initial["tif"] = string(tif)
	return s
}

// SetText sets the initial order's source/user-defined text.
func (s *CreatePriceTriggeredOrderService) SetText(text string) *CreatePriceTriggeredOrderService {
	s.initial["text"] = text
	return s
}

// SetReduceOnly marks the initial order as reduce-only.
func (s *CreatePriceTriggeredOrderService) SetReduceOnly(reduceOnly bool) *CreatePriceTriggeredOrderService {
	s.initial["reduce_only"] = reduceOnly
	return s
}

// SetClose marks the initial order as a full-position close (single-position mode).
func (s *CreatePriceTriggeredOrderService) SetClose(closePosition bool) *CreatePriceTriggeredOrderService {
	s.initial["close"] = closePosition
	return s
}

// SetAutoSize closes one leg of a dual-mode position (size must be 0).
func (s *CreatePriceTriggeredOrderService) SetAutoSize(autoSize AutoSize) *CreatePriceTriggeredOrderService {
	s.initial["auto_size"] = string(autoSize)
	return s
}

// SetStrategyType selects the trigger strategy (0 price trigger, 1 spread trigger).
func (s *CreatePriceTriggeredOrderService) SetStrategyType(strategyType int) *CreatePriceTriggeredOrderService {
	s.trigger["strategy_type"] = strategyType
	return s
}

// SetPriceType selects the reference price (0 last, 1 mark, 2 index).
func (s *CreatePriceTriggeredOrderService) SetPriceType(priceType int) *CreatePriceTriggeredOrderService {
	s.trigger["price_type"] = priceType
	return s
}

// SetExpiration sets the maximum wait time (seconds) before the trigger is cancelled.
func (s *CreatePriceTriggeredOrderService) SetExpiration(expiration int) *CreatePriceTriggeredOrderService {
	s.trigger["expiration"] = expiration
	return s
}

// SetOrderType selects the take-profit/stop-loss order type (e.g. close-long-position).
func (s *CreatePriceTriggeredOrderService) SetOrderType(orderType string) *CreatePriceTriggeredOrderService {
	s.body["order_type"] = orderType
	return s
}

func (s *CreatePriceTriggeredOrderService) Do(ctx context.Context) (*FuturesTriggerOrderResponse, error) {
	s.body["initial"] = s.initial
	s.body["trigger"] = s.trigger
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/price_orders", s.body).WithSign()
	return request.Do[FuturesTriggerOrderResponse](req)
}

// CancelPriceTriggeredOrderListService -- DELETE /api/v4/futures/{settle}/price_orders (private)
//
// Cancels all of the account's price-triggered orders, optionally limited to one
// contract, returning the orders that were cancelled.
type CancelPriceTriggeredOrderListService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewCancelPriceTriggeredOrderListService(settle Settle) *CancelPriceTriggeredOrderListService {
	return &CancelPriceTriggeredOrderListService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract restricts the cancellation to a single futures contract.
func (s *CancelPriceTriggeredOrderListService) SetContract(contract string) *CancelPriceTriggeredOrderListService {
	s.params["contract"] = contract
	return s
}

func (s *CancelPriceTriggeredOrderListService) Do(ctx context.Context) ([]FuturesPriceTriggeredOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/price_orders", s.params).WithSign()
	resp, err := request.Do[[]FuturesPriceTriggeredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetPriceTriggeredOrderService -- GET /api/v4/futures/{settle}/price_orders/{order_id} (private)
//
// Returns the details of a single price-triggered order.
type GetPriceTriggeredOrderService struct {
	c       *FuturesClient
	settle  Settle
	orderID string
}

func (c *FuturesClient) NewGetPriceTriggeredOrderService(settle Settle, orderID string) *GetPriceTriggeredOrderService {
	return &GetPriceTriggeredOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *GetPriceTriggeredOrderService) Do(ctx context.Context) (*FuturesPriceTriggeredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/price_orders/"+s.orderID).WithSign()
	return request.Do[FuturesPriceTriggeredOrder](req)
}

// CancelPriceTriggeredOrderService -- DELETE /api/v4/futures/{settle}/price_orders/{order_id} (private)
//
// Cancels a single price-triggered order, returning its final state.
type CancelPriceTriggeredOrderService struct {
	c       *FuturesClient
	settle  Settle
	orderID string
}

func (c *FuturesClient) NewCancelPriceTriggeredOrderService(settle Settle, orderID string) *CancelPriceTriggeredOrderService {
	return &CancelPriceTriggeredOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *CancelPriceTriggeredOrderService) Do(ctx context.Context) (*FuturesPriceTriggeredOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/price_orders/"+s.orderID).WithSign()
	return request.Do[FuturesPriceTriggeredOrder](req)
}

// FuturesTriggerOrderResponse is the id returned when a price-triggered order is
// successfully created.
type FuturesTriggerOrderResponse struct {
	ID int64 `json:"id"`
}

// FuturesPriceTriggeredOrder is a price-triggered (auto) futures order: the
// initial order to submit and the trigger condition that submits it, plus its
// lifecycle metadata.
type FuturesPriceTriggeredOrder struct {
	Initial    FuturesInitialOrder `json:"initial"`
	Trigger    FuturesPriceTrigger `json:"trigger"`
	ID         int64               `json:"id"`
	User       int64               `json:"user"`
	CreateTime time.Time           `json:"create_time,format:unix"`
	FinishTime time.Time           `json:"finish_time,format:unix"`
	TradeID    int64               `json:"trade_id"`
	Status     OrderStatus         `json:"status"`
	FinishAs   FinishAs            `json:"finish_as"`
	Reason     string              `json:"reason"`
	OrderType  string              `json:"order_type"`
	MeOrderID  int64               `json:"me_order_id"`
}

// FuturesInitialOrder is the order that a price-triggered order submits once its
// trigger fires.
type FuturesInitialOrder struct {
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
	Close        bool            `json:"close"`
	TimeInForce  TimeInForce     `json:"tif"`
	Text         string          `json:"text"`
	ReduceOnly   bool            `json:"reduce_only"`
	AutoSize     AutoSize        `json:"auto_size"`
	IsReduceOnly bool            `json:"is_reduce_only"`
	IsClose      bool            `json:"is_close"`
}

// FuturesPriceTrigger is the price condition that fires a price-triggered order.
type FuturesPriceTrigger struct {
	StrategyType int             `json:"strategy_type"`
	PriceType    int             `json:"price_type"`
	Price        decimal.Decimal `json:"price"`
	Rule         int             `json:"rule"`
	Expiration   int             `json:"expiration"`
}
