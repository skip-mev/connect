package integration

import (
	"context"
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/static"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

type ConnectOracleValidatorIntegrationSuite struct {
	*ConnectIntegrationSuite
}

func NewConnectOracleValidatorIntegrationSuite(suite *ConnectIntegrationSuite) *ConnectOracleValidatorIntegrationSuite {
	return &ConnectOracleValidatorIntegrationSuite{
		ConnectIntegrationSuite: suite,
	}
}

func (s *ConnectOracleValidatorIntegrationSuite) TestUnbonding() {
	// Set up some price feeds
	ethusdcCP := connecttypes.NewCurrencyPair("ETH", "USDC")
	ethusdtCP := connecttypes.NewCurrencyPair("ETH", "USDT")
	ethusdCP := connecttypes.NewCurrencyPair("ETH", "USD")

	// add multiple currency pairs
	tickers := []mmtypes.Ticker{
		enabledTicker(ethusdcCP),
		enabledTicker(ethusdtCP),
		enabledTicker(ethusdCP),
	}

	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, tickers...))

	cc, closeFn, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer closeFn()

	// get the currency pair ids
	ctx := context.Background()
	_, err = getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdcCP)
	s.Require().NoError(err)

	_, err = getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdtCP)
	s.Require().NoError(err)

	_, err = getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), ethusdCP)
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

	// wait for the oracle to start
	s.T().Log("Waiting for oracle to start...")
	time.Sleep(45 * time.Second)
	s.T().Log("Done waiting for oracle to start")

	vals, err := s.chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	s.Require().NoError(err)
	val := vals[0].OperatorAddress

	wasErr := true
	for _, node := range s.chain.Validators {
		height, err := s.chain.Height(ctx)
		s.Require().NoError(err)

		s.T().Logf("attempting to unbond at height: %d", height)
		err = node.StakingUnbond(ctx, validatorKey, val, vals[0].BondedTokens().String()+s.denom)
		if err == nil {
			wasErr = false
		}
	}
	s.Require().False(wasErr)

	height, err := s.chain.Height(ctx)
	s.Require().NoError(err)
	s.T().Logf("unbond successful after height: %d", height)

	s.Eventually(
		func() bool {
			next, err := s.chain.Height(ctx)
			s.Require().NoError(err)
			s.T().Logf("current height: %d", next)
			return next > height+5
		},
		5*time.Minute,
		5*time.Second,
	)
}
