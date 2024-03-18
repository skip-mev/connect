package types

import (
	"fmt"
	"strings"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

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
	base := first.CurrencyPair.Base
	if first.Invert {
		base = first.CurrencyPair.Quote
	}

	last := p.Operations[len(p.Operations)-1]
	quote := last.CurrencyPair.Quote
	if last.Invert {
		quote = last.CurrencyPair.Base
	}

	return ticker == fmt.Sprintf("%s/%s", base, quote)
}

// GetCurrencyPairs returns the set of CurrencyPairs in the path. Note that some of the CurrencyPairs
// may need to be inverted. This function does NOT return the inverted CurrencyPairs.
func (p *Path) GetCurrencyPairs() []slinkytypes.CurrencyPair {
	cps := make([]slinkytypes.CurrencyPair, len(p.Operations))
	for i, op := range p.Operations {
		cps[i] = op.CurrencyPair
	}
	return cps
}

// ShowRoute returns the route of the path in human-readable format.
func (p *Path) ShowRoute() string {
	hops := make([]string, len(p.Operations))
	for i, op := range p.Operations {
		base := op.CurrencyPair.Base
		if op.Invert {
			base = op.CurrencyPair.Quote
		}

		quote := op.CurrencyPair.Quote
		if op.Invert {
			quote = op.CurrencyPair.Base
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

	if len(p.Operations) > MaxPathLength {
		return fmt.Errorf("length of path cannot be longer than %d, got %d", MaxPathLength, len(p.Operations))
	}

	quote := first.CurrencyPair.Quote
	if first.Invert {
		quote = first.CurrencyPair.Base
	}

	// Ensure that the path is a directed acyclic graph.
	seen := map[slinkytypes.CurrencyPair]struct{}{
		first.CurrencyPair: {},
	}
	for _, op := range p.Operations[1:] {
		if err := op.ValidateBasic(); err != nil {
			return err
		}

		if op.Invert {
			if _, ok := seen[op.CurrencyPair.Invert()]; ok {
				return fmt.Errorf("duplicated pair found")
			}
		}

		if _, ok := seen[op.CurrencyPair]; ok {
			return fmt.Errorf("path is not a directed acyclic graph")
		}
		seen[op.CurrencyPair] = struct{}{}

		switch {
		case !op.Invert && quote != op.CurrencyPair.Base:
			return fmt.Errorf("invalid path; expected %s, got %s", quote, op.CurrencyPair.Base)
		case !op.Invert && quote == op.CurrencyPair.Base:
			quote = op.CurrencyPair.Quote
		case op.Invert && quote != op.CurrencyPair.Quote:
			return fmt.Errorf("invalid path; expected %s, got %s", quote, op.CurrencyPair.Quote)
		case op.Invert && quote == op.CurrencyPair.Quote:
			quote = op.CurrencyPair.Base
		}
	}

	return nil
}

// NewOperation returns a new Operation instance. An Operation is a single step
// in a path that represents a conversion from one ticker to another. The operation's
// ticker is a price feed that is supported by a set of providers and may be inverted
// if necessary.
func NewOperation(cp slinkytypes.CurrencyPair, invert bool) (Operation, error) {
	o := Operation{
		CurrencyPair: cp,
		Invert:       invert,
	}

	if err := o.ValidateBasic(); err != nil {
		return Operation{}, err
	}

	return o, nil
}

// ValidateBasic performs basic validation on the Operation.
func (o *Operation) ValidateBasic() error {
	return o.CurrencyPair.ValidateBasic()
}

// ValidateBasic performs basic validation on the Paths.
func (p *Paths) ValidateBasic(cp slinkytypes.CurrencyPair) error {
	if len(p.Paths) == 0 {
		return fmt.Errorf("at least one path is required for a ticker to be calculated")
	}

	// Track the routes to ensure that there are no duplicates.
	routes := make(map[string]struct{})
	for _, path := range p.Paths {
		if err := path.ValidateBasic(); err != nil {
			return err
		}

		route := path.ShowRoute()
		if _, ok := routes[route]; ok {
			return fmt.Errorf("duplicate route found: %s", route)
		}
		routes[route] = struct{}{}

		// Ensure that the path ends up converting to the ticker.
		if !path.Match(cp.String()) {
			return fmt.Errorf("path does not match ticker")
		}
	}

	return nil
}

// UniqueTickers returns all unique tickers across all paths that
// are part of the Paths. This is particularly useful for determining the
// set of markets that are required for a given ticker as well as ensuring
// that a given set of providers can provide the required markets.
func (p *Paths) UniqueTickers() map[slinkytypes.CurrencyPair]struct{} {
	seen := make(map[slinkytypes.CurrencyPair]struct{})

	for _, path := range p.Paths {
		for _, cp := range path.GetCurrencyPairs() {
			seen[cp] = struct{}{}
		}
	}

	return seen
}

// Equal returns true iff the Paths is equal to the given Paths.
func (p *Paths) Equal(other Paths) bool {
	if len(p.Paths) != len(other.Paths) {
		return false
	}

	for i, path := range p.Paths {
		if !path.Equal(other.Paths[i]) {
			return false
		}
	}

	return true
}

// Equal returns true iff the Path is equal to the given Path.
func (p *Path) Equal(other Path) bool {
	if len(p.Operations) != len(other.Operations) {
		return false
	}

	for i, op := range p.Operations {
		if !op.Equal(other.Operations[i]) {
			return false
		}
	}

	return true
}

// Equal returns true iff the Operation is equal to the given Operation.
func (o *Operation) Equal(other Operation) bool {
	return o.CurrencyPair.Equal(other.CurrencyPair) &&
		o.Invert == other.Invert && o.Provider == other.Provider
}
