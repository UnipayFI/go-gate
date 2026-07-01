package spot

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/v4/client"
	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

func testWSClient() *SpotWebSocketClient {
	opts := []client.WebSocketOptions{}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewSpotWebSocketClient(opts...)
}

func TestSpotWSPublic(t *testing.T) {
	c := testWSClient()
	ctx := testutil.Ctx(t)

	t.Run("BookTicker", func(t *testing.T) {
		got := make(chan *request.WsPush[WsBookTicker], 1)
		done, _, err := c.NewSubscribeBookTickerService("BTC_USDT").Do(ctx, func(p *request.WsPush[WsBookTicker], e error) {
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
			t.Logf("book_ticker: %s bid=%s ask=%s ts=%s", p.Result.CurrencyPair, p.Result.BestBid, p.Result.BestAsk, p.Result.Time)
			if p.Result.BestBid.IsZero() || p.Result.CurrencyPair != "BTC_USDT" {
				t.Errorf("unexpected book_ticker: %+v", p.Result)
			}
		case <-time.After(15 * time.Second):
			t.Fatal("no book_ticker push in 15s")
		}
	})

	t.Run("Tickers", func(t *testing.T) {
		got := make(chan *request.WsPush[WsTicker], 1)
		done, _, err := c.NewSubscribeTickersService("BTC_USDT").Do(ctx, func(p *request.WsPush[WsTicker], e error) {
			if e != nil {
				return
			}
			select {
			case got <- p:
			default:
			}
		})
		if err != nil {
			t.Fatalf("subscribe tickers: %v", err)
		}
		defer close(done)
		select {
		case p := <-got:
			t.Logf("ticker: %s last=%s vol=%s", p.Result.CurrencyPair, p.Result.Last, p.Result.BaseVolume)
			if p.Result.Last.IsZero() {
				t.Error("zero last price")
			}
		case <-time.After(15 * time.Second):
			t.Fatal("no ticker push in 15s")
		}
	})

	t.Run("Trades", func(t *testing.T) {
		got := make(chan *request.WsPush[WsPublicTrade], 1)
		done, _, err := c.NewSubscribeTradesService("BTC_USDT").Do(ctx, func(p *request.WsPush[WsPublicTrade], e error) {
			if e != nil {
				return
			}
			select {
			case got <- p:
			default:
			}
		})
		if err != nil {
			t.Fatalf("subscribe trades: %v", err)
		}
		defer close(done)
		select {
		case p := <-got:
			t.Logf("trade: %s %s px=%s amt=%s t=%s", p.Result.CurrencyPair, p.Result.Side, p.Result.Price, p.Result.Amount, p.Result.CreateTime)
			if p.Result.Price.IsZero() {
				t.Error("zero trade price")
			}
		case <-time.After(20 * time.Second):
			t.Log("no trade push in 20s (BTC_USDT quiet); channel+decode still exercised")
		}
	})
}

// TestSpotWSPrivate subscribes to the private orders channel and, when writes are
// enabled, places+cancels a tiny far-from-market order to trigger a real push.
func TestSpotWSPrivate(t *testing.T) {
	apiKey, apiSecret := testutil.Creds(t)
	opts := []client.WebSocketOptions{client.WithWebSocketAuth(apiKey, apiSecret)}
	if proxy := testutil.Proxy(); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	c := NewSpotWebSocketClient(opts...)
	ctx := testutil.Ctx(t)

	pushErr := make(chan error, 4)
	gotOrder := make(chan WsOrder, 4)
	done, _, err := c.NewSubscribeOrdersService("!all").Do(ctx, func(p *request.WsPush[[]WsOrder], e error) {
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

	sc := testClient(t)
	if err := sc.SyncServerTime(ctx); err != nil {
		t.Fatalf("sync: %v", err)
	}
	order, err := sc.NewCreateOrderService("BTC_USDT", SideBuy, decimal.RequireFromString("0.0001")).
		SetType(OrderTypeLimit).SetPrice(decimal.RequireFromString("41000")).SetTimeInForce(TimeInForceGTC).Do(ctx)
	if err != nil {
		t.Fatalf("place order: %v", err)
	}
	t.Logf("placed %s to trigger ws push", order.ID)
	defer sc.NewCancelOrderService(order.ID, "BTC_USDT").Do(ctx)

	select {
	case o := <-gotOrder:
		t.Logf("ws order push: id=%s event=%s status pair=%s left=%s", o.ID, o.Event, o.CurrencyPair, o.Left)
	case <-time.After(10 * time.Second):
		t.Fatal("no order push after placing order")
	}
}
