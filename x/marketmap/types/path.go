package types

import (
	fmt "fmt"
	"strings"
)

// NewPathsConfig returns a new PathsConfig instance. Given a set of paths, this will construct a new
// PathsConfig instance.
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
		return fmt.Errorf("paths cannot be empty")
	}

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

		if !path.Match(c.Ticker.String()) {
			return fmt.Errorf("path does not match ticker")
		}
	}

	return nil
}

// NewPath returns a new Path instance. This constructs a new path from a list of operations.
// The set of operations are valid if they form a directed acyclic graph.
func NewPath(ops ...Operation) (Path, error) {
	p := Path{
		Operations: ops,
	}

	if err := p.ValidateBasic(); err != nil {
		return Path{}, err
	}

	return p, nil
}

// Match returns true if the path matches the provided ticker. This is useful for determining
// if a path is valid for a given market.
func (p *Path) Match(ticker string) bool {
	if len(p.Operations) == 0 {
		return false
	}

	first := p.Operations[0]
	base := first.Ticker.Base
	if first.Invert {
		base = first.Ticker.Quote
	}

	last := p.Operations[len(p.Operations)-1]
	quote := last.Ticker.Quote
	if last.Invert {
		quote = last.Ticker.Base
	}

	return fmt.Sprintf("%s/%s", base, quote) == ticker
}

// ValidateBasic performs basic validation on the Path. Specifically this will check that order
// is topologically sorted for each market. For example, if the oracle receives a price for
// BTC/USDT and USDT/USD, the order must be BTC/USDT -> USDT/USD. Alternatively, if the oracle
// receives a price for BTC/USDT and USD/USDT, the order must be BTC/USDT -> USD/USDT (inverted == true).
func (p *Path) ValidateBasic() error {
	if len(p.Operations) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	first := p.Operations[0]
	if err := first.ValidateBasic(); err != nil {
		return err
	}

	quote := first.Ticker.Quote
	if first.Invert {
		quote = first.Ticker.Base
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
		case !op.Invert && quote != op.Ticker.Base:
			return fmt.Errorf("invalid path; expected %s, got %s", quote, op.Ticker.Base)
		case !op.Invert && quote == op.Ticker.Base:
			quote = op.Ticker.Quote
		case op.Invert && quote != op.Ticker.Quote:
			return fmt.Errorf("invalid path; expected %s, got %s", quote, op.Ticker.Quote)
		case op.Invert && quote == op.Ticker.Quote:
			quote = op.Ticker.Base
		}
	}

	return nil
}

// ShowRoute returns the route of the path in human readable format.
func (p *Path) ShowRoute() string {
	hops := make([]string, len(p.Operations))
	for i, op := range p.Operations {
		base := op.Ticker.Base
		if op.Invert {
			base = op.Ticker.Quote
		}

		quote := op.Ticker.Quote
		if op.Invert {
			quote = op.Ticker.Base
		}

		hops[i] = fmt.Sprintf("%s/%s", base, quote)
	}

	return strings.Join(hops, " -> ")
}

// NewOperation returns a new Operation instance.
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
	if err := o.Ticker.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
