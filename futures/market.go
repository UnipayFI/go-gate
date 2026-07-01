package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListFuturesContractsService -- GET /api/v4/futures/{settle}/contracts
//
// Returns every perpetual-futures contract and its trading rules for a
// settlement currency.
type ListFuturesContractsService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesContractsService(settle Settle) *ListFuturesContractsService {
	return &ListFuturesContractsService{c: c, settle: settle, params: map[string]string{}}
}

// SetLimit caps the number of contracts returned in one page.
func (s *ListFuturesContractsService) SetLimit(limit int) *ListFuturesContractsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the first offset contracts (pagination, starting from 0).
func (s *ListFuturesContractsService) SetOffset(offset int) *ListFuturesContractsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListFuturesContractsService) Do(ctx context.Context) ([]FuturesContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/contracts", s.params)
	resp, err := request.Do[[]FuturesContract](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetFuturesContractService -- GET /api/v4/futures/{settle}/contracts/{contract}
//
// Returns the trading rules for a single perpetual-futures contract.
type GetFuturesContractService struct {
	c        *FuturesClient
	settle   Settle
	contract string
}

func (c *FuturesClient) NewGetFuturesContractService(settle Settle, contract string) *GetFuturesContractService {
	return &GetFuturesContractService{c: c, settle: settle, contract: contract}
}

func (s *GetFuturesContractService) Do(ctx context.Context) (*FuturesContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/contracts/"+s.contract)
	return request.Do[FuturesContract](req)
}

// FuturesContract is a perpetual-futures contract and its trading rules.
type FuturesContract struct {
	Name                  string          `json:"name"`
	Type                  string          `json:"type"`
	QuantoMultiplier      decimal.Decimal `json:"quanto_multiplier"`
	LeverageMin           decimal.Decimal `json:"leverage_min"`
	LeverageMax           decimal.Decimal `json:"leverage_max"`
	CrossLeverageDefault  decimal.Decimal `json:"cross_leverage_default"`
	MaintenanceRate       decimal.Decimal `json:"maintenance_rate"`
	MarkType              string          `json:"mark_type"`
	MarkPrice             decimal.Decimal `json:"mark_price"`
	IndexPrice            decimal.Decimal `json:"index_price"`
	LastPrice             decimal.Decimal `json:"last_price"`
	MakerFeeRate          decimal.Decimal `json:"maker_fee_rate"`
	TakerFeeRate          decimal.Decimal `json:"taker_fee_rate"`
	OrderPriceRound       decimal.Decimal `json:"order_price_round"`
	MarkPriceRound        decimal.Decimal `json:"mark_price_round"`
	FundingRate           decimal.Decimal `json:"funding_rate"`
	FundingInterval       int             `json:"funding_interval"`
	FundingNextApply      time.Time       `json:"funding_next_apply,format:unix"`
	FundingOffset         int             `json:"funding_offset"`
	FundingRateIndicative decimal.Decimal `json:"funding_rate_indicative"`
	FundingImpactValue    decimal.Decimal `json:"funding_impact_value"`
	FundingCapRatio       decimal.Decimal `json:"funding_cap_ratio"`
	FundingRateLimit      decimal.Decimal `json:"funding_rate_limit"`
	InterestRate          decimal.Decimal `json:"interest_rate"`
	RiskLimitBase         decimal.Decimal `json:"risk_limit_base"`
	RiskLimitStep         decimal.Decimal `json:"risk_limit_step"`
	RiskLimitMax          decimal.Decimal `json:"risk_limit_max"`
	OrderSizeMin          int64           `json:"order_size_min"`
	OrderSizeMax          int64           `json:"order_size_max"`
	OrderPriceDeviate     decimal.Decimal `json:"order_price_deviate"`
	OrdersLimit           int64           `json:"orders_limit"`
	RefDiscountRate       decimal.Decimal `json:"ref_discount_rate"`
	RefRebateRate         decimal.Decimal `json:"ref_rebate_rate"`
	OrderbookID           int64           `json:"orderbook_id"`
	TradeID               int64           `json:"trade_id"`
	TradeSize             int64           `json:"trade_size"`
	PositionSize          int64           `json:"position_size"`
	LongUsers             int64           `json:"long_users"`
	ShortUsers            int64           `json:"short_users"`
	ConfigChangeTime      time.Time       `json:"config_change_time,format:unix"`
	CreateTime            time.Time       `json:"create_time,format:unix"`
	LaunchTime            time.Time       `json:"launch_time,format:unix"`
	DelistingTime         time.Time       `json:"delisting_time,format:unix"`
	DelistedTime          time.Time       `json:"delisted_time,format:unix"`
	InDelisting           bool            `json:"in_delisting"`
	IsPreMarket           bool            `json:"is_pre_market"`
	EnableBonus           bool            `json:"enable_bonus"`
	EnableCredit          bool            `json:"enable_credit"`
	EnableDecimal         bool            `json:"enable_decimal"`
	EnableCircuitBreaker  bool            `json:"enable_circuit_breaker"`
	VoucherLeverage       decimal.Decimal `json:"voucher_leverage"`
	MarketOrderSlipRatio  decimal.Decimal `json:"market_order_slip_ratio"`
	MarketOrderSizeMax    decimal.Decimal `json:"market_order_size_max"`
	ContractType          string          `json:"contract_type"`
	Status                string          `json:"status"`
}

// ListFuturesOrderBookService -- GET /api/v4/futures/{settle}/order_book
//
// Returns the current bid/ask depth for a contract. Bids are sorted by price
// high-to-low, asks low-to-high.
type ListFuturesOrderBookService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesOrderBookService(settle Settle, contract string) *ListFuturesOrderBookService {
	return &ListFuturesOrderBookService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetInterval aggregates price levels by the given tick (e.g. "0", "0.1").
func (s *ListFuturesOrderBookService) SetInterval(interval string) *ListFuturesOrderBookService {
	s.params["interval"] = interval
	return s
}

// SetLimit caps the number of levels returned per side.
func (s *ListFuturesOrderBookService) SetLimit(limit int) *ListFuturesOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetWithID includes the order-book update id in the response when true.
func (s *ListFuturesOrderBookService) SetWithID(withID bool) *ListFuturesOrderBookService {
	s.params["with_id"] = strconv.FormatBool(withID)
	return s
}

func (s *ListFuturesOrderBookService) Do(ctx context.Context) (*FuturesOrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/order_book", s.params)
	return request.Do[FuturesOrderBook](req)
}

// FuturesOrderBook is a snapshot of the futures market depth.
type FuturesOrderBook struct {
	ID      int64                   `json:"id"`
	Current time.Time               `json:"current,format:unix"`
	Update  time.Time               `json:"update,format:unix"`
	Asks    []FuturesOrderBookEntry `json:"asks"`
	Bids    []FuturesOrderBookEntry `json:"bids"`
}

// FuturesOrderBookEntry is one price level in the futures order book. Unlike
// spot (which sends [price, size] arrays), futures rows are objects.
type FuturesOrderBookEntry struct {
	P decimal.Decimal `json:"p"`
	S int64           `json:"s"`
}

// ListFuturesTradesService -- GET /api/v4/futures/{settle}/trades
//
// Returns recent public trades for a contract.
type ListFuturesTradesService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesTradesService(settle Settle, contract string) *ListFuturesTradesService {
	return &ListFuturesTradesService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetLimit caps the number of trades returned.
func (s *ListFuturesTradesService) SetLimit(limit int) *ListFuturesTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the first offset trades (pagination, starting from 0).
func (s *ListFuturesTradesService) SetOffset(offset int) *ListFuturesTradesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetLastID returns trades after the given fill id (cursor pagination).
func (s *ListFuturesTradesService) SetLastID(lastID string) *ListFuturesTradesService {
	s.params["last_id"] = lastID
	return s
}

// SetFrom limits results to trades at or after t (unix seconds).
func (s *ListFuturesTradesService) SetFrom(t time.Time) *ListFuturesTradesService {
	s.params["from"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

// SetTo limits results to trades at or before t (unix seconds).
func (s *ListFuturesTradesService) SetTo(t time.Time) *ListFuturesTradesService {
	s.params["to"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

func (s *ListFuturesTradesService) Do(ctx context.Context) ([]FuturesTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/trades", s.params)
	resp, err := request.Do[[]FuturesTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesTrade is a single public trade. size is signed (negative = taker sell);
// there is no separate side field. create_time and create_time_ms are both unix
// seconds — create_time_ms simply carries millisecond precision (3 decimals).
type FuturesTrade struct {
	ID           int64           `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,format:unix"`
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
	IsInternal   bool            `json:"is_internal"`
}

// ListFuturesCandlesticksService -- GET /api/v4/futures/{settle}/candlesticks
//
// Returns OHLC candlesticks for a contract. Prefix contract with "mark_" for
// mark-price or "index_" for index-price candlesticks.
type ListFuturesCandlesticksService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesCandlesticksService(settle Settle, contract string) *ListFuturesCandlesticksService {
	return &ListFuturesCandlesticksService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetFrom sets the start of the candlestick window (unix seconds).
func (s *ListFuturesCandlesticksService) SetFrom(t time.Time) *ListFuturesCandlesticksService {
	s.params["from"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

// SetTo sets the end of the candlestick window (unix seconds).
func (s *ListFuturesCandlesticksService) SetTo(t time.Time) *ListFuturesCandlesticksService {
	s.params["to"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

// SetLimit caps the number of candlesticks returned (max 2000).
func (s *ListFuturesCandlesticksService) SetLimit(limit int) *ListFuturesCandlesticksService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetInterval selects the candlestick period (e.g. "1m", "1h", "1d").
func (s *ListFuturesCandlesticksService) SetInterval(interval string) *ListFuturesCandlesticksService {
	s.params["interval"] = interval
	return s
}

// SetTimezone selects the timezone used to align daily candlesticks.
func (s *ListFuturesCandlesticksService) SetTimezone(timezone string) *ListFuturesCandlesticksService {
	s.params["timezone"] = timezone
	return s
}

func (s *ListFuturesCandlesticksService) Do(ctx context.Context) ([]FuturesCandlestick, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/candlesticks", s.params)
	resp, err := request.Do[[]FuturesCandlestick](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesCandlestick is a single OHLC data point. Futures candlesticks are
// objects (not the positional arrays spot uses).
type FuturesCandlestick struct {
	T   time.Time       `json:"t,format:unix"`
	V   int64           `json:"v"`
	O   decimal.Decimal `json:"o"`
	H   decimal.Decimal `json:"h"`
	L   decimal.Decimal `json:"l"`
	C   decimal.Decimal `json:"c"`
	Sum decimal.Decimal `json:"sum"`
}

// ListFuturesPremiumIndexService -- GET /api/v4/futures/{settle}/premium_index
//
// Returns premium-index OHLC candlesticks for a contract.
type ListFuturesPremiumIndexService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesPremiumIndexService(settle Settle, contract string) *ListFuturesPremiumIndexService {
	return &ListFuturesPremiumIndexService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetFrom sets the start of the premium-index window (unix seconds).
func (s *ListFuturesPremiumIndexService) SetFrom(t time.Time) *ListFuturesPremiumIndexService {
	s.params["from"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

// SetTo sets the end of the premium-index window (unix seconds).
func (s *ListFuturesPremiumIndexService) SetTo(t time.Time) *ListFuturesPremiumIndexService {
	s.params["to"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

// SetLimit caps the number of points returned (max 1000).
func (s *ListFuturesPremiumIndexService) SetLimit(limit int) *ListFuturesPremiumIndexService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetInterval selects the point period (e.g. "1m", "1h", "1d").
func (s *ListFuturesPremiumIndexService) SetInterval(interval string) *ListFuturesPremiumIndexService {
	s.params["interval"] = interval
	return s
}

func (s *ListFuturesPremiumIndexService) Do(ctx context.Context) ([]FuturesPremiumIndex, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/premium_index", s.params)
	resp, err := request.Do[[]FuturesPremiumIndex](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesPremiumIndex is a single premium-index OHLC data point.
type FuturesPremiumIndex struct {
	T time.Time       `json:"t,format:unix"`
	O decimal.Decimal `json:"o"`
	H decimal.Decimal `json:"h"`
	L decimal.Decimal `json:"l"`
	C decimal.Decimal `json:"c"`
}

// ListFuturesTickersService -- GET /api/v4/futures/{settle}/tickers
//
// Returns 24h ticker statistics for all contracts, or one when contract is set.
type ListFuturesTickersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesTickersService(settle Settle) *ListFuturesTickersService {
	return &ListFuturesTickersService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single contract.
func (s *ListFuturesTickersService) SetContract(contract string) *ListFuturesTickersService {
	s.params["contract"] = contract
	return s
}

func (s *ListFuturesTickersService) Do(ctx context.Context) ([]FuturesTicker, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/tickers", s.params)
	resp, err := request.Do[[]FuturesTicker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesTicker is 24h rolling market statistics for a futures contract.
type FuturesTicker struct {
	Contract              string          `json:"contract"`
	Last                  decimal.Decimal `json:"last"`
	ChangePercentage      decimal.Decimal `json:"change_percentage"`
	ChangePrice           decimal.Decimal `json:"change_price"`
	ChangeUTC0            decimal.Decimal `json:"change_utc0"`
	ChangeUTC8            decimal.Decimal `json:"change_utc8"`
	ChangeUTC0Price       decimal.Decimal `json:"change_utc0_price"`
	ChangeUTC8Price       decimal.Decimal `json:"change_utc8_price"`
	TotalSize             decimal.Decimal `json:"total_size"`
	Low24h                decimal.Decimal `json:"low_24h"`
	High24h               decimal.Decimal `json:"high_24h"`
	Volume24h             decimal.Decimal `json:"volume_24h"`
	Volume24hBTC          decimal.Decimal `json:"volume_24h_btc"`
	Volume24hUSD          decimal.Decimal `json:"volume_24h_usd"`
	Volume24hBase         decimal.Decimal `json:"volume_24h_base"`
	Volume24hQuote        decimal.Decimal `json:"volume_24h_quote"`
	Volume24hSettle       decimal.Decimal `json:"volume_24h_settle"`
	MarkPrice             decimal.Decimal `json:"mark_price"`
	FundingRate           decimal.Decimal `json:"funding_rate"`
	FundingRateIndicative decimal.Decimal `json:"funding_rate_indicative"`
	IndexPrice            decimal.Decimal `json:"index_price"`
	QuantoBaseRate        decimal.Decimal `json:"quanto_base_rate"`
	QuantoMultiplier      decimal.Decimal `json:"quanto_multiplier"`
	LowestAsk             decimal.Decimal `json:"lowest_ask"`
	LowestSize            decimal.Decimal `json:"lowest_size"`
	HighestBid            decimal.Decimal `json:"highest_bid"`
	HighestSize           decimal.Decimal `json:"highest_size"`
}

// ListFuturesFundingRateHistoryService -- GET /api/v4/futures/{settle}/funding_rate
//
// Returns the historical funding rate for a contract.
type ListFuturesFundingRateHistoryService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesFundingRateHistoryService(settle Settle, contract string) *ListFuturesFundingRateHistoryService {
	return &ListFuturesFundingRateHistoryService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetLimit caps the number of funding-rate records returned.
func (s *ListFuturesFundingRateHistoryService) SetLimit(limit int) *ListFuturesFundingRateHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom limits results to records at or after t (unix seconds).
func (s *ListFuturesFundingRateHistoryService) SetFrom(t time.Time) *ListFuturesFundingRateHistoryService {
	s.params["from"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

// SetTo limits results to records at or before t (unix seconds).
func (s *ListFuturesFundingRateHistoryService) SetTo(t time.Time) *ListFuturesFundingRateHistoryService {
	s.params["to"] = strconv.FormatInt(t.Unix(), 10)
	return s
}

func (s *ListFuturesFundingRateHistoryService) Do(ctx context.Context) ([]FuturesFundingRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/funding_rate", s.params)
	resp, err := request.Do[[]FuturesFundingRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesFundingRate is one historical funding-rate observation.
type FuturesFundingRate struct {
	T time.Time       `json:"t,format:unix"`
	R decimal.Decimal `json:"r"`
}
