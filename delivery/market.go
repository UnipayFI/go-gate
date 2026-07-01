package delivery

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListDeliveryContractsService -- GET /api/v4/delivery/{settle}/contracts
//
// Returns every delivery (dated-futures) contract and its trading rules.
type ListDeliveryContractsService struct {
	c      *DeliveryClient
	settle Settle
}

func (c *DeliveryClient) NewListDeliveryContractsService(settle Settle) *ListDeliveryContractsService {
	return &ListDeliveryContractsService{c: c, settle: settle}
}

func (s *ListDeliveryContractsService) Do(ctx context.Context) ([]DeliveryContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/contracts")
	resp, err := request.Do[[]DeliveryContract](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetDeliveryContractService -- GET /api/v4/delivery/{settle}/contracts/{contract}
//
// Returns the trading rules and current market snapshot for a single delivery
// contract.
type GetDeliveryContractService struct {
	c        *DeliveryClient
	settle   Settle
	contract string
}

func (c *DeliveryClient) NewGetDeliveryContractService(settle Settle, contract string) *GetDeliveryContractService {
	return &GetDeliveryContractService{c: c, settle: settle, contract: contract}
}

func (s *GetDeliveryContractService) Do(ctx context.Context) (*DeliveryContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/contracts/"+s.contract)
	return request.Do[DeliveryContract](req)
}

// DeliveryContract is a dated-futures contract and its trading rules. ExpireTime
// and ConfigChangeTime are unix-second epochs sent as bare numbers.
type DeliveryContract struct {
	Name                string          `json:"name"`
	Underlying          string          `json:"underlying"`
	Cycle               string          `json:"cycle"`
	Type                string          `json:"type"`
	QuantoMultiplier    decimal.Decimal `json:"quanto_multiplier"`
	LeverageMin         decimal.Decimal `json:"leverage_min"`
	LeverageMax         decimal.Decimal `json:"leverage_max"`
	MaintenanceRate     decimal.Decimal `json:"maintenance_rate"`
	MarkType            string          `json:"mark_type"`
	MarkPrice           decimal.Decimal `json:"mark_price"`
	IndexPrice          decimal.Decimal `json:"index_price"`
	LastPrice           decimal.Decimal `json:"last_price"`
	MakerFeeRate        decimal.Decimal `json:"maker_fee_rate"`
	TakerFeeRate        decimal.Decimal `json:"taker_fee_rate"`
	OrderPriceRound     decimal.Decimal `json:"order_price_round"`
	MarkPriceRound      decimal.Decimal `json:"mark_price_round"`
	BasisRate           decimal.Decimal `json:"basis_rate"`
	BasisValue          decimal.Decimal `json:"basis_value"`
	BasisImpactValue    decimal.Decimal `json:"basis_impact_value"`
	SettlePrice         decimal.Decimal `json:"settle_price"`
	SettleFeeRate       decimal.Decimal `json:"settle_fee_rate"`
	SettlePriceInterval int             `json:"settle_price_interval"`
	SettlePriceDuration int             `json:"settle_price_duration"`
	ExpireTime          time.Time       `json:"expire_time,format:unix"`
	RiskLimitBase       decimal.Decimal `json:"risk_limit_base"`
	RiskLimitStep       decimal.Decimal `json:"risk_limit_step"`
	RiskLimitMax        decimal.Decimal `json:"risk_limit_max"`
	OrderSizeMin        int64           `json:"order_size_min"`
	OrderSizeMax        int64           `json:"order_size_max"`
	OrderPriceDeviate   decimal.Decimal `json:"order_price_deviate"`
	RefDiscountRate     decimal.Decimal `json:"ref_discount_rate"`
	RefRebateRate       decimal.Decimal `json:"ref_rebate_rate"`
	OrderbookID         int64           `json:"orderbook_id"`
	TradeID             int64           `json:"trade_id"`
	TradeSize           int64           `json:"trade_size"`
	PositionSize        int64           `json:"position_size"`
	ConfigChangeTime    time.Time       `json:"config_change_time,format:unix"`
	InDelisting         bool            `json:"in_delisting"`
	OrdersLimit         int             `json:"orders_limit"`
}

// ListDeliveryOrderBookService -- GET /api/v4/delivery/{settle}/order_book
//
// Returns the current market depth for a delivery contract: bids sorted
// high-to-low and asks low-to-high.
type ListDeliveryOrderBookService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryOrderBookService(settle Settle, contract string) *ListDeliveryOrderBookService {
	return &ListDeliveryOrderBookService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetInterval sets the price precision used to aggregate depth levels
// ("0" means no aggregation).
func (s *ListDeliveryOrderBookService) SetInterval(interval string) *ListDeliveryOrderBookService {
	s.params["interval"] = interval
	return s
}

// SetLimit caps the number of depth levels returned per side.
func (s *ListDeliveryOrderBookService) SetLimit(limit int) *ListDeliveryOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetWithID requests the order book update ID (populates the id field).
func (s *ListDeliveryOrderBookService) SetWithID(withID bool) *ListDeliveryOrderBookService {
	s.params["with_id"] = strconv.FormatBool(withID)
	return s
}

func (s *ListDeliveryOrderBookService) Do(ctx context.Context) (*DeliveryOrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/order_book", s.params)
	return request.Do[DeliveryOrderBook](req)
}

// DeliveryOrderBook is a snapshot of market depth. ID is valid only when with_id
// is set. Current/Update are unix-second epochs sent as bare numbers.
type DeliveryOrderBook struct {
	ID      int64                    `json:"id"`
	Current time.Time                `json:"current,format:unix"`
	Update  time.Time                `json:"update,format:unix"`
	Asks    []DeliveryOrderBookEntry `json:"asks"`
	Bids    []DeliveryOrderBookEntry `json:"bids"`
}

// DeliveryOrderBookEntry is a single price level: Price (quote currency) and Size
// (contract count).
type DeliveryOrderBookEntry struct {
	Price decimal.Decimal `json:"p"`
	Size  int64           `json:"s"`
}

// ListDeliveryTradesService -- GET /api/v4/delivery/{settle}/trades
//
// Returns recent public market trades for a delivery contract, most recent first.
type ListDeliveryTradesService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryTradesService(settle Settle, contract string) *ListDeliveryTradesService {
	return &ListDeliveryTradesService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetLimit caps the number of trades returned in a single list.
func (s *ListDeliveryTradesService) SetLimit(limit int) *ListDeliveryTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetLastID uses the id of the last record from the previous page as the cursor
// (deprecated by Gate in favour of from/to).
func (s *ListDeliveryTradesService) SetLastID(lastID string) *ListDeliveryTradesService {
	s.params["last_id"] = lastID
	return s
}

// SetFrom sets the start of the query time range (unix seconds).
func (s *ListDeliveryTradesService) SetFrom(from time.Time) *ListDeliveryTradesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range (unix seconds, defaults to now).
func (s *ListDeliveryTradesService) SetTo(to time.Time) *ListDeliveryTradesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListDeliveryTradesService) Do(ctx context.Context) ([]DeliveryTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/trades", s.params)
	resp, err := request.Do[[]DeliveryTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeliveryTrade is a single public market trade. Size is signed (negative for a
// taker sell). CreateTime and CreateTimeMs are unix-second epochs sent as bare
// numbers, the latter carrying millisecond precision as a fraction. IsInternal is
// only returned for insurance/ADL takeover trades.
type DeliveryTrade struct {
	ID           int64           `json:"id"`
	CreateTime   time.Time       `json:"create_time,format:unix"`
	CreateTimeMs time.Time       `json:"create_time_ms,format:unix"`
	Contract     string          `json:"contract"`
	Size         int64           `json:"size"`
	Price        decimal.Decimal `json:"price"`
	IsInternal   bool            `json:"is_internal"`
}

// ListDeliveryCandlesticksService -- GET /api/v4/delivery/{settle}/candlesticks
//
// Returns OHLC candlestick data for a delivery contract. Prefix the contract with
// "mark_" for mark-price candles or "index_" for index-price candles. At most
// 2000 points per query.
type ListDeliveryCandlesticksService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryCandlesticksService(settle Settle, contract string) *ListDeliveryCandlesticksService {
	return &ListDeliveryCandlesticksService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetFrom sets the start time of the candlestick range (unix seconds).
func (s *ListDeliveryCandlesticksService) SetFrom(from time.Time) *ListDeliveryCandlesticksService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the candlestick range (unix seconds, defaults to now).
func (s *ListDeliveryCandlesticksService) SetTo(to time.Time) *ListDeliveryCandlesticksService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of recent points returned (conflicts with from/to).
func (s *ListDeliveryCandlesticksService) SetLimit(limit int) *ListDeliveryCandlesticksService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetInterval selects the time interval between points (e.g. "1m", "1h", "1d";
// "1w" is a natural week).
func (s *ListDeliveryCandlesticksService) SetInterval(interval string) *ListDeliveryCandlesticksService {
	s.params["interval"] = interval
	return s
}

func (s *ListDeliveryCandlesticksService) Do(ctx context.Context) ([]DeliveryCandlestick, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/candlesticks", s.params)
	resp, err := request.Do[[]DeliveryCandlestick](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeliveryCandlestick is a single OHLC point. Timestamp is a unix-second epoch
// sent as a bare number. Volume (contract size) is only returned when the
// contract is not prefixed with mark_/index_.
type DeliveryCandlestick struct {
	Timestamp time.Time       `json:"t,format:unix"`
	Volume    int64           `json:"v"`
	Close     decimal.Decimal `json:"c"`
	High      decimal.Decimal `json:"h"`
	Low       decimal.Decimal `json:"l"`
	Open      decimal.Decimal `json:"o"`
}

// ListDeliveryTickersService -- GET /api/v4/delivery/{settle}/tickers
//
// Returns 24h trading statistics for all delivery contracts, or one when contract
// is set.
type ListDeliveryTickersService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryTickersService(settle Settle) *ListDeliveryTickersService {
	return &ListDeliveryTickersService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single delivery contract.
func (s *ListDeliveryTickersService) SetContract(contract string) *ListDeliveryTickersService {
	s.params["contract"] = contract
	return s
}

func (s *ListDeliveryTickersService) Do(ctx context.Context) ([]DeliveryTicker, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/tickers", s.params)
	resp, err := request.Do[[]DeliveryTicker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeliveryTicker is 24h rolling market statistics for a delivery contract. The
// funding-rate/quanto fields are inherited from the futures ticker shape and are
// absent from delivery responses.
type DeliveryTicker struct {
	Contract              string          `json:"contract"`
	Last                  decimal.Decimal `json:"last"`
	ChangePercentage      decimal.Decimal `json:"change_percentage"`
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
	BasisRate             decimal.Decimal `json:"basis_rate"`
	BasisValue            decimal.Decimal `json:"basis_value"`
	SettlePrice           decimal.Decimal `json:"settle_price"`
	LowestAsk             decimal.Decimal `json:"lowest_ask"`
	LowestSize            decimal.Decimal `json:"lowest_size"`
	HighestBid            decimal.Decimal `json:"highest_bid"`
	HighestSize           decimal.Decimal `json:"highest_size"`
}

// ListDeliveryInsuranceLedgerService -- GET /api/v4/delivery/{settle}/insurance
//
// Returns the delivery insurance-fund balance history.
type ListDeliveryInsuranceLedgerService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryInsuranceLedgerService(settle Settle) *ListDeliveryInsuranceLedgerService {
	return &ListDeliveryInsuranceLedgerService{c: c, settle: settle, params: map[string]string{}}
}

// SetLimit caps the number of records returned in a single list.
func (s *ListDeliveryInsuranceLedgerService) SetLimit(limit int) *ListDeliveryInsuranceLedgerService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListDeliveryInsuranceLedgerService) Do(ctx context.Context) ([]DeliveryInsurance, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/insurance", s.params)
	resp, err := request.Do[[]DeliveryInsurance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeliveryInsurance is a single insurance-fund balance snapshot. Timestamp is a
// unix-second epoch sent as a bare number.
type DeliveryInsurance struct {
	Timestamp time.Time       `json:"t,format:unix"`
	Balance   decimal.Decimal `json:"b"`
}
