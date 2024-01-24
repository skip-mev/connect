package proposals

import (
	"fmt"
)

type InvalidExtendedCommitInfoError struct {
	Err error
}

func (e InvalidExtendedCommitInfoError) Error() string {
	return fmt.Sprintf("invalid extended commit info: %s", e.Err.Error())
}

func (e InvalidExtendedCommitInfoError) Label() string {
	return "InvalidExtendedCommitInfoError"
}

type MissingCommitInfoError struct{}

func (e MissingCommitInfoError) Error() string {
	return "missing commit info"
}

func (e MissingCommitInfoError) Label() string {
	return "MissingCommitInfoError"
}
