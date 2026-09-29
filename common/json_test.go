package common

import (
	"encoding/json/v2"
	"strings"
	"sync"
	"testing"
	"time"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

func TestJSONSupport(t *testing.T) {
	if errJSONSupport != nil {
		t.Fatal(errJSONSupport)
	}
}

// TestUnixMatchesStandard checks that values in the declared unit decode
// exactly like the standard `format` option (same instant, same UTC location,
// same errors), whether Gate sends them quoted or bare.
func TestUnixMatchesStandard(t *testing.T) {
	std := jsonexp.ExperimentalSupportFormatTag(true)
	values := map[string][]string{
		"unix":      {"1790646162", "1790646080.498", "1546569968.123456", "0", "1", "-1500", "12.000001", "999999999", "99999999999"},
		"unixmilli": {"1790646162123", "1790646118431.553000", "1620696347990", "0", "1", "-1500", "12.000001", "999999999999"},
		"unixmicro": {"1767262881251901", "1767262881251901.5", "0", "1", "-1500", "999999999999999"},
		"unixnano":  {"1790646162123456789", "9223372036854775807", "0", "1", "-1500"},
	}
	invalid := []string{"1e3", "01", "+1", "1.", ".5", "", "-", "1.2.3", "18446744073709551616", "99999999999999999999999"}
	compared := 0
	for format, vals := range values {
		for _, n := range append(vals, invalid...) {
			for _, stringTag := range []bool{false, true} {
				want, wantErr := decodeTagged(std, format, stringTag, quoteIf(stringTag, n))
				for _, quoted := range []bool{false, true} {
					if n == "" && !quoted {
						continue
					}
					got, err := decodeTagged(nil, format, stringTag, quoteIf(quoted, n))
					if (wantErr != nil) != (err != nil) {
						t.Errorf("%s string=%v %s quoted=%v: errors differ: std=%v codec=%v", format, stringTag, n, quoted, wantErr, err)
						continue
					}
					if wantErr != nil {
						continue
					}
					compared++
					if want != got || got.Location() != time.UTC {
						t.Errorf("%s string=%v %s quoted=%v: std=%v codec=%v", format, stringTag, n, quoted, want, got)
					}
				}
			}
		}
	}
	if compared < 100 {
		t.Errorf("only %d successful comparisons", compared)
	}
}

func quoteIf(quote bool, s string) string {
	if quote {
		return `"` + s + `"`
	}
	return s
}

// decodeTagged decodes {"t":raw} into a field tagged with format (plus
// ,string when stringTag), using the standard library alone when std is set
// and the Gate codec otherwise.
func decodeTagged(std json.Options, format string, stringTag bool, raw string) (time.Time, error) {
	in := []byte(`{"t":` + raw + `}`)
	decode := func(v any) error {
		if std != nil {
			return json.Unmarshal(in, v, std)
		}
		return JSONUnmarshal(in, v)
	}
	switch {
	case format == "unix" && !stringTag:
		var v struct {
			T time.Time `json:"t,format:unix"`
		}
		err := decode(&v)
		return v.T, err
	case format == "unix":
		var v struct {
			T time.Time `json:"t,string,format:unix"`
		}
		err := decode(&v)
		return v.T, err
	case format == "unixmilli" && !stringTag:
		var v struct {
			T time.Time `json:"t,format:unixmilli"`
		}
		err := decode(&v)
		return v.T, err
	case format == "unixmilli":
		var v struct {
			T time.Time `json:"t,string,format:unixmilli"`
		}
		err := decode(&v)
		return v.T, err
	case format == "unixmicro" && !stringTag:
		var v struct {
			T time.Time `json:"t,format:unixmicro"`
		}
		err := decode(&v)
		return v.T, err
	case format == "unixmicro":
		var v struct {
			T time.Time `json:"t,string,format:unixmicro"`
		}
		err := decode(&v)
		return v.T, err
	case format == "unixnano" && !stringTag:
		var v struct {
			T time.Time `json:"t,format:unixnano"`
		}
		err := decode(&v)
		return v.T, err
	default:
		var v struct {
			T time.Time `json:"t,string,format:unixnano"`
		}
		err := decode(&v)
		return v.T, err
	}
}

// TestUnitByDigits: a 10/13/16/19-digit value is read in s/ms/µs/ns whatever
// the declared unit, since the declared unit would misdate it by decades or
// millennia (e.g. Gate's microsecond create_time when cancelling a futures
// price-triggered order, a field that is seconds everywhere else).
func TestUnitByDigits(t *testing.T) {
	want := map[string]time.Time{
		"1767262881":          time.Unix(1767262881, 0),
		"1767262881.5":        time.Unix(1767262881, 500_000_000),
		"1767262881251":       time.UnixMilli(1767262881251),
		"1767262881251.901":   time.UnixMicro(1767262881251901),
		"1767262881251901":    time.UnixMicro(1767262881251901),
		"1767262881251901234": time.Unix(0, 1767262881251901234),
	}
	for _, format := range []string{"unix", "unixmilli", "unixmicro", "unixnano"} {
		for raw, w := range want {
			for _, stringTag := range []bool{false, true} {
				got, err := decodeTagged(nil, format, stringTag, raw)
				if err != nil || !got.Equal(w) || got.Location() != time.UTC {
					t.Errorf("%s %s = %v, %v; want %v", format, raw, got, err, w.UTC())
				}
			}
		}
	}
	// Other widths, and negative values, keep the declared unit.
	for raw, w := range map[string]time.Time{"0": time.Unix(0, 0), "999999999": time.Unix(999999999, 0),
		"99999999999": time.Unix(99999999999, 0), "-1767262881251": time.Unix(-1767262881251, 0)} {
		if got, err := decodeTagged(nil, "unix", false, raw); err != nil || !got.Equal(w) {
			t.Errorf("unix %s = %v, %v; want %v", raw, got, err, w.UTC())
		}
	}
	if got, err := decodeTagged(nil, "unixmilli", false, "5000"); err != nil || !got.Equal(time.UnixMilli(5000)) {
		t.Errorf("unixmilli 5000 = %v, %v", got, err)
	}
}

// TestGateWireShapes pins the timestamp shapes Gate actually sends (captured
// from the live API) and their round trip.
func TestGateWireShapes(t *testing.T) {
	var v struct {
		Sec      time.Time `json:"sec,format:unix"`                   // futures create_time
		SecFrac  time.Time `json:"sec_frac,format:unix"`              // futures trade create_time_ms
		SecStr   time.Time `json:"sec_str,string,format:unix"`        // spot order create_time
		MS       time.Time `json:"ms,format:unixmilli"`               // server_time
		MSStr    time.Time `json:"ms_str,string,format:unixmilli"`    // crossex delist_time
		MSFrac   time.Time `json:"ms_frac,string,format:unixmilli"`   // spot trade create_time_ms
		Micro    time.Time `json:"micro,format:unix"`                 // price-triggered order cancel
		Epoch    time.Time `json:"epoch,format:unix"`                 // currency pair sell_start: 0
		EpochStr time.Time `json:"epoch_str,string,format:unixmilli"` // crossex delist_time: "0"
		Null     time.Time `json:"null,format:unix"`
	}
	const payload = `{"sec":1790646162,"sec_frac":1790646080.498,"sec_str":"1790646162","ms":1790646162123,` +
		`"ms_str":"1790646162123","ms_frac":"1790646118431.553000","micro":1767262881251901,"epoch":0,"epoch_str":"0","null":null}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	checks := []struct {
		name string
		got  time.Time
		want time.Time
	}{
		{"sec", v.Sec, time.Unix(1790646162, 0)},
		{"sec_frac", v.SecFrac, time.Unix(1790646080, 498_000_000)},
		{"sec_str", v.SecStr, time.Unix(1790646162, 0)},
		{"ms", v.MS, time.UnixMilli(1790646162123)},
		{"ms_str", v.MSStr, time.UnixMilli(1790646162123)},
		{"ms_frac", v.MSFrac, time.UnixMicro(1790646118431553)},
		{"micro", v.Micro, time.UnixMicro(1767262881251901)},
		{"epoch", v.Epoch, time.Unix(0, 0)},
		{"epoch_str", v.EpochStr, time.Unix(0, 0)},
	}
	for _, c := range checks {
		if !c.got.Equal(c.want) || c.got.Location() != time.UTC {
			t.Errorf("%s = %v, want %v in UTC", c.name, c.got, c.want.UTC())
		}
	}
	if !v.Null.IsZero() {
		t.Errorf("null = %v, want zero time", v.Null)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// Encoding is the standard one: the declared unit and quoting, sub-unit
	// precision kept, trailing zeros trimmed, and the zero time as its
	// (negative) instant.
	const want = `{"sec":1790646162,"sec_frac":1790646080.498,"sec_str":"1790646162","ms":1790646162123,` +
		`"ms_str":"1790646162123","ms_frac":"1790646118431.553","micro":1767262881.251901,"epoch":0,"epoch_str":"0","null":-62135596800}`
	if string(out) != want {
		t.Errorf("marshal =\n%s\nwant\n%s", out, want)
	}
}

// TestQuotingIsLenient: the value is read quoted or bare whatever the field's
// ,string option says, but an empty string or a non-numeric value is still an
// error.
func TestQuotingIsLenient(t *testing.T) {
	var v struct {
		Bare   time.Time `json:"bare,format:unix"`
		Quoted time.Time `json:"quoted,string,format:unixmilli"`
	}
	if err := JSONUnmarshal([]byte(`{"bare":"1790646162","quoted":1790646162123}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.Bare.Unix() != 1790646162 || v.Quoted.UnixMilli() != 1790646162123 {
		t.Errorf("got %v %v", v.Bare, v.Quoted)
	}
	for _, in := range []string{`{"bare":""}`, `{"quoted":""}`, `{"bare":"abc"}`, `{"bare":true}`, `{"bare":{}}`, `{"bare":[1]}`, `{"bare":"1"}`} {
		err := JSONUnmarshal([]byte(in), &v)
		if (err == nil) != (in == `{"bare":"1"}`) {
			t.Errorf("unmarshal %s: err = %v", in, err)
		}
	}
}

// TestTimeNotSet pins how Gate's "not set" values decode: null is the zero
// time, 0 is the Unix epoch (the standard reading, unchanged from earlier
// releases), and an absent key leaves the field untouched.
func TestTimeNotSet(t *testing.T) {
	var v struct {
		Null  time.Time  `json:"null,format:unix"`
		Zero  time.Time  `json:"zero,format:unix"`
		ZeroS time.Time  `json:"zero_s,string,format:unixmilli"`
		P     *time.Time `json:"p,format:unixmilli"`
		Miss  time.Time  `json:"miss,format:unix"`
	}
	if err := JSONUnmarshal([]byte(`{"null":null,"zero":0,"zero_s":"0","p":null}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	epoch := time.Unix(0, 0)
	if !v.Null.IsZero() || !v.Zero.Equal(epoch) || !v.ZeroS.Equal(epoch) || v.P != nil || !v.Miss.IsZero() {
		t.Errorf("got %v %v %v %v %v", v.Null, v.Zero, v.ZeroS, v.P, v.Miss)
	}
}

func TestPointerFields(t *testing.T) {
	var v struct {
		P    *time.Time `json:"p,format:unixmilli"`
		PStr *time.Time `json:"p_str,string,format:unix"`
		Nil  *time.Time `json:"nil,format:unix"`
	}
	if err := JSONUnmarshal([]byte(`{"p":"1790646162123","p_str":1790646162,"nil":null}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.P == nil || v.P.UnixMilli() != 1790646162123 || v.PStr == nil || v.PStr.Unix() != 1790646162 || v.Nil != nil {
		t.Errorf("got %v %v %v", v.P, v.PStr, v.Nil)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if want := `{"p":1790646162123,"p_str":"1790646162","nil":null}`; string(out) != want {
		t.Errorf("marshal = %s, want %s", out, want)
	}
}

// TestFormatOfScope: times nested in untagged values must not inherit a
// sibling field's format.
func TestFormatOfScope(t *testing.T) {
	type inner struct {
		T time.Time `json:"t"`
	}
	var v struct {
		P  *time.Time           `json:"p,format:unixmilli"`
		In inner                `json:"in"`
		Ts []time.Time          `json:"ts"`
		M  map[string]time.Time `json:"m"`
		S  time.Time            `json:"s,format:unix"`
	}
	in := `{"p":"1790646162123","in":{"t":"2026-09-29T00:00:00Z"},"ts":["2026-09-29T00:00:00Z"],"m":{"a":"2026-09-29T00:00:00Z"},"s":1790646162}`
	if err := JSONUnmarshal([]byte(in), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	day := time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)
	if v.P == nil || v.P.UnixMilli() != 1790646162123 || !v.In.T.Equal(day) || !v.Ts[0].Equal(day) || !v.M["a"].Equal(day) || v.S.Unix() != 1790646162 {
		t.Errorf("got %v %v %v %v %v", v.P, v.In.T, v.Ts, v.M, v.S)
	}
	if err := JSONUnmarshal([]byte(`{"ts":[1790646162]}`), &v); err == nil {
		t.Errorf("unix number in an untagged []time.Time: want error")
	}
}

func TestLayoutFormats(t *testing.T) {
	var v struct {
		RFC time.Time `json:"rfc,format:RFC3339"`
		D   time.Time `json:"d,format:DateOnly"`
		DT  time.Time `json:"dt,format:DateTime"`
		C   time.Time `json:"c,format:'2006/01/02'"`
	}
	const payload = `{"rfc":"2026-09-29T08:00:00+08:00","d":"2026-09-29","dt":"2026-09-29 00:00:00","c":"2026/09/29"}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)
	if !v.RFC.Equal(want) || !v.D.Equal(want) || !v.DT.Equal(want) || !v.C.Equal(want) {
		t.Errorf("got %v %v %v %v", v.RFC, v.D, v.DT, v.C)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal = %s, want %s", out, payload)
	}
	if err := JSONUnmarshal([]byte(`{"d":20260929}`), &v); err == nil {
		t.Errorf("bare number into DateOnly: want error")
	}
	// Gate's account tier_expire_time, including the "not set" form.
	var tier struct {
		T time.Time `json:"tier_expire_time,format:RFC3339"`
	}
	if err := JSONUnmarshal([]byte(`{"tier_expire_time":"0001-01-01T00:00:00Z"}`), &tier); err != nil || !tier.T.IsZero() {
		t.Errorf("tier_expire_time = %v, %v", tier.T, err)
	}
}

// TestStandardFallback checks that fields without a format keep the standard
// RFC 3339 behaviour, and invalid formats are reported.
func TestStandardFallback(t *testing.T) {
	var v struct {
		Plain time.Time `json:"plain"`
	}
	if err := JSONUnmarshal([]byte(`{"plain":"2026-09-29T08:00:00Z"}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := JSONUnmarshal([]byte(`{"plain":1790646162}`), &v); err == nil {
		t.Errorf("untagged time.Time accepted a unix number; it must require an explicit format")
	}
	var bad struct {
		T time.Time `json:"t,format:bogus"`
	}
	if err := JSONUnmarshal([]byte(`{"t":1}`), &bad); err == nil || !strings.Contains(err.Error(), "format") {
		t.Errorf("invalid format unmarshal error = %v", err)
	}
	if _, err := JSONMarshal(bad); err == nil || !strings.Contains(err.Error(), "format") {
		t.Errorf("invalid format marshal error = %v", err)
	}
}

func TestOtherStandardFormats(t *testing.T) {
	var v struct {
		D  time.Duration     `json:"d,format:units"`
		DS time.Duration     `json:"ds,format:milli"`
		B  []byte            `json:"b,format:hex"`
		F  float64           `json:"f,format:nonfinite"`
		M  map[string]string `json:"m,format:emitnull"`
		S  []int             `json:"s,format:emitempty"`
	}
	const payload = `{"d":"1h30m0s","ds":5000,"b":"0102","f":"NaN","m":null,"s":[]}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.DS != 5*time.Second {
		t.Errorf("ds = %v, want 5s", v.DS)
	}
	v.S = nil
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal = %s, want %s", out, payload)
	}
}

// TestEncodeKeepsPrecision: encoding is the standard one, so sub-unit
// precision is written as a fraction rather than truncated (Gate itself sends
// fractional seconds and milliseconds).
func TestEncodeKeepsPrecision(t *testing.T) {
	var v struct {
		S  time.Time `json:"s,format:unix"`
		MS time.Time `json:"ms,string,format:unixmilli"`
	}
	v.S = time.Unix(1790646080, 498_123_456)
	v.MS = v.S
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if want := `{"s":1790646080.498123456,"ms":"1790646080498.123456"}`; string(out) != want {
		t.Errorf("marshal = %s, want %s", out, want)
	}
}

func TestDecimalCodec(t *testing.T) {
	var v struct {
		Quoted decimal.Decimal `json:"quoted"`
		Bare   decimal.Decimal `json:"bare"`
		Empty  decimal.Decimal `json:"empty"`
		Null   decimal.Decimal `json:"null"`
	}
	if err := JSONUnmarshal([]byte(`{"quoted":"65000.5","bare":0.001,"empty":"","null":null}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.Quoted.String() != "65000.5" || v.Bare.String() != "0.001" || !v.Empty.IsZero() || !v.Null.IsZero() {
		t.Errorf("got %v %v %v %v", v.Quoted, v.Bare, v.Empty, v.Null)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if want := `{"quoted":"65000.5","bare":"0.001","empty":"0","null":"0"}`; string(out) != want {
		t.Errorf("marshal = %s, want %s", out, want)
	}
	if err := JSONUnmarshal([]byte(`{"quoted":true}`), &v); err == nil {
		t.Errorf("bool into decimal: want error")
	}
}

// TestConcurrentUse exercises the shared options from many goroutines (run
// with -race).
func TestConcurrentUse(t *testing.T) {
	type rec struct {
		T time.Time       `json:"t,string,format:unixmilli"`
		S time.Time       `json:"s,format:unix"`
		D decimal.Decimal `json:"d"`
	}
	const payload = `{"t":"1790646118431.553","s":1790646080.498,"d":"1.5"}`
	var wg sync.WaitGroup
	for range 16 {
		wg.Go(func() {
			for range 200 {
				var v rec
				if err := JSONUnmarshal([]byte(payload), &v); err != nil {
					t.Error(err)
					return
				}
				out, err := JSONMarshal(v)
				if err != nil || string(out) != payload {
					t.Errorf("round trip = %s, %v", out, err)
					return
				}
			}
		})
	}
	wg.Wait()
}
