package types

import (
	"fmt"

	servicemetrics "github.com/skip-mev/slinky/service/metrics"
)

type NilRequestError struct {
	Handler servicemetrics.ABCIMethod
}

func (e NilRequestError) Error() string {
	return fmt.Sprintf("nil request for %s", e.Handler)
}

func (e NilRequestError) Label() string {
	return "NilRequestError"
}

type WrappedHandlerError struct {
	Handler servicemetrics.ABCIMethod
	Err     error
}

func (e WrappedHandlerError) Error() string {
	return fmt.Sprintf("wrapped %s failed: %s", e.Handler, e.Err.Error())
}

func (e WrappedHandlerError) Label() string {
	return "WrappedHandlerError"
}

type CodecError struct {
	Err error
}

func (e CodecError) Error() string {
	return fmt.Sprintf("codec error: %s", e.Err.Error())
}

func (e CodecError) Label() string {
	return "CodecError"
}
