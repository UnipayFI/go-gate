package wallet

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListWithdrawalsService -- GET /api/v4/wallet/withdrawals (private)
//
// Returns the authenticated account's withdrawal history. The query time range
// cannot exceed 30 days.
type ListWithdrawalsService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListWithdrawalsService() *ListWithdrawalsService {
	return &ListWithdrawalsService{c: c, params: map[string]string{}}
}

// SetCurrency filters the result to a single currency (e.g. USDT).
func (s *ListWithdrawalsService) SetCurrency(currency string) *ListWithdrawalsService {
	s.params["currency"] = currency
	return s
}

// SetWithdrawID filters the result to a single withdrawal record ID.
func (s *ListWithdrawalsService) SetWithdrawID(withdrawID string) *ListWithdrawalsService {
	s.params["withdraw_id"] = withdrawID
	return s
}

// SetAssetClass filters the result to a single asset class.
func (s *ListWithdrawalsService) SetAssetClass(assetClass string) *ListWithdrawalsService {
	s.params["asset_class"] = assetClass
	return s
}

// SetWithdrawOrderID filters the result to a single client withdraw order ID.
func (s *ListWithdrawalsService) SetWithdrawOrderID(withdrawOrderID string) *ListWithdrawalsService {
	s.params["withdraw_order_id"] = withdrawOrderID
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListWithdrawalsService) SetFrom(from time.Time) *ListWithdrawalsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range.
func (s *ListWithdrawalsService) SetTo(to time.Time) *ListWithdrawalsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *ListWithdrawalsService) SetLimit(limit int) *ListWithdrawalsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListWithdrawalsService) SetOffset(offset int) *ListWithdrawalsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListWithdrawalsService) Do(ctx context.Context) ([]WithdrawalRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/withdrawals", s.params).WithSign()
	resp, err := request.Do[[]WithdrawalRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// WithdrawalRecord is a single withdrawal history entry.
type WithdrawalRecord struct {
	ID              string          `json:"id"`
	TxID            string          `json:"txid"`
	BlockNumber     string          `json:"block_number"`
	WithdrawOrderID string          `json:"withdraw_order_id"`
	Timestamp       time.Time       `json:"timestamp,string,format:unix"`
	Amount          decimal.Decimal `json:"amount"`
	Fee             decimal.Decimal `json:"fee"`
	Currency        string          `json:"currency"`
	Address         string          `json:"address"`
	FailReason      string          `json:"fail_reason"`
	// Timestamp2 is the withdrawal final time (cancellation or success). Gate
	// returns "" while the withdrawal is still in progress, which a unix time.Time
	// cannot decode, so it is kept as a raw string.
	Timestamp2 string `json:"timestamp2"`
	Memo       string `json:"memo"`
	Status     string `json:"status"`
	Chain      string `json:"chain"`
}

// ListDepositsService -- GET /api/v4/wallet/deposits (private)
//
// Returns the authenticated account's deposit history. The query time range
// cannot exceed 30 days.
type ListDepositsService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListDepositsService() *ListDepositsService {
	return &ListDepositsService{c: c, params: map[string]string{}}
}

// SetCurrency filters the result to a single currency (e.g. USDT).
func (s *ListDepositsService) SetCurrency(currency string) *ListDepositsService {
	s.params["currency"] = currency
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListDepositsService) SetFrom(from time.Time) *ListDepositsService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range.
func (s *ListDepositsService) SetTo(to time.Time) *ListDepositsService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned (max 500).
func (s *ListDepositsService) SetLimit(limit int) *ListDepositsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListDepositsService) SetOffset(offset int) *ListDepositsService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *ListDepositsService) Do(ctx context.Context) ([]DepositRecord, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/deposits", s.params).WithSign()
	resp, err := request.Do[[]DepositRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DepositRecord is a single deposit history entry.
type DepositRecord struct {
	ID              string          `json:"id"`
	TxID            string          `json:"txid"`
	WithdrawOrderID string          `json:"withdraw_order_id"`
	Timestamp       time.Time       `json:"timestamp,string,format:unix"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	Address         string          `json:"address"`
	Memo            string          `json:"memo"`
	Status          string          `json:"status"`
	Chain           string          `json:"chain"`
}

// ListWithdrawStatusService -- GET /api/v4/wallet/withdraw_status (private)
//
// Returns per-currency deposit/withdrawal fees and limits for the authenticated
// account.
type ListWithdrawStatusService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListWithdrawStatusService() *ListWithdrawStatusService {
	return &ListWithdrawStatusService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListWithdrawStatusService) SetCurrency(currency string) *ListWithdrawStatusService {
	s.params["currency"] = currency
	return s
}

func (s *ListWithdrawStatusService) Do(ctx context.Context) ([]WithdrawStatus, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/withdraw_status", s.params).WithSign()
	resp, err := request.Do[[]WithdrawStatus](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// WithdrawStatus is the per-currency withdrawal fee/limit schedule.
type WithdrawStatus struct {
	Currency string          `json:"currency"`
	Name     string          `json:"name"`
	NameCn   string          `json:"name_cn"`
	Deposit  decimal.Decimal `json:"deposit"`
	// WithdrawPercent is a percentage rate that Gate returns with a "%" suffix
	// (e.g. "0%"), so it is kept as a raw string rather than a decimal.
	WithdrawPercent         string            `json:"withdraw_percent"`
	WithdrawFix             decimal.Decimal   `json:"withdraw_fix"`
	WithdrawDayLimit        decimal.Decimal   `json:"withdraw_day_limit"`
	WithdrawAmountMini      decimal.Decimal   `json:"withdraw_amount_mini"`
	WithdrawDayLimitRemain  decimal.Decimal   `json:"withdraw_day_limit_remain"`
	WithdrawEachtimeLimit   decimal.Decimal   `json:"withdraw_eachtime_limit"`
	WithdrawFixOnChains     map[string]string `json:"withdraw_fix_on_chains"`
	WithdrawPercentOnChains map[string]string `json:"withdraw_percent_on_chains"`
}

// ListSubAccountBalancesService -- GET /api/v4/wallet/sub_account_balances (private)
//
// Returns the spot balances of the main account's sub-accounts.
type ListSubAccountBalancesService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSubAccountBalancesService() *ListSubAccountBalancesService {
	return &ListSubAccountBalancesService{c: c, params: map[string]string{}}
}

// SetSubUID narrows the result to one or more sub-account user IDs (comma-separated).
func (s *ListSubAccountBalancesService) SetSubUID(subUID string) *ListSubAccountBalancesService {
	s.params["sub_uid"] = subUID
	return s
}

func (s *ListSubAccountBalancesService) Do(ctx context.Context) ([]SubAccountBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/sub_account_balances", s.params).WithSign()
	resp, err := request.Do[[]SubAccountBalance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubAccountBalance is a sub-account's spot balances, keyed by currency.
type SubAccountBalance struct {
	UID       string                     `json:"uid"`
	Available map[string]decimal.Decimal `json:"available"`
}

// ListSubAccountMarginBalancesService -- GET /api/v4/wallet/sub_account_margin_balances (private)
//
// Returns the isolated-margin balances of the main account's sub-accounts.
type ListSubAccountMarginBalancesService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSubAccountMarginBalancesService() *ListSubAccountMarginBalancesService {
	return &ListSubAccountMarginBalancesService{c: c, params: map[string]string{}}
}

// SetSubUID narrows the result to one or more sub-account user IDs (comma-separated).
func (s *ListSubAccountMarginBalancesService) SetSubUID(subUID string) *ListSubAccountMarginBalancesService {
	s.params["sub_uid"] = subUID
	return s
}

func (s *ListSubAccountMarginBalancesService) Do(ctx context.Context) ([]SubAccountMarginBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/sub_account_margin_balances", s.params).WithSign()
	resp, err := request.Do[[]SubAccountMarginBalance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubAccountMarginBalance is a sub-account's isolated-margin balances.
type SubAccountMarginBalance struct {
	UID       string             `json:"uid"`
	Available []SubMarginAccount `json:"available"`
}

// SubMarginAccount is one trading pair's isolated-margin account.
type SubMarginAccount struct {
	CurrencyPair string                   `json:"currency_pair"`
	AccountType  string                   `json:"account_type"`
	Leverage     decimal.Decimal          `json:"leverage"`
	Locked       bool                     `json:"locked"`
	Risk         decimal.Decimal          `json:"risk"`
	MMR          decimal.Decimal          `json:"mmr"`
	Base         SubMarginAccountCurrency `json:"base"`
	Quote        SubMarginAccountCurrency `json:"quote"`
}

// SubMarginAccountCurrency is one currency side (base/quote) of a margin account.
type SubMarginAccountCurrency struct {
	Currency  string          `json:"currency"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
	Borrowed  decimal.Decimal `json:"borrowed"`
	Interest  decimal.Decimal `json:"interest"`
}

// ListSubAccountFuturesBalancesService -- GET /api/v4/wallet/sub_account_futures_balances (private)
//
// Returns the perpetual-futures balances of the main account's sub-accounts.
type ListSubAccountFuturesBalancesService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSubAccountFuturesBalancesService() *ListSubAccountFuturesBalancesService {
	return &ListSubAccountFuturesBalancesService{c: c, params: map[string]string{}}
}

// SetSubUID narrows the result to one or more sub-account user IDs (comma-separated).
func (s *ListSubAccountFuturesBalancesService) SetSubUID(subUID string) *ListSubAccountFuturesBalancesService {
	s.params["sub_uid"] = subUID
	return s
}

// SetSettle narrows the result to a single settlement currency (e.g. usdt).
func (s *ListSubAccountFuturesBalancesService) SetSettle(settle string) *ListSubAccountFuturesBalancesService {
	s.params["settle"] = settle
	return s
}

func (s *ListSubAccountFuturesBalancesService) Do(ctx context.Context) ([]SubAccountFuturesBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/sub_account_futures_balances", s.params).WithSign()
	resp, err := request.Do[[]SubAccountFuturesBalance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubAccountFuturesBalance is a sub-account's perpetual-futures balances, keyed
// by settlement currency.
type SubAccountFuturesBalance struct {
	UID       string                       `json:"uid"`
	Available map[string]SubFuturesAccount `json:"available"`
}

// SubFuturesAccount is a perpetual-futures account for one settlement currency.
type SubFuturesAccount struct {
	Total                  decimal.Decimal          `json:"total"`
	UnrealisedPnL          decimal.Decimal          `json:"unrealised_pnl"`
	PositionMargin         decimal.Decimal          `json:"position_margin"`
	OrderMargin            decimal.Decimal          `json:"order_margin"`
	Available              decimal.Decimal          `json:"available"`
	Point                  decimal.Decimal          `json:"point"`
	Currency               string                   `json:"currency"`
	InDualMode             bool                     `json:"in_dual_mode"`
	PositionMode           string                   `json:"position_mode"`
	EnableCredit           bool                     `json:"enable_credit"`
	PositionInitialMargin  decimal.Decimal          `json:"position_initial_margin"`
	MaintenanceMargin      decimal.Decimal          `json:"maintenance_margin"`
	Bonus                  decimal.Decimal          `json:"bonus"`
	EnableEvolvedClassic   bool                     `json:"enable_evolved_classic"`
	CrossOrderMargin       decimal.Decimal          `json:"cross_order_margin"`
	CrossInitialMargin     decimal.Decimal          `json:"cross_initial_margin"`
	CrossMaintenanceMargin decimal.Decimal          `json:"cross_maintenance_margin"`
	CrossUnrealisedPnL     decimal.Decimal          `json:"cross_unrealised_pnl"`
	CrossAvailable         decimal.Decimal          `json:"cross_available"`
	CrossMarginBalance     decimal.Decimal          `json:"cross_margin_balance"`
	CrossMMR               decimal.Decimal          `json:"cross_mmr"`
	CrossIMR               decimal.Decimal          `json:"cross_imr"`
	IsolatedPositionMargin decimal.Decimal          `json:"isolated_position_margin"`
	EnableNewDualMode      bool                     `json:"enable_new_dual_mode"`
	MarginMode             int                      `json:"margin_mode"`
	EnableTieredMM         bool                     `json:"enable_tiered_mm"`
	History                SubFuturesAccountHistory `json:"history"`
}

// SubFuturesAccountHistory is the cumulative statistics of a futures account.
type SubFuturesAccountHistory struct {
	DNW         decimal.Decimal `json:"dnw"`
	PnL         decimal.Decimal `json:"pnl"`
	Fee         decimal.Decimal `json:"fee"`
	Refr        decimal.Decimal `json:"refr"`
	Fund        decimal.Decimal `json:"fund"`
	PointDNW    decimal.Decimal `json:"point_dnw"`
	PointFee    decimal.Decimal `json:"point_fee"`
	PointRefr   decimal.Decimal `json:"point_refr"`
	BonusDNW    decimal.Decimal `json:"bonus_dnw"`
	BonusOffset decimal.Decimal `json:"bonus_offset"`
}

// ListSubAccountCrossMarginBalancesService -- GET /api/v4/wallet/sub_account_cross_margin_balances (private)
//
// Returns the cross-margin balances of the main account's sub-accounts.
type ListSubAccountCrossMarginBalancesService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSubAccountCrossMarginBalancesService() *ListSubAccountCrossMarginBalancesService {
	return &ListSubAccountCrossMarginBalancesService{c: c, params: map[string]string{}}
}

// SetSubUID narrows the result to one or more sub-account user IDs (comma-separated).
func (s *ListSubAccountCrossMarginBalancesService) SetSubUID(subUID string) *ListSubAccountCrossMarginBalancesService {
	s.params["sub_uid"] = subUID
	return s
}

func (s *ListSubAccountCrossMarginBalancesService) Do(ctx context.Context) ([]SubAccountCrossMarginBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/sub_account_cross_margin_balances", s.params).WithSign()
	resp, err := request.Do[[]SubAccountCrossMarginBalance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubAccountCrossMarginBalance is a sub-account's cross-margin account.
type SubAccountCrossMarginBalance struct {
	UID       string                `json:"uid"`
	Available SubCrossMarginAccount `json:"available"`
}

// SubCrossMarginAccount is a cross-margin account and its per-currency balances.
type SubCrossMarginAccount struct {
	UserID                     int64                            `json:"user_id"`
	Locked                     bool                             `json:"locked"`
	Balances                   map[string]SubCrossMarginBalance `json:"balances"`
	Total                      decimal.Decimal                  `json:"total"`
	Borrowed                   decimal.Decimal                  `json:"borrowed"`
	BorrowedNet                decimal.Decimal                  `json:"borrowed_net"`
	Net                        decimal.Decimal                  `json:"net"`
	Leverage                   decimal.Decimal                  `json:"leverage"`
	Interest                   decimal.Decimal                  `json:"interest"`
	Risk                       decimal.Decimal                  `json:"risk"`
	TotalInitialMargin         decimal.Decimal                  `json:"total_initial_margin"`
	TotalMarginBalance         decimal.Decimal                  `json:"total_margin_balance"`
	TotalMaintenanceMargin     decimal.Decimal                  `json:"total_maintenance_margin"`
	TotalInitialMarginRate     decimal.Decimal                  `json:"total_initial_margin_rate"`
	TotalMaintenanceMarginRate decimal.Decimal                  `json:"total_maintenance_margin_rate"`
	TotalAvailableMargin       decimal.Decimal                  `json:"total_available_margin"`
}

// SubCrossMarginBalance is one currency's cross-margin balance.
type SubCrossMarginBalance struct {
	Available decimal.Decimal `json:"available"`
	Freeze    decimal.Decimal `json:"freeze"`
	Borrowed  decimal.Decimal `json:"borrowed"`
	Interest  decimal.Decimal `json:"interest"`
}

// ListSmallBalanceService -- GET /api/v4/wallet/small_balance (private)
//
// Returns the currencies whose small (dust) balances can be converted to GT.
type ListSmallBalanceService struct {
	c *WalletClient
}

func (c *WalletClient) NewListSmallBalanceService() *ListSmallBalanceService {
	return &ListSmallBalanceService{c: c}
}

func (s *ListSmallBalanceService) Do(ctx context.Context) ([]SmallBalance, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/small_balance").WithSign()
	resp, err := request.Do[[]SmallBalance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SmallBalance is a convertible dust balance in a single currency.
type SmallBalance struct {
	Currency         string          `json:"currency"`
	AvailableBalance decimal.Decimal `json:"available_balance"`
	EstimatedAsBTC   decimal.Decimal `json:"estimated_as_btc"`
	ConvertibleToGT  decimal.Decimal `json:"convertible_to_gt"`
}

// ConvertSmallBalanceService -- POST /api/v4/wallet/small_balance (private)
//
// Converts the dust balances of the given currencies to GT.
type ConvertSmallBalanceService struct {
	c    *WalletClient
	body map[string]any
}

func (c *WalletClient) NewConvertSmallBalanceService(currency []string) *ConvertSmallBalanceService {
	return &ConvertSmallBalanceService{c: c, body: map[string]any{"currency": currency}}
}

// SetIsAll converts every convertible dust currency, ignoring the currency list.
func (s *ConvertSmallBalanceService) SetIsAll(isAll bool) *ConvertSmallBalanceService {
	s.body["is_all"] = isAll
	return s
}

func (s *ConvertSmallBalanceService) Do(ctx context.Context) (*ConvertSmallBalanceResult, error) {
	req := request.Post(ctx, s.c, "/api/v4/wallet/small_balance", s.body).WithSign()
	return request.Do[ConvertSmallBalanceResult](req)
}

// ConvertSmallBalanceResult is the (empty) response of a dust conversion: Gate
// replies 200 with no body.
type ConvertSmallBalanceResult struct{}

// ListSmallBalanceHistoryService -- GET /api/v4/wallet/small_balance_history (private)
//
// Returns the history of dust-to-GT conversions.
type ListSmallBalanceHistoryService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListSmallBalanceHistoryService() *ListSmallBalanceHistoryService {
	return &ListSmallBalanceHistoryService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency (e.g. USDT).
func (s *ListSmallBalanceHistoryService) SetCurrency(currency string) *ListSmallBalanceHistoryService {
	s.params["currency"] = currency
	return s
}

// SetPage sets the page of records to return, starting from 1.
func (s *ListSmallBalanceHistoryService) SetPage(page int) *ListSmallBalanceHistoryService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned per page.
func (s *ListSmallBalanceHistoryService) SetLimit(limit int) *ListSmallBalanceHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *ListSmallBalanceHistoryService) Do(ctx context.Context) ([]SmallBalanceHistory, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/small_balance_history", s.params).WithSign()
	resp, err := request.Do[[]SmallBalanceHistory](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SmallBalanceHistory is a single dust-to-GT conversion record.
type SmallBalanceHistory struct {
	ID         string          `json:"id"`
	Currency   string          `json:"currency"`
	Amount     decimal.Decimal `json:"amount"`
	GTAmount   decimal.Decimal `json:"gt_amount"`
	CreateTime time.Time       `json:"create_time,format:unix"`
}

// ListPushOrdersService -- GET /api/v4/wallet/push (private)
//
// Returns the account's UID push (peer-to-peer transfer) orders.
type ListPushOrdersService struct {
	c      *WalletClient
	params map[string]string
}

func (c *WalletClient) NewListPushOrdersService() *ListPushOrdersService {
	return &ListPushOrdersService{c: c, params: map[string]string{}}
}

// SetID filters the result to a single push order ID.
func (s *ListPushOrdersService) SetID(id int64) *ListPushOrdersService {
	s.params["id"] = strconv.FormatInt(id, 10)
	return s
}

// SetFrom sets the start of the query time range.
func (s *ListPushOrdersService) SetFrom(from time.Time) *ListPushOrdersService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end of the query time range.
func (s *ListPushOrdersService) SetTo(to time.Time) *ListPushOrdersService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned.
func (s *ListPushOrdersService) SetLimit(limit int) *ListPushOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *ListPushOrdersService) SetOffset(offset int) *ListPushOrdersService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetTransactionType filters the result to a single transaction type.
func (s *ListPushOrdersService) SetTransactionType(transactionType string) *ListPushOrdersService {
	s.params["transaction_type"] = transactionType
	return s
}

func (s *ListPushOrdersService) Do(ctx context.Context) ([]PushOrder, error) {
	req := request.Get(ctx, s.c, "/api/v4/wallet/push", s.params).WithSign()
	resp, err := request.Do[[]PushOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PushOrder is a single UID push (peer-to-peer transfer) order.
type PushOrder struct {
	ID              int64           `json:"id"`
	PushUID         int64           `json:"push_uid"`
	ReceiveUID      int64           `json:"receive_uid"`
	Currency        string          `json:"currency"`
	Amount          decimal.Decimal `json:"amount"`
	CreateTime      time.Time       `json:"create_time,format:unix"`
	Status          string          `json:"status"`
	Message         string          `json:"message"`
	TransactionType string          `json:"transaction_type"`
}
