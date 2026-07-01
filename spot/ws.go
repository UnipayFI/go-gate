package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/common"
	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// SpotWebSocketClient is the stream client for Gate's spot / margin / unified
// channels on the wss://api.gateio.ws/ws/v4/ gateway. Public channels need no
// credentials; private channels (orders, usertrades, balances) carry a per-
// subscription auth object built from WithWebSocketAuth.
type SpotWebSocketClient struct {
	*client.WebSocketClient
}

// NewSpotWebSocketClient constructs a spot stream client.
func NewSpotWebSocketClient(options ...client.WebSocketOptions) *SpotWebSocketClient {
	return &SpotWebSocketClient{client.NewWebSocketClient(common.DEFAULT_WS_SPOT_URL, options...)}
}

// WsHandler is invoked for every push (or error) on a subscription. The push's
// Result field is already decoded into T.
type WsHandler[T any] func(*request.WsPush[T], error)

// --- Public channels ---

// SubscribeTickersService -- spot.tickers channel (24h ticker updates).
type SubscribeTickersService struct {
	c     *SpotWebSocketClient
	pairs []string
}

func (c *SpotWebSocketClient) NewSubscribeTickersService(pairs ...string) *SubscribeTickersService {
	return &SubscribeTickersService{c: c, pairs: pairs}
}

func (s *SubscribeTickersService) Do(ctx context.Context, cb WsHandler[WsTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsTicker](ctx, s.c, "spot.tickers", s.pairs, false, cb)
}

// WsTicker is a spot ticker push.
type WsTicker struct {
	CurrencyPair     string          `json:"currency_pair"`
	Last             decimal.Decimal `json:"last"`
	LowestAsk        decimal.Decimal `json:"lowest_ask"`
	HighestBid       decimal.Decimal `json:"highest_bid"`
	ChangePercentage decimal.Decimal `json:"change_percentage"`
	BaseVolume       decimal.Decimal `json:"base_volume"`
	QuoteVolume      decimal.Decimal `json:"quote_volume"`
	High24h          decimal.Decimal `json:"high_24h"`
	Low24h           decimal.Decimal `json:"low_24h"`
}

// SubscribeTradesService -- spot.trades channel (public tick-by-tick fills).
type SubscribeTradesService struct {
	c     *SpotWebSocketClient
	pairs []string
}

func (c *SpotWebSocketClient) NewSubscribeTradesService(pairs ...string) *SubscribeTradesService {
	return &SubscribeTradesService{c: c, pairs: pairs}
}

func (s *SubscribeTradesService) Do(ctx context.Context, cb WsHandler[WsPublicTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsPublicTrade](ctx, s.c, "spot.trades", s.pairs, false, cb)
}

// WsPublicTrade is a spot public-trade push.
type WsPublicTrade struct {
	ID           int64           `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,string,format:unixmilli"`
	Side         Side            `json:"side"`
	CurrencyPair string          `json:"currency_pair"`
	Amount       decimal.Decimal `json:"amount"`
	Price        decimal.Decimal `json:"price"`
}

// SubscribeBookTickerService -- spot.book_ticker channel (best bid/ask updates).
type SubscribeBookTickerService struct {
	c     *SpotWebSocketClient
	pairs []string
}

func (c *SpotWebSocketClient) NewSubscribeBookTickerService(pairs ...string) *SubscribeBookTickerService {
	return &SubscribeBookTickerService{c: c, pairs: pairs}
}

func (s *SubscribeBookTickerService) Do(ctx context.Context, cb WsHandler[WsBookTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsBookTicker](ctx, s.c, "spot.book_ticker", s.pairs, false, cb)
}

// WsBookTicker is a best bid/ask push.
type WsBookTicker struct {
	Time         time.Time       `json:"t,format:unixmilli"`
	LastID       int64           `json:"u"`
	CurrencyPair string          `json:"s"`
	BestBid      decimal.Decimal `json:"b"`
	BestBidSize  decimal.Decimal `json:"B"`
	BestAsk      decimal.Decimal `json:"a"`
	BestAskSize  decimal.Decimal `json:"A"`
}

// SubscribeCandlesticksService -- spot.candlesticks channel. interval is e.g.
// "10s", "1m", "1h"; the payload is [interval, currency_pair].
type SubscribeCandlesticksService struct {
	c        *SpotWebSocketClient
	interval string
	pair     string
}

func (c *SpotWebSocketClient) NewSubscribeCandlesticksService(interval, pair string) *SubscribeCandlesticksService {
	return &SubscribeCandlesticksService{c: c, interval: interval, pair: pair}
}

func (s *SubscribeCandlesticksService) Do(ctx context.Context, cb WsHandler[WsCandlestick]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsCandlestick](ctx, s.c, "spot.candlesticks", []string{s.interval, s.pair}, false, cb)
}

// WsCandlestick is a candlestick push. Name ("n") is "<interval>_<pair>".
type WsCandlestick struct {
	Time        time.Time       `json:"t,string,format:unix"`
	Volume      decimal.Decimal `json:"v"`
	Close       decimal.Decimal `json:"c"`
	High        decimal.Decimal `json:"h"`
	Low         decimal.Decimal `json:"l"`
	Open        decimal.Decimal `json:"o"`
	Name        string          `json:"n"`
	Amount      decimal.Decimal `json:"a"`
	WindowClose bool            `json:"w"`
}

// SubscribeOrderBookUpdateService -- spot.order_book_update channel (incremental
// depth). interval is "100ms" or "1000ms".
type SubscribeOrderBookUpdateService struct {
	c        *SpotWebSocketClient
	pair     string
	interval string
}

func (c *SpotWebSocketClient) NewSubscribeOrderBookUpdateService(pair, interval string) *SubscribeOrderBookUpdateService {
	return &SubscribeOrderBookUpdateService{c: c, pair: pair, interval: interval}
}

func (s *SubscribeOrderBookUpdateService) Do(ctx context.Context, cb WsHandler[WsDepthUpdate]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsDepthUpdate](ctx, s.c, "spot.order_book_update", []string{s.pair, s.interval}, false, cb)
}

// WsDepthUpdate is an incremental order-book change.
type WsDepthUpdate struct {
	Time         time.Time           `json:"t,format:unixmilli"`
	Event        string              `json:"e"`
	EventTime    time.Time           `json:"E,format:unix"`
	CurrencyPair string              `json:"s"`
	FirstID      int64               `json:"U"`
	LastID       int64               `json:"u"`
	Bids         [][]decimal.Decimal `json:"b"`
	Asks         [][]decimal.Decimal `json:"a"`
}

// SubscribeOrderBookService -- spot.order_book channel (limited-depth snapshots).
// level is "5", "10", "20", "50", "100"; interval is "100ms" or "1000ms".
type SubscribeOrderBookService struct {
	c        *SpotWebSocketClient
	pair     string
	level    string
	interval string
}

func (c *SpotWebSocketClient) NewSubscribeOrderBookService(pair, level, interval string) *SubscribeOrderBookService {
	return &SubscribeOrderBookService{c: c, pair: pair, level: level, interval: interval}
}

func (s *SubscribeOrderBookService) Do(ctx context.Context, cb WsHandler[WsDepthSnapshot]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[WsDepthSnapshot](ctx, s.c, "spot.order_book", []string{s.pair, s.level, s.interval}, false, cb)
}

// WsDepthSnapshot is a limited-depth order-book snapshot.
type WsDepthSnapshot struct {
	Time         time.Time           `json:"t,format:unixmilli"`
	LastUpdateID int64               `json:"lastUpdateId"`
	CurrencyPair string              `json:"s"`
	Bids         [][]decimal.Decimal `json:"bids"`
	Asks         [][]decimal.Decimal `json:"asks"`
}

// --- Private channels (require WithWebSocketAuth) ---

// SubscribeOrdersService -- spot.orders channel (own order updates). Pass "!all"
// to receive updates for every pair.
type SubscribeOrdersService struct {
	c     *SpotWebSocketClient
	pairs []string
}

func (c *SpotWebSocketClient) NewSubscribeOrdersService(pairs ...string) *SubscribeOrdersService {
	return &SubscribeOrdersService{c: c, pairs: pairs}
}

func (s *SubscribeOrdersService) Do(ctx context.Context, cb WsHandler[[]WsOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrder](ctx, s.c, "spot.orders", s.pairs, true, cb)
}

// WsOrder is an own-order update push. Event is "put", "update" or "finish".
type WsOrder struct {
	ID                 string          `json:"id"`
	Text               string          `json:"text"`
	CreateTime         time.Time       `json:"create_time,string,format:unix"`
	UpdateTime         time.Time       `json:"update_time,string,format:unix"`
	CreateTimeMs       time.Time       `json:"create_time_ms,string,format:unixmilli"`
	UpdateTimeMs       time.Time       `json:"update_time_ms,string,format:unixmilli"`
	CurrencyPair       string          `json:"currency_pair"`
	Type               OrderType       `json:"type"`
	Account            Account         `json:"account"`
	Side               Side            `json:"side"`
	Amount             decimal.Decimal `json:"amount"`
	Price              decimal.Decimal `json:"price"`
	TimeInForce        TimeInForce     `json:"time_in_force"`
	Left               decimal.Decimal `json:"left"`
	FilledTotal        decimal.Decimal `json:"filled_total"`
	AvgDealPrice       decimal.Decimal `json:"avg_deal_price"`
	Fee                decimal.Decimal `json:"fee"`
	FeeCurrency        string          `json:"fee_currency"`
	PointFee           decimal.Decimal `json:"point_fee"`
	GtFee              decimal.Decimal `json:"gt_fee"`
	GtDiscount         bool            `json:"gt_discount"`
	RebatedFee         decimal.Decimal `json:"rebated_fee"`
	RebatedFeeCurrency string          `json:"rebated_fee_currency"`
	StpID              int64           `json:"stp_id"`
	StpAct             StpAct          `json:"stp_act"`
	FinishAs           string          `json:"finish_as"`
	BizInfo            string          `json:"biz_info"`
	AmendText          string          `json:"amend_text"`
	User               int64           `json:"user"`
	Event              string          `json:"event"`
}

// SubscribeUserTradesService -- spot.usertrades channel (own fills).
type SubscribeUserTradesService struct {
	c     *SpotWebSocketClient
	pairs []string
}

func (c *SpotWebSocketClient) NewSubscribeUserTradesService(pairs ...string) *SubscribeUserTradesService {
	return &SubscribeUserTradesService{c: c, pairs: pairs}
}

func (s *SubscribeUserTradesService) Do(ctx context.Context, cb WsHandler[[]WsUserTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsUserTrade](ctx, s.c, "spot.usertrades", s.pairs, true, cb)
}

// WsUserTrade is an own-fill push.
type WsUserTrade struct {
	ID           int64           `json:"id"`
	UserID       int64           `json:"user_id"`
	OrderID      string          `json:"order_id"`
	CurrencyPair string          `json:"currency_pair"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,string,format:unixmilli"`
	Side         Side            `json:"side"`
	Amount       decimal.Decimal `json:"amount"`
	Role         string          `json:"role"`
	Price        decimal.Decimal `json:"price"`
	Fee          decimal.Decimal `json:"fee"`
	FeeCurrency  string          `json:"fee_currency"`
	PointFee     decimal.Decimal `json:"point_fee"`
	GtFee        decimal.Decimal `json:"gt_fee"`
	Text         string          `json:"text"`
	AmendText    string          `json:"amend_text"`
	BizInfo      string          `json:"biz_info"`
}

// SubscribeBalancesService -- spot.balances channel (spot balance updates).
type SubscribeBalancesService struct {
	c *SpotWebSocketClient
}

func (c *SpotWebSocketClient) NewSubscribeBalancesService() *SubscribeBalancesService {
	return &SubscribeBalancesService{c: c}
}

func (s *SubscribeBalancesService) Do(ctx context.Context, cb WsHandler[[]WsBalance]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsBalance](ctx, s.c, "spot.balances", []string{}, true, cb)
}

// WsBalance is a spot balance-change push.
type WsBalance struct {
	Timestamp    time.Time       `json:"timestamp,string,format:unix"`
	TimestampMs  time.Time       `json:"timestamp_ms,string,format:unixmilli"`
	User         string          `json:"user"`
	Currency     string          `json:"currency"`
	Change       decimal.Decimal `json:"change"`
	Total        decimal.Decimal `json:"total"`
	Available    decimal.Decimal `json:"available"`
	Freeze       decimal.Decimal `json:"freeze"`
	FreezeChange decimal.Decimal `json:"freeze_change"`
	ChangeType   string          `json:"change_type"`
}
