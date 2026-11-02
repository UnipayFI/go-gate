package bot

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// AIHubRebalanceAsset is one asset and its target weight in a smart-position
// strategy. Rate has two decimals (minimum 0.01); assets are unique and the
// rates sum to exactly 1.
type AIHubRebalanceAsset struct {
	Asset string `json:"asset"`
	Rate  string `json:"rate"`
}

// AIHubFuturesComboAsset is one contract member of a futures-combo strategy.
// Side is "buy" or "sell"; Rate is 0.01-1 in steps of 0.01 and the rates sum to
// 1.
type AIHubFuturesComboAsset struct {
	Market string `json:"market"`
	Side   string `json:"side"`
	Rate   string `json:"rate"`
}

// CreateSmartPositionService -- POST /api/v4/bot/smart-position/create (private)
//
// Creates a smart-position (rebalance) strategy over 2-10 assets. Exactly one of
// SetPeriod or SetThreshold must be set. beginMoneyCurrency is the investment
// currency (USDT).
type CreateSmartPositionService struct {
	c            *BotClient
	createParams map[string]any
}

func (c *BotClient) NewCreateSmartPositionService(beginMoneyAmount decimal.Decimal, beginMoneyCurrency string, assetInfo []AIHubRebalanceAsset) *CreateSmartPositionService {
	return &CreateSmartPositionService{c: c, createParams: map[string]any{
		"begin_money_amount":   beginMoneyAmount.String(),
		"begin_money_currency": beginMoneyCurrency,
		"asset_info":           assetInfo,
	}}
}

// SetPeriod rebalances on a fixed interval between 1m and 28d (exclusive with
// threshold).
func (s *CreateSmartPositionService) SetPeriod(period string) *CreateSmartPositionService {
	s.createParams["period"] = period
	return s
}

// SetThreshold rebalances when a weight drifts by this positive ratio below 1,
// with at most four decimals (exclusive with period).
func (s *CreateSmartPositionService) SetThreshold(threshold decimal.Decimal) *CreateSmartPositionService {
	s.createParams["threshold"] = threshold.String()
	return s
}

// SetCloseFlag sets the close flag (0 or 1).
func (s *CreateSmartPositionService) SetCloseFlag(closeFlag int) *CreateSmartPositionService {
	s.createParams["close_flag"] = closeFlag
	return s
}

// SetIsUseBase invests using held base assets when true.
func (s *CreateSmartPositionService) SetIsUseBase(isUseBase bool) *CreateSmartPositionService {
	s.createParams["is_use_base"] = isUseBase
	return s
}

// SetStopProfit sets the take-profit that ends the strategy.
func (s *CreateSmartPositionService) SetStopProfit(stopProfit decimal.Decimal) *CreateSmartPositionService {
	s.createParams["stop_profit"] = stopProfit.String()
	return s
}

// SetStopLoss sets the stop-loss that ends the strategy.
func (s *CreateSmartPositionService) SetStopLoss(stopLoss decimal.Decimal) *CreateSmartPositionService {
	s.createParams["stop_loss"] = stopLoss.String()
	return s
}

func (s *CreateSmartPositionService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": "smart_position",
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/smart-position/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}

// CreateFuturesComboService -- POST /api/v4/bot/futures-combo/create (private)
//
// Creates a futures-combo strategy over 2-10 contracts. Exactly one of SetPeriod
// or SetThreshold must be set. beginMoneyCurrency is the investment currency
// (USDT).
type CreateFuturesComboService struct {
	c            *BotClient
	createParams map[string]any
}

func (c *BotClient) NewCreateFuturesComboService(beginMoneyAmount decimal.Decimal, beginMoneyCurrency string, assetInfo []AIHubFuturesComboAsset, leverage int) *CreateFuturesComboService {
	return &CreateFuturesComboService{c: c, createParams: map[string]any{
		"begin_money_amount":   beginMoneyAmount.String(),
		"begin_money_currency": beginMoneyCurrency,
		"asset_info":           assetInfo,
		"leverage":             leverage,
	}}
}

// SetPeriod rebalances on a fixed interval: 30m, 1h, 4h, 8h, 12h, 1d, 3d, 7d,
// 14d or 28d (exclusive with threshold).
func (s *CreateFuturesComboService) SetPeriod(period string) *CreateFuturesComboService {
	s.createParams["period"] = period
	return s
}

// SetThreshold rebalances when a weight drifts by this ratio, 0.01-0.50
// (exclusive with period).
func (s *CreateFuturesComboService) SetThreshold(threshold decimal.Decimal) *CreateFuturesComboService {
	s.createParams["threshold"] = threshold.String()
	return s
}

// SetStopProfit sets the take-profit that ends the strategy.
func (s *CreateFuturesComboService) SetStopProfit(stopProfit decimal.Decimal) *CreateFuturesComboService {
	s.createParams["stop_profit"] = stopProfit.String()
	return s
}

// SetStopLoss sets the stop-loss that ends the strategy.
func (s *CreateFuturesComboService) SetStopLoss(stopLoss decimal.Decimal) *CreateFuturesComboService {
	s.createParams["stop_loss"] = stopLoss.String()
	return s
}

func (s *CreateFuturesComboService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": "futures_combo",
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/futures-combo/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}
