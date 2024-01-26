package proposals

import (
	"fmt"
)

// InvalidExtendedCommitInfoError is an error that is returned when a proposed ExtendedCommitInfo is invalid.
type InvalidExtendedCommitInfoError struct {
	Err error
}

func (e InvalidExtendedCommitInfoError) Error() string {
	return fmt.Sprintf("invalid extended commit info: %s", e.Err.Error())
}

func (e InvalidExtendedCommitInfoError) Label() string {
	return "InvalidExtendedCommitInfoError"
}

// MissingCommitInfoError is an error that is returned when a proposal is missing the CommitInfo from the previous
// height.
type MissingCommitInfoError struct{}

func (e MissingCommitInfoError) Error() string {
	return "missing commit info"
}

func (e MissingCommitInfoError) Label() string {
	return "MissingCommitInfoError"
}
