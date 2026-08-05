package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListSpotPOVOrdersService -- GET /api/v4/spot/pov_orders (private)
//
// Returns the account's POV (percentage-of-volume) strategy orders in the given
// status ("open" or "finished").
type ListSpotPOVOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListSpotPOVOrdersService(status string) *ListSpotPOVOrdersService {
	return &ListSpotPOVOrdersService{c: c, params: map[string]string{"status": status}}
}

// SetCurrencyPair filters the result to a single trading pair (e.g. BTC_USDT).
func (s *ListSpotPOVOrdersService) SetCurrencyPair(currencyPair string) *ListSpotPOVOrdersService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetSide narrows the result to all bids or all asks (both if unset).
func (s *ListSpotPOVOrdersService) SetSide(side Side) *ListSpotPOVOrdersService {
	s.params["side"] = string(side)
	return s
}

// SetPage sets the page number (up to 100).
func (s *ListSpotPOVOrdersService) SetPage(page int) *ListSpotPOVOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListSpotPOVOrdersService) SetLimit(limit int) *ListSpotPOVOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListSpotPOVOrdersService) Do(ctx context.Context) ([]SpotPOVOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/pov_orders", s.params).WithSign()
	resp, err := request.Do[[]SpotPOVOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateSpotPOVOrderService -- POST /api/v4/spot/pov_orders (private)
//
// Places a POV (percentage-of-volume) strategy order, which works the amount
// into the market at the target participation rate until it fills or the time
// to live expires. participationRate is a percentage (5, 10, 20 or 40); ttl is
// one of 1h, 6h, 12h, 1d, 2d, 3d, 4d, 5d, 6d or 7d.
type CreateSpotPOVOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCreateSpotPOVOrderService(currencyPair string, side Side, amount decimal.Decimal, participationRate int, ttl string) *CreateSpotPOVOrderService {
	return &CreateSpotPOVOrderService{c: c, body: map[string]any{
		"currency_pair":      currencyPair,
		"side":               string(side),
		"amount":             amount.String(),
		"participation_rate": participationRate,
		"ttl":                ttl,
	}}
}

// SetLimitPrice caps the execution price; the market price is used if unset.
func (s *CreateSpotPOVOrderService) SetLimitPrice(limitPrice decimal.Decimal) *CreateSpotPOVOrderService {
	s.body["limit_price"] = limitPrice.String()
	return s
}

// SetTriggerPrice delays the start of the strategy until the market reaches the
// given price; the order starts immediately if unset.
func (s *CreateSpotPOVOrderService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateSpotPOVOrderService {
	s.body["trigger_price"] = triggerPrice.String()
	return s
}

// SetText sets a custom order ID. It must start with "t-", be at most 28 bytes
// after that prefix, and use only letters, digits, "_", "-" or ".".
func (s *CreateSpotPOVOrderService) SetText(text string) *CreateSpotPOVOrderService {
	s.body["text"] = text
	return s
}

func (s *CreateSpotPOVOrderService) Do(ctx context.Context) (*SpotPOVOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/pov_orders", s.body).WithSign()
	return request.Do[SpotPOVOrder](req)
}

// CancelSpotPOVOrdersService -- POST /api/v4/spot/pov_orders/cancel (private)
//
// Cancels the account's running POV orders in bulk, optionally scoped to a
// single trading pair. The request only reports that the batch was accepted;
// success is determined from the returned order list.
type CancelSpotPOVOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewCancelSpotPOVOrdersService() *CancelSpotPOVOrdersService {
	return &CancelSpotPOVOrdersService{c: c, params: map[string]string{}}
}

// SetCurrencyPair limits the cancellation to a single trading pair.
func (s *CancelSpotPOVOrdersService) SetCurrencyPair(currencyPair string) *CancelSpotPOVOrdersService {
	s.params["currency_pair"] = currencyPair
	return s
}

func (s *CancelSpotPOVOrdersService) Do(ctx context.Context) ([]SpotPOVOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/pov_orders/cancel")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	resp, err := request.Do[[]SpotPOVOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetSpotPOVOrderService -- GET /api/v4/spot/pov_orders/{order_id} (private)
//
// Returns a single POV order by the ID returned at creation, or by the custom
// ID set through the order's text field.
type GetSpotPOVOrderService struct {
	c       *SpotClient
	orderID string
}

func (c *SpotClient) NewGetSpotPOVOrderService(orderID string) *GetSpotPOVOrderService {
	return &GetSpotPOVOrderService{c: c, orderID: orderID}
}

func (s *GetSpotPOVOrderService) Do(ctx context.Context) (*SpotPOVOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/pov_orders/"+s.orderID).WithSign()
	return request.Do[SpotPOVOrder](req)
}

// CancelSpotPOVOrderService -- POST /api/v4/spot/pov_orders/{order_id}/cancel (private)
//
// Cancels a single POV order, addressed by the ID returned at creation or by
// the custom ID set through the order's text field.
type CancelSpotPOVOrderService struct {
	c       *SpotClient
	orderID string
}

func (c *SpotClient) NewCancelSpotPOVOrderService(orderID string) *CancelSpotPOVOrderService {
	return &CancelSpotPOVOrderService{c: c, orderID: orderID}
}

func (s *CancelSpotPOVOrderService) Do(ctx context.Context) (*SpotPOVOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/spot/pov_orders/"+s.orderID+"/cancel").WithSign()
	return request.Do[SpotPOVOrder](req)
}

// SpotPOVOrder is a POV (percentage-of-volume) strategy order: an amount worked
// into the market at a target participation rate. Status is one of CREATED,
// CANCELING, RUNNING, COMPLETED, EXPIRED or TERMINATED; LimitPrice and
// TriggerPrice are empty when the order was created without them.
type SpotPOVOrder struct {
	ID                string          `json:"id"`
	CurrencyPair      string          `json:"currency_pair"`
	Side              Side            `json:"side"`
	Amount            decimal.Decimal `json:"amount"`
	ParticipationRate int             `json:"participation_rate"`
	TTL               string          `json:"ttl"`
	LimitPrice        decimal.Decimal `json:"limit_price"`
	TriggerPrice      decimal.Decimal `json:"trigger_price"`
	Status            string          `json:"status"`
	TerminatedAs      string          `json:"terminated_as"`
	StartTimeMs       time.Time       `json:"start_time_ms,format:unixmilli"`
	EndTimeMs         time.Time       `json:"end_time_ms,format:unixmilli"`
	ExpireTimeMs      time.Time       `json:"expire_time_ms,format:unixmilli"`
	CreateTimeMs      time.Time       `json:"create_time_ms,format:unixmilli"`
	UpdateTimeMs      time.Time       `json:"update_time_ms,format:unixmilli"`
	Text              string          `json:"text"`
}
