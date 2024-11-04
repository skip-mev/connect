package types

import "fmt"

// MarketAlreadyExistsError is an error indicating the given Market exists in state.
type MarketAlreadyExistsError struct {
	ticker TickerString
}

func NewMarketAlreadyExistsError(ticker TickerString) MarketAlreadyExistsError {
	return MarketAlreadyExistsError{ticker: ticker}
}

// Error returns the error string for MarketAlreadyExistsError.
func (e MarketAlreadyExistsError) Error() string {
	return fmt.Sprintf("market already exists for ticker %s", e.ticker)
}

// MarketDoesNotExistsError is an error indicating the given Market does not exist in state.
type MarketDoesNotExistsError struct {
	ticker TickerString
}

func NewMarketDoesNotExistsError(ticker TickerString) MarketDoesNotExistsError {
	return MarketDoesNotExistsError{ticker: ticker}
}

// Error returns the error string for MarketDoesNotExistsError.
func (e MarketDoesNotExistsError) Error() string {
	return fmt.Sprintf("market does not exist for ticker %s", e.ticker)
}

// MarketIsEnabledError is an error indicating the given Market does not exist in state.
type MarketIsEnabledError struct {
	ticker TickerString
}

func NewMarketIsEnabledError(ticker TickerString) MarketIsEnabledError {
	return MarketIsEnabledError{ticker: ticker}
}

// Error returns the error string for MarketIsEnabledError.
func (e MarketIsEnabledError) Error() string {
	return fmt.Sprintf("market is currently enabled %s", e.ticker)
}

// CannotChangeDecimalOfEnabledMarketError is an error indicating that the update to a
// given enabled market is trying to change the Decimals value.
type CannotChangeDecimalOfEnabledMarketError struct {
	ticker TickerString
}

func NewCannotChangeDecimalOfEnabledMarketError(ticker TickerString) CannotChangeDecimalOfEnabledMarketError {
	return CannotChangeDecimalOfEnabledMarketError{ticker: ticker}
}

// Error returns the error string for CannotChangeDecimalOfEnabledMarketError.
func (e CannotChangeDecimalOfEnabledMarketError) Error() string {
	return fmt.Sprintf("Decimals value cannot be changed for an enabled market: %s", e.ticker)
}
