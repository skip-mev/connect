package types

import (
	"fmt"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// DefaultEnabled is the default value for the enabled flag.
	DefaultEnabled = true

	// DefaultPruningEnabled is the default value for the pruning enabled flag.
	DefaultPruningEnabled = true

	// DefaultBondAmount is the default value for the bond amount.
	DefaultBondAmount = math.NewInt(1000000)

	// DefaultAlertExpiry is the maximum age of an alert.
	DefaultAlertExpiry = uint64(10000)

	// DefaultBlocksToPrune is the default number of blocks an alert is kept in state
	// until it is pruned.
	DefaultBlocksToPrune = uint64(1000)
)

// NewParams creates a new Params object, this method will panic if the provided ConclusionVerificationParams
// cannot be encoded to an Any.
//
// NOTICE: The MaxBlockAge + BlocksToPrune < UnbondingPeriod. This inequality is required to ensure that
// any infracting validator cannot unbond to avoid being slashed.
func NewParams(ap AlertParams, cvp ConclusionVerificationParams, pp PruningParams) Params {
	params := Params{
		AlertParams:   ap,
		PruningParams: pp,
	}

	var err error
	if cvp != nil {
		params.ConclusionVerificationParams, err = codectypes.NewAnyWithValue(cvp)
		if err != nil {
			panic(err)
		}
	}

	return params
}

// DefaultParams returns a default set of parameters, i.e. BondAmount is set to
// DefaultBondAmount.
func DefaultParams(denom string, cvp ConclusionVerificationParams) Params {
	return NewParams(
		AlertParams{
			Enabled:     DefaultEnabled,
			BondAmount:  sdk.NewCoin(denom, DefaultBondAmount),
			MaxBlockAge: DefaultAlertExpiry,
		},
		cvp,
		PruningParams{
			Enabled:       DefaultPruningEnabled,
			BlocksToPrune: DefaultBlocksToPrune,
		},
	)
}

// Validate performs a basic validation of the AlertParams, i.e. if Alerts are enabled, that the
// bond amount is non-zero, and that the MaxBlockAge is non-zero.
func (ap *AlertParams) Validate() error {
	if !ap.Enabled {
		if !ap.BondAmount.IsZero() || !(ap.MaxBlockAge == 0) {
			return fmt.Errorf("invalid alert params: bond amount must be zero if alerts are disabled")
		}

		return nil
	}

	if ap.BondAmount.IsZero() || ap.BondAmount.IsNegative() {
		return fmt.Errorf("invalid alert params: bond amount must be non-zero")
	}

	if ap.MaxBlockAge == 0 {
		return fmt.Errorf("invalid alert params: max block age must be non-zero")
	}

	return nil
}

// Validate performs basic validation of the Pruning Params, specifically, that the BlocksToPrune is
// non-zero if pruning is enabled, and zero if disabled.
func (pp *PruningParams) Validate() error {
	if !pp.Enabled {
		if pp.BlocksToPrune != 0 {
			return fmt.Errorf("invalid pruning params: blocks to prune must be zero if pruning is disabled")
		}

		return nil
	}

	if pp.BlocksToPrune == 0 {
		return fmt.Errorf("invalid pruning params: blocks to prune must be non-zero if pruning is enabled")
	}

	return nil
}

// Validate performs a basic validation of the Params, i.e. that the AlertParams are valid,
// and the ConclusionVerificationParams are valid (if present).
func (p *Params) Validate() error {
	if p.ConclusionVerificationParams != nil {
		// Unmarshal the Any into a ConclusionVerificationParams
		var params ConclusionVerificationParams
		if err := pc.UnpackAny(p.ConclusionVerificationParams, &params); err != nil {
			return err
		}

		// Validate the ConclusionVerificationParams
		if err := params.ValidateBasic(); err != nil {
			return err
		}
	}

	if err := p.PruningParams.Validate(); err != nil {
		return err
	}

	return p.AlertParams.Validate()
}
