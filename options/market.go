package options

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListOptionsUnderlyingsService -- GET /api/v4/options/underlyings
//
// Returns every options underlying and its current index price.
type ListOptionsUnderlyingsService struct {
	c *OptionsClient
}

func (c *OptionsClient) NewListOptionsUnderlyingsService() *ListOptionsUnderlyingsService {
	return &ListOptionsUnderlyingsService{c: c}
}

func (s *ListOptionsUnderlyingsService) Do(ctx context.Context) ([]OptionsUnderlying, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/underlyings")
	resp, err := request.Do[[]OptionsUnderlying](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OptionsUnderlying is an options underlying and its spot index price.
type OptionsUnderlying struct {
	Name       string          `json:"name"`
	IndexTime  time.Time       `json:"index_time,format:unix"`
	IndexPrice decimal.Decimal `json:"index_price"`
}

// ListOptionsExpirationsService -- GET /api/v4/options/expirations
//
// Returns the list of expiration timestamps (unix seconds) for an underlying.
type ListOptionsExpirationsService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsExpirationsService(underlying string) *ListOptionsExpirationsService {
	return &ListOptionsExpirationsService{c: c, params: map[string]string{"underlying": underlying}}
}

func (s *ListOptionsExpirationsService) Do(ctx context.Context) ([]int64, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/expirations", s.params)
	resp, err := request.Do[[]int64](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListOptionsContractsService -- GET /api/v4/options/contracts
//
// Returns all option contracts for an underlying, optionally narrowed to a
// single expiration date.
type ListOptionsContractsService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsContractsService(underlying string) *ListOptionsContractsService {
	return &ListOptionsContractsService{c: c, params: map[string]string{"underlying": underlying}}
}

// SetExpiration narrows the result to contracts expiring at the given time.
func (s *ListOptionsContractsService) SetExpiration(expiration time.Time) *ListOptionsContractsService {
	s.params["expiration"] = strconv.FormatInt(expiration.Unix(), 10)
	return s
}

func (s *ListOptionsContractsService) Do(ctx context.Context) ([]OptionsContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/contracts", s.params)
	resp, err := request.Do[[]OptionsContract](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOptionsContractService -- GET /api/v4/options/contracts/{contract}
//
// Returns the full detail of a single option contract.
type GetOptionsContractService struct {
	c        *OptionsClient
	contract string
}

func (c *OptionsClient) NewGetOptionsContractService(contract string) *GetOptionsContractService {
	return &GetOptionsContractService{c: c, contract: contract}
}

func (s *GetOptionsContractService) Do(ctx context.Context) (*OptionsContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/contracts/"+s.contract)
	return request.Do[OptionsContract](req)
}

// OptionsContract is an option contract and its trading rules / market snapshot.
type OptionsContract struct {
	Name                 string          `json:"name"`
	Tag                  string          `json:"tag"`
	IsActive             bool            `json:"is_active"`
	IsCall               bool            `json:"is_call"`
	Underlying           string          `json:"underlying"`
	UnderlyingPrice      decimal.Decimal `json:"underlying_price"`
	StrikePrice          decimal.Decimal `json:"strike_price"`
	CreateTime           time.Time       `json:"create_time,format:unix"`
	ExpirationTime       time.Time       `json:"expiration_time,format:unix"`
	Multiplier           decimal.Decimal `json:"multiplier"`
	LastPrice            decimal.Decimal `json:"last_price"`
	MarkPrice            decimal.Decimal `json:"mark_price"`
	MarkPriceRound       decimal.Decimal `json:"mark_price_round"`
	MarkPriceUp          decimal.Decimal `json:"mark_price_up"`
	MarkPriceDown        decimal.Decimal `json:"mark_price_down"`
	OrderPriceRound      decimal.Decimal `json:"order_price_round"`
	OrderPriceDeviate    decimal.Decimal `json:"order_price_deviate"`
	MarketOrderSlipRatio decimal.Decimal `json:"market_order_slip_ratio"`
	Ask1Price            decimal.Decimal `json:"ask1_price"`
	Ask1Size             int64           `json:"ask1_size"`
	Bid1Price            decimal.Decimal `json:"bid1_price"`
	Bid1Size             int64           `json:"bid1_size"`
	MakerFeeRate         decimal.Decimal `json:"maker_fee_rate"`
	TakerFeeRate         decimal.Decimal `json:"taker_fee_rate"`
	SettleFeeRate        decimal.Decimal `json:"settle_fee_rate"`
	SettleLimitFeeRate   decimal.Decimal `json:"settle_limit_fee_rate"`
	PriceLimitFeeRate    decimal.Decimal `json:"price_limit_fee_rate"`
	RefDiscountRate      decimal.Decimal `json:"ref_discount_rate"`
	RefRebateRate        decimal.Decimal `json:"ref_rebate_rate"`
	InitMarginHigh       decimal.Decimal `json:"init_margin_high"`
	InitMarginLow        decimal.Decimal `json:"init_margin_low"`
	MaintMarginBase      decimal.Decimal `json:"maint_margin_base"`
	OrderSizeMin         int64           `json:"order_size_min"`
	OrderSizeMax         int64           `json:"order_size_max"`
	OrdersLimit          int             `json:"orders_limit"`
	PositionSize         int64           `json:"position_size"`
	PositionLimit        int64           `json:"position_limit"`
	TradeSize            int64           `json:"trade_size"`
	TradeID              int64           `json:"trade_id"`
	OrderbookID          int64           `json:"orderbook_id"`
}

// ListOptionsSettlementsService -- GET /api/v4/options/settlements
//
// Returns the settlement history of an underlying's option contracts.
type ListOptionsSettlementsService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsSettlementsService(underlying string) *ListOptionsSettlementsService {
	return &ListOptionsSettlementsService{c: c, params: map[string]string{"underlying": underlying}}
}

// SetLimit caps the number of records returned in a single page.
func (s *ListOptionsSettlementsService) SetLimit(limit int) *ListOptionsSettlementsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListOptionsSettlementsService) SetOffset(offset int) *ListOptionsSettlementsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom filters records at or after the given time.
func (s *ListOptionsSettlementsService) SetFrom(from time.Time) *ListOptionsSettlementsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo filters records at or before the given time.
func (s *ListOptionsSettlementsService) SetTo(to time.Time) *ListOptionsSettlementsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListOptionsSettlementsService) Do(ctx context.Context) ([]OptionsSettlement, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/settlements", s.params)
	resp, err := request.Do[[]OptionsSettlement](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOptionsSettlementService -- GET /api/v4/options/settlements/{contract}
//
// Returns a single settlement record identified by contract and settlement time.
type GetOptionsSettlementService struct {
	c        *OptionsClient
	contract string
	params   map[string]string
}

func (c *OptionsClient) NewGetOptionsSettlementService(contract, underlying string, at time.Time) *GetOptionsSettlementService {
	return &GetOptionsSettlementService{c: c, contract: contract, params: map[string]string{
		"underlying": underlying,
		"at":         strconv.FormatInt(at.Unix(), 10),
	}}
}

func (s *GetOptionsSettlementService) Do(ctx context.Context) (*OptionsSettlement, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/settlements/"+s.contract, s.params)
	return request.Do[OptionsSettlement](req)
}

// OptionsSettlement is one contract's settlement result.
type OptionsSettlement struct {
	Time        time.Time       `json:"time,format:unix"`
	Contract    string          `json:"contract"`
	Profit      decimal.Decimal `json:"profit"`
	Fee         decimal.Decimal `json:"fee"`
	StrikePrice decimal.Decimal `json:"strike_price"`
	SettlePrice decimal.Decimal `json:"settle_price"`
}

// ListOptionsOrderBookService -- GET /api/v4/options/order_book
//
// Returns the current order book (asks/bids) for an option contract.
type ListOptionsOrderBookService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsOrderBookService(contract string) *ListOptionsOrderBookService {
	return &ListOptionsOrderBookService{c: c, params: map[string]string{"contract": contract}}
}

// SetInterval sets the price precision for depth aggregation ("0" means none).
func (s *ListOptionsOrderBookService) SetInterval(interval string) *ListOptionsOrderBookService {
	s.params["interval"] = interval
	return s
}

// SetLimit sets the number of depth levels returned per side.
func (s *ListOptionsOrderBookService) SetLimit(limit int) *ListOptionsOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetWithID controls whether the order book update id is returned.
func (s *ListOptionsOrderBookService) SetWithID(withID bool) *ListOptionsOrderBookService {
	s.params["with_id"] = strconv.FormatBool(withID)
	return s
}

func (s *ListOptionsOrderBookService) Do(ctx context.Context) (*OptionsOrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/order_book", s.params)
	return request.Do[OptionsOrderBook](req)
}

// OptionsOrderBook is a snapshot of an option contract's order book.
type OptionsOrderBook struct {
	ID      int64                  `json:"id"`
	Current time.Time              `json:"current,format:unix"`
	Update  time.Time              `json:"update,format:unix"`
	Asks    []OptionsOrderBookItem `json:"asks"`
	Bids    []OptionsOrderBookItem `json:"bids"`
}

// OptionsOrderBookItem is one price level in the order book.
type OptionsOrderBookItem struct {
	Price decimal.Decimal `json:"p"`
	Size  int64           `json:"s"`
}

// ListOptionsTickersService -- GET /api/v4/options/tickers
//
// Returns the ticker (mark price, greeks, best bid/ask) of every contract of an
// underlying.
type ListOptionsTickersService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsTickersService(underlying string) *ListOptionsTickersService {
	return &ListOptionsTickersService{c: c, params: map[string]string{"underlying": underlying}}
}

func (s *ListOptionsTickersService) Do(ctx context.Context) ([]OptionsTicker, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/tickers", s.params)
	resp, err := request.Do[[]OptionsTicker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OptionsTicker is a contract's market snapshot including its option greeks.
type OptionsTicker struct {
	Name            string          `json:"name"`
	LastPrice       decimal.Decimal `json:"last_price"`
	MarkPrice       decimal.Decimal `json:"mark_price"`
	IndexPrice      decimal.Decimal `json:"index_price"`
	UnderlyingPrice decimal.Decimal `json:"underlying_price"`
	ExpirationTime  time.Time       `json:"expiration_time,format:unix"`
	Ask1Size        int64           `json:"ask1_size"`
	Ask1Price       decimal.Decimal `json:"ask1_price"`
	Bid1Size        int64           `json:"bid1_size"`
	Bid1Price       decimal.Decimal `json:"bid1_price"`
	PositionSize    int64           `json:"position_size"`
	MarkIV          decimal.Decimal `json:"mark_iv"`
	BidIV           decimal.Decimal `json:"bid_iv"`
	AskIV           decimal.Decimal `json:"ask_iv"`
	Leverage        decimal.Decimal `json:"leverage"`
	Delta           decimal.Decimal `json:"delta"`
	Gamma           decimal.Decimal `json:"gamma"`
	Vega            decimal.Decimal `json:"vega"`
	Theta           decimal.Decimal `json:"theta"`
	Rho             decimal.Decimal `json:"rho"`
}

// ListOptionsUnderlyingTickersService -- GET /api/v4/options/underlying/tickers/{underlying}
//
// Returns aggregate 24h option activity for an underlying.
type ListOptionsUnderlyingTickersService struct {
	c          *OptionsClient
	underlying string
}

func (c *OptionsClient) NewListOptionsUnderlyingTickersService(underlying string) *ListOptionsUnderlyingTickersService {
	return &ListOptionsUnderlyingTickersService{c: c, underlying: underlying}
}

func (s *ListOptionsUnderlyingTickersService) Do(ctx context.Context) (*OptionsUnderlyingTicker, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/underlying/tickers/"+s.underlying)
	return request.Do[OptionsUnderlyingTicker](req)
}

// OptionsUnderlyingTicker is 24h aggregate option statistics for an underlying.
type OptionsUnderlyingTicker struct {
	TradePut   int64           `json:"trade_put"`
	TradeCall  int64           `json:"trade_call"`
	IndexPrice decimal.Decimal `json:"index_price"`
}

// ListOptionsCandlesticksService -- GET /api/v4/options/candlesticks
//
// Returns the OHLC candlesticks of a single option contract.
type ListOptionsCandlesticksService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsCandlesticksService(contract string) *ListOptionsCandlesticksService {
	return &ListOptionsCandlesticksService{c: c, params: map[string]string{"contract": contract}}
}

// SetLimit caps the number of candlesticks returned.
func (s *ListOptionsCandlesticksService) SetLimit(limit int) *ListOptionsCandlesticksService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start time of the candlestick window.
func (s *ListOptionsCandlesticksService) SetFrom(from time.Time) *ListOptionsCandlesticksService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the candlestick window.
func (s *ListOptionsCandlesticksService) SetTo(to time.Time) *ListOptionsCandlesticksService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetInterval sets the candlestick interval (e.g. "1m", "1h", "1d").
func (s *ListOptionsCandlesticksService) SetInterval(interval string) *ListOptionsCandlesticksService {
	s.params["interval"] = interval
	return s
}

func (s *ListOptionsCandlesticksService) Do(ctx context.Context) ([]OptionsCandlestick, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/candlesticks", s.params)
	resp, err := request.Do[[]OptionsCandlestick](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OptionsCandlestick is one option contract OHLC data point.
type OptionsCandlestick struct {
	Timestamp time.Time       `json:"t,format:unix"`
	Volume    int64           `json:"v"`
	Close     decimal.Decimal `json:"c"`
	High      decimal.Decimal `json:"h"`
	Low       decimal.Decimal `json:"l"`
	Open      decimal.Decimal `json:"o"`
}

// ListOptionsUnderlyingCandlesticksService -- GET /api/v4/options/underlying/candlesticks
//
// Returns the OHLC candlesticks of an underlying's index price.
type ListOptionsUnderlyingCandlesticksService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsUnderlyingCandlesticksService(underlying string) *ListOptionsUnderlyingCandlesticksService {
	return &ListOptionsUnderlyingCandlesticksService{c: c, params: map[string]string{"underlying": underlying}}
}

// SetLimit caps the number of candlesticks returned.
func (s *ListOptionsUnderlyingCandlesticksService) SetLimit(limit int) *ListOptionsUnderlyingCandlesticksService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start time of the candlestick window.
func (s *ListOptionsUnderlyingCandlesticksService) SetFrom(from time.Time) *ListOptionsUnderlyingCandlesticksService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the candlestick window.
func (s *ListOptionsUnderlyingCandlesticksService) SetTo(to time.Time) *ListOptionsUnderlyingCandlesticksService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetInterval sets the candlestick interval (e.g. "1m", "1h", "1d").
func (s *ListOptionsUnderlyingCandlesticksService) SetInterval(interval string) *ListOptionsUnderlyingCandlesticksService {
	s.params["interval"] = interval
	return s
}

func (s *ListOptionsUnderlyingCandlesticksService) Do(ctx context.Context) ([]OptionsUnderlyingCandlestick, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/underlying/candlesticks", s.params)
	resp, err := request.Do[[]OptionsUnderlyingCandlestick](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OptionsUnderlyingCandlestick is one underlying index-price OHLC data point.
type OptionsUnderlyingCandlestick struct {
	Timestamp time.Time       `json:"t,format:unix"`
	Close     decimal.Decimal `json:"c"`
	High      decimal.Decimal `json:"h"`
	Low       decimal.Decimal `json:"l"`
	Open      decimal.Decimal `json:"o"`
}

// ListOptionsTradesService -- GET /api/v4/options/trades
//
// Returns recent public trades, filtered by contract or by underlying and
// option type.
type ListOptionsTradesService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsTradesService() *ListOptionsTradesService {
	return &ListOptionsTradesService{c: c, params: map[string]string{}}
}

// SetContract narrows the result to a single option contract.
func (s *ListOptionsTradesService) SetContract(contract string) *ListOptionsTradesService {
	s.params["contract"] = contract
	return s
}

// SetType filters by option type: "C" for call, "P" for put.
func (s *ListOptionsTradesService) SetType(optionType string) *ListOptionsTradesService {
	s.params["type"] = optionType
	return s
}

// SetLimit caps the number of trades returned.
func (s *ListOptionsTradesService) SetLimit(limit int) *ListOptionsTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListOptionsTradesService) SetOffset(offset int) *ListOptionsTradesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom filters trades at or after the given time.
func (s *ListOptionsTradesService) SetFrom(from time.Time) *ListOptionsTradesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo filters trades at or before the given time.
func (s *ListOptionsTradesService) SetTo(to time.Time) *ListOptionsTradesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListOptionsTradesService) Do(ctx context.Context) ([]OptionsTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/trades", s.params)
	resp, err := request.Do[[]OptionsTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OptionsTrade is a single public options trade. Size is signed: negative means
// the taker sold.
type OptionsTrade struct {
	ID         int64           `json:"id"`
	CreateTime time.Time       `json:"create_time,format:unix"`
	Contract   string          `json:"contract"`
	Size       int64           `json:"size"`
	Price      decimal.Decimal `json:"price"`
}
