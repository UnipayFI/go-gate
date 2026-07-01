// Command graw signs and executes a single Gate APIv4 REST call and pretty
// prints the raw response. It is a development aid for capturing the exact
// shape of private endpoints (which cannot be curled without HMAC-SHA512
// signing) so the typed response structs can be reconciled against reality.
//
// Usage:
//
//	GATE_API_KEY=... GATE_API_SECRET=... \
//	  go run ./cmd/graw GET  /api/v4/spot/accounts
//	  go run ./cmd/graw GET  /api/v4/spot/orders "currency_pair=BTC_USDT&status=open"
//	  go run ./cmd/graw POST /api/v4/spot/orders '{"currency_pair":"BTC_USDT", ...}'
//
// The third argument is the query string (GET/DELETE) or JSON body (POST/PUT).
// Set GATE_PROXY to route through a proxy.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	gateCommon "github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/request"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: graw <GET|POST|PUT|DELETE> <path> [query-or-jsonbody]")
		os.Exit(2)
	}
	method := strings.ToUpper(os.Args[1])
	path := os.Args[2]
	arg := ""
	if len(os.Args) > 3 {
		arg = os.Args[3]
	}

	opts := []client.Options{
		client.WithAuth(os.Getenv("GATE_API_KEY"), os.Getenv("GATE_API_SECRET")),
	}
	if proxy := os.Getenv("GATE_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	c := client.NewClient(opts...)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var req *request.Request
	switch method {
	case http.MethodGet:
		req = request.Get(ctx, c, path, parseQuery(arg)).WithSign()
	case http.MethodDelete:
		req = request.Delete(ctx, c, path, parseQuery(arg)).WithSign()
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		builder := request.Post
		if method == http.MethodPut {
			builder = request.Put
		} else if method == http.MethodPatch {
			builder = request.Patch
		}
		body := map[string]any{}
		if arg != "" {
			if err := gateCommon.JSONUnmarshal([]byte(arg), &body); err != nil {
				fail("invalid json body: %v", err)
			}
		}
		req = builder(ctx, c, path, body).WithSign()
	default:
		fail("unsupported method %q", method)
	}

	body, err := request.DoRaw(req)
	if err != nil {
		fail("request error: %v", err)
	}
	fmt.Println(pretty(body))
}

func parseQuery(q string) map[string]string {
	out := map[string]string{}
	q = strings.TrimPrefix(q, "?")
	for pair := range strings.SplitSeq(q, "&") {
		if pair == "" {
			continue
		}
		k, v, _ := strings.Cut(pair, "=")
		out[k] = v
	}
	return out
}

func pretty(b []byte) string {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(out)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
