package metrics

import "fmt"

const (
	// StatusLabel is a label for the status of a provider API response.
	StatusLabel = "internal_status"
	// StatusCodeLabel is a label for the status code of a provider API response.
	StatusCodeLabel = "status_code"
	// StatusCodeExactLabel is a label for the exact status code of a provider API response.
	StatusCodeExactLabel = "status_code_exact"
	// EndpointLabel is a label for the endpoint of a provider API response.
	EndpointLabel = "endpoint"
	// RedactedURL is a label for the redacted URL of a provider API response.
	RedactedURL = "redacted_url"
)

type (
	// RPCCode is the status code a RPC request.
	RPCCode string
)

const (
	// RPCCodeOK is the status code for a successful RPC request.
	RPCCodeOK RPCCode = "ok"
	// RPCCodeError is the status code for a failed RPC request.
	RPCCodeError RPCCode = "request_error"
)

// RedactedEndpointURL returns a redacted version of the given URL.
func RedactedEndpointURL(index int) string {
	return fmt.Sprintf("redacted_endpoint_index=%d", index)
}
