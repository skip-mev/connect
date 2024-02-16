package types

import (
	fmt "fmt"
	"strings"
)

// NewPathsConfig returns a new PathsConfig instance. PathsConfig represents
// the list of convertable markets (paths) that will be used to convert the
// prices of a set of tickers to a common ticker.
//
// For example, if the oracle receives a price for BTC/USDT and USDT/USD, one
// possible path to get the price of BTC/USD would would be BTC/USDT -> USDT/USD.
func NewPathsConfig(ticker Ticker, paths ...Path) (PathsConfig, error) {
	c := PathsConfig{
		Ticker: ticker,
		Paths:  paths,
	}

	if err := c.ValidateBasic(); err != nil {
		return PathsConfig{}, err
	}

	return c, nil
}

// ValidateBasic performs basic validation on the PathsConfig.
func (c *PathsConfig) ValidateBasic() error {
	if err := c.Ticker.ValidateBasic(); err != nil {
		return err
	}

	if len(c.Paths) == 0 {
		return fmt.Errorf("at least one path is required for a ticker to be calculated")
	}

	// Track the routes to ensure that there are no duplicates.
	routes := make(map[string]struct{})
	for _, path := range c.Paths {
		if err := path.ValidateBasic(); err != nil {
			return err
		}

		route := path.ShowRoute()
		if _, ok := routes[route]; ok {
			return fmt.Errorf("duplicate path found: %s", route)
		}
		routes[route] = struct{}{}

		// Ensure that the path ends up converting to the ticker.
		if !path.Match(c.Ticker.String()) {
			return fmt.Errorf("path does not match ticker")
		}
	}

	return nil
}

// UniqueTickers returns all unique tickers across all paths that
// are part of the PathsConfig. This is particularly useful for determining the
// set of markets that are required for a given ticker as well as ensuring
// that a given set of providers can provide the required markets.
func (c *PathsConfig) UniqueTickers() map[Ticker]struct{} {
	seen := make(map[Ticker]struct{})

	for _, path := range c.Paths {
		for _, ticker := range path.GetTickers() {
			seen[ticker] = struct{}{}
		}
	}

	return seen
}

// NewPath returns a new Path instance. A Path is a list of convertable markets
// that will be used to convert the prices of a set of tickers to a common ticker.
func NewPath(ops ...Operation) (Path, error) {
	p := Path{
		Operations: ops,
	}

	if err := p.ValidateBasic(); err != nil {
		return Path{}, err
	}

	return p, nil
}

// Match returns true if the path matches the provided ticker.
func (p *Path) Match(ticker string) bool {
	if len(p.Operations) == 0 {
		return false
	}

	first := p.Operations[0]
	base := first.Ticker.CurrencyPair.Base
	if first.Invert {
		base = first.Ticker.CurrencyPair.Quote
	}

	last := p.Operations[len(p.Operations)-1]
	quote := last.Ticker.CurrencyPair.Quote
	if last.Invert {
		quote = last.Ticker.CurrencyPair.Base
	}

	return ticker == fmt.Sprintf("%s/%s", base, quote)
}

// GetTickers returns the set of tickers in the path. Note that some of the tickers
// may need to be inverted. This function does NOT return the inverted tickers.
func (p *Path) GetTickers() []Ticker {
	tickers := make([]Ticker, len(p.Operations))
	for i, op := range p.Operations {
		tickers[i] = op.Ticker
	}
	return tickers
}

// ShowRoute returns the route of the path in human-readable format.
func (p *Path) ShowRoute() string {
	hops := make([]string, len(p.Operations))
	for i, op := range p.Operations {
		base := op.Ticker.CurrencyPair.Base
		if op.Invert {
			base = op.Ticker.CurrencyPair.Quote
		}

		quote := op.Ticker.CurrencyPair.Quote
		if op.Invert {
			quote = op.Ticker.CurrencyPair.Base
		}

		hops[i] = fmt.Sprintf("%s/%s", base, quote)
	}

	return strings.Join(hops, " -> ")
}

// ValidateBasic performs basic validation on the Path. Specifically this will check
// that order is topologically sorted for each market. For example, if the oracle
// receives a price for BTC/USDT and USDT/USD, the order must be BTC/USDT -> USDT/USD.
// Alternatively, if the oracle receives a price for BTC/USDT and USD/USDT, the order
// must be BTC/USDT -> USD/USDT (inverted == true).
func (p *Path) ValidateBasic() error {
	if len(p.Operations) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	first := p.Operations[0]
	if err := first.ValidateBasic(); err != nil {
		return err
	}

	if len(p.Operations) == 1 {
		return nil
	}

	quote := first.Ticker.CurrencyPair.Quote
	if first.Invert {
		quote = first.Ticker.CurrencyPair.Base
	}

	// Ensure that the path is a directed acyclic graph.
	seen := map[Ticker]struct{}{
		first.Ticker: {},
	}
	for _, op := range p.Operations[1:] {
		if err := op.ValidateBasic(); err != nil {
			return err
		}

		if _, ok := seen[op.Ticker]; ok {
			return fmt.Errorf("path is not a directed acyclic graph")
		}
		seen[op.Ticker] = struct{}{}

		switch {
		case !op.Invert && quote != op.Ticker.CurrencyPair.Base:
			return fmt.Errorf("invalid path; expected %s, got %s", quote, op.Ticker.CurrencyPair.Base)
		case !op.Invert && quote == op.Ticker.CurrencyPair.Base:
			quote = op.Ticker.CurrencyPair.Quote
		case op.Invert && quote != op.Ticker.CurrencyPair.Quote:
			return fmt.Errorf("invalid path; expected %s, got %s", quote, op.Ticker.CurrencyPair.Quote)
		case op.Invert && quote == op.Ticker.CurrencyPair.Quote:
			quote = op.Ticker.CurrencyPair.Base
		}
	}

	return nil
}

// NewOperation returns a new Operation instance. An Operation is a single step
// in a path that represents a conversion from one ticker to another. The operation's
// ticker is a price feed that is supported by a set of providers and may be inverted
// if necessary.
func NewOperation(ticker Ticker, invert bool) (Operation, error) {
	o := Operation{
		Ticker: ticker,
		Invert: invert,
	}

	if err := o.ValidateBasic(); err != nil {
		return Operation{}, err
	}

	return o, nil
}

// ValidateBasic performs basic validation on the Operation.
func (o *Operation) ValidateBasic() error {
	return o.Ticker.ValidateBasic()
}
