package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Strategy defines the callback function that will be executed by the
// incentives module with the given context and incentive. The strategy is
// responsible for the following:
//  1. If the strategy desires to update the incentive, it must return the
//     updated incentive.
//  2. If the strategy desires to delete the incentive, it must return nil.
//  3. If the strategy desires to leave the incentive unchanged, it must return
//     the same incentive.
//  4. Applying any desired state transitions such as minting rewards, or
//     slashing.
//
// For an example implementation, please see the examples/ directory.
type Strategy func(ctx sdk.Context, incentive Incentive) (updatedIncentive Incentive, err error)
