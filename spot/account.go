package spot

import (
	"context"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListSpotAccountsService -- GET /api/v4/spot/accounts (private)
//
// Returns the spot balances of the authenticated account, optionally filtered
// to a single currency.
type ListSpotAccountsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewListSpotAccountsService() *ListSpotAccountsService {
	return &ListSpotAccountsService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListSpotAccountsService) SetCurrency(currency string) *ListSpotAccountsService {
	s.params["currency"] = currency
	return s
}

func (s *ListSpotAccountsService) Do(ctx context.Context) ([]SpotAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/spot/accounts", s.params).WithSign()
	resp, err := request.Do[[]SpotAccount](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotAccount is a single currency's spot balance.
type SpotAccount struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
	UpdateID  int64           `json:"update_id"`
}
