package metrics

const (
	// StatusLabel is a label for the status of a provider API response.
	StatusLabel = "internal_status"
	// StatusCodeLabel is a label for the status code of a provider API response.
	StatusCodeLabel = "status_code"
	// StatusCodeExactLabel is a label for the exact status code of a provider API response.
	StatusCodeExactLabel = "status_code_exact"
)

type (
	// RPCCode is the status code a RPC request.
	RPCCode string
)

const (
	// RPCCodeOK is the status code for a successful RPC request.
	RPCCodeOK RPCCode = "OK"
	// RPCCodeError is the status code for a failed RPC request.
	RPCCodeError RPCCode = "ERROR"
)
