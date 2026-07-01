package delivery

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListPriceTriggeredDeliveryOrdersService -- GET /api/v4/delivery/{settle}/price_orders (private)
//
// Returns the authenticated account's price-triggered (auto) delivery orders in
// the given status ("open" or "finished").
type ListPriceTriggeredDeliveryOrdersService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListPriceTriggeredDeliveryOrdersService(settle Settle, status string) *ListPriceTriggeredDeliveryOrdersService {
	return &ListPriceTriggeredDeliveryOrdersService{c: c, settle: settle, params: map[string]string{"status": status}}
}

// SetContract filters the result to a single delivery contract (e.g. BTC_USDT_20241227).
func (s *ListPriceTriggeredDeliveryOrdersService) SetContract(contract string) *ListPriceTriggeredDeliveryOrdersService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListPriceTriggeredDeliveryOrdersService) SetLimit(limit int) *ListPriceTriggeredDeliveryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListPriceTriggeredDeliveryOrdersService) SetOffset(offset int) *ListPriceTriggeredDeliveryOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListPriceTriggeredDeliveryOrdersService) Do(ctx context.Context) ([]DeliveryPriceTriggeredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/price_orders", s.params).WithSign()
	resp, err := request.Do[[]DeliveryPriceTriggeredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreatePriceTriggeredDeliveryOrderService -- POST /api/v4/delivery/{settle}/price_orders (private)
//
// Places a price-triggered (auto) delivery order: once the market price
// satisfies the trigger rule, the queued initial order is submitted.
type CreatePriceTriggeredDeliveryOrderService struct {
	c         *DeliveryClient
	settle    Settle
	initial   map[string]any
	trigger   map[string]any
	orderType string
}

func (c *DeliveryClient) NewCreatePriceTriggeredDeliveryOrderService(settle Settle, contract string, size int64, price, triggerPrice decimal.Decimal, triggerRule, triggerExpiration int) *CreatePriceTriggeredDeliveryOrderService {
	return &CreatePriceTriggeredDeliveryOrderService{
		c:      c,
		settle: settle,
		initial: map[string]any{
			"contract": contract,
			"size":     size,
			"price":    price,
		},
		trigger: map[string]any{
			"price":      triggerPrice,
			"rule":       triggerRule,
			"expiration": triggerExpiration,
		},
	}
}

// SetInitialClose closes the whole position (single-position mode) when the
// queued order fires; size must be 0 in that case.
func (s *CreatePriceTriggeredDeliveryOrderService) SetInitialClose(closePosition bool) *CreatePriceTriggeredDeliveryOrderService {
	s.initial["close"] = closePosition
	return s
}

// SetInitialTif sets the queued order's time-in-force (gtc or ioc; market orders
// support ioc only).
func (s *CreatePriceTriggeredDeliveryOrderService) SetInitialTif(tif TimeInForce) *CreatePriceTriggeredDeliveryOrderService {
	s.initial["tif"] = tif
	return s
}

// SetInitialText attaches a user-defined label to the queued order.
func (s *CreatePriceTriggeredDeliveryOrderService) SetInitialText(text string) *CreatePriceTriggeredDeliveryOrderService {
	s.initial["text"] = text
	return s
}

// SetInitialReduceOnly restricts the queued order to reducing or closing the
// position, never opening a new one.
func (s *CreatePriceTriggeredDeliveryOrderService) SetInitialReduceOnly(reduceOnly bool) *CreatePriceTriggeredDeliveryOrderService {
	s.initial["reduce_only"] = reduceOnly
	return s
}

// SetInitialAutoSize selects which side to close in dual-position mode
// ("close_long" or "close_short").
func (s *CreatePriceTriggeredDeliveryOrderService) SetInitialAutoSize(autoSize string) *CreatePriceTriggeredDeliveryOrderService {
	s.initial["auto_size"] = autoSize
	return s
}

// SetTriggerStrategyType selects the trigger strategy (0: price trigger,
// 1: price-spread trigger).
func (s *CreatePriceTriggeredDeliveryOrderService) SetTriggerStrategyType(strategyType int) *CreatePriceTriggeredDeliveryOrderService {
	s.trigger["strategy_type"] = strategyType
	return s
}

// SetTriggerPriceType selects the reference price (0: last trade, 1: mark,
// 2: index).
func (s *CreatePriceTriggeredDeliveryOrderService) SetTriggerPriceType(priceType int) *CreatePriceTriggeredDeliveryOrderService {
	s.trigger["price_type"] = priceType
	return s
}

// SetOrderType selects the take-profit/stop-loss order type (e.g.
// "close-long-position", "plan-close-short-position").
func (s *CreatePriceTriggeredDeliveryOrderService) SetOrderType(orderType string) *CreatePriceTriggeredDeliveryOrderService {
	s.orderType = orderType
	return s
}

func (s *CreatePriceTriggeredDeliveryOrderService) Do(ctx context.Context) (*DeliveryTriggerOrderResponse, error) {
	body := map[string]any{
		"initial": s.initial,
		"trigger": s.trigger,
	}
	if s.orderType != "" {
		body["order_type"] = s.orderType
	}
	req := request.Post(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/price_orders").WithSign().SetBody(body)
	return request.Do[DeliveryTriggerOrderResponse](req)
}

// CancelPriceTriggeredDeliveryOrderListService -- DELETE /api/v4/delivery/{settle}/price_orders (private)
//
// Cancels all of the account's running price-triggered orders on a contract.
type CancelPriceTriggeredDeliveryOrderListService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewCancelPriceTriggeredDeliveryOrderListService(settle Settle, contract string) *CancelPriceTriggeredDeliveryOrderListService {
	return &CancelPriceTriggeredDeliveryOrderListService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

func (s *CancelPriceTriggeredDeliveryOrderListService) Do(ctx context.Context) ([]DeliveryPriceTriggeredOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/price_orders", s.params).WithSign()
	resp, err := request.Do[[]DeliveryPriceTriggeredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetPriceTriggeredDeliveryOrderService -- GET /api/v4/delivery/{settle}/price_orders/{order_id} (private)
//
// Returns a single price-triggered delivery order by its auto order ID.
type GetPriceTriggeredDeliveryOrderService struct {
	c       *DeliveryClient
	settle  Settle
	orderID string
}

func (c *DeliveryClient) NewGetPriceTriggeredDeliveryOrderService(settle Settle, orderID string) *GetPriceTriggeredDeliveryOrderService {
	return &GetPriceTriggeredDeliveryOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *GetPriceTriggeredDeliveryOrderService) Do(ctx context.Context) (*DeliveryPriceTriggeredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/price_orders/"+s.orderID).WithSign()
	return request.Do[DeliveryPriceTriggeredOrder](req)
}

// CancelPriceTriggeredDeliveryOrderService -- DELETE /api/v4/delivery/{settle}/price_orders/{order_id} (private)
//
// Cancels a single price-triggered delivery order and returns its final state.
type CancelPriceTriggeredDeliveryOrderService struct {
	c       *DeliveryClient
	settle  Settle
	orderID string
}

func (c *DeliveryClient) NewCancelPriceTriggeredDeliveryOrderService(settle Settle, orderID string) *CancelPriceTriggeredDeliveryOrderService {
	return &CancelPriceTriggeredDeliveryOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *CancelPriceTriggeredDeliveryOrderService) Do(ctx context.Context) (*DeliveryPriceTriggeredOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/price_orders/"+s.orderID).WithSign()
	return request.Do[DeliveryPriceTriggeredOrder](req)
}

// DeliveryPriceTriggeredOrder is a price-triggered (auto) delivery order: a
// queued initial order plus the market-price condition that releases it.
type DeliveryPriceTriggeredOrder struct {
	Initial    DeliveryInitialOrder `json:"initial"`
	Trigger    DeliveryPriceTrigger `json:"trigger"`
	ID         int64                `json:"id"`
	User       int64                `json:"user"`
	CreateTime time.Time            `json:"create_time,format:unix"`
	FinishTime time.Time            `json:"finish_time,format:unix"`
	TradeID    int64                `json:"trade_id"`
	Status     string               `json:"status"`
	FinishAs   string               `json:"finish_as"`
	Reason     string               `json:"reason"`
	OrderType  string               `json:"order_type"`
	MeOrderID  int64                `json:"me_order_id"`
}

// DeliveryPriceTrigger is the market-price condition that releases a queued
// initial order.
type DeliveryPriceTrigger struct {
	StrategyType int             `json:"strategy_type"`
	PriceType    int             `json:"price_type"`
	Price        decimal.Decimal `json:"price"`
	Rule         int             `json:"rule"`
	Expiration   int             `json:"expiration"`
}

// DeliveryInitialOrder is the order queued for submission once the trigger fires.
type DeliveryInitialOrder struct {
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
	Close        bool            `json:"close"`
	Tif          TimeInForce     `json:"tif"`
	Text         string          `json:"text"`
	ReduceOnly   bool            `json:"reduce_only"`
	AutoSize     string          `json:"auto_size"`
	IsReduceOnly bool            `json:"is_reduce_only"`
	IsClose      bool            `json:"is_close"`
}

// DeliveryTriggerOrderResponse is the auto order ID assigned to a newly created
// price-triggered order.
type DeliveryTriggerOrderResponse struct {
	ID int64 `json:"id"`
}
