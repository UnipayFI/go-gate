package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListAllFuturesContractsService -- GET /api/v4/futures/{settle}/contracts_all (public)
//
// Returns every perpetual-futures contract for a settlement currency, including
// contracts that have already been delisted. Unlike NewListFuturesContractsService
// (which only lists live contracts), this endpoint paginates over the full history.
type ListAllFuturesContractsService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListAllFuturesContractsService(settle Settle) *ListAllFuturesContractsService {
	return &ListAllFuturesContractsService{c: c, settle: settle, params: map[string]string{}}
}

// SetLimit caps the number of contracts returned in one page.
func (s *ListAllFuturesContractsService) SetLimit(limit int) *ListAllFuturesContractsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the first offset contracts (pagination, starting from 0).
func (s *ListAllFuturesContractsService) SetOffset(offset int) *ListAllFuturesContractsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListAllFuturesContractsService) Do(ctx context.Context) ([]FuturesContract, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/contracts_all", s.params)
	resp, err := request.Do[[]FuturesContract](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BatchFundingRatesService -- POST /api/v4/futures/{settle}/funding_rates (public)
//
// Batch-queries the historical funding-rate series of one or more perpetual
// contracts in a single request. The result groups the rate points per contract.
type BatchFundingRatesService struct {
	c         *FuturesClient
	settle    Settle
	contracts []string
}

func (c *FuturesClient) NewBatchFundingRatesService(settle Settle, contracts []string) *BatchFundingRatesService {
	return &BatchFundingRatesService{c: c, settle: settle, contracts: contracts}
}

func (s *BatchFundingRatesService) Do(ctx context.Context) ([]FuturesBatchFundingRate, error) {
	body := map[string]any{"contracts": s.contracts}
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/funding_rates", body)
	resp, err := request.Do[[]FuturesBatchFundingRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesBatchFundingRate is one contract's historical funding-rate series.
type FuturesBatchFundingRate struct {
	Contract string                    `json:"contract"`
	Data     []FuturesFundingRatePoint `json:"data"`
}

// FuturesFundingRatePoint is a single funding-rate observation. t is the funding
// time (unix seconds) and r is the funding rate applied at that time.
type FuturesFundingRatePoint struct {
	Time time.Time       `json:"t,format:unix"`
	Rate decimal.Decimal `json:"r"`
}
