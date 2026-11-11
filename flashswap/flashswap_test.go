package flashswap

import (
	"strconv"
	"strings"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestFlashSwap(t *testing.T) {
	t.Run("CurrencyPairs", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		list, err := c.NewListFlashSwapCurrencyPairService().SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("currency pairs: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no currency pairs returned")
		}
		t.Logf("pair: %+v", list[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/flash_swap/currency_pairs",
			map[string]string{"limit": "2"}, false)
		testutil.AssertCovers(t, "flash_swap/currency_pairs", raw, list)
	})

	t.Run("Orders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}

		list, err := c.NewListFlashSwapOrdersService().SetLimit(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "flash_swap/orders", err) {
				return
			}
			t.Fatalf("list orders: %v", err)
		}
		t.Logf("orders=%d", len(list))
		if len(list) == 0 {
			return
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/flash_swap/orders",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "flash_swap/orders", raw, list)

		one, err := c.NewGetFlashSwapOrderService(list[0].ID).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "flash_swap/orders/{order_id}", err) {
				return
			}
			t.Fatalf("get order: %v", err)
		}
		t.Logf("order: %+v", one)
		rawOne := testutil.FetchRawGet(t, c, cx,
			"/api/v4/flash_swap/orders/"+strconv.FormatInt(list[0].ID, 10), nil, true)
		testutil.AssertCovers(t, "flash_swap/orders/{order_id}", rawOne, one)
	})

	t.Run("CurrencyPairV1", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}

		pairs, err := c.NewListFlashSwapCurrencyPairV1Service("BTC", FlashSwapSideSell).Do(cx)
		if err != nil {
			// Announced in v4.106.180, but as of 2026-11-12 the route is a bare
			// openresty 404 on every API host.
			if strings.Contains(err.Error(), "status 404") {
				t.Skipf("flash-swap/v1/currency-pair not deployed yet: %v", err)
			}
			if testutil.Tolerable(t, "flash-swap/v1/currency-pair", err) {
				return
			}
			t.Fatalf("currency pair v1: %v", err)
		}
		t.Logf("pairs=%d", len(pairs.CurrencyPair))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/flash-swap/v1/currency-pair",
			map[string]string{"asset_name": "BTC", "side": "1"}, true)
		testutil.AssertCovers(t, "flash-swap/v1/currency-pair", raw, pairs)
	})

	t.Run("Preview", func(t *testing.T) {
		// Preview only quotes a swap; it does not execute one. CreateFlashSwapOrder
		// would actually swap funds irreversibly, so it is implemented but never
		// called in tests. The whole block is still gated behind GATE_TEST_WRITE.
		if !testutil.WriteEnabled() {
			t.Skip("flash swap preview requires GATE_TEST_WRITE=1")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}

		sellAmount := decimal.RequireFromString("0.0001")
		preview, err := c.NewPreviewFlashSwapOrderService("BTC", "USDT").
			SetSellAmount(sellAmount).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "flash_swap/orders/preview", err) {
				return
			}
			t.Fatalf("preview: %v", err)
		}
		t.Logf("preview: %+v", preview)
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/flash_swap/orders/preview",
			map[string]any{
				"sell_currency": "BTC",
				"buy_currency":  "USDT",
				"sell_amount":   sellAmount.String(),
			}, true)
		testutil.AssertCovers(t, "flash_swap/orders/preview", raw, preview)
	})
}
