package crossex

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/common"
	"github.com/shopspring/decimal"
)

func TestCrossexV4106186Decode(t *testing.T) {
	var coin CrossexTransferCoin
	if err := common.JSONUnmarshal([]byte(`{"coin":"USDT","min_trans_amount":"10.5","est_fee":"0.25","precision":6,"is_disabled":0}`), &coin); err != nil {
		t.Fatal(err)
	}
	if !coin.MinTransAmount.Equal(decimal.RequireFromString("10.5")) || !coin.EstFee.Equal(decimal.RequireFromString("0.25")) {
		t.Errorf("transfer coin: min_trans_amount=%s est_fee=%s", coin.MinTransAmount, coin.EstFee)
	}

	var tr CrossexTrade
	if err := common.JSONUnmarshal([]byte(`{"order_id":"1","symbol":"US_STOCK_AAPL_USDT","qty":"2","price":"200","fx_rate":"1.0005","create_time":"1700000000000"}`), &tr); err != nil {
		t.Fatal(err)
	}
	if !tr.FXRate.Equal(decimal.RequireFromString("1.0005")) {
		t.Errorf("trade: fx_rate=%s", tr.FXRate)
	}

	var pos CrossexPosition
	if err := common.JSONUnmarshal([]byte(`{"symbol":"BINANCE_FUTURE_BTC_USDT","isolated_margin":"0","margin_mode":"CROSS","liq_price":"0","risk_limits":[{"symbol":"BINANCE_FUTURE_BTC_USDT","tiers":[{"min_risk_limit_value":"0","max_risk_limit_value":"50000","quick_cal_amount":"0","leverage_max":"50","maintenance_rate":"0.005","tier":"1"}]}]}`), &pos); err != nil {
		t.Fatal(err)
	}
	if len(pos.RiskLimits) != 1 || len(pos.RiskLimits[0].Tiers) != 1 || pos.RiskLimits[0].Tiers[0].Tier != "1" {
		t.Errorf("position risk_limits: %+v", pos.RiskLimits)
	}

	var acct CrossexAccount
	if err := common.JSONUnmarshal([]byte(`{"user_id":"1","assets":[{"coin":"BTC","price":"95000.5","balance":"1"}]}`), &acct); err != nil {
		t.Fatal(err)
	}
	if len(acct.Assets) != 1 || !acct.Assets[0].Price.Equal(decimal.RequireFromString("95000.5")) {
		t.Errorf("account asset price: %+v", acct.Assets)
	}
}
