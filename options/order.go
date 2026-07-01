package options

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListOptionsOrdersService -- GET /api/v4/options/orders (private)
//
// Lists the account's options orders in a given status (open or finished).
type ListOptionsOrdersService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsOrdersService(status string) *ListOptionsOrdersService {
	return &ListOptionsOrdersService{c: c, params: map[string]string{
		"status": status,
	}}
}

// SetContract narrows the result to a single options contract.
func (s *ListOptionsOrdersService) SetContract(contract string) *ListOptionsOrdersService {
	s.params["contract"] = contract
	return s
}

// SetUnderlying narrows the result to a single underlying.
func (s *ListOptionsOrdersService) SetUnderlying(underlying string) *ListOptionsOrdersService {
	s.params["underlying"] = underlying
	return s
}

// SetLimit caps the number of records returned.
func (s *ListOptionsOrdersService) SetLimit(limit int) *ListOptionsOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListOptionsOrdersService) SetOffset(offset int) *ListOptionsOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive, Unix seconds).
func (s *ListOptionsOrdersService) SetFrom(from time.Time) *ListOptionsOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive, Unix seconds).
func (s *ListOptionsOrdersService) SetTo(to time.Time) *ListOptionsOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListOptionsOrdersService) Do(ctx context.Context) ([]OptionsOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/orders", s.params).WithSign()
	resp, err := request.Do[[]OptionsOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateOptionsOrderService -- POST /api/v4/options/orders (private)
//
// Places a single options order. size is the signed number of contracts
// (positive to buy, negative to sell, 0 with close=true to close a position).
type CreateOptionsOrderService struct {
	c    *OptionsClient
	body map[string]any
}

func (c *OptionsClient) NewCreateOptionsOrderService(contract string, size int64) *CreateOptionsOrderService {
	return &CreateOptionsOrderService{c: c, body: map[string]any{
		"contract": contract,
		"size":     size,
	}}
}

// SetPrice sets the order price. A price of "0" together with tif=ioc places a
// market order.
func (s *CreateOptionsOrderService) SetPrice(price decimal.Decimal) *CreateOptionsOrderService {
	s.body["price"] = price.String()
	return s
}

// SetTimeInForce selects how long the order stays active (gtc/ioc/poc).
func (s *CreateOptionsOrderService) SetTimeInForce(tif string) *CreateOptionsOrderService {
	s.body["tif"] = tif
	return s
}

// SetText attaches custom order information (must be prefixed with "t-").
func (s *CreateOptionsOrderService) SetText(text string) *CreateOptionsOrderService {
	s.body["text"] = text
	return s
}

// SetReduceOnly marks the order as reduce-only so it can only shrink a position.
func (s *CreateOptionsOrderService) SetReduceOnly(reduceOnly bool) *CreateOptionsOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetClose closes the position, with size set to 0.
func (s *CreateOptionsOrderService) SetClose(closePosition bool) *CreateOptionsOrderService {
	s.body["close"] = closePosition
	return s
}

// SetIceberg sets the display size for iceberg orders (0 for non-iceberg).
func (s *CreateOptionsOrderService) SetIceberg(iceberg int64) *CreateOptionsOrderService {
	s.body["iceberg"] = iceberg
	return s
}

func (s *CreateOptionsOrderService) Do(ctx context.Context) (*OptionsOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/options/orders", s.body).WithSign()
	return request.Do[OptionsOrder](req)
}

// CancelOptionsOrdersService -- DELETE /api/v4/options/orders (private)
//
// Cancels every open order, optionally limited to a contract, an underlying
// and/or a single side.
type CancelOptionsOrdersService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewCancelOptionsOrdersService() *CancelOptionsOrdersService {
	return &CancelOptionsOrdersService{c: c, params: map[string]string{}}
}

// SetContract narrows the cancellation to a single options contract.
func (s *CancelOptionsOrdersService) SetContract(contract string) *CancelOptionsOrdersService {
	s.params["contract"] = contract
	return s
}

// SetUnderlying narrows the cancellation to a single underlying.
func (s *CancelOptionsOrdersService) SetUnderlying(underlying string) *CancelOptionsOrdersService {
	s.params["underlying"] = underlying
	return s
}

// SetSide cancels only buy ("bid") or only sell ("ask") orders.
func (s *CancelOptionsOrdersService) SetSide(side string) *CancelOptionsOrdersService {
	s.params["side"] = side
	return s
}

func (s *CancelOptionsOrdersService) Do(ctx context.Context) ([]OptionsOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/options/orders", s.params).WithSign()
	resp, err := request.Do[[]OptionsOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOptionsOrderService -- GET /api/v4/options/orders/{order_id} (private)
//
// Returns a single options order by id.
type GetOptionsOrderService struct {
	c       *OptionsClient
	orderID string
}

func (c *OptionsClient) NewGetOptionsOrderService(orderID string) *GetOptionsOrderService {
	return &GetOptionsOrderService{c: c, orderID: orderID}
}

func (s *GetOptionsOrderService) Do(ctx context.Context) (*OptionsOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/orders/"+s.orderID).WithSign()
	return request.Do[OptionsOrder](req)
}

// CancelOptionsOrderService -- DELETE /api/v4/options/orders/{order_id} (private)
//
// Cancels a single options order by id.
type CancelOptionsOrderService struct {
	c       *OptionsClient
	orderID string
}

func (c *OptionsClient) NewCancelOptionsOrderService(orderID string) *CancelOptionsOrderService {
	return &CancelOptionsOrderService{c: c, orderID: orderID}
}

func (s *CancelOptionsOrderService) Do(ctx context.Context) (*OptionsOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/options/orders/"+s.orderID).WithSign()
	return request.Do[OptionsOrder](req)
}

// CountdownCancelAllOptionsService -- POST /api/v4/options/countdown_cancel_all (private)
//
// Arms a dead-man's-switch: if the countdown (>= 5 seconds) is not refreshed in
// time, Gate cancels all open options orders. A timeout of 0 disarms it.
type CountdownCancelAllOptionsService struct {
	c    *OptionsClient
	body map[string]any
}

func (c *OptionsClient) NewCountdownCancelAllOptionsService(timeout int) *CountdownCancelAllOptionsService {
	return &CountdownCancelAllOptionsService{c: c, body: map[string]any{
		"timeout": timeout,
	}}
}

// SetContract limits the countdown cancellation to a single contract.
func (s *CountdownCancelAllOptionsService) SetContract(contract string) *CountdownCancelAllOptionsService {
	s.body["contract"] = contract
	return s
}

// SetUnderlying limits the countdown cancellation to a single underlying.
func (s *CountdownCancelAllOptionsService) SetUnderlying(underlying string) *CountdownCancelAllOptionsService {
	s.body["underlying"] = underlying
	return s
}

func (s *CountdownCancelAllOptionsService) Do(ctx context.Context) (*OptionsCountdownStatus, error) {
	req := request.Post(ctx, s.c, "/api/v4/options/countdown_cancel_all", s.body).WithSign()
	return request.Do[OptionsCountdownStatus](req)
}

// OptionsOrder is an options order and its live state. create_time / finish_time
// are float-second Unix timestamps; size is signed (positive buy, negative sell)
// so there is no separate side field.
type OptionsOrder struct {
	ID           int64           `json:"id"`
	User         int64           `json:"user"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	FinishTime   time.Time       `json:"finish_time,format:unix"`
	FinishAs     string          `json:"finish_as"`
	Status       string          `json:"status"`
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Iceberg      int64           `json:"iceberg"`
	Price        decimal.Decimal `json:"price"`
	Close        bool            `json:"close"`
	IsClose      bool            `json:"is_close"`
	ReduceOnly   bool            `json:"reduce_only"`
	IsReduceOnly bool            `json:"is_reduce_only"`
	IsLiq        bool            `json:"is_liq"`
	Mmp          bool            `json:"mmp"`
	IsMmp        bool            `json:"is_mmp"`
	Tif          string          `json:"tif"`
	Left         int64           `json:"left"`
	FillPrice    decimal.Decimal `json:"fill_price"`
	Text         string          `json:"text"`
	Tkfr         decimal.Decimal `json:"tkfr"`
	Mkfr         decimal.Decimal `json:"mkfr"`
	Refu         int64           `json:"refu"`
	Refr         decimal.Decimal `json:"refr"`
}

// OptionsCountdownStatus reports when the armed countdown will fire. triggerTime
// is a millisecond epoch (0 when the countdown is disarmed).
type OptionsCountdownStatus struct {
	TriggerTime time.Time `json:"triggerTime,format:unixmilli"`
}
