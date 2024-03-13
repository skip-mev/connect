package voteweighted_test

import (
	"math/big"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/pkg/math/voteweighted"
	"github.com/skip-mev/slinky/pkg/math/voteweighted/mocks"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

var btcUsd = slinkytypes.CurrencyPair{
	Base:  "BTC",
	Quote: "USD",
}

func (s *MathTestSuite) TestConstrainedSWMedian() {
	cases := []struct {
		name              string
		existingPrices    map[slinkytypes.CurrencyPair]*big.Int
		providerPrices    aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]
		validators        []validator
		totalBondedTokens sdkmath.Int
		expectedPrices    map[slinkytypes.CurrencyPair]*big.Int
	}{
		{
			name:              "no providers",
			existingPrices:    make(map[slinkytypes.CurrencyPair]*big.Int),
			providerPrices:    aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{},
			validators:        []validator{},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[slinkytypes.CurrencyPair]*big.Int{},
		},
		{
			name: "single provider entire stake + single price, no price change",
			existingPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(100),
			},
			providerPrices: aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				validator1.String(): map[slinkytypes.CurrencyPair]*big.Int{
					btcUsd: big.NewInt(100),
				},
			},
			validators: []validator{
				{stake: sdkmath.NewInt(100), consAddr: validator1},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(100),
			},
		},
		{
			name: "single provider with less than 2/3 stake but small price movement posts update",
			existingPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(100),
			},
			providerPrices: aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				validator1.String(): map[slinkytypes.CurrencyPair]*big.Int{
					btcUsd: big.NewInt(101),
				},
			},
			validators: []validator{
				{stake: sdkmath.NewInt(50), consAddr: validator1},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(101),
			},
		},
		{
			name: "3 providers with equal stake, but large movement posts no update",
			existingPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(100),
			},
			providerPrices: aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				validator1.String(): map[slinkytypes.CurrencyPair]*big.Int{
					btcUsd: big.NewInt(100),
				},
				validator2.String(): map[slinkytypes.CurrencyPair]*big.Int{
					btcUsd: big.NewInt(200),
				},
				validator3.String(): map[slinkytypes.CurrencyPair]*big.Int{
					btcUsd: big.NewInt(300),
				},
			},
			validators: []validator{
				{stake: sdkmath.NewInt(33), consAddr: validator1},
				{stake: sdkmath.NewInt(33), consAddr: validator2},
				{stake: sdkmath.NewInt(33), consAddr: validator3},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[slinkytypes.CurrencyPair]*big.Int{},
		},
		{
			name: "3 providers with equal stake, mostly low movement posts update",
			existingPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(100),
			},
			providerPrices: aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				validator1.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(105)},
				validator2.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(110)},
				validator3.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(300)},
			},
			validators: []validator{
				{stake: sdkmath.NewInt(33), consAddr: validator1},
				{stake: sdkmath.NewInt(33), consAddr: validator2},
				{stake: sdkmath.NewInt(33), consAddr: validator3},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(110),
			},
		},
		{
			name: "4 providers with equal stake, mostly low movement, different directions posts update",
			existingPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(100),
			},
			providerPrices: aggregator.AggregatedProviderData[string, map[slinkytypes.CurrencyPair]*big.Int]{
				validator1.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(101)},
				validator2.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(99)},
				validator3.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(101)},
				validator4.String(): map[slinkytypes.CurrencyPair]*big.Int{btcUsd: big.NewInt(99)},
			},
			validators: []validator{
				{stake: sdkmath.NewInt(25), consAddr: validator1},
				{stake: sdkmath.NewInt(25), consAddr: validator2},
				{stake: sdkmath.NewInt(25), consAddr: validator3},
				{stake: sdkmath.NewInt(25), consAddr: validator4},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[slinkytypes.CurrencyPair]*big.Int{
				btcUsd: big.NewInt(99),
			},
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			mockValidatorStore := s.createMockValidatorStore(tc.validators, tc.totalBondedTokens)
			logger := log.NewNopLogger()
			mockOracleKeeper := s.createMockOracleKeeper(tc.existingPrices)
			result := voteweighted.ConstrainedSWMedian(
				logger,
				mockValidatorStore,
				voteweighted.DefaultPowerThreshold,
				mockOracleKeeper,
				s.defaultThresholdFunction(),
			)(s.ctx)(tc.providerPrices)
			// Verify the result.
			s.Require().Len(result, len(tc.expectedPrices))
			for currencyPair, expectedPrice := range tc.expectedPrices {
				s.Require().Equal(expectedPrice, result[currencyPair])
			}
		})
	}
}

func (s *MathTestSuite) defaultThresholdFunction() voteweighted.ThresholdDetermination {
	return func(currentPrice *big.Int, proposedPrice *big.Int, priceInfo voteweighted.PriceInfo) sdkmath.Int {
		return voteweighted.ThresholdWeightCalc(currentPrice, proposedPrice, big.NewInt(500_000), priceInfo)
	}
}

func (s *MathTestSuite) createMockOracleKeeper(existingPrices map[slinkytypes.CurrencyPair]*big.Int) *mocks.OracleKeeper {
	mockOracleKeeper := mocks.NewOracleKeeper(s.T())
	for cp, price := range existingPrices {
		mockOracleKeeper.On("GetPriceForCurrencyPair", mock.Anything, cp).
			Return(types.QuotePrice{
				Price:          sdkmath.NewIntFromBigInt(price),
				BlockTimestamp: time.Now(),
			}, nil)
	}
	return mockOracleKeeper
}
