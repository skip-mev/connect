package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	connecthttp "github.com/skip-mev/connect/v2/pkg/http"
)

func TestRoundTripperWithHeaders(t *testing.T) {
	expectedHeaderFields := map[string]string{
		"X-Api-Key": "test",
	}

	rt := &customRoundTripper{
		expectedHeaderFields: expectedHeaderFields,
	}

	rtWithHeaders := connecthttp.NewRoundTripperWithHeaders(rt, connecthttp.WithAuthentication("X-Api-Key", "test"))

	client := &http.Client{
		Transport: rtWithHeaders,
	}

	req, err := http.NewRequest(http.MethodGet, "http://test.com", nil)
	require.NoError(t, err)

	// Make the request
	_, err = client.Do(req)
	require.NoError(t, err)
}

type customRoundTripper struct {
	expectedHeaderFields map[string]string
}

func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range c.expectedHeaderFields {
		if req.Header.Get(k) != v {
			return nil, fmt.Errorf("expected header %s to be %s, got %s", k, v, req.Header.Get(k))
		}
	}
	return &http.Response{}, nil
}
