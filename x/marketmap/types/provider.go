package types

import (
	"fmt"
)

// ValidateBasic performs basic validation on Providers.
func (p *Providers) ValidateBasic() error {
	for _, provider := range p.Providers {
		if err := provider.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic performs basic validation on a ProviderConfig.
func (pc *ProviderConfig) ValidateBasic() error {
	if len(pc.Name) == 0 {
		return fmt.Errorf("provider name must not be empty")
	}

	if len(pc.OffChainTicker) == 0 {
		return fmt.Errorf("provider offchain ticker must not be empty")
	}

	return nil
}
