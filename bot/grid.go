package bot

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateSpotGridService -- POST /api/v4/bot/spot-grid/create (private)
//
// Creates a spot grid strategy.
type CreateSpotGridService struct {
	c            *BotClient
	strategyType string
	market       string
	createParams map[string]any
}

func (c *BotClient) NewCreateSpotGridService(strategyType, market string, money, lowPrice, highPrice decimal.Decimal, gridNum, priceType int) *CreateSpotGridService {
	return &CreateSpotGridService{c: c, strategyType: strategyType, market: market, createParams: map[string]any{
		"money":      money.String(),
		"low_price":  lowPrice.String(),
		"high_price": highPrice.String(),
		"grid_num":   gridNum,
		"price_type": priceType,
	}}
}

// SetTriggerPrice sets the price at which the strategy activates.
func (s *CreateSpotGridService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateSpotGridService {
	s.createParams["trigger_price"] = triggerPrice.String()
	return s
}

// SetStopProfit sets the take-profit price that ends the strategy.
func (s *CreateSpotGridService) SetStopProfit(stopProfit decimal.Decimal) *CreateSpotGridService {
	s.createParams["stop_profit"] = stopProfit.String()
	return s
}

// SetStopLoss sets the stop-loss price that ends the strategy.
func (s *CreateSpotGridService) SetStopLoss(stopLoss decimal.Decimal) *CreateSpotGridService {
	s.createParams["stop_loss"] = stopLoss.String()
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio.
func (s *CreateSpotGridService) SetProfitSharingRatio(profitSharingRatio decimal.Decimal) *CreateSpotGridService {
	s.createParams["profit_sharing_ratio"] = profitSharingRatio.String()
	return s
}

// SetIsUseBase invests using the base currency when true.
func (s *CreateSpotGridService) SetIsUseBase(isUseBase bool) *CreateSpotGridService {
	s.createParams["is_use_base"] = isUseBase
	return s
}

func (s *CreateSpotGridService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": s.strategyType,
		"market":        s.market,
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/spot-grid/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}

// CreateMarginGridService -- POST /api/v4/bot/margin-grid/create (private)
//
// Creates a leverage (margin) grid strategy.
type CreateMarginGridService struct {
	c            *BotClient
	strategyType string
	market       string
	createParams map[string]any
}

func (c *BotClient) NewCreateMarginGridService(strategyType, market string, money, lowPrice, highPrice decimal.Decimal, gridNum, priceType int, leverage decimal.Decimal) *CreateMarginGridService {
	return &CreateMarginGridService{c: c, strategyType: strategyType, market: market, createParams: map[string]any{
		"money":      money.String(),
		"low_price":  lowPrice.String(),
		"high_price": highPrice.String(),
		"grid_num":   gridNum,
		"price_type": priceType,
		"leverage":   leverage.String(),
	}}
}

// SetDirection sets the position direction of the leverage grid.
func (s *CreateMarginGridService) SetDirection(direction FuturesDirection) *CreateMarginGridService {
	s.createParams["direction"] = string(direction)
	return s
}

// SetTriggerPrice sets the price at which the strategy activates.
func (s *CreateMarginGridService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateMarginGridService {
	s.createParams["trigger_price"] = triggerPrice.String()
	return s
}

// SetStopProfit sets the take-profit price that ends the strategy.
func (s *CreateMarginGridService) SetStopProfit(stopProfit decimal.Decimal) *CreateMarginGridService {
	s.createParams["stop_profit"] = stopProfit.String()
	return s
}

// SetStopLoss sets the stop-loss price that ends the strategy.
func (s *CreateMarginGridService) SetStopLoss(stopLoss decimal.Decimal) *CreateMarginGridService {
	s.createParams["stop_loss"] = stopLoss.String()
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio.
func (s *CreateMarginGridService) SetProfitSharingRatio(profitSharingRatio decimal.Decimal) *CreateMarginGridService {
	s.createParams["profit_sharing_ratio"] = profitSharingRatio.String()
	return s
}

// SetIsUseBase invests using the base currency when true.
func (s *CreateMarginGridService) SetIsUseBase(isUseBase bool) *CreateMarginGridService {
	s.createParams["is_use_base"] = isUseBase
	return s
}

func (s *CreateMarginGridService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": s.strategyType,
		"market":        s.market,
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/margin-grid/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}

// CreateInfiniteGridService -- POST /api/v4/bot/infinite-grid/create (private)
//
// Creates an infinite grid strategy.
type CreateInfiniteGridService struct {
	c            *BotClient
	strategyType string
	market       string
	createParams map[string]any
}

func (c *BotClient) NewCreateInfiniteGridService(strategyType, market string, money, priceFloor, profitPerGrid decimal.Decimal) *CreateInfiniteGridService {
	return &CreateInfiniteGridService{c: c, strategyType: strategyType, market: market, createParams: map[string]any{
		"money":           money.String(),
		"price_floor":     priceFloor.String(),
		"profit_per_grid": profitPerGrid.String(),
	}}
}

// SetGridNum sets the number of grids (optional).
func (s *CreateInfiniteGridService) SetGridNum(gridNum int) *CreateInfiniteGridService {
	s.createParams["grid_num"] = gridNum
	return s
}

// SetPriceType sets the grid spacing type: 0 arithmetic, 1 geometric.
func (s *CreateInfiniteGridService) SetPriceType(priceType int) *CreateInfiniteGridService {
	s.createParams["price_type"] = priceType
	return s
}

// SetTriggerPrice sets the price at which the strategy activates.
func (s *CreateInfiniteGridService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateInfiniteGridService {
	s.createParams["trigger_price"] = triggerPrice.String()
	return s
}

// SetStopProfit sets the take-profit price that ends the strategy.
func (s *CreateInfiniteGridService) SetStopProfit(stopProfit decimal.Decimal) *CreateInfiniteGridService {
	s.createParams["stop_profit"] = stopProfit.String()
	return s
}

// SetStopLoss sets the stop-loss price that ends the strategy.
func (s *CreateInfiniteGridService) SetStopLoss(stopLoss decimal.Decimal) *CreateInfiniteGridService {
	s.createParams["stop_loss"] = stopLoss.String()
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio.
func (s *CreateInfiniteGridService) SetProfitSharingRatio(profitSharingRatio decimal.Decimal) *CreateInfiniteGridService {
	s.createParams["profit_sharing_ratio"] = profitSharingRatio.String()
	return s
}

// SetIsUseBase invests using the base currency when true.
func (s *CreateInfiniteGridService) SetIsUseBase(isUseBase bool) *CreateInfiniteGridService {
	s.createParams["is_use_base"] = isUseBase
	return s
}

func (s *CreateInfiniteGridService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": s.strategyType,
		"market":        s.market,
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/infinite-grid/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}

// CreateFuturesGridService -- POST /api/v4/bot/futures-grid/create (private)
//
// Creates a contract (futures) grid strategy.
type CreateFuturesGridService struct {
	c            *BotClient
	strategyType string
	market       string
	createParams map[string]any
}

func (c *BotClient) NewCreateFuturesGridService(strategyType, market string, money, lowPrice, highPrice decimal.Decimal, gridNum, priceType int, leverage decimal.Decimal) *CreateFuturesGridService {
	return &CreateFuturesGridService{c: c, strategyType: strategyType, market: market, createParams: map[string]any{
		"money":      money.String(),
		"low_price":  lowPrice.String(),
		"high_price": highPrice.String(),
		"grid_num":   gridNum,
		"price_type": priceType,
		"leverage":   leverage.String(),
	}}
}

// SetDirection sets the position direction of the contract grid.
func (s *CreateFuturesGridService) SetDirection(direction FuturesDirection) *CreateFuturesGridService {
	s.createParams["direction"] = string(direction)
	return s
}

// SetTriggerPrice sets the price at which the strategy activates.
func (s *CreateFuturesGridService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateFuturesGridService {
	s.createParams["trigger_price"] = triggerPrice.String()
	return s
}

// SetStopProfit sets the take-profit price that ends the strategy.
func (s *CreateFuturesGridService) SetStopProfit(stopProfit decimal.Decimal) *CreateFuturesGridService {
	s.createParams["stop_profit"] = stopProfit.String()
	return s
}

// SetStopLoss sets the stop-loss price that ends the strategy.
func (s *CreateFuturesGridService) SetStopLoss(stopLoss decimal.Decimal) *CreateFuturesGridService {
	s.createParams["stop_loss"] = stopLoss.String()
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio.
func (s *CreateFuturesGridService) SetProfitSharingRatio(profitSharingRatio decimal.Decimal) *CreateFuturesGridService {
	s.createParams["profit_sharing_ratio"] = profitSharingRatio.String()
	return s
}

// SetIsUseBase invests using the base currency when true.
func (s *CreateFuturesGridService) SetIsUseBase(isUseBase bool) *CreateFuturesGridService {
	s.createParams["is_use_base"] = isUseBase
	return s
}

func (s *CreateFuturesGridService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": s.strategyType,
		"market":        s.market,
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/futures-grid/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}
