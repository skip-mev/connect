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

