package account

import (
	"context"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
)

// GetMainKeysService -- GET /api/v4/account/main_keys (private)
//
// Returns the API-key information of all main-account API keys: status, user
// mode, remark, whitelists, permissions and timestamps.
type GetMainKeysService struct {
	c *AccountClient
}

func (c *AccountClient) NewGetMainKeysService() *GetMainKeysService {
	return &GetMainKeysService{c: c}
}

func (s *GetMainKeysService) Do(ctx context.Context) ([]AccountMainKey, error) {
	req := request.Get(ctx, s.c, "/api/v4/account/main_keys").WithSign()
	resp, err := request.Do[[]AccountMainKey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountMainKey is a single main-account API key and its configuration.
type AccountMainKey struct {
	// State is the API key status: 1 - Normal, 2 - Locked, 3 - Frozen.
	State int `json:"state"`
	// Mode is the user mode: 1 - Classic mode, 2 - Legacy unified mode.
	Mode          int              `json:"mode"`
	Name          string           `json:"name"`
	CurrencyPairs []string         `json:"currency_pairs"`
	UserID        int64            `json:"user_id"`
	IPWhitelist   []string         `json:"ip_whitelist"`
	Perms         []AccountKeyPerm `json:"perms"`
	Key           string           `json:"key"`
	CreatedAt     time.Time        `json:"created_at,format:unix"`
	UpdateAt      time.Time        `json:"update_at,format:unix"`
	LastAccess    time.Time        `json:"last_access,format:unix"`
}

// AccountKeyPerm is one permission entry of a main-account API key.
type AccountKeyPerm struct {
	// Name is the permission function name (empty value clears it), e.g.
	// wallet, spot, futures, delivery, earn, options, account, unified, loan.
	Name     string `json:"name"`
	ReadOnly bool   `json:"read_only"`
}
