package futures

import (
	"strconv"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestFuturesAutoorder(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// ---- Trail auto orders ----

	t.Run("CreateTrailOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise trail create")
		}
		_, err := c.NewCreateTrailOrderService(SettleUSDT, "BTC_USDT", decimal.NewFromInt(1)).
			SetActivationPrice(decimal.NewFromInt(1000000)).
			SetIsGte(true).
			SetPriceOffset("0.1%").
			SetReduceOnly(true).
			Do(cx)
		if err != nil {
			t.Logf("create trail order: %v (tolerable)", err)
			return
		}
		t.Log("create trail order accepted")
	})

	t.Run("StopTrailOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise trail stop")
		}
		_, err := c.NewStopTrailOrderService(SettleUSDT).SetID(1).Do(cx)
		if err != nil {
			t.Logf("stop trail order: %v (tolerable)", err)
			return
		}
		t.Log("stop trail order accepted")
	})

	t.Run("StopAllTrailOrders", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise trail stop_all")
		}
		_, err := c.NewStopAllTrailOrdersService(SettleUSDT).SetContract("BTC_USDT").Do(cx)
		if err != nil {
			t.Logf("stop all trail orders: %v (tolerable)", err)
			return
		}
		t.Log("stop all trail orders accepted")
	})

	t.Run("ListTrailOrders", func(t *testing.T) {
		got, err := c.NewListTrailOrdersService(SettleUSDT).SetContract("BTC_USDT").SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/trail/list", err) {
				return
			}
			t.Fatalf("list trail orders: %v", err)
		}
		t.Logf("trailOrders=%d", len(got.Data.Orders))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/autoorder/v1/trail/list",
			map[string]string{"contract": "BTC_USDT", "page_size": "2"}, true)
		testutil.AssertCovers(t, "futures/autoorder/v1/trail/list", raw, got)
	})

	t.Run("GetTrailOrder", func(t *testing.T) {
		list, err := c.NewListTrailOrdersService(SettleUSDT).SetPageSize(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/trail/detail", err) {
				return
			}
			t.Fatalf("list trail orders: %v", err)
		}
		if len(list.Data.Orders) == 0 {
			t.Skip("no trail orders to fetch detail for")
		}
		id := list.Data.Orders[0].ID
		got, err := c.NewGetTrailOrderService(SettleUSDT, id).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/trail/detail", err) {
				return
			}
			t.Fatalf("get trail order %d: %v", id, err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/autoorder/v1/trail/detail",
			map[string]string{"id": strconv.FormatInt(id, 10)}, true)
		testutil.AssertCovers(t, "futures/autoorder/v1/trail/detail", raw, got)
	})

	t.Run("UpdateTrailOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise trail update")
		}
		_, err := c.NewUpdateTrailOrderService(SettleUSDT, 1).
			SetPriceOffset("0.2%").
			Do(cx)
		if err != nil {
			t.Logf("update trail order: %v (tolerable)", err)
			return
		}
		t.Log("update trail order accepted")
	})

	t.Run("GetTrailChangeLog", func(t *testing.T) {
		list, err := c.NewListTrailOrdersService(SettleUSDT).SetPageSize(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/trail/change_log", err) {
				return
			}
			t.Fatalf("list trail orders: %v", err)
		}
		if len(list.Data.Orders) == 0 {
			t.Skip("no trail orders to fetch change log for")
		}
		id := list.Data.Orders[0].ID
		got, err := c.NewGetTrailChangeLogService(SettleUSDT, id).SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/trail/change_log", err) {
				return
			}
			t.Fatalf("get trail change log %d: %v", id, err)
		}
		t.Logf("trailChangeLog=%d", len(got.Data.ChangeLog))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/autoorder/v1/trail/change_log",
			map[string]string{"id": strconv.FormatInt(id, 10), "page_size": "2"}, true)
		testutil.AssertCovers(t, "futures/autoorder/v1/trail/change_log", raw, got)
	})

	// ---- Chase auto orders ----

	t.Run("CreateChaseOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise chase create")
		}
		_, err := c.NewCreateChaseOrderService(SettleUSDT, "BTC_USDT", decimal.NewFromInt(1), decimal.NewFromInt(1)).
			SetReduceOnly(true).
			Do(cx)
		if err != nil {
			t.Logf("create chase order: %v (tolerable)", err)
			return
		}
		t.Log("create chase order accepted")
	})

	t.Run("StopChaseOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise chase stop")
		}
		_, err := c.NewStopChaseOrderService(SettleUSDT).SetID("1").Do(cx)
		if err != nil {
			t.Logf("stop chase order: %v (tolerable)", err)
			return
		}
		t.Log("stop chase order accepted")
	})

	t.Run("StopAllChaseOrders", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("set GATE_TEST_WRITE=1 to exercise chase stop_all")
		}
		_, err := c.NewStopAllChaseOrdersService(SettleUSDT).SetContract("BTC_USDT").Do(cx)
		if err != nil {
			t.Logf("stop all chase orders: %v (tolerable)", err)
			return
		}
		t.Log("stop all chase orders accepted")
	})

	t.Run("ListChaseOrders", func(t *testing.T) {
		got, err := c.NewListChaseOrdersService(SettleUSDT, 1).SetContract("BTC_USDT").SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/chase/list", err) {
				return
			}
			t.Fatalf("list chase orders: %v", err)
		}
		t.Logf("chaseOrders=%d", len(got.Data.Orders))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/autoorder/v1/chase/list",
			map[string]string{"sort_by": "1", "contract": "BTC_USDT", "page_size": "2"}, true)
		testutil.AssertCovers(t, "futures/autoorder/v1/chase/list", raw, got)
	})

	t.Run("GetChaseOrder", func(t *testing.T) {
		list, err := c.NewListChaseOrdersService(SettleUSDT, 1).SetPageSize(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/chase/detail", err) {
				return
			}
			t.Fatalf("list chase orders: %v", err)
		}
		if len(list.Data.Orders) == 0 {
			t.Skip("no chase orders to fetch detail for")
		}
		id := list.Data.Orders[0].ID
		got, err := c.NewGetChaseOrderService(SettleUSDT, id).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/autoorder/v1/chase/detail", err) {
				return
			}
			t.Fatalf("get chase order %s: %v", id, err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/autoorder/v1/chase/detail",
			map[string]string{"id": id}, true)
		testutil.AssertCovers(t, "futures/autoorder/v1/chase/detail", raw, got)
	})
}
