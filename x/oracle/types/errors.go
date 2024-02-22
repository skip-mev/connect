package types

import (
	"fmt"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

func NewCurrencyPairNotExistError(cp slinkytypes.CurrencyPair) CurrencyPairNotExistError {
	return CurrencyPairNotExistError{cp.String()}
}

type CurrencyPairNotExistError struct {
	cp string
}

func (e CurrencyPairNotExistError) Error() string {
	return fmt.Sprintf("nonce is not stored for CurrencyPair: %s", e.cp)
}

func NewQuotePriceNotExistError(cp slinkytypes.CurrencyPair) QuotePriceNotExistError {
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

func NewCurrencyPairAlreadyExistsError(cp slinkytypes.CurrencyPair) CurrencyPairAlreadyExistsError {
	return CurrencyPairAlreadyExistsError{cp.String()}
}

func (e CurrencyPairAlreadyExistsError) Error() string {
	return fmt.Sprintf("currency pair already exists: %s", e.cp)
}
