package spot

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetServerTimeService -- GET /api/v4/spot/time
//
// Returns the exchange server time in milliseconds.
type GetServerTimeService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTime, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/time")
	return request.Do[ServerTime](req)
}

// ServerTime is the exchange clock. server_time is a millisecond epoch number.
type ServerTime struct {
	ServerTime time.Time `json:"server_time,format:unixmilli"`
}

// ListCurrencyPairsService -- GET /api/v4/spot/currency_pairs
//
// Returns every spot trading pair and its trading rules.
type ListCurrencyPairsService struct {
	c *SpotClient
}

func (c *SpotClient) NewListCurrencyPairsService() *ListCurrencyPairsService {
	return &ListCurrencyPairsService{c: c}
}

func (s *ListCurrencyPairsService) Do(ctx context.Context) ([]CurrencyPair, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/currency_pairs")
	resp, err := request.Do[[]CurrencyPair](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetCurrencyPairService -- GET /api/v4/spot/currency_pairs/{currency_pair}
//
// Returns the trading rules for a single spot pair.
type GetCurrencyPairService struct {
	c            *SpotClient
	currencyPair string
}

func (c *SpotClient) NewGetCurrencyPairService(currencyPair string) *GetCurrencyPairService {
	return &GetCurrencyPairService{c: c, currencyPair: currencyPair}
}

func (s *GetCurrencyPairService) Do(ctx context.Context) (*CurrencyPair, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/currency_pairs/"+s.currencyPair)
	return request.Do[CurrencyPair](req)
}

// CurrencyPair is a spot trading pair and its trading rules.
type CurrencyPair struct {
	ID                  string          `json:"id"`
	Base                string          `json:"base"`
	BaseName            string          `json:"base_name"`
	Quote               string          `json:"quote"`
	QuoteName           string          `json:"quote_name"`
	Fee                 decimal.Decimal `json:"fee"`
	MinBaseAmount       decimal.Decimal `json:"min_base_amount"`
	MinQuoteAmount      decimal.Decimal `json:"min_quote_amount"`
	MaxBaseAmount       decimal.Decimal `json:"max_base_amount"`
	MaxQuoteAmount      decimal.Decimal `json:"max_quote_amount"`
	AmountPrecision     int             `json:"amount_precision"`
	Precision           int             `json:"precision"`
	TradeStatus         TradeStatus     `json:"trade_status"`
	SellStart           time.Time       `json:"sell_start,format:unix"`
	BuyStart            time.Time       `json:"buy_start,format:unix"`
	Type                string          `json:"type"`
	STTag               bool            `json:"st_tag"`
	Slippage            decimal.Decimal `json:"slippage"`
	TradeURL            string          `json:"trade_url"`
	MarketOrderMaxStock decimal.Decimal `json:"market_order_max_stock"`
	MarketOrderMaxMoney decimal.Decimal `json:"market_order_max_money"`
	UpRate              decimal.Decimal `json:"up_rate"`
	DownRate            decimal.Decimal `json:"down_rate"`
}

// GetTickersService -- GET /api/v4/spot/tickers
//
// Returns 24h ticker statistics for all pairs, or one when currency_pair is set.
type GetTickersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetTickersService() *GetTickersService {
	return &GetTickersService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *GetTickersService) SetCurrencyPair(currencyPair string) *GetTickersService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetTimezone selects the timezone for the change_utc0 / change_utc8 fields
// ("utc0", "utc8", or "all").
func (s *GetTickersService) SetTimezone(timezone string) *GetTickersService {
	s.params["timezone"] = timezone
	return s
}

func (s *GetTickersService) Do(ctx context.Context) ([]Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/tickers", s.params)
	resp, err := request.Do[[]Ticker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Ticker is 24h rolling market statistics for a spot pair.
type Ticker struct {
	CurrencyPair     string          `json:"currency_pair"`
	Last             decimal.Decimal `json:"last"`
	LowestAsk        decimal.Decimal `json:"lowest_ask"`
	LowestSize       decimal.Decimal `json:"lowest_size"`
	HighestBid       decimal.Decimal `json:"highest_bid"`
	HighestSize      decimal.Decimal `json:"highest_size"`
	ChangePercentage decimal.Decimal `json:"change_percentage"`
	ChangeUTC0       decimal.Decimal `json:"change_utc0"`
	ChangeUTC8       decimal.Decimal `json:"change_utc8"`
	BaseVolume       decimal.Decimal `json:"base_volume"`
	QuoteVolume      decimal.Decimal `json:"quote_volume"`
	High24h          decimal.Decimal `json:"high_24h"`
	Low24h           decimal.Decimal `json:"low_24h"`
	ETFNetValue      decimal.Decimal `json:"etf_net_value"`
	ETFPreNetValue   decimal.Decimal `json:"etf_pre_net_value"`
	ETFPreTimestamp  time.Time       `json:"etf_pre_timestamp,format:unix"`
	ETFLeverage      decimal.Decimal `json:"etf_leverage"`
}
