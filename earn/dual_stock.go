package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// The stock Dual Investment endpoints mirror the crypto ones under
// /earn/dual/*, sharing their response schemas, so they decode into the same
// types. Gate documents the product lists as public, but live they reject
// unsigned requests, so every service here signs.

// ListDualStockInvestmentPlansService -- GET /api/v4/earn/dual/stock/investment_plan (private)
//
// Returns the available stock Dual Investment products.
type ListDualStockInvestmentPlansService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListDualStockInvestmentPlansService() *ListDualStockInvestmentPlansService {
	return &ListDualStockInvestmentPlansService{c: c, params: map[string]string{}}
}

// SetPlanID narrows the result to a single product.
func (s *ListDualStockInvestmentPlansService) SetPlanID(planID int64) *ListDualStockInvestmentPlansService {
	s.params["plan_id"] = strconv.FormatInt(planID, 10)
	return s
}

// SetCoin filters by underlying currency.
func (s *ListDualStockInvestmentPlansService) SetCoin(coin string) *ListDualStockInvestmentPlansService {
	s.params["coin"] = coin
	return s
}

// SetType filters by direction: "call" (sell high) or "put" (buy low).
func (s *ListDualStockInvestmentPlansService) SetType(planType string) *ListDualStockInvestmentPlansService {
	s.params["type"] = planType
	return s
}

// SetQuoteCurrency sets the settlement currency (default USDT).
func (s *ListDualStockInvestmentPlansService) SetQuoteCurrency(quoteCurrency string) *ListDualStockInvestmentPlansService {
	s.params["quote_currency"] = quoteCurrency
	return s
}

// SetSort selects the sort order: "apy" (highest APY first), "short-period"
// (shortest tenor first) or "multiple" (highest premium first).
func (s *ListDualStockInvestmentPlansService) SetSort(sort string) *ListDualStockInvestmentPlansService {
	s.params["sort"] = sort
	return s
}

// SetPage selects the result page (1-based).
func (s *ListDualStockInvestmentPlansService) SetPage(page int) *ListDualStockInvestmentPlansService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize sets the number of items per page (max 100).
func (s *ListDualStockInvestmentPlansService) SetPageSize(pageSize int) *ListDualStockInvestmentPlansService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListDualStockInvestmentPlansService) Do(ctx context.Context) ([]DualInvestmentPlan, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/stock/investment_plan", s.params).WithSign()
	resp, err := request.Do[[]DualInvestmentPlan](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListDualStockProjectRecommendService -- GET /api/v4/earn/dual/stock/project-recommend (private)
//
// Returns recommended stock Dual Investment products.
type ListDualStockProjectRecommendService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListDualStockProjectRecommendService() *ListDualStockProjectRecommendService {
	return &ListDualStockProjectRecommendService{c: c, params: map[string]string{}}
}

// SetCoin filters by underlying currency.
func (s *ListDualStockProjectRecommendService) SetCoin(coin string) *ListDualStockProjectRecommendService {
	s.params["coin"] = coin
	return s
}

// SetType filters by direction: "call" (sell high) or "put" (buy low).
func (s *ListDualStockProjectRecommendService) SetType(projectType string) *ListDualStockProjectRecommendService {
	s.params["type"] = projectType
	return s
}

// SetMode selects the recommendation mode: "normal" (default, regular users) or
// "senior" (experienced users).
func (s *ListDualStockProjectRecommendService) SetMode(mode string) *ListDualStockProjectRecommendService {
	s.params["mode"] = mode
	return s
}

// SetHistoryPids excludes already-recommended projects (comma-separated IDs).
func (s *ListDualStockProjectRecommendService) SetHistoryPids(historyPids string) *ListDualStockProjectRecommendService {
	s.params["history_pids"] = historyPids
	return s
}

func (s *ListDualStockProjectRecommendService) Do(ctx context.Context) ([]DualProjectRecommend, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/stock/project-recommend", s.params).WithSign()
	resp, err := request.Do[[]DualProjectRecommend](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ListDualStockOrdersService -- GET /api/v4/earn/dual/stock/orders (private)
//
// Returns the authenticated user's stock Dual Investment orders.
type ListDualStockOrdersService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewListDualStockOrdersService() *ListDualStockOrdersService {
	return &ListDualStockOrdersService{c: c, params: map[string]string{}}
}

// SetCoin filters by investment token.
func (s *ListDualStockOrdersService) SetCoin(coin string) *ListDualStockOrdersService {
	s.params["coin"] = coin
	return s
}

// SetStatus filters by order status: "HOLD" (current holdings), "REPAY"
// (historical holdings), "PROCESSING", "SETTLEMENT_PROCESSING" or "ALL"
// (default).
func (s *ListDualStockOrdersService) SetStatus(status string) *ListDualStockOrdersService {
	s.params["status"] = status
	return s
}

// SetType filters by direction: "call" (sell high) or "put" (buy low).
func (s *ListDualStockOrdersService) SetType(orderType string) *ListDualStockOrdersService {
	s.params["type"] = orderType
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListDualStockOrdersService) SetFrom(from time.Time) *ListDualStockOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range.
func (s *ListDualStockOrdersService) SetTo(to time.Time) *ListDualStockOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetPage selects the result page (1-based).
func (s *ListDualStockOrdersService) SetPage(page int) *ListDualStockOrdersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize sets the number of records per page (max 100, default 100).
func (s *ListDualStockOrdersService) SetPageSize(pageSize int) *ListDualStockOrdersService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *ListDualStockOrdersService) SetLimit(limit int) *ListDualStockOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListDualStockOrdersService) Do(ctx context.Context) ([]DualOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/stock/orders", s.params).WithSign()
	resp, err := request.Do[[]DualOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PlaceDualStockOrderService -- POST /api/v4/earn/dual/stock/orders (private)
//
// Subscribes to a stock Dual Investment product. For buy-low (put) products the
// hit_* flags collect funds from the spot, stock or Simple Earn account before
// subscribing.
type PlaceDualStockOrderService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewPlaceDualStockOrderService(planID string, amount decimal.Decimal) *PlaceDualStockOrderService {
	return &PlaceDualStockOrderService{c: c, body: map[string]any{
		"plan_id": planID,
		"amount":  amount.String(),
	}}
}

// SetText attaches a custom order id. It must start with "t-", be at most 28
// bytes after the prefix, and contain only letters, digits, "_", "-" or ".".
func (s *PlaceDualStockOrderService) SetText(text string) *PlaceDualStockOrderService {
	s.body["text"] = text
	return s
}

// SetHitSpot funds a put subscription from the spot account when 1.
func (s *PlaceDualStockOrderService) SetHitSpot(hitSpot int) *PlaceDualStockOrderService {
	s.body["hit_spot"] = hitSpot
	return s
}

// SetHitStock funds a put subscription from the stock account when 1.
func (s *PlaceDualStockOrderService) SetHitStock(hitStock int) *PlaceDualStockOrderService {
	s.body["hit_stock"] = hitStock
	return s
}

// SetHitSimple funds a put subscription from the Simple Earn account when 1.
func (s *PlaceDualStockOrderService) SetHitSimple(hitSimple int) *PlaceDualStockOrderService {
	s.body["hit_simple"] = hitSimple
	return s
}

// Do places the order. Gate documents the created order as the response; an
// empty body decodes to a zero DualOrder.
func (s *PlaceDualStockOrderService) Do(ctx context.Context) (*DualOrder, error) {
	req := request.Post(ctx, s.c, "/api/v4/earn/dual/stock/orders", s.body).WithSign()
	return request.Do[DualOrder](req)
}

// GetDualStockBalanceService -- GET /api/v4/earn/dual/stock/balance (private)
//
// Returns the authenticated user's aggregated stock Dual Investment assets and
// interest, valued in both USDT and BTC.
type GetDualStockBalanceService struct {
	c *EarnClient
}

func (c *EarnClient) NewGetDualStockBalanceService() *GetDualStockBalanceService {
	return &GetDualStockBalanceService{c: c}
}

func (s *GetDualStockBalanceService) Do(ctx context.Context) (*DualBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/stock/balance").WithSign()
	return request.Do[DualBalance](req)
}

// GetDualStockOrderRefundPreviewService -- GET /api/v4/earn/dual/stock/order-refund-preview (private)
//
// Previews the settlement terms of an early stock Dual Investment redemption;
// the returned req_id is required to submit the actual redemption.
type GetDualStockOrderRefundPreviewService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetDualStockOrderRefundPreviewService(orderID string) *GetDualStockOrderRefundPreviewService {
	return &GetDualStockOrderRefundPreviewService{c: c, params: map[string]string{
		"order_id": orderID,
	}}
}

func (s *GetDualStockOrderRefundPreviewService) Do(ctx context.Context) (*DualOrderRefundPreview, error) {
	req := request.Get(ctx, s.c, "/api/v4/earn/dual/stock/order-refund-preview", s.params).WithSign()
	return request.Do[DualOrderRefundPreview](req)
}

// RefundDualStockOrderService -- POST /api/v4/earn/dual/stock/order-refund (private)
//
// Submits an early redemption of a stock Dual Investment order using the req_id
// returned by the refund preview. Gate returns an empty body on success.
type RefundDualStockOrderService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewRefundDualStockOrderService(orderID, reqID string) *RefundDualStockOrderService {
	return &RefundDualStockOrderService{c: c, body: map[string]any{
		"order_id": orderID,
		"req_id":   reqID,
	}}
}

func (s *RefundDualStockOrderService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/dual/stock/order-refund", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ModifyDualStockOrderReinvestService -- POST /api/v4/earn/dual/stock/modify-order-reinvest (private)
//
// Toggles auto-reinvestment for a stock Dual Investment order. All parameters
// are optional; Gate returns an empty body on success.
type ModifyDualStockOrderReinvestService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewModifyDualStockOrderReinvestService() *ModifyDualStockOrderReinvestService {
	return &ModifyDualStockOrderReinvestService{c: c, body: map[string]any{}}
}

// SetOrderID selects the order to modify.
func (s *ModifyDualStockOrderReinvestService) SetOrderID(orderID int64) *ModifyDualStockOrderReinvestService {
	s.body["order_id"] = orderID
	return s
}

// SetStatus toggles reinvestment: 0 — off; 1 — on.
func (s *ModifyDualStockOrderReinvestService) SetStatus(status int) *ModifyDualStockOrderReinvestService {
	s.body["status"] = status
	return s
}

// SetEffectiveTimeDuration sets the effective duration in seconds (default 1 day,
// 86400).
func (s *ModifyDualStockOrderReinvestService) SetEffectiveTimeDuration(seconds int64) *ModifyDualStockOrderReinvestService {
	s.body["effective_time_duration"] = seconds
	return s
}

func (s *ModifyDualStockOrderReinvestService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/earn/dual/stock/modify-order-reinvest", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}
