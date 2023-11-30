package types

import (
	"fmt"
)

func NewCurrencyPairNotExistError(cp CurrencyPair) CurrencyPairNotExistError {
	return CurrencyPairNotExistError{cp.ToString()}
}

type CurrencyPairNotExistError struct {
	cp string
}

func (e CurrencyPairNotExistError) Error() string {
	return fmt.Sprintf("nonce is not stored for CurrencyPair: %s", e.cp)
}

func NewQuotePriceNotExistError(cp CurrencyPair) QuotePriceNotExistError {
	return QuotePriceNotExistError{cp.ToString()}
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
	return CurrencyPairAlreadyExistsError{cp.ToString()}
}

func (e CurrencyPairAlreadyExistsError) Error() string {
	return fmt.Sprintf("currency pair already exists: %s", e.cp)
}
