package futures

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateBBOOrderService -- POST /api/v4/futures/{settle}/bbo_orders (private)
//
// Places a level-based BBO (best-bid/offer) futures order: the order price is
// taken from the order book at the requested depth level rather than being set
// explicitly. size is the signed contract quantity (positive buy, negative sell);
// direction selects which book side to price against ("buy" prices against the
// ask, "sell" against the bid); level is the depth level (max 20).
type CreateBBOOrderService struct {
	c       *FuturesClient
	settle  Settle
	body    map[string]any
	expTime string
}

func (c *FuturesClient) NewCreateBBOOrderService(settle Settle, contract string, size int64, direction string, level int64) *CreateBBOOrderService {
	return &CreateBBOOrderService{c: c, settle: settle, body: map[string]any{
		"contract":  contract,
		"size":      size,
		"direction": direction,
		"level":     level,
	}}
}

// SetIceberg sets the display size for iceberg orders (0 for non-iceberg).
func (s *CreateBBOOrderService) SetIceberg(iceberg int64) *CreateBBOOrderService {
	s.body["iceberg"] = iceberg
	return s
}

// SetClose closes the position in single-position mode (size must be 0).
func (s *CreateBBOOrderService) SetClose(closePosition bool) *CreateBBOOrderService {
	s.body["close"] = closePosition
	return s
}

// SetReduceOnly marks the order as reduce-only so it can only shrink a position.
func (s *CreateBBOOrderService) SetReduceOnly(reduceOnly bool) *CreateBBOOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetTimeInForce selects how long the order stays active (gtc/ioc/poc/fok).
func (s *CreateBBOOrderService) SetTimeInForce(tif TimeInForce) *CreateBBOOrderService {
	s.body["tif"] = string(tif)
	return s
}

// SetText attaches custom order information (must be prefixed with "t-").
func (s *CreateBBOOrderService) SetText(text string) *CreateBBOOrderService {
	s.body["text"] = text
	return s
}

// SetAutoSize closes one leg of a dual-mode position (size must be 0).
func (s *CreateBBOOrderService) SetAutoSize(autoSize AutoSize) *CreateBBOOrderService {
	s.body["auto_size"] = string(autoSize)
	return s
}

// SetStpAct sets the self-trade-prevention action.
func (s *CreateBBOOrderService) SetStpAct(stpAct StpAct) *CreateBBOOrderService {
	s.body["stp_act"] = string(stpAct)
	return s
}

// SetPID sets the position id the order applies to.
func (s *CreateBBOOrderService) SetPID(pid int64) *CreateBBOOrderService {
	s.body["pid"] = pid
	return s
}

// SetExpireTime sets the X-Gate-Exptime header: the order is discarded if it does
// not reach the matching engine before this Unix-millisecond deadline.
func (s *CreateBBOOrderService) SetExpireTime(expireTimeMs int64) *CreateBBOOrderService {
	s.expTime = decimal.NewFromInt(expireTimeMs).String()
	return s
}

func (s *CreateBBOOrderService) Do(ctx context.Context) (*FuturesOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/bbo_orders", s.body).WithSign()
	req.SetHeader("X-Gate-Exptime", s.expTime)
	return request.Do[FuturesOrder](req)
}

// AmendPriceTriggeredOrderService -- PUT /api/v4/futures/{settle}/price_orders/amend (private)
//
// Amends a single pending price-triggered (auto) order in place. Only the fields
// you set are changed; orderID identifies the pending trigger order to modify.
type AmendPriceTriggeredOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewAmendPriceTriggeredOrderService(settle Settle, orderID int64) *AmendPriceTriggeredOrderService {
	return &AmendPriceTriggeredOrderService{c: c, settle: settle, body: map[string]any{
		"order_id": orderID,
	}}
}

// SetSize sets the modified contract quantity (0 to fully close, or a partial size).
func (s *AmendPriceTriggeredOrderService) SetSize(size int64) *AmendPriceTriggeredOrderService {
	s.body["size"] = size
	return s
}

// SetAmount sets the modified size for decimal-size contracts (same role as size).
func (s *AmendPriceTriggeredOrderService) SetAmount(amount decimal.Decimal) *AmendPriceTriggeredOrderService {
	s.body["amount"] = amount.String()
	return s
}

// SetPrice sets the modified initial order price (a value of "0" places a market order).
func (s *AmendPriceTriggeredOrderService) SetPrice(price decimal.Decimal) *AmendPriceTriggeredOrderService {
	s.body["price"] = price.String()
	return s
}

// SetTriggerPrice sets the modified trigger price.
func (s *AmendPriceTriggeredOrderService) SetTriggerPrice(triggerPrice decimal.Decimal) *AmendPriceTriggeredOrderService {
	s.body["trigger_price"] = triggerPrice.String()
	return s
}

// SetPriceType sets the reference price type (0 last, 1 mark, 2 index).
func (s *AmendPriceTriggeredOrderService) SetPriceType(priceType int) *AmendPriceTriggeredOrderService {
	s.body["price_type"] = priceType
	return s
}

// SetAutoSize sets the hedge-mode close side ("close_long" or "close_short").
func (s *AmendPriceTriggeredOrderService) SetAutoSize(autoSize AutoSize) *AmendPriceTriggeredOrderService {
	s.body["auto_size"] = string(autoSize)
	return s
}

// SetClose marks the amendment as a full-position close in single-position mode.
func (s *AmendPriceTriggeredOrderService) SetClose(closePosition bool) *AmendPriceTriggeredOrderService {
	s.body["close"] = closePosition
	return s
}

func (s *AmendPriceTriggeredOrderService) Do(ctx context.Context) (*FuturesAmendPriceOrderResult, error) {
	req := request.Put(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/price_orders/amend", s.body).WithSign()
	return request.Do[FuturesAmendPriceOrderResult](req)
}

// FuturesAmendPriceOrderResult is the id of an amended price-triggered order,
// returned both as an integer and as its string form.
type FuturesAmendPriceOrderResult struct {
	ID       int64  `json:"id"`
	IDString string `json:"id_string"`
}
