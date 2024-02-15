package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// DefaultParams returns default marketmap parameters.
func DefaultParams() Params {
	return Params{
		MarketAuthority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	}
}

// NewParams returns a new Params instance.
func NewParams(authority string) Params {
	return Params{
		MarketAuthority: authority,
	}
}

// ValidateBasic performs stateless validation of the Params.
func (p *Params) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.MarketAuthority); err != nil {
		return fmt.Errorf("invalid market authority string: %w", err)
	}

	return nil
}
