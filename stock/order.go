package stock

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListOrdersService -- GET /api/v4/stock/orders (private)
//
// Lists the caller's open (active) orders, optionally narrowed to one symbol.
type ListOrdersService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c, params: map[string]string{}}
}

// SetSymbol narrows the result to a single symbol.
func (s *ListOrdersService) SetSymbol(symbol string) *ListOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *ListOrdersService) Do(ctx context.Context) (*StockOrdersResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/orders", s.params).WithSign()
	return request.Do[StockOrdersResponse](req)
}

// StockOrdersResponse is the envelope of the open-order query.
type StockOrdersResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		List []StockOrder `json:"list"`
	} `json:"data"`
}

// StockOrder is a single open order. side is 1=sell, 2=buy; price_type is
// "market" or "limit"; time_setup / time_update are integer-second Unix
// timestamps.
type StockOrder struct {
	OrderID             string          `json:"order_id"`
	Symbol              string          `json:"symbol"`
	Exchange            string          `json:"exchange"`
	QuoteCurrency       string          `json:"quote_currency"`
	FXRate              decimal.Decimal `json:"fx_rate"`
	SymbolDesc          string          `json:"symbol_desc"`
	TradeStatus         string          `json:"trade_status"`
	TradeMode           int             `json:"trade_mode"`
	PriceType           string          `json:"price_type"`
	Side                int             `json:"side"`
	Status              int             `json:"status"`
	Volume              decimal.Decimal `json:"volume"`
	FillVolume          decimal.Decimal `json:"fill_volume"`
	Price               decimal.Decimal `json:"price"`
	TimeSetup           time.Time       `json:"time_setup,format:unix"`
	TimeUpdate          time.Time       `json:"time_update,format:unix"`
	MaxOrderVolume      decimal.Decimal `json:"max_order_volume"`
	StepOrderVolume     decimal.Decimal `json:"step_order_volume"`
	MinOrderVolume      decimal.Decimal `json:"min_order_volume"`
	PricePrecision      int             `json:"price_precision"`
	PriceProtection     decimal.Decimal `json:"price_protection"`
	SellPriceProtection decimal.Decimal `json:"sell_price_protection"`
	BuyPriceProtection  decimal.Decimal `json:"buy_price_protection"`
	CommissionRate      decimal.Decimal `json:"commission_rate"`
	SlippageRate        decimal.Decimal `json:"slippage_rate"`
}

// CreateOrderService -- POST /api/v4/stock/orders (private)
//
// Places a single order. side is 1=sell, 2=buy; priceType is "market" or
// "limit"; tradingSession is "regular" (market orders only support regular) or
// "all" (limit orders only); timeInForce is "day" or "gtc". Price is required
// for limit orders.
type CreateOrderService struct {
	c    *StockClient
	body map[string]any
}

func (c *StockClient) NewCreateOrderService(symbol string, side int, priceType, tradingSession, timeInForce string, volume decimal.Decimal) *CreateOrderService {
	return &CreateOrderService{c: c, body: map[string]any{
		"symbol":          symbol,
		"side":            side,
		"price_type":      priceType,
		"trading_session": tradingSession,
		"time_in_force":   timeInForce,
		"volume":          volume.String(),
	}}
}

// SetPrice sets the order price (required for limit orders).
func (s *CreateOrderService) SetPrice(price decimal.Decimal) *CreateOrderService {
	s.body["price"] = price.String()
	return s
}

// SetClientOrderID sets a client-defined order id.
func (s *CreateOrderService) SetClientOrderID(clientOrderID string) *CreateOrderService {
	s.body["client_order_id"] = clientOrderID
	return s
}

func (s *CreateOrderService) Do(ctx context.Context) (*StockCreateOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/stock/orders", s.body).WithSign()
	return request.Do[StockCreateOrderResponse](req)
}

// StockCreateOrderResponse is the envelope of the create-order request.
type StockCreateOrderResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		ID string `json:"id"`
	} `json:"data"`
}

// CancelAllOrdersService -- DELETE /api/v4/stock/orders (private)
//
// Cancels all of the caller's open orders. Gate returns an empty object on
// success.
type CancelAllOrdersService struct {
	c *StockClient
}

func (c *StockClient) NewCancelAllOrdersService() *CancelAllOrdersService {
	return &CancelAllOrdersService{c: c}
}

func (s *CancelAllOrdersService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v4/stock/orders").WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListOrderHistoryService -- GET /api/v4/stock/orders/history (private)
//
// Lists the caller's historical orders, paginated. The queryable range must not
// exceed 3 months.
type ListOrderHistoryService struct {
	c      *StockClient
	params map[string]string
}

func (c *StockClient) NewListOrderHistoryService() *ListOrderHistoryService {
	return &ListOrderHistoryService{c: c, params: map[string]string{}}
}

// SetSymbol narrows the result to a single symbol.
func (s *ListOrderHistoryService) SetSymbol(symbol string) *ListOrderHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetOrderIDs narrows the result to a comma-separated list of order ids (max
// 20, each a positive integer).
func (s *ListOrderHistoryService) SetOrderIDs(orderIDs string) *ListOrderHistoryService {
	s.params["order_ids"] = orderIDs
	return s
}

// SetBeginTime bounds the result to orders at or after this time.
func (s *ListOrderHistoryService) SetBeginTime(beginTime time.Time) *ListOrderHistoryService {
	s.params["begin_time"] = strconv.FormatInt(beginTime.Unix(), 10)
	return s
}

// SetEndTime bounds the result to orders at or before this time.
func (s *ListOrderHistoryService) SetEndTime(endTime time.Time) *ListOrderHistoryService {
	s.params["end_time"] = strconv.FormatInt(endTime.Unix(), 10)
	return s
}

// SetSide narrows the result to an order side (1=sell, 2=buy).
func (s *ListOrderHistoryService) SetSide(side int) *ListOrderHistoryService {
	s.params["side"] = strconv.Itoa(side)
	return s
}

// SetPage selects the result page (defaults to 1).
func (s *ListOrderHistoryService) SetPage(page int) *ListOrderHistoryService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize caps the number of records per page (defaults to 10, max 500).
func (s *ListOrderHistoryService) SetPageSize(pageSize int) *ListOrderHistoryService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListOrderHistoryService) Do(ctx context.Context) (*StockOrderHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/stock/orders/history", s.params).WithSign()
	return request.Do[StockOrderHistoryResponse](req)
}

// StockOrderHistoryResponse is the envelope of the historical-order query.
type StockOrderHistoryResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		Total     int64               `json:"total"`
		TotalPage int                 `json:"total_page"`
		List      []StockOrderHistory `json:"list"`
	} `json:"data"`
}

// StockOrderHistory is a single historical order. side is 1=sell, 2=buy;
// time_setup / time_done are integer-second Unix timestamps; status_detail is
// null unless the server attaches a title/message; avg_fill_price is null until
// the order fills.
type StockOrderHistory struct {
	OrderID       string                  `json:"order_id"`
	Symbol        string                  `json:"symbol"`
	Exchange      string                  `json:"exchange"`
	QuoteCurrency string                  `json:"quote_currency"`
	FXRate        decimal.Decimal         `json:"fx_rate"`
	SymbolDesc    string                  `json:"symbol_desc"`
	PriceType     string                  `json:"price_type"`
	Status        int                     `json:"status"`
	StatusDesc    string                  `json:"status_desc"`
	StatusDetail  *StockOrderStatusDetail `json:"status_detail"`
	FinishAs      int                     `json:"finish_as"`
	Side          int                     `json:"side"`
	TimeInForce   string                  `json:"time_in_force"`
	Volume        decimal.Decimal         `json:"volume"`
	FillVolume    decimal.Decimal         `json:"fill_volume"`
	Price         decimal.Decimal         `json:"price"`
	AvgFillPrice  decimal.Decimal         `json:"avg_fill_price"`
	Commission    decimal.Decimal         `json:"commission"`
	TimeSetup     time.Time               `json:"time_setup,format:unix"`
	TimeDone      time.Time               `json:"time_done,format:unix"`
}

// StockOrderStatusDetail is the optional title/message describing an order's
// status.
type StockOrderStatusDetail struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// ModifyOrderService -- PUT /api/v4/stock/orders/{order_id} (private)
//
// Modifies an active order's quantity and price.
type ModifyOrderService struct {
	c       *StockClient
	orderID int64
	body    map[string]any
}

func (c *StockClient) NewModifyOrderService(orderID int64, volume, price decimal.Decimal) *ModifyOrderService {
	return &ModifyOrderService{c: c, orderID: orderID, body: map[string]any{
		"volume": volume.String(),
		"price":  price.String(),
	}}
}

func (s *ModifyOrderService) Do(ctx context.Context) (*StockModifyOrderResponse, error) {
	req := request.Put(ctx, s.c, "/api/v4/stock/orders/"+strconv.FormatInt(s.orderID, 10), s.body).WithSign()
	return request.Do[StockModifyOrderResponse](req)
}

// StockModifyOrderResponse is the envelope of the modify-order request.
type StockModifyOrderResponse struct {
	Label     string    `json:"label"`
	Timestamp time.Time `json:"timestamp,format:unixmilli"`
	Data      struct {
		OrderID int64 `json:"order_id"`
	} `json:"data"`
}

// CancelOrderService -- DELETE /api/v4/stock/orders/{order_id} (private)
//
// Cancels a single order. Gate returns an empty object on success.
type CancelOrderService struct {
	c       *StockClient
	orderID int64
}

func (c *StockClient) NewCancelOrderService(orderID int64) *CancelOrderService {
	return &CancelOrderService{c: c, orderID: orderID}
}

func (s *CancelOrderService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v4/stock/orders/"+strconv.FormatInt(s.orderID, 10)).WithSign()
	_, err := request.DoRaw(req)
	return err
}
