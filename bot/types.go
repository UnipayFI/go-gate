package bot

// StrategyType is the complete enumeration of policy types supported by AIHub
// (e.g. spot_grid, margin_grid, infinite_grid, futures_grid, spot_martingale,
// contract_martingale).
type StrategyType string

// DiscoverScene enumerates the scenarios supported by the strategy
// recommendation interface (e.g. filter, refresh).
type DiscoverScene string

// FuturesDirection is the direction enumeration supported by contract-based
// grid strategies.
type FuturesDirection string

// ContractMartingaleDirection is the direction enumeration supported by the
// contract Martingale strategy.
type ContractMartingaleDirection string

// AIHubCreateResponse is the envelope returned by every strategy-creation
// endpoint (spot-grid, margin-grid, infinite-grid, futures-grid,
// spot-martingale and contract-martingale).
type AIHubCreateResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    AIHubCreateData `json:"data"`
	TraceID string          `json:"trace_id"`
}

// AIHubCreateData is the policy information returned after a strategy is
// successfully created.
type AIHubCreateData struct {
	StrategyID   string       `json:"strategy_id"`
	StrategyType StrategyType `json:"strategy_type"`
	Market       string       `json:"market"`
	Status       string       `json:"status"`
	JumpURL      string       `json:"jump_url"`
}
