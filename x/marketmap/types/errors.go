package types

import "fmt"

// TickerAlreadyExistsError is an error indicating the given Ticker exists in state.
type TickerAlreadyExistsError struct {
	ticker TickerString
}

func NewTickerAlreadyExistsError(ticker TickerString) TickerAlreadyExistsError {
	return TickerAlreadyExistsError{ticker: ticker}
}

// Error returns the error string for TickerAlreadyExistsError.
func (e TickerAlreadyExistsError) Error() string {
	return fmt.Sprintf("market already exists for ticker %s", e.ticker)
}
