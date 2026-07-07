package bot

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetStrategyRecommendService -- GET /api/v4/bot/strategy/recommend (private)
//
// Returns AIHub strategy recommendations, optionally narrowed by market,
// strategy type, direction, target scenario and back-test filters.
type GetStrategyRecommendService struct {
	c      *BotClient
	params map[string]string
}

func (c *BotClient) NewGetStrategyRecommendService() *GetStrategyRecommendService {
	return &GetStrategyRecommendService{c: c, params: map[string]string{}}
}

// SetMarket narrows recommendations to a trading pair, such as BTC_USDT.
func (s *GetStrategyRecommendService) SetMarket(market string) *GetStrategyRecommendService {
	s.params["market"] = market
	return s
}

// SetStrategyType narrows recommendations to a target policy type
// (contract_martingale is not allowed here).
func (s *GetStrategyRecommendService) SetStrategyType(strategyType string) *GetStrategyRecommendService {
	s.params["strategy_type"] = strategyType
	return s
}

// SetDirection narrows recommendations to a market direction.
func (s *GetStrategyRecommendService) SetDirection(direction string) *GetStrategyRecommendService {
	s.params["direction"] = direction
	return s
}

// SetInvestAmount sets the investment amount used to size recommendations.
func (s *GetStrategyRecommendService) SetInvestAmount(investAmount decimal.Decimal) *GetStrategyRecommendService {
	s.params["invest_amount"] = investAmount.String()
	return s
}

// SetScene selects the recommendation scenario; when empty the service infers
// it automatically.
func (s *GetStrategyRecommendService) SetScene(scene string) *GetStrategyRecommendService {
	s.params["scene"] = scene
	return s
}

// SetRefreshRecommendationID passes the recommendation context to refresh, used
// when scene=refresh.
func (s *GetStrategyRecommendService) SetRefreshRecommendationID(refreshRecommendationID string) *GetStrategyRecommendService {
	s.params["refresh_recommendation_id"] = refreshRecommendationID
	return s
}

// SetLimit caps the number of recommendations returned (up to 10 when
// scene=filter).
func (s *GetStrategyRecommendService) SetLimit(limit int) *GetStrategyRecommendService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetMaxDrawdownLTE bounds recommendations to a maximum drawdown limit.
func (s *GetStrategyRecommendService) SetMaxDrawdownLTE(maxDrawdownLTE decimal.Decimal) *GetStrategyRecommendService {
	s.params["max_drawdown_lte"] = maxDrawdownLTE.String()
	return s
}

// SetBacktestAPRGTE bounds recommendations to a back-test annualized lower limit.
func (s *GetStrategyRecommendService) SetBacktestAPRGTE(backtestAPRGTE decimal.Decimal) *GetStrategyRecommendService {
	s.params["backtest_apr_gte"] = backtestAPRGTE.String()
	return s
}

func (s *GetStrategyRecommendService) Do(ctx context.Context) (*AIHubStrategyRecommendResponse, error) {
	req := request.Get(ctx, s.c, "/api/v4/bot/strategy/recommend", s.params).WithSign()
	return request.Do[AIHubStrategyRecommendResponse](req)
}

// AIHubStrategyRecommendResponse is the envelope returned by the strategy
// recommendation endpoint.
type AIHubStrategyRecommendResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    AIHubDiscoverData `json:"data"`
}

// AIHubDiscoverData is the strategy recommendation result data.
type AIHubDiscoverData struct {
	Scene              DiscoverScene         `json:"scene"`
	Recommendations    []AIHubRecommendation `json:"recommendations"`
	UnsupportedFilters []string              `json:"unsupported_filters"`
	TraceID            string                `json:"trace_id"`
}

// AIHubRecommendation is a single piece of strategy recommendation information.
type AIHubRecommendation struct {
	RecommendationID      string          `json:"recommendation_id"`
	BacktestID            int64           `json:"backtest_id"`
	Market                string          `json:"market"`
	StrategyType          StrategyType    `json:"strategy_type"`
	StrategyName          string          `json:"strategy_name"`
	BacktestAPR           decimal.Decimal `json:"backtest_apr"`
	MaxDrawdown           decimal.Decimal `json:"max_drawdown"`
	Summary               string          `json:"summary"`
	StrategyParamsPreview string          `json:"strategy_params_preview"`
}
