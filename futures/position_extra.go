package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListPositionsWithTimeRangeService -- GET /api/v4/futures/{settle}/positions_timerange (private)
//
// Lists the account's historical (closed) positions on a contract within a
// Unix-second time range.
type ListPositionsWithTimeRangeService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewListPositionsWithTimeRangeService(settle Settle, contract string) *ListPositionsWithTimeRangeService {
	return &ListPositionsWithTimeRangeService{c: c, settle: settle, params: map[string]string{"contract": contract}}
}

// SetFrom sets the start time (inclusive).
func (s *ListPositionsWithTimeRangeService) SetFrom(from time.Time) *ListPositionsWithTimeRangeService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time (inclusive).
func (s *ListPositionsWithTimeRangeService) SetTo(to time.Time) *ListPositionsWithTimeRangeService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned.
func (s *ListPositionsWithTimeRangeService) SetLimit(limit int) *ListPositionsWithTimeRangeService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListPositionsWithTimeRangeService) SetOffset(offset int) *ListPositionsWithTimeRangeService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListPositionsWithTimeRangeService) Do(ctx context.Context) ([]FuturesPositionHistory, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions_timerange", s.params).WithSign()
	resp, err := request.Do[[]FuturesPositionHistory](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FuturesPositionHistory is one historical (closed) position record. The newer
// history endpoint encodes size and all money fields as strings; time is a
// Unix-second timestamp.
type FuturesPositionHistory struct {
	Contract           string          `json:"contract"`
	Size               decimal.Decimal `json:"size"`
	Leverage           decimal.Decimal `json:"leverage"`
	RiskLimit          decimal.Decimal `json:"risk_limit"`
	LeverageMax        decimal.Decimal `json:"leverage_max"`
	MaintenanceRate    decimal.Decimal `json:"maintenance_rate"`
	Margin             decimal.Decimal `json:"margin"`
	LiqPrice           decimal.Decimal `json:"liq_price"`
	RealisedPnL        decimal.Decimal `json:"realised_pnl"`
	HistoryPnL         decimal.Decimal `json:"history_pnl"`
	LastClosePnL       decimal.Decimal `json:"last_close_pnl"`
	RealisedPoint      decimal.Decimal `json:"realised_point"`
	HistoryPoint       decimal.Decimal `json:"history_point"`
	Mode               string          `json:"mode"`
	CrossLeverageLimit decimal.Decimal `json:"cross_leverage_limit"`
	EntryPrice         decimal.Decimal `json:"entry_price"`
	Time               time.Time       `json:"time,format:unix"`
}

// GetPositionLeverageService -- GET /api/v4/futures/{settle}/get_leverage/{contract} (private)
//
// Returns the configured leverage of a contract in a specific position-margin
// mode. posMarginMode is required for the split-position mode; dualSide selects
// the leg ("dual_long" or "dual_short").
type GetPositionLeverageService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewGetPositionLeverageService(settle Settle, contract, posMarginMode, dualSide string) *GetPositionLeverageService {
	return &GetPositionLeverageService{
		c:        c,
		settle:   settle,
		contract: contract,
		params:   map[string]string{"pos_margin_mode": posMarginMode, "dual_side": dualSide},
	}
}

func (s *GetPositionLeverageService) Do(ctx context.Context) (*FuturesLeverageInfo, error) {
	req := request.Get(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/get_leverage/"+s.contract, s.params).WithSign()
	return request.Do[FuturesLeverageInfo](req)
}

// FuturesLeverageInfo is the leverage configuration returned for a contract.
type FuturesLeverageInfo struct {
	Leverage decimal.Decimal `json:"leverage"`
}

// SetPositionLeverageService -- POST /api/v4/futures/{settle}/positions/{contract}/set_leverage (private)
//
// Updates the leverage of a contract for a specific margin mode. marginMode is
// "isolated" or "cross"; dualSide (optional) selects the leg in hedge mode.
type SetPositionLeverageService struct {
	c        *FuturesClient
	settle   Settle
	contract string
	params   map[string]string
}

func (c *FuturesClient) NewSetPositionLeverageService(settle Settle, contract, leverage, marginMode string) *SetPositionLeverageService {
	return &SetPositionLeverageService{
		c:        c,
		settle:   settle,
		contract: contract,
		params:   map[string]string{"leverage": leverage, "margin_mode": marginMode},
	}
}

// SetDualSide selects the leg to update in hedge mode ("dual_long" or "dual_short").
func (s *SetPositionLeverageService) SetDualSide(dualSide string) *SetPositionLeverageService {
	s.params["dual_side"] = dualSide
	return s
}

func (s *SetPositionLeverageService) Do(ctx context.Context) (*FuturesSplitPosition, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/positions/"+s.contract+"/set_leverage")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	return request.Do[FuturesSplitPosition](req)
}

// FuturesSplitPosition is the position snapshot returned by the split-position
// (set_leverage) endpoint. It carries hedge-mode fields and encodes size / money
// values as strings; update_time / open_time are Unix-second timestamps.
type FuturesSplitPosition struct {
	User                   int64               `json:"user"`
	Contract               string              `json:"contract"`
	Size                   decimal.Decimal     `json:"size"`
	HedgeStatus            string              `json:"hedge_status"`
	HedgedSize             decimal.Decimal     `json:"hedged_size"`
	UnhedgedSize           decimal.Decimal     `json:"unhedged_size"`
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
	PosMarginMode          string              `json:"pos_margin_mode"`
	Lever                  decimal.Decimal     `json:"lever"`
}

// SetPositionModeService -- POST /api/v4/futures/{settle}/set_position_mode (private)
//
// Sets the account's position-holding mode for a settlement currency. positionMode
// is one of "single", "dual", "dual_plus" or "repeat". This replaces the older
// dual_mode toggle and returns the resulting futures account snapshot.
type SetPositionModeService struct {
	c      *FuturesClient
	settle Settle
	params map[string]string
}

func (c *FuturesClient) NewSetPositionModeService(settle Settle, positionMode string) *SetPositionModeService {
	return &SetPositionModeService{c: c, settle: settle, params: map[string]string{"position_mode": positionMode}}
}

func (s *SetPositionModeService) Do(ctx context.Context) (*FuturesAccount, error) {
	req := request.Post(ctx, s.c, "/api/v4/futures/"+string(s.settle)+"/set_position_mode")
	for k, v := range s.params {
		req.SetQuery(k, v)
	}
	req.WithSign()
	return request.Do[FuturesAccount](req)
}
