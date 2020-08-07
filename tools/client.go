package tools

import (
	"net/http"
)

// HTTPRequester can make REST calls via HTTP protocol
type HTTPRequester interface {
	Do(request *http.Request) (*http.Response, error)
}

// NewRestRequester Creates a new http.Client
func NewRestRequester() HTTPRequester {
	return &http.Client{}
}

// IsErrorCode Check if the given http status is an error code
func IsErrorCode(code int) bool {
	return code < 200 || code > 299
}
