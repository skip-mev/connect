package integration

import (
	"context"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"fmt"
	"time"
)

var (
	providerChainID = "provider-1"
	providerNumValidators = int(4)
	providerVersion = "v5.0.0-rc0"
)

// CCVChainConstructor is a constructor for the CCV chain
func CCVChainConstructor(t *testing.T, spec *interchaintest.ChainSpec) []*cosmos.CosmosChain {
	// require that we only have 4 validators
	require.Equal(t, 4, *spec.NumValidators)

	cf := interchaintest.NewBuiltinChainFactory(
		zaptest.NewLogger(t),
		[]*interchaintest.ChainSpec{
			spec,
			{Name: "ics-provider", Version: providerVersion, NumValidators: &providerNumValidators, ChainConfig: ibc.ChainConfig{
				GasPrices: "0.0uatom",
				ChainID: providerChainID,
				TrustingPeriod: "336h",
				ModifyGenesis: cosmos.ModifyGenesis(
					[]cosmos.GenesisKV{
						cosmos.NewGenesisKV("app_state.provider.params.blocks_per_epoch", "1"),
					},
				),
			}},
		},
	)

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	require.Len(t, chains, 2)

	return []*cosmos.CosmosChain{chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)}
}

type CCVInterchain struct {
	relayer ibc.Relayer
	reporter *testreporter.RelayerExecReporter
	ibcPath string
}

func (c *CCVInterchain) Relayer() ibc.Relayer {
	return c.relayer
}

func (c *CCVInterchain) Reporter() *testreporter.RelayerExecReporter {
	return c.reporter
}

func (c *CCVInterchain) IBCPath() string {
	return c.ibcPath
}

// CCVInterchainConstructor is a constructor for the CCV interchain
func CCVInterchainConstructor(ctx context.Context, t *testing.T, chains []*cosmos.CosmosChain) Interchain {
	// expect 2 chains
	require.Len(t, chains, 2)

	// create a relayer
	client, network := interchaintest.DockerSetup(t)
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	).Build(t, client, network)

	path := "slinky-ibc-path"
	// create the interchain
	ic := interchaintest.NewInterchain().
		AddChain(chains[0]).
		AddChain(chains[1]).
		AddRelayer(r, "relayer").
		AddProviderConsumerLink(interchaintest.ProviderConsumerLink{
			Provider: chains[1],
			Consumer: chains[0],
			Relayer: r,
			Path: path,

		})
	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName: t.Name(),
		Client: client,
		NetworkID: network,
		SkipPathCreation: false,
	}))

	require.NoError(t, chains[1].FinishICSProviderSetup(ctx, r, eRep, path))

	return &CCVInterchain{
		relayer: r,
		reporter: eRep,
		ibcPath: path,
	}
}
