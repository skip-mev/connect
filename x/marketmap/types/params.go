package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// DefaultVersion is the default value for the Version Param.
	DefaultVersion = 0
)

// DefaultParams returns default marketmap parameters.
func DefaultParams() Params {
	return Params{
		MarketAuthority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Version:         DefaultVersion,
	}
}

// NewParams returns a new Params instance.
func NewParams(authority string, version uint64) Params {
	return Params{
		MarketAuthority: authority,
		Version:         version,
	}
}

// ValidateBasic performs stateless validation of the Params.
func (p *Params) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.MarketAuthority); err != nil {
		return fmt.Errorf("invalid market authority string: %w", err)
	}

	return nil
}
