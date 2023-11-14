package coinmarketcap

import (
	"sync"
)

// getSymbolForPair returns the symbol for a currency pair.
func (p *Provider) getSymbolForTokenName(tokenName string) string {
	if symbol, ok := p.config.TokenNameToSymbol[tokenName]; ok {
		return symbol
	}

	return tokenName
}

// finish takes a wait-group, and returns a channel that is sent on when the
// Waitgroup is finished.
func finish(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})

	// non-blocing wait for waitgroup to finish, and return channel
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}
