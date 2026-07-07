package earn

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetDualBalanceService -- GET /api/v4/earn/dual/balance (private)
//
// Returns the authenticated user's aggregated Dual-Currency Earning assets and
// interest, valued in both USDT and BTC.
type GetDualBalanceService struct {
	c *EarnClient
}

func (c *EarnClient) NewGetDualBalanceService() *GetDualBalanceService {
	return &GetDualBalanceService{c: c}
}

func (s *GetDualBalanceService) Do(ctx context.Context) (*DualBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/balance").WithSign()
	return request.Do[DualBalance](req)
}

// DualBalance is the user's Dual-Currency Earning assets and interest, valued in
// USDT and BTC equivalents.
type DualBalance struct {
	UserAssetUSDT         decimal.Decimal `json:"user_asset_usdt"`
	UserAssetBTC          decimal.Decimal `json:"user_asset_btc"`
	UserTotalInterestUSDT decimal.Decimal `json:"user_total_interest_usdt"`
	UserTotalInterestBTC  decimal.Decimal `json:"user_total_interest_btc"`
}

// GetDualOrderRefundPreviewService -- GET /api/v4/earn/dual/order-refund-preview (private)
//
// Previews the settlement terms of an early Dual-Currency order redemption; the
// returned req_id is required to submit the actual redemption.
type GetDualOrderRefundPreviewService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetDualOrderRefundPreviewService(orderID string) *GetDualOrderRefundPreviewService {
	return &GetDualOrderRefundPreviewService{c: c, params: map[string]string{
		"order_id": orderID,
	}}
}

func (s *GetDualOrderRefundPreviewService) Do(ctx context.Context) (*DualOrderRefundPreview, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/order-refund-preview", s.params).WithSign()
	return request.Do[DualOrderRefundPreview](req)
}

// DualOrderRefundPreview is the previewed settlement of an early Dual-Currency
// redemption.
type DualOrderRefundPreview struct {
	CreateTimest        time.Time       `json:"create_timest,format:unix"`
	DeliveryTimest      time.Time       `json:"delivery_timest,format:unix"`
	ExercisePrice       decimal.Decimal `json:"exercise_price"`
	InvestAmount        decimal.Decimal `json:"invest_amount"`
	InvestCurrency      string          `json:"invest_currency"`
	Name                string          `json:"name"`
	OrderID             int64           `json:"order_id"`
	ReqID               string          `json:"req_id"`
	RefundServiceCharge int64           `json:"refund_service_charge"`
	SettlePrice         decimal.Decimal `json:"settle_price"`
	SettlementAmount    decimal.Decimal `json:"settlement_amount"`
	SettlementCurrency  string          `json:"settlement_currency"`
	SettlementInterest  decimal.Decimal `json:"settlement_interest"`
	SettlementPrinciple decimal.Decimal `json:"settlement_principle"`
	Type                string          `json:"type"`
	MoneyBackTimest     time.Time       `json:"money_back_timest,format:unix"`
}

// RefundDualOrderService -- POST /api/v4/earn/dual/order-refund (private)
//
// Submits an early redemption of a Dual-Currency order using the req_id returned
// by the refund preview. Gate returns an empty body on success.
type RefundDualOrderService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewRefundDualOrderService(orderID, reqID string) *RefundDualOrderService {
	return &RefundDualOrderService{c: c, body: map[string]any{
		"order_id": orderID,
		"req_id":   reqID,
	}}
}

func (s *RefundDualOrderService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/dual/order-refund", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ModifyDualOrderReinvestService -- POST /api/v4/earn/dual/modify-order-reinvest (private)
//
// Toggles auto-reinvestment for a Dual-Currency order. All parameters are
// optional; Gate returns an empty body on success.
type ModifyDualOrderReinvestService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewModifyDualOrderReinvestService() *ModifyDualOrderReinvestService {
	return &ModifyDualOrderReinvestService{c: c, body: map[string]any{}}
}

// SetOrderID selects the order to modify.
func (s *ModifyDualOrderReinvestService) SetOrderID(orderID int64) *ModifyDualOrderReinvestService {
	s.body["order_id"] = orderID
	return s
}

// SetStatus toggles reinvestment: 0 — off; 1 — on.
func (s *ModifyDualOrderReinvestService) SetStatus(status int) *ModifyDualOrderReinvestService {
	s.body["status"] = status
	return s
}

// SetEffectiveTimeDuration sets the effective duration in seconds (default 1 day,
// 86400).
func (s *ModifyDualOrderReinvestService) SetEffectiveTimeDuration(seconds int64) *ModifyDualOrderReinvestService {
	s.body["effective_time_duration"] = seconds
	return s
}

func (s *ModifyDualOrderReinvestService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/dual/modify-order-reinvest", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListDualProjectRecommendService -- GET /api/v4/earn/dual/project-recommend (private)
//
// Returns recommended Dual-Currency projects, optionally filtered by sort mode,
// investment token or direction.
type ListDualProjectRecommendService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListDualProjectRecommendService() *ListDualProjectRecommendService {
	return &ListDualProjectRecommendService{c: c, params: map[string]string{}}
}

// SetMode selects the sort mode (default "normal"; e.g. "senior", "apy_up",
// "ep_d").
func (s *ListDualProjectRecommendService) SetMode(mode string) *ListDualProjectRecommendService {
	s.params["mode"] = mode
	return s
}

// SetCoin filters by investment token.
func (s *ListDualProjectRecommendService) SetCoin(coin string) *ListDualProjectRecommendService {
	s.params["coin"] = coin
	return s
}

// SetType filters by direction: "call" (sell high) or "put" (buy low).
func (s *ListDualProjectRecommendService) SetType(projectType string) *ListDualProjectRecommendService {
	s.params["type"] = projectType
	return s
}

// SetHistoryPids excludes already-recommended projects (comma-separated IDs).
func (s *ListDualProjectRecommendService) SetHistoryPids(historyPids string) *ListDualProjectRecommendService {
	s.params["history_pids"] = historyPids
	return s
}

func (s *ListDualProjectRecommendService) Do(ctx context.Context) ([]DualProjectRecommend, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/project-recommend", s.params).WithSign()
	resp, err := request.Do[[]DualProjectRecommend](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DualProjectRecommend is a single recommended Dual-Currency project.
type DualProjectRecommend struct {
	ID               int64           `json:"id"`
	Category         int             `json:"category"`
	Type             string          `json:"type"`
	InvestCurrency   string          `json:"invest_currency"`
	ExerciseCurrency string          `json:"exercise_currency"`
	APYDisplay       decimal.Decimal `json:"apy_display"`
	ExercisePrice    decimal.Decimal `json:"exercise_price"`
	DeliveryTimest   time.Time       `json:"delivery_timest,format:unix"`
	MinAmount        decimal.Decimal `json:"min_amount"`
	MaxAmount        decimal.Decimal `json:"max_amount"`
	MinCopies        int64           `json:"min_copies"`
	MaxCopies        int64           `json:"max_copies"`
	InvestDays       int64           `json:"invest_days"`
	InvestHours      string          `json:"invest_hours"`
}
