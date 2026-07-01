package account

import (
	"testing"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/internal/testutil"
)

// testPublicClient builds an unauthenticated client for public endpoints.
func testPublicClient() *AccountClient {
	opts := []client.Options{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewAccountClient(opts...)
}

// testClient builds an authenticated client, skipping when creds are unset.
func testClient(t *testing.T) *AccountClient {
	t.Helper()
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.Options{client.WithAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewAccountClient(opts...)
}
