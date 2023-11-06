package integration

import (
	"context"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	oracleservicetypes "github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/stretchr/testify/suite"
	slinkyabci "github.com/skip-mev/slinky/abci/ve/types"
)

const (
	envKeepAlive  = "ORACLE_INTEGRATION_KEEPALIVE"
	genesisAmount = 1000000000
	defaultDenom  = "stake"
	validatorKey  = "validator"
	yes           = "yes"
	deposit       = 1000000
)

func DefaultOracleSidecar(image ibc.DockerImage) ibc.SidecarConfig {
	return ibc.SidecarConfig{
		ProcessName: "oracle",
		Image:       image,
		HomeDir:     "/oracle",
		Ports:       []string{"8080", "8081"},
		StartCmd: []string{
			"oracle", "-config", "/oracle/oracle.toml", "-host", "0.0.0.0", "-port", "8080",
		},
		ValidatorProcess: true,
		PreStart:         true,
	}
}

func DefaultOracleConfig() oracleconfig.Config {
	return oracleconfig.Config{
		Oracle: oracleconfig.Oracle{
			UpdateInterval: time.Millisecond,
		},
		Metrics: oracleconfig.Metrics{
			PrometheusServerAddress: "0.0.0.0:8081",
			OracleMetrics: oraclemetrics.Config{
				Enabled: true,
			},
		},
	}
}

type SlinkyIntegrationSuite struct {
	suite.Suite

	spec *interchaintest.ChainSpec

	// add more fields here as necessary
	chain *cosmos.CosmosChain

	// oracle side-car config
	oracleConfig ibc.SidecarConfig

	// user
	user cosmos.User

	// default token denom
	denom string

	// authority address
	authority sdk.AccAddress

	// block time
	blockTime time.Duration
}

func NewSlinkyIntegrationSuite(spec *interchaintest.ChainSpec, oracleImage ibc.DockerImage) *SlinkyIntegrationSuite {
	return &SlinkyIntegrationSuite{
		spec:         spec,
		oracleConfig: DefaultOracleSidecar(oracleImage),
		denom:        defaultDenom,
		authority:    authtypes.NewModuleAddress(govtypes.ModuleName),
		blockTime:    10 * time.Second,
	}
}

func (s *SlinkyIntegrationSuite) WithDenom(denom string) *SlinkyIntegrationSuite {
	s.denom = denom
	return s
}

func (s *SlinkyIntegrationSuite) WithAuthority(addr sdk.AccAddress) *SlinkyIntegrationSuite {
	s.authority = addr
	return s
}

func (s *SlinkyIntegrationSuite) WithBlockTime(t time.Duration) *SlinkyIntegrationSuite {
	s.blockTime = t
	return s
}

func (s *SlinkyIntegrationSuite) SetupSuite() {
	// create the chain
	s.chain = ChainBuilderFromChainSpec(s.T(), s.spec)

	s.chain.WithPrestartNodes(func(c *cosmos.CosmosChain) {
		// for each node in the chain, set the sidecars
		for i := range c.Nodes() {
			// pin
			node := c.Nodes()[i]
			// add sidecars to node
			AddSidecarToNode(node, s.oracleConfig)
			// set config for the oracle
			SetOracleConfig(node, DefaultOracleConfig())
			// set the out-of-process oracle config for all nodes
			node.WithPrestartNode(func(n *cosmos.ChainNode) {
				SetOracleOutOfProcess(n)
			})
		}
	})

	// start the chain
	BuildPOBInterchain(s.T(), context.Background(), s.chain)
	users := interchaintest.GetAndFundTestUsers(s.T(), context.Background(), s.T().Name(), genesisAmount, s.chain)
	s.user = users[0]
}

func (s *SlinkyIntegrationSuite) TearDownSuite() {
	// get the oracle integration-test suite keep alive env
	if ok := os.Getenv(envKeepAlive); ok == "" {
		return
	}

	// keep the chain running
	s.T().Log("Keeping the chain running")
	for {
	}
}

func (s *SlinkyIntegrationSuite) SetupTest() {
	// query for all currency-pairs
	resp, err := QueryCurrencyPairs(s.chain)
	s.Require().NoError(err)

	s.T().Log("Removing all currency-pairs", resp.CurrencyPairs)

	if len(resp.CurrencyPairs) == 0 {
		return
	}

	ids := make([]string, len(resp.CurrencyPairs))
	for i, cp := range resp.CurrencyPairs {
		ids[i] = cp.ToString()
	}

	// remove all currency-pairs
	s.Require().NoError(RemoveCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, ids...))
}

type SlinkyOracleIntegrationSuite struct {
	*SlinkyIntegrationSuite
}

func NewSlinkyOracleIntegrationSuite(suite *SlinkyIntegrationSuite) *SlinkyOracleIntegrationSuite {
	return &SlinkyOracleIntegrationSuite{
		SlinkyIntegrationSuite: suite,
	}
}

func (s *SlinkyOracleIntegrationSuite) TestOracleModule() {
	// query the oracle module grpc service for any CurrencyPairs
	s.Run("QueryCurrencyPairs - no currency-pairs reported", func() {
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 0)
	})

	// pass a governance proposal to approve a new currency-pair, and check prices are reported
	s.Run("Add a currency-pair and check prices", func() {
		s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []oracletypes.CurrencyPair{
			{
				Base:  "BTC",
				Quote: "USD",
			},
		}...))

		// check that the currency-pair is added to state
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 1)
		s.Require().Equal(resp.CurrencyPairs[0].Base, "BTC")
		s.Require().Equal(resp.CurrencyPairs[0].Quote, "USD")
	})

	// remove the currency-pair from state and check the prices for that currency-pair are no longer reported
	s.Run("Remove a currency-pair and check prices", func() {
		s.Require().NoError(RemoveCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []string{oracletypes.NewCurrencyPair("BTC", "USD").ToString()}...))

		// check that the currency-pair is added to state
		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 0)
	})

	s.Run("Add multiple Currency Pairs and remove 1", func() {
		cp1 := oracletypes.NewCurrencyPair("ETH", "USD")
		cp2 := oracletypes.NewCurrencyPair("BTC", "USD")
		s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []oracletypes.CurrencyPair{
			cp1, cp2,
		}...))

		resp, err := QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 2)

		// remove btc from state
		s.Require().NoError(RemoveCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []string{cp2.ToString()}...))

		// check that the currency-pair is removed from state
		resp, err = QueryCurrencyPairs(s.chain)
		s.Require().NoError(err)
		s.Require().True(len(resp.CurrencyPairs) == 1)
		s.Require().Equal(resp.CurrencyPairs[0].Base, "ETH")
		s.Require().Equal(resp.CurrencyPairs[0].Quote, "USD")
	})
}

func (s *SlinkyOracleIntegrationSuite) TestNodeFailures() {
	cp := oracletypes.CurrencyPair{
		Base:  "ETHEREUM",
		Quote: "USDC",
	}

	s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, []oracletypes.CurrencyPair{
		cp,
	}...))

	// configure failing providers for various sets of nodes
	s.Run("all nodes report prices", func() {
		// update all oracle configs
		for _, node := range s.chain.Nodes() {
			oCfg := DefaultOracleConfig()
			oCfg.Providers = append(oCfg.Providers, oracleservicetypes.ProviderConfig{
				Name: "static-mock-provider",
				TokenNameToMetadata: map[string]oracleservicetypes.TokenMetadata{
					"ETHEREUM/USDC": {
						Symbol: "1140",
					},
				},
			})

			SetOracleConfig(node, oCfg)
			RestartOracle(node)
		}

		height, err := ExpectVoteExtensions(s.chain, s.blockTime * 3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},				
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},			
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},				
			},
		})
		s.Require().NoError(err)
		// query for the given currency pair
		resp, _, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)
		s.Require().Equal(resp.Price.Int64(), int64(1140))
	})

	s.Run("single oracle down, price updates", func() {
		// stop single node's oracle process and check that all prices are reported
		node := s.chain.Nodes()[0]
		StopOracle(node)

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime * 3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[string]string{},				
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},			
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},				
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height - 1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)

		// expect update for height
		s.Require().Equal(newNonce, oldNonce + 1)

		// start the oracle again
		StartOracle(node)
	})

	s.Run("single node down, price updates", func() {
		// stop single node's oracle process and check that all prices are reported
		node := s.chain.Nodes()[3]
		node.StopContainer(context.Background())
		node.RemoveContainer(context.Background())

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime * 3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[string]string{},				
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},			
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},				
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height - 1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)

		// expect update for height
		s.Require().Equal(newNonce, oldNonce + 1)

		s.Require().NoError(node.CreateNodeContainer(context.Background()))
		s.Require().NoError(node.StartContainer(context.Background()))
	})

	s.Run("only 1 node reports a price (oracles are down)", func() {
		// shut down all oracles except for one
		for _, node := range s.chain.Nodes()[1:] {
			StopOracle(node)
		}

		// expect the given oracle reports
		height, err := ExpectVoteExtensions(s.chain ,s.blockTime * 3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[string]string{},
			},
			{
				Prices: map[string]string{},			
			},
			{
				Prices: map[string]string{},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
				},			
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height - 1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)

		// expect no update for the height
		s.Require().Equal(newNonce, oldNonce)

		// start all oracles again
		for _, node := range s.chain.Nodes()[1:] {
			StartOracle(node)
		}
	})
}

func (s *SlinkyOracleIntegrationSuite) TestMultiplePriceFeeds() {
	cp1 := oracletypes.NewCurrencyPair("ETHEREUM", "USDC")
	cp2 := oracletypes.NewCurrencyPair("ETHEREUM", "USDT")
	cp3 := oracletypes.NewCurrencyPair("ETHEREUM", "USD")

	// add multiple currency pairs
	cps := []oracletypes.CurrencyPair{
		cp1,
		cp2,
		cp3,
	}

	s.Require().NoError(AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, cps...))

	// start all oracles
	for _, node := range s.chain.Nodes() {
		oCfg := DefaultOracleConfig()
		oCfg.Providers = append(oCfg.Providers, oracleservicetypes.ProviderConfig{
			Name: "static-mock-provider",
			TokenNameToMetadata: map[string]oracleservicetypes.TokenMetadata{
				"ETHEREUM/USDC": {
					Symbol: "1140",
				},
				"ETHEREUM/USDT": {
					Symbol: "1141",
				},
				"ETHEREUM/USD": {
					Symbol: "1142",
				},
			},
		})

		SetOracleConfig(node, oCfg)
		RestartOracle(node)
	}

	s.Run("all oracles running for multiple price feeds", func() {
		height, err := ExpectVoteExtensions(s.chain, s.blockTime * 3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
		})
		s.Require().NoError(err)

		// query for the given currency pair
		for i, cp := range cps {
			resp, _, err := QueryCurrencyPair(s.chain, cp, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(1140+i), resp.Price.Int64())
		}
	})

	s.Run("all oracles running for multiple price feeds, except for one", func() {
		// stop first node's oracle, and update prices to report
		node := s.chain.Nodes()[0]
		StopOracle(node)

		oCfg := DefaultOracleConfig()
		oCfg.Providers = append(oCfg.Providers, oracleservicetypes.ProviderConfig{
			Name: "static-mock-provider",
			TokenNameToMetadata: map[string]oracleservicetypes.TokenMetadata{
				cp2.ToString(): {
					Symbol: "1141",
				},
				cp3.ToString(): {
					Symbol: "1142",
				},
			},
		})
		SetOracleConfig(node, oCfg)
		RestartOracle(node)

		height, err := ExpectVoteExtensions(s.chain, s.blockTime * 3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[string]string{
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},
			{
				Prices: map[string]string{
					"ETHEREUM/USDC": "0x474",
					"ETHEREUM/USDT": "0x475",
					"ETHEREUM/USD": "0x476",
				},
			},	
		})
		s.Require().NoError(err)

		// query for the given currency pair
		for i, cp := range cps {
			resp, _, err := QueryCurrencyPair(s.chain, cp, height)
			s.Require().NoError(err)
			s.Require().Equal(int64(1140+i), resp.Price.Int64())
		}
	})
}
