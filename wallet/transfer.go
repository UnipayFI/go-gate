package wallet

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListCurrencyChainsService -- GET /api/v4/wallet/currency_chains
//
// Returns the deposit/withdrawal chains supported for a currency.
type ListCurrencyChainsService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListCurrencyChainsService(currency string) *ListCurrencyChainsService {
	return &ListCurrencyChainsService{c: c, params: map[string]string{"currency": currency}}
}

func (s *ListCurrencyChainsService) Do(ctx context.Context) ([]CurrencyChain, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/currency_chains", s.params)
	resp, err := request.Do[[]CurrencyChain](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetDepositAddressService -- GET /api/v4/wallet/deposit_address (private)
//
// Generates (or returns) the deposit address of a currency for the account.
type GetDepositAddressService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewGetDepositAddressService(currency string) *GetDepositAddressService {
	return &GetDepositAddressService{c: c, params: map[string]string{"currency": currency}}
}

func (s *GetDepositAddressService) Do(ctx context.Context) (*DepositAddress, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/deposit_address", s.params).WithSign()
	return request.Do[DepositAddress](req)
}

// TransferService -- POST /api/v4/wallet/transfers (private)
//
// Transfers a balance between the account's own trading accounts (spot, margin,
// perpetual/delivery futures, options).
type TransferService struct {
	c    *WalletClient
	body map[string]any
}

func (c *WalletClient) NewTransferService(currency, from, to string, amount decimal.Decimal) *TransferService {
	return &TransferService{c: c, body: map[string]any{
		"currency": currency,
		"from":     from,
		"to":       to,
		"amount":   amount.String(),
	}}
}

// SetCurrencyPair sets the margin trading pair, required when transferring to or
// from a margin account.
func (s *TransferService) SetCurrencyPair(currencyPair string) *TransferService {
	s.body["currency_pair"] = currencyPair
	return s
}

// SetSettle sets the contract settlement currency, required when transferring to
// or from a contract account.
func (s *TransferService) SetSettle(settle string) *TransferService {
	s.body["settle"] = settle
	return s
}

func (s *TransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/wallet/transfers", s.body).WithSign()
	return request.Do[TransferResult](req)
}

// ListSubAccountTransfersService -- GET /api/v4/wallet/sub_account_transfers (private)
//
// Returns the transfer records between the main account and its sub-accounts.
type ListSubAccountTransfersService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSubAccountTransfersService() *ListSubAccountTransfersService {
	return &ListSubAccountTransfersService{c: c, params: map[string]string{}}
}

// SetSubUID narrows the result to one or more sub-account user IDs (comma-separated).
func (s *ListSubAccountTransfersService) SetSubUID(subUID string) *ListSubAccountTransfersService {
	s.params["sub_uid"] = subUID
	return s
}

// SetFrom sets the start time of the query window.
func (s *ListSubAccountTransfersService) SetFrom(from time.Time) *ListSubAccountTransfersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *ListSubAccountTransfersService) SetTo(to time.Time) *ListSubAccountTransfersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned in a single page.
func (s *ListSubAccountTransfersService) SetLimit(limit int) *ListSubAccountTransfersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListSubAccountTransfersService) SetOffset(offset int) *ListSubAccountTransfersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListSubAccountTransfersService) Do(ctx context.Context) ([]SubAccountTransfer, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/sub_account_transfers", s.params).WithSign()
	resp, err := request.Do[[]SubAccountTransfer](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// TransferWithSubAccountService -- POST /api/v4/wallet/sub_account_transfers (private)
//
// Transfers a balance between the main account and a sub-account. Only the main
// account's spot balance is used, regardless of which sub-account is operated.
type TransferWithSubAccountService struct {
	c    *WalletClient
	body map[string]any
}

func (c *WalletClient) NewTransferWithSubAccountService(subAccount, currency string, amount decimal.Decimal, direction string) *TransferWithSubAccountService {
	return &TransferWithSubAccountService{c: c, body: map[string]any{
		"sub_account": subAccount,
		"currency":    currency,
		"amount":      amount.String(),
		"direction":   direction,
	}}
}

// SetSubAccountType selects the sub-account trading account (spot, futures,
// delivery or options). Defaults to spot server-side.
func (s *TransferWithSubAccountService) SetSubAccountType(subAccountType string) *TransferWithSubAccountService {
	s.body["sub_account_type"] = subAccountType
	return s
}

// SetClientOrderID sets a customer-defined ID to prevent duplicate transfers.
func (s *TransferWithSubAccountService) SetClientOrderID(clientOrderID string) *TransferWithSubAccountService {
	s.body["client_order_id"] = clientOrderID
	return s
}

func (s *TransferWithSubAccountService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/wallet/sub_account_transfers", s.body).WithSign()
	return request.Do[TransferResult](req)
}

// SubAccountToSubAccountService -- POST /api/v4/wallet/sub_account_to_sub_account (private)
//
// Transfers a balance directly between two sub-accounts.
type SubAccountToSubAccountService struct {
	c    *WalletClient
	body map[string]any
}

func (c *WalletClient) NewSubAccountToSubAccountService(currency, subAccountFrom, subAccountFromType, subAccountTo, subAccountToType string, amount decimal.Decimal) *SubAccountToSubAccountService {
	return &SubAccountToSubAccountService{c: c, body: map[string]any{
		"currency":              currency,
		"sub_account_from":      subAccountFrom,
		"sub_account_from_type": subAccountFromType,
		"sub_account_to":        subAccountTo,
		"sub_account_to_type":   subAccountToType,
		"amount":                amount.String(),
	}}
}

// SetSubAccountType sets the deprecated transfer account type (prefer the
// sub_account_from_type / sub_account_to_type constructor arguments).
func (s *SubAccountToSubAccountService) SetSubAccountType(subAccountType string) *SubAccountToSubAccountService {
	s.body["sub_account_type"] = subAccountType
	return s
}

func (s *SubAccountToSubAccountService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/wallet/sub_account_to_sub_account", s.body).WithSign()
	return request.Do[TransferResult](req)
}

// GetTransferOrderStatusService -- GET /api/v4/wallet/order_status (private)
//
// Queries the status of a transfer by the tx_id returned from the transfer
// endpoint (or by the customer-defined client_order_id).
type GetTransferOrderStatusService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewGetTransferOrderStatusService(txID string) *GetTransferOrderStatusService {
	return &GetTransferOrderStatusService{c: c, params: map[string]string{"tx_id": txID}}
}

// SetClientOrderID queries by the customer-defined transfer ID instead of tx_id.
func (s *GetTransferOrderStatusService) SetClientOrderID(clientOrderID string) *GetTransferOrderStatusService {
	s.params["client_order_id"] = clientOrderID
	return s
}

func (s *GetTransferOrderStatusService) Do(ctx context.Context) (*TransferOrderStatus, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/order_status", s.params).WithSign()
	return request.Do[TransferOrderStatus](req)
}

// GetTradeFeeService -- GET /api/v4/wallet/fee (private)
//
// Returns the account's personal trading fee rates.
type GetTradeFeeService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewGetTradeFeeService() *GetTradeFeeService {
	return &GetTradeFeeService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the query to a specific spot pair for more accurate fees.
func (s *GetTradeFeeService) SetCurrencyPair(currencyPair string) *GetTradeFeeService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetSettle narrows the query to a specific contract settlement currency.
func (s *GetTradeFeeService) SetSettle(settle string) *GetTradeFeeService {
	s.params["settle"] = settle
	return s
}

func (s *GetTradeFeeService) Do(ctx context.Context) (*TradeFee, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/fee", s.params).WithSign()
	return request.Do[TradeFee](req)
}

// GetTotalBalanceService -- GET /api/v4/wallet/total_balance (private)
//
// Returns the estimated total value of every account, converted to the target
// currency. Values may be cached for up to a minute.
type GetTotalBalanceService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewGetTotalBalanceService() *GetTotalBalanceService {
	return &GetTotalBalanceService{c: c, params: map[string]string{}}
}

// SetCurrency selects the conversion currency (BTC, CNY, USD or USDT; USDT default).
func (s *GetTotalBalanceService) SetCurrency(currency string) *GetTotalBalanceService {
	s.params["currency"] = currency
	return s
}

func (s *GetTotalBalanceService) Do(ctx context.Context) (*TotalBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/total_balance", s.params).WithSign()
	return request.Do[TotalBalance](req)
}

// ListSavedAddressService -- GET /api/v4/wallet/saved_address (private)
//
// Returns the withdrawal address whitelist for a currency.
type ListSavedAddressService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSavedAddressService(currency string) *ListSavedAddressService {
	return &ListSavedAddressService{c: c, params: map[string]string{"currency": currency}}
}

// SetChain narrows the whitelist to a single chain.
func (s *ListSavedAddressService) SetChain(chain string) *ListSavedAddressService {
	s.params["chain"] = chain
	return s
}

// SetLimit caps the number of addresses returned (up to 100).
func (s *ListSavedAddressService) SetLimit(limit int) *ListSavedAddressService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetPage sets the page number.
func (s *ListSavedAddressService) SetPage(page int) *ListSavedAddressService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

func (s *ListSavedAddressService) Do(ctx context.Context) ([]SavedAddress, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/saved_address", s.params).WithSign()
	resp, err := request.Do[[]SavedAddress](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CurrencyChain is a deposit/withdrawal chain supported for a currency.
type CurrencyChain struct {
	Chain              string `json:"chain"`
	NameCN             string `json:"name_cn"`
	NameEN             string `json:"name_en"`
	ContractAddress    string `json:"contract_address"`
	IsDisabled         int    `json:"is_disabled"`
	IsDepositDisabled  int    `json:"is_deposit_disabled"`
	IsWithdrawDisabled int    `json:"is_withdraw_disabled"`
	IsTag              int    `json:"is_tag"`
	Decimal            string `json:"decimal"`
}

// DepositAddress is the deposit address of a currency, including the per-chain
// addresses for multi-chain currencies.
type DepositAddress struct {
	Currency            string              `json:"currency"`
	Address             string              `json:"address"`
	MinConfirms         int                 `json:"min_confirms"`
	MinDepositAmount    decimal.Decimal     `json:"min_deposit_amount"`
	MultichainAddresses []MultiChainAddress `json:"multichain_addresses"`
}

// MultiChainAddress is a single chain's deposit address for a multi-chain currency.
type MultiChainAddress struct {
	Chain        string `json:"chain"`
	Address      string `json:"address"`
	PaymentID    string `json:"payment_id"`
	PaymentName  string `json:"payment_name"`
	ObtainFailed int    `json:"obtain_failed"`
	MinConfirms  int    `json:"min_confirms"`
}

// TransferResult is the tx_id assigned to an accepted transfer.
type TransferResult struct {
	TxID int64 `json:"tx_id"`
}

// SubAccountTransfer is a single transfer record between the main and a sub-account.
type SubAccountTransfer struct {
	Timest         time.Time       `json:"timest,string,format:unix"`
	UID            string          `json:"uid"`
	SubAccount     string          `json:"sub_account"`
	SubAccountType string          `json:"sub_account_type"`
	Currency       string          `json:"currency"`
	Amount         decimal.Decimal `json:"amount"`
	Direction      string          `json:"direction"`
	Source         string          `json:"source"`
	ClientOrderID  string          `json:"client_order_id"`
	Status         string          `json:"status"`
}

// TransferOrderStatus is the status of a transfer looked up by tx_id.
type TransferOrderStatus struct {
	TxID   string `json:"tx_id"`
	Status string `json:"status"`
}

// TradeFee is the account's personal trading fee rates across products.
type TradeFee struct {
	UserID           int64           `json:"user_id"`
	TakerFee         decimal.Decimal `json:"taker_fee"`
	MakerFee         decimal.Decimal `json:"maker_fee"`
	GTDiscount       bool            `json:"gt_discount"`
	GTTakerFee       decimal.Decimal `json:"gt_taker_fee"`
	GTMakerFee       decimal.Decimal `json:"gt_maker_fee"`
	LoanFee          decimal.Decimal `json:"loan_fee"`
	PointType        string          `json:"point_type"`
	FuturesTakerFee  decimal.Decimal `json:"futures_taker_fee"`
	FuturesMakerFee  decimal.Decimal `json:"futures_maker_fee"`
	DeliveryTakerFee decimal.Decimal `json:"delivery_taker_fee"`
	DeliveryMakerFee decimal.Decimal `json:"delivery_maker_fee"`
	DebitFee         int             `json:"debit_fee"`
	RPIMakerFee      decimal.Decimal `json:"rpi_maker_fee"`
	RPIMM            decimal.Decimal `json:"rpi_mm"`
}

// TotalBalance is the estimated total value of every account, converted to a
// single currency.
type TotalBalance struct {
	Total   AccountBalance            `json:"total"`
	Details map[string]AccountBalance `json:"details"`
}

// AccountBalance is a total balance calculated in the requested currency unit.
type AccountBalance struct {
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	UnrealisedPnL decimal.Decimal `json:"unrealised_pnl"`
	Borrowed      decimal.Decimal `json:"borrowed"`
}

// SavedAddress is a whitelisted withdrawal address.
type SavedAddress struct {
	Currency string `json:"currency"`
	Chain    string `json:"chain"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	Verified string `json:"verified"`
}
