package abci_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/skip-mev/slinky/abci"
	"github.com/skip-mev/slinky/abci/mocks"
	abcitypes "github.com/skip-mev/slinky/abci/types"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	oraclemoduletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/suite"
)

type ABCITestSuite struct {
	suite.Suite
	ctx sdk.Context

	// ProposalHandler set up.
	proposalHandler        *abci.ProposalHandler
	prepareProposalHandler sdk.PrepareProposalHandler
	processProposalHandler sdk.ProcessProposalHandler
	aggregateFn            oracletypes.AggregateFn
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ABCITestSuite))
}

func (suite *ABCITestSuite) SetupTest() {
	testCtx := testutil.DefaultContextWithDB(
		suite.T(),
		storetypes.NewKVStoreKey(oraclemoduletypes.StoreKey),
		storetypes.NewTransientStoreKey("transient_test"),
	)
	suite.ctx = testCtx.Ctx.WithBlockHeight(1)

	// Use the default no-op prepare and process proposal handlers from the sdk.
	suite.prepareProposalHandler = baseapp.NoOpPrepareProposal()
	suite.processProposalHandler = baseapp.NoOpProcessProposal()
	suite.aggregateFn = oracletypes.ComputeMedian()

	suite.proposalHandler = abci.NewProposalHandler(
		log.NewNopLogger(),
		suite.prepareProposalHandler,
		suite.processProposalHandler,
		suite.aggregateFn,
	)
}

func (suite *ABCITestSuite) createMockValidatorStore(
	validators []validator,
	totalTokens math.Int,
) *mocks.ValidatorStore {
	store := mocks.NewValidatorStore(suite.T())
	if len(validators) != 0 {
		for _, val := range validators {
			store.On(
				"GetValidator",
				suite.ctx,
				val.address,
			).Return(
				stakingtypes.Validator{
					Tokens: val.stake,
					Status: stakingtypes.Bonded,
				},
				true,
			)
		}
	}

	store.On(
		"TotalBondedTokens",
		suite.ctx,
	).Return(
		totalTokens,
	)

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
	valAddress sdk.ValAddress,
	prices map[string]string,
	timestamp time.Time,
	height int64,
) cometabci.ExtendedVoteInfo {
	return cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: valAddress.Bytes(),
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

func (suite *ABCITestSuite) createValAddress(prefix string) sdk.ValAddress {
	return sdk.ValAddress(prefix + suite.T().Name())
}
