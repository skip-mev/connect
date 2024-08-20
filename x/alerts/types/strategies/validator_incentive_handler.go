package strategies

import (
	"math/big"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkyabci "github.com/skip-mev/connect/v2/abci/ve/types"
	"github.com/skip-mev/connect/v2/x/alerts/types"
	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

// DefaultHandleValidatorIncentive returns a ValidatorIncentiveHandler which creates a ValidatorAlertIncentive if a validator
// reported a price that lied outside the bounds of what was expected off-chain. If a validator failed to report a price for
// the ticker, or their price was within the bounds, no incentive is issued.
//
// NOTICE: no signature checks are performed on the vote-extension, as it is expected that the caller has verified the
// ExtendedVoteInfo's signature before calling this function.
func DefaultHandleValidatorIncentive() ValidatorIncentiveHandler {
	return func(ve cmtabci.ExtendedVoteInfo, pb types.PriceBound, a types.Alert, cpID uint64) (incentivetypes.Incentive, error) {
		// validate the alert
		if err := a.ValidateBasic(); err != nil {
			return nil, err
		}

		// validate the price-bound
		if err := pb.ValidateBasic(); err != nil {
			return nil, err
		}

		// unmarshal the vote-extension
		var voteExt slinkyabci.OracleVoteExtension
		if err := voteExt.Unmarshal(ve.VoteExtension); err != nil {
			return nil, err
		}

		// check for existence, if it doesn't exist, return nil
		priceBz, ok := voteExt.Prices[cpID]
		if !ok {
			return nil, nil
		}

		var price big.Int
		price.SetBytes(priceBz)

		// check bounds
		low, err := pb.GetLowInt()
		if err != nil {
			return nil, err
		}

		high, err := pb.GetHighInt()
		if err != nil {
			return nil, err
		}

		// check if the price is outside the bounds
		if !(price.Cmp(low) >= 0) || !(price.Cmp(high) <= 0) {
			// get signer address
			signer, err := sdk.AccAddressFromBech32(a.Signer)
			if err != nil {
				return nil, err
			}

			return NewValidatorAlertIncentive(ve.Validator, a.Height, signer), nil
		}

		return nil, nil
	}
}
