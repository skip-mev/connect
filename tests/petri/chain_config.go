package petri

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"go.uber.org/zap"

	"github.com/skip-mev/petri/chain/v2"
	"github.com/skip-mev/petri/node/v2"
	"github.com/skip-mev/petri/provider/v2"
	"github.com/skip-mev/petri/provider/v2/docker"
	"github.com/skip-mev/petri/types/v2"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	"github.com/skip-mev/connect/v2/tests/simapp"
)

func GetChainConfig() types.ChainConfig {
	return types.ChainConfig{
		Denom:         "stake",
		Decimals:      6,
		NumValidators: 4,
		NumNodes:      2,
		BinaryName:    "connectd",
		Image: provider.ImageDefinition{
			Image: "skip-mev/connect-e2e",
			UID:   "1000",
			GID:   "1000",
		},
		SidecarImage: provider.ImageDefinition{
			Image: "skip-mev/connect-e2e-oracle",
			UID:   "1000",
			GID:   "1000",
		},
		GasPrices:      "0stake",
		GasAdjustment:  1.5,
		Bech32Prefix:   "cosmos",
		EncodingConfig: GetEncodingConfig(),
		HomeDir:        "/petri-test",
		SidecarHomeDir: "/petri-test",
		SidecarPorts:   []string{"8080"},
		CoinType:       "118",
		ChainId:        "skip-1",
		ModifyGenesis:  GetGenesisModifier(),
		WalletConfig: types.WalletConfig{
			DerivationFn:     hd.Secp256k1.Derive(),
			GenerationFn:     hd.Secp256k1.Generate(),
			Bech32Prefix:     "cosmos",
			HDPath:           hd.CreateHDPath(0, 0, 0),
			SigningAlgorithm: "secp256k1",
		},
		UseGenesisSubCommand: true,
		NodeCreator:          node.CreateNode,
	}
}

func GetEncodingConfig() testutil.TestEncodingConfig {
	moduleBasics := make([]module.AppModuleBasic, len(simapp.ModuleBasics))
	i := 0
	for _, mb := range simapp.ModuleBasics {
		moduleBasics[i] = mb
		i++
	}
	return testutil.MakeTestEncodingConfig(
		moduleBasics...,
	)
}

func GetProvider(ctx context.Context, logger *zap.Logger) (provider.Provider, error) {
	return docker.NewDockerProvider(
		ctx,
		logger,
		"connect-docker",
	)
}

func GetChain(ctx context.Context, logger *zap.Logger) (types.ChainI, error) {
	prov, err := GetProvider(ctx, logger)
	if err != nil {
		return nil, err
	}
	return chain.CreateChain(
		ctx,
		logger,
		prov,
		GetChainConfig(),
	)
}

func GetGenesisModifier() types.GenesisModifier {
	genKVs := []chain.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: "10s",
		},
		{
			Key:   "app_state.gov.params.expedited_voting_period",
			Value: "5s",
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: "1s",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: "stake",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.amount",
			Value: "1",
		},
		{
			Key:   "app_state.gov.params.threshold",
			Value: "0.1",
		},
		{
			Key:   "app_state.gov.params.quorum",
			Value: "0",
		},
		{
			Key:   "consensus.params.abci.vote_extensions_enable_height",
			Value: "2",
		},
		{
			Key: "app_state.oracle.currency_pair_genesis",
			Value: []oracletypes.CurrencyPairGenesis{
				{
					CurrencyPair: connecttypes.CurrencyPair{
						Base:  "BITCOIN",
						Quote: "USD",
					},
					Id: 0,
				},
			},
		},
		{
			Key:   "app_state.oracle.next_id",
			Value: "1",
		},
	}
	return chain.ModifyGenesis(genKVs)
}
