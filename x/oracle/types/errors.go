package types

import (
	"fmt"
)

func NewCurrencyPairNotExistError(cp string) *CurrencyPairNotExistError {
	return &CurrencyPairNotExistError{cp}
}

type CurrencyPairNotExistError struct {
	cp string
}

func (e CurrencyPairNotExistError) Error() string {
	return fmt.Sprintf("nonce is not stored for CurrencyPair: %s", e.cp)
}
