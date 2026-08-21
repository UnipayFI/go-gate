package futures

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestFuturesStats(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// GET /api/v4/futures/usdt/insurance
	{
		list, err := c.NewListFuturesInsuranceLedgerService(SettleUSDT).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("insurance: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no insurance records returned")
		}
		t.Logf("insurance[0]: t=%s balance=%s", list[0].Timestamp, list[0].Balance)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/insurance",
			map[string]string{"limit": "2"}, false)
		testutil.AssertCovers(t, "futures/usdt/insurance", raw, list)
	}

	// GET /api/v4/futures/usdt/contract_stats
	{
		list, err := c.NewListContractStatsService(SettleUSDT, "BTC_USDT").SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("contract_stats: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no contract stats returned")
		}
		t.Logf("contract_stat[0]: %+v", list[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/contract_stats",
			map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
		testutil.AssertCovers(t, "futures/usdt/contract_stats", raw, list)
	}

	// GET /api/v4/futures/usdt/index_constituents/{index}
	{
		ic, err := c.NewGetIndexConstituentsService(SettleUSDT, "BTC_USDT").Do(cx)
		if err != nil {
			t.Fatalf("index_constituents: %v", err)
		}
		t.Logf("index=%s constituents=%d", ic.Index, len(ic.Constituents))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/index_constituents/BTC_USDT", nil, false)
		testutil.AssertCovers(t, "futures/usdt/index_constituents", raw, ic)
	}

	// GET /api/v4/futures/usdt/liq_orders
	{
		list, err := c.NewListLiquidatedOrdersService(SettleUSDT).SetContract("BTC_USDT").SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("liq_orders: %v", err)
		}
		t.Logf("liq_orders=%d", len(list))
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/liq_orders",
			map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
		testutil.AssertCovers(t, "futures/usdt/liq_orders", raw, list)
	}

	// GET /api/v4/futures/usdt/risk_limit_tiers
	{
		list, err := c.NewListFuturesRiskLimitTiersService(SettleUSDT).SetContract("BTC_USDT").SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("risk_limit_tiers: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no risk limit tiers returned")
		}
		t.Logf("risk_limit_tier[0]: %+v", list[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/risk_limit_tiers",
			map[string]string{"contract": "BTC_USDT", "limit": "2"}, false)
		testutil.AssertCovers(t, "futures/usdt/risk_limit_tiers", raw, list)
	}

	// GET /api/v4/futures/usdt/risk_limit_table
	// table_id has the form CONTRACT_YYYYMMDD; it is not exposed on the contract,
	// so use a known-good example and stay tolerant if it is ever retired.
	{
		const tableID = "CYBER_USDT_20241122"
		list, err := c.NewGetFuturesRiskLimitTableService(SettleUSDT, tableID).Do(cx)
		if err != nil {
			if testutil.Tolerable(t, "futures/usdt/risk_limit_table", err) {
				return
			}
			t.Logf("risk_limit_table (table_id=%s) unavailable: %v", tableID, err)
			return
		}
		if len(list) == 0 {
			t.Logf("risk_limit_table returned no rows for table_id=%s", tableID)
			return
		}
		t.Logf("risk_limit_table[0]: %+v", list[0])
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/risk_limit_table",
			map[string]string{"table_id": tableID}, false)
		testutil.AssertCovers(t, "futures/usdt/risk_limit_table", raw, list)
	}
}

func TestFuturesADLRiskStates(t *testing.T) {
	c := testPublicClient()
	cx := testutil.Ctx(t)

	// GET /api/v4/futures/usdt/adl_risk_states
	states, err := c.NewListFuturesADLRiskStatesService(SettleUSDT).Do(cx)
	if err != nil {
		t.Fatalf("adl_risk_states: %v", err)
	}
	if states.Settle != SettleUSDT {
		t.Fatalf("settle = %q, want %q", states.Settle, SettleUSDT)
	}
	if len(states.States) == 0 {
		t.Fatal("no adl risk states returned")
	}
	one := states.States["BTC_USDT"]
	t.Logf("adl states=%d BTC_USDT: state=%s calculated_at=%s", len(states.States), one.State, one.CalculatedAtMs)
	if one.State == "" {
		t.Fatal("BTC_USDT adl risk state is empty")
	}
	if one.CalculatedAtMs.IsZero() {
		t.Fatal("BTC_USDT calculated_at_ms is zero")
	}
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/futures/usdt/adl_risk_states", nil, false)
	testutil.AssertCovers(t, "futures/usdt/adl_risk_states", raw, states)
}
