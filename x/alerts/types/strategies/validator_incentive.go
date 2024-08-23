package strategies

import (
	"fmt"

	"cosmossdk.io/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

const (
	// ValidatorAlertIncentiveType is the type of incentive issued for validators upon being referenced in a Conclusion.
	ValidatorAlertIncentiveType = "validator_alert"
)

var (
	_                  incentivetypes.Incentive = (*ValidatorAlertIncentive)(nil)
	defaultSlashFactor                          = math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 1)
)

// NewValidatorAlertIncentive returns a new ValidatorAlertIncentive. ValidatorAlertIncentive defines the incentive strategy to
// be executed for a validator that has been confirmed to have at fault for an x/alerts alert. This strategy is expected to
// slash half of the validator's stake, and reward the slashed stake to the alerter.
func NewValidatorAlertIncentive(validator cmtabci.Validator, alertHeight uint64, signer sdk.AccAddress) incentivetypes.Incentive {
	return &ValidatorAlertIncentive{
		Validator:   validator,
		AlertHeight: alertHeight,
		AlertSigner: signer.String(),
	}
}

// ValidateBasic does a basic stateless validation check on the ValidatorAlertIncentive. Specifically, this method
// checks that the validator's address is valid, and it's power is non-negative.
func (b *ValidatorAlertIncentive) ValidateBasic() error {
	// the only check we can do statelessly is that the validator address is non-nil
	if b.Validator.Address == nil {
		return fmt.Errorf("validator address cannot be nil")
	}

	// the only check we can do statelessly is that the validator power is non-negative
	if b.Validator.Power <= 0 {
		return fmt.Errorf("validator power must be > 0: %d", b.Validator.Power)
	}

	// check signer validity
	if _, err := sdk.AccAddressFromBech32(b.AlertSigner); err != nil {
		return fmt.Errorf("validator address is invalid: %w", err)
	}

	return nil
}

// Type returns the type of the incentive.
func (b *ValidatorAlertIncentive) Type() string {
	return ValidatorAlertIncentiveType
}

// Copy returns a copy of the incentive.
func (b *ValidatorAlertIncentive) Copy() incentivetypes.Incentive {
	val := b.Validator

	val.Address = make([]byte, len(b.Validator.Address))

	// we need to copy the address, since it's a reference type
	copy(val.Address, b.Validator.Address)

	return &ValidatorAlertIncentive{
		Validator:   val,
		AlertSigner: b.AlertSigner,
		AlertHeight: b.AlertHeight,
	}
}

// DefaultValidatorAlertIncentiveStrategy is the default strategy for issuing incentives to validators upon being deemed
// at fault for an x/alerts alert. This method returns a Strategy that executes wrt. the given StakingKeeper / BankKeeper.
//
// NOTICE:
// The DefaultSlashFactor is 50% of each validator's stake. See NewValidatorAlertIncentiveStrategy for more details.
func DefaultValidatorAlertIncentiveStrategy(sk types.StakingKeeper, bk types.BankKeeper) incentivetypes.Strategy {
	return NewValidatorAlertIncentiveStrategy(sk, bk, defaultSlashFactor)
}

// NewValidatorAlertIncentiveStrategy is the default strategy for issuing incentives to validators upon being
// referenced in a Conclusion. This method returns a Strategy that executes wrt. the given StakingKeeper / BankKeeper.
// Notice, this strategy will slash half of the validator's stake, and mint the amount slashed to the alerter.
//
// CONTRACT: as of v0.50.0-rc2 of the Cosmos SDK, the Slash method will burn staked tokens, this is crucial to our logic
// in order for this operation to not inflate the bond-denom's total supply.
func NewValidatorAlertIncentiveStrategy(sk types.StakingKeeper, bk types.BankKeeper, slashFactor math.LegacyDec) incentivetypes.Strategy {
	return func(ctx sdk.Context, incentive incentivetypes.Incentive) (_ incentivetypes.Incentive, err error) {
		// assert type of incentive
		validatorAlertIncentive, ok := incentive.(*ValidatorAlertIncentive)
		if !ok {
			return nil, fmt.Errorf("incentive must be of type ValidatorAlertIncentive, got %T", incentive)
		}

		ctx.Logger().Info("validator alert incentive executed", "incentive", validatorAlertIncentive)

		ca := sdk.ConsAddress(validatorAlertIncentive.Validator.Address)
		// check that the validator exists
		if _, err := sk.GetValidatorByConsAddr(ctx, ca); err != nil {
			return nil, fmt.Errorf("validator with address %s does not exist", ca)
		}

		// adjust the alert height to account for the validator update delay to comet, notice
		// comet's view of the validator set is always buffed, in that updates returned from the app
		// are applied after the validator update delay blocks from when they are given to comet.
		infractionHeight := validatorAlertIncentive.AlertHeight - uint64(sdk.ValidatorUpdateDelay)

		// slash the validator
		//nolint:gosec
		amountSlashed, err := sk.Slash(ctx, ca, int64(infractionHeight), validatorAlertIncentive.Validator.Power,
			slashFactor) //nolint:gosec
		if err != nil {
			return nil, fmt.Errorf("failed to slash validator: %w", err)
		}

		ctx.Logger().Info("slashed validator", "validator", validatorAlertIncentive.Validator, "amount_slashed", amountSlashed)

		// get bond denom to mint to alerter from slashed validator
		denom, err := sk.BondDenom(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get bond denom: %w", err)
		}
		coinsToMint := sdk.NewCoins(sdk.NewCoin(denom, amountSlashed))

		// mint the slashed tokens (burned) to the signer of the alert
		if err := bk.MintCoins(ctx, types.ModuleName, coinsToMint); err != nil {
			return nil, fmt.Errorf("failed to mint coins: %w", err)
		}

		alertSigner, err := sdk.AccAddressFromBech32(validatorAlertIncentive.AlertSigner)
		if err != nil {
			return nil, fmt.Errorf("failed to parse alert signer address: %w", err)
		}

		// send the slashed tokens to the signer of the alert
		if err := bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, alertSigner, coinsToMint); err != nil {
			return nil, fmt.Errorf("failed to send coins: %w", err)
		}

		ctx.Logger().Info("minted coins to alert signer", "signer", alertSigner, "amount_minted", amountSlashed)

		return nil, nil
	}
}
