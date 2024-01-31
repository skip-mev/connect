package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrCreateURL is returned when the APIDataHandler cannot create a url.
	// This can occur if the provider does not have the necessary information
	// to create the url i.e. a malformed config.
	ErrCreateURL = errors.New("api data handler failed to create request")

	// ErrDoRequest is returned when the APIQueryHandler was unable to make
	// a request.
	ErrDoRequest = errors.New("api query handler failed to make the request")

	// ErrParseResponse is returned when the APIDataHandler cannot parse the response
	// from the request handler.
	ErrParseResponse = errors.New("api data handler failed to parse response")

	// ErrRateLimit is returned when the APIQueryHandler encounters a rate limit.
	ErrRateLimit = errors.New("api query handler encountered a rate limit")

	// ErrUnexpectedStatusCode is returned when the APIQueryHandler encounters an unexpected status code.
	ErrUnexpectedStatusCode = errors.New("api query handler encountered an unexpected status code")
)

// ErrCreateURLWithErr is used to create a new ErrCreateRequest with the given error.
// Provider's that implement the APIDataHandler interface should use this function to
// create the error.
func ErrCreateURLWithErr(err error) error {
	return errors.Join(ErrCreateURL, err)
}

// ErrDoRequestWithErr is used to create a new ErrDoRequest with the given error.
// Provider's that implement the APIQueryHandler interface should use this function to
// create the error.
func ErrDoRequestWithErr(err error) error {
	return errors.Join(ErrDoRequest, err)
}

// ErrParseResponseWithErr is used to create a new ErrParseResponse with the given error.
// Provider's that implement the APIDataHandler interface should use this function to
// create the error.
func ErrParseResponseWithErr(err error) error {
	return errors.Join(ErrParseResponse, err)
}

// ErrUnexpectedStatusCodeWithCode is used to create a new ErrUnexpectedStatusCode with the given code.
// Provider's that implement the APIQueryHandler interface should use this function to
// create the error.
func ErrUnexpectedStatusCodeWithCode(code int) error {
	return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, code)
}
