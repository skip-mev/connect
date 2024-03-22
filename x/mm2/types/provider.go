package types

import (
	"fmt"

	"github.com/skip-mev/slinky/pkg/json"
)

// ValidateBasic performs basic validation on a ProviderConfig.
func (pc *ProviderConfig) ValidateBasic() error {
	if len(pc.Name) == 0 {
		return fmt.Errorf("provider name must not be empty")
	}

	if len(pc.OffChainTicker) == 0 {
		return fmt.Errorf("provider offchain ticker must not be empty")
	}

	// index is allowed to be empty

	if len(pc.Metadata_JSON) > MaxMetadataJSONFieldLength {
		return fmt.Errorf("metadata json field is longer than maximum length of %d", MaxMetadataJSONFieldLength)
	}

	return json.IsValid([]byte(pc.Metadata_JSON))
}

// Equal returns true iff the ProviderConfig is equal to the given ProviderConfig.
func (pc *ProviderConfig) Equal(other ProviderConfig) bool {
	if pc.Name != other.Name {
		return false
	}

	if pc.OffChainTicker != other.OffChainTicker {
		return false
	}

	if pc.Invert != other.Invert {
		return false
	}

	if pc.Index != other.Index {
		return false
	}

	return pc.Metadata_JSON == other.Metadata_JSON
}
