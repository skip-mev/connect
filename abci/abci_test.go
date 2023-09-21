package abci_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/abci"
	"github.com/skip-mev/slinky/abci/mocks"
	abcitypes "github.com/skip-mev/slinky/abci/types"
	oracleservice "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/suite"
)

type ABCITestSuite struct {
	suite.Suite
	ctx                         sdk.Context
	voteExtensionsEnabledHeight int64

	// ProposalHandler set up.
	proposalHandler        *abci.ProposalHandler
	prepareProposalHandler sdk.PrepareProposalHandler
	processProposalHandler sdk.ProcessProposalHandler
	aggregateFn            oracleservice.AggregateFnFromContext

	// oracle keeper set up.
	oracleKeeper   keeper.Keeper
	currencyPairs  []oracletypes.CurrencyPair
	genesis        oracletypes.GenesisState
	validatorStore abci.ValidatorStore
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ABCITestSuite))
}

func (suite *ABCITestSuite) SetupTest() {
	suite.setUpOracleKeeper()

	// Use the default no-op prepare and process proposal handlers from the sdk.
	suite.prepareProposalHandler = baseapp.NoOpPrepareProposal()
	suite.processProposalHandler = baseapp.NoOpProcessProposal()
	suite.aggregateFn = oracleservice.ComputeMedianWithContext

	// Create an empty default mock validator store.
	suite.validatorStore = suite.createMockValidatorStore(
		[]validator{},
		math.NewIntFromUint64(0),
	)

	oracle := abci.NewOracle(
		log.NewTestLogger(suite.T()),
		suite.aggregateFn,
		suite.oracleKeeper,
		suite.NoOpValidateVEFn(),
		suite.validatorStore,
	)

	suite.proposalHandler = abci.NewProposalHandler(
		log.NewTestLogger(suite.T()),
		suite.prepareProposalHandler,
		suite.processProposalHandler,
		oracle,
	)
}

func (suite *ABCITestSuite) NoOpValidateVEFn() abci.ValidateVoteExtensionsFn {
	return abci.NoOpValidateVoteExtensions
}

func (suite *ABCITestSuite) setUpOracleKeeper() {
	key := storetypes.NewKVStoreKey(oracletypes.StoreKey)
	suite.oracleKeeper = keeper.NewKeeper(
		key,
		sdk.AccAddress([]byte("authority")),
	)

	testCtx := testutil.DefaultContextWithDB(
		suite.T(),
		key,
		storetypes.NewTransientStoreKey("transient_test"),
	)
	suite.ctx = testCtx.Ctx.WithBlockHeight(3)

	suite.voteExtensionsEnabledHeight = 1
	params := cmtproto.ConsensusParams{
		Abci: &cmtproto.ABCIParams{
			VoteExtensionsEnableHeight: suite.voteExtensionsEnabledHeight,
		},
	}
	suite.ctx = suite.ctx.WithConsensusParams(params)

	suite.currencyPairs = []oracletypes.CurrencyPair{
		{
			Base:  "BTC",
			Quote: "ETH",
		},
		{
			Base:  "BTC",
			Quote: "USD",
		},
		{
			Base:  "ETH",
			Quote: "USD",
		},
	}
	genesisCPs := []oracletypes.CurrencyPairGenesis{
		{
			CurrencyPair: suite.currencyPairs[0],
			Nonce:        0,
		},
		{
			CurrencyPair: suite.currencyPairs[1],
			Nonce:        0,
		},
		{
			CurrencyPair: suite.currencyPairs[2],
			Nonce:        0,
		},
	}
	suite.genesis = oracletypes.GenesisState{
		CurrencyPairGenesis: genesisCPs,
	}

	suite.oracleKeeper.InitGenesis(suite.ctx, suite.genesis)
}

func (suite *ABCITestSuite) createMockValidatorStore(
	validators []validator,
	totalTokens math.Int,
) *mocks.ValidatorStore {
	store := mocks.NewValidatorStore(suite.T())
	if len(validators) != 0 {
		mockVals := make([]*mocks.ValidatorI, len(validators))
		valPubKeys := make([]crypto.PublicKey, len(validators))

		for i, val := range validators {
			mockVals[i] = mocks.NewValidatorI(suite.T())
			mockVals[i].On(
				"GetBondedTokens",
			).Return(
				val.stake,
			).Maybe()

			store.On(
				"ValidatorByConsAddr",
				suite.ctx,
				val.consAddr,
			).Return(
				mockVals[i],
				nil,
			).Maybe()

			var err error
			valPubKeys[i], err = cryptocodec.ToCmtProtoPublicKey(ed25519.GenPrivKey().PubKey())
			if err != nil {
				panic(err)
			}

			store.On(
				"GetPubKeyByConsAddr",
				suite.ctx,
				val.consAddr,
			).Return(
				valPubKeys[i],
				nil,
			).Maybe()
		}
	}

	store.On(
		"TotalBondedTokens",
		suite.ctx,
	).Return(
		totalTokens, nil,
	).Maybe()

	return store
}

func (suite *ABCITestSuite) createRequestPrepareProposal(
	extendedCommitInfo cometabci.ExtendedCommitInfo,
	txs [][]byte,
) *cometabci.RequestPrepareProposal {
	return &cometabci.RequestPrepareProposal{
		Txs:             txs,
		LocalLastCommit: extendedCommitInfo,
	}
}

func (suite *ABCITestSuite) createExtendedCommitInfo(
	commitInfo []cometabci.ExtendedVoteInfo,
) cometabci.ExtendedCommitInfo {
	return cometabci.ExtendedCommitInfo{
		Votes: commitInfo,
	}
}

func (suite *ABCITestSuite) createExtendedVoteInfo(
	consAddr sdk.ConsAddress,
	prices map[string]string,
	timestamp time.Time,
	height int64,
) cometabci.ExtendedVoteInfo {
	return cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: consAddr,
		},
		VoteExtension: suite.createVoteExtensionBytes(prices, timestamp, height),
	}
}

func (suite *ABCITestSuite) createVoteExtensionBytes(
	prices map[string]string,
	timestamp time.Time,
	height int64,
) []byte {
	voteExtension := suite.createVoteExtension(prices, timestamp, height)
	voteExtensionBz, err := voteExtension.Marshal()
	suite.Require().NoError(err)

	return voteExtensionBz
}

func (suite *ABCITestSuite) createVoteExtension(
	prices map[string]string,
	timestamp time.Time,
	height int64,
) *abcitypes.OracleVoteExtension {
	return &abcitypes.OracleVoteExtension{
		Prices:    prices,
		Timestamp: timestamp,
		Height:    height,
	}
}

func (suite *ABCITestSuite) createValAddress(prefix string) sdk.ConsAddress {
	return sdk.ConsAddress(prefix + suite.T().Name())
}

func (suite *ABCITestSuite) createOracleData(
	finalPrices map[string]string,
	timestamp time.Time,
	height int64,
	extendedVoteInfos []cometabci.ExtendedVoteInfo,
) abcitypes.OracleData {
	extendedCommitInfo := suite.createExtendedCommitInfo(extendedVoteInfos)
	infoBz, err := extendedCommitInfo.Marshal()
	suite.Require().NoError(err)

	return abcitypes.OracleData{
		Prices:             finalPrices,
		Timestamp:          timestamp,
		Height:             height,
		ExtendedCommitInfo: infoBz,
	}
}
