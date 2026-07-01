package delivery

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListDeliveryAccountsService -- GET /api/v4/delivery/{settle}/accounts (private)
//
// Returns the delivery (dated-futures) account balance and margin breakdown for
// the given settlement currency.
type ListDeliveryAccountsService struct {
	c      *DeliveryClient
	settle Settle
}

func (c *DeliveryClient) NewListDeliveryAccountsService(settle Settle) *ListDeliveryAccountsService {
	return &ListDeliveryAccountsService{c: c, settle: settle}
}

func (s *ListDeliveryAccountsService) Do(ctx context.Context) (*DeliveryAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/accounts").WithSign()
	return request.Do[DeliveryAccount](req)
}

// DeliveryAccount is the delivery account balance and margin breakdown.
type DeliveryAccount struct {
	Total                  decimal.Decimal        `json:"total"`
	UnrealisedPnL          decimal.Decimal        `json:"unrealised_pnl"`
	PositionMargin         decimal.Decimal        `json:"position_margin"`
	OrderMargin            decimal.Decimal        `json:"order_margin"`
	Available              decimal.Decimal        `json:"available"`
	Point                  decimal.Decimal        `json:"point"`
	Currency               string                 `json:"currency"`
	InDualMode             bool                   `json:"in_dual_mode"`
	PositionMode           string                 `json:"position_mode"`
	EnableCredit           bool                   `json:"enable_credit"`
	PositionInitialMargin  decimal.Decimal        `json:"position_initial_margin"`
	MaintenanceMargin      decimal.Decimal        `json:"maintenance_margin"`
	Bonus                  decimal.Decimal        `json:"bonus"`
	EnableEvolvedClassic   bool                   `json:"enable_evolved_classic"`
	CrossOrderMargin       decimal.Decimal        `json:"cross_order_margin"`
	CrossInitialMargin     decimal.Decimal        `json:"cross_initial_margin"`
	CrossMaintenanceMargin decimal.Decimal        `json:"cross_maintenance_margin"`
	CrossUnrealisedPnL     decimal.Decimal        `json:"cross_unrealised_pnl"`
	CrossAvailable         decimal.Decimal        `json:"cross_available"`
	CrossMarginBalance     decimal.Decimal        `json:"cross_margin_balance"`
	CrossMMR               decimal.Decimal        `json:"cross_mmr"`
	CrossIMR               decimal.Decimal        `json:"cross_imr"`
	IsolatedPositionMargin decimal.Decimal        `json:"isolated_position_margin"`
	EnableNewDualMode      bool                   `json:"enable_new_dual_mode"`
	MarginMode             int                    `json:"margin_mode"`
	EnableTieredMM         bool                   `json:"enable_tiered_mm"`
	UpdateTime             time.Time              `json:"update_time,format:unix"`
	History                DeliveryAccountHistory `json:"history"`
}

// DeliveryAccountHistory is the account's cumulative deposit/withdraw and
// realized profit-and-loss statistics.
type DeliveryAccountHistory struct {
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
}

// ListDeliveryAccountBookService -- GET /api/v4/delivery/{settle}/account_book (private)
//
// Returns the delivery account balance change history (deposits, withdrawals,
// realized P&L, fees, funding and referral rebates).
type ListDeliveryAccountBookService struct {
	c      *DeliveryClient
	settle Settle
	params map[string]string
}

func (c *DeliveryClient) NewListDeliveryAccountBookService(settle Settle) *ListDeliveryAccountBookService {
	return &ListDeliveryAccountBookService{c: c, settle: settle, params: map[string]string{}}
}

// SetLimit caps the number of records returned in a single list.
func (s *ListDeliveryAccountBookService) SetLimit(limit int) *ListDeliveryAccountBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetFrom sets the start time of the query window.
func (s *ListDeliveryAccountBookService) SetFrom(from time.Time) *ListDeliveryAccountBookService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window (defaults to now).
func (s *ListDeliveryAccountBookService) SetTo(to time.Time) *ListDeliveryAccountBookService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetType narrows the result to a single change type (e.g. dnw, pnl, fee, refr,
// fund, point_dnw, point_fee, point_refr).
func (s *ListDeliveryAccountBookService) SetType(changeType string) *ListDeliveryAccountBookService {
	s.params["type"] = changeType
	return s
}

func (s *ListDeliveryAccountBookService) Do(ctx context.Context) ([]DeliveryAccountBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/account_book", s.params).WithSign()
	resp, err := request.Do[[]DeliveryAccountBook](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeliveryAccountBook is a single delivery account balance change record.
type DeliveryAccountBook struct {
	Time     time.Time       `json:"time,format:unix"`
	Change   decimal.Decimal `json:"change"`
	Balance  decimal.Decimal `json:"balance"`
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Contract string          `json:"contract"`
	TradeID  string          `json:"trade_id"`
	ID       string          `json:"id"`
}
