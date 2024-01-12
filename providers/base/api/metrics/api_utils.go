package metrics

import (
	"errors"

	providererrors "github.com/skip-mev/slinky/providers/base/api/errors"
)

const (
	// CreateURL indicates that the provider count not construct a valid url to query.
	CreateURL Status = iota
	// DoRequest indicates that the request handler could not make the request.
	DoRequest
	// ParseResponse indicates that the provider could not parse the response.
	ParseResponse
	// RateLimit indicates that the request handler encountered a rate limit.
	RateLimit
	// UnexpectedStatusCode indicates that the request handler encountered an unexpected
	// status code.
	UnexpectedStatusCode
	// Success indicates that the provider successfully queried the data.
	Success
	// Unknown indicates that the provider encountered an unknown error.
	Unknown
)

// Status is a type that represents the status of a provider response.
type Status int

// String returns a string representation of the status.
func (s Status) String() string {
	switch s {
	case CreateURL:
		return "create_url_err"
	case DoRequest:
		return "request_err"
	case ParseResponse:
		return "parse_response_err"
	case RateLimit:
		return "rate_limit_err"
	case UnexpectedStatusCode:
		return "unexpected_status_code_err"
	case Success:
		return "success"
	default:
		return "unknown_err"
	}
}

// StatusFromError returns a Status based on the error. If the error is nil, StatusSuccess is returned.
func StatusFromError(err error) Status {
	switch {
	case err == nil:
		return Success
	case errors.Is(err, providererrors.ErrCreateURL):
		return CreateURL
	case errors.Is(err, providererrors.ErrDoRequest):
		return DoRequest
	case errors.Is(err, providererrors.ErrParseResponse):
		return ParseResponse
	case errors.Is(err, providererrors.ErrRateLimit):
		return RateLimit
	case errors.Is(err, providererrors.ErrUnexpectedStatusCode):
		return UnexpectedStatusCode
	default:
		return Unknown
	}
}
