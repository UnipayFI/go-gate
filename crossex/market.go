package crossex

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// QueryMarketTickersService -- GET /api/v4/crossex/market/tickers (private)
//
// Returns exchange market tickers across venues. Margin trading pairs cannot be
// passed directly (e.g. GATE_MARGIN_BTC_USDT is invalid).
type QueryMarketTickersService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQueryMarketTickersService() *QueryMarketTickersService {
	return &QueryMarketTickersService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of trading pairs.
func (s *QueryMarketTickersService) SetSymbols(symbols string) *QueryMarketTickersService {
	s.params["symbols"] = symbols
	return s
}

func (s *QueryMarketTickersService) Do(ctx context.Context) ([]CrossexMarketTicker, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/market/tickers", s.params).WithSign()
	resp, err := request.Do[[]CrossexMarketTicker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexMarketTicker is one venue/symbol ticker. Spot symbols leave the
// futures-only fields (mark_price, index_price, open_interest*) empty. timestamp
// is a millisecond Unix timestamp.
type CrossexMarketTicker struct {
	Symbol            string          `json:"symbol"`
	LastPrice         decimal.Decimal `json:"last_price"`
	Open24h           decimal.Decimal `json:"open_24h"`
	Low24h            decimal.Decimal `json:"low_24h"`
	High24h           decimal.Decimal `json:"high_24h"`
	Volume24hBase     decimal.Decimal `json:"volume_24h_base"`
	Volume24hQuote    decimal.Decimal `json:"volume_24h_quote"`
	MarkPrice         decimal.Decimal `json:"mark_price"`
	IndexPrice        decimal.Decimal `json:"index_price"`
	OpenInterest      decimal.Decimal `json:"open_interest"`
	OpenInterestQuote decimal.Decimal `json:"open_interest_quote"`
	Timestamp         time.Time       `json:"timestamp,string,format:unixmilli"`
}

// QueryMarketFundingInfoService -- GET /api/v4/crossex/market/funding_info (private)
//
// Returns futures funding rates, funding intervals and funding timestamps across
// exchanges. For Deribit, funding_rate is the current real-time rate calculated
// over an 8-hour period.
type QueryMarketFundingInfoService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQueryMarketFundingInfoService() *QueryMarketFundingInfoService {
	return &QueryMarketFundingInfoService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of trading pairs.
func (s *QueryMarketFundingInfoService) SetSymbols(symbols string) *QueryMarketFundingInfoService {
	s.params["symbols"] = symbols
	return s
}

func (s *QueryMarketFundingInfoService) Do(ctx context.Context) ([]CrossexMarketFundingInfo, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/market/funding_info", s.params).WithSign()
	resp, err := request.Do[[]CrossexMarketFundingInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexMarketFundingInfo is one venue/symbol futures funding snapshot.
// funding_time is the next funding time as a millisecond Unix timestamp;
// funding_interval is the funding cycle length in seconds.
type CrossexMarketFundingInfo struct {
	Symbol          string          `json:"symbol"`
	FundingRate     decimal.Decimal `json:"funding_rate"`
	FundingTime     time.Time       `json:"funding_time,string,format:unixmilli"`
	FundingInterval int64           `json:"funding_interval,string"`
}
