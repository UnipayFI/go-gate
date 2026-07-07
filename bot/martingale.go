package bot

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// CreateSpotMartingaleService -- POST /api/v4/bot/spot-martingale/create (private)
//
// Creates a spot Martingale strategy.
type CreateSpotMartingaleService struct {
	c            *BotClient
	strategyType string
	market       string
	createParams map[string]any
}

func (c *BotClient) NewCreateSpotMartingaleService(strategyType, market string, investAmount, priceDeviation decimal.Decimal, maxOrders int, takeProfitRatio decimal.Decimal) *CreateSpotMartingaleService {
	return &CreateSpotMartingaleService{c: c, strategyType: strategyType, market: market, createParams: map[string]any{
		"invest_amount":     investAmount.String(),
		"price_deviation":   priceDeviation.String(),
		"max_orders":        maxOrders,
		"take_profit_ratio": takeProfitRatio.String(),
	}}
}

// SetStopLossPerCycle sets the per-round stop-loss ratio as a decimal.
func (s *CreateSpotMartingaleService) SetStopLossPerCycle(stopLossPerCycle decimal.Decimal) *CreateSpotMartingaleService {
	s.createParams["stop_loss_per_cycle"] = stopLossPerCycle.String()
	return s
}

// SetTriggerPrice sets the price at which the strategy activates.
func (s *CreateSpotMartingaleService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateSpotMartingaleService {
	s.createParams["trigger_price"] = triggerPrice.String()
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio.
func (s *CreateSpotMartingaleService) SetProfitSharingRatio(profitSharingRatio decimal.Decimal) *CreateSpotMartingaleService {
	s.createParams["profit_sharing_ratio"] = profitSharingRatio.String()
	return s
}

func (s *CreateSpotMartingaleService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": s.strategyType,
		"market":        s.market,
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/spot-martingale/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}

// CreateContractMartingaleService -- POST /api/v4/bot/contract-martingale/create (private)
//
// Creates a contract Martingale strategy.
type CreateContractMartingaleService struct {
	c            *BotClient
	strategyType string
	market       string
	createParams map[string]any
}

func (c *BotClient) NewCreateContractMartingaleService(strategyType, market string, investAmount, priceDeviation decimal.Decimal, maxOrders int, takeProfitRatio decimal.Decimal, direction ContractMartingaleDirection, leverage decimal.Decimal) *CreateContractMartingaleService {
	return &CreateContractMartingaleService{c: c, strategyType: strategyType, market: market, createParams: map[string]any{
		"invest_amount":     investAmount.String(),
		"price_deviation":   priceDeviation.String(),
		"max_orders":        maxOrders,
		"take_profit_ratio": takeProfitRatio.String(),
		"direction":         string(direction),
		"leverage":          leverage.String(),
	}}
}

// SetStopLossPrice sets the legacy stop-loss price field.
func (s *CreateContractMartingaleService) SetStopLossPrice(stopLossPrice decimal.Decimal) *CreateContractMartingaleService {
	s.createParams["stop_loss_price"] = stopLossPrice.String()
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio.
func (s *CreateContractMartingaleService) SetProfitSharingRatio(profitSharingRatio decimal.Decimal) *CreateContractMartingaleService {
	s.createParams["profit_sharing_ratio"] = profitSharingRatio.String()
	return s
}

func (s *CreateContractMartingaleService) Do(ctx context.Context) (*AIHubCreateResponse, error) {
	body := map[string]any{
		"strategy_type": s.strategyType,
		"market":        s.market,
		"create_params": s.createParams,
	}
	req := request.Post(ctx, s.c, "/api/v4/bot/contract-martingale/create").WithSign().SetBody(body)
	return request.Do[AIHubCreateResponse](req)
}
