package ve

import (
	"fmt"
)

type PreBlockError struct {
	Err error
}

func (e PreBlockError) Error() string {
	return fmt.Sprintf("finalize block error: %s", e.Err.Error())
}

func (e PreBlockError) Label() string {
	return "PreBlockError"
}

type Panic struct {
	Err error
}

func (e Panic) Error() string {
	return fmt.Sprintf("panic: %s", e.Err.Error())
}

func (e Panic) Label() string {
	return "Panic"
}

type OracleClientError struct {
	Err error
}

func (e OracleClientError) Error() string {
	return fmt.Sprintf("oracle client error: %s", e.Err.Error())
}

func (e OracleClientError) Label() string {
	return "OracleClientError"
}

type TransformPricesError struct {
	Err error
}

func (e TransformPricesError) Error() string {
	return fmt.Sprintf("prices transform error: %s", e.Err.Error())
}

func (e TransformPricesError) Label() string {
	return "TransformPricesError"
}

type ValidateVoteExtensionError struct {
	Err error
}

func (e ValidateVoteExtensionError) Error() string {
	return fmt.Sprintf("validate vote extension error: %s", e.Err.Error())
}

func (e ValidateVoteExtensionError) Label() string {
	return "ValidateVoteExtensionError"
}
