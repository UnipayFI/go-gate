package client

import "fmt"

// APIError is the error body Gate returns on a failed request, paired with the
// HTTP status code. Gate has no success envelope: a 2xx response carries the
// endpoint payload directly, while a non-2xx response carries
// {"label":...,"message":...,"detail":...}.
type APIError struct {
	Status  int    `json:"-"`
	Label   string `json:"label"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error returns the status, label and message (and detail when present).
func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("<APIError> status=%d, label=%s, msg=%s, detail=%s", e.Status, e.Label, e.Message, e.Detail)
	}
	return fmt.Sprintf("<APIError> status=%d, label=%s, msg=%s", e.Status, e.Label, e.Message)
}

// IsAPIError reports whether err is a Gate *APIError.
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}
