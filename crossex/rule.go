package crossex

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// QuerySymbolsService -- GET /api/v4/crossex/rule/symbols (private)
//
// Returns cross-exchange trading-pair (symbol) rule information such as size,
// notional and price steps.
type QuerySymbolsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQuerySymbolsService() *QuerySymbolsService {
	return &QuerySymbolsService{c: c, params: map[string]string{}}
}

// SetSymbols narrows the result to a comma-separated list of trading pairs.
func (s *QuerySymbolsService) SetSymbols(symbols string) *QuerySymbolsService {
	s.params["symbols"] = symbols
	return s
}

func (s *QuerySymbolsService) Do(ctx context.Context) ([]CrossexSymbol, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/rule/symbols", s.params).WithSign()
	resp, err := request.Do[[]CrossexSymbol](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexSymbol is a single cross-exchange trading-pair rule. delist_time is a
// millisecond Unix timestamp (0 means not delisted).
type CrossexSymbol struct {
	Symbol          string          `json:"symbol"`
	ExchangeType    string          `json:"exchange_type"`
	BusinessType    string          `json:"business_type"`
	State           string          `json:"state"`
	MinSize         decimal.Decimal `json:"min_size"`
	MinNotional     decimal.Decimal `json:"min_notional"`
	LotSize         decimal.Decimal `json:"lot_size"`
	TickSize        decimal.Decimal `json:"tick_size"`
	MaxNumOrders    decimal.Decimal `json:"max_num_orders"`
	MaxMarketSize   decimal.Decimal `json:"max_market_size"`
	MaxLimitSize    decimal.Decimal `json:"max_limit_size"`
	ContractSize    decimal.Decimal `json:"contract_size"`
	LiquidationFee  decimal.Decimal `json:"liquidation_fee"`
	DefaultLeverage decimal.Decimal `json:"default_leverage"`
	DelistTime      time.Time       `json:"delist_time,string,format:unixmilli"`
}

// QueryRiskLimitsService -- GET /api/v4/crossex/rule/risk_limits (private)
//
// Returns the tiered risk-limit table for the requested trading pairs.
type QueryRiskLimitsService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQueryRiskLimitsService(symbols string) *QueryRiskLimitsService {
	return &QueryRiskLimitsService{c: c, params: map[string]string{
		"symbols": symbols,
	}}
}

func (s *QueryRiskLimitsService) Do(ctx context.Context) ([]CrossexRiskLimit, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/rule/risk_limits", s.params).WithSign()
	resp, err := request.Do[[]CrossexRiskLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexRiskLimit is the tiered risk-limit table for one trading pair.
type CrossexRiskLimit struct {
	Symbol string                 `json:"symbol"`
	Tiers  []CrossexRiskLimitTier `json:"tiers"`
}

// CrossexRiskLimitTier is a single risk-limit tier and its margin/leverage caps.
type CrossexRiskLimitTier struct {
	MinRiskLimitValue decimal.Decimal `json:"min_risk_limit_value"`
	MaxRiskLimitValue decimal.Decimal `json:"max_risk_limit_value"`
	QuickCalAmount    decimal.Decimal `json:"quick_cal_amount"`
	LeverageMax       decimal.Decimal `json:"leverage_max"`
	MaintenanceRate   decimal.Decimal `json:"maintenance_rate"`
	Tier              string          `json:"tier"`
}

// QueryInterestRateService -- GET /api/v4/crossex/interest_rate (private)
//
// Returns the margin-asset hourly interest rates.
type QueryInterestRateService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQueryInterestRateService() *QueryInterestRateService {
	return &QueryInterestRateService{c: c, params: map[string]string{}}
}

// SetCoin narrows the result to a single currency.
func (s *QueryInterestRateService) SetCoin(coin string) *QueryInterestRateService {
	s.params["coin"] = coin
	return s
}

// SetExchangeType narrows the result to a single venue.
func (s *QueryInterestRateService) SetExchangeType(exchangeType string) *QueryInterestRateService {
	s.params["exchange_type"] = exchangeType
	return s
}

func (s *QueryInterestRateService) Do(ctx context.Context) ([]CrossexInterestRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/interest_rate", s.params).WithSign()
	resp, err := request.Do[[]CrossexInterestRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexInterestRate is a margin-asset hourly interest rate. time is a
// millisecond Unix timestamp.
type CrossexInterestRate struct {
	Coin             string          `json:"coin"`
	ExchangeType     string          `json:"exchange_type"`
	HourInterestRate decimal.Decimal `json:"hour_interest_rate"`
	Time             time.Time       `json:"time,string,format:unixmilli"`
}

// GetFeeService -- GET /api/v4/crossex/fee (private)
//
// Returns the authenticated user's per-venue maker/taker fee rates.
type GetFeeService struct {
	c *CrossexClient
}

func (c *CrossexClient) NewGetFeeService() *GetFeeService {
	return &GetFeeService{c: c}
}

func (s *GetFeeService) Do(ctx context.Context) ([]CrossexFee, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/fee").WithSign()
	resp, err := request.Do[[]CrossexFee](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexFee is the maker/taker fee rates for one venue, plus any per-symbol
// special fee overrides.
type CrossexFee struct {
	ExchangeType   string              `json:"exchange_type"`
	SpotMakerFee   decimal.Decimal     `json:"spot_maker_fee"`
	SpotTakerFee   decimal.Decimal     `json:"spot_taker_fee"`
	FutureMakerFee decimal.Decimal     `json:"future_maker_fee"`
	FutureTakerFee decimal.Decimal     `json:"future_taker_fee"`
	SpecialFeeList []CrossexSpecialFee `json:"special_fee_list"`
}

// CrossexSpecialFee is a per-symbol fee-rate override.
type CrossexSpecialFee struct {
	Symbol       string          `json:"symbol"`
	TakerFeeRate decimal.Decimal `json:"taker_fee_rate"`
	MakerFeeRate decimal.Decimal `json:"maker_fee_rate"`
}

// QueryCoinDiscountRateService -- GET /api/v4/crossex/coin_discount_rate (private)
//
// Returns the tiered collateral discount rates by currency and venue.
type QueryCoinDiscountRateService struct {
	c      *CrossexClient
	params map[string]string
}

func (c *CrossexClient) NewQueryCoinDiscountRateService() *QueryCoinDiscountRateService {
	return &QueryCoinDiscountRateService{c: c, params: map[string]string{}}
}

// SetCoin narrows the result to a single currency.
func (s *QueryCoinDiscountRateService) SetCoin(coin string) *QueryCoinDiscountRateService {
	s.params["coin"] = coin
	return s
}

// SetExchangeType narrows the result to a single venue.
func (s *QueryCoinDiscountRateService) SetExchangeType(exchangeType string) *QueryCoinDiscountRateService {
	s.params["exchange_type"] = exchangeType
	return s
}

func (s *QueryCoinDiscountRateService) Do(ctx context.Context) ([]CrossexCoinDiscountRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/crossex/coin_discount_rate", s.params).WithSign()
	resp, err := request.Do[[]CrossexCoinDiscountRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossexCoinDiscountRate is one tier of a currency's collateral discount rate.
type CrossexCoinDiscountRate struct {
	Coin         string          `json:"coin"`
	ExchangeType string          `json:"exchange_type"`
	Tier         string          `json:"tier"`
	MinValue     decimal.Decimal `json:"min_value"`
	MaxValue     decimal.Decimal `json:"max_value"`
	DiscountRate decimal.Decimal `json:"discount_rate"`
}
