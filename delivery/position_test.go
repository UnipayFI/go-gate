package delivery

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
	"github.com/shopspring/decimal"
)

// positionTestContract returns a live delivery contract name (expiry-suffixed)
// so single-position lookups have a real {contract} path segment to target.
func positionTestContract(t *testing.T, c *DeliveryClient, cx context.Context) string {
	t.Helper()
	raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/contracts", nil, false)
	var contracts []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &contracts); err != nil {
		t.Fatalf("decode delivery contracts: %v", err)
	}
	if len(contracts) == 0 {
		t.Skip("no delivery contracts available")
	}
	return contracts[0].Name
}

func TestDeliveryPosition(t *testing.T) {
	c := testClient(t)
	cx := testutil.Ctx(t)
	if err := c.SyncServerTime(cx); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	contract := positionTestContract(t, testPublicClient(), cx)
	t.Logf("using delivery contract %s", contract)

	// ListDeliveryPositions -- private read, tolerable when the account holds none.
	list, err := c.NewListDeliveryPositionsService(SettleUSDT).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/positions", err) {
			t.Fatalf("list positions: %v", err)
		}
	} else {
		t.Logf("positions=%d", len(list))
		for _, p := range list {
			t.Logf("  %s size=%d entry=%s upl=%s", p.Contract, p.Size, p.EntryPrice, p.UnrealisedPnL)
		}
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/positions", nil, true)
		testutil.AssertCovers(t, "delivery/positions", raw, list)
	}

	// GetDeliveryPosition -- private read of a single contract's position.
	pos, err := c.NewGetDeliveryPositionService(SettleUSDT, contract).Do(cx)
	if err != nil {
		if !testutil.Tolerable(t, "delivery/positions/"+contract, err) {
			t.Fatalf("get position: %v", err)
		}
	} else {
		t.Logf("position %s size=%d leverage=%s mark=%s", pos.Contract, pos.Size, pos.Leverage, pos.MarkPrice)
		raw := testutil.FetchRawGet(t, c, cx, "/api/v4/delivery/usdt/positions/"+contract, nil, true)
		testutil.AssertCovers(t, "delivery/positions/"+contract, raw, pos)
	}

	// Margin / leverage / risk-limit updates mutate account state, so they only
	// run when writes are explicitly enabled; even then errors on an unheld
	// contract are tolerated.
	if !testutil.WriteEnabled() {
		t.Log("GATE_TEST_WRITE not set; skipping margin/leverage/risk_limit updates")
		return
	}

	// These return the same DeliveryPosition shape already covered by the read
	// endpoints above, so here we only exercise the request path + signing.
	if lev, err := c.NewUpdateDeliveryPositionLeverageService(SettleUSDT, contract, decimal.NewFromInt(0)).Do(cx); err != nil {
		if !testutil.Tolerable(t, "delivery/positions/leverage", err) {
			t.Fatalf("update leverage: %v", err)
		}
	} else {
		t.Logf("leverage updated: %s", lev.Leverage)
	}

	if rl, err := c.NewUpdateDeliveryPositionRiskLimitService(SettleUSDT, contract, decimal.NewFromInt(0)).Do(cx); err != nil {
		if !testutil.Tolerable(t, "delivery/positions/risk_limit", err) {
			t.Fatalf("update risk_limit: %v", err)
		}
	} else {
		t.Logf("risk_limit updated: %s", rl.RiskLimit)
	}

	if mg, err := c.NewUpdateDeliveryPositionMarginService(SettleUSDT, contract, decimal.NewFromInt(1)).Do(cx); err != nil {
		if !testutil.Tolerable(t, "delivery/positions/margin", err) {
			t.Fatalf("update margin: %v", err)
		}
	} else {
		t.Logf("margin updated: %s", mg.Margin)
	}
}
