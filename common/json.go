package common

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"time"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

// Gate encodes monetary amounts, prices, rates and ratios as JSON strings, and
// returns "" (occasionally null) when a value is not set. The stock shopspring
// decimal codec rejects the empty-string form, so we teach the JSON codec how to
// read/write decimals tolerantly, once, globally: every decimal.Decimal field in
// this SDK is a plain field with a plain json tag.
//
// Gate mixes second- and millisecond-, integer- and string-encoded times across
// different fields (futures send integer unix seconds, spot/wallet send quoted
// strings, server_time is a number in ms), so each time.Time field declares its
// wire format with the standard `format` tag option — ,format:unix /
// ,format:unixmilli, plus the ,string option when the wire value is quoted —
// which Go 1.27's encoding/json/v2 only honours when ExperimentalSupportFormatTag
// is set. Encoding follows that option exactly; decoding keeps its semantics and
// only adds the Gate quirks described on decodeTime.
var (
	unmarshalOptions json.Options
	marshalOptions   json.Options

	// errJSONSupport is non-nil when the running Go release no longer
	// provides the `format` tag hooks this codec relies on.
	errJSONSupport = initJSON()
)

// JSONMarshal marshals v with Gate's time and decimal conventions applied.
func JSONMarshal(v any) ([]byte, error) {
	if errJSONSupport != nil {
		return nil, errJSONSupport
	}
	return json.Marshal(v, marshalOptions)
}

// JSONUnmarshal unmarshals data into v with Gate's time and decimal
// conventions applied.
func JSONUnmarshal(data []byte, v any) error {
	if errJSONSupport != nil {
		return errJSONSupport
	}
	return json.Unmarshal(data, v, unmarshalOptions)
}

// initJSON builds the codec options and round-trips a probe through them, so
// a Go release that drops the experimental `format` tag support (by panicking
// on the unknown option or by rejecting the tag) or changes the options layout
// formatOf reads fails loudly instead of misdating fields.
func initJSON() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("gate: encoding/json/v2 format tag support unavailable: %v", r)
		}
	}()
	if errFormatOf != nil {
		return errFormatOf
	}
	formatTag := jsonexp.ExperimentalSupportFormatTag(true)
	unmarshalOptions = json.JoinOptions(formatTag, json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(decodeTime),
		json.UnmarshalFromFunc(decodeDecimal),
	)))
	marshalOptions = json.JoinOptions(formatTag, json.WithMarshalers(json.MarshalToFunc(encodeDecimal)))

	// The quoting is the opposite of what the tags declare and "us" is in
	// microseconds, so this only decodes if decodeTime sees a unix format on
	// each field.
	const payload = `{"s":"1790646080.498","ms":1790646118431.553,"us":1767262881251901}`
	var probe struct {
		S  time.Time `json:"s,format:unix"`
		MS time.Time `json:"ms,string,format:unixmilli"`
		US time.Time `json:"us,format:unix"`
	}
	if err := json.Unmarshal([]byte(payload), &probe, unmarshalOptions); err != nil {
		return fmt.Errorf("gate: encoding/json/v2 format tag support unavailable: %w", err)
	}
	if s, ms, us := probe.S.UnixMilli(), probe.MS.UnixMicro(), probe.US.UnixMicro(); s != 1790646080498 || ms != 1790646118431553 || us != 1767262881251901 {
		return fmt.Errorf("gate: encoding/json/v2 format tag probe decoded %d ms, %d µs and %d µs, want 1790646080498, 1790646118431553 and 1767262881251901", s, ms, us)
	}
	const want = `{"s":1790646080.498,"ms":"1790646118431.553","us":1767262881.251901}`
	if out, err := json.Marshal(probe, marshalOptions); err != nil || string(out) != want {
		return fmt.Errorf("gate: encoding/json/v2 format tag probe encoded %s (%v), want %s", out, err, want)
	}
	return nil
}

func decodeDecimal(dec *jsontext.Decoder, d *decimal.Decimal) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	var s string
	switch tok.Kind() {
	case 'n': // null
		*d = decimal.Zero
		return nil
	case '"': // quoted string
		s = tok.String()
	case '0': // bare number
		s = tok.String()
	default:
		return fmt.Errorf("gate: cannot decode %v token into decimal", tok.Kind())
	}
	if s == "" {
		*d = decimal.Zero
		return nil
	}
	v, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("gate: invalid decimal %q: %w", s, err)
	}
	*d = v
	return nil
}

func encodeDecimal(enc *jsontext.Encoder, d decimal.Decimal) error {
	return enc.WriteToken(jsontext.String(d.String()))
}
