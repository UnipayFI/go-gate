package options

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestOptionsAccount(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	const underlying = "BTC_USDT"

	// GET /api/v4/options/accounts
	acct, err := c.NewListOptionsAccountService().Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/accounts", err) {
			t.Fatalf("options accounts: %v", err)
		}
	} else {
		t.Logf("account: total=%s equity=%s available=%s currency=%s marginMode=%d",
			acct.Total, acct.Equity, acct.Available, acct.Currency, acct.MarginMode)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/accounts", nil, true)
		testutil.AssertCovers(t, "options/accounts", raw, acct)
	}

	// GET /api/v4/options/account_book
	book, err := c.NewListOptionsAccountBookService().SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/account_book", err) {
			t.Fatalf("options account_book: %v", err)
		}
	} else if len(book) == 0 {
		t.Logf("options account_book: empty")
	} else {
		t.Logf("account_book[0]: type=%s change=%s balance=%s",
			book[0].Type, book[0].Change, book[0].Balance)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/account_book",
			map[string]string{"limit": "2"}, true)
		testutil.AssertCovers(t, "options/account_book", raw, book)
	}

	// GET /api/v4/options/positions
	positions, err := c.NewListOptionsPositionsService().SetUnderlying(underlying).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/positions", err) {
			t.Fatalf("options positions: %v", err)
		}
	} else if len(positions) == 0 {
		t.Logf("options positions: empty")
	} else {
		t.Logf("positions[0]: contract=%s size=%d markPrice=%s",
			positions[0].Contract, positions[0].Size, positions[0].MarkPrice)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/positions",
			map[string]string{"underlying": underlying}, true)
		testutil.AssertCovers(t, "options/positions", raw, positions)

		// GET /api/v4/options/positions/{contract} — use a contract we actually hold.
		pc := positions[0].Contract
		pos, err := c.NewGetOptionsPositionService(pc).Do(cx)
		if err != nil {
			if !testutil.Tolerable(t, "options/positions/"+pc, err) {
				t.Fatalf("options position %s: %v", pc, err)
			}
		} else {
			t.Logf("position %s size=%d markPrice=%s", pos.Contract, pos.Size, pos.MarkPrice)
			raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/positions/"+pc, nil, true)
			testutil.AssertCovers(t, "options/positions/"+pc, raw, pos)
		}
	}

	// GET /api/v4/options/position_close
	closes, err := c.NewListOptionsPositionCloseService(underlying).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/position_close", err) {
			t.Fatalf("options position_close: %v", err)
		}
	} else if len(closes) == 0 {
		t.Logf("options position_close: empty")
	} else {
		t.Logf("position_close[0]: contract=%s side=%s pnl=%s",
			closes[0].Contract, closes[0].Side, closes[0].PnL)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/position_close",
			map[string]string{"underlying": underlying}, true)
		testutil.AssertCovers(t, "options/position_close", raw, closes)
	}

	// GET /api/v4/options/my_settlements
	settlements, err := c.NewListMyOptionsSettlementsService(underlying).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/my_settlements", err) {
			t.Fatalf("options my_settlements: %v", err)
		}
	} else if len(settlements) == 0 {
		t.Logf("options my_settlements: empty")
	} else {
		t.Logf("my_settlements[0]: contract=%s settleProfit=%s realisedPnL=%s",
			settlements[0].Contract, settlements[0].SettleProfit, settlements[0].RealisedPnL)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/my_settlements",
			map[string]string{"underlying": underlying, "limit": "2"}, true)
		testutil.AssertCovers(t, "options/my_settlements", raw, settlements)
	}

	// GET /api/v4/options/my_trades
	trades, err := c.NewListMyOptionsTradesService(underlying).SetLimit(2).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/my_trades", err) {
			t.Fatalf("options my_trades: %v", err)
		}
	} else if len(trades) == 0 {
		t.Logf("options my_trades: empty")
	} else {
		t.Logf("my_trades[0]: contract=%s size=%d price=%s role=%s",
			trades[0].Contract, trades[0].Size, trades[0].Price, trades[0].Role)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/my_trades",
			map[string]string{"underlying": underlying, "limit": "2"}, true)
		testutil.AssertCovers(t, "options/my_trades", raw, trades)
	}

	// GET /api/v4/options/mmp
	mmp, err := c.NewGetOptionsMMPService(underlying).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "options/mmp", err) {
			t.Fatalf("options mmp: %v", err)
		}
	} else {
		t.Logf("mmp: underlying=%s window=%d frozenPeriod=%d qtyLimit=%s deltaLimit=%s",
			mmp.Underlying, mmp.Window, mmp.FrozenPeriod, mmp.QtyLimit, mmp.DeltaLimit)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/options/mmp",
			map[string]string{"underlying": underlying}, true)
		testutil.AssertCovers(t, "options/mmp", raw, mmp)
	}

	// POST /api/v4/options/mmp and /api/v4/options/mmp/reset mutate account state,
	// so they only run when writes are explicitly enabled; even then errors on an
	// account without options/MMP enabled are tolerated. window=0 disables MMP,
	// so this leaves the account in a benign, reversible state.
	if !testutil.WriteEnabled() {
		t.Log("GATE_TEST_WRITE not set; skipping MMP set/reset")
		return
	}

	if set, err := c.NewSetOptionsMMPService(underlying, 0, 0, decimal.NewFromInt(0), decimal.NewFromInt(0)).Do(cx); err != nil {
		if !testutil.Tolerable(t, "options/mmp/set", err) {
			t.Fatalf("set mmp: %v", err)
		}
	} else {
		t.Logf("mmp set: window=%d frozenPeriod=%d", set.Window, set.FrozenPeriod)
	}

	if reset, err := c.NewResetOptionsMMPService(underlying).Do(cx); err != nil {
		if !testutil.Tolerable(t, "options/mmp/reset", err) {
			t.Fatalf("reset mmp: %v", err)
		}
	} else {
		t.Logf("mmp reset: window=%d frozenPeriod=%d", reset.Window, reset.FrozenPeriod)
	}
}
