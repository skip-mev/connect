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
		MarketAuthorities: []string{authtypes.NewModuleAddress(govtypes.ModuleName).String()},
		Admin:             authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	}
}

// NewParams returns a new Params instance.
func NewParams(authorities []string, admin string) (Params, error) {
	if authorities == nil {
		return Params{}, fmt.Errorf("cannot create Params with nil authority")
	}

	return Params{
		MarketAuthorities: authorities,
		Admin:             admin,
	}, nil
}

// ValidateBasic performs stateless validation of the Params.
func (p *Params) ValidateBasic() error {
	if p.MarketAuthorities == nil {
		return fmt.Errorf("cannot create Params with empty market authorities")
	}

	seenAuthorities := make(map[string]struct{}, len(p.MarketAuthorities))
	for _, authority := range p.MarketAuthorities {
		if _, seen := seenAuthorities[authority]; seen {
			return fmt.Errorf("duplicate authority %s found", authority)
		}

		if _, err := sdk.AccAddressFromBech32(authority); err != nil {
			return fmt.Errorf("invalid market authority string: %w", err)
		}

		seenAuthorities[authority] = struct{}{}
	}

	if _, err := sdk.AccAddressFromBech32(p.Admin); err != nil {
		return fmt.Errorf("invalid marketmap admin string: %w", err)
	}

	return nil
}
