package delivery

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/UnipayFI/go-gate/v4/request"
)

// testWSClient builds an unauthenticated delivery stream client for public tests.
func testWSClient() *DeliveryWebSocketClient {
	opts := []client.WebSocketOptions{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewDeliveryWebSocketClient(opts...)
}

// liveContract returns a currently-listed delivery contract name, skipping when
// none are listed (delivery contracts have expiry-suffixed names that rotate).
func liveContract(t *testing.T) string {
	t.Helper()
	c := testPublicClient()
	cx := testutil.Ctx(t)
	list, err := c.NewListDeliveryContractsService(SettleUSDT).Do(cx)
	if err != nil {
		t.Fatalf("list delivery contracts: %v", err)
	}
	if len(list) == 0 {
		t.Skip("no delivery contracts listed; skipping delivery WS test")
	}
	return list[0].Name
}

func TestDeliveryWSPublic(t *testing.T) {
	contract := liveContract(t)
	c := testWSClient()
	ctx := testutil.Ctx(t)

	// Delivery is low-volume: a book_ticker push may not arrive quickly, so a
	// quiet 15s window is logged (not failed) — the subscribe + decode path is
	// still exercised.
	got := make(chan *request.WsPush[WsDeliveryBookTicker], 1)
	done, _, err := c.NewSubscribeBookTickerService(contract).Do(ctx, func(p *request.WsPush[WsDeliveryBookTicker], e error) {
		if e != nil {
			return
		}
		select {
		case got <- p:
		default:
		}
	})
	if err != nil {
		t.Fatalf("subscribe book_ticker %s: %v", contract, err)
	}
	defer close(done)

	select {
	case p := <-got:
		t.Logf("book_ticker: %s bid=%s ask=%s ts=%s", p.Result.Contract, p.Result.BestBidPrice, p.Result.BestAskPrice, p.Result.Time)
		if p.Result.Contract == "" {
			t.Errorf("empty contract in book_ticker push: %+v", p.Result)
		}
	case <-time.After(15 * time.Second):
		t.Logf("no book_ticker push for %s in 15s (delivery is low-volume); channel+decode path exercised", contract)
	}
}

// TestDeliveryWSPrivate subscribes to the private positions channel and fails only
// on an auth/error push; a quiet window is a pass since delivery accounts rarely
// hold live positions.
func TestDeliveryWSPrivate(t *testing.T) {
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.WebSocketOptions{client.WithWebSocketAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	c := NewDeliveryWebSocketClient(opts...)
	ctx := testutil.Ctx(t)

	// Resolve the account's user id via REST; skip when the account has no
	// delivery wallet.
	rc := testClient(t)
	if err := rc.SyncServerTime(ctx); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	acc, err := rc.NewListDeliveryAccountsService(SettleUSDT).Do(ctx)
	if err != nil {
		if testutil.Tolerable(t, "delivery/accounts", err) {
			t.Skip("no delivery account for these creds; skipping private WS test")
		}
		t.Fatalf("delivery accounts: %v", err)
	}
	t.Logf("delivery account: currency=%s total=%s", acc.Currency, acc.Total)

	// The accounts struct does not surface the user id, so read it from the raw
	// accounts body (the delivery/futures WS private payload wants it as a string).
	raw := testutil.FetchRawGet(t, rc, ctx, "/api/v4/delivery/usdt/accounts", nil, true)
	var meta struct {
		User int64 `json:"user"`
	}
	if uerr := common.JSONUnmarshal(raw, &meta); uerr != nil {
		t.Fatalf("parse user id: %v", uerr)
	}
	userID := strconv.FormatInt(meta.User, 10)

	contract := liveContract(t)

	pushErr := make(chan error, 4)
	done, _, err := c.NewSubscribePositionsService(userID, contract).Do(ctx, func(p *request.WsPush[[]WsDeliveryPosition], e error) {
		if e != nil {
			pushErr <- e
			return
		}
		for _, pos := range p.Result {
			t.Logf("ws position push: contract=%s size=%d entry=%s", pos.Contract, pos.Size, pos.EntryPrice)
		}
	})
	if err != nil {
		t.Fatalf("subscribe positions: %v", err)
	}
	defer close(done)

	// A bad auth surfaces as an error push; otherwise a quiet window is a pass.
	select {
	case e := <-pushErr:
		t.Fatalf("positions subscription error (auth?): %v", e)
	case <-time.After(2 * time.Second):
		t.Log("positions subscription accepted (no error push); channel+auth path exercised")
	}
}
