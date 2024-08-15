package integration_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/tests/integration"
	"github.com/skip-mev/slinky/x/alerts"
	"github.com/skip-mev/slinky/x/incentives"
	marketmapmodule "github.com/skip-mev/slinky/x/marketmap"
	"github.com/skip-mev/slinky/x/oracle"
	"github.com/skip-mev/slinky/x/sla"
)

var (
	image = ibc.DockerImage{
		Repository: "skip-mev/connect-e2e",
		Version:    "latest",
		UidGid:     "1000:1000",
	}

	numValidators = 4
	numFullNodes  = 0
	noHostMount   = false
	gasAdjustment = 1.5

	oracleImage = ibc.DockerImage{
		Repository: "skip-mev/connect-e2e-oracle",
		Version:    "latest",
		UidGid:     "1000:1000",
	}
	encodingConfig = testutil.MakeTestEncodingConfig(
		bank.AppModuleBasic{},
		oracle.AppModuleBasic{},
		gov.AppModuleBasic{},
		alerts.AppModuleBasic{},
		auth.AppModuleBasic{},
		marketmapmodule.AppModuleBasic{},
		incentives.AppModuleBasic{},
		sla.AppModuleBasic{},
	)

	VotingPeriod     = "10s"
	MaxDepositPeriod = "1s"
	UnbondingTime    = "10s"

	defaultGenesisKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: VotingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: MaxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: denom,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.amount",
			Value: "0",
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
			Key:   "app_state.staking.params.unbonding_time",
			Value: UnbondingTime,
		},
	}

	denom = "stake"
	spec  = &interchaintest.ChainSpec{
		ChainName:     "connect",
		Name:          "connect",
		NumValidators: &numValidators,
		NumFullNodes:  &numFullNodes,
		Version:       "latest",
		NoHostMount:   &noHostMount,
		ChainConfig: ibc.ChainConfig{
			EncodingConfig: &encodingConfig,
			Images: []ibc.DockerImage{
				image,
			},
			Type:           "cosmos",
			Name:           "connect",
			Denom:          denom,
			ChainID:        "chain-id-0",
			Bin:            "connectd",
			Bech32Prefix:   "cosmos",
			CoinType:       "118",
			GasAdjustment:  gasAdjustment,
			GasPrices:      fmt.Sprintf("0%s", denom),
			TrustingPeriod: "48h",
			NoHostMount:    noHostMount,
			ModifyGenesis:  cosmos.ModifyGenesis(defaultGenesisKV),
		},
	}
)

func TestSlinkyOracleIntegration(t *testing.T) {
	baseSuite := integration.NewSlinkyIntegrationSuite(
		spec,
		oracleImage,
	)

	suite.Run(t, integration.NewSlinkyOracleIntegrationSuite(baseSuite))
}

func TestSlinkySlashingIntegration(t *testing.T) {
	baseSuite := integration.NewSlinkyIntegrationSuite(
		spec,
		oracleImage,
	)

	suite.Run(t, integration.NewSlinkySlashingIntegrationSuite(baseSuite))
}

func TestSlinkyOracleValidatorIntegration(t *testing.T) {
	baseSuite := integration.NewSlinkyIntegrationSuite(
		spec,
		oracleImage,
	)

	suite.Run(t, integration.NewSlinkyOracleValidatorIntegrationSuite(baseSuite))
}
