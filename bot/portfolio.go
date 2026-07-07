package bot

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListPortfolioRunningService -- GET /api/v4/bot/portfolio/running (private)
//
// Lists the authenticated account's currently running strategies.
type ListPortfolioRunningService struct {
	c      *BotClient
	params map[string]string
}

func (c *BotClient) NewListPortfolioRunningService() *ListPortfolioRunningService {
	return &ListPortfolioRunningService{c: c, params: map[string]string{}}
}

// SetStrategyType filters the list by policy type.
func (s *ListPortfolioRunningService) SetStrategyType(strategyType string) *ListPortfolioRunningService {
	s.params["strategy_type"] = strategyType
	return s
}

// SetMarket filters the list by trading pair.
func (s *ListPortfolioRunningService) SetMarket(market string) *ListPortfolioRunningService {
	s.params["market"] = market
	return s
}

// SetPage selects the result page (default 1).
func (s *ListPortfolioRunningService) SetPage(page int) *ListPortfolioRunningService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetPageSize sets the page size (default 20, maximum 50).
func (s *ListPortfolioRunningService) SetPageSize(pageSize int) *ListPortfolioRunningService {
	s.params["page_size"] = strconv.Itoa(pageSize)
	return s
}

func (s *ListPortfolioRunningService) Do(ctx context.Context) (*AIHubPortfolioRunningResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/bot/portfolio/running", s.params).WithSign()
	return request.Do[AIHubPortfolioRunningResponse](req)
}

// AIHubPortfolioRunningResponse is the envelope returned by the running-strategy
// list endpoint.
type AIHubPortfolioRunningResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    AIHubPortfolioRunningData `json:"data"`
}

// AIHubPortfolioRunningData is the paginated running-strategy list data.
type AIHubPortfolioRunningData struct {
	Items    []AIHubPortfolioRunningItem `json:"items"`
	Page     int                         `json:"page"`
	PageSize int                         `json:"page_size"`
	Total    int                         `json:"total"`
	TraceID  string                      `json:"trace_id"`
}

// AIHubPortfolioRunningItem is a single record in the list of running strategies.
type AIHubPortfolioRunningItem struct {
	StrategyID   string          `json:"strategy_id"`
	StrategyType StrategyType    `json:"strategy_type"`
	StrategyName string          `json:"strategy_name"`
	Market       string          `json:"market"`
	Status       string          `json:"status"`
	PnL          decimal.Decimal `json:"pnl"`
	PnLRate      decimal.Decimal `json:"pnl_rate"`
	InvestAmount decimal.Decimal `json:"invest_amount"`
	CreatedAt    string          `json:"created_at"`
}

// GetPortfolioDetailService -- GET /api/v4/bot/portfolio/detail (private)
//
// Returns the detail of a single strategy, keyed by strategy id and type.
type GetPortfolioDetailService struct {
	c      *BotClient
	params map[string]string
}

func (c *BotClient) NewGetPortfolioDetailService(strategyID, strategyType string) *GetPortfolioDetailService {
	return &GetPortfolioDetailService{c: c, params: map[string]string{
		"strategy_id":   strategyID,
		"strategy_type": strategyType,
	}}
}

func (s *GetPortfolioDetailService) Do(ctx context.Context) (*AIHubPortfolioDetailResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/bot/portfolio/detail", s.params).WithSign()
	return request.Do[AIHubPortfolioDetailResponse](req)
}

// AIHubPortfolioDetailResponse is the envelope returned by the strategy-detail
// endpoint.
type AIHubPortfolioDetailResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    AIHubPortfolioDetailData `json:"data"`
	TraceID string                   `json:"trace_id"`
}

// AIHubPortfolioDetailData is the strategy detail data. The metrics and position
// fields returned depend on the strategy type.
type AIHubPortfolioDetailData struct {
	StrategyID    string                  `json:"strategy_id"`
	StrategyType  StrategyType            `json:"strategy_type"`
	Market        string                  `json:"market"`
	Status        string                  `json:"status"`
	BaseInfo      AIHubPortfolioBaseInfo  `json:"base_info"`
	Metrics       AIHubPortfolioMetrics   `json:"metrics"`
	Position      *AIHubPortfolioPosition `json:"position"`
	StopSupported bool                    `json:"stop_supported"`
}

// AIHubPortfolioBaseInfo is the base info block of a strategy detail.
type AIHubPortfolioBaseInfo struct {
	StrategyName    string          `json:"strategy_name"`
	CreatedAt       string          `json:"created_at"`
	RunningDuration int64           `json:"running_duration"`
	InvestAmount    decimal.Decimal `json:"invest_amount"`
	TotalProfit     decimal.Decimal `json:"total_profit"`
	ProfitRate      decimal.Decimal `json:"profit_rate"`
}

// AIHubPortfolioMetrics is the metrics block of a strategy detail; which fields
// are populated depends on the strategy type.
type AIHubPortfolioMetrics struct {
	GridProfit                decimal.Decimal `json:"grid_profit"`
	FloatingPnL               decimal.Decimal `json:"floating_pnl"`
	ArbitrageCount            int64           `json:"arbitrage_count"`
	PriceRange                string          `json:"price_range"`
	GridCount                 int64           `json:"grid_count"`
	EstimatedLiquidationPrice decimal.Decimal `json:"estimated_liquidation_price"`
	PriceFloor                decimal.Decimal `json:"price_floor"`
	GridProfitRate            decimal.Decimal `json:"grid_profit_rate"`
	RealizedPnL               decimal.Decimal `json:"realized_pnl"`
	FinishedRounds            int64           `json:"finished_rounds"`
	AvgCost                   decimal.Decimal `json:"avg_cost"`
	TakeProfitPrice           decimal.Decimal `json:"take_profit_price"`
	MaintenanceMarginRatio    decimal.Decimal `json:"maintenance_margin_ratio"`
}

// AIHubPortfolioPosition is the position block of a strategy detail; which fields
// are populated depends on the strategy type.
type AIHubPortfolioPosition struct {
	Amount        decimal.Decimal `json:"amount"`
	EntryPrice    decimal.Decimal `json:"entry_price"`
	QuoteAmount   decimal.Decimal `json:"quote_amount"`
	PositionValue decimal.Decimal `json:"position_value"`
	Margin        decimal.Decimal `json:"margin"`
	Side          string          `json:"side"`
}

// StopPortfolioService -- POST /api/v4/bot/portfolio/stop (private)
//
// Terminates a single running strategy.
type StopPortfolioService struct {
	c    *BotClient
	body map[string]any
}

func (c *BotClient) NewStopPortfolioService(strategyID string, strategyType StrategyType) *StopPortfolioService {
	return &StopPortfolioService{c: c, body: map[string]any{
		"strategy_id":   strategyID,
		"strategy_type": string(strategyType),
	}}
}

func (s *StopPortfolioService) Do(ctx context.Context) (*AIHubPortfolioStopResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/bot/portfolio/stop", s.body).WithSign()
	return request.Do[AIHubPortfolioStopResponse](req)
}

// AIHubPortfolioStopResponse is the envelope returned by the strategy-stop
// endpoint.
type AIHubPortfolioStopResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    AIHubPortfolioStopData `json:"data"`
	TraceID string                 `json:"trace_id"`
}

// AIHubPortfolioStopData is the result information returned after a strategy is
// successfully terminated.
type AIHubPortfolioStopData struct {
	StrategyID    string       `json:"strategy_id"`
	StrategyType  StrategyType `json:"strategy_type"`
	Status        string       `json:"status"`
	ResultMessage string       `json:"result_message"`
}
