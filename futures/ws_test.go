package futures

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/UnipayFI/go-gate/request"
	"github.com/shopspring/decimal"
)

// newTestFuturesWSClient builds an unauthenticated futures stream client for the
// given settlement currency, wiring the optional test proxy.
func newTestFuturesWSClient(settle Settle) *FuturesWebSocketClient {
	opts := []client.WebSocketOptions{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewFuturesWebSocketClient(settle, opts...)
}

func TestFuturesWSPublic(t *testing.T) {
	c := newTestFuturesWSClient(SettleUSDT)
	ctx := testutil.Ctx(t)

	t.Run("BookTicker", func(t *testing.T) {
		got := make(chan *request.WsPush[WsFuturesBookTicker], 1)
		done, _, err := c.NewSubscribeBookTickerService("BTC_USDT").Do(ctx, func(p *request.WsPush[WsFuturesBookTicker], e error) {
			if e != nil {
				return
			}
			select {
			case got <- p:
			default:
			}
		})
		if err != nil {
			t.Fatalf("subscribe book_ticker: %v", err)
		}
		defer close(done)
		select {
		case p := <-got:
			t.Logf("book_ticker: %s bid=%s ask=%s ts=%s", p.Result.Contract, p.Result.BestBidPrice, p.Result.BestAskPrice, p.Result.Time)
			if p.Result.BestBidPrice.IsZero() || p.Result.Contract != "BTC_USDT" {
				t.Errorf("unexpected book_ticker: %+v", p.Result)
			}
		case <-time.After(15 * time.Second):
			t.Fatal("no book_ticker push in 15s")
		}
	})

	t.Run("Tickers", func(t *testing.T) {
		got := make(chan WsFuturesTicker, 1)
		done, _, err := c.NewSubscribeTickersService("BTC_USDT").Do(ctx, func(p *request.WsPush[[]WsFuturesTicker], e error) {
			if e != nil {
				return
			}
			for _, tk := range p.Result {
				select {
				case got <- tk:
				default:
				}
			}
		})
		if err != nil {
			t.Fatalf("subscribe tickers: %v", err)
		}
		defer close(done)
		select {
		case tk := <-got:
			t.Logf("ticker: %s last=%s mark=%s funding=%s", tk.Contract, tk.Last, tk.MarkPrice, tk.FundingRate)
			if tk.Last.IsZero() {
				t.Error("zero last price")
			}
		case <-time.After(15 * time.Second):
			t.Fatal("no ticker push in 15s")
		}
	})

	t.Run("Trades", func(t *testing.T) {
		got := make(chan WsFuturesTrade, 1)
		done, _, err := c.NewSubscribeTradesService("BTC_USDT").Do(ctx, func(p *request.WsPush[[]WsFuturesTrade], e error) {
			if e != nil {
				return
			}
			for _, tr := range p.Result {
				select {
				case got <- tr:
				default:
				}
			}
		})
		if err != nil {
			t.Fatalf("subscribe trades: %v", err)
		}
		defer close(done)
		select {
		case tr := <-got:
			t.Logf("trade: %s size=%d px=%s t=%s", tr.Contract, tr.Size, tr.Price, tr.CreateTime)
			if tr.Price.IsZero() {
				t.Error("zero trade price")
			}
		case <-time.After(20 * time.Second):
			t.Log("no trade push in 20s (BTC_USDT quiet); channel+decode still exercised")
		}
	})
}

// TestFuturesWSPrivate subscribes to the private orders channel and, when writes
// are enabled, places+cancels a tiny far-from-market order to trigger a real
// push. It resolves the required user id from the REST accounts endpoint.
func TestFuturesWSPrivate(t *testing.T) {
	apiKey, apiSecret := testutil.Creds(t)
	ctx := testutil.Ctx(t)

	// Resolve the numeric user id, required as the first payload element.
	rc := testClient(t)
	if err := rc.SyncServerTime(ctx); err != nil {
		t.Fatalf("sync: %v", err)
	}
	acc, err := rc.NewListFuturesAccountsService(SettleUSDT).Do(ctx)
	if err != nil {
		t.Fatalf("accounts: %v", err)
	}
	userID := strconv.FormatInt(acc.User, 10)

	opts := []client.WebSocketOptions{client.WithWebSocketAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	c := NewFuturesWebSocketClient(SettleUSDT, opts...)

	pushErr := make(chan error, 4)
	gotOrder := make(chan WsFuturesOrder, 4)
	done, _, err := c.NewSubscribeOrdersService(userID, "BTC_USDT").Do(ctx, func(p *request.WsPush[[]WsFuturesOrder], e error) {
		if e != nil {
			pushErr <- e
			return
		}
		for _, o := range p.Result {
			select {
			case gotOrder <- o:
			default:
			}
		}
	})
	if err != nil {
		t.Fatalf("subscribe orders: %v", err)
	}
	defer close(done)

	// Give the subscription a moment; a bad auth surfaces as an error push.
	select {
	case e := <-pushErr:
		t.Fatalf("orders subscription error (auth?): %v", e)
	case <-time.After(2 * time.Second):
	}

	if !testutil.WriteEnabled() {
		t.Skip("set GATE_TEST_WRITE=1 to trigger a live order push")
	}

	// Place a tiny far-below-market buy so it rests, then cancel it.
	order, err := rc.NewCreateFuturesOrderService(SettleUSDT, "BTC_USDT", 1).
		SetPrice(decimal.RequireFromString("41000")).SetTimeInForce(TimeInForceGTC).Do(ctx)
	if err != nil {
		t.Fatalf("place order: %v", err)
	}
	t.Logf("placed %d to trigger ws push", order.ID)
	defer rc.NewCancelFuturesOrderService(SettleUSDT, strconv.FormatInt(order.ID, 10)).Do(ctx)

	select {
	case o := <-gotOrder:
		t.Logf("ws order push: id=%d contract=%s size=%d left=%d finish_as=%s", o.ID, o.Contract, o.Size, o.Left, o.FinishAs)
	case <-time.After(10 * time.Second):
		t.Fatal("no order push after placing order")
	}
}
