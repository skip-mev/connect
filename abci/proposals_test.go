package abci_test

import (
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/skip-mev/slinky/abci"
	"github.com/skip-mev/slinky/abci/types"
	oracleservicetypes "github.com/skip-mev/slinky/oracle/types"
)

func (suite *ABCITestSuite) TestPrepareProposal() {
	cases := []struct {
		name string
		// returns the request and expected response for the prepare proposal
		getReq         func() *cometabci.RequestPrepareProposal
		validators     []validator
		totalBonded    math.Int
		expectedPrices map[string]string
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
		},
		{
			name: "no txs single vote extension with price updates for single asset",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH": "0x5",
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
				"BTC/ETH": "0x5",
			},
		},
		{
			name: "multiple validators with price updates for single asset",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH": "0x5",
				}
				extendVoteInfoA := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				// Create a vote extension with price updates for validator b
				prices = map[string]string{
					"BTC/ETH": "0x1",
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
				"BTC/ETH": "0x5",
			},
		},
		{
			name: "multiple validators with price updates for multiple assets",
			getReq: func() *cometabci.RequestPrepareProposal {
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH": "0x5",
					"BTC/USD": "0x1",
				}
				extendVoteInfoA := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				// Create a vote extension with price updates for validator b
				prices = map[string]string{
					"BTC/ETH": "0x1",
					"BTC/USD": "0x3",
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
				"BTC/ETH": "0x5",
				"BTC/USD": "0x1",
			},
		},
		{
			name: "multiple validators with price updates for multiple assets and txs",
			getReq: func() *cometabci.RequestPrepareProposal {
				// Create a vote extension with price updates for validator a
				timestamps := time.Now()
				heights := suite.ctx.BlockHeight()
				prices := map[string]string{
					"BTC/ETH": "0x5",
					"BTC/USD": "0x1",
				}
				extendVoteInfoA := suite.createExtendedVoteInfo(
					validator1,
					prices,
					timestamps,
					heights,
				)

				// Create a vote extension with price updates for validator b
				prices = map[string]string{
					"BTC/ETH": "0x1",
					"BTC/USD": "0x3",
				}
				extendVoteInfoB := suite.createExtendedVoteInfo(
					validator2,
					prices,
					timestamps,
					heights,
				)

				prices = map[string]string{
					"BTC/ETH": "0x3",
					"BTC/USD": "0x2",
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
				"BTC/ETH": "0x3",
				"BTC/USD": "0x2",
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

			oracle := abci.NewOracle(
				log.NewTestLogger(suite.T()),
				aggregateFn,
				suite.oracleKeeper,
				suite.NoOpValidateVEFn(),
			)

			// Create a proposal handler.
			suite.proposalHandler = abci.NewProposalHandler(
				log.NewTestLogger(suite.T()),
				suite.prepareProposalHandler,
				suite.processProposalHandler,
				oracle,
			)
			prepareProposalHandler := suite.proposalHandler.PrepareProposalHandler()

			// Create a proposal.
			resp, err := prepareProposalHandler(suite.ctx, tc.getReq())
			suite.Require().NoError(err)

			// first index must be the oracle data
			suite.Require().GreaterOrEqual(len(resp.Txs), 1)

			// There cannot be any errors when unmarshalling the oracle data.
			if suite.ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight > suite.voteExtensionsEnabledHeight {
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
	cases := []struct {
		name          string
		getTxs        func() [][]byte
		expectedError bool
	}{
		{
			name: "no txs before vote extensions enabled",
			getTxs: func() [][]byte {
				suite.ctx = suite.ctx.WithBlockHeight(0)

				return nil
			},
			expectedError: false,
		},
		{
			name: "no txs after vote extensions enabled",
			getTxs: func() [][]byte {
				return nil
			},
			expectedError: true,
		},
		{
			name: "single tx (no oracle data)",
			getTxs: func() [][]byte {
				return [][]byte{
					[]byte("tx"),
				}
			},
			expectedError: true,
		},
		{
			name: "empty oracle data and vote extensions",
			getTxs: func() [][]byte {
				oracleData := types.OracleData{}
				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
				}
			},
			expectedError: false,
		},
		{
			name: "empty oracle data and non-empty vote extensions",
			getTxs: func() [][]byte {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{
						"BTC/ETH":   "0x1",
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1},
				)

				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
				}
			},
			expectedError: true,
		},
		{
			name: "non-empty oracle data and empty vote extensions",
			getTxs: func() [][]byte {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH":   "0x1",
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1},
				)

				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
				}
			},
			expectedError: true,
		},
		{
			name: "non-empty oracle data and non-empty vote extensions",
			getTxs: func() [][]byte {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{
						"BTC/ETH":   "0x1",
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH":   "0x1",
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1},
				)

				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
				}
			},
			expectedError: false,
		},
		{
			name: "non-empty oracle data and non-empty vote extensions (multiple validators)",
			getTxs: func() [][]byte {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{
						"BTC/ETH": "0x1",
						"ETH/USD": "0x2",
					},
					time.Now(),
					100,
				)

				voteExtension2 := suite.createExtendedVoteInfo(
					validator2,
					map[string]string{
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH":   "0x1",
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1, voteExtension2},
				)

				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
				}
			},
			expectedError: false,
		},
		{
			name: "non-empty oracle data and non-empty vote extensions (multiple validators, multiple txs)",
			getTxs: func() [][]byte {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{
						"BTC/ETH": "0x1",
						"ETH/USD": "0x2",
					},
					time.Now(),
					100,
				)

				voteExtension2 := suite.createExtendedVoteInfo(
					validator2,
					map[string]string{
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH":   "0x1",
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1, voteExtension2},
				)

				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
					[]byte("tx"),
					[]byte("tx2"),
				}
			},
			expectedError: false,
		},
		{
			name: "multiple validators posting prices for same asset",
			getTxs: func() [][]byte {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{
						"BTC/ETH": "0x1",
					},
					time.Now(),
					100,
				)

				voteExtension2 := suite.createExtendedVoteInfo(
					validator2,
					map[string]string{
						"BTC/ETH": "0x2",
					},
					time.Now(),
					100,
				)

				voteExtension3 := suite.createExtendedVoteInfo(
					validator3,
					map[string]string{
						"BTC/ETH": "0x3",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH": "0x2",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{
						voteExtension1,
						voteExtension2,
						voteExtension3,
					},
				)

				oracleDataBytes, err := oracleData.Marshal()
				suite.Require().NoError(err)

				return [][]byte{
					oracleDataBytes,
				}
			},
			expectedError: false,
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			oracle := abci.NewOracle(
				log.NewTestLogger(suite.T()),
				oracleservicetypes.ComputeMedian(),
				suite.oracleKeeper,
				suite.NoOpValidateVEFn(),
			)

			// Create a proposal handler.
			suite.proposalHandler = abci.NewProposalHandler(
				log.NewTestLogger(suite.T()),
				suite.prepareProposalHandler,
				suite.processProposalHandler,
				oracle,
			)
			processProposalHandler := suite.proposalHandler.ProcessProposalHandler()

			// Create a proposal.
			req := &cometabci.RequestProcessProposal{
				Txs: tc.getTxs(),
			}
			resp, err := processProposalHandler(suite.ctx, req)
			if tc.expectedError {
				suite.Require().Error(err)
				suite.Require().Nil(resp)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
			}
		})
	}
}
