package types

import (
	"errors"
	"fmt"
)

// ErrorCode is a type alias for an int error code.
type ErrorCode int

const (
	OK                         ErrorCode = 0
	ErrorRateLimitExceeded     ErrorCode = 1
	ErrorUnknown               ErrorCode = 2
	ErrorUnknownPair           ErrorCode = 3
	ErrorSubscriptionFailed    ErrorCode = 4
	ErrorNotSubscribed         ErrorCode = 5
	ErrorPingFailed            ErrorCode = 6
	ErrorPongFailed            ErrorCode = 7
	ErrorInvalidRequest        ErrorCode = 8
	ErrorInvalidArgument       ErrorCode = 9
	ErrorUnableToCreateURL     ErrorCode = 10
	ErrorWebsocketStartFail    ErrorCode = 11
	ErrorInvalidAPIChains      ErrorCode = 12
	ErrorNoResponse            ErrorCode = 13
	ErrorInvalidResponse       ErrorCode = 14
	ErrorInvalidChainID        ErrorCode = 15
	ErrorFailedToParsePrice    ErrorCode = 16
	ErrorInvalidWebSocketTopic ErrorCode = 17
	ErrorFailedToDecode        ErrorCode = 18
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
	case ErrorSubscriptionFailed:
		return errors.New("failed to make subscription")
	case ErrorNotSubscribed:
		return errors.New("not subscribed to object")
	case ErrorPingFailed:
		return errors.New("ping failed")
	case ErrorPongFailed:
		return errors.New("pong failed")
	case ErrorInvalidRequest:
		return errors.New("invalid body in request")
	case ErrorInvalidArgument:
		return errors.New("invalid argument in request")
	case ErrorFailedToDecode:
		return errors.New("failed to decode message")
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
		return errors.New("failed to start websocker connection")
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
	return fmt.Sprintf("%s: %s", ec.code.Error(), ec.internalErr.Error())
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
