package voteweighted_test

import (
	"crypto"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/pkg/math/voteweighted"
	"github.com/skip-mev/slinky/pkg/math/voteweighted/mocks"
	"github.com/skip-mev/slinky/x/oracle/types"
)

type MathTestSuite struct {
	suite.Suite

	ctx sdk.Context
}

type validator struct {
	stake    sdkmath.Int
	consAddr sdk.ConsAddress
}

var (
	validator1 = sdk.ConsAddress("validator1")
	validator2 = sdk.ConsAddress("validator2")
	validator3 = sdk.ConsAddress("validator3")
)

func (s *MathTestSuite) SetupTest() {
	s.ctx = testutils.CreateBaseSDKContext(s.T())
}

func TestMathTestSuite(t *testing.T) {
	suite.Run(t, new(MathTestSuite))
}

func (s *MathTestSuite) TestMedian() {
	cases := []struct {
		name              string
		providerPrices    aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]
		validators        []validator
		totalBondedTokens sdkmath.Int
		expectedPrices    map[types.CurrencyPair]*big.Int
	}{
		{
			name:              "no providers",
			providerPrices:    aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{},
			validators:        []validator{},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[types.CurrencyPair]*big.Int{},
		},
		{
			name: "single provider entire stake + single price",
			providerPrices: aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{
				validator1.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(100),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(100),
			},
		},
		{
			name: "single provider with not enough stake + single price",
			providerPrices: aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{
				validator1.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(50),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[types.CurrencyPair]*big.Int{},
		},
		{
			name: "single provider with just enough stake + multiple prices",
			providerPrices: aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{
				validator1.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
					{
						Base:  "ETH",
						Quote: "USD",
					}: big.NewInt(200),
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(68),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(100),
				{
					Base:  "ETH",
					Quote: "USD",
				}: big.NewInt(200),
			},
		},
		{
			name: "2 providers with equal stake + single asset",
			providerPrices: aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{
				validator1.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
				},
				validator2.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(200),
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(50),
					consAddr: validator1,
				},
				{
					stake:    sdkmath.NewInt(50),
					consAddr: validator2,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(100),
			},
		},
		{
			name: "3 providers with equal stake + single asset",
			providerPrices: aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{
				validator1.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
				},
				validator2.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(200),
				},
				validator3.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(300),
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(33),
					consAddr: validator1,
				},
				{
					stake:    sdkmath.NewInt(33),
					consAddr: validator2,
				},
				{
					stake:    sdkmath.NewInt(33),
					consAddr: validator3,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(200),
			},
		},
		{
			name: "3 providers with equal stake + multiple assets",
			providerPrices: aggregator.AggregatedProviderData[string, map[types.CurrencyPair]*big.Int]{
				validator1.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
					{
						Base:  "ETH",
						Quote: "USD",
					}: big.NewInt(200),
				},
				validator2.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(300),
					{
						Base:  "ETH",
						Quote: "USD",
					}: big.NewInt(400),
				},
				validator3.String(): map[types.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(500),
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(33),
					consAddr: validator1,
				},
				{
					stake:    sdkmath.NewInt(33),
					consAddr: validator2,
				},
				{
					stake:    sdkmath.NewInt(33),
					consAddr: validator3,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*big.Int{ // only btc/usd should be included
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(300),
			},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			// Create a mock validator store.
			mockValidatorStore := s.createMockValidatorStore(tc.validators, tc.totalBondedTokens)

			// Compute the stake weighted median.
			aggregateFn := voteweighted.Median(s.ctx, log.NewTestLogger(s.T()), mockValidatorStore, voteweighted.DefaultPowerThreshold)
			result := aggregateFn(tc.providerPrices)

			// Verify the result.
			s.Require().Len(result, len(tc.expectedPrices))
			for currencyPair, expectedPrice := range tc.expectedPrices {
				s.Require().Equal(expectedPrice, result[currencyPair])
			}
		})
	}
}

func (s *MathTestSuite) TestComputeMedian() {
	cases := []struct {
		name      string
		priceInfo voteweighted.PriceInfo
		expected  *big.Int
	}{
		{
			name: "single price",
			priceInfo: voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      big.NewInt(100),
					},
				},
				TotalWeight: sdkmath.NewInt(1),
			},
			expected: big.NewInt(100),
		},
		{
			name: "two prices that are equal",
			priceInfo: voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      big.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      big.NewInt(100),
					},
				},
				TotalWeight: sdkmath.NewInt(2),
			},
			expected: big.NewInt(100),
		},
		{
			name: "two prices that are not equal",
			priceInfo: voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      big.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      big.NewInt(200),
					},
				},
				TotalWeight: sdkmath.NewInt(2),
			},
			expected: big.NewInt(100),
		},
		{
			name: "two prices that are not equal with different weights",
			priceInfo: voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(10),
						Price:      big.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(20),
						Price:      big.NewInt(200),
					},
				},
				TotalWeight: sdkmath.NewInt(30),
			},
			expected: big.NewInt(200),
		},
		{
			name: "three prices that are not equal with different weights",
			priceInfo: voteweighted.PriceInfo{
				Prices: []voteweighted.PricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(10),
						Price:      big.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(20),
						Price:      big.NewInt(200),
					},
					{
						VoteWeight: sdkmath.NewInt(30),
						Price:      big.NewInt(300),
					},
				},
				TotalWeight: sdkmath.NewInt(60),
			},
			expected: big.NewInt(200),
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			result := voteweighted.ComputeMedian(tc.priceInfo)
			s.Require().Equal(tc.expected, result)
		})
	}
}

func (s *MathTestSuite) createMockValidatorStore(
	validators []validator,
	totalTokens sdkmath.Int,
) *mocks.ValidatorStore {
	store := mocks.NewValidatorStore(s.T())
	if len(validators) != 0 {
		mockVals := make([]*mocks.ValidatorI, len(validators))
		valPubKeys := make([]crypto.PublicKey, len(validators))

		for i, val := range validators {
			mockVals[i] = mocks.NewValidatorI(s.T())
			mockVals[i].On(
				"GetBondedTokens",
			).Return(
				val.stake,
			).Maybe()

			store.On(
				"ValidatorByConsAddr",
				s.ctx,
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
				s.ctx,
				val.consAddr,
			).Return(
				valPubKeys[i],
				nil,
			).Maybe()
		}
	}

	store.On(
		"TotalBondedTokens",
		s.ctx,
	).Return(
		totalTokens, nil,
	).Maybe()

	return store
}
