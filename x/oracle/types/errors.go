package types

import (
	"fmt"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
)

func NewCurrencyPairNotExistError(cp connecttypes.CurrencyPair) CurrencyPairNotExistError {
	return CurrencyPairNotExistError{cp.String()}
}

type CurrencyPairNotExistError struct {
	cp string
}

func (e CurrencyPairNotExistError) Error() string {
	return fmt.Sprintf("nonce is not stored for CurrencyPair: %s", e.cp)
}

func NewQuotePriceNotExistError(cp connecttypes.CurrencyPair) QuotePriceNotExistError {
	return QuotePriceNotExistError{cp.String()}
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

func NewCurrencyPairAlreadyExistsError(cp connecttypes.CurrencyPair) CurrencyPairAlreadyExistsError {
	return CurrencyPairAlreadyExistsError{cp.String()}
}

func (e CurrencyPairAlreadyExistsError) Error() string {
	return fmt.Sprintf("currency pair already exists: %s", e.cp)
}
