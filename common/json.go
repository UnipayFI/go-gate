package common

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/shopspring/decimal"
)

// Gate encodes monetary amounts, prices, rates and ratios as JSON strings, and
// returns "" (occasionally null) when a value is not set. The stock shopspring
// decimal codec rejects the empty-string form, so we teach the JSON codec how to
// read/write decimals tolerantly, once, globally: every decimal.Decimal field in
// this SDK is a plain field with a plain json tag.
//
// Timestamps are deliberately NOT handled here. Gate mixes second- and
// millisecond-, integer- and string-encoded times across different fields
// (futures send integer unix seconds, spot/wallet send quoted strings,
// server_time is a number in ms), so each time.Time field carries its own
// go-json-experiment format tag — ,format:unix / ,format:unixmilli, plus the
// ,string option when the wire value is quoted. A global time.Time codec would
// override those per-field tags (verified), so there is none.
var (
	unmarshalers = json.WithUnmarshalers(json.UnmarshalFromFunc(decodeDecimal))
	marshalers   = json.WithMarshalers(json.MarshalToFunc(encodeDecimal))
)

// JSONMarshal marshals v with Gate's decimal-string convention applied.
func JSONMarshal(v any) ([]byte, error) {
	return json.Marshal(v, marshalers)
}

// JSONUnmarshal unmarshals data into v with Gate's decimal-string convention
// applied. time.Time fields still resolve via their own format tags.
func JSONUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v, unmarshalers)
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
