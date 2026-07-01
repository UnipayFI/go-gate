package delivery

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListDeliveryOrdersService -- GET /api/v4/delivery/{settle}/orders (private)
//
// Lists the account's delivery orders in a given status (open or finished).
type ListDeliveryOrdersService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryOrdersService(settle Settle, status OrderStatus) *ListDeliveryOrdersService {
	return &ListDeliveryOrdersService{
		c:      c,
		settle: settle,
		params: map[string]string{"status": string(status)},
	}
}

// SetContract narrows the result to a single delivery contract.
func (s *ListDeliveryOrdersService) SetContract(contract string) *ListDeliveryOrdersService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListDeliveryOrdersService) SetLimit(limit int) *ListDeliveryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the given number of records (for pagination).
func (s *ListDeliveryOrdersService) SetOffset(offset int) *ListDeliveryOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetLastID returns records after the given order id (finished orders only).
func (s *ListDeliveryOrdersService) SetLastID(lastID string) *ListDeliveryOrdersService {
	s.params["last_id"] = lastID
	return s
}

// SetCountTotal requests the total record count via the response headers
// (1 to enable, 0 to disable).
func (s *ListDeliveryOrdersService) SetCountTotal(countTotal int) *ListDeliveryOrdersService {
	s.params["count_total"] = strconv.Itoa(countTotal)
	return s
}

func (s *ListDeliveryOrdersService) Do(ctx context.Context) ([]DeliveryOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/orders", s.params).WithSign()
	resp, err := request.Do[[]DeliveryOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateDeliveryOrderService -- POST /api/v4/delivery/{settle}/orders (private)
//
// Places a delivery order. size is the signed contract quantity (positive for
// long, negative for short; 0 with close/auto_size closes a position).
type CreateDeliveryOrderService struct {
	c      *DeliveryClient
	settle Settle
	body   map[string]any
}

func (c *DeliveryClient) NewCreateDeliveryOrderService(settle Settle, contract string, size int64) *CreateDeliveryOrderService {
	return &CreateDeliveryOrderService{
		c:      c,
		settle: settle,
		body: map[string]any{
			"contract": contract,
			"size":     size,
		},
	}
}

// SetPrice sets the order price. A price of 0 with tif set to ioc is a market order.
func (s *CreateDeliveryOrderService) SetPrice(price decimal.Decimal) *CreateDeliveryOrderService {
	s.body["price"] = price
	return s
}

// SetTif sets the time-in-force strategy (gtc, ioc, poc, fok).
func (s *CreateDeliveryOrderService) SetTif(tif TimeInForce) *CreateDeliveryOrderService {
	s.body["tif"] = string(tif)
	return s
}

// SetText sets the order's source/user-defined text (must be prefixed with t-).
func (s *CreateDeliveryOrderService) SetText(text string) *CreateDeliveryOrderService {
	s.body["text"] = text
	return s
}

// SetReduceOnly marks the order as reduce-only.
func (s *CreateDeliveryOrderService) SetReduceOnly(reduceOnly bool) *CreateDeliveryOrderService {
	s.body["reduce_only"] = reduceOnly
	return s
}

// SetClose marks the order as a full-position close (size must be 0).
func (s *CreateDeliveryOrderService) SetClose(closePosition bool) *CreateDeliveryOrderService {
	s.body["close"] = closePosition
	return s
}

// SetIceberg sets the display size for an iceberg order (0 for a normal order).
func (s *CreateDeliveryOrderService) SetIceberg(iceberg int64) *CreateDeliveryOrderService {
	s.body["iceberg"] = iceberg
	return s
}

// SetAutoSize closes one leg of a dual-mode position ("close_long" or
// "close_short"; size must be 0).
func (s *CreateDeliveryOrderService) SetAutoSize(autoSize string) *CreateDeliveryOrderService {
	s.body["auto_size"] = autoSize
	return s
}

// SetStpAct sets the self-trade-prevention action (cn, co, cb).
func (s *CreateDeliveryOrderService) SetStpAct(stpAct StpAct) *CreateDeliveryOrderService {
	s.body["stp_act"] = string(stpAct)
	return s
}

func (s *CreateDeliveryOrderService) Do(ctx context.Context) (*DeliveryOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/orders", s.body).WithSign()
	return request.Do[DeliveryOrder](req)
}

// CancelDeliveryOrdersService -- DELETE /api/v4/delivery/{settle}/orders (private)
//
// Cancels all of the account's open delivery orders on a contract, returning the
// orders that were cancelled.
type CancelDeliveryOrdersService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewCancelDeliveryOrdersService(settle Settle, contract string) *CancelDeliveryOrdersService {
	return &CancelDeliveryOrdersService{
		c:      c,
		settle: settle,
		params: map[string]string{"contract": contract},
	}
}

// SetSide restricts the cancellation to one side ("ask" for all sells, "bid" for
// all buys); both are cancelled when unset.
func (s *CancelDeliveryOrdersService) SetSide(side string) *CancelDeliveryOrdersService {
	s.params["side"] = side
	return s
}

func (s *CancelDeliveryOrdersService) Do(ctx context.Context) ([]DeliveryOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/orders", s.params).WithSign()
	resp, err := request.Do[[]DeliveryOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetDeliveryOrderService -- GET /api/v4/delivery/{settle}/orders/{order_id} (private)
//
// Returns the details of a single delivery order.
type GetDeliveryOrderService struct {
	c       *DeliveryClient
	settle  Settle
	orderID string
}

func (c *DeliveryClient) NewGetDeliveryOrderService(settle Settle, orderID string) *GetDeliveryOrderService {
	return &GetDeliveryOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *GetDeliveryOrderService) Do(ctx context.Context) (*DeliveryOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/orders/"+s.orderID).WithSign()
	return request.Do[DeliveryOrder](req)
}

// CancelDeliveryOrderService -- DELETE /api/v4/delivery/{settle}/orders/{order_id} (private)
//
// Cancels a single delivery order, returning its final state.
type CancelDeliveryOrderService struct {
	c       *DeliveryClient
	settle  Settle
	orderID string
}

func (c *DeliveryClient) NewCancelDeliveryOrderService(settle Settle, orderID string) *CancelDeliveryOrderService {
	return &CancelDeliveryOrderService{c: c, settle: settle, orderID: orderID}
}

func (s *CancelDeliveryOrderService) Do(ctx context.Context) (*DeliveryOrder, error) {
	req := request.Delete(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/orders/"+s.orderID).WithSign()
	return request.Do[DeliveryOrder](req)
}

// GetMyDeliveryTradesService -- GET /api/v4/delivery/{settle}/my_trades (private)
//
// Lists the account's personal delivery trade (fill) records.
type GetMyDeliveryTradesService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewGetMyDeliveryTradesService(settle Settle) *GetMyDeliveryTradesService {
	return &GetMyDeliveryTradesService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single delivery contract.
func (s *GetMyDeliveryTradesService) SetContract(contract string) *GetMyDeliveryTradesService {
	s.params["contract"] = contract
	return s
}

// SetOrder returns only the fills belonging to the given order id.
func (s *GetMyDeliveryTradesService) SetOrder(order int64) *GetMyDeliveryTradesService {
	s.params["order"] = strconv.FormatInt(order, 10)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *GetMyDeliveryTradesService) SetLimit(limit int) *GetMyDeliveryTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the given number of records (for pagination).
func (s *GetMyDeliveryTradesService) SetOffset(offset int) *GetMyDeliveryTradesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetLastID returns records after the given fill id.
func (s *GetMyDeliveryTradesService) SetLastID(lastID string) *GetMyDeliveryTradesService {
	s.params["last_id"] = lastID
	return s
}

// SetCountTotal requests the total record count via the response headers
// (1 to enable, 0 to disable).
func (s *GetMyDeliveryTradesService) SetCountTotal(countTotal int) *GetMyDeliveryTradesService {
	s.params["count_total"] = strconv.Itoa(countTotal)
	return s
}

func (s *GetMyDeliveryTradesService) Do(ctx context.Context) ([]DeliveryMyTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/my_trades", s.params).WithSign()
	resp, err := request.Do[[]DeliveryMyTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListDeliveryPositionCloseService -- GET /api/v4/delivery/{settle}/position_close (private)
//
// Lists the account's position-close (realised PnL) history.
type ListDeliveryPositionCloseService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryPositionCloseService(settle Settle) *ListDeliveryPositionCloseService {
	return &ListDeliveryPositionCloseService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single delivery contract.
func (s *ListDeliveryPositionCloseService) SetContract(contract string) *ListDeliveryPositionCloseService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListDeliveryPositionCloseService) SetLimit(limit int) *ListDeliveryPositionCloseService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListDeliveryPositionCloseService) Do(ctx context.Context) ([]DeliveryPositionClose, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/position_close", s.params).WithSign()
	resp, err := request.Do[[]DeliveryPositionClose](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListDeliveryLiquidatesService -- GET /api/v4/delivery/{settle}/liquidates (private)
//
// Lists the account's forced-liquidation history.
type ListDeliveryLiquidatesService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryLiquidatesService(settle Settle) *ListDeliveryLiquidatesService {
	return &ListDeliveryLiquidatesService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single delivery contract.
func (s *ListDeliveryLiquidatesService) SetContract(contract string) *ListDeliveryLiquidatesService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListDeliveryLiquidatesService) SetLimit(limit int) *ListDeliveryLiquidatesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetAt filters to liquidations at the given time (unix seconds); 0 for all.
func (s *ListDeliveryLiquidatesService) SetAt(at int) *ListDeliveryLiquidatesService {
	s.params["at"] = strconv.Itoa(at)
	return s
}

func (s *ListDeliveryLiquidatesService) Do(ctx context.Context) ([]DeliveryLiquidate, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/liquidates", s.params).WithSign()
	resp, err := request.Do[[]DeliveryLiquidate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListDeliverySettlementsService -- GET /api/v4/delivery/{settle}/settlements (private)
//
// Lists the account's delivery settlement records.
type ListDeliverySettlementsService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliverySettlementsService(settle Settle) *ListDeliverySettlementsService {
	return &ListDeliverySettlementsService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single delivery contract.
func (s *ListDeliverySettlementsService) SetContract(contract string) *ListDeliverySettlementsService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListDeliverySettlementsService) SetLimit(limit int) *ListDeliverySettlementsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetAt filters to settlements at the given time (unix seconds); 0 for all.
func (s *ListDeliverySettlementsService) SetAt(at int) *ListDeliverySettlementsService {
	s.params["at"] = strconv.Itoa(at)
	return s
}

func (s *ListDeliverySettlementsService) Do(ctx context.Context) ([]DeliverySettlement, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/settlements", s.params).WithSign()
	resp, err := request.Do[[]DeliverySettlement](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListDeliveryRiskLimitTiersService -- GET /api/v4/delivery/{settle}/risk_limit_tiers
//
// Lists the risk-limit tiers of a delivery contract. When contract is unset the
// top markets are returned and limit/offset paginate at the market level.
type ListDeliveryRiskLimitTiersService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryRiskLimitTiersService(settle Settle) *ListDeliveryRiskLimitTiersService {
	return &ListDeliveryRiskLimitTiersService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single delivery contract.
func (s *ListDeliveryRiskLimitTiersService) SetContract(contract string) *ListDeliveryRiskLimitTiersService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of markets returned when contract is unset.
func (s *ListDeliveryRiskLimitTiersService) SetLimit(limit int) *ListDeliveryRiskLimitTiersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset paginates over markets when contract is unset.
func (s *ListDeliveryRiskLimitTiersService) SetOffset(offset int) *ListDeliveryRiskLimitTiersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListDeliveryRiskLimitTiersService) Do(ctx context.Context) ([]DeliveryRiskLimitTier, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/risk_limit_tiers", s.params)
	resp, err := request.Do[[]DeliveryRiskLimitTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeliveryOrder is a delivery (dated-futures) order and its lifecycle metadata.
type DeliveryOrder struct {
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
	TimeInForce  TimeInForce     `json:"tif"`
	Left         int64           `json:"left"`
	FillPrice    decimal.Decimal `json:"fill_price"`
	Text         string          `json:"text"`
	Tkfr         decimal.Decimal `json:"tkfr"`
	Mkfr         decimal.Decimal `json:"mkfr"`
	Refu         int64           `json:"refu"`
	AutoSize     string          `json:"auto_size"`
	StpID        int64           `json:"stp_id"`
	StpAct       StpAct          `json:"stp_act"`
	AmendText    string          `json:"amend_text"`
	LimitVIP     int64           `json:"limit_vip"`
	PID          int64           `json:"pid"`
}

// DeliveryMyTrade is a single personal delivery fill.
type DeliveryMyTrade struct {
	ID         int64           `json:"id"`
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

// DeliveryPositionClose is a realised-PnL record produced when a position closes.
type DeliveryPositionClose struct {
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

// DeliveryLiquidate is a forced-liquidation history record.
type DeliveryLiquidate struct {
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

// DeliverySettlement is a delivery settlement record (final mark-to-settle PnL).
type DeliverySettlement struct {
	Time        time.Time       `json:"time,format:unix"`
	Contract    string          `json:"contract"`
	Leverage    decimal.Decimal `json:"leverage"`
	Size        int64           `json:"size"`
	Margin      decimal.Decimal `json:"margin"`
	EntryPrice  decimal.Decimal `json:"entry_price"`
	SettlePrice decimal.Decimal `json:"settle_price"`
	Profit      decimal.Decimal `json:"profit"`
	Fee         decimal.Decimal `json:"fee"`
}

// DeliveryRiskLimitTier is one risk-limit tier of a delivery contract.
type DeliveryRiskLimitTier struct {
	Tier            int             `json:"tier"`
	RiskLimit       decimal.Decimal `json:"risk_limit"`
	InitialRate     decimal.Decimal `json:"initial_rate"`
	MaintenanceRate decimal.Decimal `json:"maintenance_rate"`
	LeverageMax     decimal.Decimal `json:"leverage_max"`
	Contract        string          `json:"contract"`
	Deduction       decimal.Decimal `json:"deduction"`
}
