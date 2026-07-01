package delivery

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestDeliveryPriceOrder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// List running auto orders (private read).
	list, err := c.NewListPriceTriggeredDeliveryOrdersService(SettleUSDT, "open").SetLimit(2).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "delivery/price_orders list", err) {
			return
		}
		t.Fatalf("list price orders: %v", err)
	}
	t.Logf("price orders=%d", len(list))
	for _, o := range list {
		t.Logf("  id=%d contract=%s status=%s size=%d", o.ID, o.Initial.Contract, o.Status, o.Initial.Size)
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/price_orders",
		map[string]string{"status": "open", "limit": "2"}, true)
	testutil.AssertCovers(t, "delivery/price_orders", raw, list)

	// Query a single auto order if one already exists.
	if len(list) > 0 {
		id := strconv.FormatInt(list[0].ID, 10)
		one, err := c.NewGetPriceTriggeredDeliveryOrderService(SettleUSDT, id).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "delivery/price_orders/{id}", err) {
				t.Fatalf("get price order: %v", err)
			}
		} else {
			t.Logf("order %s status=%s", id, one.Status)
			oneRaw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/price_orders/"+id, nil, true)
			testutil.AssertCovers(t, "delivery/price_orders/{id}", oneRaw, one)
		}
	}

	// State-changing: place -> query -> cancel a tiny, far-from-market auto order.
	if !testutil.WriteEnabled() {
		t.Skip("set GATE_TEST_WRITE=1 to exercise create/cancel delivery price orders")
	}

	// Delivery contracts carry expiry-suffixed names, so fetch a live BTC_USDT
	// contract and its last price to build a trigger that never fires.
	contract, last := fetchDeliveryContract(t, c, cx)
	t.Logf("using contract=%s last=%s", contract, last)

	// Cancel any leftover auto orders on this contract first (exercises the
	// cancel-all endpoint; tolerable when there is nothing to cancel).
	if _, err := c.NewCancelPriceTriggeredDeliveryOrderListService(SettleUSDT, contract).Do(cx); err != nil {
		if !testutil.Tolerable(t, "delivery/price_orders cancel-all", err) {
			t.Fatalf("cancel-all price orders: %v", err)
		}
	}

	// Trigger a market buy only if price rises to 3x market (far above), so the
	// order stays queued and never fires; rule 1 = trigger when price >= value.
	triggerPrice := last.Mul(decimal.NewFromInt(3)).Ceil()
	created, err := c.NewCreatePriceTriggeredDeliveryOrderService(
		SettleUSDT, contract, 1, decimal.Zero, triggerPrice, 1, 86400,
	).SetInitialTif(TimeInForceIOC).Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "delivery/price_orders create", err) {
			return
		}
		t.Fatalf("create price order: %v", err)
	}
	t.Logf("created price order id=%d", created.ID)

	id := strconv.FormatInt(created.ID, 10)
	got, err := c.NewGetPriceTriggeredDeliveryOrderService(SettleUSDT, id).Do(cx)
	if err != nil {
		t.Fatalf("get created price order: %v", err)
	}
	t.Logf("queried order status=%s", got.Status)
	gotRaw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/price_orders/"+id, nil, true)
	testutil.AssertCovers(t, "delivery/price_orders/{id}", gotRaw, got)

	cancelled, err := c.NewCancelPriceTriggeredDeliveryOrderService(SettleUSDT, id).Do(cx)
	if err != nil {
		t.Fatalf("cancel price order: %v", err)
	}
	t.Logf("cancelled order status=%s", cancelled.Status)
}

// fetchDeliveryContract returns a live USDT delivery contract name (preferring
// BTC_USDT) and its last price, for building reversible write tests.
func fetchDeliveryContract(t *testing.T, c *DeliveryClient, cx context.Context) (string, decimal.Decimal) {
	t.Helper()
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/contracts", nil, false)
	var contracts []struct {
		Name      string `json:"name"`
		LastPrice string `json:"last_price"`
	}
	if err := json.Unmarshal(raw, &contracts); err != nil {
		t.Fatalf("decode contracts: %v", err)
	}
	if len(contracts) == 0 {
		t.Skip("no live delivery contracts available")
	}
	pick := contracts[0]
	for _, ct := range contracts {
		if strings.HasPrefix(ct.Name, "BTC_USDT") {
			pick = ct
			break
		}
	}
	last, err := decimal.NewFromString(pick.LastPrice)
	if err != nil || last.IsZero() {
		last = decimal.NewFromInt(100)
	}
	return pick.Name, last
}
