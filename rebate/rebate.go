package rebate

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// AgencyTransactionHistoryService -- GET /api/v4/rebate/agency/transaction_history (private)
//
// Broker obtains the transaction history of its recommended users. The query
// window may not exceed 30 days.
type AgencyTransactionHistoryService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewAgencyTransactionHistoryService() *AgencyTransactionHistoryService {
	return &AgencyTransactionHistoryService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *AgencyTransactionHistoryService) SetCurrencyPair(currencyPair string) *AgencyTransactionHistoryService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetUserID narrows the result to a single recommended user.
func (s *AgencyTransactionHistoryService) SetUserID(userID int64) *AgencyTransactionHistoryService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetFrom sets the start time of the query window.
func (s *AgencyTransactionHistoryService) SetFrom(from time.Time) *AgencyTransactionHistoryService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *AgencyTransactionHistoryService) SetTo(to time.Time) *AgencyTransactionHistoryService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *AgencyTransactionHistoryService) SetLimit(limit int) *AgencyTransactionHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *AgencyTransactionHistoryService) SetOffset(offset int) *AgencyTransactionHistoryService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *AgencyTransactionHistoryService) Do(ctx context.Context) ([]AgencyTransaction, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/agency/transaction_history", s.params).WithSign()
	resp, err := request.Do[[]AgencyTransaction](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AgencyTransaction is one trading-pair bucket of an agency's recommended-user
// transaction history.
type AgencyTransaction struct {
	CurrencyPair string                    `json:"currency_pair"`
	Total        int64                     `json:"total"`
	List         []AgencyTransactionRecord `json:"list"`
}

// AgencyTransactionRecord is a single recommended-user transaction.
type AgencyTransactionRecord struct {
	TransactionTime time.Time       `json:"transaction_time,format:unix"`
	UserID          int64           `json:"user_id"`
	GroupName       string          `json:"group_name"`
	Fee             decimal.Decimal `json:"fee"`
	FeeAsset        string          `json:"fee_asset"`
	CurrencyPair    string          `json:"currency_pair"`
	Amount          decimal.Decimal `json:"amount"`
	AmountAsset     string          `json:"amount_asset"`
	Source          string          `json:"source"`
}

// AgencyCommissionsHistoryService -- GET /api/v4/rebate/agency/commission_history (private)
//
// Broker obtains the rebate history of its recommended users. The query window
// may not exceed 30 days.
type AgencyCommissionsHistoryService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewAgencyCommissionsHistoryService() *AgencyCommissionsHistoryService {
	return &AgencyCommissionsHistoryService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency.
func (s *AgencyCommissionsHistoryService) SetCurrency(currency string) *AgencyCommissionsHistoryService {
	s.params["currency"] = currency
	return s
}

// SetCommissionType filters by rebate type (1 - direct, 2 - indirect, 3 - self).
func (s *AgencyCommissionsHistoryService) SetCommissionType(commissionType int) *AgencyCommissionsHistoryService {
	s.params["commission_type"] = strconv.Itoa(commissionType)
	return s
}

// SetUserID narrows the result to a single recommended user.
func (s *AgencyCommissionsHistoryService) SetUserID(userID int64) *AgencyCommissionsHistoryService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetFrom sets the start time of the query window.
func (s *AgencyCommissionsHistoryService) SetFrom(from time.Time) *AgencyCommissionsHistoryService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *AgencyCommissionsHistoryService) SetTo(to time.Time) *AgencyCommissionsHistoryService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *AgencyCommissionsHistoryService) SetLimit(limit int) *AgencyCommissionsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *AgencyCommissionsHistoryService) SetOffset(offset int) *AgencyCommissionsHistoryService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *AgencyCommissionsHistoryService) Do(ctx context.Context) ([]AgencyCommission, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/agency/commission_history", s.params).WithSign()
	resp, err := request.Do[[]AgencyCommission](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AgencyCommission is one trading-pair bucket of an agency's recommended-user
// rebate history.
type AgencyCommission struct {
	CurrencyPair string                   `json:"currency_pair"`
	Total        int64                    `json:"total"`
	List         []AgencyCommissionRecord `json:"list"`
}

// AgencyCommissionRecord is a single recommended-user rebate.
type AgencyCommissionRecord struct {
	CommissionTime   time.Time       `json:"commission_time,format:unix"`
	UserID           int64           `json:"user_id"`
	GroupName        string          `json:"group_name"`
	CommissionAmount decimal.Decimal `json:"commission_amount"`
	CommissionAsset  string          `json:"commission_asset"`
	Source           string          `json:"source"`
}

// PartnerTransactionHistoryService -- GET /api/v4/rebate/partner/transaction_history (private)
//
// Partner obtains the transaction history of its recommended users. The query
// window may not exceed 30 days.
type PartnerTransactionHistoryService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewPartnerTransactionHistoryService() *PartnerTransactionHistoryService {
	return &PartnerTransactionHistoryService{c: c, params: map[string]string{}}
}

// SetCurrencyPair narrows the result to a single trading pair.
func (s *PartnerTransactionHistoryService) SetCurrencyPair(currencyPair string) *PartnerTransactionHistoryService {
	s.params["currency_pair"] = currencyPair
	return s
}

// SetUserID narrows the result to a single recommended user.
func (s *PartnerTransactionHistoryService) SetUserID(userID int64) *PartnerTransactionHistoryService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetFrom sets the start time of the query window.
func (s *PartnerTransactionHistoryService) SetFrom(from time.Time) *PartnerTransactionHistoryService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *PartnerTransactionHistoryService) SetTo(to time.Time) *PartnerTransactionHistoryService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *PartnerTransactionHistoryService) SetLimit(limit int) *PartnerTransactionHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *PartnerTransactionHistoryService) SetOffset(offset int) *PartnerTransactionHistoryService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *PartnerTransactionHistoryService) Do(ctx context.Context) (*PartnerTransaction, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/partner/transaction_history", s.params).WithSign()
	return request.Do[PartnerTransaction](req)
}

// PartnerTransaction is a partner's recommended-user transaction history.
type PartnerTransaction struct {
	Total int64                     `json:"total"`
	List  []AgencyTransactionRecord `json:"list"`
}

// PartnerCommissionsHistoryService -- GET /api/v4/rebate/partner/commission_history (private)
//
// Partner obtains the rebate records of its recommended users. The query window
// may not exceed 30 days.
type PartnerCommissionsHistoryService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewPartnerCommissionsHistoryService() *PartnerCommissionsHistoryService {
	return &PartnerCommissionsHistoryService{c: c, params: map[string]string{}}
}

// SetCurrency narrows the result to a single currency.
func (s *PartnerCommissionsHistoryService) SetCurrency(currency string) *PartnerCommissionsHistoryService {
	s.params["currency"] = currency
	return s
}

// SetUserID narrows the result to a single recommended user.
func (s *PartnerCommissionsHistoryService) SetUserID(userID int64) *PartnerCommissionsHistoryService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetFrom sets the start time of the query window.
func (s *PartnerCommissionsHistoryService) SetFrom(from time.Time) *PartnerCommissionsHistoryService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *PartnerCommissionsHistoryService) SetTo(to time.Time) *PartnerCommissionsHistoryService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *PartnerCommissionsHistoryService) SetLimit(limit int) *PartnerCommissionsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *PartnerCommissionsHistoryService) SetOffset(offset int) *PartnerCommissionsHistoryService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *PartnerCommissionsHistoryService) Do(ctx context.Context) (*PartnerCommission, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/partner/commission_history", s.params).WithSign()
	return request.Do[PartnerCommission](req)
}

// PartnerCommission is a partner's recommended-user rebate history.
type PartnerCommission struct {
	Total int64                    `json:"total"`
	List  []AgencyCommissionRecord `json:"list"`
}

// PartnerSubListService -- GET /api/v4/rebate/partner/sub_list (private)
//
// Partner subordinate list, including sub-agents, direct customers and indirect
// customers.
type PartnerSubListService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewPartnerSubListService() *PartnerSubListService {
	return &PartnerSubListService{c: c, params: map[string]string{}}
}

// SetUserID narrows the result to a single subordinate user.
func (s *PartnerSubListService) SetUserID(userID int64) *PartnerSubListService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetLimit caps the number of records returned in a single list.
func (s *PartnerSubListService) SetLimit(limit int) *PartnerSubListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *PartnerSubListService) SetOffset(offset int) *PartnerSubListService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

func (s *PartnerSubListService) Do(ctx context.Context) (*PartnerSubItem, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/partner/sub_list", s.params).WithSign()
	return request.Do[PartnerSubItem](req)
}

// PartnerSubItem is a partner's subordinate list.
type PartnerSubItem struct {
	Total int64             `json:"total"`
	List  []PartnerSubEntry `json:"list"`
}

// PartnerSubEntry is a single subordinate of a partner.
type PartnerSubEntry struct {
	UserID       int64     `json:"user_id"`
	UserJoinTime time.Time `json:"user_join_time,format:unix"`
	// Type is 1 - sub-agent, 2 - indirect direct customer, 3 - direct direct customer.
	Type int64 `json:"type"`
}

// RebateBrokerCommissionHistoryService -- GET /api/v4/rebate/broker/commission_history (private)
//
// Broker obtains its users' rebate records. The query window may not exceed 30
// days.
type RebateBrokerCommissionHistoryService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewRebateBrokerCommissionHistoryService() *RebateBrokerCommissionHistoryService {
	return &RebateBrokerCommissionHistoryService{c: c, params: map[string]string{}}
}

// SetLimit caps the number of records returned in a single list.
func (s *RebateBrokerCommissionHistoryService) SetLimit(limit int) *RebateBrokerCommissionHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *RebateBrokerCommissionHistoryService) SetOffset(offset int) *RebateBrokerCommissionHistoryService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetUserID narrows the result to a single user.
func (s *RebateBrokerCommissionHistoryService) SetUserID(userID int64) *RebateBrokerCommissionHistoryService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetFrom sets the start time of the query window.
func (s *RebateBrokerCommissionHistoryService) SetFrom(from time.Time) *RebateBrokerCommissionHistoryService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *RebateBrokerCommissionHistoryService) SetTo(to time.Time) *RebateBrokerCommissionHistoryService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *RebateBrokerCommissionHistoryService) Do(ctx context.Context) ([]BrokerCommission, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/broker/commission_history", s.params).WithSign()
	resp, err := request.Do[[]BrokerCommission](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerCommission is a page of a broker's user rebate records.
type BrokerCommission struct {
	Total int64                    `json:"total"`
	List  []BrokerCommissionRecord `json:"list"`
}

// BrokerCommissionRecord is a single broker rebate record.
type BrokerCommissionRecord struct {
	CommissionTime time.Time           `json:"commission_time,format:unix"`
	UserID         int64               `json:"user_id"`
	GroupName      string              `json:"group_name"`
	Amount         decimal.Decimal     `json:"amount"`
	Fee            decimal.Decimal     `json:"fee"`
	FeeAsset       string              `json:"fee_asset"`
	RebateFee      decimal.Decimal     `json:"rebate_fee"`
	Source         string              `json:"source"`
	CurrencyPair   string              `json:"currency_pair"`
	SubBrokerInfo  BrokerSubBrokerInfo `json:"sub_broker_info"`
	// AlphaContractAddr is the Alpha contract address.
	AlphaContractAddr string `json:"alpha_contract_addr"`
}

// BrokerSubBrokerInfo is the sub-broker breakdown attached to a broker record.
type BrokerSubBrokerInfo struct {
	UserID                 int64           `json:"user_id"`
	OriginalCommissionRate decimal.Decimal `json:"original_commission_rate"`
	RelativeCommissionRate decimal.Decimal `json:"relative_commission_rate"`
	CommissionRate         decimal.Decimal `json:"commission_rate"`
}

// RebateBrokerTransactionHistoryService -- GET /api/v4/rebate/broker/transaction_history (private)
//
// Broker obtains its users' trading history. The query window may not exceed 30
// days.
type RebateBrokerTransactionHistoryService struct {
	c      *RebateClient
	params map[string]string
}

func (c *RebateClient) NewRebateBrokerTransactionHistoryService() *RebateBrokerTransactionHistoryService {
	return &RebateBrokerTransactionHistoryService{c: c, params: map[string]string{}}
}

// SetLimit caps the number of records returned in a single list.
func (s *RebateBrokerTransactionHistoryService) SetLimit(limit int) *RebateBrokerTransactionHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOffset sets the list offset, starting from 0.
func (s *RebateBrokerTransactionHistoryService) SetOffset(offset int) *RebateBrokerTransactionHistoryService {
	s.params["offset"] = strconv.Itoa(offset)
	return s
}

// SetUserID narrows the result to a single user.
func (s *RebateBrokerTransactionHistoryService) SetUserID(userID int64) *RebateBrokerTransactionHistoryService {
	s.params["user_id"] = strconv.FormatInt(userID, 10)
	return s
}

// SetFrom sets the start time of the query window.
func (s *RebateBrokerTransactionHistoryService) SetFrom(from time.Time) *RebateBrokerTransactionHistoryService {
	s.params["from"] = strconv.FormatInt(from.Unix(), 10)
	return s
}

// SetTo sets the end time of the query window.
func (s *RebateBrokerTransactionHistoryService) SetTo(to time.Time) *RebateBrokerTransactionHistoryService {
	s.params["to"] = strconv.FormatInt(to.Unix(), 10)
	return s
}

func (s *RebateBrokerTransactionHistoryService) Do(ctx context.Context) ([]BrokerTransaction, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/broker/transaction_history", s.params).WithSign()
	resp, err := request.Do[[]BrokerTransaction](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerTransaction is a page of a broker's user trading history.
type BrokerTransaction struct {
	Total int64                     `json:"total"`
	List  []BrokerTransactionRecord `json:"list"`
}

// BrokerTransactionRecord is a single broker trading record.
type BrokerTransactionRecord struct {
	TransactionTime time.Time           `json:"transaction_time,format:unix"`
	UserID          int64               `json:"user_id"`
	GroupName       string              `json:"group_name"`
	Fee             decimal.Decimal     `json:"fee"`
	CurrencyPair    string              `json:"currency_pair"`
	Amount          decimal.Decimal     `json:"amount"`
	FeeAsset        string              `json:"fee_asset"`
	Source          string              `json:"source"`
	SubBrokerInfo   BrokerSubBrokerInfo `json:"sub_broker_info"`
	// AlphaContractAddr is the Alpha contract address.
	AlphaContractAddr string `json:"alpha_contract_addr"`
}

// RebateUserInfoService -- GET /api/v4/rebate/user/info (private)
//
// User obtains its own rebate information.
type RebateUserInfoService struct {
	c *RebateClient
}

func (c *RebateClient) NewRebateUserInfoService() *RebateUserInfoService {
	return &RebateUserInfoService{c: c}
}

func (s *RebateUserInfoService) Do(ctx context.Context) ([]RebateUserInfo, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/user/info").WithSign()
	resp, err := request.Do[[]RebateUserInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// RebateUserInfo is a user's rebate information.
type RebateUserInfo struct {
	// InviteUID is the UID of the account that invited this user.
	InviteUID int64 `json:"invite_uid"`
}

// UserSubRelationService -- GET /api/v4/rebate/user/sub_relation (private)
//
// Queries whether the specified users are within the caller's referral system.
type UserSubRelationService struct {
	c      *RebateClient
	params map[string]string
}

// NewUserSubRelationService takes a comma-separated user ID list (only the first
// 100 IDs are honored).
func (c *RebateClient) NewUserSubRelationService(userIDList string) *UserSubRelationService {
	return &UserSubRelationService{c: c, params: map[string]string{"user_id_list": userIDList}}
}

func (s *UserSubRelationService) Do(ctx context.Context) (*UserSubRelation, error) {
	req := request.Get(ctx, s.c, "/api/v4/rebate/user/sub_relation", s.params).WithSign()
	return request.Do[UserSubRelation](req)
}

// UserSubRelation is the subordinate-relationship result for the queried users.
type UserSubRelation struct {
	List []UserSubEntry `json:"list"`
}

// UserSubEntry describes one queried user's position in the referral system.
type UserSubEntry struct {
	UID int64 `json:"uid"`
	// Belong is the system the user belongs to (partner/referral); empty means none.
	Belong string `json:"belong"`
	// Type is 0 - not in system, 1 - direct subordinate agent, 2 - indirect
	// subordinate agent, 3 - direct direct customer, 4 - indirect direct
	// customer, 5 - regular user.
	Type int64 `json:"type"`
	// RefUID is the inviter user ID.
	RefUID int64 `json:"ref_uid"`
}
