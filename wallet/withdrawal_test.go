package wallet

import "testing"

// TestWalletWithdrawal intentionally does not exercise the withdrawal endpoints:
// Withdraw, WithdrawPushOrder and CancelWithdrawal move funds irreversibly, so
// they are implemented and compiled but never sent from the test suite.
func TestWalletWithdrawal(t *testing.T) {
	t.Skip("withdrawal endpoints are implemented but never executed")
}
