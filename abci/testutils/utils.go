package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/strategies"
	"github.com/skip-mev/slinky/abci/ve/types"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// CreateTestOracleKeeperWithGenesis creates a test oracle keeper with the given genesis state.
func CreateTestOracleKeeperWithGenesis(ctx sdk.Context, key storetypes.StoreKey, genesis oracletypes.GenesisState) keeper.Keeper {
	keeper := keeper.NewKeeper(
		key,
		sdk.AccAddress([]byte("authority")),
	)

	keeper.InitGenesis(ctx, genesis)

	return keeper
}

// CreateExtendedCommitInfo creates an extended commit info with the given commit info.
func CreateExtendedCommitInfo(commitInfo []cometabci.ExtendedVoteInfo, codec strategies.ExtendedCommitCodec) (cometabci.ExtendedCommitInfo, []byte, error) {
	extendedCommitInfo := cometabci.ExtendedCommitInfo{
		Votes: commitInfo,
	}

	bz, err := codec.Encode(extendedCommitInfo)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, nil, err
	}

	return extendedCommitInfo, bz, nil
}

// CreateExtendedVoteInfo creates an extended vote info with the given prices, timestamp and height.
func CreateExtendedVoteInfo(
	consAddr sdk.ConsAddress,
	prices map[uint64][]byte,
	codec strategies.VoteExtensionCodec,
) (cometabci.ExtendedVoteInfo, error) {
	ve, err := CreateVoteExtensionBytes(prices, codec)
	if err != nil {
		return cometabci.ExtendedVoteInfo{}, err
	}

	voteInfo := cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: consAddr,
		},
		VoteExtension: ve,
	}

	return voteInfo, nil
}

// UpdateContextWithVEHeight updates the context with the given height and enables vote extensions
// for the given height.
func UpdateContextWithVEHeight(ctx sdk.Context, height int64) sdk.Context {
	params := cmtproto.ConsensusParams{
		Abci: &cmtproto.ABCIParams{
			VoteExtensionsEnableHeight: height,
		},
	}

	ctx = ctx.WithConsensusParams(params)
	return ctx
}

// CreateBaseSDKContextWithKeys creates a base sdk context with the given store key and transient key.
func CreateBaseSDKContextWithKeys(t *testing.T, storekey storetypes.StoreKey, transientkey *storetypes.TransientStoreKey) sdk.Context {
	testCtx := testutil.DefaultContextWithDB(
		t,
		storekey,
		transientkey,
	)

	return testCtx.Ctx
}

// CreateBaseSDKContext creates a base sdk context with the default store key and transient key.
func CreateBaseSDKContext(t *testing.T) sdk.Context {
	key := storetypes.NewKVStoreKey(oracletypes.StoreKey)

	testCtx := testutil.DefaultContextWithDB(
		t,
		key,
		storetypes.NewTransientStoreKey("transient_test"),
	)

	return testCtx.Ctx
}

// CreateVoteExtensionBytes creates a vote extension bytes with the given prices, timestamp and height.
func CreateVoteExtensionBytes(
	prices map[uint64][]byte,
	codec strategies.VoteExtensionCodec,
) ([]byte, error) {
	voteExtension := CreateVoteExtension(prices)
	voteExtensionBz, err := codec.Encode(voteExtension)
	if err != nil {
		return nil, err
	}

	return voteExtensionBz, nil
}

// CreateVoteExtension creates a vote extension with the given prices, timestamp and height.
func CreateVoteExtension(
	prices map[uint64][]byte,
) types.OracleVoteExtension {
	return types.OracleVoteExtension{
		Prices: prices,
	}
}
