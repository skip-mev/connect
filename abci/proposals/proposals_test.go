package proposals_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	servicemetrics "github.com/skip-mev/slinky/service/metrics"

	"github.com/cometbft/cometbft/types"

	"cosmossdk.io/log"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/abci/proposals"
	"github.com/skip-mev/slinky/abci/strategies/codec"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	currencypairmocks "github.com/skip-mev/slinky/abci/strategies/currencypair/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/abci/ve"
	servicemetricsmocks "github.com/skip-mev/slinky/service/metrics/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	noRizz = append([]byte("no_rizz"), make([]byte, 33)...)

	oneHundred   = big.NewInt(100)
	twoHundred   = big.NewInt(200)
	threeHundred = big.NewInt(300)
	fourHundred  = big.NewInt(400)

	btcUSD = oracletypes.CurrencyPair{ // id 0
		Base:  "BTC",
		Quote: "USD",
	}

	ethUSD = oracletypes.CurrencyPair{ // id 1
		Base:  "ETH",
		Quote: "USD",
	}

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
				proposal := [][]byte{}

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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				return cpStrategy
			},
			expectedProposalTxns: 3,
			expectedError:        false,
		},
		{
			name: "vote extensions enabled with multiple vote extensions",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{}

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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				cpStrategy.On("FromID", mock.Anything, uint64(1)).Return(ethUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, ethUSD, mock.Anything).Return(big.NewInt(20), nil)

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				cpStrategy.On("FromID", mock.Anything, uint64(1)).Return(ethUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, ethUSD, mock.Anything).Return(big.NewInt(20), nil)

				return cpStrategy
			},
			expectedProposalTxns: 1,
			expectedError:        false,
		},
		{
			name: "cannot build block with invalid a vote extension",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{}

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
			expectedProposalTxns: 0,
			expectedError:        true,
		},
		{
			name: "can reject a block with malformed prices",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{}

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
				return currencypairmocks.NewCurrencyPairStrategy(s.T())
			},
			expectedProposalTxns: 1,
			expectedError:        true,
		},
		{
			name: "can reject a request with ve that contains invalid currency pair id",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{}

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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, fmt.Errorf("no rizz error ha"))

				return cpStrategy
			},
			expectedProposalTxns: 1,
			expectedError:        true,
		},
		{
			name: "can reject a request with ve that contains invalid price bytes",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{}

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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), fmt.Errorf(">:("))

				return cpStrategy
			},
			expectedProposalTxns: 1,
			expectedError:        true,
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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				return cpStrategy
			},
			expectedProposalTxns: 3,
			expectedError:        false,
		},
		{
			name: "can re-inject removed VE Txn",
			request: func() *cometabci.RequestPrepareProposal {
				proposal := [][]byte{}

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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				return cpStrategy
			},
			expectedProposalTxns:   1,
			expectedError:          false,
			prepareProposalHandler: &removeFirstTxn,
		},
		{
			name: "can exclude VE Txn when it's too large",
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

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				return cpStrategy
			},
			expectedProposalTxns: 2,
			expectedError:        false,
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

			s.Require().Equal(len(response.Txs), tc.expectedProposalTxns)

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

func (s *ProposalsTestSuite) TestProcessProposal() {
	testCases := []struct {
		name                 string
		request              func() *cometabci.RequestProcessProposal
		veEnabled            bool
		currencyPairStrategy func() currencypair.CurrencyPairStrategy
		expectedError        bool
		expectedResp         *cometabci.ResponseProcessProposal
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

				return s.createRequestProcessProposal(
					proposal,
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
				proposal := [][]byte{}

				return s.createRequestProcessProposal(
					proposal,
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

				_, commitInfoBz, err := testutils.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo}, s.extCommitCodec)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				return s.createRequestProcessProposal(
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				return cpStrategy
			},
			expectedError: false,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_ACCEPT,
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

				_, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo1, valVoteInfo2, valVoteInfo3},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				return s.createRequestProcessProposal(
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				cpStrategy.On("FromID", mock.Anything, uint64(1)).Return(ethUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, ethUSD, mock.Anything).Return(big.NewInt(20), nil)

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), nil)

				cpStrategy.On("FromID", mock.Anything, uint64(1)).Return(ethUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, ethUSD, mock.Anything).Return(big.NewInt(20), nil)

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
				valVoteInfo1, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				valVoteInfo1.VoteExtension = []byte("bad vote extension")

				_, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo1},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				return s.createRequestProcessProposal(
					proposal,
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

				_, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				return s.createRequestProcessProposal(
					proposal,
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
			name: "rejects a request with ve that contains invalid currency pair id",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				_, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				return s.createRequestProcessProposal(
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, fmt.Errorf("no rizz error ha"))

				return cpStrategy
			},
			expectedError: true,
			expectedResp: &cometabci.ResponseProcessProposal{
				Status: cometabci.ResponseProcessProposal_REJECT,
			},
		},
		{
			name: "rejects a request with ve that contains invalid price bytes",
			request: func() *cometabci.RequestProcessProposal {
				valVoteInfo, err := testutils.CreateExtendedVoteInfo(val1, prices1, s.codec)
				s.Require().NoError(err)

				_, commitInfoBz, err := testutils.CreateExtendedCommitInfo(
					[]cometabci.ExtendedVoteInfo{valVoteInfo},
					s.extCommitCodec,
				)
				s.Require().NoError(err)

				proposal := [][]byte{
					commitInfoBz,
					[]byte("tx1"),
				}

				return s.createRequestProcessProposal(
					proposal,
					3,
				)
			},
			veEnabled: true,
			currencyPairStrategy: func() currencypair.CurrencyPairStrategy {
				cpStrategy := currencypairmocks.NewCurrencyPairStrategy(s.T())

				cpStrategy.On("FromID", mock.Anything, uint64(0)).Return(btcUSD, nil)
				cpStrategy.On("GetDecodedPrice", mock.Anything, btcUSD, mock.Anything).Return(big.NewInt(10), fmt.Errorf(">:("))

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
			response, err := s.proposalHandler.ProcessProposalHandler()(s.ctx, req)

			s.Require().Equal(tc.expectedResp, response)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *ProposalsTestSuite) TestProposalLatency() {
	// check that no latency is reported for a failed PrepareProposal
	metricsmocks := servicemetricsmocks.NewMetrics(s.T())

	s.Run("failed prepare proposal does not report latency", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				return nil, fmt.Errorf("prepare proposal failed")
			},
			nil,
			ve.NoOpValidateVoteExtensions,
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsmocks,
		)

		// make vote-extensions not enabled
		cp := s.ctx.ConsensusParams()
		cp.Abci.VoteExtensionsEnableHeight = 3
		s.ctx = s.ctx.WithConsensusParams(cp)

		req := &cometabci.RequestPrepareProposal{
			Height: 2,
		}

		// expect no metric to be reported
		_, err := propHandler.PrepareProposalHandler()(s.ctx, req)
		s.Require().Error(err, "prepare proposal failed")
	})

	// check that latency reported in upstream logic is ignored
	s.Run("wrapped prepare proposal latency is ignored", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			func(ctx sdk.Context, rpp *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
				// simulate a long-running prepare proposal
				time.Sleep(200 * time.Millisecond)
				return &cometabci.ResponsePrepareProposal{
					Txs: nil,
				}, nil
			},
			nil,
			ve.NoOpValidateVoteExtensions,
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsmocks,
		)

		req := s.createRequestPrepareProposal( // the votes here are invalid, but that's fine
			cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{
					{
						VoteExtension: []byte("vote-extension"),
					},
				},
			},
			nil,
			4, // vote extensions will be enabled
		)
		metricsmocks.On("ObserveABCIMethodLatency", servicemetrics.PrepareProposal, mock.Anything).Return().Run(func(args mock.Arguments) {
			// the second arg shld be a duration
			latency := args.Get(1).(time.Duration)
			s.Require().True(latency < 100*time.Millisecond) // shld have ignored the wrapped prepareProposal latency
		}).Once()

		_, err := propHandler.PrepareProposalHandler()(s.ctx, req)
		s.Require().NoError(err)
	})

	s.Run("failed process proposal does not report latency", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			baseapp.NoOpPrepareProposal(),
			func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				return nil, fmt.Errorf("process proposal failed")
			},
			ve.NoOpValidateVoteExtensions,
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsmocks,
		)

		req := s.createRequestProcessProposal(nil, 2)

		_, err := propHandler.ProcessProposalHandler()(s.ctx, req)
		s.Require().Error(err, "process proposal failed")
	})

	s.Run("wrapped process proposal latency is ignored", func() {
		propHandler := proposals.NewProposalHandler(
			log.NewTestLogger(s.T()),
			baseapp.NoOpPrepareProposal(),
			func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
				// simulate a long-running process proposal
				time.Sleep(200 * time.Millisecond)
				return &cometabci.ResponseProcessProposal{}, nil
			},
			ve.NoOpValidateVoteExtensions,
			codec.NewDefaultVoteExtensionCodec(),
			codec.NewDefaultExtendedCommitCodec(),
			currencypairmocks.NewCurrencyPairStrategy(s.T()),
			metricsmocks,
		)

		_, extInfoBz, err := testutils.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{},
			codec.NewDefaultExtendedCommitCodec(),
		)
		s.Require().NoError(err)

		req := s.createRequestProcessProposal([][]byte{extInfoBz}, 4)
		metricsmocks.On("ObserveABCIMethodLatency", servicemetrics.ProcessProposal, mock.Anything).Return().Run(func(args mock.Arguments) {
			// the second arg shld be a duration
			latency := args.Get(1).(time.Duration)
			s.Require().True(latency < 100*time.Millisecond) // shld have ignored the wrapped processProposal latency
		}).Once()

		_, err = propHandler.ProcessProposalHandler()(s.ctx, req)
		s.Require().NoError(err)
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
		// Use the same default MaxTxBytes that comet does
		MaxTxBytes: types.DefaultBlockParams().MaxBytes,
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
