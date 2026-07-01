package account

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// GetAccountDetailService -- GET /api/v4/account/detail (private)
//
// Returns the authenticated account's profile: user id, IP / trading-pair
// whitelists, VIP tier and API-key mode.
type GetAccountDetailService struct {
	c *AccountClient
}

func (c *AccountClient) NewGetAccountDetailService() *GetAccountDetailService {
	return &GetAccountDetailService{c: c}
}

func (s *GetAccountDetailService) Do(ctx context.Context) (*AccountDetail, error) {
	req := request.Get(ctx, s.c, "/api/v4/account/detail").WithSign()
	return request.Do[AccountDetail](req)
}

// AccountDetail is the authenticated account's profile information.
type AccountDetail struct {
	IPWhitelist         []string         `json:"ip_whitelist"`
	CurrencyPairs       []string         `json:"currency_pairs"`
	UserID              int64            `json:"user_id"`
	Tier                int64            `json:"tier"`
	VIPTier             int64            `json:"vip_tier"`
	Key                 AccountDetailKey `json:"key"`
	CopyTradingRole     int              `json:"copy_trading_role"`
	SpotCopyTradingRole int              `json:"spot_copy_trading_role"`
	TierExpireTime      time.Time        `json:"tier_expire_time"` // RFC3339
}

// AccountDetailKey describes the API key used for the request.
type AccountDetailKey struct {
	// Mode: 1 - Classic mode, 2 - Legacy unified mode.
	Mode int `json:"mode"`
}

// GetAccountRateLimitService -- GET /api/v4/account/rate_limit (private)
//
// Returns the account's fill-ratio based transaction rate-limit tiers.
type GetAccountRateLimitService struct {
	c *AccountClient
}

func (c *AccountClient) NewGetAccountRateLimitService() *GetAccountRateLimitService {
	return &GetAccountRateLimitService{c: c}
}

func (s *GetAccountRateLimitService) Do(ctx context.Context) ([]AccountRateLimit, error) {
	req := request.Get(ctx, s.c, "/api/v4/account/rate_limit").WithSign()
	resp, err := request.Do[[]AccountRateLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountRateLimit is a fill-ratio rate-limit tier for the account.
type AccountRateLimit struct {
	Tier      string          `json:"tier"`
	Type      string          `json:"type"`
	Ratio     decimal.Decimal `json:"ratio"`
	MainRatio decimal.Decimal `json:"main_ratio"`
	UpdatedAt time.Time       `json:"updated_at,string,format:unix"`
}

// ListSTPGroupsService -- GET /api/v4/account/stp_groups (private)
//
// Lists the self-trade-prevention (STP) user groups created by the current main
// account.
type ListSTPGroupsService struct {
	c      *AccountClient
	params map[string]string
}

func (c *AccountClient) NewListSTPGroupsService() *ListSTPGroupsService {
	return &ListSTPGroupsService{c: c, params: map[string]string{}}
}

// SetName fuzzy-searches STP groups by name.
func (s *ListSTPGroupsService) SetName(name string) *ListSTPGroupsService {
	s.params["name"] = name
	return s
}

func (s *ListSTPGroupsService) Do(ctx context.Context) ([]STPGroup, error) {
	req := request.Get(ctx, s.c, "/api/v4/account/stp_groups", s.params).WithSign()
	resp, err := request.Do[[]STPGroup](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// STPGroup is a self-trade-prevention user group.
type STPGroup struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatorID  int64     `json:"creator_id"`
	CreateTime time.Time `json:"create_time,format:unix"`
}

// CreateSTPGroupService -- POST /api/v4/account/stp_groups (private)
//
// Creates a new STP user group (main account only).
type CreateSTPGroupService struct {
	c    *AccountClient
	body map[string]any
}

func (c *AccountClient) NewCreateSTPGroupService(name string) *CreateSTPGroupService {
	return &CreateSTPGroupService{c: c, body: map[string]any{"name": name}}
}

func (s *CreateSTPGroupService) Do(ctx context.Context) (*STPGroup, error) {
	req := request.Post(ctx, s.c, "/api/v4/account/stp_groups", s.body).WithSign()
	return request.Do[STPGroup](req)
}

// ListSTPGroupsUsersService -- GET /api/v4/account/stp_groups/{stp_id}/users (private)
//
// Lists the users belonging to an STP group (creator main account only).
type ListSTPGroupsUsersService struct {
	c     *AccountClient
	stpID int64
}

func (c *AccountClient) NewListSTPGroupsUsersService(stpID int64) *ListSTPGroupsUsersService {
	return &ListSTPGroupsUsersService{c: c, stpID: stpID}
}

func (s *ListSTPGroupsUsersService) Do(ctx context.Context) ([]STPGroupUser, error) {
	req := request.Get(ctx, s.c, "/api/v4/account/stp_groups/"+strconv.FormatInt(s.stpID, 10)+"/users").WithSign()
	resp, err := request.Do[[]STPGroupUser](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// STPGroupUser is a member of an STP group.
type STPGroupUser struct {
	UserID     int64     `json:"user_id"`
	STPID      int64     `json:"stp_id"`
	CreateTime time.Time `json:"create_time,format:unix"`
}

// AddSTPGroupUsersService -- POST /api/v4/account/stp_groups/{stp_id}/users (private)
//
// Adds users to an STP group. Only sub-accounts under the current main account
// are allowed.
type AddSTPGroupUsersService struct {
	c       *AccountClient
	stpID   int64
	userIDs []int64
}

func (c *AccountClient) NewAddSTPGroupUsersService(stpID int64, userIDs []int64) *AddSTPGroupUsersService {
	return &AddSTPGroupUsersService{c: c, stpID: stpID, userIDs: userIDs}
}

func (s *AddSTPGroupUsersService) Do(ctx context.Context) ([]STPGroupUser, error) {
	req := request.Post(ctx, s.c, "/api/v4/account/stp_groups/"+strconv.FormatInt(s.stpID, 10)+"/users").WithSign().SetBody(s.userIDs)
	resp, err := request.Do[[]STPGroupUser](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DeleteSTPGroupUsersService -- DELETE /api/v4/account/stp_groups/{stp_id}/users (private)
//
// Removes a user from an STP group (creator main account only).
type DeleteSTPGroupUsersService struct {
	c      *AccountClient
	stpID  int64
	params map[string]string
}

func (c *AccountClient) NewDeleteSTPGroupUsersService(stpID, userID int64) *DeleteSTPGroupUsersService {
	return &DeleteSTPGroupUsersService{c: c, stpID: stpID, params: map[string]string{
		"user_id": strconv.FormatInt(userID, 10),
	}}
}

func (s *DeleteSTPGroupUsersService) Do(ctx context.Context) ([]STPGroupUser, error) {
	req := request.Delete(ctx, s.c, "/api/v4/account/stp_groups/"+strconv.FormatInt(s.stpID, 10)+"/users", s.params).WithSign()
	resp, err := request.Do[[]STPGroupUser](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetDebitFeeService -- GET /api/v4/account/debit_fee (private)
//
// Returns the account's GT fee-deduction configuration.
type GetDebitFeeService struct {
	c *AccountClient
}

func (c *AccountClient) NewGetDebitFeeService() *GetDebitFeeService {
	return &GetDebitFeeService{c: c}
}

func (s *GetDebitFeeService) Do(ctx context.Context) (*DebitFee, error) {
	req := request.Get(ctx, s.c, "/api/v4/account/debit_fee").WithSign()
	return request.Do[DebitFee](req)
}

// DebitFee is the account's GT fee-deduction configuration.
type DebitFee struct {
	Enabled bool `json:"enabled"`
	// DebitFee is the GT fee-deduction switch (1 - enabled, 0 - disabled).
	DebitFee int `json:"debit_fee"`
}

// SetDebitFeeService -- POST /api/v4/account/debit_fee (private)
//
// Configures the account's GT fee deduction. Returns no content.
type SetDebitFeeService struct {
	c    *AccountClient
	body map[string]any
}

func (c *AccountClient) NewSetDebitFeeService(debitFee int) *SetDebitFeeService {
	return &SetDebitFeeService{c: c, body: map[string]any{"debit_fee": debitFee}}
}

func (s *SetDebitFeeService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/account/debit_fee", s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}
