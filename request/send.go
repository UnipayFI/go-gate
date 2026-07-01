package request

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/UnipayFI/go-gate/client"
	"github.com/UnipayFI/go-gate/common"
)

// Do executes the request and decodes the response body into *T. Gate has no
// success envelope, so the whole 2xx body is the payload. A non-2xx response is
// returned as a *client.APIError.
func Do[T any](r *Request) (resp *T, err error) {
	body, err := DoRaw(r)
	if err != nil {
		return nil, err
	}
	var out T
	if len(body) == 0 {
		return &out, nil
	}
	if uerr := common.JSONUnmarshal(body, &out); uerr != nil {
		return nil, fmt.Errorf("gate: decode response: %w (body: %s)", uerr, common.BytesToString(body))
	}
	return &out, nil
}

// DoRaw executes the request and returns the raw response body after verifying
// the HTTP status. Because Gate returns the payload directly, this raw body is
// exactly what tests diff against the typed structs.
func DoRaw(r *Request) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}

	r.client.GetLogger().Debugf("request: %s %s body=%s", r.method, r.r.URL, r.bodyJSON)
	response, err := r.r.Send()
	if err != nil {
		r.client.GetLogger().Errorf("request %s %s failed: %s", r.method, r.r.URL, err)
		return nil, err
	}
	body := response.Body()
	r.client.GetLogger().Debugf("response: %s", common.BytesToString(body))

	if response.StatusCode() >= 300 {
		return nil, parseAPIError(response.StatusCode(), body)
	}
	return body, nil
}

// prepare finalizes the URL, body and (when private) the KEY / SIGN / Timestamp
// signing headers. The signed prehash is
// method\npath\nquery\nhex(SHA512(body))\ntimestamp, using the exact bytes that
// go on the wire; the timestamp is whole seconds.
func (r *Request) prepare() error {
	r.r.URL = r.fullURL()
	r.r.Method = r.method
	if r.bodyJSON != "" {
		r.r.SetHeader("Content-Type", "application/json")
		r.r.SetBody(r.bodyJSON)
	}
	if !r.needSign {
		return nil
	}

	apiKey := r.client.GetAPIKey()
	secret := r.client.GetAPISecret()
	if apiKey == "" || secret == "" {
		return errors.New("missing credentials: configure client.WithAuth(apiKey, apiSecret)")
	}

	ts := strconv.FormatInt(r.client.TimestampMs()/1000, 10)
	// Gate signs the URL-UNESCAPED query string (e.g. "a=x,y"), while the wire
	// carries the percent-encoded form ("a=x%2Cy"). Mirror the official SDK:
	// unescape queryString() for the prehash, keep the encoded form in the URL.
	signedQuery, err := url.QueryUnescape(r.queryString())
	if err != nil {
		signedQuery = r.queryString()
	}
	prehash := r.method + "\n" + r.path + "\n" + signedQuery + "\n" + hashPayload(r.bodyJSON) + "\n" + ts

	var sign string
	if fn := r.client.GetSignFn(); fn != nil {
		sign, err = fn(secret, prehash)
	} else {
		sign, err = HMACSign(secret, prehash)
	}
	if err != nil {
		return err
	}

	r.r.SetHeader("KEY", apiKey)
	r.r.SetHeader("SIGN", sign)
	r.r.SetHeader("Timestamp", ts)
	return nil
}

// parseAPIError decodes a non-2xx body as a Gate error envelope, falling back to
// a raw-body error when it is not the expected {label,message} shape.
func parseAPIError(status int, body []byte) error {
	apiErr := &client.APIError{Status: status}
	if e := common.JSONUnmarshal(body, apiErr); e != nil || apiErr.Label == "" {
		return fmt.Errorf("request failed (status %d): %s", status, common.BytesToString(body))
	}
	return apiErr
}
