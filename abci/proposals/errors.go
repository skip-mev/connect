package proposals

import (
	"fmt"
)

type ABCIHandlerMethod int

const (
	PrepareProposal ABCIHandlerMethod = iota
	ProcessProposal
)

func (a ABCIHandlerMethod) String() string {
	switch a {
	case PrepareProposal:
		return "prepare_proposal"
	case ProcessProposal:
		return "process_proposal"
	default:
		return "not_implemented"
	}
}

type NilRequestError struct {
	Handler ABCIHandlerMethod
}

func (e NilRequestError) Error() string {
	return fmt.Sprintf("ABCIHandler: %s received a nil request", e.Handler)
}

func (e NilRequestError) Label() string {
	return "NilRequestError"
}

type WrappedHandlerError struct {
	Handler ABCIHandlerMethod
	Err     error
}

func (e WrappedHandlerError) Error() string {
	return fmt.Sprintf("wrapped %s failed: %s", e.Handler, e.Err.Error())
}

func (e WrappedHandlerError) Label() string {
	return "WrappedHandlerError"
}

type InvalidExtendedCommitInfoError struct {
	Err error
}

func (e InvalidExtendedCommitInfoError) Error() string {
	return fmt.Sprintf("invalid extended commit info: %s", e.Err.Error())
}

func (e InvalidExtendedCommitInfoError) Label() string {
	return "InvalidExtendedCommitInfoError"
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

type MissingCommitInfoError struct{}

func (e MissingCommitInfoError) Error() string {
	return "missing commit info"
}

func (e MissingCommitInfoError) Label() string {
	return "MissingCommitInfoError"
}
