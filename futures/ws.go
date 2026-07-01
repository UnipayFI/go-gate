package futures

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// FuturesWebSocketClient is the stream client for Gate's perpetual-futures
// channels on the wss://fx-ws.gateio.ws/v4/ws/{settle} gateways. Public channels
// need no credentials; private channels (orders, usertrades, positions,
// balances) carry a per-subscription auth object built from WithWebSocketAuth
// and take the numeric user id as the first payload element.
type FuturesWebSocketClient struct {
	*client.WebSocketClient
}

// NewFuturesWebSocketClient constructs a futures stream client for a settlement
// currency: SettleBTC targets the btc gateway, everything else the usdt gateway.
func NewFuturesWebSocketClient(settle Settle, options ...client.WebSocketOptions) *FuturesWebSocketClient {
	url := common.DEFAULT_WS_FUTURES_USDT_URL
	if settle == SettleBTC {
		url = common.DEFAULT_WS_FUTURES_BTC_URL
	}
	return &FuturesWebSocketClient{client.NewWebSocketClient(url, options...)}
}

// WsHandler is invoked for every push (or error) on a subscription. The push's
// Result field is already decoded into T.
type WsHandler[T any] func(*request.WsPush[T], error)

// --- Public channels ---

// SubscribeTickersService -- futures.tickers channel (24h ticker updates). The
// payload is the list of contracts; each push carries an array of tickers.
type SubscribeTickersService struct {
	c         *FuturesWebSocketClient
	contracts []string
}

func (c *FuturesWebSocketClient) NewSubscribeTickersService(contracts ...string) *SubscribeTickersService {
	return &SubscribeTickersService{c: c, contracts: contracts}
}

func (s *SubscribeTickersService) Do(ctx context.Context, cb WsHandler[[]WsFuturesTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFuturesTicker](ctx, s.c, "futures.tickers", s.contracts, false, cb)
}

// WsFuturesTicker is a futures ticker push. Note the volume_24_usd / volume_24_btc
// keys drop the "h" that the REST ticker uses.
type WsFuturesTicker struct {
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
// The payload is [contract]; each push carries an array of trades.
type SubscribeTradesService struct {
	c        *FuturesWebSocketClient
	contract string
}

func (c *FuturesWebSocketClient) NewSubscribeTradesService(contract string) *SubscribeTradesService {
	return &SubscribeTradesService{c: c, contract: contract}
}

func (s *SubscribeTradesService) Do(ctx context.Context, cb WsHandler[[]WsFuturesTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFuturesTrade](ctx, s.c, "futures.trades", []string{s.contract}, false, cb)
}

// WsFuturesTrade is a futures public-trade push. size is signed (negative = taker
// sell); there is no separate side field.
type WsFuturesTrade struct {
	ID           int64           `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,format:unixmilli"`
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
}

// SubscribeBookTickerService -- futures.book_ticker channel (best bid/ask
// updates). The payload is [contract].
type SubscribeBookTickerService struct {
	c        *FuturesWebSocketClient
	contract string
}

func (c *FuturesWebSocketClient) NewSubscribeBookTickerService(contract string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, contract: contract}
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, cb WsHandler[WsFuturesBookTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsFuturesBookTicker](ctx, s.c, "futures.book_ticker", []string{s.contract}, false, cb)
}

// WsFuturesBookTicker is a best bid/ask push. "t" is a millisecond timestamp.
type WsFuturesBookTicker struct {
	Time         time.Time       `json:"t,format:unixmilli"`
	Contract     string          `json:"s"`
	UpdateID     int64           `json:"u"`
	BestBidPrice decimal.Decimal `json:"b"`
	BestBidSize  int64           `json:"B"`
	BestAskPrice decimal.Decimal `json:"a"`
	BestAskSize  int64           `json:"A"`
}

// SubscribeOrderBookService -- futures.order_book channel (limited-depth
// snapshots). level is the depth (e.g. "20"); interval is "0" or "100ms". The
// payload is [contract, level, interval].
type SubscribeOrderBookService struct {
	c        *FuturesWebSocketClient
	contract string
	level    string
	interval string
}

func (c *FuturesWebSocketClient) NewSubscribeOrderBookService(contract, level, interval string) *SubscribeOrderBookService {
	return &SubscribeOrderBookService{c: c, contract: contract, level: level, interval: interval}
}

func (s *SubscribeOrderBookService) Do(ctx context.Context, cb WsHandler[WsFuturesOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsFuturesOrderBook](ctx, s.c, "futures.order_book", []string{s.contract, s.level, s.interval}, false, cb)
}

// WsFuturesOrderBook is a limited-depth order-book snapshot. "t" is a
// millisecond timestamp; rows are {p, s} objects (unlike spot's [price, size]).
type WsFuturesOrderBook struct {
	ID       int64                    `json:"id"`
	Time     time.Time                `json:"t,format:unixmilli"`
	Contract string                   `json:"contract"`
	Asks     []WsFuturesOrderBookItem `json:"asks"`
	Bids     []WsFuturesOrderBookItem `json:"bids"`
}

// WsFuturesOrderBookItem is one price level. Size is signed contract count.
type WsFuturesOrderBookItem struct {
	P decimal.Decimal `json:"p"`
	S int64           `json:"s"`
}

// SubscribeOrderBookUpdateService -- futures.order_book_update channel
// (incremental depth). interval is "20ms", "100ms" or "1000ms"; level is
// optional ("20", "50", "100"). The payload is [contract, interval] or
// [contract, interval, level].
type SubscribeOrderBookUpdateService struct {
	c        *FuturesWebSocketClient
	contract string
	interval string
	level    string
}

func (c *FuturesWebSocketClient) NewSubscribeOrderBookUpdateService(contract, interval string) *SubscribeOrderBookUpdateService {
	return &SubscribeOrderBookUpdateService{c: c, contract: contract, interval: interval}
}

// SetLevel adds the optional depth level to the subscription payload.
func (s *SubscribeOrderBookUpdateService) SetLevel(level string) *SubscribeOrderBookUpdateService {
	s.level = level
	return s
}

func (s *SubscribeOrderBookUpdateService) Do(ctx context.Context, cb WsHandler[WsFuturesOrderBookUpdate]) (chan<- struct{}, <-chan struct{}, error) {
	payload := []string{s.contract, s.interval}
	if s.level != "" {
		payload = append(payload, s.level)
	}
	return request.Subscribe[WsFuturesOrderBookUpdate](ctx, s.c, "futures.order_book_update", payload, false, cb)
}

// WsFuturesOrderBookUpdate is an incremental order-book change. "t" is a
// millisecond timestamp.
type WsFuturesOrderBookUpdate struct {
	Time     time.Time                `json:"t,format:unixmilli"`
	Contract string                   `json:"s"`
	FirstID  int64                    `json:"U"`
	LastID   int64                    `json:"u"`
	Asks     []WsFuturesOrderBookItem `json:"a"`
	Bids     []WsFuturesOrderBookItem `json:"b"`
}

// SubscribeCandlesticksService -- futures.candlesticks channel. interval is e.g.
// "10s", "1m", "1h"; the payload is [interval, contract]. Prefix the contract
// with "mark_" / "index_" / "premium_index_" to stream those series instead.
type SubscribeCandlesticksService struct {
	c        *FuturesWebSocketClient
	interval string
	contract string
}

func (c *FuturesWebSocketClient) NewSubscribeCandlesticksService(interval, contract string) *SubscribeCandlesticksService {
	return &SubscribeCandlesticksService{c: c, interval: interval, contract: contract}
}

func (s *SubscribeCandlesticksService) Do(ctx context.Context, cb WsHandler[WsFuturesCandlestick]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsFuturesCandlestick](ctx, s.c, "futures.candlesticks", []string{s.interval, s.contract}, false, cb)
}

// WsFuturesCandlestick is a candlestick push. Name ("n") is "<interval>_<contract>".
type WsFuturesCandlestick struct {
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
// payload is [userID, contract]; pass "!all" as the contract to receive updates
// for every contract.
type SubscribeOrdersService struct {
	c        *FuturesWebSocketClient
	userID   string
	contract string
}

func (c *FuturesWebSocketClient) NewSubscribeOrdersService(userID, contract string) *SubscribeOrdersService {
	return &SubscribeOrdersService{c: c, userID: userID, contract: contract}
}

func (s *SubscribeOrdersService) Do(ctx context.Context, cb WsHandler[[]WsFuturesOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFuturesOrder](ctx, s.c, "futures.orders", []string{s.userID, s.contract}, true, cb)
}

// WsFuturesOrder is an own-order update push. Note the WS shape differs from the
// REST FuturesOrder: user / stp_id are strings, update_time is milliseconds.
type WsFuturesOrder struct {
	ID                   int64           `json:"id"`
	IDString             string          `json:"id_string"`
	User                 string          `json:"user"`
	Contract             string          `json:"contract"`
	CreateTime           time.Time       `json:"create_time,format:unix"`
	CreateTimeMs         time.Time       `json:"create_time_ms,format:unixmilli"`
	UpdateTime           time.Time       `json:"update_time,format:unixmilli"`
	FinishTime           time.Time       `json:"finish_time,format:unix"`
	FinishTimeMs         time.Time       `json:"finish_time_ms,format:unixmilli"`
	FinishAs             FinishAs        `json:"finish_as"`
	Status               OrderStatus     `json:"status"`
	Size                 int64           `json:"size"`
	Iceberg              int64           `json:"iceberg"`
	Price                decimal.Decimal `json:"price"`
	IsClose              bool            `json:"is_close"`
	IsReduceOnly         bool            `json:"is_reduce_only"`
	IsLiq                bool            `json:"is_liq"`
	IsVoucher            bool            `json:"is_voucher"`
	Tif                  TimeInForce     `json:"tif"`
	Left                 int64           `json:"left"`
	FillPrice            decimal.Decimal `json:"fill_price"`
	Text                 string          `json:"text"`
	Tkfr                 decimal.Decimal `json:"tkfr"`
	Mkfr                 decimal.Decimal `json:"mkfr"`
	Fee                  decimal.Decimal `json:"fee"`
	PointFee             decimal.Decimal `json:"point_fee"`
	Refu                 int64           `json:"refu"`
	Refr                 decimal.Decimal `json:"refr"`
	StopProfitPrice      decimal.Decimal `json:"stop_profit_price"`
	StopLossPrice        decimal.Decimal `json:"stop_loss_price"`
	StpID                string          `json:"stp_id"`
	StpAct               StpAct          `json:"stp_act"`
	BizInfo              string          `json:"biz_info"`
	AmendText            string          `json:"amend_text"`
	UpdateID             int64           `json:"update_id"`
	Role                 string          `json:"role"`
	BBO                  string          `json:"bbo"`
	MarketOrderSlipRatio decimal.Decimal `json:"market_order_slip_ratio"`
	PosMarginMode        string          `json:"pos_margin_mode"`
	Leverage             decimal.Decimal `json:"leverage"`
}

// SubscribeUserTradesService -- futures.usertrades channel (own fills). The
// payload is [userID, contract].
type SubscribeUserTradesService struct {
	c        *FuturesWebSocketClient
	userID   string
	contract string
}

func (c *FuturesWebSocketClient) NewSubscribeUserTradesService(userID, contract string) *SubscribeUserTradesService {
	return &SubscribeUserTradesService{c: c, userID: userID, contract: contract}
}

func (s *SubscribeUserTradesService) Do(ctx context.Context, cb WsHandler[[]WsFuturesUserTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFuturesUserTrade](ctx, s.c, "futures.usertrades", []string{s.userID, s.contract}, true, cb)
}

// WsFuturesUserTrade is an own-fill push.
type WsFuturesUserTrade struct {
	ID           string          `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,format:unixmilli"`
	Contract     string          `json:"contract"`
	OrderID      string          `json:"order_id"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
	Role         string          `json:"role"`
	Text         string          `json:"text"`
	Fee          decimal.Decimal `json:"fee"`
	PointFee     decimal.Decimal `json:"point_fee"`
}

// SubscribePositionsService -- futures.positions channel (own position updates).
// The payload is [userID, contract].
type SubscribePositionsService struct {
	c        *FuturesWebSocketClient
	userID   string
	contract string
}

func (c *FuturesWebSocketClient) NewSubscribePositionsService(userID, contract string) *SubscribePositionsService {
	return &SubscribePositionsService{c: c, userID: userID, contract: contract}
}

func (s *SubscribePositionsService) Do(ctx context.Context, cb WsHandler[[]WsFuturesPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFuturesPosition](ctx, s.c, "futures.positions", []string{s.userID, s.contract}, true, cb)
}

// WsFuturesPosition is an own-position update push.
type WsFuturesPosition struct {
	Contract           string          `json:"contract"`
	User               string          `json:"user"`
	Size               int64           `json:"size"`
	EntryPrice         decimal.Decimal `json:"entry_price"`
	Leverage           decimal.Decimal `json:"leverage"`
	LeverageMax        decimal.Decimal `json:"leverage_max"`
	CrossLeverageLimit decimal.Decimal `json:"cross_leverage_limit"`
	LiqPrice           decimal.Decimal `json:"liq_price"`
	MaintenanceRate    decimal.Decimal `json:"maintenance_rate"`
	Margin             decimal.Decimal `json:"margin"`
	Mode               string          `json:"mode"`
	RiskLimit          decimal.Decimal `json:"risk_limit"`
	RealisedPnl        decimal.Decimal `json:"realised_pnl"`
	RealisedPoint      decimal.Decimal `json:"realised_point"`
	HistoryPnl         decimal.Decimal `json:"history_pnl"`
	HistoryPoint       decimal.Decimal `json:"history_point"`
	LastClosePnl       decimal.Decimal `json:"last_close_pnl"`
	Time               time.Time       `json:"time,format:unix"`
	TimeMs             time.Time       `json:"time_ms,format:unixmilli"`
}

// SubscribeBalancesService -- futures.balances channel (own balance updates).
// The payload is [userID].
type SubscribeBalancesService struct {
	c      *FuturesWebSocketClient
	userID string
}

func (c *FuturesWebSocketClient) NewSubscribeBalancesService(userID string) *SubscribeBalancesService {
	return &SubscribeBalancesService{c: c, userID: userID}
}

func (s *SubscribeBalancesService) Do(ctx context.Context, cb WsHandler[[]WsFuturesBalance]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFuturesBalance](ctx, s.c, "futures.balances", []string{s.userID}, true, cb)
}

// WsFuturesBalance is an own balance-change push.
type WsFuturesBalance struct {
	Balance  decimal.Decimal `json:"balance"`
	Change   decimal.Decimal `json:"change"`
	Text     string          `json:"text"`
	Type     string          `json:"type"`
	Currency string          `json:"currency"`
	User     string          `json:"user"`
	Time     time.Time       `json:"time,format:unix"`
	TimeMs   time.Time       `json:"time_ms,format:unixmilli"`
}
