package voteweighted

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	consumerkeeper "github.com/cosmos/interchain-security/v5/x/ccv/consumer/keeper"
	ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
)

var _ ValidatorStore = CCVConsumerCompatKeeper{}
var _ stakingtypes.ValidatorI = CCVCompat{}

// CCVCompat is used for compatibility between stakingtypes.ValidatorI and CrossChainValidator.
type CCVCompat struct {
	stakingtypes.ValidatorI
	ccv ccvtypes.CrossChainValidator
}

// GetBondedTokens returns the power of the validator as math.Int.
func (c CCVCompat) GetBondedTokens() math.Int {
	return math.NewInt(c.ccv.Power)
}

// CCVConsumerCompatKeeper is used for compatibility between the consumer keeper and the ValidatorStore interface.
type CCVConsumerCompatKeeper struct {
	ccvKeeper consumerkeeper.Keeper
}

// ValidatorByConsAddr returns a compat validator from the consumer keeper.
func (c CCVConsumerCompatKeeper) ValidatorByConsAddr(ctx context.Context, addr sdk.ConsAddress) (stakingtypes.ValidatorI, error) {
	sdkCtx, ok := ctx.(sdk.Context)
	if !ok {
		return nil, fmt.Errorf("could not convert context to sdk.Context")
	}
	ccv, found := c.ccvKeeper.GetCCValidator(sdkCtx, addr.Bytes())
	if !found {
		return nil, fmt.Errorf("could not find validator %s", addr.String())
	}
	return CCVCompat{ccv: ccv}, nil
}

// TotalBondedTokens iterates through all CCVs and returns the sum of all validator power.
func (c CCVConsumerCompatKeeper) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	total := math.NewInt(0)
	sdkCtx, ok := ctx.(sdk.Context)
	if !ok {
		return total, fmt.Errorf("could not convert context to sdk.Context")
	}
	for _, ccVal := range c.ccvKeeper.GetAllCCValidator(sdkCtx) {
		total = total.Add(math.NewInt(ccVal.Power))
	}
	return total, nil
}
