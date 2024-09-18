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
	ccvtypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/connect/v2/abci/testutils"
	"github.com/skip-mev/connect/v2/aggregator"
	"github.com/skip-mev/connect/v2/pkg/math/voteweighted"
	"github.com/skip-mev/connect/v2/pkg/math/voteweighted/mocks"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
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
		providerPrices    aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]
		validators        []validator
		totalBondedTokens sdkmath.Int
		expectedPrices    map[connecttypes.CurrencyPair]*big.Int
	}{
		{
			name:           "no providers",
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{},
			validators: []validator{
				{
					stake:    sdkmath.NewInt(100),
					consAddr: validator1,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "single provider entire stake + single price",
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{
				validator1.String(): map[connecttypes.CurrencyPair]*big.Int{
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
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(100),
			},
		},
		{
			name: "single provider with not enough stake + single price",
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{
				validator1.String(): map[connecttypes.CurrencyPair]*big.Int{
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
				{
					stake:    sdkmath.NewInt(50),
					consAddr: validator2,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[connecttypes.CurrencyPair]*big.Int{},
		},
		{
			name: "single provider with just enough stake + multiple prices",
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{
				validator1.String(): map[connecttypes.CurrencyPair]*big.Int{
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
				{
					stake:    sdkmath.NewInt(32),
					consAddr: validator2,
				},
			},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
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
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{
				validator1.String(): map[connecttypes.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
				},
				validator2.String(): map[connecttypes.CurrencyPair]*big.Int{
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
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(100),
			},
		},
		{
			name: "3 providers with equal stake + single asset",
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{
				validator1.String(): map[connecttypes.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
				},
				validator2.String(): map[connecttypes.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(200),
				},
				validator3.String(): map[connecttypes.CurrencyPair]*big.Int{
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
			totalBondedTokens: sdkmath.NewInt(99),
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{
				{
					Base:  "BTC",
					Quote: "USD",
				}: big.NewInt(200),
			},
		},
		{
			name: "3 providers with equal stake + multiple assets",
			providerPrices: aggregator.AggregatedProviderData[string, map[connecttypes.CurrencyPair]*big.Int]{
				validator1.String(): map[connecttypes.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(100),
					{
						Base:  "ETH",
						Quote: "USD",
					}: big.NewInt(200),
				},
				validator2.String(): map[connecttypes.CurrencyPair]*big.Int{
					{
						Base:  "BTC",
						Quote: "USD",
					}: big.NewInt(300),
					{
						Base:  "ETH",
						Quote: "USD",
					}: big.NewInt(400),
				},
				validator3.String(): map[connecttypes.CurrencyPair]*big.Int{
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
			totalBondedTokens: sdkmath.NewInt(99),
			expectedPrices: map[connecttypes.CurrencyPair]*big.Int{ // only btc/usd should be included
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
			// Also test ICS based val keeper
			ccvConsumerCompatKeeper := s.createMockCCVConsumerCompatKeeper(tc.validators)

			// Compute the stake weighted median for both staking keeper based and ICS based val stores.
			defaultAggregateFn := voteweighted.Median(s.ctx, log.NewTestLogger(s.T()), mockValidatorStore, voteweighted.DefaultPowerThreshold)
			defaultResult := defaultAggregateFn(tc.providerPrices)
			ccvAggregateFn := voteweighted.Median(s.ctx, log.NewTestLogger(s.T()), ccvConsumerCompatKeeper, voteweighted.DefaultPowerThreshold)
			ccvResult := ccvAggregateFn(tc.providerPrices)

			// Verify the results.
			s.Require().Len(defaultResult, len(tc.expectedPrices))
			s.Require().Len(ccvResult, len(tc.expectedPrices))
			for currencyPair, expectedPrice := range tc.expectedPrices {
				s.Require().Equal(expectedPrice, defaultResult[currencyPair])
				s.Require().Equal(expectedPrice, ccvResult[currencyPair])
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

func (s *MathTestSuite) createMockCCVConsumerCompatKeeper(
	validators []validator,
) voteweighted.CCVConsumerCompatKeeper {
	valStore := mocks.NewCCValidatorStore(s.T())
	ccvCompatKeeper := voteweighted.NewCCVConsumerCompatKeeper(valStore)
	mockVals := make([]ccvtypes.CrossChainValidator, len(validators))
	if len(validators) != 0 {
		for i, val := range validators {
			valPubKey := ed25519.GenPrivKey().PubKey()
			ccVal, err := ccvtypes.NewCCValidator(
				validators[i].consAddr,
				validators[i].stake.Int64(),
				valPubKey,
			)
			if err != nil {
				panic(err)
			}
			mockVals[i] = ccVal

			valStore.On(
				"GetCCValidator",
				s.ctx,
				val.consAddr.Bytes(),
			).Return(
				ccVal,
				true,
			).Maybe()
			valStore.On(
				"GetPubKeyByConsAddr",
				s.ctx,
				val.consAddr,
			).Return(
				valPubKey,
				nil,
			).Maybe()
		}
	}
	valStore.On(
		"GetAllCCValidator",
		s.ctx,
	).Return(
		mockVals,
	)

	return ccvCompatKeeper
}
