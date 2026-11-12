package futures

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-gate/v4/common"
)

// TestMicrosecondFields decodes the v4.106.181 docs samples offline: live
// payloads did not carry the *_us keys yet when they were added.
func TestMicrosecondFields(t *testing.T) {
	var ob FuturesOrderBook
	if err := common.JSONUnmarshal([]byte(`{"id":123456,"current":1623898993.123456,"current_us":1623898993123456,
		"update":1623898993.121456,"update_us":1623898993121456,"asks":[],"bids":[]}`), &ob); err != nil {
		t.Fatal(err)
	}
	if !ob.CurrentUs.Equal(time.UnixMicro(1623898993123456)) || !ob.UpdateUs.Equal(time.UnixMicro(1623898993121456)) {
		t.Errorf("order book: current_us=%s update_us=%s", ob.CurrentUs, ob.UpdateUs)
	}
	if ob.CurrentUs.Year() != 2021 {
		t.Errorf("order book: current_us decoded to %s, want 2021", ob.CurrentUs)
	}

	var o FuturesOrder
	if err := common.JSONUnmarshal([]byte(`{"id":1,"create_time":1546569968.123456,"create_time_us":1546569968123456,
		"update_time":1546569970.456789,"update_time_us":1546569970456789}`), &o); err != nil {
		t.Fatal(err)
	}
	if !o.CreateTimeUs.Equal(time.UnixMicro(1546569968123456)) || !o.UpdateTimeUs.Equal(time.UnixMicro(1546569970456789)) {
		t.Errorf("order: create_time_us=%s update_time_us=%s", o.CreateTimeUs, o.UpdateTimeUs)
	}
	if !o.FinishTimeUs.IsZero() {
		t.Errorf("order: omitted finish_time_us decoded to %s, want zero", o.FinishTimeUs)
	}
}
