package common

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

// decodeTime decodes a time.Time field whose `format` is unix, unixmilli,
// unixmicro or unixnano exactly like the standard option (same parser, UTC
// result), plus two Gate quirks:
//
//   - The value may be quoted or bare whatever the field's ,string option says.
//     Gate's docs warn that time fields "may differ in formats (int64, number or
//     string)", and its docs and live API disagree on the quoting of some
//     fields.
//   - A non-negative value whose integer part has exactly 10, 13, 16 or 19
//     digits is read in seconds, milliseconds, microseconds or nanoseconds
//     respectively. Such values are only ever timestamps between 2001 and 2286
//     in that unit, while the declared unit would place them within weeks of
//     1970 or tens of thousands of years ahead. Gate does send them: cancelling
//     a futures price-triggered order reports create_time in microseconds while
//     every other endpoint uses seconds.
//
// Anything else (no format, RFC3339 and other layouts) is left to the
// standard library.
func decodeTime(dec *jsontext.Decoder, t *time.Time) error {
	pow10, ok := unixPow10(formatOf(dec.Options()))
	if !ok {
		return errors.ErrUnsupported
	}
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}
	var b []byte
	switch val.Kind() {
	case 'n': // null
		*t = time.Time{}
		return nil
	case '"': // quoted string
		if b, err = unquote(val); err != nil {
			return err
		}
	case '0': // bare number
		b = val
	default:
		return fmt.Errorf("gate: cannot decode JSON %v into a unix time", val.Kind())
	}
	if p, ok := digitsPow10(b); ok {
		pow10 = p
	}
	*t, err = parseTimeUnix(b, pow10)
	return err
}

// unixPow10 reports the units per second of a unix* `format` value.
func unixPow10(format string) (uint64, bool) {
	switch format {
	case "unix":
		return 1e0, true
	case "unixmilli":
		return 1e3, true
	case "unixmicro":
		return 1e6, true
	case "unixnano":
		return 1e9, true
	}
	return 0, false
}

// digitsPow10 maps the width of a non-negative timestamp's integer part to the
// unit whose timestamps between 2001-09-09 and 2286-11-20 have that width.
func digitsPow10(b []byte) (uint64, bool) {
	n := 0
	for n < len(b) && '0' <= b[n] && b[n] <= '9' {
		n++
	}
	if n < len(b) && b[n] != '.' {
		return 0, false // sign or garbage; left to parseTimeUnix
	}
	switch n {
	case 10:
		return 1e0, true
	case 13:
		return 1e3, true
	case 16:
		return 1e6, true
	case 19:
		return 1e9, true
	}
	return 0, false
}

// unquote returns the contents of a JSON string, avoiding a copy when it
// contains no escape sequences (timestamps never do).
func unquote(val jsontext.Value) ([]byte, error) {
	if b := val[1 : len(val)-1]; !containsByte(b, '\\') {
		return b, nil
	}
	return jsontext.AppendUnquote(nil, val)
}

func containsByte(b []byte, c byte) bool {
	for _, x := range b {
		if x == c {
			return true
		}
	}
	return false
}

// formatOf returns the `format` tag option of the struct field currently
// being marshaled or unmarshaled, or "" if it has none. Pass it the Options
// of the Encoder/Decoder given to a MarshalToFunc/UnmarshalFromFunc; they are
// only valid for the duration of that call.
//
// encoding/json/v2 has no public accessor for the format, so this reads the
// Format field of the options struct those methods return (a
// *jsonopts.Struct in Go 1.27). errFormatOf reports when that layout is not
// what this code expects, in which case formatOf always returns "".
//
// Format is only reset between struct fields, so it must only be consulted for
// the value a format-tagged field holds directly (time.Time or *time.Time), not
// for values nested inside a format-tagged type with its own JSON methods.
func formatOf(opts json.Options) string {
	if errFormatOf != nil || reflect.TypeOf(opts) != optionsType {
		return ""
	}
	// An interface value is a (type, data) word pair; data points at the struct.
	p := (*[2]unsafe.Pointer)(unsafe.Pointer(&opts))[1]
	return *(*string)(unsafe.Add(p, formatOffset))
}

var (
	optionsType  = reflect.TypeOf(new(jsontext.Decoder).Options())
	formatOffset uintptr
	errFormatOf  = func() error {
		if enc := reflect.TypeOf(new(jsontext.Encoder).Options()); enc != optionsType {
			return fmt.Errorf("gate: encoder options type %v differs from decoder options type %v", enc, optionsType)
		}
		if optionsType == nil || optionsType.Kind() != reflect.Pointer || optionsType.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("gate: unexpected encoding/json/v2 options type %v", optionsType)
		}
		sf, ok := optionsType.Elem().FieldByName("Format")
		if !ok || sf.Type.Kind() != reflect.String {
			return fmt.Errorf("gate: encoding/json/v2 options type %v has no Format string field", optionsType)
		}
		// FieldByName reports the offset within the innermost embedded
		// struct; sum the offsets along the embedding path.
		st := optionsType.Elem()
		for _, i := range sf.Index {
			if st.Kind() != reflect.Struct {
				return fmt.Errorf("gate: encoding/json/v2 options field Format is not embedded by value in %v", optionsType)
			}
			f := st.Field(i)
			formatOffset += f.Offset
			st = f.Type
		}
		return nil
	}()
)
