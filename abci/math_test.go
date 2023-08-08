package abci_test

import (
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/abci"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/x/oracle/types"

	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/golang/mock/gomock"

	storetypes "cosmossdk.io/store/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdkmock "github.com/cosmos/cosmos-sdk/testutil/mock"
)

type validator struct {
	stake    sdkmath.Int
	consAddr sdk.ConsAddress
}

var (
	validator1 = sdk.ConsAddress("validator1")
	validator2 = sdk.ConsAddress("validator2")
	validator3 = sdk.ConsAddress("validator3")
)

func (suite *ABCITestSuite) TestVoteWeightedMedian() {
	cases := []struct {
		name              string
		providerPrices    oracletypes.AggregatedProviderPrices
		validators        []validator
		totalBondedTokens sdkmath.Int
		expectedPrices    map[types.CurrencyPair]*uint256.Int
	}{
		{
			name:              "no providers",
			providerPrices:    oracletypes.AggregatedProviderPrices{},
			validators:        []validator{},
			totalBondedTokens: sdkmath.NewInt(100),
			expectedPrices:    map[types.CurrencyPair]*uint256.Int{},
		},
		{
			name: "single provider entire stake + single price",
			providerPrices: oracletypes.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
			providerPrices: oracletypes.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
			providerPrices: oracletypes.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
			providerPrices: oracletypes.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
				},
				validator2.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
			providerPrices: oracletypes.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(100),
					},
				},
				validator2.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
					{
						Base:  "BTC",
						Quote: "USD",
					}: {
						Price: uint256.NewInt(200),
					},
				},
				validator3.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
			providerPrices: oracletypes.AggregatedProviderPrices{
				validator1.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
				validator2.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
				validator3.String(): map[types.CurrencyPair]oracletypes.QuotePrice{
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
		suite.Run(tc.name, func() {
			// Create a mock validator store.
			mockValidatorStore := suite.createMockValidatorStore(tc.validators, tc.totalBondedTokens)
			// Compute the stake weighted median.
			aggregateFn := abci.VoteWeightedMedian(suite.ctx, mockValidatorStore, abci.DefaultPowerThreshold)
			result := aggregateFn(tc.providerPrices)

			// Verify the result.
			suite.Require().Len(result, len(tc.expectedPrices))
			for currencyPair, expectedPrice := range tc.expectedPrices {
				suite.Require().Equal(expectedPrice, result[currencyPair])
			}
		})
	}
}

func (suite *ABCITestSuite) TestComputeVoteWeightedMedian() {
	cases := []struct {
		name      string
		priceInfo abci.VoteWeightedPriceInfo
		expected  *uint256.Int
	}{
		{
			name: "single price",
			priceInfo: abci.VoteWeightedPriceInfo{
				Prices: []abci.VoteWeightedPricePerValidator{
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
			priceInfo: abci.VoteWeightedPriceInfo{
				Prices: []abci.VoteWeightedPricePerValidator{
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
			priceInfo: abci.VoteWeightedPriceInfo{
				Prices: []abci.VoteWeightedPricePerValidator{
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
			priceInfo: abci.VoteWeightedPriceInfo{
				Prices: []abci.VoteWeightedPricePerValidator{
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
			priceInfo: abci.VoteWeightedPriceInfo{
				Prices: []abci.VoteWeightedPricePerValidator{
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
		suite.Run(tc.name, func() {
			result := abci.ComputeVoteWeightedMedian(tc.priceInfo)
			suite.Require().Equal(tc.expected, result)
		})
	}
}

func (suite *ABCITestSuite) TestAggregationWithContext() {
	// create staking-keeper + context
	key := storetypes.NewKVStoreKey(stakingtypes.StoreKey)
	ec := moduletestutil.MakeTestEncodingConfig()

	ctrl := gomock.NewController(suite.T())
	bk := stakingtestutil.NewMockBankKeeper(ctrl)
	ak := stakingtestutil.NewMockAccountKeeper(ctrl)

	ctx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient"))

	const (
		bondedPoolAccAddr = "bonded-pool"
	)
	bondedPoolAcc := authtypes.NewModuleAccount(authtypes.NewBaseAccount(sdk.AccAddress(bondedPoolAccAddr), nil, 0, 0), stakingtypes.BondedPoolName)

	ak.EXPECT().GetModuleAddress(stakingtypes.BondedPoolName).Return(bondedPoolAcc.GetAddress())
	ak.EXPECT().GetModuleAddress(stakingtypes.NotBondedPoolName).Return(sdk.AccAddress("ignore-this"))
	ak.EXPECT().AddressCodec().Return(address.NewBech32Codec("cosmos")).AnyTimes()

	ak.EXPECT().GetModuleAccount(ctx.Ctx, stakingtypes.BondedPoolName).Return(bondedPoolAcc)
	bk.EXPECT().GetBalance(ctx.Ctx, bondedPoolAcc.GetAddress(), "stake").Return(sdk.NewInt64Coin("stake", 300))

	stakingKeeper := stakingkeeper.NewKeeper(ec.Codec, runtime.NewKVStoreService(key), ak, bk, sdk.AccAddress("does-not-matter").String(), ec.Codec.InterfaceRegistry().SigningContext().AddressCodec(), ec.InterfaceRegistry.SigningContext().ValidatorAddressCodec())
	stakingKeeper.SetParams(ctx.Ctx, stakingtypes.DefaultParams())

	valCons1 := sdk.ConsAddress("valcons1")
	valCons2 := sdk.ConsAddress("valcons2")
	valCons3 := sdk.ConsAddress("valcons3")

	// set three validators in state, w/ equal stake
	pk := sdkmock.NewPV().PrivKey.PubKey()

	fakePk, err := codectypes.NewAnyWithValue(pk)
	suite.Require().NoError(err)

	stakingstore := ctx.Ctx.KVStore(key)
	stakingstore.Set(stakingtypes.GetValidatorKey(sdk.ValAddress(valCons1)), stakingtypes.MustMarshalValidator(ec.Codec, &stakingtypes.Validator{
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(100),
		ConsensusPubkey: fakePk,
	}))
	stakingstore.Set(stakingtypes.GetValidatorByConsAddrKey(valCons1), sdk.ValAddress(valCons1))

	stakingstore.Set(stakingtypes.GetValidatorKey(sdk.ValAddress(valCons2)), stakingtypes.MustMarshalValidator(ec.Codec, &stakingtypes.Validator{
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(100),
		ConsensusPubkey: fakePk,
	}))
	stakingstore.Set(stakingtypes.GetValidatorByConsAddrKey(valCons2), sdk.ValAddress(valCons2))

	stakingstore.Set(stakingtypes.GetValidatorKey(sdk.ValAddress(valCons3)), stakingtypes.MustMarshalValidator(ec.Codec, &stakingtypes.Validator{
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(100),
		ConsensusPubkey: fakePk,
	}))
	stakingstore.Set(stakingtypes.GetValidatorByConsAddrKey(valCons3), sdk.ValAddress(valCons3))

	// create commit
	valCons1Vote := suite.createExtendedVoteInfo(valCons1, map[string]string{"BTC/USD": "0x123"}, time.Now(), 1)
	valCons2Vote := suite.createExtendedVoteInfo(valCons2, map[string]string{"BTC/USD": "0x456"}, time.Now(), 1)
	valCons3Vote := suite.createExtendedVoteInfo(valCons3, map[string]string{"BTC/USD": "0x789"}, time.Now(), 1)

	oracle := abci.NewOracle(
		log.NewTestLogger(suite.T()),
		abci.VoteWeightedMedianFromContext(
			stakingKeeper,
			abci.DefaultPowerThreshold,
		),
		suite.oracleKeeper,
		abci.NoOpValidateVoteExtensions,
		suite.validatorStore,
	)

	od, err := oracle.AggregateOracleData(ctx.Ctx, suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valCons1Vote, valCons2Vote, valCons3Vote}))
	suite.Require().NoError(err)
	// expect price to be middle
	suite.Require().Equal("0x456", od.Prices["BTC/USD"])

	// create new multi-store + context w/ diff powers, and execute
	ctx = testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient"))

	ak.EXPECT().GetModuleAccount(ctx.Ctx, stakingtypes.BondedPoolName).Return(bondedPoolAcc)
	bk.EXPECT().GetBalance(ctx.Ctx, bondedPoolAcc.GetAddress(), "stake").Return(sdk.NewCoin("stake", sdkmath.NewInt(300)))

	stakingstore = ctx.Ctx.KVStore(key)

	stakingKeeper.SetParams(ctx.Ctx, stakingtypes.DefaultParams())
	stakingstore.Set(stakingtypes.GetValidatorKey(sdk.ValAddress(valCons1)), stakingtypes.MustMarshalValidator(ec.Codec, &stakingtypes.Validator{
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(151),
		ConsensusPubkey: fakePk,
	}))
	stakingstore.Set(stakingtypes.GetValidatorByConsAddrKey(valCons1), sdk.ValAddress(valCons1))

	stakingstore.Set(stakingtypes.GetValidatorKey(sdk.ValAddress(valCons2)), stakingtypes.MustMarshalValidator(ec.Codec, &stakingtypes.Validator{
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(74),
		ConsensusPubkey: fakePk,
	}))
	stakingstore.Set(stakingtypes.GetValidatorByConsAddrKey(valCons2), sdk.ValAddress(valCons2))

	stakingstore.Set(stakingtypes.GetValidatorKey(sdk.ValAddress(valCons3)), stakingtypes.MustMarshalValidator(ec.Codec, &stakingtypes.Validator{
		Status:          stakingtypes.Bonded,
		Tokens:          sdkmath.NewInt(75),
		ConsensusPubkey: fakePk,
	}))
	stakingstore.Set(stakingtypes.GetValidatorByConsAddrKey(valCons3), sdk.ValAddress(valCons3))

	// create commit
	valCons1Vote = suite.createExtendedVoteInfo(valCons1, map[string]string{"BTC/USD": "0x123"}, time.Now(), 1)
	valCons2Vote = suite.createExtendedVoteInfo(valCons2, map[string]string{"BTC/USD": "0x456"}, time.Now(), 1)
	valCons3Vote = suite.createExtendedVoteInfo(valCons3, map[string]string{"BTC/USD": "0x789"}, time.Now(), 1)

	od, err = oracle.AggregateOracleData(ctx.Ctx, suite.createExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valCons1Vote, valCons2Vote, valCons3Vote}))
	suite.Require().NoError(err)
	// expect price to be middle
	suite.Require().Equal("0x123", od.Prices["BTC/USD"])
}
