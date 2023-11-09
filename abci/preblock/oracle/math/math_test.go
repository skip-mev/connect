package math_test

import (
	"crypto"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"

	preblockmath "github.com/skip-mev/slinky/abci/preblock/oracle/math"
	"github.com/skip-mev/slinky/abci/preblock/oracle/math/mocks"
	"github.com/skip-mev/slinky/abci/testutils"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/suite"
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

func (s *MathTestSuite) TestVoteWeightedMedian() {
	cases := []struct {
		name              string
		providerPrices    aggregator.AggregatedProviderPrices
		validators        []validator
		totalBondedTokens sdkmath.Int
		expectedPrices    map[types.CurrencyPair]*uint256.Int
	}{
		{
			name:              "no providers",
			providerPrices:    aggregator.AggregatedProviderPrices{},
			validators:        []validator{},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[types.CurrencyPair]*uint256.Int{},
		},
		{
			name: "single provider entire stake + single price",
			providerPrices: aggregator.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(100),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*uint256.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: uint256.NewInt(100),
			},
		},
		{
			name: "single provider with not enough stake + single price",
			providerPrices: aggregator.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(50),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[types.CurrencyPair]*uint256.Int{},
		},
		{
			name: "single provider with just enough stake + multiple prices",
			providerPrices: aggregator.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
					{
						Base:  "ETH",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(200),
					},
				},
			},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(68),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[types.CurrencyPair]*uint256.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: uint256.NewInt(100),
				{
					Base:  "ETH",
					Quote: "USD",
				}: uint256.NewInt(200),
			},
		},
		{
			name: "2 providers with equal stake + single asset",
			providerPrices: aggregator.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
				},
				validator2.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(200),
					},
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
			expectedPrices: map[types.CurrencyPair]*uint256.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: uint256.NewInt(100),
			},
		},
		{
			name: "3 providers with equal stake + single asset",
			providerPrices: aggregator.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
				},
				validator2.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(200),
					},
				},
				validator3.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(300),
					},
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
			expectedPrices: map[types.CurrencyPair]*uint256.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: uint256.NewInt(200),
			},
		},
		{
			name: "3 providers with equal stake + multiple assets",
			providerPrices: aggregator.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
					{
						Base:  "ETH",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(200),
					},
				},
				validator2.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(300),
					},
					{
						Base:  "ETH",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(400),
					},
				},
				validator3.String(): map[types.CurrencyPair]aggregator.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(500),
					},
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
			expectedPrices: map[types.CurrencyPair]*uint256.Int{ // only btc/usd should be included
				{
					Base:  "BTC",
					Quote: "USD",
				}: uint256.NewInt(300),
			},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			// Create a mock validator store.
			mockValidatorStore := s.createMockValidatorStore(tc.validators, tc.totalBondedTokens)

			// Compute the stake weighted median.
			aggregateFn := preblockmath.VoteWeightedMedian(s.ctx, log.NewTestLogger(s.T()), mockValidatorStore, preblockmath.DefaultPowerThreshold)
			result := aggregateFn(tc.providerPrices)

			// Verify the result.
			s.Require().Len(result, len(tc.expectedPrices))
			for currencyPair, expectedPrice := range tc.expectedPrices {
				s.Require().Equal(expectedPrice, result[currencyPair])
			}
		})
	}
}

func (s *MathTestSuite) TestComputeVoteWeightedMedian() {
	cases := []struct {
		name      string
		priceInfo preblockmath.VoteWeightedPriceInfo
		expected  *uint256.Int
	}{
		{
			name: "single price",
			priceInfo: preblockmath.VoteWeightedPriceInfo{
				Prices: []preblockmath.VoteWeightedPricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      uint256.NewInt(100),
					},
				},
				TotalWeight: sdkmath.NewInt(1),
			},
			expected: uint256.NewInt(100),
		},
		{
			name: "two prices that are equal",
			priceInfo: preblockmath.VoteWeightedPriceInfo{
				Prices: []preblockmath.VoteWeightedPricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      uint256.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      uint256.NewInt(100),
					},
				},
				TotalWeight: sdkmath.NewInt(2),
			},
			expected: uint256.NewInt(100),
		},
		{
			name: "two prices that are not equal",
			priceInfo: preblockmath.VoteWeightedPriceInfo{
				Prices: []preblockmath.VoteWeightedPricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      uint256.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(1),
						Price:      uint256.NewInt(200),
					},
				},
				TotalWeight: sdkmath.NewInt(2),
			},
			expected: uint256.NewInt(100),
		},
		{
			name: "two prices that are not equal with different weights",
			priceInfo: preblockmath.VoteWeightedPriceInfo{
				Prices: []preblockmath.VoteWeightedPricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(10),
						Price:      uint256.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(20),
						Price:      uint256.NewInt(200),
					},
				},
				TotalWeight: sdkmath.NewInt(30),
			},
			expected: uint256.NewInt(200),
		},
		{
			name: "three prices that are not equal with different weights",
			priceInfo: preblockmath.VoteWeightedPriceInfo{
				Prices: []preblockmath.VoteWeightedPricePerValidator{
					{
						VoteWeight: sdkmath.NewInt(10),
						Price:      uint256.NewInt(100),
					},
					{
						VoteWeight: sdkmath.NewInt(20),
						Price:      uint256.NewInt(200),
					},
					{
						VoteWeight: sdkmath.NewInt(30),
						Price:      uint256.NewInt(300),
					},
				},
				TotalWeight: sdkmath.NewInt(60),
			},
			expected: uint256.NewInt(200),
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			result := preblockmath.ComputeVoteWeightedMedian(tc.priceInfo)
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
