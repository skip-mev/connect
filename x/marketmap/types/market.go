package types

import (
	"fmt"
)

// Providers is a type alias for an array of ProviderConfig objects.
type Providers []ProviderConfig

// ValidateBasic performs basic validation on the MarketConfig.
func (p Providers) ValidateBasic() error {
	seen := make(map[string]struct{})
	for i := range p {
		providerConfig := p[i]

		if err := providerConfig.ValidateBasic(); err != nil {
			return err
		}

		// Check for duplicate providers.
		if _, ok := seen[providerConfig.Name]; ok {
			return fmt.Errorf("duplicate provider found: %s", providerConfig.Name)
		}
		seen[providerConfig.Name] = struct{}{}
	}

	return nil
}

// ValidateBasic performs basic validation on a ProviderConfig
func (pc *ProviderConfig) ValidateBasic() error {
	if len(pc.Name) == 0 {
		return fmt.Errorf("provider name must not be empty")
	}

	if len(pc.OffChainTicker) == 0 {
		return fmt.Errorf("provider offchain ticker must not be empty")
	}

	return nil
}
