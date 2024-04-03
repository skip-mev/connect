package http

import (
	"net/http"
)

// IsValidHTTPResponse returns true if the response is a valid http response, false otherwise.
func IsValidHTTPResponse(resp *http.Response) bool {
	return resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300
}
