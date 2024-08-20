package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	compression "github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/ve/types"
	"github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/skip-mev/connect/v2/x/oracle/types/mocks"
)

// CreateTestOracleKeeperWithGenesis creates a test oracle keeper with the given genesis state.
func CreateTestOracleKeeperWithGenesis(t *testing.T, ctx sdk.Context, key *storetypes.KVStoreKey, genesis oracletypes.GenesisState) keeper.Keeper {
	t.Helper()

	ss := runtime.NewKVStoreService(key)
	encCfg := moduletestutil.MakeTestEncodingConfig()

	k := keeper.NewKeeper(
		ss,
		encCfg.Codec,
		mocks.NewMarketMapKeeper(t),
		sdk.AccAddress("authority"),
	)

	k.InitGenesis(ctx, genesis)

	return k
}

// CreateExtendedCommitInfo creates an extended commit info with the given commit info.
func CreateExtendedCommitInfo(commitInfo []cometabci.ExtendedVoteInfo, codec compression.ExtendedCommitCodec) (cometabci.ExtendedCommitInfo, []byte, error) {
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
	codec compression.VoteExtensionCodec,
) (cometabci.ExtendedVoteInfo, error) {
	return CreateExtendedVoteInfoWithPower(consAddr, 1, prices, codec)
}

// CreateExtendedVoteInfoWithPower CreateExtendedVoteInfo creates an extended vote info
// with the given power, prices, timestamp and height.
func CreateExtendedVoteInfoWithPower(
	consAddr sdk.ConsAddress,
	power int64,
	prices map[uint64][]byte,
	codec compression.VoteExtensionCodec,
) (cometabci.ExtendedVoteInfo, error) {
	ve, err := CreateVoteExtensionBytes(prices, codec)
	if err != nil {
		return cometabci.ExtendedVoteInfo{}, err
	}
	voteInfo := cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: consAddr,
			Power:   power,
		},
		VoteExtension: ve,
		BlockIdFlag:   cometproto.BlockIDFlagCommit,
	}

	return voteInfo, nil
}

// UpdateContextWithVEHeight updates the context with the given height and enables vote extensions
// for the given height.
func UpdateContextWithVEHeight(ctx sdk.Context, height int64) sdk.Context {
	params := cometproto.ConsensusParams{
		Abci: &cometproto.ABCIParams{
			VoteExtensionsEnableHeight: height,
		},
	}

	ctx = ctx.WithConsensusParams(params)
	return ctx
}

// CreateBaseSDKContextWithKeys creates a base sdk context with the given store key and transient key.
func CreateBaseSDKContextWithKeys(t *testing.T, storeKey storetypes.StoreKey, transientKey *storetypes.TransientStoreKey) sdk.Context {
	t.Helper()

	testCtx := testutil.DefaultContextWithDB(
		t,
		storeKey,
		transientKey,
	)

	return testCtx.Ctx
}

// CreateBaseSDKContext creates a base sdk context with the default store key and transient key.
func CreateBaseSDKContext(t *testing.T) sdk.Context {
	t.Helper()

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
	codec compression.VoteExtensionCodec,
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
