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

// Equal returns true iff the Providers is equal to the given Providers.
func (p *Providers) Equal(other Providers) bool {
	if len(p.Providers) != len(other.Providers) {
		return false
	}

	for i, provider := range p.Providers {
		if !provider.Equal(other.Providers[i]) {
			return false
		}
	}

	return true
}

// Equal returns true iff the ProviderConfig is equal to the given ProviderConfig.
func (pc *ProviderConfig) Equal(other ProviderConfig) bool {
	return pc.Name == other.Name &&
		pc.OffChainTicker == other.OffChainTicker
}
