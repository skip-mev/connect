package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// ExecSLA enforces the SLA criteria for all price feeds that it is maintaining.
// This function is called at the beginning of every block and already assumes that
// all price feeds have been updated for the current block via the pre-block hook.
func (k *Keeper) ExecSLA(ctx sdk.Context, sla slatypes.PriceFeedSLA) error {
	// Ensure that the SLA should be checked for the current block height.
	height := ctx.BlockHeight()
	if height == 0 || height%int64(sla.Frequency) != 0 {
		return nil
	}

	// Fetch all price feeds for the given SLA.
	feeds, err := k.GetAllPriceFeeds(ctx, sla.ID)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to get price feeds for SLA",
			"sla", sla.ID,
			"err", err,
		)

		return err
	}

	// Iterate through all price feeds and check if the price feed
	// qualifies for an SLA check + meets the SLA criteria.
	for _, priceFeed := range feeds {
		qualifies, err := sla.Qualifies(priceFeed)
		if err != nil {
			k.Logger(ctx).Error(
				"unable to determine if price feed qualifies for SLA",
				"sla", sla.ID,
				"err", err,
			)

			return err
		}
		if !qualifies {
			k.Logger(ctx).Info("price feed does not qualify for SLA check")
			continue
		}

		if err := k.EnforceSLA(ctx, sla, priceFeed); err != nil {
			k.Logger(ctx).Error(
				"failed to check SLA",
				"sla", sla.ID,
				"err", err,
			)

			return err
		}
	}

	return nil
}

// EnforceSLA checks whether the given price feed meets the criteria for
// the given SLA. If the price feed has met the expected uptime, then no action is
// taken. Otherwise, the validator is slashed by the deviation from the
// expected uptime.
func (k *Keeper) EnforceSLA(ctx sdk.Context, sla slatypes.PriceFeedSLA, priceFeed slatypes.PriceFeed) error {
	// Ensure that the validator exists. In the event that the validator
	// does not exist, we will delete the price feed from the store.
	validator := sdk.ValAddress(priceFeed.Validator)
	power, err := k.stakingKeeper.GetLastValidatorPower(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to get last validator power; removing incentive for validator",
			"validator", validator.String(),
			"err", err,
		)

		return k.RemovePriceFeed(ctx, sla.ID, priceFeed.CurrencyPair, priceFeed.Validator)
	}

	// Determine the uptime for the price feed.
	uptime, err := sla.GetUptimeFromPriceFeed(priceFeed)
	if err != nil {
		k.Logger(ctx).Error(
			"unable to get uptime from SLA",
			"err", err,
		)

		return err
	}

	// Check if the validator is subject to slashing.
	if uptime.GTE(sla.ExpectedUptime) {
		k.Logger(ctx).Info(
			"validator met SLA",
			"validator", validator.String(),
			"uptime", uptime,
			"expected_uptime", sla.ExpectedUptime,
		)

		return nil
	}

	// deviation = ((expected_uptime - uptime) / expected_uptime) * K
	deviation := (sla.ExpectedUptime.Sub(uptime)).Quo(sla.ExpectedUptime)
	slashFactor := deviation.Mul(sla.SlashConstant)

	k.Logger(ctx).Info(
		"validator did not meet SLA",
		"validator", validator.String(),
		"uptime", uptime,
		"deviation", deviation,
		"expected_uptime", sla.ExpectedUptime,
		"slash_factor", slashFactor,
	)

	return k.Slash(ctx, validator, power, slashFactor)
}

// Slash will slash the validator with the given power and slash factor.
func (k *Keeper) Slash(
	ctx sdk.Context,
	validator sdk.ValAddress,
	power int64,
	slashFactor math.LegacyDec,
) error {
	// We do height - ValidatorUpdateDelay because the vote extensions included in this block
	// were constructed in the previous block.
	height := ctx.BlockHeight() - sdk.ValidatorUpdateDelay
	amount, err := k.slashingKeeper.Slash(ctx, sdk.ConsAddress(validator), height, power, slashFactor)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to slash validator",
			"validator", validator.String(),
			"err", err,
		)

		return err
	}

	k.Logger(ctx).Info(
		"slashed validator",
		"validator", validator.String(),
		"amount", amount.String(),
	)

	return nil
}
