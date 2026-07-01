package wallet

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// WithdrawService -- POST /api/v4/withdrawals (private)
//
// Submits an on-chain withdrawal. currency, amount, address and chain are
// required; memo and a client-side withdraw_order_id are optional.
type WithdrawService struct {
	c    *WalletClient
	body map[string]any
}

func (c *WalletClient) NewWithdrawService(currency string, amount decimal.Decimal, address, chain string) *WithdrawService {
	return &WithdrawService{c: c, body: map[string]any{
		"currency": currency,
		"amount":   amount.String(),
		"address":  address,
		"chain":    chain,
	}}
}

// SetMemo sets the additional address memo/tag some chains require.
func (s *WithdrawService) SetMemo(memo string) *WithdrawService {
	s.body["memo"] = memo
	return s
}

// SetWithdrawOrderID sets a user-defined order number for the withdrawal, used
// to look the record up later.
func (s *WithdrawService) SetWithdrawOrderID(withdrawOrderID string) *WithdrawService {
	s.body["withdraw_order_id"] = withdrawOrderID
	return s
}

func (s *WithdrawService) Do(ctx context.Context) (*WithdrawalResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/withdrawals", s.body).WithSign()
	return request.Do[WithdrawalResult](req)
}

// WithdrawPushOrderService -- POST /api/v4/withdrawals/push (private)
//
// Transfers funds to another user by UID. Both parties must be main spot
// accounts (neither can be a sub-account).
type WithdrawPushOrderService struct {
	c    *WalletClient
	body map[string]any
}

func (c *WalletClient) NewWithdrawPushOrderService(receiveUID int64, currency string, amount decimal.Decimal) *WithdrawPushOrderService {
	return &WithdrawPushOrderService{c: c, body: map[string]any{
		"receive_uid": receiveUID,
		"currency":    currency,
		"amount":      amount.String(),
	}}
}

func (s *WithdrawPushOrderService) Do(ctx context.Context) (*WithdrawalRequest, error) {
	req := request.Post(ctx, s.c, "/api/v4/withdrawals/push", s.body).WithSign()
	return request.Do[WithdrawalRequest](req)
}

// CancelWithdrawalService -- DELETE /api/v4/withdrawals/{withdrawal_id} (private)
//
// Cancels a pending withdrawal by its record ID.
type CancelWithdrawalService struct {
	c            *WalletClient
	withdrawalID string
}

func (c *WalletClient) NewCancelWithdrawalService(withdrawalID string) *CancelWithdrawalService {
	return &CancelWithdrawalService{c: c, withdrawalID: withdrawalID}
}

func (s *CancelWithdrawalService) Do(ctx context.Context) (*WithdrawalResult, error) {
	req := request.Delete(ctx, s.c, "/api/v4/withdrawals/"+s.withdrawalID).WithSign()
	return request.Do[WithdrawalResult](req)
}

// WithdrawalResult is a withdrawal ledger record, returned when a withdrawal is
// submitted or cancelled.
type WithdrawalResult struct {
	ID              string          `json:"id"`
	TxID            string          `json:"txid"`
	WithdrawOrderID string          `json:"withdraw_order_id"`
	Timestamp       time.Time       `json:"timestamp,string,format:unix"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	Address         string          `json:"address"`
	Memo            string          `json:"memo"`
	WithdrawID      string          `json:"withdraw_id"`
	AssetClass      string          `json:"asset_class"`
	Status          string          `json:"status"`
	Chain           string          `json:"chain"`
}

// WithdrawalRequest is the created transfer order returned by a UID transfer.
type WithdrawalRequest struct {
	ID int64 `json:"id"`
}
