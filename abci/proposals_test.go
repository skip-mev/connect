package abci_test

import (
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/skip-mev/slinky/abci"
	"github.com/skip-mev/slinky/abci/types"
)

func (suite *ABCITestSuite) TestPrepareProposal() {
	cases := []struct {
		name string
		// returns the request and expected response for the prepare proposal
		getReq         func() *cometabci.RequestPrepareProposal
		validators     []validator
		totalBonded    math.Int
		expectedPrices map[string]string
		expectedError  bool
	}{
		{
			name: "no txs single vote extension with no price updates",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with no price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				extendVoteInfo := suite.createExtendedVoteInfo(
					validator1,
					nil,
					timestamps,
					heights,
				)

				// create a request that includes the previous commit info (vote extensions)
				req := suite.createRequestPrepareProposal(
					suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{extendVoteInfo}),
					nil,
				)

				return req
			},
			validators:     []validator{},
			totalBonded:    math.NewInt(100),
			expectedPrices: map[string]string{},
			expectedError:  false,
		},
		{
			name: "no txs single vote extension with price updates for single asset",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH/18": "0x5",
				}
				extendVoteInfo := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				req := suite.createRequestPrepareProposal(
					suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{extendVoteInfo}),
					nil,
				)

				return req
			},
			validators: []validator{
				{
					stake:   math.NewInt(100),
					address: validator1,
				},
			},
			totalBonded: math.NewInt(100),
			expectedPrices: map[string]string{
				"BTC/ETH/18": "0x5",
			},
			expectedError: false,
		},
		{
			name: "multiple validators with price updates for single asset",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH/18": "0x5",
				}
				extendVoteInfoA := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				// Create a vote extension with price updates for validator b
				prices = map[string]string{
					"BTC/ETH/18": "0x1",
				}
				extendVoteInfoB := suite.createExtendedVoteInfo(
					validator2,
					prices,
					timestamps,
					heights,
				)

				req := suite.createRequestPrepareProposal(
					suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{extendVoteInfoA, extendVoteInfoB}),
					nil,
				)

				return req
			},
			validators: []validator{
				{
					stake:   math.NewInt(80),
					address: validator1,
				},
				{
					stake:   math.NewInt(20),
					address: validator2,
				},
			},
			totalBonded: math.NewInt(100),
			expectedPrices: map[string]string{
				"BTC/ETH/18": "0x5",
			},
			expectedError: false,
		},
		{
			name: "multiple validators with price updates for multiple assets",
			getReq: func() *cometabci.RequestPrepareProposal {
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH/18": "0x5",
					"BTC/USD/8":  "0x1",
				}
				extendVoteInfoA := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				// Create a vote extension with price updates for validator b
				prices = map[string]string{
					"BTC/ETH/18": "0x1",
					"BTC/USD/8":  "0x3",
				}
				extendVoteInfoB := suite.createExtendedVoteInfo(
					validator2,
					prices,
					timestamps,
					heights,
				)

				req := suite.createRequestPrepareProposal(
					suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{extendVoteInfoA, extendVoteInfoB}),
					nil,
				)

				return req
			},
			validators: []validator{
				{
					stake:   math.NewInt(80),
					address: validator1,
				},
				{
					stake:   math.NewInt(20),
					address: validator2,
				},
			},
			totalBonded: math.NewInt(100),
			expectedPrices: map[string]string{
				"BTC/ETH/18": "0x5",
				"BTC/USD/8":  "0x1",
			},
			expectedError: false,
		},
		{
			name: "multiple validators with price updates for multiple assets and txs",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH/18": "0x5",
					"BTC/USD/8":  "0x1",
				}
				extendVoteInfoA := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				// Create a vote extension with price updates for validator b
				prices = map[string]string{
					"BTC/ETH/18": "0x1",
					"BTC/USD/8":  "0x3",
				}
				extendVoteInfoB := suite.createExtendedVoteInfo(
					validator2,
					prices,
					timestamps,
					heights,
				)

				prices = map[string]string{
					"BTC/ETH/18": "0x3",
					"BTC/USD/8":  "0x2",
				}
				extendVoteInfoC := suite.createExtendedVoteInfo(
					validator3,
					prices,
					timestamps,
					heights,
				)

				// Create a transaction that will be included in the proposal.
				tx := []byte("tx")

				req := suite.createRequestPrepareProposal(
					suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{
						extendVoteInfoA,
						extendVoteInfoB,
						extendVoteInfoC,
					}),
					[][]byte{tx},
				)

				return req
			},
			validators: []validator{
				{
					stake:   math.NewInt(30),
					address: validator1,
				},
				{
					stake:   math.NewInt(30),
					address: validator2,
				},
				{
					stake:   math.NewInt(40),
					address: validator3,
				},
			},
			totalBonded: math.NewInt(100),
			expectedPrices: map[string]string{
				"BTC/ETH/18": "0x3",
				"BTC/USD/8":  "0x2",
			},
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// Create a mock validator store.
			validatorStore := suite.createMockValidatorStore(
				tc.validators,
				tc.totalBonded,
			)

			// Create a stake weighted median aggregator.
			aggregateFn := abci.StakeWeightedMedian(suite.ctx, validatorStore, abci.DefaultPowerThreshold)

			// Create a proposal handler.
			suite.proposalHandler = abci.NewProposalHandler(
				log.NewNopLogger(),
				suite.prepareProposalHandler,
				suite.processProposalHandler,
				aggregateFn,
			)
			prepareProposalHandler := suite.proposalHandler.PrepareProposalHandler()

			// Create a proposal.
			resp, err := prepareProposalHandler(suite.ctx, tc.getReq())
			if tc.expectedError {
				suite.Require().Error(err)
				suite.Require().Nil(resp)
			} else {
				suite.Require().NoError(err)

				// first index must be the oracle data
				suite.Require().GreaterOrEqual(len(resp.Txs), 1)

				// There cannot be any errors when unmarshalling the oracle data.
				oracleData := types.OracleData{}
				err := oracleData.Unmarshal(resp.Txs[0])
				suite.Require().NoError(err)

				// Check that the oracle data contains the expected prices.
				suite.Require().Len(oracleData.Prices, len(tc.expectedPrices))
				for currencyPair, price := range tc.expectedPrices {
					suite.Require().Contains(oracleData.Prices, currencyPair)
					suite.Require().Equal(price, oracleData.Prices[currencyPair])
				}
			}
		})
	}
}

func (suite *ABCITestSuite) TestProcessProposal() {
	suite.T().Skip("TODO")
}
