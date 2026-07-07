package p2p

import (
	"context"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetUserInfoService -- POST /api/v4/p2p/merchant/account/get_user_info (private)
//
// Returns the authenticated merchant's own P2P account information and trading
// statistics.
type GetUserInfoService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetUserInfoService() *GetUserInfoService {
	return &GetUserInfoService{c: c, body: map[string]any{}}
}

func (s *GetUserInfoService) Do(ctx context.Context) (*P2PResponse[P2PAccountInfo], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/account/get_user_info", s.body).WithSign()
	return doP2P[P2PAccountInfo](req)
}

// P2PAccountInfo is the authenticated merchant's own P2P profile and stats.
type P2PAccountInfo struct {
	IsSelf                    bool            `json:"is_self"`
	UserTimest                string          `json:"user_timest"`
	CounterpartiesNum         int             `json:"counterparties_num"`
	EmailVerified             string          `json:"email_verified"`
	Verified                  string          `json:"verified"`
	HasPhone                  string          `json:"has_phone"`
	UserName                  string          `json:"user_name"`
	UserNote                  string          `json:"user_note"`
	CompleteTransactions      string          `json:"complete_transactions"`
	PaidTransactions          string          `json:"paid_transactions"`
	AcceptedTransactions      string          `json:"accepted_transactions"`
	TransactionsUsedTime      string          `json:"transactions_used_time"`
	CancelledUsedTimeMonth    string          `json:"cancelled_used_time_month"`
	CompleteTransactionsMonth string          `json:"complete_transactions_month"`
	CompleteRateMonth         decimal.Decimal `json:"complete_rate_month"`
	OrdersBuyRateMonth        decimal.Decimal `json:"orders_buy_rate_month"`
	IsBlack                   int             `json:"is_black"`
	IsFollow                  int             `json:"is_follow"`
	HaveTraded                int             `json:"have_traded"`
	BizUID                    string          `json:"biz_uid"`
	BlueVip                   int             `json:"blue_vip"`
	WorkStatus                int             `json:"work_status"`
	RegistrationDays          int             `json:"registration_days"`
	FirstTradeDays            int             `json:"first_trade_days"`
	NeedReplenish             int             `json:"need_replenish"`
	MerchantInfo              P2PMerchantInfo `json:"merchant_info"`
	OnlineStatus              int             `json:"online_status"`
	WorkHours                 *P2PWorkHours   `json:"work_hours"`
	TransactionsMonth         decimal.Decimal `json:"transactions_month"`
	TransactionsAll           decimal.Decimal `json:"transactions_all"`
	TradeVersatile            bool            `json:"trade_versatile"`
}

// P2PMerchantInfo describes the markets where a merchant can place orders.
type P2PMerchantInfo struct {
	Type   string `json:"type"`
	Market string `json:"market"`
}

// P2PWorkHours holds a merchant's custom working-hours configuration. It is null
// for merchants who use the default (always-on / manual) working status.
type P2PWorkHours struct {
	WorkStatus int    `json:"work_status"`
	CycleType  string `json:"cycle_type"`
	DayOfWeek  string `json:"day_of_week"`
	TimeZone   string `json:"time_zone"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// GetCounterpartyUserInfoService -- POST /api/v4/p2p/merchant/account/get_counterparty_user_info (private)
//
// Returns the public P2P profile and trading statistics of a counterparty,
// identified by their encrypted (crypto) UID.
type GetCounterpartyUserInfoService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetCounterpartyUserInfoService(bizUID string) *GetCounterpartyUserInfoService {
	return &GetCounterpartyUserInfoService{c: c, body: map[string]any{
		"biz_uid": bizUID,
	}}
}

func (s *GetCounterpartyUserInfoService) Do(ctx context.Context) (*P2PResponse[P2PCounterpartyUserInfo], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/account/get_counterparty_user_info", s.body).WithSign()
	return doP2P[P2PCounterpartyUserInfo](req)
}

// P2PCounterpartyUserInfo is a counterparty's public P2P profile and stats.
type P2PCounterpartyUserInfo struct {
	UserTimest                string          `json:"user_timest"`
	EmailVerified             string          `json:"email_verified"`
	Verified                  string          `json:"verified"`
	HasPhone                  string          `json:"has_phone"`
	UserName                  string          `json:"user_name"`
	UserNote                  string          `json:"user_note"`
	CompleteTransactions      string          `json:"complete_transactions"`
	PaidTransactions          string          `json:"paid_transactions"`
	AcceptedTransactions      string          `json:"accepted_transactions"`
	TransactionsUsedTime      string          `json:"transactions_used_time"`
	CancelledUsedTimeMonth    string          `json:"cancelled_used_time_month"`
	CompleteTransactionsMonth string          `json:"complete_transactions_month"`
	CompleteRateMonth         decimal.Decimal `json:"complete_rate_month"`
	IsFollow                  int             `json:"is_follow"`
	HaveTraded                int             `json:"have_traded"`
	BizUID                    string          `json:"biz_uid"`
	RegistrationDays          int             `json:"registration_days"`
	FirstTradeDays            int             `json:"first_trade_days"`
	TradeVersatile            bool            `json:"trade_versatile"`
}

// GetMyselfPaymentService -- POST /api/v4/p2p/merchant/account/get_myself_payment (private)
//
// Returns the authenticated merchant's bound payment methods, optionally filtered
// to a single fiat currency.
type GetMyselfPaymentService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewGetMyselfPaymentService() *GetMyselfPaymentService {
	return &GetMyselfPaymentService{c: c, body: map[string]any{}}
}

// SetFiat filters the payment methods to a single fiat currency.
func (s *GetMyselfPaymentService) SetFiat(fiat string) *GetMyselfPaymentService {
	s.body["fiat"] = fiat
	return s
}

func (s *GetMyselfPaymentService) Do(ctx context.Context) (*P2PResponse[[]P2PPaymentMethodGroup], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/account/get_myself_payment", s.body).WithSign()
	return doP2P[[]P2PPaymentMethodGroup](req)
}

// P2PPaymentMethodGroup groups a merchant's bound payment accounts by type.
type P2PPaymentMethodGroup struct {
	PayType string                    `json:"pay_type"`
	PayName string                    `json:"pay_name"`
	IDs     []string                  `json:"ids"`
	List    []P2PPaymentMethodAccount `json:"list"`
}

// P2PPaymentMethodAccount is a single bound payment account (bank card, e-wallet,
// etc.) belonging to the merchant.
type P2PPaymentMethodAccount struct {
	UID          int64  `json:"uid"`
	BankID       string `json:"bankid"`
	Nickname     int64  `json:"nickname"`
	Bankname     string `json:"bankname"`
	Bankbranch   string `json:"bankbranch"`
	Bankcity     string `json:"bankcity"`
	Bankprov     string `json:"bankprov"`
	Bankaddr     string `json:"bankaddr"`
	Bankdesc     string `json:"bankdesc"`
	HoldUID      int64  `json:"hold_uid"`
	HoldUsername string `json:"hold_username"`
	RealName     string `json:"real_name"`
	ID           string `json:"id"`
	AccountDes   string `json:"account_des"`
	PayType      string `json:"pay_type"`
	File         string `json:"file"`
	FileKey      string `json:"file_key"`
	Account      string `json:"account"`
	Memo         string `json:"memo"`
	Code         string `json:"code"`
	MemoExt      string `json:"memo_ext"`
	TradeTips    string `json:"trade_tips"`
}

// SetMerchantWorkHoursService -- POST /api/v4/p2p/merchant/account/set_merchant_work_hours (private)
//
// Sets the merchant's working status and, when using custom hours (work_status
// 2), the custom working-hours schedule.
type SetMerchantWorkHoursService struct {
	c    *P2PClient
	body map[string]any
}

func (c *P2PClient) NewSetMerchantWorkHoursService(workStatus int) *SetMerchantWorkHoursService {
	return &SetMerchantWorkHoursService{c: c, body: map[string]any{
		"work_status": workStatus,
	}}
}

// SetCycleType sets the custom working cycle (required when work_status is 2).
func (s *SetMerchantWorkHoursService) SetCycleType(cycleType string) *SetMerchantWorkHoursService {
	s.body["cycle_type"] = cycleType
	return s
}

// SetDayOfWeek sets the weekly working days as comma-separated values 1-7
// (Monday to Sunday); required when work_status is 2.
func (s *SetMerchantWorkHoursService) SetDayOfWeek(dayOfWeek string) *SetMerchantWorkHoursService {
	s.body["day_of_week"] = dayOfWeek
	return s
}

// SetTimeZone sets the UTC timezone offset (-12 to +14); required when
// work_status is 2.
func (s *SetMerchantWorkHoursService) SetTimeZone(timeZone string) *SetMerchantWorkHoursService {
	s.body["time_zone"] = timeZone
	return s
}

// SetStartTime sets the custom working start time in HH:mm format; required when
// work_status is 2.
func (s *SetMerchantWorkHoursService) SetStartTime(startTime string) *SetMerchantWorkHoursService {
	s.body["start_time"] = startTime
	return s
}

// SetEndTime sets the custom working end time in HH:mm format; required when
// work_status is 2.
func (s *SetMerchantWorkHoursService) SetEndTime(endTime string) *SetMerchantWorkHoursService {
	s.body["end_time"] = endTime
	return s
}

func (s *SetMerchantWorkHoursService) Do(ctx context.Context) (*P2PResponse[P2PWorkStatusData], error) {
	req := request.Post(ctx, s.c, "/api/v4/p2p/merchant/account/set_merchant_work_hours", s.body).WithSign()
	return doP2P[P2PWorkStatusData](req)
}

// P2PWorkStatusData reports the merchant's working status after an update.
type P2PWorkStatusData struct {
	WorkStatus int `json:"work_status"`
}
