package integration

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"

	connectabci "github.com/skip-mev/connect/v2/abci/ve/types"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/static"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

// ConnectCCVSuite is a testing-suite for testing Connect's integration with ics consumer chains
type ConnectCCVSuite struct {
	*ConnectIntegrationSuite
}

func NewConnectCCVIntegrationSuite(
	spec *interchaintest.ChainSpec, oracleImage ibc.DockerImage, opts ...Option,
) *ConnectCCVSuite {
	suite := NewConnectIntegrationSuite(spec, oracleImage, opts...)
	return &ConnectCCVSuite{
		ConnectIntegrationSuite: suite,
	}
}

func (s *ConnectCCVSuite) TestCCVAggregation() {
	ethusdc := connecttypes.NewCurrencyPair("ETH", "USDC")

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 3600, enabledTicker(ethusdc)))

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
		oracleConfig.Providers[static.Name] = oracleconfig.ProviderConfig{
			Name: static.Name,
			API: oracleconfig.APIConfig{
				Enabled:          true,
				Timeout:          250 * time.Millisecond,
				Interval:         250 * time.Millisecond,
				ReconnectTimeout: 250 * time.Millisecond,
				MaxQueries:       1,
				Endpoints: []oracleconfig.Endpoint{
					{
						URL: "http://un-used-url.com",
					},
				},
				Atomic: true,
				Name:   static.Name,
			},
			Type: types.ConfigType,
		}

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig)
		s.Require().NoError(RestartOracle(node))
	}

	// test that prices are reported as expected when stake-weight is the same across validators
	s.Run("expect a price-feed to be reported", func() {
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
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
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []connectabci.OracleVoteExtension{
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
