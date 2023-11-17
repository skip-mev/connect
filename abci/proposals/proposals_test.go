package proposals_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/abci/proposals"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/abci/ve"
)

var (
	prices1 = map[string]string{
		"BTC/USD": "0x100",
		"ETH/USD": "0x200",
	}
	prices2 = map[string]string{
		"BTC/USD": "0x300",
		"ETH/USD": "0x400",
	}
	prices3 = map[string]string{
		"BTC/USD": "0x500",
		"ETH/USD": "0x600",
	}
	malformedPrices = map[string]string{
		"BTC/USD": "no_rizz",
	}

	val1 = sdk.ConsAddress("val1")
	val2 = sdk.ConsAddress("val2")
	val3 = sdk.ConsAddress("val3")
)

type ProposalsTestSuite struct {
	suite.Suite

	ctx sdk.Context

	proposalHandler        *proposals.ProposalHandler
	prepareProposalHandler sdk.PrepareProposalHandler
	processProposalHandler sdk.ProcessProposalHandler
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ProposalsTestSuite))
}

func (s *ProposalsTestSuite) SetupTest() {
	s.ctx = testutils.CreateBaseSDKContext(s.T())
	s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 1)
	s.ctx = s.ctx.WithBlockHeight(1)

	// Use the default no-op prepare and process proposal handlers from the sdk.
	s.prepareProposalHandler = baseapp.NoOpPrepareProposal()
	s.processProposalHandler = baseapp.NoOpProcessProposal()
	s.proposalHandler = proposals.NewProposalHandler(
		log.NewTestLogger(s.T()),
		s.prepareProposalHandler,
		s.processProposalHandler,
		ve.NoOpValidateVoteExtensions,
	)
}

func (s *ProposalsTestSuite) TestPrepareProposal() {
	prepareHandler := s.proposalHandler.PrepareProposalHandler()

	s.Run("vote extensions not enabled", func() {
		proposal := [][]byte{
			[]byte("tx1"),
		}

		req := s.createRequestPrepareProposal(
			cometabci.ExtendedCommitInfo{},
			proposal,
			0,
		)

		response, err := prepareHandler(s.ctx, req)
		s.Require().Equal(response.Txs, proposal)
		s.Require().NoError(err)
	})

	s.Run("vote extensions disabled with multiple txs", func() {
		proposal := [][]byte{
			[]byte("tx1"),
			[]byte("tx2"),
		}

		req := s.createRequestPrepareProposal(
			cometabci.ExtendedCommitInfo{},
			proposal,
			0,
		)

		response, err := prepareHandler(s.ctx, req)
		s.Require().Equal(response.Txs, proposal)
		s.Require().NoError(err)
	})

	// Enable vote extensions for the remaining tests.
	s.ctx = s.ctx.WithBlockHeight(3)

	s.Run("vote extensions enabled with no txs and a single vote extension", func() {
		proposal := [][]byte{}

		valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		commitInfo, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		req := s.createRequestPrepareProposal(
			commitInfo,
			proposal,
			3,
		)

		response, err := prepareHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(response.Txs, 1)
		s.Require().Equal(response.Txs[0], commitInfoBz)
	})

	s.Run("vote extensions enabled with multiple txs and a single vote extension", func() {
		proposal := [][]byte{
			[]byte("tx1"),
			[]byte("tx2"),
		}

		valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		commitInfo, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		req := s.createRequestPrepareProposal(
			commitInfo,
			proposal,
			3,
		)

		response, err := prepareHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(response.Txs, 3)
		s.Require().Equal(response.Txs[0], commitInfoBz)
		s.Require().Equal(response.Txs[1], proposal[0])
		s.Require().Equal(response.Txs[2], proposal[1])
	})

	s.Run("vote extensions enabled with multiple vote extensions", func() {
		proposal := [][]byte{}

		valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, time.Now(), 2)
		s.Require().NoError(err)

		commitInfo, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3})
		s.Require().NoError(err)

		req := s.createRequestPrepareProposal(
			commitInfo,
			proposal,
			3,
		)

		response, err := prepareHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Len(response.Txs, 1)
		s.Require().Equal(response.Txs[0], commitInfoBz)
	})

	s.Run("cannot build block with invalid a vote extension", func() {
		proposal := [][]byte{}

		valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, time.Now(), 2)
		s.Require().NoError(err)

		// Set the height of the third vote extension to 3, which is invalid.
		valVoteInfo3.VoteExtension = []byte("bad vote extension")

		commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3})
		s.Require().NoError(err)

		req := s.createRequestPrepareProposal(
			commitInfo,
			proposal,
			3,
		)

		_, err = prepareHandler(s.ctx, req)
		s.Require().Error(err)
	})

	s.Run("can reject a block with malformed prices", func() {
		proposal := [][]byte{}

		valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, malformedPrices, time.Now(), 2)
		s.Require().NoError(err)

		commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		req := s.createRequestPrepareProposal(
			commitInfo,
			proposal,
			3,
		)

		response, err := prepareHandler(s.ctx, req)
		s.Require().Error(err)
		s.Require().Len(response.Txs, 0)
	})
}

func (s *ProposalsTestSuite) TestProcessProposal() {
	processHandler := s.proposalHandler.ProcessProposalHandler()

	s.Run("can process an empty block when vote extensions are disabled", func() {
		proposal := [][]byte{}

		req := s.createRequestProcessProposal(proposal, 1)

		response, err := processHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_ACCEPT)
	})

	s.Run("can process a block with a single tx", func() {
		proposal := [][]byte{
			[]byte("tx1"),
		}

		req := s.createRequestProcessProposal(proposal, 1)

		response, err := processHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_ACCEPT)
	})

	// Enable vote extensions for the remaining tests.
	s.ctx = s.ctx.WithBlockHeight(3)

	s.Run("rejects a block with missing vote extensions", func() {
		proposal := [][]byte{
			[]byte("tx1"),
		}

		req := s.createRequestProcessProposal(proposal, 3)

		response, err := processHandler(s.ctx, req)
		s.Require().Error(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_REJECT)
	})

	s.Run("can process a block with a single vote extension", func() {
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		_, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{
			commitInfoBz,
			[]byte("tx1"),
		}

		req := s.createRequestProcessProposal(proposal, 3)

		response, err := processHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_ACCEPT)
	})

	s.Run("can process a block with multiple vote extensions", func() {
		valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, time.Now(), 2)
		s.Require().NoError(err)

		_, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3})
		s.Require().NoError(err)

		proposal := [][]byte{
			commitInfoBz,
			[]byte("tx1"),
		}

		req := s.createRequestProcessProposal(proposal, 3)

		response, err := processHandler(s.ctx, req)
		s.Require().NoError(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_ACCEPT)
	})

	s.Run("rejects a block with an invalid vote extension", func() {
		valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, time.Now(), 2)
		s.Require().NoError(err)

		valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, time.Now(), 2)
		s.Require().NoError(err)

		// Set the height of the third vote extension to 3, which is invalid.
		valVoteInfo3.VoteExtension = []byte("bad vote extension")

		_, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3})
		s.Require().NoError(err)

		proposal := [][]byte{
			commitInfoBz,
			[]byte("tx1"),
		}

		req := s.createRequestProcessProposal(proposal, 3)

		response, err := processHandler(s.ctx, req)
		s.Require().Error(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_REJECT)
	})

	s.Run("rejects a block with malformed prices", func() {
		valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, malformedPrices, time.Now(), 2)
		s.Require().NoError(err)

		_, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		s.Require().NoError(err)

		proposal := [][]byte{
			commitInfoBz,
			[]byte("tx1"),
		}

		req := s.createRequestProcessProposal(proposal, 3)

		response, err := processHandler(s.ctx, req)
		s.Require().Error(err)
		s.Require().Equal(response.Status, cometabci.ResponseProcessProposal_REJECT)
	})
}

func (s *ProposalsTestSuite) createRequestPrepareProposal(
	extendedCommitInfo cometabci.ExtendedCommitInfo,
	txs [][]byte,
	height int64,
) *cometabci.RequestPrepareProposal {
	return &cometabci.RequestPrepareProposal{
		Txs:             txs,
		LocalLastCommit: extendedCommitInfo,
		Height:          height,
	}
}

func (s *ProposalsTestSuite) createRequestProcessProposal(
	proposal [][]byte,
	height int64,
) *cometabci.RequestProcessProposal {
	return &cometabci.RequestProcessProposal{
		Txs:    proposal,
		Height: height,
	}
}
