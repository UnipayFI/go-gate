package futures

import (
	"testing"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/internal/testutil"
)

// testPublicClient builds an unauthenticated futures client for public endpoints.
func testPublicClient() *FuturesClient {
	opts := []client.Options{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewFuturesClient(opts...)
}

// testClient builds an authenticated futures client, skipping when creds are unset.
func testClient(t *testing.T) *FuturesClient {
	t.Helper()
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.Options{client.WithAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewFuturesClient(opts...)
}
