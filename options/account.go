package options

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListOptionsAccountService -- GET /api/v4/options/accounts (private)
//
// Returns the options account balance, equity and margin snapshot of the
// authenticated user.
type ListOptionsAccountService struct {
	c *OptionsClient
}

func (c *OptionsClient) NewListOptionsAccountService() *ListOptionsAccountService {
	return &ListOptionsAccountService{c: c}
}

func (s *ListOptionsAccountService) Do(ctx context.Context) (*OptionsAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/accounts").WithSign()
	return request.Do[OptionsAccount](req)
}

// ListOptionsAccountBookService -- GET /api/v4/options/account_book (private)
//
// Lists the account's balance-change history (deposits/withdrawals, premiums,
// fees, rebates and settlement P&L).
type ListOptionsAccountBookService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsAccountBookService() *ListOptionsAccountBookService {
	return &ListOptionsAccountBookService{c: c, params: map[string]string{}}
}

// SetLimit caps the number of records returned in a single list.
func (s *ListOptionsAccountBookService) SetLimit(limit int) *ListOptionsAccountBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListOptionsAccountBookService) SetOffset(offset int) *ListOptionsAccountBookService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive, Unix seconds).
func (s *ListOptionsAccountBookService) SetFrom(from time.Time) *ListOptionsAccountBookService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive, Unix seconds).
func (s *ListOptionsAccountBookService) SetTo(to time.Time) *ListOptionsAccountBookService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetType filters by change type (dnw, prem, fee, refr, set, ...).
func (s *ListOptionsAccountBookService) SetType(changeType string) *ListOptionsAccountBookService {
	s.params["type"] = changeType
	return s
}

func (s *ListOptionsAccountBookService) Do(ctx context.Context) ([]OptionsAccountBook, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/account_book", s.params).WithSign()
	resp, err := request.Do[[]OptionsAccountBook](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListOptionsPositionsService -- GET /api/v4/options/positions (private)
//
// Lists the account's open options positions, optionally limited to a single
// underlying.
type ListOptionsPositionsService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsPositionsService() *ListOptionsPositionsService {
	return &ListOptionsPositionsService{c: c, params: map[string]string{}}
}

// SetUnderlying narrows the result to a single underlying.
func (s *ListOptionsPositionsService) SetUnderlying(underlying string) *ListOptionsPositionsService {
	s.params["underlying"] = underlying
	return s
}

func (s *ListOptionsPositionsService) Do(ctx context.Context) ([]OptionsPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/positions", s.params).WithSign()
	resp, err := request.Do[[]OptionsPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOptionsPositionService -- GET /api/v4/options/positions/{contract} (private)
//
// Returns the account's position in a single options contract.
type GetOptionsPositionService struct {
	c        *OptionsClient
	contract string
}

func (c *OptionsClient) NewGetOptionsPositionService(contract string) *GetOptionsPositionService {
	return &GetOptionsPositionService{c: c, contract: contract}
}

func (s *GetOptionsPositionService) Do(ctx context.Context) (*OptionsPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/positions/"+s.contract).WithSign()
	return request.Do[OptionsPosition](req)
}

// ListOptionsPositionCloseService -- GET /api/v4/options/position_close (private)
//
// Lists the account's position-close (liquidation) history for an underlying.
type ListOptionsPositionCloseService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListOptionsPositionCloseService(underlying string) *ListOptionsPositionCloseService {
	return &ListOptionsPositionCloseService{c: c, params: map[string]string{
		"underlying": underlying,
	}}
}

// SetContract narrows the result to a single options contract.
func (s *ListOptionsPositionCloseService) SetContract(contract string) *ListOptionsPositionCloseService {
	s.params["contract"] = contract
	return s
}

func (s *ListOptionsPositionCloseService) Do(ctx context.Context) ([]OptionsPositionClose, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/position_close", s.params).WithSign()
	resp, err := request.Do[[]OptionsPositionClose](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListMyOptionsSettlementsService -- GET /api/v4/options/my_settlements (private)
//
// Lists the account's personal settlement records for an underlying.
type ListMyOptionsSettlementsService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListMyOptionsSettlementsService(underlying string) *ListMyOptionsSettlementsService {
	return &ListMyOptionsSettlementsService{c: c, params: map[string]string{
		"underlying": underlying,
	}}
}

// SetContract narrows the result to a single options contract.
func (s *ListMyOptionsSettlementsService) SetContract(contract string) *ListMyOptionsSettlementsService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListMyOptionsSettlementsService) SetLimit(limit int) *ListMyOptionsSettlementsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListMyOptionsSettlementsService) SetOffset(offset int) *ListMyOptionsSettlementsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive, Unix seconds).
func (s *ListMyOptionsSettlementsService) SetFrom(from time.Time) *ListMyOptionsSettlementsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive, Unix seconds).
func (s *ListMyOptionsSettlementsService) SetTo(to time.Time) *ListMyOptionsSettlementsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListMyOptionsSettlementsService) Do(ctx context.Context) ([]OptionsMySettlement, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/my_settlements", s.params).WithSign()
	resp, err := request.Do[[]OptionsMySettlement](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListMyOptionsTradesService -- GET /api/v4/options/my_trades (private)
//
// Lists the account's personal options trading (fill) records for an underlying.
type ListMyOptionsTradesService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewListMyOptionsTradesService(underlying string) *ListMyOptionsTradesService {
	return &ListMyOptionsTradesService{c: c, params: map[string]string{
		"underlying": underlying,
	}}
}

// SetContract narrows the result to a single options contract.
func (s *ListMyOptionsTradesService) SetContract(contract string) *ListMyOptionsTradesService {
	s.params["contract"] = contract
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListMyOptionsTradesService) SetLimit(limit int) *ListMyOptionsTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListMyOptionsTradesService) SetOffset(offset int) *ListMyOptionsTradesService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetFrom sets the start time (inclusive, Unix seconds).
func (s *ListMyOptionsTradesService) SetFrom(from time.Time) *ListMyOptionsTradesService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive, Unix seconds).
func (s *ListMyOptionsTradesService) SetTo(to time.Time) *ListMyOptionsTradesService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *ListMyOptionsTradesService) Do(ctx context.Context) ([]OptionsMyTrade, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/my_trades", s.params).WithSign()
	resp, err := request.Do[[]OptionsMyTrade](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetOptionsMMPService -- GET /api/v4/options/mmp (private)
//
// Returns the account's Market Maker Protection (MMP) settings for an underlying.
type GetOptionsMMPService struct {
	c      *OptionsClient
	params map[string]string
}

func (c *OptionsClient) NewGetOptionsMMPService(underlying string) *GetOptionsMMPService {
	return &GetOptionsMMPService{c: c, params: map[string]string{
		"underlying": underlying,
	}}
}

func (s *GetOptionsMMPService) Do(ctx context.Context) (*OptionsMMP, error) {
	req := request.Get(ctx, s.c, "/api/v4/options/mmp", s.params).WithSign()
	return request.Do[OptionsMMP](req)
}

// SetOptionsMMPService -- POST /api/v4/options/mmp (private)
//
// Configures the account's Market Maker Protection for an underlying: a rolling
// window (ms), a freeze duration (ms), and the trade-quantity / net-delta limits
// that trip the protection.
type SetOptionsMMPService struct {
	c    *OptionsClient
	body map[string]any
}

func (c *OptionsClient) NewSetOptionsMMPService(underlying string, window, frozenPeriod int, qtyLimit, deltaLimit decimal.Decimal) *SetOptionsMMPService {
	return &SetOptionsMMPService{c: c, body: map[string]any{
		"underlying":    underlying,
		"window":        window,
		"frozen_period": frozenPeriod,
		"qty_limit":     qtyLimit.String(),
		"delta_limit":   deltaLimit.String(),
	}}
}

func (s *SetOptionsMMPService) Do(ctx context.Context) (*OptionsMMP, error) {
	req := request.Post(ctx, s.c, "/api/v4/options/mmp", s.body).WithSign()
	return request.Do[OptionsMMP](req)
}

// ResetOptionsMMPService -- POST /api/v4/options/mmp/reset (private)
//
// Unfreezes the account's Market Maker Protection for an underlying after it has
// been triggered.
type ResetOptionsMMPService struct {
	c    *OptionsClient
	body map[string]any
}

func (c *OptionsClient) NewResetOptionsMMPService(underlying string) *ResetOptionsMMPService {
	return &ResetOptionsMMPService{c: c, body: map[string]any{
		"underlying": underlying,
	}}
}

func (s *ResetOptionsMMPService) Do(ctx context.Context) (*OptionsMMP, error) {
	req := request.Post(ctx, s.c, "/api/v4/options/mmp/reset", s.body).WithSign()
	return request.Do[OptionsMMP](req)
}

// OptionsAccount is the options account balance, equity and margin snapshot.
// position_value is signed (positive for long, negative for short); margin_mode
// is a code (0 classic-spot, 1 cross-currency, 2 portfolio).
type OptionsAccount struct {
	User                  int64           `json:"user"`
	Total                 decimal.Decimal `json:"total"`
	PositionValue         decimal.Decimal `json:"position_value"`
	Equity                decimal.Decimal `json:"equity"`
	ShortEnabled          bool            `json:"short_enabled"`
	MMPEnabled            bool            `json:"mmp_enabled"`
	LiqTriggered          bool            `json:"liq_triggered"`
	MarginMode            int             `json:"margin_mode"`
	UnrealisedPnL         decimal.Decimal `json:"unrealised_pnl"`
	InitMargin            decimal.Decimal `json:"init_margin"`
	MaintMargin           decimal.Decimal `json:"maint_margin"`
	OrderMargin           decimal.Decimal `json:"order_margin"`
	AskOrderMargin        decimal.Decimal `json:"ask_order_margin"`
	BidOrderMargin        decimal.Decimal `json:"bid_order_margin"`
	Available             decimal.Decimal `json:"available"`
	Point                 decimal.Decimal `json:"point"`
	Currency              string          `json:"currency"`
	OrdersLimit           int             `json:"orders_limit"`
	PositionNotionalLimit decimal.Decimal `json:"position_notional_limit"`
}

// OptionsAccountBook is a single balance-change record. time is a float-second
// Unix timestamp; type is the change category (dnw, prem, fee, refr, set, ...).
type OptionsAccountBook struct {
	Time    time.Time       `json:"time,format:unix"`
	Change  decimal.Decimal `json:"change"`
	Balance decimal.Decimal `json:"balance"`
	Type    string          `json:"type"`
	Text    string          `json:"text"`
}

// OptionsPosition is the account's position in a single options contract. size
// is signed (positive long, negative short); mark_iv is the mark implied
// volatility, and delta/gamma/vega/theta are the position greeks.
type OptionsPosition struct {
	User            int64                      `json:"user"`
	Underlying      string                     `json:"underlying"`
	UnderlyingPrice decimal.Decimal            `json:"underlying_price"`
	Contract        string                     `json:"contract"`
	Size            int64                      `json:"size"`
	EntryPrice      decimal.Decimal            `json:"entry_price"`
	MarkPrice       decimal.Decimal            `json:"mark_price"`
	MarkIV          decimal.Decimal            `json:"mark_iv"`
	RealisedPnL     decimal.Decimal            `json:"realised_pnl"`
	UnrealisedPnL   decimal.Decimal            `json:"unrealised_pnl"`
	PendingOrders   int                        `json:"pending_orders"`
	CloseOrder      *OptionsPositionCloseOrder `json:"close_order"`
	Delta           decimal.Decimal            `json:"delta"`
	Gamma           decimal.Decimal            `json:"gamma"`
	Vega            decimal.Decimal            `json:"vega"`
	Theta           decimal.Decimal            `json:"theta"`
}

// OptionsPositionCloseOrder is the current close order attached to a position,
// or null when none is open.
type OptionsPositionCloseOrder struct {
	ID    int64           `json:"id"`
	Price decimal.Decimal `json:"price"`
	IsLiq bool            `json:"is_liq"`
}

// OptionsPositionClose is a position-close (liquidation) record. time is a
// float-second Unix timestamp; side is "long" or "short".
type OptionsPositionClose struct {
	Time       time.Time       `json:"time,format:unix"`
	Contract   string          `json:"contract"`
	Side       string          `json:"side"`
	PnL        decimal.Decimal `json:"pnl"`
	Text       string          `json:"text"`
	SettleSize decimal.Decimal `json:"settle_size"`
}

// OptionsMySettlement is a personal settlement record. time is a float-second
// Unix timestamp; realised_pnl is the accumulated P&L including premium, fees
// and settlement profit.
type OptionsMySettlement struct {
	Time         time.Time       `json:"time,format:unix"`
	Underlying   string          `json:"underlying"`
	Contract     string          `json:"contract"`
	StrikePrice  decimal.Decimal `json:"strike_price"`
	SettlePrice  decimal.Decimal `json:"settle_price"`
	Size         int64           `json:"size"`
	SettleProfit decimal.Decimal `json:"settle_profit"`
	Fee          decimal.Decimal `json:"fee"`
	RealisedPnL  decimal.Decimal `json:"realised_pnl"`
}

// OptionsMyTrade is a personal options fill. create_time is a float-second Unix
// timestamp; role is "taker" or "maker".
type OptionsMyTrade struct {
	ID              int64           `json:"id"`
	CreateTime      time.Time       `json:"create_time,format:unix"`
	Contract        string          `json:"contract"`
	OrderID         int64           `json:"order_id"`
	Size            int64           `json:"size"`
	Price           decimal.Decimal `json:"price"`
	UnderlyingPrice decimal.Decimal `json:"underlying_price"`
	Role            string          `json:"role"`
}

// OptionsMMP is the Market Maker Protection configuration and state. window and
// frozen_period are durations in milliseconds (window 0 disables MMP), while
// trigger_time_ms / frozen_until_ms are millisecond-epoch moments (0 when the
// protection has not been triggered / has no unfreeze time).
type OptionsMMP struct {
	Underlying    string          `json:"underlying"`
	Window        int             `json:"window"`
	FrozenPeriod  int             `json:"frozen_period"`
	QtyLimit      decimal.Decimal `json:"qty_limit"`
	DeltaLimit    decimal.Decimal `json:"delta_limit"`
	TriggerTimeMs time.Time       `json:"trigger_time_ms,format:unixmilli"`
	FrozenUntilMs time.Time       `json:"frozen_until_ms,format:unixmilli"`
}
