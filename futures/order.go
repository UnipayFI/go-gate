package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateFuturesOrderService -- POST /api/v4/futures/{settle}/orders (private)
//
// Places a single futures order. size is the signed number of contracts
// (positive to buy, negative to sell, 0 with close=true to close a position).
type CreateFuturesOrderService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewCreateFuturesOrderService(settle Settle, contract string, size int64) *CreateFuturesOrderService {
	return &CreateFuturesOrderService{c: c, settle: settle, body: map[string]any{
		"contract": contract,
		"size":     size,
	}}
}

// SetPrice sets the order price. A price of "0" together with tif=ioc places a
// market order.
func (s *CreateFuturesOrderService) SetPrice(price decimal.Decimal) *CreateFuturesOrderService {
	s.body["price"] = price.String()
	return s
}

// SetTimeInForce selects how long the order stays active (gtc/ioc/poc/fok).
func (s *CreateFuturesOrderService) SetTimeInForce(tif TimeInForce) *CreateFuturesOrderService {
	s.body["tif"] = string(tif)
	return s
}

// SetText attaches custom order information (must be prefixed with "t-").
func (s *CreateFuturesOrderService) SetText(text string) *CreateFuturesOrderService {
	s.body["text"] = text
	return s
}

// SetReduceOnly marks the order as reduce-only so it can only shrink a position.
func (s *CreateFuturesOrderService) SetReduceOnly(reduceOnly bool) *CreateFuturesOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetClose closes the position in single-position mode (size must be 0).
func (s *CreateFuturesOrderService) SetClose(closePosition bool) *CreateFuturesOrderService {
	s.body["close"] = closePosition
	return s
}

// SetIceberg sets the display size for iceberg orders (0 for non-iceberg).
func (s *CreateFuturesOrderService) SetIceberg(iceberg int64) *CreateFuturesOrderService {
	s.body["iceberg"] = iceberg
	return s
}

// SetAutoSize closes one leg of a dual-mode position (size must be 0).
func (s *CreateFuturesOrderService) SetAutoSize(autoSize AutoSize) *CreateFuturesOrderService {
	s.body["auto_size"] = string(autoSize)
	return s
}

// SetStpAct sets the self-trade-prevention action.
func (s *CreateFuturesOrderService) SetStpAct(stpAct StpAct) *CreateFuturesOrderService {
	s.body["stp_act"] = string(stpAct)
	return s
}

func (s *CreateFuturesOrderService) Do(ctx context.Context) (*FuturesOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders", s.body).WithSign()
	return request.Do[FuturesOrder](req)
}

// ListFuturesOrdersService -- GET /api/v4/futures/{settle}/orders (private)
//
// Lists the account's futures orders in a given status (open or finished).
type ListFuturesOrdersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesOrdersService(settle Settle, status OrderStatus) *ListFuturesOrdersService {
	return &ListFuturesOrdersService{c: c, settle: settle, params: map[string]string{
		"status": string(status),
	}}
}

// SetContract narrows the result to a single futures contract.
func (s *ListFuturesOrdersService) SetContract(contract string) *ListFuturesOrdersService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned.
func (s *ListFuturesOrdersService) SetLimit(limit int) *ListFuturesOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListFuturesOrdersService) SetOffset(offset int) *ListFuturesOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetLastID paginates from a previously retrieved order id (finished orders).
func (s *ListFuturesOrdersService) SetLastID(lastID string) *ListFuturesOrdersService {
	s.params["last_id"] = lastID
	return s
}

func (s *ListFuturesOrdersService) Do(ctx context.Context) ([]FuturesOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders", s.params).WithSign()
	resp, err := request.Do[[]FuturesOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CancelFuturesOrdersService -- DELETE /api/v4/futures/{settle}/orders (private)
//
// Cancels every open order on a contract, optionally limited to one side.
type CancelFuturesOrdersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewCancelFuturesOrdersService(settle Settle, contract string) *CancelFuturesOrdersService {
	return &CancelFuturesOrdersService{c: c, settle: settle, params: map[string]string{
		"contract": contract,
	}}
}

// SetSide cancels only buy ("bid") or only sell ("ask") orders.
func (s *CancelFuturesOrdersService) SetSide(side string) *CancelFuturesOrdersService {
	s.params["side"] = side
	return s
}

// SetExcludeReduceOnly excludes reduce-only orders from cancellation.
func (s *CancelFuturesOrdersService) SetExcludeReduceOnly(exclude bool) *CancelFuturesOrdersService {
	s.params["exclude_reduce_only"] = strconv.FormatBool(exclude)
	return s
}

// SetText attaches a remark for the cancellation.
func (s *CancelFuturesOrdersService) SetText(text string) *CancelFuturesOrdersService {
	s.params["text"] = text
	return s
}

func (s *CancelFuturesOrdersService) Do(ctx context.Context) ([]FuturesOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders", s.params).WithSign()
	resp, err := request.Do[[]FuturesOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOrdersWithTimeRangeService -- GET /api/v4/futures/{settle}/orders_timerange (private)
//
// Lists the account's futures orders within a Unix-second time range.
type GetOrdersWithTimeRangeService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetOrdersWithTimeRangeService(settle Settle) *GetOrdersWithTimeRangeService {
	return &GetOrdersWithTimeRangeService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *GetOrdersWithTimeRangeService) SetContract(contract string) *GetOrdersWithTimeRangeService {
	s.params["contract"] = contract
	return s
}

// SetFrom sets the start time (inclusive).
func (s *GetOrdersWithTimeRangeService) SetFrom(from time.Time) *GetOrdersWithTimeRangeService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive).
func (s *GetOrdersWithTimeRangeService) SetTo(to time.Time) *GetOrdersWithTimeRangeService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned.
func (s *GetOrdersWithTimeRangeService) SetLimit(limit int) *GetOrdersWithTimeRangeService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *GetOrdersWithTimeRangeService) SetOffset(offset int) *GetOrdersWithTimeRangeService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *GetOrdersWithTimeRangeService) Do(ctx context.Context) ([]FuturesOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders_timerange", s.params).WithSign()
	resp, err := request.Do[[]FuturesOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateBatchFuturesOrderService -- POST /api/v4/futures/{settle}/batch_orders (private)
//
// Places up to 10 futures orders in one request. The result array corresponds
// positionally to the submitted orders; inspect FuturesOrderBatchResult.Succeeded
// to tell placed orders from rejected ones.
type CreateBatchFuturesOrderService struct {
	c      *FuturesClient
	settle Settle
	orders []*CreateFuturesOrderService
}

func (c *FuturesClient) NewCreateBatchFuturesOrderService(settle Settle, orders ...*CreateFuturesOrderService) *CreateBatchFuturesOrderService {
	return &CreateBatchFuturesOrderService{c: c, settle: settle, orders: orders}
}

func (s *CreateBatchFuturesOrderService) Do(ctx context.Context) ([]FuturesOrderBatchResult, error) {
	bodies := make([]map[string]any, 0, len(s.orders))
	for _, o := range s.orders {
		bodies = append(bodies, o.body)
	}
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/batch_orders").WithSign().SetBody(bodies)
	resp, err := request.Do[[]FuturesOrderBatchResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetFuturesOrderService -- GET /api/v4/futures/{settle}/orders/{order_id} (private)
//
// Returns a single futures order by its id (or user custom text id).
type GetFuturesOrderService struct {
	c       *FuturesClient
	settle  Settle
	orderID string
}

func (c *FuturesClient) NewGetFuturesOrderService(settle Settle, orderID string) *GetFuturesOrderService {
	return &GetFuturesOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *GetFuturesOrderService) Do(ctx context.Context) (*FuturesOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders/"+s.orderID).WithSign()
	return request.Do[FuturesOrder](req)
}

// AmendFuturesOrderService -- PUT /api/v4/futures/{settle}/orders/{order_id} (private)
//
// Amends the size and/or price of an open futures order.
type AmendFuturesOrderService struct {
	c       *FuturesClient
	settle  Settle
	orderID string
	body    map[string]any
}

func (c *FuturesClient) NewAmendFuturesOrderService(settle Settle, orderID string) *AmendFuturesOrderService {
	return &AmendFuturesOrderService{c: c, settle: settle, orderID: orderID, body: map[string]any{}}
}

// SetSize sets the new order size (including the already-filled part).
func (s *AmendFuturesOrderService) SetSize(size int64) *AmendFuturesOrderService {
	s.body["size"] = size
	return s
}

// SetPrice sets the new order price.
func (s *AmendFuturesOrderService) SetPrice(price decimal.Decimal) *AmendFuturesOrderService {
	s.body["price"] = price.String()
	return s
}

// SetAmendText attaches custom info recorded with the amendment.
func (s *AmendFuturesOrderService) SetAmendText(amendText string) *AmendFuturesOrderService {
	s.body["amend_text"] = amendText
	return s
}

func (s *AmendFuturesOrderService) Do(ctx context.Context) (*FuturesOrder, error) {
	req := request.Put(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders/"+s.orderID, s.body).WithSign()
	return request.Do[FuturesOrder](req)
}

// CancelFuturesOrderService -- DELETE /api/v4/futures/{settle}/orders/{order_id} (private)
//
// Cancels a single futures order by its id (or user custom text id).
type CancelFuturesOrderService struct {
	c       *FuturesClient
	settle  Settle
	orderID string
}

func (c *FuturesClient) NewCancelFuturesOrderService(settle Settle, orderID string) *CancelFuturesOrderService {
	return &CancelFuturesOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *CancelFuturesOrderService) Do(ctx context.Context) (*FuturesOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/orders/"+s.orderID).WithSign()
	return request.Do[FuturesOrder](req)
}

// GetMyTradesService -- GET /api/v4/futures/{settle}/my_trades (private)
//
// Lists the account's personal trade fills, covering roughly the past 6 months.
type GetMyTradesService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetMyTradesService(settle Settle) *GetMyTradesService {
	return &GetMyTradesService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *GetMyTradesService) SetContract(contract string) *GetMyTradesService {
	s.params["contract"] = contract
	return s
}

// SetOrder narrows the result to fills of a single futures order.
func (s *GetMyTradesService) SetOrder(orderID int64) *GetMyTradesService {
	s.params["order"] = strconv.FormatInt(orderID, 10)
	return s
}

// SetLimit caps the number of records returned.
func (s *GetMyTradesService) SetLimit(limit int) *GetMyTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *GetMyTradesService) SetOffset(offset int) *GetMyTradesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetLastID paginates from a previously retrieved fill id (deprecated by Gate).
func (s *GetMyTradesService) SetLastID(lastID string) *GetMyTradesService {
	s.params["last_id"] = lastID
	return s
}

func (s *GetMyTradesService) Do(ctx context.Context) ([]FuturesMyTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/my_trades", s.params).WithSign()
	resp, err := request.Do[[]FuturesMyTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetMyTradesWithTimeRangeService -- GET /api/v4/futures/{settle}/my_trades_timerange (private)
//
// Lists the account's personal trade fills within a Unix-second time range.
type GetMyTradesWithTimeRangeService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetMyTradesWithTimeRangeService(settle Settle) *GetMyTradesWithTimeRangeService {
	return &GetMyTradesWithTimeRangeService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *GetMyTradesWithTimeRangeService) SetContract(contract string) *GetMyTradesWithTimeRangeService {
	s.params["contract"] = contract
	return s
}

// SetFrom sets the start time (inclusive).
func (s *GetMyTradesWithTimeRangeService) SetFrom(from time.Time) *GetMyTradesWithTimeRangeService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive).
func (s *GetMyTradesWithTimeRangeService) SetTo(to time.Time) *GetMyTradesWithTimeRangeService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned.
func (s *GetMyTradesWithTimeRangeService) SetLimit(limit int) *GetMyTradesWithTimeRangeService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *GetMyTradesWithTimeRangeService) SetOffset(offset int) *GetMyTradesWithTimeRangeService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetRole narrows the result to a trade role ("maker" or "taker").
func (s *GetMyTradesWithTimeRangeService) SetRole(role string) *GetMyTradesWithTimeRangeService {
	s.params["role"] = role
	return s
}

func (s *GetMyTradesWithTimeRangeService) Do(ctx context.Context) ([]FuturesMyTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/my_trades_timerange", s.params).WithSign()
	resp, err := request.Do[[]FuturesMyTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListPositionCloseService -- GET /api/v4/futures/{settle}/position_close (private)
//
// Lists the account's closed-position history and realized PnL.
type ListPositionCloseService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListPositionCloseService(settle Settle) *ListPositionCloseService {
	return &ListPositionCloseService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *ListPositionCloseService) SetContract(contract string) *ListPositionCloseService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned.
func (s *ListPositionCloseService) SetLimit(limit int) *ListPositionCloseService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListPositionCloseService) SetOffset(offset int) *ListPositionCloseService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive).
func (s *ListPositionCloseService) SetFrom(from time.Time) *ListPositionCloseService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive).
func (s *ListPositionCloseService) SetTo(to time.Time) *ListPositionCloseService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetSide narrows the result to a position side ("long" or "short").
func (s *ListPositionCloseService) SetSide(side string) *ListPositionCloseService {
	s.params["side"] = side
	return s
}

// SetPnL narrows the result by realized profit or loss.
func (s *ListPositionCloseService) SetPnL(pnl string) *ListPositionCloseService {
	s.params["pnl"] = pnl
	return s
}

func (s *ListPositionCloseService) Do(ctx context.Context) ([]PositionClose, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/position_close", s.params).WithSign()
	resp, err := request.Do[[]PositionClose](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListLiquidatesService -- GET /api/v4/futures/{settle}/liquidates (private)
//
// Lists the account's forced-liquidation history.
type ListLiquidatesService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListLiquidatesService(settle Settle) *ListLiquidatesService {
	return &ListLiquidatesService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *ListLiquidatesService) SetContract(contract string) *ListLiquidatesService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned.
func (s *ListLiquidatesService) SetLimit(limit int) *ListLiquidatesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListLiquidatesService) SetOffset(offset int) *ListLiquidatesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive).
func (s *ListLiquidatesService) SetFrom(from time.Time) *ListLiquidatesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive).
func (s *ListLiquidatesService) SetTo(to time.Time) *ListLiquidatesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetAt narrows the result to a specific liquidation timestamp.
func (s *ListLiquidatesService) SetAt(at time.Time) *ListLiquidatesService {
	s.params["at"] = strconv.FormatInt(at.Unix(), 10)
	return s
}

func (s *ListLiquidatesService) Do(ctx context.Context) ([]FuturesLiquidateRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/liquidates", s.params).WithSign()
	resp, err := request.Do[[]FuturesLiquidateRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListAutoDeleveragesService -- GET /api/v4/futures/{settle}/auto_deleverages (private)
//
// Lists the account's ADL (auto-deleveraging) history.
type ListAutoDeleveragesService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListAutoDeleveragesService(settle Settle) *ListAutoDeleveragesService {
	return &ListAutoDeleveragesService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single futures contract.
func (s *ListAutoDeleveragesService) SetContract(contract string) *ListAutoDeleveragesService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned.
func (s *ListAutoDeleveragesService) SetLimit(limit int) *ListAutoDeleveragesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListAutoDeleveragesService) SetOffset(offset int) *ListAutoDeleveragesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive).
func (s *ListAutoDeleveragesService) SetFrom(from time.Time) *ListAutoDeleveragesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive).
func (s *ListAutoDeleveragesService) SetTo(to time.Time) *ListAutoDeleveragesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetAt narrows the result to a specific auto-deleveraging timestamp.
func (s *ListAutoDeleveragesService) SetAt(at time.Time) *ListAutoDeleveragesService {
	s.params["at"] = strconv.FormatInt(at.Unix(), 10)
	return s
}

func (s *ListAutoDeleveragesService) Do(ctx context.Context) ([]FuturesAutoDeleverage, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/auto_deleverages", s.params).WithSign()
	resp, err := request.Do[[]FuturesAutoDeleverage](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CountdownCancelAllFuturesService -- POST /api/v4/futures/{settle}/countdown_cancel_all (private)
//
// Arms a dead-man's-switch: if the countdown (>= 5 seconds) is not refreshed in
// time, Gate cancels all open orders. A timeout of 0 disarms it.
type CountdownCancelAllFuturesService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewCountdownCancelAllFuturesService(settle Settle, timeout int) *CountdownCancelAllFuturesService {
	return &CountdownCancelAllFuturesService{c: c, settle: settle, body: map[string]any{
		"timeout": timeout,
	}}
}

// SetContract limits the countdown cancellation to a single contract.
func (s *CountdownCancelAllFuturesService) SetContract(contract string) *CountdownCancelAllFuturesService {
	s.body["contract"] = contract
	return s
}

func (s *CountdownCancelAllFuturesService) Do(ctx context.Context) (*FuturesCountdownStatus, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/countdown_cancel_all", s.body).WithSign()
	return request.Do[FuturesCountdownStatus](req)
}

// CancelBatchFutureOrdersService -- POST /api/v4/futures/{settle}/batch_cancel_orders (private)
//
// Cancels a batch of futures orders by id (up to 20 per request). Each result's
// Succeeded field reports whether that id was cancelled.
type CancelBatchFutureOrdersService struct {
	c        *FuturesClient
	settle   Settle
	orderIDs []string
}

func (c *FuturesClient) NewCancelBatchFutureOrdersService(settle Settle, orderIDs []string) *CancelBatchFutureOrdersService {
	return &CancelBatchFutureOrdersService{c: c, settle: settle, orderIDs: orderIDs}
}

func (s *CancelBatchFutureOrdersService) Do(ctx context.Context) ([]FuturesCancelBatchResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/batch_cancel_orders").WithSign().SetBody(s.orderIDs)
	resp, err := request.Do[[]FuturesCancelBatchResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AmendBatchFutureOrdersService -- POST /api/v4/futures/{settle}/batch_amend_orders (private)
//
// Amends a batch of futures orders by id (up to 10 per request).
type AmendBatchFutureOrdersService struct {
	c      *FuturesClient
	settle Settle
	items  []BatchAmendItem
}

func (c *FuturesClient) NewAmendBatchFutureOrdersService(settle Settle, items []BatchAmendItem) *AmendBatchFutureOrdersService {
	return &AmendBatchFutureOrdersService{c: c, settle: settle, items: items}
}

func (s *AmendBatchFutureOrdersService) Do(ctx context.Context) ([]FuturesOrderBatchResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/batch_amend_orders").WithSign().SetBody(s.items)
	resp, err := request.Do[[]FuturesOrderBatchResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BatchAmendItem is one entry in an AmendBatchFutureOrders request body. Only the
// non-zero fields are sent, so leave the ones you are not changing untouched.
type BatchAmendItem struct {
	OrderID   int64           `json:"order_id,omitzero"`
	Size      int64           `json:"size,omitzero"`
	Price     decimal.Decimal `json:"price,omitzero"`
	AmendText string          `json:"amend_text,omitzero"`
}

// FuturesOrder is a futures order and its live state. create_time / update_time /
// finish_time are float-second Unix timestamps; size is signed (positive long,
// negative short) so there is no separate side field.
type FuturesOrder struct {
	ID           int64           `json:"id"`
	User         int64           `json:"user"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	UpdateTime   time.Time       `json:"update_time,format:unix"`
	FinishTime   time.Time       `json:"finish_time,format:unix"`
	FinishAs     FinishAs        `json:"finish_as"`
	Status       OrderStatus     `json:"status"`
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Iceberg      int64           `json:"iceberg"`
	Price        decimal.Decimal `json:"price"`
	Close        bool            `json:"close"`
	IsClose      bool            `json:"is_close"`
	ReduceOnly   bool            `json:"reduce_only"`
	IsReduceOnly bool            `json:"is_reduce_only"`
	IsLiq        bool            `json:"is_liq"`
	Tif          TimeInForce     `json:"tif"`
	Left         int64           `json:"left"`
	FillPrice    decimal.Decimal `json:"fill_price"`
	Text         string          `json:"text"`
	Tkfr         decimal.Decimal `json:"tkfr"`
	Mkfr         decimal.Decimal `json:"mkfr"`
	Refu         int64           `json:"refu"`
	AutoSize     AutoSize        `json:"auto_size"`
	StpID        int64           `json:"stp_id"`
	StpAct       StpAct          `json:"stp_act"`
	AmendText    string          `json:"amend_text"`
	LimitVip     int64           `json:"limit_vip"`
	Pid          int64           `json:"pid"`

	// Live-only fields not present in the official model.
	BBO                  string          `json:"bbo"`
	BizInfo              string          `json:"biz_info"`
	Leverage             decimal.Decimal `json:"leverage"`
	MarketOrderSlipRatio decimal.Decimal `json:"market_order_slip_ratio"`
	PnL                  decimal.Decimal `json:"pnl"`
	PnLMargin            decimal.Decimal `json:"pnl_margin"`
	PosMarginMode        string          `json:"pos_margin_mode"`
	Refr                 decimal.Decimal `json:"refr"`
	UpdateID             int64           `json:"update_id"`
}

// FuturesOrderBatchResult is one element of a batch order/cancel/amend response.
// It embeds FuturesOrder (populated when the operation touched a real order) and
// adds the per-item execution outcome. batch_orders / batch_amend_orders fill
// Label/Detail on failure; batch_cancel_orders fills UserID/Message instead.
type FuturesOrderBatchResult struct {
	FuturesOrder
	Succeeded bool   `json:"succeeded"`
	Label     string `json:"label"`
	Detail    string `json:"detail"`
	Message   string `json:"message"`
	UserID    int64  `json:"user_id"`
}

// FuturesCancelBatchResult is one element of a batch_cancel_orders response. Unlike
// batch create/amend, cancel returns a minimal shape whose id is a STRING.
type FuturesCancelBatchResult struct {
	UserID    int64  `json:"user_id"`
	ID        string `json:"id"`
	Succeeded bool   `json:"succeeded"`
	Message   string `json:"message"`
	Label     string `json:"label"`
}

// FuturesMyTrade is one personal trade fill. GetMyTrades returns the fill id in
// "id" (int64) while GetMyTradesWithTimeRange returns it in "trade_id" (string);
// both keys are captured here. create_time is a float-second Unix timestamp.
type FuturesMyTrade struct {
	ID         int64           `json:"id"`
	TradeID    string          `json:"trade_id"`
	CreateTime time.Time       `json:"create_time,format:unix"`
	Contract   string          `json:"contract"`
	OrderID    string          `json:"order_id"`
	Size       int64           `json:"size"`
	CloseSize  int64           `json:"close_size"`
	Price      decimal.Decimal `json:"price"`
	Role       string          `json:"role"`
	Text       string          `json:"text"`
	Fee        decimal.Decimal `json:"fee"`
	PointFee   decimal.Decimal `json:"point_fee"`
}

// PositionClose is one closed-position record with its realized-PnL breakdown.
// time is a float-second Unix timestamp; first_open_time is an integer-second one.
type PositionClose struct {
	Time          time.Time       `json:"time,format:unix"`
	Contract      string          `json:"contract"`
	Side          string          `json:"side"`
	PnL           decimal.Decimal `json:"pnl"`
	PnLPnL        decimal.Decimal `json:"pnl_pnl"`
	PnLFund       decimal.Decimal `json:"pnl_fund"`
	PnLFee        decimal.Decimal `json:"pnl_fee"`
	Text          string          `json:"text"`
	MaxSize       decimal.Decimal `json:"max_size"`
	AccumSize     decimal.Decimal `json:"accum_size"`
	FirstOpenTime time.Time       `json:"first_open_time,format:unix"`
	LongPrice     decimal.Decimal `json:"long_price"`
	ShortPrice    decimal.Decimal `json:"short_price"`
}

// FuturesLiquidateRecord is one forced-liquidation record. time is an
// integer-second Unix timestamp. The margin/price fields are only returned on
// the authenticated endpoint.
type FuturesLiquidateRecord struct {
	Time       time.Time       `json:"time,format:unix"`
	Contract   string          `json:"contract"`
	Leverage   decimal.Decimal `json:"leverage"`
	Size       int64           `json:"size"`
	Margin     decimal.Decimal `json:"margin"`
	EntryPrice decimal.Decimal `json:"entry_price"`
	LiqPrice   decimal.Decimal `json:"liq_price"`
	MarkPrice  decimal.Decimal `json:"mark_price"`
	OrderID    int64           `json:"order_id"`
	OrderPrice decimal.Decimal `json:"order_price"`
	FillPrice  decimal.Decimal `json:"fill_price"`
	Left       int64           `json:"left"`
}

// FuturesAutoDeleverage is one ADL (auto-deleveraging) record. time is an
// integer-second Unix timestamp.
type FuturesAutoDeleverage struct {
	Time               time.Time       `json:"time,format:unix"`
	User               int64           `json:"user"`
	OrderID            int64           `json:"order_id"`
	Contract           string          `json:"contract"`
	Leverage           decimal.Decimal `json:"leverage"`
	CrossLeverageLimit decimal.Decimal `json:"cross_leverage_limit"`
	EntryPrice         decimal.Decimal `json:"entry_price"`
	FillPrice          decimal.Decimal `json:"fill_price"`
	TradeSize          int64           `json:"trade_size"`
	PositionSize       int64           `json:"position_size"`
}

// FuturesCountdownStatus reports when the armed countdown will fire. triggerTime
// is a millisecond Unix timestamp (0 when the countdown was disarmed).
type FuturesCountdownStatus struct {
	TriggerTime time.Time `json:"triggerTime,format:unixmilli"`
}
