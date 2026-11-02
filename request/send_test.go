package request

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UnipayFI/go-gate/v4/client"
)

// TestDoResetContentIsAPIError checks that an HTTP 205 business error (e.g. GET
// /futures/{settle}/price_orders/{order_id} for an unknown order) is returned as
// a *client.APIError instead of being decoded as an empty payload.
func TestDoResetContentIsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusResetContent)
		_, _ = w.Write([]byte(`{"code":"1034","label":"AUTO_ORDER_NOT_FOUND","message":"no orderID match"}`))
	}))
	defer srv.Close()

	c := client.NewClient(client.WithBaseURL(srv.URL))
	type payload struct {
		ID int64 `json:"id"`
	}
	got, err := Do[payload](Get(context.Background(), c, "/api/v4/futures/usdt/price_orders/1"))
	if err == nil {
		t.Fatalf("expected an error, got payload %+v", got)
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *client.APIError, got %T: %v", err, err)
	}
	if apiErr.Status != http.StatusResetContent || apiErr.Label != "AUTO_ORDER_NOT_FOUND" || apiErr.Message != "no orderID match" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
}
