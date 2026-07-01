package request

import (
	"context"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"github.com/UnipayFI/go-gate/common"
	"github.com/UnipayFI/go-gate/pkg/log"
	"github.com/go-resty/resty/v2"
)

// Client is what every endpoint Service needs from a product REST client. All
// getters are read-only; the concrete *client.Client satisfies it.
type Client interface {
	GetHttpClient() *resty.Client
	GetAPIKey() string
	GetAPISecret() string
	GetLogger() log.Logger
	GetSignFn() SignFn
	TimestampMs() int64
}

type Request struct {
	client   Client
	r        *resty.Request
	method   string
	path     string
	query    url.Values
	bodyJSON string
	needSign bool
	err      error
}

func newRequest(ctx context.Context, c Client, method, path string) *Request {
	r := c.GetHttpClient().R().
		SetHeader("User-Agent", common.GO_GATE_USER_AGENT).
		SetHeader("Accept", "application/json").
		SetContext(ctx)
	r.Method = method
	return &Request{
		client: c,
		r:      r,
		method: method,
		path:   path,
		query:  url.Values{},
	}
}

// Get builds a GET request. Any params maps are merged into the (sorted,
// url-encoded) query string, which is also part of the signed prehash.
func Get(ctx context.Context, c Client, path string, params ...map[string]string) *Request {
	r := newRequest(ctx, c, http.MethodGet, path)
	r.setQuery(params...)
	return r
}

// Post builds a POST request. Any body maps are merged and JSON-encoded once;
// the exact bytes sent are the bytes hashed into the signature.
func Post(ctx context.Context, c Client, path string, body ...map[string]any) *Request {
	r := newRequest(ctx, c, http.MethodPost, path)
	r.setBody(mergeBody(body...))
	return r
}

// Put builds a PUT request (Gate uses it for order amendment and batch upserts).
func Put(ctx context.Context, c Client, path string, body ...map[string]any) *Request {
	r := newRequest(ctx, c, http.MethodPut, path)
	r.setBody(mergeBody(body...))
	return r
}

// Delete builds a DELETE request. Any params maps become the query string;
// some Gate delete endpoints also carry a JSON body (set via SetBody).
func Delete(ctx context.Context, c Client, path string, params ...map[string]string) *Request {
	r := newRequest(ctx, c, http.MethodDelete, path)
	r.setQuery(params...)
	return r
}

// Patch builds a PATCH request.
func Patch(ctx context.Context, c Client, path string, body ...map[string]any) *Request {
	r := newRequest(ctx, c, http.MethodPatch, path)
	r.setBody(mergeBody(body...))
	return r
}

func (r *Request) setQuery(params ...map[string]string) {
	for _, p := range params {
		for k, v := range p {
			if v == "" {
				continue
			}
			r.query.Set(k, v)
		}
	}
}

// SetQuery adds or overrides a single query parameter. Empty values are ignored.
func (r *Request) SetQuery(key, value string) *Request {
	if value != "" {
		r.query.Set(key, value)
	}
	return r
}

// SetBody overrides the request body with an arbitrary JSON-serializable value
// (used for batch endpoints whose body is an array or a nested struct rather
// than a flat map). The value is marshaled once and reused for signing.
func (r *Request) SetBody(body any) *Request {
	if r.err != nil {
		return r
	}
	data, err := common.JSONMarshal(body)
	if err != nil {
		r.err = err
		return r
	}
	r.bodyJSON = common.BytesToString(data)
	return r
}

func (r *Request) setBody(body map[string]any) {
	if len(body) == 0 {
		return
	}
	data, err := common.JSONMarshal(body)
	if err != nil {
		r.err = err
		return
	}
	r.bodyJSON = common.BytesToString(data)
}

func mergeBody(body ...map[string]any) map[string]any {
	merged := make(map[string]any)
	for _, b := range body {
		maps.Copy(merged, b)
	}
	return merged
}

// WithSign marks the request as private: the KEY / SIGN / Timestamp headers are
// attached at send time. Public market endpoints omit this.
func (r *Request) WithSign() *Request {
	r.needSign = true
	return r
}

// queryString returns the sorted, url-encoded query string (without "?"). Gate
// signs exactly the query string that is sent, so the same value feeds both the
// URL and the prehash.
func (r *Request) queryString() string {
	return r.query.Encode()
}

func (r *Request) fullURL() string {
	base := strings.TrimSuffix(r.client.GetHttpClient().BaseURL, "/")
	path := r.path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	urlStr := base + path
	if q := r.queryString(); q != "" {
		urlStr += "?" + q
	}
	return urlStr
}
