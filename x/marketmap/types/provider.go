package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	MaxProviderNameFieldLength   = 128
	MaxProviderTickerFieldLength = 256
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

	if len(pc.Name) > MaxProviderNameFieldLength {
		return fmt.Errorf("provider length is longer than maximum length of %d", MaxProviderNameFieldLength)
	}

	if len(pc.OffChainTicker) == 0 {
		return fmt.Errorf("provider offchain ticker must not be empty")
	}

	if len(pc.OffChainTicker) > MaxProviderTickerFieldLength {
		return fmt.Errorf("provider offchain ticker is longer than maximum length of %d", MaxProviderTickerFieldLength)
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

func (p Providers) MarshalJSON() ([]byte, error) {
	slices.SortFunc(p.Providers, func(i, j ProviderConfig) int {
		return strings.Compare(i.Name, j.Name)
	})
	b := &bytes.Buffer{}
	_, err := b.WriteString("{\"providers\":[")
	if err != nil {
		return nil, err
	}
	for i, providerConfig := range p.Providers {
		pConfigBytes, err := json.Marshal(providerConfig)
		if err != nil {
			return nil, err
		}
		_, err = b.Write(pConfigBytes)
		if err != nil {
			return nil, err
		}
		if i < len(p.Providers)-1 {
			_, err = b.Write([]byte{','})
			if err != nil {
				return nil, err
			}
		}
	}
	_, err = b.WriteString("]}")
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
