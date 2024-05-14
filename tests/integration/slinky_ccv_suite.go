package integration

import (
	"github.com/skip-mev/slinky/oracle/constants"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/skip-mev/slinky/providers/static"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	"time"
	"context"
	"math/big"
	"github.com/skip-mev/slinky/oracle/types"
	slinkyabci "github.com/skip-mev/slinky/abci/ve/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"fmt"
	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"strings"
	"encoding/hex"
	"encoding/json"
	"github.com/strangelove-ventures/interchaintest/v8"
)

// Type SlinkyCCVSuite is a testing-suite for testing slinky's integration with ics consumer chains
type SlinkyCCVSuite struct {
	*SlinkyIntegrationSuite
}

func NewSlinkyCCVIntegrationSuite(
	spec *interchaintest.ChainSpec, oracleImage ibc.DockerImage, opts ...Option,
) *SlinkyCCVSuite {
	suite := NewSlinkyIntegrationSuite(spec, oracleImage, opts...)
	return &SlinkyCCVSuite{
		SlinkyIntegrationSuite: suite,
	}
}

func (s *SlinkyCCVSuite) TestCCVAggregation() {
	ethusdc := constants.ETHEREUM_USDC

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 3600, ethusdc))

	cc, closeFn, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer closeFn()

	// get the currency pair ids
	ctx := context.Background()
	id, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdc)
	s.Require().NoError(err)

	zero := big.NewInt(0)
	zeroBz, err := zero.GobEncode()
	s.Require().NoError(err)

	// start all oracles
	for _, node := range s.chain.Nodes() {
		oracleConfig := DefaultOracleConfig(translateGRPCAddr(s.chain))
		oracleConfig.Providers = append(oracleConfig.Providers, oracleconfig.ProviderConfig{
			Name: static.Name,
			API: oracleconfig.APIConfig{
				Enabled:          true,
				Timeout:          250 * time.Millisecond,
				Interval:         250 * time.Millisecond,
				ReconnectTimeout: 250 * time.Millisecond,
				MaxQueries:       1,
				URL:              "http://un-used-url.com",
				Atomic:           true,
				Name:             static.Name,
			},
			Type: types.ConfigType,
		})

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig)
		s.Require().NoError(RestartOracle(node))
	}

	// test that prices are reported as expected when stake-weight is the same across validators
	s.Run("expect a price-feed to be reported", func() {
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: zeroBz,
				},
			},
		})
		s.Require().NoError(err)

		// query for the given currency pair
		resp, _, err := QueryCurrencyPair(s.chain, ethusdc, height)
		s.Require().NoError(err)
		s.Require().Equal(int64(360000000000), resp.Price.Int64())
	})

	// test that when provider stake-weight changes, the price changes accordingly
	s.Run("when stake weight changes, expect no price updates", func() {
		provider := s.chain.Provider
		if provider == nil {
			s.T().Skip("provider not found")
		}

		ctx := context.Background()
		providerValidators, err := provider.StakingQueryValidators(ctx, stakingtypes.BondStatusBonded)
		s.Require().NoError(err)

		// get the validator
		validator := providerValidators[0]

		// fund a new user
		users := interchaintest.GetAndFundTestUsers(s.T(), context.Background(), s.T().Name(), validator.Tokens.Mul(math.NewInt(2)), provider)

		// double this validator's stake
		tokens := validator.Tokens
		s.Require().NoError(provider.GetNode().StakingDelegate(ctx, users[0].KeyName(), validator.OperatorAddress, fmt.Sprintf("%s%s", tokens.String(), provider.Config().Denom)))

		// expect stake to have doubled for validator
		updatedValidator, err := provider.StakingQueryValidator(ctx, validator.OperatorAddress)
		s.Require().NoError(err)
	
		s.Require().Equal(tokens.Mul(math.NewInt(2)), updatedValidator.Tokens)

		// flush packets
		s.Require().NotNil(s.ic)
		provider.FlushPendingICSPackets(ctx, s.ic.Relayer(), s.ic.Reporter(), s.ic.IBCPath())

		// consensus address
		expectedConsensusAddress, err := pubKeyToAddress(validator.ConsensusPubkey)
		s.Require().NoError(err)

		// turn off the oracle for the validator who's stake has doubled
		for _, node := range s.chain.Nodes() {
			nodePkFile, err := node.PrivValFileContent(context.Background())
			s.Require().NoError(err)

			// unmarshal the private key
			var privValFile PrivValFile
			s.Require().NoError(json.Unmarshal(nodePkFile, &privValFile))

			if privValFile.Address == expectedConsensusAddress {
				// turn off the oracle
				StopOracle(node)
			}
		}

		// update the market
		s.UpdateCurrencyPair(s.chain, []mmtypes.Market{
			{
				Ticker: mmtypes.Ticker{
					CurrencyPair:     ethusdc,
					Decimals:         8,
					MinProviderCount: 1,
					Metadata_JSON:    "",
					Enabled:          true,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           static.Name,
						OffChainTicker: ethusdc.String(),
						Metadata_JSON:  fmt.Sprintf(`{"price": %f}`, 3000.0),
					},
				},
			},
		})

		priceDelta := big.NewInt(-60000000000)
		bz, err := priceDelta.GobEncode()

		// wait for the vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: bz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: bz,
				},
			},
			{
				Prices: map[uint64][]byte{
					id: bz,
				},
			},
		})
		s.Require().NoError(err)

		// expect the price to have remained the same
		resp, _, err := QueryCurrencyPair(s.chain, ethusdc, height)
		s.Require().NoError(err)
		s.Require().Equal(int64(360000000000), resp.Price.Int64())
	})
}

type PrivValFile struct {
	Address string `json:"address"`
}

// pubKeyToAddress converts a public key to an address
func pubKeyToAddress(pubKeyAny *codectypes.Any) (string, error) {
	// create a codec
	ir := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(ir)

	cryptocodec.RegisterInterfaces(ir)

	// decode the public key
	var pubKey cryptotypes.PubKey
	if err := cdc.UnpackAny(pubKeyAny, &pubKey); err != nil {
		return "", err
	}

	// get the address
	return strings.ToUpper(hex.EncodeToString(pubKey.Address().Bytes())), nil
}