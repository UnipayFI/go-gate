package otc

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestOTC(t *testing.T) {
	// ---------- private reads ----------

	t.Run("ListBanks", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListBanksService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/bank/list", err) {
				return
			}
			t.Fatalf("list banks: %v", err)
		}
		t.Logf("banks=%d", len(got.Data.Lists))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/otc/bank/list", nil, true)
		testutil.AssertCovers(t, "otc/bank/list", raw, got)
	})

	t.Run("ListOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListOrdersService().SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/order/list", err) {
				return
			}
			t.Fatalf("list orders: %v", err)
		}
		t.Logf("fiat orders=%d", len(got.Data.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/otc/order/list",
			map[string]string{"ps": "2"}, true)
		testutil.AssertCovers(t, "otc/order/list", raw, got)
	})

	t.Run("ListStableCoinOrders", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		got, err := c.NewListStableCoinOrdersService().SetPageSize(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/stable_coin/order/list", err) {
				return
			}
			t.Fatalf("list stablecoin orders: %v", err)
		}
		t.Logf("stablecoin orders=%d", len(got.Data.List))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/otc/stable_coin/order/list",
			map[string]string{"page_size": "2"}, true)
		testutil.AssertCovers(t, "otc/stable_coin/order/list", raw, got)
	})

	t.Run("GetOrderDetail", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		list, err := c.NewListOrdersService().SetPageSize(1).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/order/list", err) {
				return
			}
			t.Fatalf("list orders: %v", err)
		}
		if len(list.Data.List) == 0 {
			t.Skip("no fiat orders to fetch detail for")
		}
		orderID := list.Data.List[0].OrderID
		got, err := c.NewGetOrderDetailService(orderID).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/order/detail", err) {
				return
			}
			t.Fatalf("order detail: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/otc/order/detail",
			map[string]string{"order_id": orderID}, true)
		testutil.AssertCovers(t, "otc/order/detail", raw, got)
	})

	t.Run("GetBankSupplementChecklist", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		banks, err := c.NewListBanksService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/bank/list", err) {
				return
			}
			t.Fatalf("list banks: %v", err)
		}
		if len(banks.Data.Lists) == 0 {
			t.Skip("no bank cards to fetch checklist for")
		}
		bankID := banks.Data.Lists[0].ID
		got, err := c.NewGetBankSupplementChecklistService(bankID).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "otc/bank/bank_supplement_checklist", err) {
				return
			}
			t.Fatalf("bank supplement checklist: %v", err)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/otc/bank/bank_supplement_checklist",
			map[string]string{"bank_id": bankID}, true)
		testutil.AssertCovers(t, "otc/bank/bank_supplement_checklist", raw, got)
	})

	// ---------- private writes ----------
	// State-changing (or token-minting) POSTs. Gated behind GATE_TEST_WRITE and
	// exercised with tiny / likely-rejected parameters so they never create a real
	// order, quote token or bank card; any error is treated as a pass.

	t.Run("Quote", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewQuoteService("PAY", "USD", "USDT").
			SetPayAmount(decimal.NewFromFloat(0.00000001)).Do(cx)
		if err != nil {
			t.Logf("quote: %v (tolerable)", err)
			return
		}
		t.Log("quote accepted")
	})

	t.Run("CreateOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateOrderService("BUY", "PAY", "USDT", "USD",
			decimal.NewFromFloat(0.00000001), decimal.NewFromFloat(0.00000001),
			"invalid-quote-token", "0").Do(cx)
		if err != nil {
			t.Logf("create order: %v (tolerable)", err)
			return
		}
		t.Log("create order accepted")
	})

	t.Run("CreateStableCoinOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateStableCoinOrderService(
			"USD", "USDT",
			decimal.NewFromFloat(0.00000001), decimal.NewFromFloat(0.00000001),
			"PAY", "invalid-quote-token").Do(cx)
		if err != nil {
			t.Logf("create stablecoin order: %v (tolerable)", err)
			return
		}
		t.Log("create stablecoin order accepted")
	})

	t.Run("CreateBank", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCreateBankService("gogate-test", "gogate-test", "US",
			"nowhere", "INVALIDIBAN", "INVALIDSWIFT", "").Do(cx)
		if err != nil {
			t.Logf("create bank: %v (tolerable)", err)
			return
		}
		t.Log("create bank accepted")
	})

	t.Run("DeleteBank", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewDeleteBankService("0").Do(cx)
		if err != nil {
			t.Logf("delete bank: %v (tolerable)", err)
			return
		}
		t.Log("delete bank accepted")
	})

	t.Run("SetDefaultBank", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSetDefaultBankService("0").Do(cx)
		if err != nil {
			t.Logf("set default bank: %v (tolerable)", err)
			return
		}
		t.Log("set default bank accepted")
	})

	t.Run("SubmitPersonalBankSupplement", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSubmitPersonalBankSupplementService("0", "", "", "").Do(cx)
		if err != nil {
			t.Logf("submit personal bank supplement: %v (tolerable)", err)
			return
		}
		t.Log("submit personal bank supplement accepted")
	})

	t.Run("SubmitEnterpriseBankSupplement", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSubmitEnterpriseBankSupplementService("0", "", "", "", "").Do(cx)
		if err != nil {
			t.Logf("submit enterprise bank supplement: %v (tolerable)", err)
			return
		}
		t.Log("submit enterprise bank supplement accepted")
	})

	t.Run("MarkOrderPaid", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewMarkOrderPaidService("0", "invalid-file-key").Do(cx)
		if err != nil {
			t.Logf("mark order paid: %v (tolerable)", err)
			return
		}
		t.Log("mark order paid accepted")
	})

	t.Run("CancelOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewCancelOrderService("0").Do(cx)
		if err != nil {
			t.Logf("cancel order: %v (tolerable)", err)
			return
		}
		t.Log("cancel order accepted")
	})
}
