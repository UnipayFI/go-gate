package delivery

import (
	"testing"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/internal/testutil"
)

// testPublicClient builds an unauthenticated delivery client for public endpoints.
func testPublicClient() *DeliveryClient {
	opts := []client.Options{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewDeliveryClient(opts...)
}

// testClient builds an authenticated delivery client, skipping when creds are unset.
func testClient(t *testing.T) *DeliveryClient {
	t.Helper()
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.Options{client.WithAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewDeliveryClient(opts...)
}
