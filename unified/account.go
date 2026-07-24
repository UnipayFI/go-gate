package unified

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListUnifiedAccountsService -- GET /api/v4/unified/accounts (private)
//
// Returns the unified account overview: per-currency balances plus the total
// asset/margin figures, all converted to USD according to each currency's
// liquidity-adjustment coefficient.
type ListUnifiedAccountsService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewListUnifiedAccountsService() *ListUnifiedAccountsService {
	return &ListUnifiedAccountsService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the balances to a single currency (e.g. USDT).
func (s *ListUnifiedAccountsService) SetCurrency(currency string) *ListUnifiedAccountsService {
	s.params["currency"] = currency
	return s
}

// SetSubUID queries the unified account of a sub account by its user ID.
func (s *ListUnifiedAccountsService) SetSubUID(subUID string) *ListUnifiedAccountsService {
	s.params["sub_uid"] = subUID
	return s
}

func (s *ListUnifiedAccountsService) Do(ctx context.Context) (*UnifiedAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/accounts", s.params).WithSign()
	return request.Do[UnifiedAccount](req)
}

// UnifiedAccount is the unified-account overview, aggregating every currency's
// balance and the account-wide margin totals.
type UnifiedAccount struct {
	UserID                     int64                     `json:"user_id"`
	RefreshTime                time.Time                 `json:"refresh_time,format:unixmilli"`
	Locked                     bool                      `json:"locked"`
	Balances                   map[string]UnifiedBalance `json:"balances"`
	Total                      decimal.Decimal           `json:"total"`
	Borrowed                   decimal.Decimal           `json:"borrowed"`
	TotalInitialMargin         decimal.Decimal           `json:"total_initial_margin"`
	TotalMarginBalance         decimal.Decimal           `json:"total_margin_balance"`
	TotalMaintenanceMargin     decimal.Decimal           `json:"total_maintenance_margin"`
	TotalInitialMarginRate     decimal.Decimal           `json:"total_initial_margin_rate"`
	TotalMaintenanceMarginRate decimal.Decimal           `json:"total_maintenance_margin_rate"`
	TotalAvailableMargin       decimal.Decimal           `json:"total_available_margin"`
	UnifiedAccountTotal        decimal.Decimal           `json:"unified_account_total"`
	UnifiedAccountTotalLiab    decimal.Decimal           `json:"unified_account_total_liab"`
	UnifiedAccountTotalEquity  decimal.Decimal           `json:"unified_account_total_equity"`
	Leverage                   decimal.Decimal           `json:"leverage"`
	SpotOrderLoss              decimal.Decimal           `json:"spot_order_loss"`
	SpotHedge                  bool                      `json:"spot_hedge"`
	UseFunding                 bool                      `json:"use_funding"`
	IsAllCollateral            bool                      `json:"is_all_collateral"`
}

// UnifiedBalance is a single currency's balance within the unified account.
// Which fields are populated depends on the account margin mode.
type UnifiedBalance struct {
	Available         decimal.Decimal `json:"available"`
	Freeze            decimal.Decimal `json:"freeze"`
	Borrowed          decimal.Decimal `json:"borrowed"`
	NegativeLiab      decimal.Decimal `json:"negative_liab"`
	FuturesPosLiab    decimal.Decimal `json:"futures_pos_liab"`
	Equity            decimal.Decimal `json:"equity"`
	TotalFreeze       decimal.Decimal `json:"total_freeze"`
	TotalLiab         decimal.Decimal `json:"total_liab"`
	SpotInUse         decimal.Decimal `json:"spot_in_use"`
	Funding           decimal.Decimal `json:"funding"`
	FundingVersion    string          `json:"funding_version"`
	CrossBalance      decimal.Decimal `json:"cross_balance"`
	IsoBalance        decimal.Decimal `json:"iso_balance"`
	IM                decimal.Decimal `json:"im"`
	MM                decimal.Decimal `json:"mm"`
	IMR               decimal.Decimal `json:"imr"`
	MMR               decimal.Decimal `json:"mmr"`
	MarginBalance     decimal.Decimal `json:"margin_balance"`
	AvailableMargin   decimal.Decimal `json:"available_margin"`
	EnabledCollateral bool            `json:"enabled_collateral"`
}

// GetUnifiedBorrowableService -- GET /api/v4/unified/borrowable (private)
//
// Returns the maximum amount of a currency the unified account can still borrow.
type GetUnifiedBorrowableService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUnifiedBorrowableService(currency string) *GetUnifiedBorrowableService {
	return &GetUnifiedBorrowableService{c: c, params: map[string]string{"currency": currency}}
}

func (s *GetUnifiedBorrowableService) Do(ctx context.Context) (*UnifiedBorrowable, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/borrowable", s.params).WithSign()
	return request.Do[UnifiedBorrowable](req)
}

// UnifiedBorrowable is a currency's maximum borrowable amount.
type UnifiedBorrowable struct {
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

// GetUnifiedTransferableService -- GET /api/v4/unified/transferable (private)
//
// Returns the maximum amount of a currency that can be transferred out of the
// unified account.
type GetUnifiedTransferableService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUnifiedTransferableService(currency string) *GetUnifiedTransferableService {
	return &GetUnifiedTransferableService{c: c, params: map[string]string{"currency": currency}}
}

func (s *GetUnifiedTransferableService) Do(ctx context.Context) (*UnifiedTransferable, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/transferable", s.params).WithSign()
	return request.Do[UnifiedTransferable](req)
}

// UnifiedTransferable is a currency's maximum transferable amount.
type UnifiedTransferable struct {
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

// GetUnifiedTransferablesService -- GET /api/v4/unified/transferables (private)
//
// Batch variant of GetUnifiedTransferable: the maximum transferable amount for
// each of up to 100 currencies.
type GetUnifiedTransferablesService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUnifiedTransferablesService(currencies []string) *GetUnifiedTransferablesService {
	return &GetUnifiedTransferablesService{c: c, params: map[string]string{"currencies": strings.Join(currencies, ",")}}
}

func (s *GetUnifiedTransferablesService) Do(ctx context.Context) ([]UnifiedTransferable, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/transferables", s.params).WithSign()
	resp, err := request.Do[[]UnifiedTransferable](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetUnifiedBorrowableListService -- GET /api/v4/unified/batch_borrowable (private)
//
// Batch variant of GetUnifiedBorrowable: the maximum borrowable amount for each
// of up to 10 currencies.
type GetUnifiedBorrowableListService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUnifiedBorrowableListService(currencies []string) *GetUnifiedBorrowableListService {
	return &GetUnifiedBorrowableListService{c: c, params: map[string]string{"currencies": strings.Join(currencies, ",")}}
}

func (s *GetUnifiedBorrowableListService) Do(ctx context.Context) ([]UnifiedBorrowable, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/batch_borrowable", s.params).WithSign()
	resp, err := request.Do[[]UnifiedBorrowable](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetUnifiedRiskUnitsService -- GET /api/v4/unified/risk_units (private)
//
// Returns the account's risk-unit breakdown. Only meaningful in portfolio
// margin mode.
type GetUnifiedRiskUnitsService struct {
	c *UnifiedClient
}

func (c *UnifiedClient) NewGetUnifiedRiskUnitsService() *GetUnifiedRiskUnitsService {
	return &GetUnifiedRiskUnitsService{c: c}
}

func (s *GetUnifiedRiskUnitsService) Do(ctx context.Context) (*UnifiedRiskUnits, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/risk_units").WithSign()
	return request.Do[UnifiedRiskUnits](req)
}

// UnifiedRiskUnits is the account's portfolio-margin risk-unit detail.
type UnifiedRiskUnits struct {
	UserID    int64             `json:"user_id"`
	SpotHedge bool              `json:"spot_hedge"`
	RiskUnits []UnifiedRiskUnit `json:"risk_units"`
}

// UnifiedRiskUnit is a single risk unit and its aggregated greeks and margins.
type UnifiedRiskUnit struct {
	Symbol         string          `json:"symbol"`
	SpotInUse      decimal.Decimal `json:"spot_in_use"`
	MaintainMargin decimal.Decimal `json:"maintain_margin"`
	InitialMargin  decimal.Decimal `json:"initial_margin"`
	Delta          decimal.Decimal `json:"delta"`
	Gamma          decimal.Decimal `json:"gamma"`
	Theta          decimal.Decimal `json:"theta"`
	Vega           decimal.Decimal `json:"vega"`
}

// GetUnifiedModeService -- GET /api/v4/unified/unified_mode (private)
//
// Returns the account's unified mode and its per-mode toggles.
type GetUnifiedModeService struct {
	c *UnifiedClient
}

func (c *UnifiedClient) NewGetUnifiedModeService() *GetUnifiedModeService {
	return &GetUnifiedModeService{c: c}
}

func (s *GetUnifiedModeService) Do(ctx context.Context) (*UnifiedMode, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/unified_mode").WithSign()
	return request.Do[UnifiedMode](req)
}

// UnifiedMode is the account's margin mode and the switches enabled under it.
// mode is one of classic, multi_currency, portfolio or single_currency.
type UnifiedMode struct {
	Mode     string              `json:"mode"`
	Settings UnifiedModeSettings `json:"settings"`
}

// UnifiedModeSettings holds the per-mode feature switches.
type UnifiedModeSettings struct {
	USDTFutures bool `json:"usdt_futures"`
	SpotHedge   bool `json:"spot_hedge"`
	UseFunding  bool `json:"use_funding"`
	Options     bool `json:"options"`
}

// SetUnifiedModeService -- PUT /api/v4/unified/unified_mode (private)
//
// Switches the account's unified mode, optionally toggling the switches
// available under the target mode.
type SetUnifiedModeService struct {
	c    *UnifiedClient
	body map[string]any
}

// NewSetUnifiedModeService switches to mode (classic, multi_currency, portfolio
// or single_currency).
func (c *UnifiedClient) NewSetUnifiedModeService(mode string) *SetUnifiedModeService {
	return &SetUnifiedModeService{c: c, body: map[string]any{"mode": mode}}
}

func (s *SetUnifiedModeService) setting(key string, value bool) *SetUnifiedModeService {
	settings, _ := s.body["settings"].(map[string]any)
	if settings == nil {
		settings = map[string]any{}
		s.body["settings"] = settings
	}
	settings[key] = value
	return s
}

// SetUSDTFutures toggles the USDT-futures switch (multi_currency mode can only
// enable it, not disable it).
func (s *SetUnifiedModeService) SetUSDTFutures(enabled bool) *SetUnifiedModeService {
	return s.setting("usdt_futures", enabled)
}

// SetSpotHedge toggles the spot-hedging switch.
func (s *SetUnifiedModeService) SetSpotHedge(enabled bool) *SetUnifiedModeService {
	return s.setting("spot_hedge", enabled)
}

// SetUseFunding toggles using Earn funds as margin.
func (s *SetUnifiedModeService) SetUseFunding(enabled bool) *SetUnifiedModeService {
	return s.setting("use_funding", enabled)
}

// SetOptions toggles the options switch (multi_currency mode can only enable it).
func (s *SetUnifiedModeService) SetOptions(enabled bool) *SetUnifiedModeService {
	return s.setting("options", enabled)
}

func (s *SetUnifiedModeService) Do(ctx context.Context) error {
	req := request.Put(ctx, s.c, "/api/v4/unified/unified_mode", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// GetUnifiedEstimateRateService -- GET /api/v4/unified/estimate_rate (private)
//
// Returns the estimated hourly borrow rate for each requested currency. Rates
// fluctuate with lending depth, and an unsupported currency comes back empty.
type GetUnifiedEstimateRateService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUnifiedEstimateRateService(currencies []string) *GetUnifiedEstimateRateService {
	return &GetUnifiedEstimateRateService{c: c, params: map[string]string{"currencies": strings.Join(currencies, ",")}}
}

func (s *GetUnifiedEstimateRateService) Do(ctx context.Context) (map[string]UnifiedEstimateRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/estimate_rate", s.params).WithSign()
	resp, err := request.Do[map[string]UnifiedEstimateRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnifiedEstimateRate is a currency's estimated hourly borrow rate. Gate returns
// it as a quoted number, or an empty string when the currency is not supported.
type UnifiedEstimateRate = decimal.Decimal

// ListCurrencyDiscountTiersService -- GET /api/v4/unified/currency_discount_tiers
//
// Returns the tiered collateral-discount schedule for each currency.
type ListCurrencyDiscountTiersService struct {
	c *UnifiedClient
}

func (c *UnifiedClient) NewListCurrencyDiscountTiersService() *ListCurrencyDiscountTiersService {
	return &ListCurrencyDiscountTiersService{c: c}
}

func (s *ListCurrencyDiscountTiersService) Do(ctx context.Context) ([]CurrencyDiscountTier, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/currency_discount_tiers")
	resp, err := request.Do[[]CurrencyDiscountTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CurrencyDiscountTier is a currency's tiered collateral-discount schedule.
type CurrencyDiscountTier struct {
	Currency      string                `json:"currency"`
	DiscountTiers []UnifiedDiscountTier `json:"discount_tiers"`
}

// UnifiedDiscountTier is a single collateral-discount tier. upper_limit is "+"
// on the last (open-ended) tier, so it stays a string.
type UnifiedDiscountTier struct {
	Tier       string          `json:"tier"`
	Discount   decimal.Decimal `json:"discount"`
	LowerLimit decimal.Decimal `json:"lower_limit"`
	UpperLimit string          `json:"upper_limit"`
	Leverage   decimal.Decimal `json:"leverage"`
}

// ListLoanMarginTiersService -- GET /api/v4/unified/loan_margin_tiers
//
// Returns the tiered loan-margin (borrow-margin-rate) schedule for each currency.
type ListLoanMarginTiersService struct {
	c *UnifiedClient
}

func (c *UnifiedClient) NewListLoanMarginTiersService() *ListLoanMarginTiersService {
	return &ListLoanMarginTiersService{c: c}
}

func (s *ListLoanMarginTiersService) Do(ctx context.Context) ([]UnifiedLoanMarginTier, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/loan_margin_tiers")
	resp, err := request.Do[[]UnifiedLoanMarginTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnifiedLoanMarginTier is a currency's tiered loan-margin schedule.
type UnifiedLoanMarginTier struct {
	Currency    string              `json:"currency"`
	MarginTiers []UnifiedMarginTier `json:"margin_tiers"`
}

// UnifiedMarginTier is a single loan-margin tier. upper_limit is "" on the last
// (open-ended) tier, so it stays a string.
type UnifiedMarginTier struct {
	Tier       string          `json:"tier"`
	MarginRate decimal.Decimal `json:"margin_rate"`
	LowerLimit decimal.Decimal `json:"lower_limit"`
	UpperLimit string          `json:"upper_limit"`
	Leverage   decimal.Decimal `json:"leverage"`
}

// CalculatePortfolioMarginService -- POST /api/v4/unified/portfolio_calculator
//
// Runs the portfolio-margin calculator over a simulated set of spot balances,
// futures/options positions and orders. It is a stateless calculator and needs
// no authentication.
type CalculatePortfolioMarginService struct {
	c    *UnifiedClient
	body map[string]any
}

func (c *UnifiedClient) NewCalculatePortfolioMarginService() *CalculatePortfolioMarginService {
	return &CalculatePortfolioMarginService{c: c, body: map[string]any{}}
}

// SetSpotBalances sets the simulated spot balances (currently BTC and ETH only).
func (s *CalculatePortfolioMarginService) SetSpotBalances(balances []PortfolioSpotBalance) *CalculatePortfolioMarginService {
	s.body["spot_balances"] = balances
	return s
}

// SetSpotOrders sets the simulated spot orders.
func (s *CalculatePortfolioMarginService) SetSpotOrders(orders []PortfolioSpotOrder) *CalculatePortfolioMarginService {
	s.body["spot_orders"] = orders
	return s
}

// SetFuturesPositions sets the simulated futures positions.
func (s *CalculatePortfolioMarginService) SetFuturesPositions(positions []PortfolioFuturesPosition) *CalculatePortfolioMarginService {
	s.body["futures_positions"] = positions
	return s
}

// SetFuturesOrders sets the simulated futures orders.
func (s *CalculatePortfolioMarginService) SetFuturesOrders(orders []PortfolioFuturesOrder) *CalculatePortfolioMarginService {
	s.body["futures_orders"] = orders
	return s
}

// SetOptionsPositions sets the simulated options positions.
func (s *CalculatePortfolioMarginService) SetOptionsPositions(positions []PortfolioOptionsPosition) *CalculatePortfolioMarginService {
	s.body["options_positions"] = positions
	return s
}

// SetOptionsOrders sets the simulated options orders.
func (s *CalculatePortfolioMarginService) SetOptionsOrders(orders []PortfolioOptionsOrder) *CalculatePortfolioMarginService {
	s.body["options_orders"] = orders
	return s
}

// SetSpotHedge enables spot hedging in the simulation.
func (s *CalculatePortfolioMarginService) SetSpotHedge(enabled bool) *CalculatePortfolioMarginService {
	s.body["spot_hedge"] = enabled
	return s
}

func (s *CalculatePortfolioMarginService) Do(ctx context.Context) (*PortfolioMargin, error) {
	req := request.Post(ctx, s.c, "/api/v4/unified/portfolio_calculator", s.body)
	return request.Do[PortfolioMargin](req)
}

// PortfolioSpotBalance is a simulated spot balance (equity = balance - borrowed).
type PortfolioSpotBalance struct {
	Currency string          `json:"currency"`
	Equity   decimal.Decimal `json:"equity"`
}

// PortfolioSpotOrder is a simulated spot order.
type PortfolioSpotOrder struct {
	CurrencyPairs string          `json:"currency_pairs"`
	OrderPrice    decimal.Decimal `json:"order_price"`
	Count         decimal.Decimal `json:"count,omitempty"`
	Left          decimal.Decimal `json:"left"`
	Type          string          `json:"type"`
}

// PortfolioFuturesPosition is a simulated futures position (size in contracts).
type PortfolioFuturesPosition struct {
	Contract string          `json:"contract"`
	Size     decimal.Decimal `json:"size"`
}

// PortfolioFuturesOrder is a simulated futures order.
type PortfolioFuturesOrder struct {
	Contract string          `json:"contract"`
	Size     decimal.Decimal `json:"size"`
	Left     decimal.Decimal `json:"left"`
}

// PortfolioOptionsPosition is a simulated options position.
type PortfolioOptionsPosition struct {
	OptionsName string          `json:"options_name"`
	Size        decimal.Decimal `json:"size"`
}

// PortfolioOptionsOrder is a simulated options order.
type PortfolioOptionsOrder struct {
	OptionsName string          `json:"options_name"`
	Size        decimal.Decimal `json:"size"`
	Left        decimal.Decimal `json:"left"`
}

// PortfolioMargin is the portfolio-margin calculator result.
type PortfolioMargin struct {
	MaintainMarginTotal decimal.Decimal     `json:"maintain_margin_total"`
	InitialMarginTotal  decimal.Decimal     `json:"initial_margin_total"`
	CalculateTime       time.Time           `json:"calculate_time,format:unixmilli"`
	RiskUnit            []PortfolioRiskUnit `json:"risk_unit"`
}

// PortfolioRiskUnit is a calculated risk unit with its per-scenario margin results.
type PortfolioRiskUnit struct {
	Symbol         string                  `json:"symbol"`
	SpotInUse      decimal.Decimal         `json:"spot_in_use"`
	MaintainMargin decimal.Decimal         `json:"maintain_margin"`
	InitialMargin  decimal.Decimal         `json:"initial_margin"`
	MarginResult   []PortfolioMarginResult `json:"margin_result"`
	Delta          decimal.Decimal         `json:"delta"`
	Gamma          decimal.Decimal         `json:"gamma"`
	Theta          decimal.Decimal         `json:"theta"`
	Vega           decimal.Decimal         `json:"vega"`
}

// PortfolioMarginResult is the margin outcome for one position-combination type.
type PortfolioMarginResult struct {
	Type             string                     `json:"type"`
	ProfitLossRanges []PortfolioProfitLossRange `json:"profit_loss_ranges"`
	MaxLoss          PortfolioProfitLossRange   `json:"max_loss"`
	MR1              decimal.Decimal            `json:"mr1"`
	MR2              decimal.Decimal            `json:"mr2"`
	MR3              decimal.Decimal            `json:"mr3"`
	MR4              decimal.Decimal            `json:"mr4"`
}

// PortfolioProfitLossRange is one stress-scenario PnL point.
type PortfolioProfitLossRange struct {
	PricePercentage             decimal.Decimal `json:"price_percentage"`
	ImpliedVolatilityPercentage decimal.Decimal `json:"implied_volatility_percentage"`
	ProfitLoss                  decimal.Decimal `json:"profit_loss"`
}

// GetUserLeverageCurrencyConfigService -- GET /api/v4/unified/leverage/user_currency_config (private)
//
// Returns the min/max leverage the account may set for a borrow currency, and
// the borrowable amounts at the current leverage.
type GetUserLeverageCurrencyConfigService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUserLeverageCurrencyConfigService(currency string) *GetUserLeverageCurrencyConfigService {
	return &GetUserLeverageCurrencyConfigService{c: c, params: map[string]string{"currency": currency}}
}

func (s *GetUserLeverageCurrencyConfigService) Do(ctx context.Context) (*UserLeverageConfig, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/leverage/user_currency_config", s.params).WithSign()
	return request.Do[UserLeverageConfig](req)
}

// UserLeverageConfig is a borrow currency's leverage bounds and borrowable amounts.
type UserLeverageConfig struct {
	CurrentLeverage          decimal.Decimal `json:"current_leverage"`
	MinLeverage              decimal.Decimal `json:"min_leverage"`
	MaxLeverage              decimal.Decimal `json:"max_leverage"`
	Debit                    decimal.Decimal `json:"debit"`
	AvailableMargin          decimal.Decimal `json:"available_margin"`
	Borrowable               decimal.Decimal `json:"borrowable"`
	ExceptLeverageBorrowable decimal.Decimal `json:"except_leverage_borrowable"`
}

// GetUserLeverageCurrencySettingService -- GET /api/v4/unified/leverage/user_currency_setting (private)
//
// Returns the account's configured borrow-currency leverage multipliers,
// optionally filtered to a single currency.
type GetUserLeverageCurrencySettingService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetUserLeverageCurrencySettingService() *GetUserLeverageCurrencySettingService {
	return &GetUserLeverageCurrencySettingService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency.
func (s *GetUserLeverageCurrencySettingService) SetCurrency(currency string) *GetUserLeverageCurrencySettingService {
	s.params["currency"] = currency
	return s
}

func (s *GetUserLeverageCurrencySettingService) Do(ctx context.Context) ([]UserLeverageSetting, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/leverage/user_currency_setting", s.params).WithSign()
	resp, err := request.Do[[]UserLeverageSetting](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UserLeverageSetting is a borrow currency's configured leverage multiplier.
type UserLeverageSetting struct {
	Currency string          `json:"currency"`
	Leverage decimal.Decimal `json:"leverage"`
}

// SetUserLeverageCurrencySettingService -- POST /api/v4/unified/leverage/user_currency_setting (private)
//
// Sets the leverage multiplier used when borrowing a currency.
type SetUserLeverageCurrencySettingService struct {
	c    *UnifiedClient
	body map[string]any
}

func (c *UnifiedClient) NewSetUserLeverageCurrencySettingService(currency, leverage string) *SetUserLeverageCurrencySettingService {
	return &SetUserLeverageCurrencySettingService{c: c, body: map[string]any{
		"currency": currency,
		"leverage": leverage,
	}}
}

func (s *SetUserLeverageCurrencySettingService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/unified/leverage/user_currency_setting", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// SetUserLeverageService -- POST /api/v4/unified/leverage/user_setting (private)
//
// Sets the leverage for all of the user's borrowed currencies. Currencies with
// outstanding loans cannot be changed; a value above a currency's leverage limit
// is capped at that limit. Failures are not rolled back and affect only the
// currencies that failed — Do returns the currencies whose update failed together
// with the reason.
type SetUserLeverageService struct {
	c    *UnifiedClient
	body map[string]any
}

func (c *UnifiedClient) NewSetUserLeverageService(leverage string) *SetUserLeverageService {
	return &SetUserLeverageService{c: c, body: map[string]any{"leverage": leverage}}
}

func (s *SetUserLeverageService) Do(ctx context.Context) ([]LeverageFailedCurrency, error) {
	req := request.Post(ctx, s.c, "/api/v4/unified/leverage/user_setting", s.body).WithSign()
	resp, err := request.Do[[]LeverageFailedCurrency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LeverageFailedCurrency is a currency whose leverage update failed, with the reason.
type LeverageFailedCurrency struct {
	Currency string `json:"currency"`
	Reason   string `json:"reason"`
}

// ListUnifiedCurrenciesService -- GET /api/v4/unified/currencies
//
// Returns the currencies the unified account can borrow and their borrow limits.
type ListUnifiedCurrenciesService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewListUnifiedCurrenciesService() *ListUnifiedCurrenciesService {
	return &ListUnifiedCurrenciesService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency.
func (s *ListUnifiedCurrenciesService) SetCurrency(currency string) *ListUnifiedCurrenciesService {
	s.params["currency"] = currency
	return s
}

func (s *ListUnifiedCurrenciesService) Do(ctx context.Context) ([]UnifiedCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/currencies", s.params)
	resp, err := request.Do[[]UnifiedCurrency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnifiedCurrency is a borrowable currency and its borrow limits.
type UnifiedCurrency struct {
	Name                 string          `json:"name"`
	Prec                 decimal.Decimal `json:"prec"`
	MinBorrowAmount      decimal.Decimal `json:"min_borrow_amount"`
	UserMaxBorrowAmount  decimal.Decimal `json:"user_max_borrow_amount"`
	TotalMaxBorrowAmount decimal.Decimal `json:"total_max_borrow_amount"`
	LoanStatus           string          `json:"loan_status"`
}

// GetHistoryLoanRateService -- GET /api/v4/unified/history_loan_rate
//
// Returns a currency's historical hourly lending rates, most recent first.
type GetHistoryLoanRateService struct {
	c      *UnifiedClient
	params map[string]string
}

func (c *UnifiedClient) NewGetHistoryLoanRateService(currency string) *GetHistoryLoanRateService {
	return &GetHistoryLoanRateService{c: c, params: map[string]string{"currency": currency}}
}

// SetTier selects the VIP level whose floating rate to return.
func (s *GetHistoryLoanRateService) SetTier(tier string) *GetHistoryLoanRateService {
	s.params["tier"] = tier
	return s
}

// SetPage sets the page number.
func (s *GetHistoryLoanRateService) SetPage(page int) *GetHistoryLoanRateService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of rate points returned (default 100, max 100).
func (s *GetHistoryLoanRateService) SetLimit(limit int) *GetHistoryLoanRateService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryLoanRateService) Do(ctx context.Context) (*HistoryLoanRate, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/history_loan_rate", s.params)
	return request.Do[HistoryLoanRate](req)
}

// HistoryLoanRate is a currency's historical hourly lending rates.
type HistoryLoanRate struct {
	Currency   string                 `json:"currency"`
	Tier       string                 `json:"tier"`
	TierUpRate decimal.Decimal        `json:"tier_up_rate"`
	Rates      []HistoryLoanRateEntry `json:"rates"`
}

// HistoryLoanRateEntry is the lending rate for a single hour.
type HistoryLoanRateEntry struct {
	Time time.Time       `json:"time,format:unixmilli"`
	Rate decimal.Decimal `json:"rate"`
}

// SetUnifiedCollateralService -- POST /api/v4/unified/collateral_currencies (private)
//
// Sets which currencies are used as collateral. collateralType is 0 (all
// currencies) or 1 (custom); the enable/disable lists apply only when custom.
type SetUnifiedCollateralService struct {
	c    *UnifiedClient
	body map[string]any
}

func (c *UnifiedClient) NewSetUnifiedCollateralService(collateralType int) *SetUnifiedCollateralService {
	return &SetUnifiedCollateralService{c: c, body: map[string]any{"collateral_type": collateralType}}
}

// SetEnableList adds currencies to the collateral set (custom mode).
func (s *SetUnifiedCollateralService) SetEnableList(currencies []string) *SetUnifiedCollateralService {
	s.body["enable_list"] = currencies
	return s
}

// SetDisableList removes currencies from the collateral set (custom mode).
func (s *SetUnifiedCollateralService) SetDisableList(currencies []string) *SetUnifiedCollateralService {
	s.body["disable_list"] = currencies
	return s
}

func (s *SetUnifiedCollateralService) Do(ctx context.Context) (*UnifiedCollateralResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/unified/collateral_currencies", s.body).WithSign()
	return request.Do[UnifiedCollateralResult](req)
}

// UnifiedCollateralResult reports whether the collateral setting was applied.
type UnifiedCollateralResult struct {
	IsSuccess bool `json:"is_success"`
}
