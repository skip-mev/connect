package types

import (
	"fmt"
)

func NewCurrencyPairNotExistError(cpID string) CurrencyPairNotExistError {
	return CurrencyPairNotExistError{cpID}
}

type CurrencyPairNotExistError struct {
	cp string
}

func (e CurrencyPairNotExistError) Error() string {
	return fmt.Sprintf("nonce is not stored for CurrencyPair: %s", e.cp)
}

func NewQuotePriceNotExistError(cpID string) QuotePriceNotExistError {
	return QuotePriceNotExistError{cpID}
}

type QuotePriceNotExistError struct {
	cp string
}

func (e QuotePriceNotExistError) Error() string {
	return fmt.Sprintf("no price updates for CurrencyPair: %s", e.cp)
}

type CurrencyPairAlreadyExistsError struct {
	cp string
}

func NewCurrencyPairAlreadyExistsError(cp CurrencyPair) CurrencyPairAlreadyExistsError {
	return CurrencyPairAlreadyExistsError{cp.String()}
}

func (e CurrencyPairAlreadyExistsError) Error() string {
	return fmt.Sprintf("currency pair already exists: %s", e.cp)
}
