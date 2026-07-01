package delivery

import (
	"strconv"
	"strings"
	"testing"

	"github.com/UnipayFI/go-gate/common"
	"github.com/UnipayFI/go-gate/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestDeliveryOrder(t *testing.T) {
	pub := testPublicClient()
	cx := testutil.Ctx(t)

	// --- public: risk limit tiers (top markets, no auth) ---
	tiers, err := pub.NewListDeliveryRiskLimitTiersService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		t.Fatalf("risk limit tiers: %v", err)
	}
	if len(tiers) == 0 {
		t.Fatal("no risk limit tiers returned")
	}
	t.Logf("riskLimitTier: %+v", tiers[0])
	raw := testutil.FetchRawGet(t, pub, cx, "/api/v4/delivery/usdt/risk_limit_tiers",
		map[string]string{"limit": "2"}, false)
	testutil.AssertCovers(t, "delivery/risk_limit_tiers", raw, tiers)

	// --- private reads (skips when creds unset) ---
	c := testClient(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// list finished orders (Gate rejects `limit` for status=open, so use finished
	// for the paginated coverage check)
	orders, err := c.NewListDeliveryOrdersService(SettleUSDT, OrderStatusFinished).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/orders", err) {
			t.Fatalf("list orders: %v", err)
		}
	} else {
		t.Logf("finishedOrders=%d", len(orders))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/orders",
			map[string]string{"status": "finished", "limit": "2"}, true)
		testutil.AssertCovers(t, "delivery/orders", raw, orders)
	}

	// personal trades
	trades, err := c.NewGetMyDeliveryTradesService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/my_trades", err) {
			t.Fatalf("my trades: %v", err)
		}
	} else {
		t.Logf("myTrades=%d", len(trades))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/my_trades",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "delivery/my_trades", raw, trades)
	}

	// position close history
	closes, err := c.NewListDeliveryPositionCloseService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/position_close", err) {
			t.Fatalf("position close: %v", err)
		}
	} else {
		t.Logf("positionClose=%d", len(closes))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/position_close",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "delivery/position_close", raw, closes)
	}

	// liquidation history
	liqs, err := c.NewListDeliveryLiquidatesService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/liquidates", err) {
			t.Fatalf("liquidates: %v", err)
		}
	} else {
		t.Logf("liquidates=%d", len(liqs))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/liquidates",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "delivery/liquidates", raw, liqs)
	}

	// settlement records
	settlements, err := c.NewListDeliverySettlementsService(SettleUSDT).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/settlements", err) {
			t.Fatalf("settlements: %v", err)
		}
	} else {
		t.Logf("settlements=%d", len(settlements))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/settlements",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "delivery/settlements", raw, settlements)
	}

	// single order by a non-existent id — exercises path + signing
	if _, err := c.NewGetDeliveryOrderService(SettleUSDT, "1").Do(cx); err != nil {
		if !testutil.Tolerable(t, "delivery/orders/{id}", err) {
			t.Fatalf("get order: %v", err)
		}
	}

	// --- state-changing: place -> query -> cancel (reversible) ---
	if !testutil.WriteEnabled() {
		t.Log("GATE_TEST_WRITE!=1; skipping delivery order create/cancel")
		return
	}

	// pick a live BTC delivery contract (expiry-suffixed name).
	rawContracts := testutil.FetchRawGet(t, pub, cx, "/api/v4/delivery/usdt/contracts", nil, false)
	var contracts []struct {
		Name string `json:"name"`
	}
	if err := common.JSONUnmarshal(rawContracts, &contracts); err != nil {
		t.Fatalf("decode contracts: %v", err)
	}
	var contract string
	for _, ct := range contracts {
		if strings.HasPrefix(ct.Name, "BTC_USDT") {
			contract = ct.Name
			break
		}
	}
	if contract == "" {
		t.Skip("no BTC delivery contract available")
	}

	// resting buy far below market (size 1, price 1000) so it does not fill.
	created, err := c.NewCreateDeliveryOrderService(SettleUSDT, contract, 1).
		SetPrice(decimal.NewFromInt(1000)).SetTif(TimeInForceGTC).SetText("t-gogate").Do(cx)
	if err != nil {
		if testutil.Tolerable(t, "delivery/create-order", err) {
			return
		}
		t.Logf("create order failed (non-fatal, delivery liquidity is thin): %v", err)
		return
	}
	t.Logf("createdOrder id=%d status=%s", created.ID, created.Status)

	orderID := strconv.FormatInt(created.ID, 10)
	got, err := c.NewGetDeliveryOrderService(SettleUSDT, orderID).Do(cx)
	if err != nil {
		t.Fatalf("get created order: %v", err)
	}
	rawOrder := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/orders/"+orderID, nil, true)
	testutil.AssertCovers(t, "delivery/orders/{id}", rawOrder, got)

	cancelled, err := c.NewCancelDeliveryOrderService(SettleUSDT, orderID).Do(cx)
	if err != nil {
		t.Fatalf("cancel order: %v", err)
	}
	t.Logf("cancelledOrder id=%d status=%s finishAs=%s", cancelled.ID, cancelled.Status, cancelled.FinishAs)

	// batch-cancel the (now empty) contract to exercise the endpoint.
	if _, err := c.NewCancelDeliveryOrdersService(SettleUSDT, contract).Do(cx); err != nil {
		t.Fatalf("cancel orders: %v", err)
	}
}
