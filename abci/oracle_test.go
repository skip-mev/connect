package abci_test

import (
	"time"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/holiman/uint256"
	abcitypes "github.com/skip-mev/slinky/abci/types"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
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

			voteExtension, err := suite.proposalHandler.GetOracleDataFromVE(voteExtensionBz)
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
					"BTC/ETH/18": "0x1",
				}
				timestamp := time.Now()
				height := int64(100)
				valAddress := suite.createValAddress("a")

				commitInfo := suite.createExtendedVoteInfo(valAddress, prices, timestamp, height)

				return []cometabci.ExtendedVoteInfo{commitInfo}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH", 18): uint256.NewInt(1),
			},
		},
		{
			name: "single valid commit info with multiple prices",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices := map[string]string{
					"BTC/ETH/18": "0x1",
					"ETH/USD/6":  "0x2",
				}
				timestamp := time.Now()
				height := int64(100)
				valAddress := suite.createValAddress("a")

				commitInfo := suite.createExtendedVoteInfo(valAddress, prices, timestamp, height)

				return []cometabci.ExtendedVoteInfo{commitInfo}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH", 18): uint256.NewInt(1),
				oracletypes.NewCurrencyPair("ETH", "USD", 6):  uint256.NewInt(2),
			},
		},
		{
			name: "multiple valid commit infos",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH/18": "0x1",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"ETH/USD/6": "0x2",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH", 18): uint256.NewInt(1),
				oracletypes.NewCurrencyPair("ETH", "USD", 6):  uint256.NewInt(2),
			},
		},
		{
			name: "multiple valid commit infos for same asset",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH/18": "0x1",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"BTC/ETH/18": "0x2",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				prices3 := map[string]string{
					"BTC/ETH/18": "0x3",
				}
				timestamp3 := time.Now()
				height3 := int64(100)
				valAddress3 := suite.createValAddress("c")

				commitInfo3 := suite.createExtendedVoteInfo(valAddress3, prices3, timestamp3, height3)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2, commitInfo3}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH", 18): uint256.NewInt(2),
			},
		},
		{
			name: "multiple valid commit infos for same asset with different decimals",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH/18": "0x1",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"BTC/ETH/6": "0x2",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				prices3 := map[string]string{
					"BTC/ETH/12": "0x3",
				}
				timestamp3 := time.Now()
				height3 := int64(100)
				valAddress3 := suite.createValAddress("c")

				commitInfo3 := suite.createExtendedVoteInfo(valAddress3, prices3, timestamp3, height3)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2, commitInfo3}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH", 18): uint256.NewInt(1),
				oracletypes.NewCurrencyPair("BTC", "ETH", 6):  uint256.NewInt(2),
				oracletypes.NewCurrencyPair("BTC", "ETH", 12): uint256.NewInt(3),
			},
		},
		{
			name: "multiple commit infos with an average",
			getCommitInfos: func() []cometabci.ExtendedVoteInfo {
				prices1 := map[string]string{
					"BTC/ETH/18": "0x2",
				}
				timestamp1 := time.Now()
				height1 := int64(100)
				valAddress1 := suite.createValAddress("a")

				commitInfo1 := suite.createExtendedVoteInfo(valAddress1, prices1, timestamp1, height1)

				prices2 := map[string]string{
					"BTC/ETH/18": "0x4",
				}
				timestamp2 := time.Now()
				height2 := int64(100)
				valAddress2 := suite.createValAddress("b")

				commitInfo2 := suite.createExtendedVoteInfo(valAddress2, prices2, timestamp2, height2)

				return []cometabci.ExtendedVoteInfo{commitInfo1, commitInfo2}
			},
			expectedPrices: map[oracletypes.CurrencyPair]*uint256.Int{
				oracletypes.NewCurrencyPair("BTC", "ETH", 18): uint256.NewInt(3),
			},
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			commitInfos := suite.createExtendedCommitInfo(tc.getCommitInfos())
			oracleInfoBz, err := suite.proposalHandler.AggregateOracleData(suite.ctx, commitInfos)
			suite.Require().NoError(err)

			oracleInfo := &abcitypes.OracleData{}
			suite.Require().NoError(oracleInfo.Unmarshal(oracleInfoBz))

			suite.Require().Equal(len(tc.expectedPrices), len(oracleInfo.Prices))

			for currencyPairStr, priceStr := range oracleInfo.Prices {
				currencyPair, err := oracletypes.NewCurrencyPairFromString(currencyPairStr)
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
