package abci_test

import (
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/abci"
	abcitypes "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func (suite *ABCITestSuite) TestGetOracleDataFromVE() {
	cases := []struct {
		name             string
		getVoteExtension func() []byte
		expectedError    bool
	}{
		{
			name: "nil vote extension",
			getVoteExtension: func() []byte {
				return nil
			},
			expectedError: true,
		},
		{
			name: "empty vote extension",
			getVoteExtension: func() []byte {
				return []byte{}
			},
			expectedError: true,
		},
		{
			name: "valid vote extension",
			getVoteExtension: func() []byte {
				prices := map[string]string{
					"BTC": "100",
					"ETH": "200",
				}
				timestamp := time.Now()
				height := int64(100)

				voteExtension := suite.createVoteExtension(prices, timestamp, height)
				voteExtensionBz, err := voteExtension.Marshal()
				suite.Require().NoError(err)

				return voteExtensionBz
			},
			expectedError: false,
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			voteExtensionBz := tc.getVoteExtension()

			oracle := abci.NewOracle(
				log.NewTestLogger(suite.T()),
				types.ComputeMedian(),
				suite.oracleKeeper,
				suite.NoOpValidateVEFn(),
				suite.validatorStore,
			)

			voteExtension, err := oracle.GetOracleDataFromVE(voteExtensionBz)
			if tc.expectedError {
				suite.Require().Error(err)
				suite.Require().Nil(voteExtension)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(voteExtension)

				oracleData := &abcitypes.OracleVoteExtension{}
				suite.Require().NoError(oracleData.Unmarshal(voteExtensionBz))
				suite.Require().Equal(oracleData, voteExtension)
			}
		})
	}
}

func (suite *ABCITestSuite) TestAggregateOracleData() {
	cases := []struct {
		name           string
		getCommitInfos func() []cometabci.ExtendedVoteInfo
		expectedPrices map[oracletypes.CurrencyPair]*uint256.Int
	}{
		{
			name: "empty commit infos",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				return nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{},
		},
		{
			name: "empty commit infos",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				return []cometabci.ExtendedVoteInfo{}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{},
		},
		{
			name: "single valid commit infos",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices := map[string]string{
					"BTC/ETH": "0x1",
				}
				timestamp := time.Now()
				height := int64(100)
				valAddress := suite.createValAddress("a")

				commitInfo := suite.createExtendedVoteInfo(valAddress, prices, timestamp, height)

				return []cometabci.ExtendedVoteInfo{commitInfo}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH"): uint256.NewInt(1),
			},
		},
		{
			name: "single valid commit info with multiple prices",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices := map[string]string{
					"BTC/ETH": "0x1",
					"ETH/USD": "0x2",
				}
				timestamp := time.Now()
				height := int64(100)
				valAddress := suite.createValAddress("a")

				commitInfo := suite.createExtendedVoteInfo(valAddress, prices, timestamp, height)

				return []cometabci.ExtendedVoteInfo{commitInfo}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH"): uint256.NewInt(1),
				oracletypes.NewCurrencyPair("ETH", "USD"): uint256.NewInt(2),
			},
		},
		{
			name: "multiple valid commit infos",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH": "0x1",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"ETH/USD": "0x2",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH"): uint256.NewInt(1),
				oracletypes.NewCurrencyPair("ETH", "USD"): uint256.NewInt(2),
			},
		},
		{
			name: "multiple valid commit infos for same asset",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH": "0x1",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"BTC/ETH": "0x2",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				prices3 := map[string]string{
					"BTC/ETH": "0x3",
				}
				timestamp3 := time.Now()
				height3 := int64(100)
				valAddress3 := suite.createValAddress("c")

				commitInfo3 := suite.createExtendedVoteInfo(valAddress3, prices3, timestamp3, height3)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2, commitInfo3}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH"): uint256.NewInt(2),
			},
		},
		{
			name: "multiple valid commit infos for same asset",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH": "0x1",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"BTC/ETH": "0x2",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				prices3 := map[string]string{
					"BTC/ETH": "0x3",
				}
				timestamp3 := time.Now()
				height3 := int64(100)
				valAddress3 := suite.createValAddress("c")

				commitInfo3 := suite.createExtendedVoteInfo(valAddress3, prices3, timestamp3, height3)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2, commitInfo3}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH"): uint256.NewInt(2),
			},
		},
		{
			name: "multiple commit infos with an average",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH": "0x2",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"BTC/ETH": "0x4",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH"): uint256.NewInt(3),
			},
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			oracle := abci.NewOracle(
				log.NewTestLogger(suite.T()),
				types.ComputeMedian(),
				suite.oracleKeeper,
				suite.NoOpValidateVEFn(),
				suite.validatorStore,
			)

			commitInfos := suite.createExtendedCommitInfo(tc.getCommitInfos())
			oracleData, err := oracle.AggregateOracleData(suite.ctx, commitInfos)
			suite.Require().NoError(err)

			suite.Require().Equal(len(tc.expectedPrices), len(oracleData.Prices))

			for currencyPairStr, priceStr := range oracleData.Prices {
				currencyPair, err := oracletypes.CurrencyPairFromString(currencyPairStr)
				suite.Require().NoError(err)

				expectedPrice, ok := tc.expectedPrices[currencyPair]
				suite.Require().True(ok)

				price, err := uint256.FromHex(priceStr)
				suite.Require().NoError(err)
				suite.Require().Equal(expectedPrice, price)
			}
		})
	}
}

func (suite *ABCITestSuite) TestVerifyOraclePrices() {
	cases := []struct {
		name          string
		getOracleInfo func() abcitypes.OracleData
		expectedError bool
	}{
		{
			name: "empty oracle info",
			getOracleInfo: func() abcitypes.OracleData {
				return abcitypes.OracleData{}
			},
			expectedError: false,
		},
		{
			name: "valid oracle info with a single price",
			getOracleInfo: func() abcitypes.OracleData {
				voteExtension1 := suite.createExtendedVoteInfo(
					validator1,
					map[string]string{
						"BTC/ETH": "0x1",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH": "0x1",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1},
				)

				return oracleData
			},
			expectedError: false,
		},
		{
			name: "valid oracle info with multiple prices from single validator",
			getOracleInfo: func() abcitypes.OracleData {
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

				return oracleData
			},
			expectedError: false,
		},
		{
			name: "valid oracle info with multiple prices from multiple validators",
			getOracleInfo: func() abcitypes.OracleData {
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

				return oracleData
			},
			expectedError: false,
		},
		{
			name: "vote extensions with multiple prices but some are missing in oracle data",
			getOracleInfo: func() abcitypes.OracleData {
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
						"ETH/USD":   "0x2",
						"ATOM/USDC": "0x3",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{voteExtension1, voteExtension2},
				)

				return oracleData
			},
			expectedError: true,
		},
		{
			name: "multiple vote extensions for the same asset but with different prices",
			getOracleInfo: func() abcitypes.OracleData {
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
					validator2,
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

				return oracleData
			},
			expectedError: false,
		},
		{
			name: "oracle data reported the wrong prices",
			getOracleInfo: func() abcitypes.OracleData {
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
					validator2,
					map[string]string{
						"BTC/ETH": "0x3",
					},
					time.Now(),
					100,
				)

				oracleData := suite.createOracleData(
					map[string]string{
						"BTC/ETH": "0x3",
					},
					time.Now(),
					100,
					[]cometabci.ExtendedVoteInfo{
						voteExtension1,
						voteExtension2,
						voteExtension3,
					},
				)

				return oracleData
			},
			expectedError: true,
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			oracle := abci.NewOracle(
				log.NewTestLogger(suite.T()),
				types.ComputeMedian(),
				suite.oracleKeeper,
				suite.NoOpValidateVEFn(),
				suite.validatorStore,
			)

			oracleInfo := tc.getOracleInfo()

			extendedCommitInfo := cometabci.ExtendedCommitInfo{}
			suite.Require().NoError(extendedCommitInfo.Unmarshal(oracleInfo.ExtendedCommitInfo))

			_, err := oracle.VerifyOraclePrices(suite.ctx, oracleInfo, extendedCommitInfo)
			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *ABCITestSuite) TestWriteOracleData() {
	cases := []struct {
		name       string
		oracleData abcitypes.OracleData
	}{
		{
			name: "empty oracle data",
			oracleData: abcitypes.OracleData{
				Prices: map[string]string{},
			},
		},
		{
			name: "single valid oracle data",
			oracleData: abcitypes.OracleData{
				Prices: map[string]string{
					"BTC/ETH": "0x1",
				},
			},
		},
		{
			name: "multiple valid oracle data",
			oracleData: abcitypes.OracleData{
				Prices: map[string]string{
					"BTC/ETH": "0x1",
					"ETH/USD": "0x2",
				},
			},
		},
		{
			name: "posting prices that are not supported by the oracle module",
			oracleData: abcitypes.OracleData{
				Prices: map[string]string{
					"BTC/ETH":   "0x1",
					"ETH/USD":   "0x2",
					"ATOM/USDC": "0x3",
				},
			},
		},
		{
			name: "posting prices that are not supported by the oracle module",
			oracleData: abcitypes.OracleData{
				Prices: map[string]string{
					"BTC/ETH": "1",
				},
			},
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.oracleKeeper.InitGenesis(suite.ctx, suite.genesis)

			oracle := abci.NewOracle(
				log.NewTestLogger(suite.T()),
				types.ComputeMedian(),
				suite.oracleKeeper,
				suite.NoOpValidateVEFn(),
				suite.validatorStore,
			)

			oracle.WriteOracleData(suite.ctx, tc.oracleData)

			// ensure that no new currency pairs were added to the module
			currencyPairs := suite.oracleKeeper.GetAllCurrencyPairs(suite.ctx)
			keeperCPs := make(map[string]struct{})
			oracleCPs := make(map[string]struct{})

			for _, cp := range suite.currencyPairs {
				keeperCPs[cp.ToString()] = struct{}{}
			}

			for _, cp := range currencyPairs {
				oracleCPs[cp.ToString()] = struct{}{}
			}
			suite.Require().Equal(keeperCPs, oracleCPs)

			// ensure that the prices were written to the store
			for _, currencyPair := range currencyPairs {
				// If the currency pair is not in the oracle data, then skip it.
				priceHex, ok := tc.oracleData.Prices[currencyPair.ToString()]
				if !ok {
					continue
				}

				price, err := uint256.FromHex(priceHex)
				if err != nil {
					continue
				}

				sdkInt := math.NewIntFromBigInt(price.ToBig())

				priceInfo, err := suite.oracleKeeper.GetPriceWithNonceForCurrencyPair(suite.ctx, currencyPair)
				suite.Require().NoError(err)

				suite.Require().Equal(sdkInt, priceInfo.Price)
				suite.Require().Equal(priceInfo.BlockHeight, uint64(suite.ctx.BlockHeight()))
				suite.Require().Equal(uint64(1), priceInfo.Nonce())
				suite.Require().Equal(priceInfo.BlockTimestamp, suite.ctx.BlockHeader().Time)
			}
		})
	}
}
