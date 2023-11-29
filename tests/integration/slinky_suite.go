package integration

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/pelletier/go-toml"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/stretchr/testify/suite"

	slinkyabci "github.com/skip-mev/slinky/abci/ve/types"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/providers/mock"
	service_metrics "github.com/skip-mev/slinky/service/metrics"
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
			"oracle",
			"--oracle-config-path", "/oracle/oracle.toml",
			"--metrics-config-path", "/oracle/metrics.toml",
			"--host", "0.0.0.0",
			"--port", "8080",
		},
		ValidatorProcess: true,
		PreStart:         true,
	}
}

func DefaultOracleConfig(node *cosmos.ChainNode) (
	oracleconfig.OracleConfig,
	oracleconfig.MetricsConfig,
) {
	oracle := GetOracleSideCar(node)

	// Create the oracle config
	oracleConfig := oracleconfig.OracleConfig{
		UpdateInterval: 500 * time.Millisecond,
		InProcess:      false,
		RemoteAddress:  fmt.Sprintf("%s:%s", oracle.HostName(), "8080"),
		Timeout:        500 * time.Millisecond,
	}

	// get the consensus address of the node
	bz, _, err := node.ExecBin(context.Background(), "cometbft", "show-address")
	if err != nil {
		panic(err)
	}
	consAddress := sdk.ConsAddress(bz)

	// Create the metrics config
	metricsConfig := oracleconfig.MetricsConfig{
		PrometheusServerAddress: fmt.Sprintf("%s:%s", oracle.HostName(), "8081"),
		OracleMetrics: metrics.Config{
			Enabled: true,
		},
		AppMetrics: service_metrics.Config{
			Enabled:              true,
			ValidatorConsAddress: consAddress.String(),
		},
	}

	return oracleConfig, metricsConfig
}

func GetOracleSideCar(node *cosmos.ChainNode) *cosmos.SidecarProcess {
	if len(node.Sidecars) == 0 {
		panic("no sidecars found")
	}
	return node.Sidecars[0]
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
			oracleCfg, metricsCfg := DefaultOracleConfig(node)
			SetOracleConfigsOnOracle(GetOracleSideCar(node), oracleCfg, metricsCfg)

			// set the out-of-process oracle config for all nodes
			node.WithPrestartNode(func(n *cosmos.ChainNode) {
				SetOracleConfigsOnApp(n, oracleCfg, metricsCfg)
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

	// reset the oracle services
	// start all oracles
	for _, node := range s.chain.Nodes() {
		oCfg, mCfg := DefaultOracleConfig(node)
		oCfg.Providers = append(oCfg.Providers, oracleconfig.ProviderConfig{
			Name: "static-mock-provider",
		})

		SetOracleConfigsOnOracle(GetOracleSideCar(node), oCfg, mCfg)
		RestartOracle(node)
	}

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

	// pass a governance proposal to approve a new currency-pair, and check Prices are reported
	s.Run("Add a currency-pair and check Prices", func() {
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

	// remove the currency-pair from state and check the Prices for that currency-pair are no longer reported
	s.Run("Remove a currency-pair and check Prices", func() {
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

	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	id, err := getIDForCurrencyPair(context.Background(), oracletypes.NewQueryClient(cc), cp)
	s.Require().NoError(err)

	// configure failing providers for various sets of nodes
	s.Run("all nodes report Prices", func() {
		// update all oracle configs
		for _, node := range s.chain.Nodes() {
			oracle := GetOracleSideCar(node)

			oracleConfig, metricsConfig := DefaultOracleConfig(node)
			oracleConfig.Providers = append(oracleConfig.Providers, config.ProviderConfig{
				Name: "static-mock-provider",
				Path: path.Join(oracle.HomeDir(), staticMockProviderConfigPath),
			})
			oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cp)

			// Write the static provider config to the node
			staticConfig := mock.StaticMockProviderConfig{
				TokenPrices: map[string]string{
					cp.ToString(): "1140",
				},
			}

			bz, err := toml.Marshal(staticConfig)
			s.Require().NoError(err)
			s.Require().NoError(oracle.WriteFile(context.Background(), bz, staticMockProviderConfigPath))

			SetOracleConfigsOnOracle(oracle, oracleConfig, metricsConfig)
			RestartOracle(node)
		}

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
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
		// stop single node's oracle process and check that all Prices are reported
		node := s.chain.Nodes()[0]
		StopOracle(node)

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)

		// expect update for height
		s.Require().Equal(newNonce, oldNonce+1)

		// start the oracle again
		StartOracle(node)
	})

	s.Run("single node down, price updates", func() {
		// stop single node's oracle process and check that all prices are reported
		node := s.chain.Nodes()[0]
		StopOracle(node)

		// expect the following vote-extensions
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height-1)
		s.Require().NoError(err)

		_, newNonce, err := QueryCurrencyPair(s.chain, cp, height)
		s.Require().NoError(err)

		// expect update for height
		s.Require().Equal(newNonce, oldNonce+1)

		StartOracle(node)
	})

	s.Run("only 1 node reports a price (oracles are down)", func() {
		// shut down all oracles except for one
		for _, node := range s.chain.Nodes()[1:] {
			StopOracle(node)
		}

		// expect the given oracle reports
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{},
			},
			{
				Prices: map[uint64][]byte{
					id: big.NewInt(1140).Bytes(),
				},
			},
		})
		s.Require().NoError(err)

		_, oldNonce, err := QueryCurrencyPair(s.chain, cp, height-1)
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

	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	// get the currency pair ids
	ctx := context.Background()
	id1, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp1)
	s.Require().NoError(err)

	id2, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp2)
	s.Require().NoError(err)

	id3, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp3)
	s.Require().NoError(err)

	// start all oracles
	for _, node := range s.chain.Nodes() {
		oracle := GetOracleSideCar(node)

		oracleConfig, metricsConfig := DefaultOracleConfig(node)
		oracleConfig.Providers = append(oracleConfig.Providers, config.ProviderConfig{
			Name: "static-mock-provider",
			Path: path.Join(oracle.HomeDir(), staticMockProviderConfigPath),
		})
		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cps...)

		// Write the static provider config to the node
		staticConfig := mock.StaticMockProviderConfig{
			TokenPrices: map[string]string{
				cp1.ToString(): "1140",
				cp2.ToString(): "1141",
				cp3.ToString(): "1142",
			},
		}

		bz, err := toml.Marshal(staticConfig)
		s.Require().NoError(err)
		s.Require().NoError(oracle.WriteFile(context.Background(), bz, staticMockProviderConfigPath))

		SetOracleConfigsOnOracle(oracle, oracleConfig, metricsConfig)
		RestartOracle(node)
	}

	s.Run("all oracles running for multiple price feeds", func() {
		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
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
		// stop first node's oracle, and update Prices to report
		node := s.chain.Nodes()[0]
		StopOracle(node)

		oracle := GetOracleSideCar(node)

		oracleConfig, metricsConfig := DefaultOracleConfig(node)
		oracleConfig.CurrencyPairs = append(oracleConfig.CurrencyPairs, cps...)
		oracleConfig.Providers = append(oracleConfig.Providers, config.ProviderConfig{
			Name: "static-mock-provider",
			Path: path.Join(oracle.HomeDir(), staticMockProviderConfigPath),
		})

		// Write the static provider config to the node
		staticConfig := mock.StaticMockProviderConfig{
			TokenPrices: map[string]string{
				cp1.ToString(): "1140",
				cp2.ToString(): "1141",
			},
		}

		bz, err := toml.Marshal(staticConfig)
		s.Require().NoError(err)
		s.Require().NoError(oracle.WriteFile(context.Background(), bz, staticMockProviderConfigPath))

		SetOracleConfigsOnOracle(oracle, oracleConfig, metricsConfig)
		RestartOracle(node)

		height, err := ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
				},
			},
			{
				Prices: map[uint64][]byte{
					id1: big.NewInt(1140).Bytes(),
					id2: big.NewInt(1141).Bytes(),
					id3: big.NewInt(1142).Bytes(),
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

func getIDForCurrencyPair(ctx context.Context, client oracletypes.QueryClient, cp oracletypes.CurrencyPair) (uint64, error) {
	// query for the given currency pair
	resp, err := client.GetPrice(ctx, &oracletypes.GetPriceRequest{
		CurrencyPairSelector: &oracletypes.GetPriceRequest_CurrencyPair{
			CurrencyPair: &cp,
		},
	})
	if err != nil {
		return 0, err
	}

	return resp.Id, nil
}
