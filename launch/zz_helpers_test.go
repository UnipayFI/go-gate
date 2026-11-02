package launch

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

// testPublicClient builds an unauthenticated client for public endpoints.
func testPublicClient() *LaunchClient {
	opts := []client.Options{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewLaunchClient(opts...)
}

// testClient builds an authenticated client, skipping when creds are unset.
func testClient(t *testing.T) *LaunchClient {
	t.Helper()
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.Options{client.WithAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewLaunchClient(opts...)
}
