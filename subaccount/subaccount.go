package subaccount

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
)

// ListSubAccountsService -- GET /api/v4/sub_accounts (private)
//
// Lists the main account's sub-accounts, optionally filtered by type.
type ListSubAccountsService struct {
	c      *SubAccountClient
	params map[string]string
}

func (c *SubAccountClient) NewListSubAccountsService() *ListSubAccountsService {
	return &ListSubAccountsService{c: c, params: map[string]string{}}
}

// SetType filters by sub-account type ("0" lists all supported types,
// "1" lists regular sub-accounts only; default is regular sub-accounts).
func (s *ListSubAccountsService) SetType(t string) *ListSubAccountsService {
	s.params["type"] = t
	return s
}

func (s *ListSubAccountsService) Do(ctx context.Context) ([]SubAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/sub_accounts", s.params).WithSign()
	resp, err := request.Do[[]SubAccount](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateSubAccountsService -- POST /api/v4/sub_accounts (private)
//
// Creates a new sub-account under the main account.
type CreateSubAccountsService struct {
	c    *SubAccountClient
	body map[string]any
}

func (c *SubAccountClient) NewCreateSubAccountsService(loginName string) *CreateSubAccountsService {
	return &CreateSubAccountsService{c: c, body: map[string]any{"login_name": loginName}}
}

// SetRemark sets the sub-account remark.
func (s *CreateSubAccountsService) SetRemark(remark string) *CreateSubAccountsService {
	s.body["remark"] = remark
	return s
}

// SetPassword sets the sub-account login password (defaults to the main account's).
func (s *CreateSubAccountsService) SetPassword(password string) *CreateSubAccountsService {
	s.body["password"] = password
	return s
}

// SetEmail sets the sub-account email address (defaults to the main account's).
func (s *CreateSubAccountsService) SetEmail(email string) *CreateSubAccountsService {
	s.body["email"] = email
	return s
}

func (s *CreateSubAccountsService) Do(ctx context.Context) (*SubAccount, error) {
	req := request.Post(ctx, s.c, "/api/v4/sub_accounts", s.body).WithSign()
	return request.Do[SubAccount](req)
}

// GetSubAccountService -- GET /api/v4/sub_accounts/{user_id} (private)
//
// Returns a single sub-account by user ID.
type GetSubAccountService struct {
	c      *SubAccountClient
	userID int64
}

func (c *SubAccountClient) NewGetSubAccountService(userID int64) *GetSubAccountService {
	return &GetSubAccountService{c: c, userID: userID}
}

func (s *GetSubAccountService) Do(ctx context.Context) (*SubAccount, error) {
	req := request.Get(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)).WithSign()
	return request.Do[SubAccount](req)
}

// ListSubAccountKeysService -- GET /api/v4/sub_accounts/{user_id}/keys (private)
//
// Lists every API key pair belonging to a sub-account.
type ListSubAccountKeysService struct {
	c      *SubAccountClient
	userID int64
}

func (c *SubAccountClient) NewListSubAccountKeysService(userID int64) *ListSubAccountKeysService {
	return &ListSubAccountKeysService{c: c, userID: userID}
}

func (s *ListSubAccountKeysService) Do(ctx context.Context) ([]SubAccountKey, error) {
	req := request.Get(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/keys").WithSign()
	resp, err := request.Do[[]SubAccountKey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CreateSubAccountKeysService -- POST /api/v4/sub_accounts/{user_id}/keys (private)
//
// Creates a new API key pair for a sub-account.
type CreateSubAccountKeysService struct {
	c      *SubAccountClient
	userID int64
	body   map[string]any
}

func (c *SubAccountClient) NewCreateSubAccountKeysService(userID int64) *CreateSubAccountKeysService {
	return &CreateSubAccountKeysService{c: c, userID: userID, body: map[string]any{}}
}

// SetName sets the API key name.
func (s *CreateSubAccountKeysService) SetName(name string) *CreateSubAccountKeysService {
	s.body["name"] = name
	return s
}

// SetPerms sets the permission scopes granted to the key.
func (s *CreateSubAccountKeysService) SetPerms(perms []SubAccountKeyPerm) *CreateSubAccountKeysService {
	s.body["perms"] = perms
	return s
}

// SetIPWhitelist restricts the key to the given source IPs.
func (s *CreateSubAccountKeysService) SetIPWhitelist(ips []string) *CreateSubAccountKeysService {
	s.body["ip_whitelist"] = ips
	return s
}

func (s *CreateSubAccountKeysService) Do(ctx context.Context) (*SubAccountKey, error) {
	req := request.Post(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/keys", s.body).WithSign()
	return request.Do[SubAccountKey](req)
}

// GetSubAccountKeyService -- GET /api/v4/sub_accounts/{user_id}/keys/{key} (private)
//
// Returns a single API key pair of a sub-account.
type GetSubAccountKeyService struct {
	c      *SubAccountClient
	userID int64
	key    string
}

func (c *SubAccountClient) NewGetSubAccountKeyService(userID int64, key string) *GetSubAccountKeyService {
	return &GetSubAccountKeyService{c: c, userID: userID, key: key}
}

func (s *GetSubAccountKeyService) Do(ctx context.Context) (*SubAccountKey, error) {
	req := request.Get(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/keys/"+s.key).WithSign()
	return request.Do[SubAccountKey](req)
}

// UpdateSubAccountKeysService -- PUT /api/v4/sub_accounts/{user_id}/keys/{key} (private)
//
// Updates a sub-account API key pair's permissions and whitelist. Returns no content.
type UpdateSubAccountKeysService struct {
	c      *SubAccountClient
	userID int64
	key    string
	body   map[string]any
}

func (c *SubAccountClient) NewUpdateSubAccountKeysService(userID int64, key string) *UpdateSubAccountKeysService {
	return &UpdateSubAccountKeysService{c: c, userID: userID, key: key, body: map[string]any{}}
}

// SetName renames the API key.
func (s *UpdateSubAccountKeysService) SetName(name string) *UpdateSubAccountKeysService {
	s.body["name"] = name
	return s
}

// SetPerms replaces the permission scopes granted to the key.
func (s *UpdateSubAccountKeysService) SetPerms(perms []SubAccountKeyPerm) *UpdateSubAccountKeysService {
	s.body["perms"] = perms
	return s
}

// SetIPWhitelist replaces the key's source-IP whitelist.
func (s *UpdateSubAccountKeysService) SetIPWhitelist(ips []string) *UpdateSubAccountKeysService {
	s.body["ip_whitelist"] = ips
	return s
}

func (s *UpdateSubAccountKeysService) Do(ctx context.Context) error {
	req := request.Put(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/keys/"+s.key, s.body).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// DeleteSubAccountKeysService -- DELETE /api/v4/sub_accounts/{user_id}/keys/{key} (private)
//
// Deletes a sub-account API key pair. Returns no content.
type DeleteSubAccountKeysService struct {
	c      *SubAccountClient
	userID int64
	key    string
}

func (c *SubAccountClient) NewDeleteSubAccountKeysService(userID int64, key string) *DeleteSubAccountKeysService {
	return &DeleteSubAccountKeysService{c: c, userID: userID, key: key}
}

func (s *DeleteSubAccountKeysService) Do(ctx context.Context) error {
	req := request.Delete(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/keys/"+s.key).WithSign()
	_, err := request.DoRaw(req)
	return err
}

// LockSubAccountService -- POST /api/v4/sub_accounts/{user_id}/lock (private)
//
// Locks a sub-account, suspending its access. Returns no content.
type LockSubAccountService struct {
	c      *SubAccountClient
	userID int64
}

func (c *SubAccountClient) NewLockSubAccountService(userID int64) *LockSubAccountService {
	return &LockSubAccountService{c: c, userID: userID}
}

func (s *LockSubAccountService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/lock").WithSign()
	_, err := request.DoRaw(req)
	return err
}

// UnlockSubAccountService -- POST /api/v4/sub_accounts/{user_id}/unlock (private)
//
// Unlocks a previously locked sub-account. Returns no content.
type UnlockSubAccountService struct {
	c      *SubAccountClient
	userID int64
}

func (c *SubAccountClient) NewUnlockSubAccountService(userID int64) *UnlockSubAccountService {
	return &UnlockSubAccountService{c: c, userID: userID}
}

func (s *UnlockSubAccountService) Do(ctx context.Context) error {
	req := request.Post(ctx, s.c, "/api/v4/sub_accounts/"+strconv.FormatInt(s.userID, 10)+"/unlock").WithSign()
	_, err := request.DoRaw(req)
	return err
}

// ListUnifiedModeService -- GET /api/v4/sub_accounts/unified_mode (private)
//
// Returns each sub-account's unified-account mode.
type ListUnifiedModeService struct {
	c *SubAccountClient
}

func (c *SubAccountClient) NewListUnifiedModeService() *ListUnifiedModeService {
	return &ListUnifiedModeService{c: c}
}

func (s *ListUnifiedModeService) Do(ctx context.Context) ([]SubAccountUnifiedMode, error) {
	req := request.Get(ctx, s.c, "/api/v4/sub_accounts/unified_mode").WithSign()
	resp, err := request.Do[[]SubAccountUnifiedMode](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubAccount is a sub-account belonging to the main account.
type SubAccount struct {
	UserID    int64  `json:"user_id"`
	LoginName string `json:"login_name"`
	Remark    string `json:"remark"`
	Email     string `json:"email"`
	// State is the sub-account status: 1-normal, 2-locked.
	State int `json:"state"`
	// Type is the sub-account type: 1-regular, 3-cross-margin.
	Type       int       `json:"type"`
	CreateTime time.Time `json:"create_time,format:unix"`
}

// SubAccountKey is an API key pair belonging to a sub-account.
type SubAccountKey struct {
	UserID int64 `json:"user_id"`
	// Mode is the account mode: 1-classic, 2-portfolio.
	Mode        int                 `json:"mode"`
	Name        string              `json:"name"`
	Perms       []SubAccountKeyPerm `json:"perms"`
	IPWhitelist []string            `json:"ip_whitelist"`
	Key         string              `json:"key"`
	// State is the key status: 1-normal, 2-frozen, 3-locked.
	State      int       `json:"state"`
	CreatedAt  time.Time `json:"created_at,format:unix"`
	UpdatedAt  time.Time `json:"updated_at,format:unix"`
	LastAccess time.Time `json:"last_access,format:unix"`
}

// SubAccountKeyPerm is a single permission scope granted to a sub-account key.
type SubAccountKeyPerm struct {
	Name     string `json:"name"`
	ReadOnly bool   `json:"read_only"`
}

// SubAccountUnifiedMode is a sub-account's unified-account mode.
type SubAccountUnifiedMode struct {
	UserID    int64 `json:"user_id"`
	IsUnified bool  `json:"is_unified"`
	// Mode is the unified account mode: classic, multi_currency, or portfolio.
	Mode string `json:"mode"`
}
