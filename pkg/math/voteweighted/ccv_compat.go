package voteweighted

import (
	"context"
	"fmt"

	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvtypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
)

var (
	_ ValidatorStore          = CCVConsumerCompatKeeper{}
	_ stakingtypes.ValidatorI = CCVCompat{}
)

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
	ccvStore CCValidatorStore
}

// NewCCVConsumerCompatKeeper constructs a CCVConsumerCompatKeeper from a consumer keeper.
func NewCCVConsumerCompatKeeper(ccvStore CCValidatorStore) CCVConsumerCompatKeeper {
	return CCVConsumerCompatKeeper{
		ccvStore: ccvStore,
	}
}

// ValidatorByConsAddr returns a compat validator from the consumer keeper.
func (c CCVConsumerCompatKeeper) ValidatorByConsAddr(ctx context.Context, addr sdk.ConsAddress) (stakingtypes.ValidatorI, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ccv, found := c.ccvStore.GetCCValidator(sdkCtx, addr.Bytes())
	if !found {
		return nil, fmt.Errorf("could not find validator %s", addr.String())
	}
	return CCVCompat{ccv: ccv}, nil
}

// TotalBondedTokens iterates through all CCVs and returns the sum of all validator power.
func (c CCVConsumerCompatKeeper) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	total := math.NewInt(0)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, ccVal := range c.ccvStore.GetAllCCValidator(sdkCtx) {
		total = total.Add(math.NewInt(ccVal.Power))
	}
	return total, nil
}

// GetPubKeyByConsAddr returns the public key of a validator given the consensus addr.
func (c CCVConsumerCompatKeeper) GetPubKeyByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (cmtprotocrypto.PublicKey, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	val, found := c.ccvStore.GetCCValidator(sdkCtx, consAddr)
	if !found {
		return cmtprotocrypto.PublicKey{}, fmt.Errorf("not found CCValidator for address: %s", consAddr.String())
	}

	consPubKey, err := val.ConsPubKey()
	if err != nil {
		return cmtprotocrypto.PublicKey{}, fmt.Errorf("could not get pubkey for val %s: %w", val.String(), err)
	}
	tmPubKey, err := cryptocodec.ToCmtProtoPublicKey(consPubKey)
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	return tmPubKey, nil
}
