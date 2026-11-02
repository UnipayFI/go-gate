package bot

// StrategyType is the complete enumeration of policy types supported by AIHub
// (spot_grid, margin_grid, infinite_grid, futures_grid, spot_martingale,
// contract_martingale, smart_position, futures_combo).
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
// spot-martingale, contract-martingale, smart-position and futures-combo).
//
// HTTP 200 can carry a business failure: check Code first. On failure Data and
// TraceID may be absent, and a failed write may return only the facts Gate
// could verify (e.g. the strategy identity).
type AIHubCreateResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    AIHubCreateData `json:"data"`
	TraceID string          `json:"trace_id"`
	Label   string          `json:"label"`
}

// AIHubCreateData is the policy information returned after a strategy is
// successfully created. Status is the observed lifecycle state (or "unknown");
// an accepted or starting strategy is not yet running. OperationStatus is
// "pending" when the write result still needs verification and "rejected" for
// a proven refusal before any effect.
type AIHubCreateData struct {
	StrategyID       string       `json:"strategy_id"`
	StrategyType     StrategyType `json:"strategy_type"`
	Market           string       `json:"market"`
	Status           string       `json:"status"`
	JumpURL          string       `json:"jump_url"`
	OperationStatus  string       `json:"operation_status"`
	RecoveryRequired bool         `json:"recovery_required"`
}
