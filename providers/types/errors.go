package types

import (
	"errors"
)

// ErrorCode is a type alias for an int error code.
type ErrorCode int

const (
	OK                          ErrorCode = 0
	ErrorRateLimitExceeded      ErrorCode = 1
	ErrorUnknown                ErrorCode = 2
	ErrorUnknownPair            ErrorCode = 3
	ErrorUnableToCreateURL      ErrorCode = 4
	ErrorWebsocketStartFail     ErrorCode = 5
	ErrorInvalidAPIChains       ErrorCode = 6
	ErrorNoResponse             ErrorCode = 7
	ErrorInvalidResponse        ErrorCode = 8
	ErrorInvalidChainID         ErrorCode = 9
	ErrorFailedToParsePrice     ErrorCode = 10
	ErrorInvalidWebSocketTopic  ErrorCode = 11
	ErrorFailedToDecode         ErrorCode = 12
	ErrorAPIGeneral             ErrorCode = 13
	ErrorWebSocketGeneral       ErrorCode = 14
	ErrorGRPCGeneral            ErrorCode = 15
	ErrorNoExistingPrice        ErrorCode = 16
	ErrorTickerMetadataNotFound ErrorCode = 17
)

// Error returns the error representation of the ErrorCode.
func (e ErrorCode) Error() error {
	switch e {
	case OK:
		return nil
	case ErrorRateLimitExceeded:
		return errors.New("rate limit exceeded")
	case ErrorUnknownPair:
		return errors.New("unknown market pair")
	case ErrorFailedToDecode:
		return errors.New("failed to decode message")
	case ErrorAPIGeneral:
		return errors.New("general api error")
	case ErrorWebSocketGeneral:
		return errors.New("general websocket error")
	case ErrorInvalidWebSocketTopic:
		return errors.New("invalid websocket topic received")
	case ErrorFailedToParsePrice:
		return errors.New("failed to parse price")
	case ErrorInvalidChainID:
		return errors.New("invalid chain ID in response")
	case ErrorInvalidResponse:
		return errors.New("invalid response")
	case ErrorNoResponse:
		return errors.New("got no response")
	case ErrorInvalidAPIChains:
		return errors.New("invalid chains for api handler")
	case ErrorUnableToCreateURL:
		return errors.New("failed to create URL for request")
	case ErrorWebsocketStartFail:
		return errors.New("failed to start websocket connection")
	case ErrorGRPCGeneral:
		return errors.New("general grpc error")
	case ErrorNoExistingPrice:
		return errors.New("no existing price")
	case ErrorTickerMetadataNotFound:
		return errors.New("ticker metadata not found")
	case ErrorUnknown:
		fallthrough
	default:
		return errors.New("unknown error")
	}
}

type ErrorWithCode struct {
	code        ErrorCode
	internalErr error
}

// Error returns an error string wrapping the internalErr and the error code.
func (ec ErrorWithCode) Error() string {
	return ec.internalErr.Error()
}

// Code returns the internal ErrorCode.
func (ec ErrorWithCode) Code() ErrorCode {
	return ec.code
}

func NewErrorWithCode(err error, ec ErrorCode) ErrorWithCode {
	return ErrorWithCode{
		code:        ec,
		internalErr: err,
	}
}
