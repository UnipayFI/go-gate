package unified

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetUnifiedDeltaNeutralService -- GET /api/v4/unified/delta_neutral (private)
//
// Queries whether the account's Delta-neutral strategy mode is enabled.
type GetUnifiedDeltaNeutralService struct {
	c *UnifiedClient
}

func (c *UnifiedClient) NewGetUnifiedDeltaNeutralService() *GetUnifiedDeltaNeutralService {
	return &GetUnifiedDeltaNeutralService{c: c}
}

func (s *GetUnifiedDeltaNeutralService) Do(ctx context.Context) (*UnifiedDeltaNeutral, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/delta_neutral").WithSign()
	return request.Do[UnifiedDeltaNeutral](req)
}

// SetUnifiedDeltaNeutralService -- POST /api/v4/unified/delta_neutral (private)
//
// Enables or disables the account's Delta-neutral strategy mode. Requires VIP
// level >= 4 and the account in cross-currency margin mode.
type SetUnifiedDeltaNeutralService struct {
	c    *UnifiedClient
	body map[string]any
}

// NewSetUnifiedDeltaNeutralService toggles the Delta-neutral strategy mode.
func (c *UnifiedClient) NewSetUnifiedDeltaNeutralService(enabled bool) *SetUnifiedDeltaNeutralService {
	return &SetUnifiedDeltaNeutralService{c: c, body: map[string]any{"enabled": enabled}}
}

func (s *SetUnifiedDeltaNeutralService) Do(ctx context.Context) (*UnifiedDeltaNeutral, error) {
	req := request.Post(ctx, s.c, "/api/v4/unified/delta_neutral", s.body).WithSign()
	return request.Do[UnifiedDeltaNeutral](req)
}

// UnifiedDeltaNeutral is the account's Delta-neutral strategy mode setting.
type UnifiedDeltaNeutral struct {
	Enabled bool `json:"enabled"`
}

// GetEstimatedQuickRepaymentService -- GET /api/v4/unified/estimated_quick_repayment (private)
//
// Returns the estimated quick-repayment details: each outstanding liability and
// the currencies available to repay it. Applies only to unified accounts in
// cross-currency or portfolio margin mode.
type GetEstimatedQuickRepaymentService struct {
	c *UnifiedClient
}

func (c *UnifiedClient) NewGetEstimatedQuickRepaymentService() *GetEstimatedQuickRepaymentService {
	return &GetEstimatedQuickRepaymentService{c: c}
}

func (s *GetEstimatedQuickRepaymentService) Do(ctx context.Context) (*EstimatedQuickRepayment, error) {
	req := request.Get(ctx, s.c, "/api/v4/unified/estimated_quick_repayment").WithSign()
	return request.Do[EstimatedQuickRepayment](req)
}

// EstimatedQuickRepayment is the estimated quick-repayment overview.
type EstimatedQuickRepayment struct {
	DebtCurrencies []UnifiedDebtCurrencies `json:"debt_currencies"`
}

// UnifiedDebtCurrencies is a single liability currency and the currencies
// available to repay it.
type UnifiedDebtCurrencies struct {
	Currency            string                       `json:"currency"`
	DebtAmount          decimal.Decimal              `json:"debt_amount"`
	EstimatedUSD        decimal.Decimal              `json:"estimated_usd"`
	Borrowed            decimal.Decimal              `json:"borrowed"`
	NegBalance          decimal.Decimal              `json:"neg_balance"`
	AvailableCurrencies []UnifiedAvailableCurrencies `json:"available_currencies"`
}

// UnifiedAvailableCurrencies is a currency available to repay a liability.
type UnifiedAvailableCurrencies struct {
	Currency     string          `json:"currency"`
	Available    decimal.Decimal `json:"available"`
	EstimatedUSD decimal.Decimal `json:"estimated_usd"`
}

// CreateQuickRepaymentService -- POST /api/v4/unified/quick_repayment (private)
//
// Executes a quick repayment, clearing the listed liability currencies using the
// listed available currencies.
type CreateQuickRepaymentService struct {
	c    *UnifiedClient
	body map[string]any
}

// NewCreateQuickRepaymentService repays debtCurrencies (liability currencies)
// using availableCurrencies (currencies to repay with).
func (c *UnifiedClient) NewCreateQuickRepaymentService(debtCurrencies, availableCurrencies []string) *CreateQuickRepaymentService {
	return &CreateQuickRepaymentService{c: c, body: map[string]any{
		"debt_currencies":      debtCurrencies,
		"available_currencies": availableCurrencies,
	}}
}

func (s *CreateQuickRepaymentService) Do(ctx context.Context) (*QuickRepaymentResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/unified/quick_repayment", s.body).WithSign()
	return request.Do[QuickRepaymentResult](req)
}

// QuickRepaymentResult is the result of a quick repayment.
type QuickRepaymentResult struct {
	OrderID     string       `json:"order_id"`
	RepaidInfos []RepaidInfo `json:"repaid_infos"`
}

// RepaidInfo is one repaid currency's details, including the currencies used to
// repay it.
type RepaidInfo struct {
	Currency  string          `json:"currency"`
	Repaid    decimal.Decimal `json:"repaid"`
	Left      decimal.Decimal `json:"left"`
	UsedInfos []UsedInfo      `json:"used_infos"`
}

// UsedInfo is a currency consumed to repay a liability.
type UsedInfo struct {
	Currency string          `json:"currency"`
	Used     decimal.Decimal `json:"used"`
}
