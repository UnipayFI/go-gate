package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// Position is a single futures position (one contract, one settlement currency).
// size is signed: a positive value is a long, a negative value a short; there is
// no separate side field.
type Position struct {
	User                   int64               `json:"user"`
	Contract               string              `json:"contract"`
	Size                   int64               `json:"size"`
	Leverage               decimal.Decimal     `json:"leverage"`
	RiskLimit              decimal.Decimal     `json:"risk_limit"`
	LeverageMax            decimal.Decimal     `json:"leverage_max"`
	MaintenanceRate        decimal.Decimal     `json:"maintenance_rate"`
	Value                  decimal.Decimal     `json:"value"`
	Margin                 decimal.Decimal     `json:"margin"`
	EntryPrice             decimal.Decimal     `json:"entry_price"`
	LiqPrice               decimal.Decimal     `json:"liq_price"`
	MarkPrice              decimal.Decimal     `json:"mark_price"`
	InitialMargin          decimal.Decimal     `json:"initial_margin"`
	MaintenanceMargin      decimal.Decimal     `json:"maintenance_margin"`
	UnrealisedPnL          decimal.Decimal     `json:"unrealised_pnl"`
	RealisedPnL            decimal.Decimal     `json:"realised_pnl"`
	PnLPnL                 decimal.Decimal     `json:"pnl_pnl"`
	PnLFund                decimal.Decimal     `json:"pnl_fund"`
	PnLFee                 decimal.Decimal     `json:"pnl_fee"`
	HistoryPnL             decimal.Decimal     `json:"history_pnl"`
	LastClosePnL           decimal.Decimal     `json:"last_close_pnl"`
	RealisedPoint          decimal.Decimal     `json:"realised_point"`
	HistoryPoint           decimal.Decimal     `json:"history_point"`
	ADLRanking             int                 `json:"adl_ranking"`
	PendingOrders          int                 `json:"pending_orders"`
	CloseOrder             *PositionCloseOrder `json:"close_order"`
	Mode                   string              `json:"mode"`
	CrossLeverageLimit     decimal.Decimal     `json:"cross_leverage_limit"`
	UpdateTime             time.Time           `json:"update_time,format:unix"`
	UpdateID               int64               `json:"update_id"`
	OpenTime               time.Time           `json:"open_time,format:unix"`
	RiskLimitTable         string              `json:"risk_limit_table"`
	AverageMaintenanceRate decimal.Decimal     `json:"average_maintenance_rate"`
	PID                    int64               `json:"pid"`
	ID                     int64               `json:"id"`
	Lever                  decimal.Decimal     `json:"lever"`
	LiqLock                bool                `json:"liq_lock"`
	PosMarginMode          string              `json:"pos_margin_mode"`
	TradeLongSize          int64               `json:"trade_long_size"`
	TradeShortSize         int64               `json:"trade_short_size"`
	VoucherID              int64               `json:"voucher_id"`
	VoucherMargin          decimal.Decimal     `json:"voucher_margin"`
	VoucherSize            decimal.Decimal     `json:"voucher_size"`
}

// PositionCloseOrder is the pending close order attached to a position, or null
// when there is none.
type PositionCloseOrder struct {
	ID    int64           `json:"id"`
	Price decimal.Decimal `json:"price"`
	IsLiq bool            `json:"is_liq"`
}

// DualModeResult is the futures account snapshot returned when toggling dual
// (hedge) position mode. It mirrors the fields most relevant to the switch; the
// full account view lives on the account service.
type DualModeResult struct {
	Total                decimal.Decimal `json:"total"`
	UnrealisedPnL        decimal.Decimal `json:"unrealised_pnl"`
	PositionMargin       decimal.Decimal `json:"position_margin"`
	OrderMargin          decimal.Decimal `json:"order_margin"`
	Available            decimal.Decimal `json:"available"`
	Point                decimal.Decimal `json:"point"`
	Currency             string          `json:"currency"`
	InDualMode           bool            `json:"in_dual_mode"`
	PositionMode         string          `json:"position_mode"`
	EnableNewDualMode    bool            `json:"enable_new_dual_mode"`
	EnableEvolvedClassic bool            `json:"enable_evolved_classic"`
	MarginMode           int             `json:"margin_mode"`
}

// ListPositionsService -- GET /api/v4/futures/{settle}/positions (private)
//
// Lists the authenticated user's positions in the given settlement currency.
type ListPositionsService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListPositionsService(settle Settle) *ListPositionsService {
	return &ListPositionsService{c: c, settle: settle, params: map[string]string{}}
}

// SetHolding, when true, returns only contracts currently holding a position.
func (s *ListPositionsService) SetHolding(holding bool) *ListPositionsService {
	s.params["holding"] = strconv.FormatBool(holding)
	return s
}

// SetLimit caps the number of records returned; valid values are 1-100. When
// unset, there is no default and the full current position list is returned.
func (s *ListPositionsService) SetLimit(limit int) *ListPositionsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset skips the given number of records (pagination).
func (s *ListPositionsService) SetOffset(offset int) *ListPositionsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListPositionsService) Do(ctx context.Context) ([]Position, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions", s.params).WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetPositionService -- GET /api/v4/futures/{settle}/positions/{contract} (private)
//
// Returns a single position for one contract.
type GetPositionService struct {
	c        *FuturesClient
	settle   Settle
	contract string
}

func (c *FuturesClient) NewGetPositionService(settle Settle, contract string) *GetPositionService {
	return &GetPositionService{c: c, settle: settle, contract: contract}
}

func (s *GetPositionService) Do(ctx context.Context) (*Position, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions/"+s.contract).WithSign()
	return request.Do[Position](req)
}

// UpdatePositionMarginService -- POST /api/v4/futures/{settle}/positions/{contract}/margin (private)
//
// Adds or removes isolated margin on a position. change is a signed amount:
// positive increases the margin, negative decreases it.
type UpdatePositionMarginService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewUpdatePositionMarginService(settle Settle, contract, change string) *UpdatePositionMarginService {
	return &UpdatePositionMarginService{c: c, settle: settle, contract: contract, params: map[string]string{"change": change}}
}

func (s *UpdatePositionMarginService) Do(ctx context.Context) (*Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions/"+s.contract+"/margin")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	return request.Do[Position](req)
}

// UpdatePositionLeverageService -- POST /api/v4/futures/{settle}/positions/{contract}/leverage (private)
//
// Sets the leverage of a position. A leverage of "0" selects cross margin, in
// which case cross_leverage_limit bounds the effective leverage.
type UpdatePositionLeverageService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewUpdatePositionLeverageService(settle Settle, contract, leverage string) *UpdatePositionLeverageService {
	return &UpdatePositionLeverageService{c: c, settle: settle, contract: contract, params: map[string]string{"leverage": leverage}}
}

// SetCrossLeverageLimit bounds cross-margin leverage (valid only when leverage is 0).
func (s *UpdatePositionLeverageService) SetCrossLeverageLimit(limit string) *UpdatePositionLeverageService {
	s.params["cross_leverage_limit"] = limit
	return s
}

func (s *UpdatePositionLeverageService) Do(ctx context.Context) (*Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions/"+s.contract+"/leverage")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	return request.Do[Position](req)
}

// UpdatePositionRiskLimitService -- POST /api/v4/futures/{settle}/positions/{contract}/risk_limit (private)
//
// Sets the risk limit of a position.
type UpdatePositionRiskLimitService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewUpdatePositionRiskLimitService(settle Settle, contract, riskLimit string) *UpdatePositionRiskLimitService {
	return &UpdatePositionRiskLimitService{c: c, settle: settle, contract: contract, params: map[string]string{"risk_limit": riskLimit}}
}

func (s *UpdatePositionRiskLimitService) Do(ctx context.Context) (*Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions/"+s.contract+"/risk_limit")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	return request.Do[Position](req)
}

// UpdatePositionCrossModeService -- POST /api/v4/futures/{settle}/positions/cross_mode (private)
//
// Switches a contract between isolated ("ISOLATED") and cross ("CROSS") margin
// mode in one-way (single) position mode.
type UpdatePositionCrossModeService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewUpdatePositionCrossModeService(settle Settle, mode, contract string) *UpdatePositionCrossModeService {
	return &UpdatePositionCrossModeService{c: c, settle: settle, body: map[string]any{"mode": mode, "contract": contract}}
}

func (s *UpdatePositionCrossModeService) Do(ctx context.Context) (*Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions/cross_mode", s.body).WithSign()
	return request.Do[Position](req)
}

// SetDualModeService -- POST /api/v4/futures/{settle}/dual_mode (private)
//
// Enables or disables dual (hedge) position mode for the whole settlement
// currency. The switch is only allowed with no open positions or orders.
type SetDualModeService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewSetDualModeService(settle Settle, dualMode bool) *SetDualModeService {
	return &SetDualModeService{c: c, settle: settle, params: map[string]string{"dual_mode": strconv.FormatBool(dualMode)}}
}

func (s *SetDualModeService) Do(ctx context.Context) (*DualModeResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/dual_mode")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	return request.Do[DualModeResult](req)
}

// GetDualModePositionService -- GET /api/v4/futures/{settle}/dual_comp/positions/{contract} (private)
//
// Returns both legs (long and short) of a contract's position in dual mode.
type GetDualModePositionService struct {
	c        *FuturesClient
	settle   Settle
	contract string
}

func (c *FuturesClient) NewGetDualModePositionService(settle Settle, contract string) *GetDualModePositionService {
	return &GetDualModePositionService{c: c, settle: settle, contract: contract}
}

func (s *GetDualModePositionService) Do(ctx context.Context) ([]Position, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/dual_comp/positions/"+s.contract).WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UpdateDualModePositionMarginService -- POST /api/v4/futures/{settle}/dual_comp/positions/{contract}/margin (private)
//
// Adds or removes margin on one leg of a dual-mode position. dual_side selects
// the leg ("dual_long" or "dual_short"); change is signed.
type UpdateDualModePositionMarginService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewUpdateDualModePositionMarginService(settle Settle, contract, change, dualSide string) *UpdateDualModePositionMarginService {
	return &UpdateDualModePositionMarginService{
		c:        c,
		settle:   settle,
		contract: contract,
		params:   map[string]string{"change": change, "dual_side": dualSide},
	}
}

func (s *UpdateDualModePositionMarginService) Do(ctx context.Context) ([]Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/dual_comp/positions/"+s.contract+"/margin")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UpdateDualModePositionLeverageService -- POST /api/v4/futures/{settle}/dual_comp/positions/{contract}/leverage (private)
//
// Sets the leverage of both legs of a dual-mode position.
type UpdateDualModePositionLeverageService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewUpdateDualModePositionLeverageService(settle Settle, contract, leverage string) *UpdateDualModePositionLeverageService {
	return &UpdateDualModePositionLeverageService{c: c, settle: settle, contract: contract, params: map[string]string{"leverage": leverage}}
}

// SetCrossLeverageLimit bounds cross-margin leverage (valid only when leverage is 0).
func (s *UpdateDualModePositionLeverageService) SetCrossLeverageLimit(limit string) *UpdateDualModePositionLeverageService {
	s.params["cross_leverage_limit"] = limit
	return s
}

func (s *UpdateDualModePositionLeverageService) Do(ctx context.Context) ([]Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/dual_comp/positions/"+s.contract+"/leverage")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UpdateDualModePositionRiskLimitService -- POST /api/v4/futures/{settle}/dual_comp/positions/{contract}/risk_limit (private)
//
// Sets the risk limit of both legs of a dual-mode position.
type UpdateDualModePositionRiskLimitService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewUpdateDualModePositionRiskLimitService(settle Settle, contract, riskLimit string) *UpdateDualModePositionRiskLimitService {
	return &UpdateDualModePositionRiskLimitService{c: c, settle: settle, contract: contract, params: map[string]string{"risk_limit": riskLimit}}
}

func (s *UpdateDualModePositionRiskLimitService) Do(ctx context.Context) ([]Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/dual_comp/positions/"+s.contract+"/risk_limit")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UpdateDualCompPositionCrossModeService -- POST /api/v4/futures/{settle}/dual_comp/positions/cross_mode (private)
//
// Switches a contract between isolated ("ISOLATED") and cross ("CROSS") margin
// mode while in dual (hedge) position mode.
type UpdateDualCompPositionCrossModeService struct {
	c      *FuturesClient
	settle Settle
	body   map[string]any
}

func (c *FuturesClient) NewUpdateDualCompPositionCrossModeService(settle Settle, mode, contract string) *UpdateDualCompPositionCrossModeService {
	return &UpdateDualCompPositionCrossModeService{c: c, settle: settle, body: map[string]any{"mode": mode, "contract": contract}}
}

func (s *UpdateDualCompPositionCrossModeService) Do(ctx context.Context) ([]Position, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/dual_comp/positions/cross_mode", s.body).WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
