package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListFuturesAccountsService -- GET /api/v4/futures/{settle}/accounts (private)
//
// Returns the authenticated futures account balance and margin summary for a
// settlement currency.
type ListFuturesAccountsService struct {
	c      *FuturesClient
	settle Settle
}

func (c *FuturesClient) NewListFuturesAccountsService(settle Settle) *ListFuturesAccountsService {
	return &ListFuturesAccountsService{c: c, settle: settle}
}

func (s *ListFuturesAccountsService) Do(ctx context.Context) (*FuturesAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/accounts").WithSign()
	return request.Do[FuturesAccount](req)
}

// FuturesAccount is the futures wallet balance and margin breakdown.
type FuturesAccount struct {
	User                   int64                 `json:"user"`
	Total                  decimal.Decimal       `json:"total"`
	UnrealisedPnL          decimal.Decimal       `json:"unrealised_pnl"`
	PositionMargin         decimal.Decimal       `json:"position_margin"`
	OrderMargin            decimal.Decimal       `json:"order_margin"`
	Available              decimal.Decimal       `json:"available"`
	Point                  decimal.Decimal       `json:"point"`
	Currency               string                `json:"currency"`
	InDualMode             bool                  `json:"in_dual_mode"`
	PositionMode           string                `json:"position_mode"`
	EnableCredit           bool                  `json:"enable_credit"`
	PositionInitialMargin  decimal.Decimal       `json:"position_initial_margin"`
	MaintenanceMargin      decimal.Decimal       `json:"maintenance_margin"`
	Bonus                  decimal.Decimal       `json:"bonus"`
	EnableEvolvedClassic   bool                  `json:"enable_evolved_classic"`
	CrossOrderMargin       decimal.Decimal       `json:"cross_order_margin"`
	CrossInitialMargin     decimal.Decimal       `json:"cross_initial_margin"`
	CrossMaintenanceMargin decimal.Decimal       `json:"cross_maintenance_margin"`
	CrossUnrealisedPnL     decimal.Decimal       `json:"cross_unrealised_pnl"`
	CrossAvailable         decimal.Decimal       `json:"cross_available"`
	CrossMarginBalance     decimal.Decimal       `json:"cross_margin_balance"`
	CrossMMR               decimal.Decimal       `json:"cross_mmr"`
	CrossIMR               decimal.Decimal       `json:"cross_imr"`
	IsolatedPositionMargin decimal.Decimal       `json:"isolated_position_margin"`
	EnableNewDualMode      bool                  `json:"enable_new_dual_mode"`
	MarginMode             int                   `json:"margin_mode"`
	MarginModeName         string                `json:"margin_mode_name"`
	PositionVoucherTotal   decimal.Decimal       `json:"position_voucher_total"`
	EnableTieredMM         bool                  `json:"enable_tiered_mm"`
	UpdateTime             time.Time             `json:"update_time,format:unix"`
	UpdateID               int64                 `json:"update_id"`
	History                FuturesAccountHistory `json:"history"`
}

// FuturesAccountHistory is the cumulative statistics that make up the balance.
type FuturesAccountHistory struct {
	Dnw         decimal.Decimal `json:"dnw"`
	PnL         decimal.Decimal `json:"pnl"`
	Fee         decimal.Decimal `json:"fee"`
	Refr        decimal.Decimal `json:"refr"`
	Fund        decimal.Decimal `json:"fund"`
	PointDnw    decimal.Decimal `json:"point_dnw"`
	PointFee    decimal.Decimal `json:"point_fee"`
	PointRefr   decimal.Decimal `json:"point_refr"`
	BonusDnw    decimal.Decimal `json:"bonus_dnw"`
	BonusOffset decimal.Decimal `json:"bonus_offset"`
	CrossSettle decimal.Decimal `json:"cross_settle"`
}

// ListFuturesAccountBookService -- GET /api/v4/futures/{settle}/account_book (private)
//
// Returns the account balance change history (deposits, PnL, fees, funding).
type ListFuturesAccountBookService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListFuturesAccountBookService(settle Settle) *ListFuturesAccountBookService {
	return &ListFuturesAccountBookService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract limits the history to a single contract (data after 2023-10-30).
func (s *ListFuturesAccountBookService) SetContract(contract string) *ListFuturesAccountBookService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListFuturesAccountBookService) SetLimit(limit int) *ListFuturesAccountBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListFuturesAccountBookService) SetOffset(offset int) *ListFuturesAccountBookService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start of the query window (unix seconds).
func (s *ListFuturesAccountBookService) SetFrom(from time.Time) *ListFuturesAccountBookService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query window (unix seconds).
func (s *ListFuturesAccountBookService) SetTo(to time.Time) *ListFuturesAccountBookService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetType filters by change type (dnw, pnl, fee, refr, fund, point_dnw,
// point_fee, point_refr, bonus_offset).
func (s *ListFuturesAccountBookService) SetType(changeType string) *ListFuturesAccountBookService {
	s.params["type"] = changeType
	return s
}

func (s *ListFuturesAccountBookService) Do(ctx context.Context) ([]FuturesAccountBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/account_book", s.params).WithSign()
	resp, err := request.Do[[]FuturesAccountBook](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesAccountBook is a single balance-change record.
type FuturesAccountBook struct {
	Time     time.Time       `json:"time,format:unix"`
	Change   decimal.Decimal `json:"change"`
	Balance  decimal.Decimal `json:"balance"`
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Contract string          `json:"contract"`
	TradeID  string          `json:"trade_id"`
	ID       string          `json:"id"`
	UpdateID string          `json:"update_id"`
	BizInfo  string          `json:"biz_info"`
}

// GetFuturesFeeService -- GET /api/v4/futures/{settle}/fee (private)
//
// Returns the account's taker/maker fee rates, keyed by contract.
type GetFuturesFeeService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewGetFuturesFeeService(settle Settle) *GetFuturesFeeService {
	return &GetFuturesFeeService{c: c, settle: settle, params: map[string]string{}}
}

// SetContract narrows the fee rates to a single contract.
func (s *GetFuturesFeeService) SetContract(contract string) *GetFuturesFeeService {
	s.params["contract"] = contract
	return s
}

func (s *GetFuturesFeeService) Do(ctx context.Context) (map[string]FuturesFee, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/fee", s.params).WithSign()
	resp, err := request.Do[map[string]FuturesFee](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesFee is the taker/maker fee rate for one contract.
type FuturesFee struct {
	TakerFee    decimal.Decimal `json:"taker_fee"`
	MakerFee    decimal.Decimal `json:"maker_fee"`
	RPIMakerFee decimal.Decimal `json:"rpi_maker_fee"`
}
