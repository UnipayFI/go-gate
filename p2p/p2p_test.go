package p2p

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestP2P(t *testing.T) {
	// ---- Private read endpoints ----

	t.Run("GetUserInfo", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewGetUserInfoService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/account/get_user_info", err) {
				return
			}
			t.Fatalf("get user info: %v", err)
		}
		t.Logf("get_user_info: code=%d msg=%s", resp.Code, resp.Message)
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/account/get_user_info",
			map[string]any{}, true)
		testutil.AssertCovers(t, "p2p/account/get_user_info", raw, resp)
	})

	t.Run("GetCounterpartyUserInfo", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		bizUID := "0"
		resp, err := c.NewGetCounterpartyUserInfoService(bizUID).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/account/get_counterparty_user_info", err) {
				return
			}
			t.Fatalf("get counterparty user info: %v", err)
		}
		t.Logf("get_counterparty_user_info: code=%d msg=%s", resp.Code, resp.Message)
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/account/get_counterparty_user_info",
			map[string]any{"biz_uid": bizUID}, true)
		testutil.AssertCovers(t, "p2p/account/get_counterparty_user_info", raw, resp)
	})

	t.Run("GetMyselfPayment", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewGetMyselfPaymentService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/account/get_myself_payment", err) {
				return
			}
			t.Fatalf("get myself payment: %v", err)
		}
		t.Logf("get_myself_payment: code=%d groups=%d", resp.Code, len(resp.Data))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/account/get_myself_payment",
			map[string]any{}, true)
		testutil.AssertCovers(t, "p2p/account/get_myself_payment", raw, resp)
	})

	t.Run("GetPendingTransactionList", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewGetPendingTransactionListService("USDT", "USD").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/transaction/get_pending_transaction_list", err) {
				return
			}
			t.Fatalf("get pending transaction list: %v", err)
		}
		t.Logf("get_pending_transaction_list: code=%d list=%d", resp.Code, len(resp.Data.List))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/transaction/get_pending_transaction_list",
			map[string]any{"crypto_currency": "USDT", "fiat_currency": "USD"}, true)
		testutil.AssertCovers(t, "p2p/transaction/get_pending_transaction_list", raw, resp)
	})

	t.Run("GetCompletedTransactionList", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewGetCompletedTransactionListService("USDT", "USD").SetPerPage(2).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/transaction/get_completed_transaction_list", err) {
				return
			}
			t.Fatalf("get completed transaction list: %v", err)
		}
		t.Logf("get_completed_transaction_list: code=%d list=%d", resp.Code, len(resp.Data.List))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/transaction/get_completed_transaction_list",
			map[string]any{"crypto_currency": "USDT", "fiat_currency": "USD", "per_page": 2}, true)
		testutil.AssertCovers(t, "p2p/transaction/get_completed_transaction_list", raw, resp)
	})

	t.Run("GetTransactionDetails", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewGetTransactionDetailsService(0).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/transaction/get_transaction_details", err) {
				return
			}
			t.Fatalf("get transaction details: %v", err)
		}
		t.Logf("get_transaction_details: code=%d msg=%s", resp.Code, resp.Message)
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/transaction/get_transaction_details",
			map[string]any{"txid": int64(0)}, true)
		testutil.AssertCovers(t, "p2p/transaction/get_transaction_details", raw, resp)
	})

	t.Run("AdsDetail", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewAdsDetailService("0").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/books/ads_detail", err) {
				return
			}
			t.Fatalf("ads detail: %v", err)
		}
		t.Logf("ads_detail: code=%d msg=%s", resp.Code, resp.Message)
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/books/ads_detail",
			map[string]any{"adv_no": "0"}, true)
		testutil.AssertCovers(t, "p2p/books/ads_detail", raw, resp)
	})

	t.Run("MyAdsList", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewMyAdsListService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/books/my_ads_list", err) {
				return
			}
			t.Fatalf("my ads list: %v", err)
		}
		t.Logf("my_ads_list: code=%d lists=%d", resp.Code, len(resp.Data.Lists))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/books/my_ads_list",
			map[string]any{}, true)
		testutil.AssertCovers(t, "p2p/books/my_ads_list", raw, resp)
	})

	t.Run("AdsList", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewAdsListService("USDT", "USD", "buy").Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/books/ads_list", err) {
				return
			}
			t.Fatalf("ads list: %v", err)
		}
		t.Logf("ads_list: code=%d ads=%d", resp.Code, len(resp.Data))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/books/ads_list",
			map[string]any{"asset": "USDT", "fiat_unit": "USD", "trade_type": "buy"}, true)
		testutil.AssertCovers(t, "p2p/books/ads_list", raw, resp)
	})

	t.Run("GetChatsList", func(t *testing.T) {
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		resp, err := c.NewGetChatsListService().Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "p2p/chat/get_chats_list", err) {
				return
			}
			t.Fatalf("get chats list: %v", err)
		}
		t.Logf("get_chats_list: code=%d messages=%d", resp.Code, len(resp.Data.Messages))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/p2p/merchant/chat/get_chats_list",
			map[string]any{}, true)
		testutil.AssertCovers(t, "p2p/chat/get_chats_list", raw, resp)
	})

	// ---- Private write endpoints ----
	//
	// State-changing calls. Gated behind GATE_TEST_WRITE and exercised with
	// tiny / invalid parameters so they are rejected by Gate and never actually
	// publish an ad, confirm a payment or change merchant state; any error is a
	// pass.

	t.Run("SetMerchantWorkHours", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSetMerchantWorkHoursService(0).Do(cx)
		if err != nil {
			t.Logf("set merchant work hours: %v (tolerable)", err)
			return
		}
		t.Log("set merchant work hours accepted")
	})

	t.Run("ConfirmPayment", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewConfirmPaymentService("0").Do(cx)
		if err != nil {
			t.Logf("confirm payment: %v (tolerable)", err)
			return
		}
		t.Log("confirm payment accepted")
	})

	t.Run("ConfirmReceipt", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewConfirmReceiptService("0").Do(cx)
		if err != nil {
			t.Logf("confirm receipt: %v (tolerable)", err)
			return
		}
		t.Log("confirm receipt accepted")
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

	t.Run("PlaceBizPushOrder", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		tiny := decimal.NewFromFloat(0.00000001)
		_, err := c.NewPlaceBizPushOrderService("USDT", "USD", "0", tiny, tiny, "bank").Do(cx)
		if err != nil {
			t.Logf("place biz push order: %v (tolerable)", err)
			return
		}
		t.Log("place biz push order accepted")
	})

	t.Run("AdsUpdateStatus", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewAdsUpdateStatusService(0, 3).Do(cx)
		if err != nil {
			t.Logf("ads update status: %v (tolerable)", err)
			return
		}
		t.Log("ads update status accepted")
	})

	t.Run("SendChatMessage", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewSendChatMessageService(0, "test").Do(cx)
		if err != nil {
			t.Logf("send chat message: %v (tolerable)", err)
			return
		}
		t.Log("send chat message accepted")
	})

	t.Run("UploadChatFile", func(t *testing.T) {
		if !testutil.WriteEnabled() {
			t.Skip("write disabled; set GATE_TEST_WRITE=1 to run")
		}
		c := testClient(t)
		cx := testutil.Ctx(t)
		if err := c.SyncServerTime(cx); err != nil {
			t.Fatalf("sync time: %v", err)
		}
		_, err := c.NewUploadChatFileService("image/png", "dGVzdA==").Do(cx)
		if err != nil {
			t.Logf("upload chat file: %v (tolerable)", err)
			return
		}
		t.Log("upload chat file accepted")
	})
}
