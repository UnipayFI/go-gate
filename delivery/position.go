package delivery

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// ListDeliveryPositionsService -- GET /api/v4/delivery/{settle}/positions (private)
//
// Returns every delivery position held by the authenticated account for a
// settlement currency.
type ListDeliveryPositionsService struct {
	c      *DeliveryClient
	settle Settle
}

func (c *DeliveryClient) NewListDeliveryPositionsService(settle Settle) *ListDeliveryPositionsService {
	return &ListDeliveryPositionsService{c: c, settle: settle}
}

func (s *ListDeliveryPositionsService) Do(ctx context.Context) ([]DeliveryPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/positions").WithSign()
	resp, err := request.Do[[]DeliveryPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetDeliveryPositionService -- GET /api/v4/delivery/{settle}/positions/{contract} (private)
//
// Returns the authenticated account's position in a single delivery contract.
type GetDeliveryPositionService struct {
	c        *DeliveryClient
	settle   Settle
	contract string
}

func (c *DeliveryClient) NewGetDeliveryPositionService(settle Settle, contract string) *GetDeliveryPositionService {
	return &GetDeliveryPositionService{c: c, settle: settle, contract: contract}
}

func (s *GetDeliveryPositionService) Do(ctx context.Context) (*DeliveryPosition, error) {
	req := request.Get(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/positions/"+s.contract).WithSign()
	return request.Do[DeliveryPosition](req)
}

// UpdateDeliveryPositionMarginService -- POST /api/v4/delivery/{settle}/positions/{contract}/margin (private)
//
// Adjusts the isolated margin of a delivery position; a positive change adds
// margin, a negative change removes it.
type UpdateDeliveryPositionMarginService struct {
	c        *DeliveryClient
	settle   Settle
	contract string
	change   decimal.Decimal
}

func (c *DeliveryClient) NewUpdateDeliveryPositionMarginService(settle Settle, contract string, change decimal.Decimal) *UpdateDeliveryPositionMarginService {
	return &UpdateDeliveryPositionMarginService{c: c, settle: settle, contract: contract, change: change}
}

func (s *UpdateDeliveryPositionMarginService) Do(ctx context.Context) (*DeliveryPosition, error) {
	req := request.Post(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/positions/"+s.contract+"/margin").
		WithSign().
		SetQuery("change", s.change.String())
	return request.Do[DeliveryPosition](req)
}

// UpdateDeliveryPositionLeverageService -- POST /api/v4/delivery/{settle}/positions/{contract}/leverage (private)
//
// Sets the leverage of a delivery position; 0 selects cross margin, a positive
// number selects isolated margin.
type UpdateDeliveryPositionLeverageService struct {
	c        *DeliveryClient
	settle   Settle
	contract string
	leverage decimal.Decimal
}

func (c *DeliveryClient) NewUpdateDeliveryPositionLeverageService(settle Settle, contract string, leverage decimal.Decimal) *UpdateDeliveryPositionLeverageService {
	return &UpdateDeliveryPositionLeverageService{c: c, settle: settle, contract: contract, leverage: leverage}
}

func (s *UpdateDeliveryPositionLeverageService) Do(ctx context.Context) (*DeliveryPosition, error) {
	req := request.Post(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/positions/"+s.contract+"/leverage").
		WithSign().
		SetQuery("leverage", s.leverage.String())
	return request.Do[DeliveryPosition](req)
}

// UpdateDeliveryPositionRiskLimitService -- POST /api/v4/delivery/{settle}/positions/{contract}/risk_limit (private)
//
// Sets the risk limit of a delivery position.
type UpdateDeliveryPositionRiskLimitService struct {
	c         *DeliveryClient
	settle    Settle
	contract  string
	riskLimit decimal.Decimal
}

func (c *DeliveryClient) NewUpdateDeliveryPositionRiskLimitService(settle Settle, contract string, riskLimit decimal.Decimal) *UpdateDeliveryPositionRiskLimitService {
	return &UpdateDeliveryPositionRiskLimitService{c: c, settle: settle, contract: contract, riskLimit: riskLimit}
}

func (s *UpdateDeliveryPositionRiskLimitService) Do(ctx context.Context) (*DeliveryPosition, error) {
	req := request.Post(ctx, s.c, "/api/v4/delivery/"+string(s.settle)+"/positions/"+s.contract+"/risk_limit").
		WithSign().
		SetQuery("risk_limit", s.riskLimit.String())
	return request.Do[DeliveryPosition](req)
}

// DeliveryPosition is a single delivery contract position held by the account.
type DeliveryPosition struct {
	User                   int64                      `json:"user"`
	Contract               string                     `json:"contract"`
	Size                   int64                      `json:"size"`
	Leverage               decimal.Decimal            `json:"leverage"`
	RiskLimit              decimal.Decimal            `json:"risk_limit"`
	LeverageMax            decimal.Decimal            `json:"leverage_max"`
	MaintenanceRate        decimal.Decimal            `json:"maintenance_rate"`
	Value                  decimal.Decimal            `json:"value"`
	Margin                 decimal.Decimal            `json:"margin"`
	EntryPrice             decimal.Decimal            `json:"entry_price"`
	LiqPrice               decimal.Decimal            `json:"liq_price"`
	MarkPrice              decimal.Decimal            `json:"mark_price"`
	InitialMargin          decimal.Decimal            `json:"initial_margin"`
	MaintenanceMargin      decimal.Decimal            `json:"maintenance_margin"`
	UnrealisedPnL          decimal.Decimal            `json:"unrealised_pnl"`
	RealisedPnL            decimal.Decimal            `json:"realised_pnl"`
	PnLPnL                 decimal.Decimal            `json:"pnl_pnl"`
	PnLFund                decimal.Decimal            `json:"pnl_fund"`
	PnLFee                 decimal.Decimal            `json:"pnl_fee"`
	HistoryPnL             decimal.Decimal            `json:"history_pnl"`
	LastClosePnL           decimal.Decimal            `json:"last_close_pnl"`
	RealisedPoint          decimal.Decimal            `json:"realised_point"`
	HistoryPoint           decimal.Decimal            `json:"history_point"`
	ADLRanking             int                        `json:"adl_ranking"`
	PendingOrders          int                        `json:"pending_orders"`
	CloseOrder             DeliveryPositionCloseOrder `json:"close_order"`
	Mode                   string                     `json:"mode"`
	CrossLeverageLimit     decimal.Decimal            `json:"cross_leverage_limit"`
	UpdateTime             time.Time                  `json:"update_time,format:unix"`
	UpdateID               int64                      `json:"update_id"`
	OpenTime               time.Time                  `json:"open_time,format:unix"`
	RiskLimitTable         string                     `json:"risk_limit_table"`
	AverageMaintenanceRate decimal.Decimal            `json:"average_maintenance_rate"`
	PID                    int64                      `json:"pid"`
}

// DeliveryPositionCloseOrder is the position's active close order, or the zero
// value when none is pending.
type DeliveryPositionCloseOrder struct {
	ID    int64           `json:"id"`
	Price decimal.Decimal `json:"price"`
	IsLiq bool            `json:"is_liq"`
}
