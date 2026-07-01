package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListFuturesInsuranceLedgerService -- GET /api/v4/futures/{settle}/insurance
//
// Returns the futures market insurance-fund balance history.
type ListFuturesInsuranceLedgerService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesInsuranceLedgerService(settle Settle) *ListFuturesInsuranceLedgerService {
	return &ListFuturesInsuranceLedgerService{c: c, settle: settle, params: map[string]string{}}
}

// SetLimit caps the number of records returned in a single list.
func (s *ListFuturesInsuranceLedgerService) SetLimit(limit int) *ListFuturesInsuranceLedgerService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListFuturesInsuranceLedgerService) Do(ctx context.Context) ([]FuturesInsurance, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/insurance", s.params)
	resp, err := request.Do[[]FuturesInsurance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesInsurance is one snapshot of the insurance-fund balance.
type FuturesInsurance struct {
	Timestamp time.Time       `json:"t,format:unix"`
	Balance   decimal.Decimal `json:"b"`
}

// ListContractStatsService -- GET /api/v4/futures/{settle}/contract_stats
//
// Returns time-series trading statistics (long/short ratios, liquidations, open
// interest) for a single contract.
type ListContractStatsService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListContractStatsService(settle Settle, contract string) *ListContractStatsService {
	return &ListContractStatsService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetFrom sets the start timestamp of the query window.
func (s *ListContractStatsService) SetFrom(from time.Time) *ListContractStatsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetInterval selects the sampling interval (e.g. "5m", "15m", "1h", "4h", "1d").
func (s *ListContractStatsService) SetInterval(interval string) *ListContractStatsService {
	s.params["interval"] = interval
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListContractStatsService) SetLimit(limit int) *ListContractStatsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListContractStatsService) Do(ctx context.Context) ([]ContractStat, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/contract_stats", s.params)
	resp, err := request.Do[[]ContractStat](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ContractStat is one sampling interval of a contract's trading statistics.
//
// Size fields are contract counts, but for enable_decimal contracts (e.g.
// ETH_USDT) the live API returns them as fractional numbers (e.g.
// "long_liq_size":280.8), so they are decoded as decimal.Decimal rather than
// int64. Only the account/user counters stay int64. The *_account / lsr fields
// are ratios and remain decimal.
type ContractStat struct {
	Time            time.Time       `json:"time,format:unix"`
	LsrTaker        decimal.Decimal `json:"lsr_taker"`
	LsrAccount      decimal.Decimal `json:"lsr_account"`
	LongLiqSize     decimal.Decimal `json:"long_liq_size"`
	ShortLiqSize    decimal.Decimal `json:"short_liq_size"`
	OpenInterest    decimal.Decimal `json:"open_interest"`
	ShortLiqUSD     decimal.Decimal `json:"short_liq_usd"`
	MarkPrice       decimal.Decimal `json:"mark_price"`
	TopLsrSize      decimal.Decimal `json:"top_lsr_size"`
	TopLongSize     decimal.Decimal `json:"top_long_size"`
	TopShortSize    decimal.Decimal `json:"top_short_size"`
	ShortLiqAmount  decimal.Decimal `json:"short_liq_amount"`
	LongLiqAmount   decimal.Decimal `json:"long_liq_amount"`
	ShortLiqUSDNew  decimal.Decimal `json:"short_liq_usd_new"`
	LongLiqUSDNew   decimal.Decimal `json:"long_liq_usd_new"`
	OpenInterestUSD decimal.Decimal `json:"open_interest_usd"`
	TopLsrAccount   decimal.Decimal `json:"top_lsr_account"`
	TopLongAccount  int64           `json:"top_long_account"`
	TopShortAccount int64           `json:"top_short_account"`
	LongLiqUSD      decimal.Decimal `json:"long_liq_usd"`
	LongTakerSize   decimal.Decimal `json:"long_taker_size"`
	ShortTakerSize  decimal.Decimal `json:"short_taker_size"`
	LongUsers       int64           `json:"long_users"`
	ShortUsers      int64           `json:"short_users"`
}

// GetIndexConstituentsService -- GET /api/v4/futures/{settle}/index_constituents/{index}
//
// Returns the exchanges and symbols that make up a futures price index.
type GetIndexConstituentsService struct {
	c      *FuturesClient
	settle Settle
	index  string
}

func (c *FuturesClient) NewGetIndexConstituentsService(settle Settle, index string) *GetIndexConstituentsService {
	return &GetIndexConstituentsService{c: c, settle: settle, index: index}
}

func (s *GetIndexConstituentsService) Do(ctx context.Context) (*IndexConstituents, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/index_constituents/"+s.index)
	return request.Do[IndexConstituents](req)
}

// IndexConstituents lists the reference exchanges behind a price index.
type IndexConstituents struct {
	Index        string             `json:"index"`
	Constituents []IndexConstituent `json:"constituents"`
}

// IndexConstituent is one exchange's contribution to a price index.
type IndexConstituent struct {
	Exchange string          `json:"exchange"`
	Symbols  []string        `json:"symbols"`
	Weight   decimal.Decimal `json:"weight"`
}

// ListLiquidatedOrdersService -- GET /api/v4/futures/{settle}/liq_orders
//
// Returns the public liquidation-order history. The from/to window spans at most
// one hour.
type ListLiquidatedOrdersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListLiquidatedOrdersService(settle Settle) *ListLiquidatedOrdersService {
	return &ListLiquidatedOrdersService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the result to a single contract.
func (s *ListLiquidatedOrdersService) SetContract(contract string) *ListLiquidatedOrdersService {
	s.params["contract"] = contract
	return s
}

// SetFrom sets the start timestamp of the query window.
func (s *ListLiquidatedOrdersService) SetFrom(from time.Time) *ListLiquidatedOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end timestamp of the query window.
func (s *ListLiquidatedOrdersService) SetTo(to time.Time) *ListLiquidatedOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListLiquidatedOrdersService) SetLimit(limit int) *ListLiquidatedOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListLiquidatedOrdersService) Do(ctx context.Context) ([]FuturesLiquidate, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/liq_orders", s.params)
	resp, err := request.Do[[]FuturesLiquidate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesLiquidate is one public liquidation order.
type FuturesLiquidate struct {
	Time       time.Time       `json:"time,format:unix"`
	Contract   string          `json:"contract"`
	Size       decimal.Decimal `json:"size"`
	OrderSize  decimal.Decimal `json:"order_size"`
	OrderPrice decimal.Decimal `json:"order_price"`
	FillPrice  decimal.Decimal `json:"fill_price"`
	Left       decimal.Decimal `json:"left"`
}

// ListFuturesRiskLimitTiersService -- GET /api/v4/futures/{settle}/risk_limit_tiers
//
// Returns the gradient risk-limit tiers. Without a contract, limit/offset paginate
// across the top markets and each row carries its contract.
type ListFuturesRiskLimitTiersService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesRiskLimitTiersService(settle Settle) *ListFuturesRiskLimitTiersService {
	return &ListFuturesRiskLimitTiersService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract returns tiers for a single contract only.
func (s *ListFuturesRiskLimitTiersService) SetContract(contract string) *ListFuturesRiskLimitTiersService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of market-level records returned (contract param empty).
func (s *ListFuturesRiskLimitTiersService) SetLimit(limit int) *ListFuturesRiskLimitTiersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the market-level list offset, starting from 0 (contract param empty).
func (s *ListFuturesRiskLimitTiersService) SetOffset(offset int) *ListFuturesRiskLimitTiersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListFuturesRiskLimitTiersService) Do(ctx context.Context) ([]FuturesRiskLimitTier, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/risk_limit_tiers", s.params)
	resp, err := request.Do[[]FuturesRiskLimitTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesRiskLimitTier is one tier of a contract's gradient risk-limit table.
type FuturesRiskLimitTier struct {
	Tier            int             `json:"tier"`
	RiskLimit       decimal.Decimal `json:"risk_limit"`
	InitialRate     decimal.Decimal `json:"initial_rate"`
	MaintenanceRate decimal.Decimal `json:"maintenance_rate"`
	LeverageMax     decimal.Decimal `json:"leverage_max"`
	Contract        string          `json:"contract"`
	Deduction       decimal.Decimal `json:"deduction"`
}

// GetFuturesRiskLimitTableService -- GET /api/v4/futures/{settle}/risk_limit_table
//
// Returns every tier of a named risk-limit table (table_id, e.g. "BTC_USDT_20241122").
type GetFuturesRiskLimitTableService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetFuturesRiskLimitTableService(settle Settle, tableID string) *GetFuturesRiskLimitTableService {
	return &GetFuturesRiskLimitTableService{c: c, settle: settle, params: map[string]string{"table_id": tableID}}
}

func (s *GetFuturesRiskLimitTableService) Do(ctx context.Context) ([]FuturesRiskLimitTable, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/risk_limit_table", s.params)
	resp, err := request.Do[[]FuturesRiskLimitTable](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesRiskLimitTable is one tier row of a named risk-limit table.
type FuturesRiskLimitTable struct {
	Tier            int             `json:"tier"`
	RiskLimit       decimal.Decimal `json:"risk_limit"`
	InitialRate     decimal.Decimal `json:"initial_rate"`
	MaintenanceRate decimal.Decimal `json:"maintenance_rate"`
	LeverageMax     decimal.Decimal `json:"leverage_max"`
	Deduction       decimal.Decimal `json:"deduction"`
}
