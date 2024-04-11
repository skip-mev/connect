package integration

import (
	"context"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"

	slinkyabci "github.com/skip-mev/slinky/abci/ve/types"
	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/static"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type SlinkyABCIIntegrationSuite struct {
	*SlinkyIntegrationSuite
}

func NewSlinkyABCIIntegrationSuite(suite *SlinkyIntegrationSuite) *SlinkyABCIIntegrationSuite {
	return &SlinkyABCIIntegrationSuite{
		suite,
	}
}

// no op to prevent multiple calls to SetupSuite
func (s *SlinkyABCIIntegrationSuite) SetupSuite() {}

func (s *SlinkyABCIIntegrationSuite) TestCurrencyPairRemoval() {
	// initialize the oracle module with a currency-pair
	cps := []slinkytypes.CurrencyPair{
		slinkytypes.NewCurrencyPair("BTC", "USD"),
		slinkytypes.NewCurrencyPair("ETH", "USD"),
		slinkytypes.NewCurrencyPair("BTC", "ETH"),
	}
	s.addCurrencyPairs(context.Background(), s.chain, cps...)

	id1 := s.GetIDForCurrencyPair(context.Background(), cps[0])
	id2 := s.GetIDForCurrencyPair(context.Background(), cps[1])
	id3 := s.GetIDForCurrencyPair(context.Background(), cps[2])

	zeroBz, err := big.NewInt(0).GobEncode()
	s.Require().NoError(err)

	// expect prices to post
	_, err = ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
		{
			Prices: map[uint64][]byte{
				id1: zeroBz,
				id2: zeroBz,
				id3: zeroBz,
			},
		},
		{
			Prices: map[uint64][]byte{
				id1: zeroBz,
				id2: zeroBz,
				id3: zeroBz,
			},
		},
		{
			Prices: map[uint64][]byte{
				id1: zeroBz,
				id2: zeroBz,
				id3: zeroBz,
			},
		},
		{
			Prices: map[uint64][]byte{
				id1: zeroBz,
				id2: zeroBz,
				id3: zeroBz,
			},
		},
	})
	s.Require().NoError(err)

	_, err = RemoveCurrencyPairs(s.chain, s.authority.String(), sdk.NewCoin(s.denom, math.NewInt(deposit)), s.user.KeyName(), cps[0])
	s.Require().NoError(err)

	_, err = ExpectVoteExtensions(s.chain, s.blockTime*3, []slinkyabci.OracleVoteExtension{
		{
			Prices: map[uint64][]byte{
				id2: zeroBz,
				id3: zeroBz,
			},
		},
		{
			Prices: map[uint64][]byte{
				id2: zeroBz,
				id3: zeroBz,
			},
		},
		{
			Prices: map[uint64][]byte{
				id2: zeroBz,
				id3: zeroBz,
			},
		},
		{
			Prices: map[uint64][]byte{
				id2: zeroBz,
				id3: zeroBz,
			},
		},
	})
}

func (s *SlinkyABCIIntegrationSuite) GetIDForCurrencyPair(ctx context.Context, cp slinkytypes.CurrencyPair) uint64 {
	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	id, err := getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp)
	s.Require().NoError(err)

	return id
}

func (s *SlinkyABCIIntegrationSuite) addCurrencyPairs(ctx context.Context, chain *cosmos.CosmosChain, cps ...slinkytypes.CurrencyPair) {
	s.Require().NoError(s.AddCurrencyPairs(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.user, cps...))

	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	ids := make([]uint64, len(cps))

	for i, cp := range cps {
		ids[i], err = getIDForCurrencyPair(ctx, oracletypes.NewQueryClient(cc), cp)
		s.Require().NoError(err)
	}

	// start oracles
	for i := range s.chain.Nodes() {
		node := s.chain.Nodes()[i]
		oracleConfig := DefaultOracleConfig()
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
		tickers := make(map[string]mmtypes.Ticker)
		marketConfig := mmtypes.MarketMap{
			Markets: make(map[string]mmtypes.Market),
		}

		for _, cp := range cps {
			tickers[cp.String()] = mmtypes.Ticker{
				CurrencyPair:     cp,
				Decimals:         18,
				MinProviderCount: 1,
			}

			marketConfig.Markets[cp.String()] = mmtypes.Market{
				Ticker: tickers[cp.String()],
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           static.Name,
						OffChainTicker: cp.String(),
						Metadata_JSON:  `{"price": 0}`,
					},
				},
			}
		}

		oracle := GetOracleSideCar(node)
		SetOracleConfigsOnOracle(oracle, oracleConfig, marketConfig)
		s.Require().NoError(RestartOracle(node))
	}
}
