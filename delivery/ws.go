package delivery

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/common"
	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// DeliveryWebSocketClient is the stream client for Gate's delivery (dated-futures)
// channels on the wss://fx-ws.gateio.ws/v4/ws/delivery/usdt gateway. Delivery
// channels share the perpetual-futures shapes and are prefixed "futures."
// (futures.tickers, futures.orders, ...) even though they flow over the delivery
// endpoint. Public channels need no credentials; private channels (orders,
// usertrades, positions, balances) carry a per-subscription auth object built
// from WithWebSocketAuth and take the account's user id as the first payload
// element.
type DeliveryWebSocketClient struct {
	*client.WebSocketClient
}

// NewDeliveryWebSocketClient constructs a delivery stream client.
func NewDeliveryWebSocketClient(options ...client.WebSocketOptions) *DeliveryWebSocketClient {
	return &DeliveryWebSocketClient{client.NewWebSocketClient(common.DEFAULT_WS_DELIVERY_USDT_URL, options...)}
}

// WsHandler is invoked for every push (or error) on a subscription. The push's
// Result field is already decoded into T.
type WsHandler[T any] func(*request.WsPush[T], error)

// --- Public channels ---

// SubscribeTickersService -- futures.tickers channel (24h ticker updates). The
// payload is the list of contracts; each push carries a batch of tickers.
type SubscribeTickersService struct {
	c         *DeliveryWebSocketClient
	contracts []string
}

func (c *DeliveryWebSocketClient) NewSubscribeTickersService(contracts ...string) *SubscribeTickersService {
	return &SubscribeTickersService{c: c, contracts: contracts}
}

func (s *SubscribeTickersService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryTicker](ctx, s.c, "futures.tickers", s.contracts, false, cb)
}

// WsDeliveryTicker is a delivery ticker push. The funding-rate/quanto fields are
// inherited from the shared futures ticker shape and are typically absent (zero)
// on delivery contracts.
type WsDeliveryTicker struct {
	Contract              string          `json:"contract"`
	Last                  decimal.Decimal `json:"last"`
	ChangePercentage      decimal.Decimal `json:"change_percentage"`
	TotalSize             decimal.Decimal `json:"total_size"`
	Volume24h             decimal.Decimal `json:"volume_24h"`
	Volume24hBase         decimal.Decimal `json:"volume_24h_base"`
	Volume24hQuote        decimal.Decimal `json:"volume_24h_quote"`
	Volume24hSettle       decimal.Decimal `json:"volume_24h_settle"`
	Volume24USD           decimal.Decimal `json:"volume_24_usd"`
	Volume24BTC           decimal.Decimal `json:"volume_24_btc"`
	MarkPrice             decimal.Decimal `json:"mark_price"`
	FundingRate           decimal.Decimal `json:"funding_rate"`
	FundingRateIndicative decimal.Decimal `json:"funding_rate_indicative"`
	IndexPrice            decimal.Decimal `json:"index_price"`
	QuantoBaseRate        decimal.Decimal `json:"quanto_base_rate"`
	Low24h                decimal.Decimal `json:"low_24h"`
	High24h               decimal.Decimal `json:"high_24h"`
}

// SubscribeTradesService -- futures.trades channel (public tick-by-tick fills).
type SubscribeTradesService struct {
	c        *DeliveryWebSocketClient
	contract string
}

func (c *DeliveryWebSocketClient) NewSubscribeTradesService(contract string) *SubscribeTradesService {
	return &SubscribeTradesService{c: c, contract: contract}
}

func (s *SubscribeTradesService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryTrade](ctx, s.c, "futures.trades", []string{s.contract}, false, cb)
}

// WsDeliveryTrade is a delivery public-trade push. Size is signed (negative for a
// taker sell). CreateTime is unix seconds, CreateTimeMs unix milliseconds, both
// bare numbers.
type WsDeliveryTrade struct {
	ID           int64           `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,format:unixmilli"`
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
}

// SubscribeBookTickerService -- futures.book_ticker channel (best bid/ask updates).
type SubscribeBookTickerService struct {
	c        *DeliveryWebSocketClient
	contract string
}

func (c *DeliveryWebSocketClient) NewSubscribeBookTickerService(contract string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, contract: contract}
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, cb WsHandler[WsDeliveryBookTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsDeliveryBookTicker](ctx, s.c, "futures.book_ticker", []string{s.contract}, false, cb)
}

// WsDeliveryBookTicker is a best bid/ask push. Time (t) is unix milliseconds;
// bid/ask sizes are signed contract counts.
type WsDeliveryBookTicker struct {
	Time         time.Time       `json:"t,format:unixmilli"`
	UpdateID     int64           `json:"u"`
	Contract     string          `json:"s"`
	BestBidPrice decimal.Decimal `json:"b"`
	BestBidSize  int64           `json:"B"`
	BestAskPrice decimal.Decimal `json:"a"`
	BestAskSize  int64           `json:"A"`
}

// SubscribeOrderBookService -- futures.order_book channel (limited-depth
// snapshots). level is "5", "10", "20", "50", "100"; interval is "0" or "100ms".
type SubscribeOrderBookService struct {
	c        *DeliveryWebSocketClient
	contract string
	level    string
	interval string
}

func (c *DeliveryWebSocketClient) NewSubscribeOrderBookService(contract, level, interval string) *SubscribeOrderBookService {
	return &SubscribeOrderBookService{c: c, contract: contract, level: level, interval: interval}
}

func (s *SubscribeOrderBookService) Do(ctx context.Context, cb WsHandler[WsDeliveryOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsDeliveryOrderBook](ctx, s.c, "futures.order_book", []string{s.contract, s.level, s.interval}, false, cb)
}

// WsDeliveryOrderBook is a limited-depth order-book snapshot. Time (t) is unix
// milliseconds; ID increments on every order-book change.
type WsDeliveryOrderBook struct {
	Time     time.Time                 `json:"t,format:unixmilli"`
	Contract string                    `json:"contract"`
	ID       int64                     `json:"id"`
	Asks     []WsDeliveryOrderBookItem `json:"asks"`
	Bids     []WsDeliveryOrderBookItem `json:"bids"`
}

// WsDeliveryOrderBookItem is a single depth level: Price (quote currency) and
// Size (contract count).
type WsDeliveryOrderBookItem struct {
	Price decimal.Decimal `json:"p"`
	Size  int64           `json:"s"`
}

// SubscribeOrderBookUpdateService -- futures.order_book_update channel
// (incremental depth). interval is "100ms" or "1000ms".
type SubscribeOrderBookUpdateService struct {
	c        *DeliveryWebSocketClient
	contract string
	interval string
}

func (c *DeliveryWebSocketClient) NewSubscribeOrderBookUpdateService(contract, interval string) *SubscribeOrderBookUpdateService {
	return &SubscribeOrderBookUpdateService{c: c, contract: contract, interval: interval}
}

func (s *SubscribeOrderBookUpdateService) Do(ctx context.Context, cb WsHandler[WsDeliveryOrderBookUpdate]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsDeliveryOrderBookUpdate](ctx, s.c, "futures.order_book_update", []string{s.contract, s.interval}, false, cb)
}

// WsDeliveryOrderBookUpdate is an incremental order-book change. Time (t) is unix
// milliseconds; FirstID/LastID bound the change ids covered by this delta.
type WsDeliveryOrderBookUpdate struct {
	Time     time.Time                 `json:"t,format:unixmilli"`
	Contract string                    `json:"s"`
	FirstID  int64                     `json:"U"`
	LastID   int64                     `json:"u"`
	Bids     []WsDeliveryOrderBookItem `json:"b"`
	Asks     []WsDeliveryOrderBookItem `json:"a"`
}

// SubscribeCandlesticksService -- futures.candlesticks channel. interval is e.g.
// "10s", "1m", "1h"; the payload is [interval, contract].
type SubscribeCandlesticksService struct {
	c        *DeliveryWebSocketClient
	interval string
	contract string
}

func (c *DeliveryWebSocketClient) NewSubscribeCandlesticksService(interval, contract string) *SubscribeCandlesticksService {
	return &SubscribeCandlesticksService{c: c, interval: interval, contract: contract}
}

func (s *SubscribeCandlesticksService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryCandlestick]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryCandlestick](ctx, s.c, "futures.candlesticks", []string{s.interval, s.contract}, false, cb)
}

// WsDeliveryCandlestick is a candlestick push. Name ("n") is "<interval>_<contract>";
// Time (t) is unix seconds; Volume (v, contract size) is only present when the
// contract is not prefixed with mark_/index_.
type WsDeliveryCandlestick struct {
	Time   time.Time       `json:"t,format:unix"`
	Volume int64           `json:"v"`
	Close  decimal.Decimal `json:"c"`
	High   decimal.Decimal `json:"h"`
	Low    decimal.Decimal `json:"l"`
	Open   decimal.Decimal `json:"o"`
	Name   string          `json:"n"`
	Amount decimal.Decimal `json:"a"`
}

// --- Private channels (require WithWebSocketAuth) ---

// SubscribeOrdersService -- futures.orders channel (own order updates). The
// payload is [user_id, contract]; pass "!all" as the contract to receive updates
// for every contract.
type SubscribeOrdersService struct {
	c        *DeliveryWebSocketClient
	userID   string
	contract string
}

func (c *DeliveryWebSocketClient) NewSubscribeOrdersService(userID, contract string) *SubscribeOrdersService {
	return &SubscribeOrdersService{c: c, userID: userID, contract: contract}
}

func (s *SubscribeOrdersService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryOrder](ctx, s.c, "futures.orders", []string{s.userID, s.contract}, true, cb)
}

// WsDeliveryOrder is an own-order update push. User is the account id sent as a
// string; timestamps are bare numbers (create_time/finish_time in seconds, the
// _ms variants in milliseconds). Size is signed.
type WsDeliveryOrder struct {
	ID              int64           `json:"id"`
	User            string          `json:"user"`
	Contract        string          `json:"contract"`
	CreateTime      time.Time       `json:"create_time,format:unix"`
	CreateTimeMs    time.Time       `json:"create_time_ms,format:unixmilli"`
	FinishTime      time.Time       `json:"finish_time,format:unix"`
	FinishTimeMs    time.Time       `json:"finish_time_ms,format:unixmilli"`
	FinishAs        FinishAs        `json:"finish_as"`
	Status          OrderStatus     `json:"status"`
	Size            int64           `json:"size"`
	Iceberg         int64           `json:"iceberg"`
	Price           decimal.Decimal `json:"price"`
	IsClose         bool            `json:"is_close"`
	IsReduceOnly    bool            `json:"is_reduce_only"`
	IsLiq           bool            `json:"is_liq"`
	TimeInForce     TimeInForce     `json:"tif"`
	Left            int64           `json:"left"`
	FillPrice       decimal.Decimal `json:"fill_price"`
	Text            string          `json:"text"`
	Tkfr            decimal.Decimal `json:"tkfr"`
	Mkfr            decimal.Decimal `json:"mkfr"`
	Refu            int64           `json:"refu"`
	Refr            decimal.Decimal `json:"refr"`
	StopProfitPrice decimal.Decimal `json:"stop_profit_price"`
	StopLossPrice   decimal.Decimal `json:"stop_loss_price"`
	StpID           int64           `json:"stp_id"`
	StpAct          StpAct          `json:"stp_act"`
	BizInfo         string          `json:"biz_info"`
	AmendText       string          `json:"amend_text"`
}

// SubscribeUserTradesService -- futures.usertrades channel (own fills). The
// payload is [user_id, contract].
type SubscribeUserTradesService struct {
	c        *DeliveryWebSocketClient
	userID   string
	contract string
}

func (c *DeliveryWebSocketClient) NewSubscribeUserTradesService(userID, contract string) *SubscribeUserTradesService {
	return &SubscribeUserTradesService{c: c, userID: userID, contract: contract}
}

func (s *SubscribeUserTradesService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryUserTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryUserTrade](ctx, s.c, "futures.usertrades", []string{s.userID, s.contract}, true, cb)
}

// WsDeliveryUserTrade is an own-fill push. ID and OrderID are sent as strings;
// Size is signed; CreateTime is unix seconds, CreateTimeMs unix milliseconds.
type WsDeliveryUserTrade struct {
	ID           string          `json:"id"`
	OrderID      string          `json:"order_id"`
	Contract     string          `json:"contract"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,format:unixmilli"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
	Role         string          `json:"role"`
	Text         string          `json:"text"`
	Fee          decimal.Decimal `json:"fee"`
	PointFee     decimal.Decimal `json:"point_fee"`
}

// SubscribePositionsService -- futures.positions channel (own position updates).
// The payload is [user_id, contract].
type SubscribePositionsService struct {
	c        *DeliveryWebSocketClient
	userID   string
	contract string
}

func (c *DeliveryWebSocketClient) NewSubscribePositionsService(userID, contract string) *SubscribePositionsService {
	return &SubscribePositionsService{c: c, userID: userID, contract: contract}
}

func (s *SubscribePositionsService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryPosition](ctx, s.c, "futures.positions", []string{s.userID, s.contract}, true, cb)
}

// WsDeliveryPosition is an own-position update push. User is the account id sent
// as a string; Size is signed; Time is unix seconds, TimeMs unix milliseconds.
type WsDeliveryPosition struct {
	Contract           string          `json:"contract"`
	User               string          `json:"user"`
	Size               int64           `json:"size"`
	Leverage           decimal.Decimal `json:"leverage"`
	LeverageMax        decimal.Decimal `json:"leverage_max"`
	RiskLimit          decimal.Decimal `json:"risk_limit"`
	CrossLeverageLimit decimal.Decimal `json:"cross_leverage_limit"`
	MaintenanceRate    decimal.Decimal `json:"maintenance_rate"`
	EntryPrice         decimal.Decimal `json:"entry_price"`
	LiqPrice           decimal.Decimal `json:"liq_price"`
	Margin             decimal.Decimal `json:"margin"`
	Mode               string          `json:"mode"`
	RealisedPnL        decimal.Decimal `json:"realised_pnl"`
	RealisedPoint      decimal.Decimal `json:"realised_point"`
	HistoryPnL         decimal.Decimal `json:"history_pnl"`
	HistoryPoint       decimal.Decimal `json:"history_point"`
	LastClosePnL       decimal.Decimal `json:"last_close_pnl"`
	Time               time.Time       `json:"time,format:unix"`
	TimeMs             time.Time       `json:"time_ms,format:unixmilli"`
}

// SubscribeBalancesService -- futures.balances channel (delivery balance
// updates). The payload is [user_id].
type SubscribeBalancesService struct {
	c      *DeliveryWebSocketClient
	userID string
}

func (c *DeliveryWebSocketClient) NewSubscribeBalancesService(userID string) *SubscribeBalancesService {
	return &SubscribeBalancesService{c: c, userID: userID}
}

func (s *SubscribeBalancesService) Do(ctx context.Context, cb WsHandler[[]WsDeliveryBalance]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsDeliveryBalance](ctx, s.c, "futures.balances", []string{s.userID}, true, cb)
}

// WsDeliveryBalance is a delivery balance-change push. User is the account id
// sent as a string; Time is unix seconds, TimeMs unix milliseconds.
type WsDeliveryBalance struct {
	Balance  decimal.Decimal `json:"balance"`
	Change   decimal.Decimal `json:"change"`
	Text     string          `json:"text"`
	Time     time.Time       `json:"time,format:unix"`
	TimeMs   time.Time       `json:"time_ms,format:unixmilli"`
	Type     string          `json:"type"`
	User     string          `json:"user"`
	Currency string          `json:"currency"`
}
