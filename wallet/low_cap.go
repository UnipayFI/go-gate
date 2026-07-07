package wallet

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
)

// ListLowCapExchangeService -- GET /api/v4/wallet/getLowCapExchangeList (private)
//
// Retrieves the list of low-liquidity / low-cap token symbols whose balances can
// be exchanged. The endpoint returns a bare array of currency symbols.
type ListLowCapExchangeService struct {
	c *WalletClient
}

func (c *WalletClient) NewListLowCapExchangeService() *ListLowCapExchangeService {
	return &ListLowCapExchangeService{c: c}
}

func (s *ListLowCapExchangeService) Do(ctx context.Context) ([]string, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/getLowCapExchangeList").WithSign()
	resp, err := request.Do[[]string](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
