package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListSpotPriceTriggeredOrdersService -- GET /api/v4/spot/price_orders (private)
//
// Returns the authenticated account's price-triggered (auto) spot orders in the
// given status ("open" or "finished").
type ListSpotPriceTriggeredOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListSpotPriceTriggeredOrdersService(status string) *ListSpotPriceTriggeredOrdersService {
	return &ListSpotPriceTriggeredOrdersService{c: c, params: map[string]string{"status": status}}
}

// SetMarket filters the result to a single trading market (e.g. BTC_USDT).
func (s *ListSpotPriceTriggeredOrdersService) SetMarket(market string) *ListSpotPriceTriggeredOrdersService {
	s.params["market"] = market
	return s
}

// SetAccount filters by trading account type (unified accounts must pass "unified").
func (s *ListSpotPriceTriggeredOrdersService) SetAccount(account Account) *ListSpotPriceTriggeredOrdersService {
	s.params["account"] = string(account)
	return s
}

// SetLimit caps the number of records returned (default 100).
func (s *ListSpotPriceTriggeredOrdersService) SetLimit(limit int) *ListSpotPriceTriggeredOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListSpotPriceTriggeredOrdersService) SetOffset(offset int) *ListSpotPriceTriggeredOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListSpotPriceTriggeredOrdersService) Do(ctx context.Context) ([]SpotPriceTriggeredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/price_orders", s.params).WithSign()
	resp, err := request.Do[[]SpotPriceTriggeredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateSpotPriceTriggeredOrderService -- POST /api/v4/spot/price_orders (private)
//
// Places a price-triggered (auto) spot order: once the market price satisfies
// the trigger rule, the queued put order is submitted.
type CreateSpotPriceTriggeredOrderService struct {
	c       *SpotClient
	market  string
	trigger map[string]any
	put     map[string]any
}

func (c *SpotClient) NewCreateSpotPriceTriggeredOrderService(market string, triggerPrice decimal.Decimal, triggerRule string, triggerExpiration int, putSide Side, putPrice, putAmount decimal.Decimal, putAccount Account) *CreateSpotPriceTriggeredOrderService {
	return &CreateSpotPriceTriggeredOrderService{
		c:      c,
		market: market,
		trigger: map[string]any{
			"price":      triggerPrice,
			"rule":       triggerRule,
			"expiration": triggerExpiration,
		},
		put: map[string]any{
			"side":    putSide,
			"price":   putPrice,
			"amount":  putAmount,
			"account": putAccount,
		},
	}
}

// SetPutType sets the queued order's execution type (defaults to limit).
func (s *CreateSpotPriceTriggeredOrderService) SetPutType(orderType OrderType) *CreateSpotPriceTriggeredOrderService {
	s.put["type"] = orderType
	return s
}

// SetPutTimeInForce sets the queued order's time-in-force (gtc or ioc).
func (s *CreateSpotPriceTriggeredOrderService) SetPutTimeInForce(tif TimeInForce) *CreateSpotPriceTriggeredOrderService {
	s.put["time_in_force"] = tif
	return s
}

// SetPutAutoBorrow toggles automatic borrowing for margin puts.
func (s *CreateSpotPriceTriggeredOrderService) SetPutAutoBorrow(autoBorrow bool) *CreateSpotPriceTriggeredOrderService {
	s.put["auto_borrow"] = autoBorrow
	return s
}

// SetPutAutoRepay toggles automatic loan repayment for margin puts.
func (s *CreateSpotPriceTriggeredOrderService) SetPutAutoRepay(autoRepay bool) *CreateSpotPriceTriggeredOrderService {
	s.put["auto_repay"] = autoRepay
	return s
}

// SetPutText attaches a user-defined label to the queued order.
func (s *CreateSpotPriceTriggeredOrderService) SetPutText(text string) *CreateSpotPriceTriggeredOrderService {
	s.put["text"] = text
	return s
}

func (s *CreateSpotPriceTriggeredOrderService) Do(ctx context.Context) (*SpotTriggerOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/price_orders").WithSign().SetBody(map[string]any{
		"market":  s.market,
		"trigger": s.trigger,
		"put":     s.put,
	})
	return request.Do[SpotTriggerOrderResponse](req)
}

// CancelSpotPriceTriggeredOrderListService -- DELETE /api/v4/spot/price_orders (private)
//
// Cancels all of the account's running price-triggered orders, optionally
// scoped to a single market and/or account type.
type CancelSpotPriceTriggeredOrderListService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelSpotPriceTriggeredOrderListService() *CancelSpotPriceTriggeredOrderListService {
	return &CancelSpotPriceTriggeredOrderListService{c: c, params: map[string]string{}}
}

// SetMarket limits the cancellation to a single trading market (e.g. BTC_USDT).
func (s *CancelSpotPriceTriggeredOrderListService) SetMarket(market string) *CancelSpotPriceTriggeredOrderListService {
	s.params["market"] = market
	return s
}

// SetAccount limits the cancellation by trading account type (unified accounts
// must pass "unified").
func (s *CancelSpotPriceTriggeredOrderListService) SetAccount(account Account) *CancelSpotPriceTriggeredOrderListService {
	s.params["account"] = string(account)
	return s
}

func (s *CancelSpotPriceTriggeredOrderListService) Do(ctx context.Context) ([]SpotPriceTriggeredOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/spot/price_orders", s.params).WithSign()
	resp, err := request.Do[[]SpotPriceTriggeredOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetSpotPriceTriggeredOrderService -- GET /api/v4/spot/price_orders/{order_id} (private)
//
// Returns a single price-triggered order by its auto order ID.
type GetSpotPriceTriggeredOrderService struct {
	c       *SpotClient
	orderID string
}

func (c *SpotClient) NewGetSpotPriceTriggeredOrderService(orderID string) *GetSpotPriceTriggeredOrderService {
	return &GetSpotPriceTriggeredOrderService{c: c, orderID: orderID}
}

func (s *GetSpotPriceTriggeredOrderService) Do(ctx context.Context) (*SpotPriceTriggeredOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/price_orders/"+s.orderID).WithSign()
	return request.Do[SpotPriceTriggeredOrder](req)
}

// CancelSpotPriceTriggeredOrderService -- DELETE /api/v4/spot/price_orders/{order_id} (private)
//
// Cancels a single price-triggered order and returns its final state.
type CancelSpotPriceTriggeredOrderService struct {
	c       *SpotClient
	orderID string
}

func (c *SpotClient) NewCancelSpotPriceTriggeredOrderService(orderID string) *CancelSpotPriceTriggeredOrderService {
	return &CancelSpotPriceTriggeredOrderService{c: c, orderID: orderID}
}

func (s *CancelSpotPriceTriggeredOrderService) Do(ctx context.Context) (*SpotPriceTriggeredOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/spot/price_orders/"+s.orderID).WithSign()
	return request.Do[SpotPriceTriggeredOrder](req)
}

// SpotPriceTriggeredOrder is a price-triggered (auto) spot order: a queued put
// order plus the market-price condition that releases it.
type SpotPriceTriggeredOrder struct {
	Trigger      SpotPriceTrigger  `json:"trigger"`
	Put          SpotPricePutOrder `json:"put"`
	ID           int64             `json:"id"`
	User         int64             `json:"user"`
	Market       string            `json:"market"`
	Ctime        time.Time         `json:"ctime,format:unix"`
	Ftime        time.Time         `json:"ftime,format:unix"`
	FiredOrderID int64             `json:"fired_order_id"`
	Status       string            `json:"status"`
	Reason       string            `json:"reason"`
}

// SpotPriceTrigger is the market-price condition that releases a queued put order.
type SpotPriceTrigger struct {
	Price      decimal.Decimal `json:"price"`
	Rule       string          `json:"rule"`
	Expiration int             `json:"expiration"`
}

// SpotPricePutOrder is the order queued for submission once the trigger fires.
type SpotPricePutOrder struct {
	Type        OrderType       `json:"type"`
	Side        Side            `json:"side"`
	Price       decimal.Decimal `json:"price"`
	Amount      decimal.Decimal `json:"amount"`
	Account     Account         `json:"account"`
	TimeInForce TimeInForce     `json:"time_in_force"`
	AutoBorrow  bool            `json:"auto_borrow"`
	AutoRepay   bool            `json:"auto_repay"`
	Text        string          `json:"text"`
}

// SpotTriggerOrderResponse is the auto order ID assigned to a newly created
// price-triggered order.
type SpotTriggerOrderResponse struct {
	ID int64 `json:"id"`
}
