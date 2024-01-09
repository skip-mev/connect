package testutils

import (
	"bytes"
	"io"
	"net/http"
)

// CreateResponseFromJSON creates a http response from a json string.
func CreateResponseFromJSON(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}
