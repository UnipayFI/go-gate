package margin

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

// testPublicClient builds an unauthenticated client for public endpoints.
func testPublicClient() *MarginClient {
	opts := []client.Options{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewMarginClient(opts...)
}

// testClient builds an authenticated client, skipping when creds are unset.
func testClient(t *testing.T) *MarginClient {
	t.Helper()
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.Options{client.WithAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewMarginClient(opts...)
}
