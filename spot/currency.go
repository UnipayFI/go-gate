package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListCurrenciesService -- GET /api/v4/spot/currencies
//
// Returns every listed currency and its per-chain deposit/withdrawal status.
type ListCurrenciesService struct {
	c *SpotClient
}

func (c *SpotClient) NewListCurrenciesService() *ListCurrenciesService {
	return &ListCurrenciesService{c: c}
}

func (s *ListCurrenciesService) Do(ctx context.Context) ([]Currency, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/currencies")
	resp, err := request.Do[[]Currency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetCurrencyService -- GET /api/v4/spot/currencies/{currency}
//
// Returns the details of a single currency, including all of its chains.
type GetCurrencyService struct {
	c        *SpotClient
	currency string
}

func (c *SpotClient) NewGetCurrencyService(currency string) *GetCurrencyService {
	return &GetCurrencyService{c: c, currency: currency}
}

func (s *GetCurrencyService) Do(ctx context.Context) (*Currency, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/currencies/"+s.currency)
	return request.Do[Currency](req)
}

// Currency is a listed asset and its per-chain deposit/withdrawal metadata.
type Currency struct {
	Currency         string              `json:"currency"`
	Name             string              `json:"name"`
	Delisted         bool                `json:"delisted"`
	WithdrawDisabled bool                `json:"withdraw_disabled"`
	WithdrawDelayed  bool                `json:"withdraw_delayed"`
	DepositDisabled  bool                `json:"deposit_disabled"`
	TradeDisabled    bool                `json:"trade_disabled"`
	FixedRate        decimal.Decimal     `json:"fixed_rate"`
	Chain            string              `json:"chain"`
	Chains           []SpotCurrencyChain `json:"chains"`
	Category         []string            `json:"category"`
	// TotalSupply is a plain string because uncapped assets report it as "∞",
	// which is not a parseable decimal.
	TotalSupply string          `json:"total_supply"`
	MarketCap   decimal.Decimal `json:"market_cap"`
}

// SpotCurrencyChain is one blockchain a currency can be deposited/withdrawn on.
type SpotCurrencyChain struct {
	Name             string `json:"name"`
	Addr             string `json:"addr"`
	WithdrawDisabled bool   `json:"withdraw_disabled"`
	WithdrawDelayed  bool   `json:"withdraw_delayed"`
	DepositDisabled  bool   `json:"deposit_disabled"`
}

// GetSpotInsuranceHistoryService -- GET /api/v4/spot/insurance_history
//
// Returns the historical balance of the spot leverage insurance fund for a
// business over a time range.
type GetSpotInsuranceHistoryService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetSpotInsuranceHistoryService(business, currency string, from, to time.Time) *GetSpotInsuranceHistoryService {
	return &GetSpotInsuranceHistoryService{c: c, params: map[string]string{
		"business": business,
		"currency": currency,
		"from":     strconv.FormatInt(from.Unix(), 10),
		"to":       strconv.FormatInt(to.Unix(), 10),
	}}
}

// SetPage selects the page number of the paginated result.
func (s *GetSpotInsuranceHistoryService) SetPage(page int) *GetSpotInsuranceHistoryService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (default 30).
func (s *GetSpotInsuranceHistoryService) SetLimit(limit int) *GetSpotInsuranceHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSpotInsuranceHistoryService) Do(ctx context.Context) ([]SpotInsurance, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/insurance_history", s.params)
	resp, err := request.Do[[]SpotInsurance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotInsurance is a single snapshot of the insurance fund balance. Time is a
// millisecond epoch number.
type SpotInsurance struct {
	Currency string          `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
	Time     time.Time       `json:"time,format:unixmilli"`
}
