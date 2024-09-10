package proposals_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/connect/v2/abci/proposals"
	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	codecmocks "github.com/skip-mev/connect/v2/abci/strategies/codec/mocks"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	currencypairmocks "github.com/skip-mev/connect/v2/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/connect/v2/abci/testutils"
	"github.com/skip-mev/connect/v2/abci/types"
	"github.com/skip-mev/connect/v2/abci/ve"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
	servicemetricsmocks "github.com/skip-mev/connect/v2/service/metrics/mocks"
)

var (
	noRizz = append([]byte("no_rizz"), make([]byte, 33)...)

	oneHundred   = big.NewInt(100)
	twoHundred   = big.NewInt(200)
	threeHundred = big.NewInt(300)
	fourHundred  = big.NewInt(400)

	prices1 = map[uint64][]byte{
		0: oneHundred.Bytes(),
	}
	prices2 = map[uint64][]byte{
		1: twoHundred.Bytes(),
	}
	prices3 = map[uint64][]byte{
		0: threeHundred.Bytes(),
		1: fourHundred.Bytes(),
	}
	malformedPrices = map[uint64][]byte{
		0: noRizz,
	}

	val1 = sdk.ConsAddress("val1")
	val2 = sdk.ConsAddress("val2")
	val3 = sdk.ConsAddress("val3")

	removeFirstTxn = sdk.PrepareProposalHandler(func(_ sdk.Context, proposal *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
		if len(proposal.Txs) > 0 {
			return &cometabci.ResponsePrepareProposal{Txs: proposal.Txs[1:]}, nil
		}
		return &cometabci.ResponsePrepareProposal{Txs: proposal.Txs}, nil
	})
)

type ProposalsTestSuite struct {
	suite.Suite

	ctx sdk.Context

	proposalHandler        *proposals.ProposalHandler
	prepareProposalHandler sdk.PrepareProposalHandler
	processProposalHandler sdk.ProcessProposalHandler
	codec                  codec.VoteExtensionCodec
	extCommitCodec         codec.ExtendedCommitCodec
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ProposalsTestSuite))
}

func (s *ProposalsTestSuite) SetupTest() {
	s.ctx = testutils.CreateBaseSDKContext(s.T())
	s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 1)
	s.ctx = s.ctx.WithBlockHeight(1)

	s.codec = codec.NewCompressionVoteExtensionCodec(
		codec.NewDefaultVoteExtensionCodec(),
		codec.NewZLibCompressor(),
	)

	s.extCommitCodec = codec.NewCompressionExtendedCommitCodec(
		codec.NewDefaultExtendedCommitCodec(),
		codec.NewZLibCompressor(),
	)
}

func (s *ProposalsTestSuite) TestPrepareProposal() {
	testCases := []struct {
		name                   string
		request                func() *cometabci.RequestPrepareProposal
		veEnabled              bool
		currencyPairStrategy   func() currencypair.CurrencyPairStrategy
		expectedProposalTxns   int
		expectedError          bool
		prepareProposalHandler *sdk.PrepareProposalHandler
	}{
		{
			name: "nil request returns an error",
			request: func() *cometabci.RequestPrepareProposal {
				return nil
			},
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedError: true,
		},
		{
			name: "vote extensions not enabled",
			request: func() *cometabci.RequestPrepareProposal {
				return s.createRequestPrepareProposal(
					cometabci.ExtendedCommitInfo{},
					nil,
					0,
				)
			},
			veEnabled: false,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedProposalTxns: 0,
			expectedError:        false,
		},
		{
			name: "vote extensions disabled with multiple txs",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{
					[]byte("tx1"),
					[]byte("tx2"),
				}

				return s.createRequestPrepareProposal(
					cometabci.ExtendedCommitInfo{},
					proposal,
					0,
				)
			},
			veEnabled: false,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedProposalTxns: 2,
			expectedError:        false,
		},
		{
			name: "vote extensions enabled with no txs and a single vote extension",
			request: func() *cometabci.RequestPrepareProposal {
				var proposal [][]byte

				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				return s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(0), nil).Once()
				return cpStrategy
			},
			expectedProposalTxns: 1,
			expectedError:        false,
		},
		{
			name: "vote extensions enabled with multiple txs and a single vote extension",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{
					[]byte("tx1"),
					[]byte("tx2"),
				}

				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				return s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedProposalTxns: 3,
			expectedError:        false,
		},
		{
			name: "vote extensions enabled with multiple vote extensions",
			request: func() *cometabci.RequestPrepareProposal {
				var proposal [][]byte

				valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, s.codec)
				s.Require().NoError(err)

				valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				return s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(2), nil).Times(3)
				return cpStrategy
			},
			expectedProposalTxns: 1,
			expectedError:        false,
		},
		{
			name: "cannot build block with invalid a vote extension - will be pruned",
			request: func() *cometabci.RequestPrepareProposal {
				var proposal [][]byte

				valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				// Set the height of the third vote extension to 3, which is invalid.
				valVoteInfo1.VoteExtension = []byte("bad vote extension")

				commitInfo, _, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo1},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				return s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedProposalTxns: 1,
			expectedError:        false,
		},
		{
			name: "can reject a block with malformed prices - will be pruned",
			request: func() *cometabci.RequestPrepareProposal {
				var proposal [][]byte

				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, malformedPrices, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				return s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedProposalTxns: 1,
			expectedError:        false,
		},
		{
			name: "can limit tx inclusion based on MaxTxBytes",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{
					[]byte("tx1"),
					[]byte("tx2"),
					make([]byte, 500),
				}

				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				prop := s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
				prop.MaxTxBytes = 500
				return prop
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedProposalTxns: 3,
			expectedError:        false,
		},
		{
			name: "can re-inject removed VE Txn",
			request: func() *cometabci.RequestPrepareProposal {
				var proposal [][]byte

				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				prop := s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
				prop.MaxTxBytes = 500
				return prop
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedProposalTxns:   1,
			expectedError:          false,
			prepareProposalHandler: &removeFirstTxn,
		},
		{
			name: "will fail if VE Txn is too large",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{
					[]byte("one"),
					[]byte("two"),
				}

				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				commitInfo, _, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				prop := s.createRequestPrepareProposal(
					commitInfo,
					proposal,
					3,
				)
				prop.MaxTxBytes = 40
				return prop
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedProposalTxns: 0,
			expectedError:        true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.prepareProposalHandler = baseapp.NoOpPrepareProposal()
			if tc.prepareProposalHandler != nil {
				s.prepareProposalHandler = *tc.prepareProposalHandler
			}
			s.processProposalHandler = baseapp.NoOpProcessProposal()
			s.proposalHandler = proposals.NewProposalHandler(
				log.NewTestLogger(s.T()),
				s.prepareProposalHandler,
				s.processProposalHandler,
				ve.NoOpValidateVoteExtensions,
				s.codec,
				s.extCommitCodec,
				tc.currencyPairStrategy(),
				servicemetrics.NewNopMetrics(),
			)

			if tc.veEnabled {
				s.ctx = s.ctx.WithBlockHeight(3)
			}

			req := tc.request()
			response, err := s.proposalHandler.PrepareProposalHandler()(s.ctx, req)
			if tc.expectedError {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)

			s.Require().Equal(tc.expectedProposalTxns, len(response.Txs))

			if tc.veEnabled {
				bz, err := s.extCommitCodec.Encode(req.LocalLastCommit)
				s.Require().NoError(err)
				if int64(len(bz)) < req.MaxTxBytes {
					s.Require().Equal(response.Txs[0], bz)
				}
			}
		})
	}
}

func (s *ProposalsTestSuite) TestPrepareProposalRetainOracleData() {
	// If retain option is given we feed oracle-data into prepare / process
	s.Run("test RetainOracleDataInWrappedProposalHandler", func() {
		exCodec := codec.NewDefaultExtendedCommitCodec()
		veCodec := codec.NewDefaultVoteExtensionCodec()

		cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(0), nil).Once()

		emptyVote, err := testutils.CreateExtendedVoteInfo(val1, map[uint64][]byte{}, veCodec)
		s.Require().NoError(err)

		extendedCommit := cometabci.ExtendedCommitInfo{
			Round: 1,
			Votes: []cometabci.ExtendedVoteInfo{
				emptyVote,
			},
		}
		bz, err := exCodec.Encode(extendedCommit)
		s.Require().NoError(err)

		handler := proposals.NewProposalHandler(
			log.NewNopLogger(),
			func(_ sdk.Context, rpp *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				// assert that the oracle data is retained
				s.Require().Equal(bz, rpp.Txs[types.OracleInfoIndex])
				return &cometabci.ResponsePrepareProposal{}, nil
			},
			func(_ sdk.Context, rpp *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				// assert that the oracle data is retained
				s.Require().Equal(bz, rpp.Txs[types.OracleInfoIndex])
				return &cometabci.ResponseProcessProposal{}, nil
			},
			ve.NoOpValidateVoteExtensions,
			veCodec,
			exCodec,
			cpStrategy,
			servicemetrics.NewNopMetrics(),
			proposals.RetainOracleDataInWrappedProposalHandler(),
		)

		// enable VE
		s.ctx = testutils.CreateBaseSDKContext(s.T())
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 3)
		s.ctx = s.ctx.WithBlockHeight(4)

		// prepare proposal
		req := &cometabci.RequestPrepareProposal{
			LocalLastCommit: extendedCommit,
			MaxTxBytes:      100, // arbitrary
		}

		_, err = handler.PrepareProposalHandler()(s.ctx, req)
		s.Require().NoError(err)

		// process proposal
		req2 := &cometabci.RequestProcessProposal{
			ProposedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: nil,
			},
			Txs: [][]byte{bz},
		}
		_, err = handler.ProcessProposalHandler()(s.ctx, req2)
		s.Require().NoError(err)
	})

	// Otherwise, we don't
	s.Run("test that oracle-data is not passed if not RetainOracleDataInWrappedProposalHandler", func() {
		exCodec := codec.NewDefaultExtendedCommitCodec()
		veCodec := codec.NewDefaultVoteExtensionCodec()

		cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(0), nil).Once()

		emptyVote, err := testutils.CreateExtendedVoteInfo(val1, map[uint64][]byte{}, veCodec)
		s.Require().NoError(err)

		extendedCommit := cometabci.ExtendedCommitInfo{
			Round: 1,
			Votes: []cometabci.ExtendedVoteInfo{
				emptyVote,
			},
		}
		bz, err := exCodec.Encode(extendedCommit)
		s.Require().NoError(err)

		handler := proposals.NewProposalHandler(
			log.NewNopLogger(),
			func(_ sdk.Context, rpp *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				// assert that the oracle data is retained
				s.Require().Len(rpp.Txs, 0)
				return &cometabci.ResponsePrepareProposal{}, nil
			},
			func(_ sdk.Context, rpp *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				// assert that the oracle data is retained
				s.Require().Len(rpp.Txs, 0)
				return &cometabci.ResponseProcessProposal{}, nil
			},
			ve.NoOpValidateVoteExtensions,
			veCodec,
			exCodec,
			cpStrategy,
			servicemetrics.NewNopMetrics(),
		)

		// enable VE
		s.ctx = testutils.CreateBaseSDKContext(s.T())
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 3)
		s.ctx = s.ctx.WithBlockHeight(4)

		// prepare proposal
		req := &cometabci.RequestPrepareProposal{
			LocalLastCommit: extendedCommit,
			MaxTxBytes:      100, // arbitrary
		}

		_, err = handler.PrepareProposalHandler()(s.ctx, req)
		s.Require().NoError(err)

		// process proposal
		req2 := &cometabci.RequestProcessProposal{
			ProposedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: nil,
			},
			Txs: [][]byte{bz},
		}
		_, err = handler.ProcessProposalHandler()(s.ctx, req2)
		s.Require().NoError(err)
	})
}

func (s *ProposalsTestSuite) TestProcessProposal() {
	testCases := []struct {
		name                 string
		request              func() *cometabci.RequestProcessProposal
		veEnabled            bool
		lastCommit           cometabci.CommitInfo
		currencyPairStrategy func() currencypair.CurrencyPairStrategy
		expectedError        bool
		expectedResp         *cometabci.ResponseProcessProposal
		checkTxs             func(before, after [][]byte)
	}{
		{
			name: "returns an error on nil request",
			request: func() *cometabci.RequestProcessProposal {
				return nil
			},
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedError: true,
		},
		{
			name: "can process any empty block when vote extensions are disabled",
			request: func() *cometabci.RequestProcessProposal {
				return s.createRequestProcessProposal(
					[][]byte{},
					cometabci.CommitInfo{},
					1,
				)
			},
			veEnabled: false,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedError: false,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_ACCEPT,
			},
		},
		{
			name: "can process a block with a single tx",
			request: func() *cometabci.RequestProcessProposal {
				proposal := [][]byte{
					[]byte("tx1"),
				}

				lastCommit := cometabci.CommitInfo{
					Round: 3,
					Votes: []cometabci.VoteInfo{},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,
					1,
				)
			},
			veEnabled: false,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedError: false,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_ACCEPT,
			},
		},
		{
			name: "rejects a block with missing vote extensions",
			request: func() *cometabci.RequestProcessProposal {
				var proposal [][]byte

				lastCommit := cometabci.CommitInfo{
					Round: 3,
					Votes: []cometabci.VoteInfo{},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedError: true,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_REJECT,
			},
		},
		{
			name: "can process a block with a single vote extension",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				ext, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				lastCommit := cometabci.CommitInfo{
					Round: ext.Round,
					Votes: []cometabci.VoteInfo{
						{
							Validator:   valVoteInfo.Validator,
							BlockIdFlag: valVoteInfo.BlockIdFlag,
						},
					},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedError: false,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_ACCEPT,
			},
			checkTxs: func(before, after [][]byte) {
				s.Require().Equal(before, after)
			},
		},
		{
			name: "can process a block with multiple vote extensions",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, s.codec)
				s.Require().NoError(err)

				valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, s.codec)
				s.Require().NoError(err)

				ext, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				lastCommit := cometabci.CommitInfo{
					Round: ext.Round,
					Votes: []cometabci.VoteInfo{
						{
							Validator:   valVoteInfo1.Validator,
							BlockIdFlag: valVoteInfo1.BlockIdFlag,
						},
						{
							Validator:   valVoteInfo2.Validator,
							BlockIdFlag: valVoteInfo2.BlockIdFlag,
						},
						{
							Validator:   valVoteInfo3.Validator,
							BlockIdFlag: valVoteInfo3.BlockIdFlag,
						},
					},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,

					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(2), nil).Times(3)
				return cpStrategy
			},
			expectedError: false,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_ACCEPT,
			},
		},
		{
			name: "can process a block with valid pruned vote extension",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				valVoteInfo2, err := testutils.CreateExtendedVoteInfo(val2, prices2, s.codec)
				s.Require().NoError(err)

				valVoteInfo3, err := testutils.CreateExtendedVoteInfo(val3, prices3, s.codec)
				s.Require().NoError(err)

				valVoteInfo3Modified := valVoteInfo3
				valVoteInfo3Modified.BlockIdFlag = cometproto.BlockIDFlagAbsent
				valVoteInfo3Modified.VoteExtension = nil
				valVoteInfo3Modified.ExtensionSignature = nil

				ext, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3Modified},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				lastCommit := cometabci.CommitInfo{
					Round: ext.Round,
					Votes: []cometabci.VoteInfo{
						{
							Validator:   valVoteInfo1.Validator,
							BlockIdFlag: valVoteInfo1.BlockIdFlag,
						},
						{
							Validator:   valVoteInfo2.Validator,
							BlockIdFlag: valVoteInfo2.BlockIdFlag,
						},
						{
							Validator:   valVoteInfo3.Validator,
							BlockIdFlag: valVoteInfo3.BlockIdFlag,
						},
					},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,

					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(2), nil).Times(2)
				return cpStrategy
			},
			expectedError: false,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_ACCEPT,
			},
		},
		{
			name: "rejects a block with an invalid vote extension",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				valVoteInfo.VoteExtension = []byte("bad vote extension")

				ext, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				lastCommit := cometabci.CommitInfo{
					Round: ext.Round,
					Votes: []cometabci.VoteInfo{
						{
							Validator:   valVoteInfo.Validator,
							BlockIdFlag: valVoteInfo.BlockIdFlag,
						},
					},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedError: true,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_REJECT,
			},
		},
		{
			name: "rejects a block with malformed prices",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, malformedPrices, s.codec)
				s.Require().NoError(err)

				ext, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				lastCommit := cometabci.CommitInfo{
					Round: ext.Round,
					Votes: []cometabci.VoteInfo{
						{
							Validator:   valVoteInfo.Validator,
							BlockIdFlag: valVoteInfo.BlockIdFlag,
						},
					},
				}

				return s.createRequestProcessProposal(
					proposal,
					lastCommit,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
				cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Once()
				return cpStrategy
			},
			expectedError: true,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_REJECT,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.prepareProposalHandler = baseapp.NoOpPrepareProposal()
			s.processProposalHandler = baseapp.NoOpProcessProposal()
			s.proposalHandler = proposals.NewProposalHandler(
				log.NewTestLogger(s.T()),
				s.prepareProposalHandler,
				s.processProposalHandler,
				ValidateVoteExtensionsAgainstLastCommit,
				s.codec,
				s.extCommitCodec,
				tc.currencyPairStrategy(),
				servicemetrics.NewNopMetrics(),
			)

			req := tc.request()
			if tc.veEnabled {
				s.ctx = s.ctx.WithBlockHeight(3).WithCometInfo(baseapp.NewBlockInfo(
					nil,
					nil,
					nil,
					req.ProposedLastCommit,
				))
			}

			// make a copy of the txs before we run the proposal
			var before [][]byte
			if req != nil { // some tests use a nil request.
				before = make([][]byte, len(req.Txs))
				for i, tx := range req.Txs {
					before[i] = make([]byte, len(tx))
					copy(before[i], tx)
				}
			}

			response, err := s.proposalHandler.ProcessProposalHandler()(s.ctx, req)

			s.Require().Equal(tc.expectedResp, response)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			if tc.checkTxs != nil {
				tc.checkTxs(before, req.Txs)
			}
		})
	}
}

func ValidateVoteExtensionsAgainstLastCommit(
	ctx sdk.Context,
	extCommit cometabci.ExtendedCommitInfo,
) error {
	commitInfo := ctx.CometInfo().GetLastCommit()
	return ve.ValidateExtendedCommitAgainstLastCommit(extCommit, commitInfo)
}

func (s *ProposalsTestSuite) TestProposalLatency() {
	// check that no latency is reported for a failed PrepareProposal
	metricsMocks := servicemetricsmocks.NewMetrics(s.T())
	cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())
	cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(0), nil).Once()

	// check that latency reported in upstream logic is reported
	s.Run("wrapped prepare proposal latency is reported", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context, _ *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				// simulate a long-running prepare proposal
				time.Sleep(200 * time.Millisecond)
				return &cometabci.ResponsePrepareProposal{
					Txs: nil,
				}, nil
			},
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			cpStrategy,
			metricsMocks,
		)

		vote, err := testutils.CreateExtendedVoteInfo(val1, map[uint64][]byte{}, codec.NewDefaultVoteExtensionCodec())
		s.Require().NoError(err)

		req := s.createRequestPrepareProposal( // the votes here are invalid, but that's fine
			cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{vote},
			},
			nil,
			4, // vote extensions will be enabled
		)

		s.ctx = s.ctx.WithBlockHeight(4)
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return().Run(func(args mock.Arguments) {
			// the second arg should be a duration
			latency := args.Get(1).(time.Duration)
			s.Require().True(latency >= 100*time.Millisecond) // should have included latency from validate vote extensions
			s.Require().True(latency < 300*time.Millisecond)  // should have ignored wrapped prepare-proposal latency
		}).Once()
		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, servicemetrics.Success{}).Once()

		_, err = propHandler.PrepareProposalHandler()(s.ctx, req)
		s.Require().NoError(err)
	})

	s.Run("prepare proposal latency is reported in the case of failures", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				time.Sleep(100 * time.Millisecond)
				return fmt.Errorf("error in validate vote extensions")
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)

		req := s.createRequestPrepareProposal( // the votes here are invalid, but that's fine
			cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: nil,
			},
			nil,
			4, // vote extensions will be enabled
		)

		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return().Run(func(args mock.Arguments) {
			// the second arg should be a duration
			latency := args.Get(1).(time.Duration)
			s.Require().True(latency >= 100*time.Millisecond) // should have included latency from validate vote extensions
		}).Once()

		expErr := proposals.InvalidExtendedCommitInfoError{
			Err: fmt.Errorf("error in validate vote extensions"),
		}
		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, expErr).Once()
		_, err := propHandler.PrepareProposalHandler()(s.ctx, req)
		s.Require().Error(err, expErr)
	})

	s.Run("wrapped process proposal latency is reported", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			baseapp.NoOpPrepareProposal(),
			func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				// simulate a long-running process proposal
				time.Sleep(200 * time.Millisecond)
				return &cometabci.ResponseProcessProposal{}, nil
			},
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				// simulate a long-running validate vote extensions
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)

		ext, extInfoBz, err := testutils.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{},
			codec.NewDefaultExtendedCommitCodec(),
		)
		s.Require().NoError(err)

		lastCommit := cometabci.CommitInfo{
			Round: ext.Round,
			Votes: []cometabci.VoteInfo{},
		}

		req := s.createRequestProcessProposal([][]byte{extInfoBz}, lastCommit, 4)
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return().Run(func(args mock.Arguments) {
			// the second arg should be a duration
			latency := args.Get(1).(time.Duration)
			s.Require().True(latency >= 100*time.Millisecond) // should have included validate vote extensions latency
			s.Require().True(latency < 300*time.Millisecond)  // should have ignored the wrapped processProposal latency
		}).Once()
		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, servicemetrics.Success{}).Once()
		metricsMocks.On("ObserveMessageSize", servicemetrics.ExtendedCommit, mock.Anything)

		_, err = propHandler.ProcessProposalHandler()(s.ctx, req)
		s.Require().NoError(err)
	})

	s.Run("process proposal latency is reported in the case of failures", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				time.Sleep(100 * time.Millisecond)
				return fmt.Errorf("error in validate vote extensions")
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)

		ext, extInfoBz, err := testutils.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{},
			codec.NewDefaultExtendedCommitCodec(),
		)
		s.Require().NoError(err)

		expErr := proposals.InvalidExtendedCommitInfoError{
			Err: fmt.Errorf("error in validate vote extensions"),
		}

		lastCommit := cometabci.CommitInfo{
			Round: ext.Round,
			Votes: []cometabci.VoteInfo{},
		}

		req := s.createRequestProcessProposal([][]byte{extInfoBz}, lastCommit, 4)
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return().Run(func(args mock.Arguments) {
			// the second arg should be a duration
			latency := args.Get(1).(time.Duration)
			s.Require().True(latency >= 100*time.Millisecond) // should have included validate vote extensions latency
		}).Once()
		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, expErr).Once()
		metricsMocks.On("ObserveMessageSize", servicemetrics.ExtendedCommit, mock.Anything)

		_, err = propHandler.ProcessProposalHandler()(s.ctx, req)
		s.Require().Error(err, expErr)
	})
}

func (s *ProposalsTestSuite) TestPrepareProposalStatus() {
	// test nil request
	s.Run("test nil request", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context, _ *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				return nil, nil
			},
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				return nil
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)

		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, types.NilRequestError{
			Handler: servicemetrics.PrepareProposal,
		}).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return()

		_, err := propHandler.PrepareProposalHandler()(s.ctx, nil)
		s.Require().Error(err, types.NilRequestError{
			Handler: servicemetrics.PrepareProposal,
		})
	})
	// test failing wrapped prepare proposal
	s.Run("test failing wrapped prepare proposal", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		prepareErr := fmt.Errorf("error in prepare proposal")
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context, _ *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				return nil, prepareErr
			},
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				return nil
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)
		expErr := types.WrappedHandlerError{
			Handler: servicemetrics.PrepareProposal,
			Err:     prepareErr,
		}
		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, expErr).Once()

		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return()
		// make vote-extensions not enabled to skip validate vote extensions
		s.ctx = s.ctx.WithBlockHeight(1)
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 3)

		_, err := propHandler.PrepareProposalHandler()(s.ctx, &cometabci.RequestPrepareProposal{})
		s.Require().Error(err, expErr)
	})

	// test invalid extended commit
	s.Run("test invalid extended commit", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		extCommitError := fmt.Errorf("error in validating extended commit")
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context, _ *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				return nil, nil
			},
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				return extCommitError
			},
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)
		expErr := proposals.InvalidExtendedCommitInfoError{
			Err: extCommitError,
		}
		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, expErr).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return()

		// make vote-extensions enabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 1)
		s.ctx = s.ctx.WithBlockHeight(4)

		_, err := propHandler.PrepareProposalHandler()(s.ctx, &cometabci.RequestPrepareProposal{})
		s.Require().Error(err, expErr)
	})

	// test codec failure
	s.Run("test invalid extended commit", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())
		codecErr := fmt.Errorf("error in codec")
		c := codecmocks.NewExtendedCommitCodec(s.T())
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context, _ *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				return nil, nil
			},
			nil,
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				return nil
			},
			codec.NewDefaultVoteExtensionCodec(),
			c,
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsMocks,
		)
		expErr := types.CodecError{
			Err: codecErr,
		}
		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, expErr).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return()

		c.On("Encode", mock.Anything).Return(nil, codecErr)

		// make vote-extensions enabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 1)
		s.ctx = s.ctx.WithBlockHeight(4)

		_, err := propHandler.PrepareProposalHandler()(s.ctx, &cometabci.RequestPrepareProposal{})
		s.Require().Error(err, expErr)
	})

	// test success
	s.Run("test success", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(_ sdk.Context, _ *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				return &cometabci.ResponsePrepareProposal{}, nil
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			metricsMocks,
		)

		metricsMocks.On("AddABCIRequest", servicemetrics.PrepareProposal, servicemetrics.Success{}).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return()

		// make vote-extensions enabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 4)
		s.ctx = s.ctx.WithBlockHeight(1)

		_, err := propHandler.PrepareProposalHandler()(s.ctx, &cometabci.RequestPrepareProposal{})
		s.Require().NoError(err)
	})
}

func (s *ProposalsTestSuite) TestProcessProposalStatus() {
	// test nil request
	s.Run("test nil request", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			metricsMocks,
		)

		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, types.NilRequestError{
			Handler: servicemetrics.ProcessProposal,
		}).Once()

		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()

		_, err := propHandler.ProcessProposalHandler()(s.ctx, nil)
		s.Require().Error(err, types.NilRequestError{
			Handler: servicemetrics.ProcessProposal,
		})
	})
	// test failed wrapped process-proposal
	s.Run("test failed wrapped process-proposal", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		processErr := fmt.Errorf("error in process")
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				return nil, processErr
			},
			nil,
			nil,
			nil,
			nil,
			metricsMocks,
		)
		expErr := types.WrappedHandlerError{
			Handler: servicemetrics.ProcessProposal,
			Err:     processErr,
		}

		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, expErr).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()

		// make vote-extensions disabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 4)
		s.ctx = s.ctx.WithBlockHeight(3)

		_, err := propHandler.ProcessProposalHandler()(s.ctx, &cometabci.RequestProcessProposal{})
		s.Require().Error(err, expErr)
	})
	// test success
	s.Run("test success", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				return nil, nil
			},
			nil,
			nil,
			nil,
			nil,
			metricsMocks,
		)

		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, servicemetrics.Success{}).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()

		// make vote-extensions disabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 4)
		s.ctx = s.ctx.WithBlockHeight(3)

		_, err := propHandler.ProcessProposalHandler()(s.ctx, &cometabci.RequestProcessProposal{})
		s.Require().NoError(err)
	})
	// test failing w/ missing commit info
	s.Run("test failing w/ missing commit info", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())

		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				return nil, nil
			},
			nil,
			nil,
			nil,
			nil,
			metricsMocks,
		)
		expErr := types.MissingCommitInfoError{}
		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, expErr).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()

		// make vote-extensions disabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(3)

		_, err := propHandler.ProcessProposalHandler()(s.ctx, &cometabci.RequestProcessProposal{})
		s.Require().Error(err, expErr)
	})
	// test codec failure
	s.Run("test codec failure", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())
		codecErr := fmt.Errorf("error in codec")
		c := codecmocks.NewExtendedCommitCodec(s.T())
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				return nil, nil
			},
			nil,
			nil,
			c,
			nil,
			metricsMocks,
		)
		expErr := types.CodecError{
			Err: codecErr,
		}
		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, expErr).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()

		c.On("Decode", mock.Anything).Return(cometabci.ExtendedCommitInfo{}, codecErr)

		// make vote-extensions disabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(3)

		_, err := propHandler.ProcessProposalHandler()(s.ctx, &cometabci.RequestProcessProposal{
			Txs: [][]byte{{1, 2, 3}},
		})
		s.Require().Error(err, expErr)
	})
	// test invalid extended commit
	s.Run("test invalid extended commit", func() {
		metricsMocks := servicemetricsmocks.NewMetrics(s.T())
		validateErr := fmt.Errorf("error in validateExtendedCommit")
		c := codecmocks.NewExtendedCommitCodec(s.T())
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			nil,
			func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				return nil, nil
			},
			func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
				return validateErr
			},
			nil,
			c,
			nil,
			metricsMocks,
		)
		expErr := proposals.InvalidExtendedCommitInfoError{
			Err: validateErr,
		}
		metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, expErr).Once()
		metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()
		c.On("Decode", mock.Anything).Return(cometabci.ExtendedCommitInfo{}, nil)

		// make vote-extensions disabled
		s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
		s.ctx = s.ctx.WithBlockHeight(3)

		_, err := propHandler.ProcessProposalHandler()(s.ctx, &cometabci.RequestProcessProposal{
			Txs: [][]byte{{1, 2, 3}},
		})
		s.Require().Error(err, expErr)
	})
}

func (s *ProposalsTestSuite) TestExtendedCommitSize() {
	metricsMocks := servicemetricsmocks.NewMetrics(s.T())
	cdc := codecmocks.NewExtendedCommitCodec(s.T())

	propHandler := proposals.NewProposalHandler(
		log.NewTestLogger(s.T()),
		nil,
		func(_ sdk.Context, _ *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
			return nil, nil
		},
		func(_ sdk.Context, _ cometabci.ExtendedCommitInfo) error {
			return nil
		},
		nil,
		cdc,
		nil,
		metricsMocks,
	)

	extendedCommit := make([]byte, 100)

	metricsMocks.On("AddABCIRequest", servicemetrics.ProcessProposal, servicemetrics.Success{}).Once()
	metricsMocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return()
	metricsMocks.On("ObserveMessageSize", servicemetrics.ExtendedCommit, 100)

	// mock codec
	cdc.On("Decode", extendedCommit).Return(cometabci.ExtendedCommitInfo{
		Votes: []cometabci.ExtendedVoteInfo{},
	}, nil)

	// make vote-extensions enabled
	s.ctx = testutils.UpdateContextWithVEHeight(s.ctx, 2)
	s.ctx = s.ctx.WithBlockHeight(3)

	_, err := propHandler.ProcessProposalHandler()(s.ctx, &cometabci.RequestProcessProposal{
		Txs: [][]byte{extendedCommit},
	})
	s.Require().NoError(err)
}

func (s *ProposalsTestSuite) TestValidateExtendedCommitInfoProcess() {
	s.Run("should fail for nil request", func() {
	})
}

func (s *ProposalsTestSuite) TestPruning() {
	cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())

	ph := proposals.NewProposalHandler(
		log.NewNopLogger(),
		func(sdk.Context, *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
			return &cometabci.ResponsePrepareProposal{
				Txs: [][]byte{{1, 2, 3}},
			}, nil
		},
		func(sdk.Context, *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
			return nil, nil
		},
		ve.NoOpValidateVoteExtensions,
		codec.NewDefaultVoteExtensionCodec(),
		codec.NewDefaultExtendedCommitCodec(),
		cpStrategy,
		servicemetrics.NewNopMetrics(),
	)

	s.Run("no invalid votes to be pruned", func() {
		ve1, err := testutils.CreateExtendedVoteInfoWithPower(
			val1,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve2, err := testutils.CreateExtendedVoteInfoWithPower(
			val2,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve3 := cometabci.ExtendedVoteInfo{
			BlockIdFlag: cometproto.BlockIDFlagNil,
			Validator: cometabci.Validator{
				Address: val3,
				Power:   30,
			},
		}

		// mocks
		ctx := testutils.UpdateContextWithVEHeight(s.ctx, 2)
		ctx = ctx.WithBlockHeight(3)

		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Twice()

		extInfo, err := ph.PruneAndValidateExtendedCommitInfo(ctx, cometabci.ExtendedCommitInfo{
			Votes: []cometabci.ExtendedVoteInfo{ve1, ve2, ve3},
		})
		s.Require().NoError(err)
		s.Require().Len(extInfo.Votes, 3)

		// check that the votes are in the same order
		s.Require().Equal(ve1, extInfo.Votes[0])
		s.Require().Equal(ve2, extInfo.Votes[1])
		s.Require().Equal(ve3, extInfo.Votes[2])
	})

	s.Run("invalid votes to be pruned", func() {
		ve1, err := testutils.CreateExtendedVoteInfoWithPower(
			val1,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve2, err := testutils.CreateExtendedVoteInfoWithPower(
			val2,
			1,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
				2: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve3, err := testutils.CreateExtendedVoteInfoWithPower(
			val3,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		// mocks
		ctx := testutils.UpdateContextWithVEHeight(s.ctx, 2)
		ctx = ctx.WithBlockHeight(3)

		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Times(3)

		extInfo, err := ph.PruneAndValidateExtendedCommitInfo(ctx, cometabci.ExtendedCommitInfo{
			Votes: []cometabci.ExtendedVoteInfo{ve1, ve2, ve3},
		})
		s.Require().NoError(err)
		s.Require().Len(extInfo.Votes, 3)

		// check that the votes are in the same order
		s.Require().Equal(ve1, extInfo.Votes[0])
		updatedVe := extInfo.Votes[1]
		s.Require().Equal(0, len(updatedVe.ExtensionSignature))
		s.Require().Equal(0, len(updatedVe.VoteExtension))
		s.Require().Equal(cometproto.BlockIDFlagAbsent, updatedVe.BlockIdFlag)
		s.Require().Equal(ve3, extInfo.Votes[2])
	})

	s.Run("pruning votes results in lack of super-majority", func() {
		ve1, err := testutils.CreateExtendedVoteInfoWithPower(
			val1,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve2, err := testutils.CreateExtendedVoteInfoWithPower(
			val2,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
				2: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve3, err := testutils.CreateExtendedVoteInfoWithPower(
			val3,
			1,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		// mocks
		ctx := testutils.UpdateContextWithVEHeight(s.ctx, 2)
		ctx = ctx.WithBlockHeight(3)

		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Times(3)

		extInfo, err := ph.PruneAndValidateExtendedCommitInfo(ctx, cometabci.ExtendedCommitInfo{
			Votes: []cometabci.ExtendedVoteInfo{ve1, ve2, ve3},
		})
		s.Require().NoError(err)
		s.Require().Len(extInfo.Votes, 3)

		// ensure that voting power is now invalid
		s.Require().Error(s.checkVotingPowerValid(extInfo))
	})

	s.Run("prepare-proposal w/ invalid ves in LastCommit", func() {
		ve1, err := testutils.CreateExtendedVoteInfoWithPower(
			val1,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve2, err := testutils.CreateExtendedVoteInfoWithPower(
			val2,
			1,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
				2: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		ve3, err := testutils.CreateExtendedVoteInfoWithPower(
			val3,
			33,
			map[uint64][]byte{
				1: twoHundred.Bytes(),
			},
			codec.NewDefaultVoteExtensionCodec(),
		)
		s.Require().NoError(err)

		// mocks
		ctx := testutils.UpdateContextWithVEHeight(s.ctx, 2)
		ctx = ctx.WithBlockHeight(3)

		cpStrategy.On("GetMaxNumCP", mock.Anything).Return(uint64(1), nil).Times(3)

		res, err := ph.PrepareProposalHandler()(ctx, s.createRequestPrepareProposal(cometabci.ExtendedCommitInfo{
			Votes: []cometabci.ExtendedVoteInfo{ve1, ve2, ve3},
		}, [][]byte{}, 3))
		s.Require().NoError(err)

		// check that the first tx is an extended, commit constructed as expected
		extInfo, err := codec.NewDefaultExtendedCommitCodec().Decode(res.Txs[0])
		s.Require().NoError(err)

		s.Require().Len(extInfo.Votes, 3)

		// check that the votes are in the same order
		s.Require().Equal(ve1, extInfo.Votes[0])
		updatedVe := extInfo.Votes[1]
		s.Require().Equal(0, len(updatedVe.ExtensionSignature))
		s.Require().Equal(0, len(updatedVe.VoteExtension))
		s.Require().Equal(cometproto.BlockIDFlagAbsent, updatedVe.BlockIdFlag)
		s.Require().Equal(ve3, extInfo.Votes[2])
	})
}

func (s *ProposalsTestSuite) createRequestPrepareProposal(
	extendedCommitInfo cometabci.ExtendedCommitInfo,
	txs [][]byte,
	height int64,
) *cometabci.RequestPrepareProposal {
	s.T().Helper()

	return &cometabci.RequestPrepareProposal{
		Txs:             txs,
		LocalLastCommit: extendedCommitInfo,
		Height:          height,
		// Use the same default MaxTxBytes that comet does
		MaxTxBytes: cmttypes.DefaultBlockParams().MaxBytes,
	}
}

func (s *ProposalsTestSuite) createRequestProcessProposal(
	proposal [][]byte,
	lastCommit cometabci.CommitInfo,
	height int64,
) *cometabci.RequestProcessProposal {
	s.T().Helper()

	return &cometabci.RequestProcessProposal{
		Txs:                proposal,
		ProposedLastCommit: lastCommit,
		Height:             height,
	}
}

func (s *ProposalsTestSuite) checkVotingPowerValid(
	extCommit cometabci.ExtendedCommitInfo,
) error {
	s.T().Helper()

	var (
		// Total voting power of all vote extensions.
		totalVP int64
		// Total voting power of all validators that submitted valid vote extensions.
		sumVP int64
	)

	for _, vote := range extCommit.Votes {
		totalVP += vote.Validator.Power

		// Only check + include power if the vote is a commit vote. There must be super-majority, otherwise the
		// previous block (the block vote is for) could not have been committed.
		if vote.BlockIdFlag != cometproto.BlockIDFlagCommit {
			continue
		}

		sumVP += vote.Validator.Power
	}

	// This check is probably unnecessary, but better safe than sorry.
	if totalVP <= 0 {
		return fmt.Errorf("total voting power must be positive, got: %d", totalVP)
	}

	// If the sum of the voting power has not reached (2/3 + 1) we need to error.
	if requiredVP := ((totalVP * 2) / 3) + 1; sumVP < requiredVP {
		return fmt.Errorf(
			"insufficient cumulative voting power received to verify vote extensions; got: %d, expected: >=%d",
			sumVP, requiredVP,
		)
	}

	return nil
}
