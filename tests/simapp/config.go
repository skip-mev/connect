package simapp

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	vestingmodulev1 "cosmossdk.io/api/cosmos/vesting/module/v1"
	"cosmossdk.io/depinject"

	alertmodulev1 "github.com/skip-mev/slinky/api/slinky/alerts/module/v1"
	incentivesmodulev1 "github.com/skip-mev/slinky/api/slinky/incentives/module/v1"
	oraclemodulev1 "github.com/skip-mev/slinky/api/slinky/oracle/module/v1"

	_ "cosmossdk.io/x/circuit"                        // import for side-effects
	_ "cosmossdk.io/x/upgrade"                        // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/vesting"   // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/authz/module"   // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"      // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/distribution"   // import for side-effects
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	_ "github.com/cosmos/cosmos-sdk/x/group/module" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/mint"         // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/params"       // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/slashing"     // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/staking"      // import for side-effects

	_ "github.com/skip-mev/slinky/x/incentives"       // import for side-effects
	_ "github.com/skip-mev/slinky/x/marketmap/module" // import for side-effects
	_ "github.com/skip-mev/slinky/x/oracle"           // import for side-effects

	"cosmossdk.io/core/appconfig"
	circuittypes "cosmossdk.io/x/circuit/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	marketmapmodulev1 "github.com/skip-mev/slinky/api/slinky/marketmap/module/v1"
	alerttypes "github.com/skip-mev/slinky/x/alerts/types"
	"github.com/skip-mev/slinky/x/alerts/types/strategies"
	incentivetypes "github.com/skip-mev/slinky/x/incentives/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// ProvideIncentives provides the incentive strategies for the incentive module, wrt the expected Keeper dependencies for
// incentive handler.
func ProvideIncentives(bk alerttypes.BankKeeper, sk alerttypes.StakingKeeper) map[incentivetypes.Incentive]incentivetypes.Strategy {
	return map[incentivetypes.Incentive]incentivetypes.Strategy{
		&strategies.ValidatorAlertIncentive{}: strategies.DefaultValidatorAlertIncentiveStrategy(sk, bk),
	}
}

var (
	// module account permissions.
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: oracletypes.ModuleName, Permissions: []string{}},
		{Account: incentivetypes.ModuleName, Permissions: []string{}},
		{Account: alerttypes.ModuleName, Permissions: []string{authtypes.Burner, authtypes.Minter}},
	}

	// blocked account addresses.
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		// We allow the following module accounts to receive funds:
		// govtypes.ModuleName
	}

	// application configuration (used by depinject).
	AppConfig = depinject.Configs(appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: "SimApp",
					// During begin block slashing happens after distr.BeginBlocker so that
					// there is nothing left over in the validator fee pool, so as to keep the
					// CanWithdrawInvariant invariant.
					// NOTE: staking module is required if HistoricalEntries param > 0
					BeginBlockers: []string{
						upgradetypes.ModuleName,
						minttypes.ModuleName,
						distrtypes.ModuleName,
						slashingtypes.ModuleName,
						stakingtypes.ModuleName,
						genutiltypes.ModuleName,
						authz.ModuleName,
						oracletypes.ModuleName,
						incentivetypes.ModuleName,
						alerttypes.ModuleName,
						marketmaptypes.ModuleName,
					},
					EndBlockers: []string{
						govtypes.ModuleName,
						stakingtypes.ModuleName,
						genutiltypes.ModuleName,
						group.ModuleName,
						oracletypes.ModuleName,
						// alert Endblock must precede incentives types EndBlocker (issued incentives should be executed same block)
						alerttypes.ModuleName,
						incentivetypes.ModuleName,
						marketmaptypes.ModuleName,
					},
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					// NOTE: The genutils module must occur after staking so that pools are
					// properly initialized with tokens from genesis accounts.
					// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
					InitGenesis: []string{
						authtypes.ModuleName,
						banktypes.ModuleName,
						distrtypes.ModuleName,
						stakingtypes.ModuleName,
						slashingtypes.ModuleName,
						govtypes.ModuleName,
						minttypes.ModuleName,
						genutiltypes.ModuleName,
						authz.ModuleName,
						group.ModuleName,
						paramstypes.ModuleName,
						upgradetypes.ModuleName,
						vestingtypes.ModuleName,
						consensustypes.ModuleName,
						circuittypes.ModuleName,
						oracletypes.ModuleName,
						incentivetypes.ModuleName,
						alerttypes.ModuleName,
						marketmaptypes.ModuleName,
					},
					// When ExportGenesis is not specified, the export genesis module order
					// is equal to the init genesis order
					// ExportGenesis: []string{},
					// Uncomment if you want to set a custom migration order here.
					// OrderMigrations: []string{},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             "cosmos",
					ModuleAccountPermissions: moduleAccPerms,
					// By default modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name:   vestingtypes.ModuleName,
				Config: appconfig.WrapAny(&vestingmodulev1.Module{}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   paramstypes.ModuleName,
				Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   oracletypes.ModuleName,
				Config: appconfig.WrapAny(&oraclemodulev1.Module{}),
			},
			{
				Name:   incentivetypes.ModuleName,
				Config: appconfig.WrapAny(&incentivesmodulev1.Module{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: appconfig.WrapAny(&authzmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: appconfig.WrapAny(&mintmodulev1.Module{}),
			},
			{
				Name: group.ModuleName,
				Config: appconfig.WrapAny(&groupmodulev1.Module{
					MaxExecutionPeriod: durationpb.New(time.Second * 1209600),
					MaxMetadataLen:     255,
				}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: appconfig.WrapAny(&govmodulev1.Module{}),
			},

			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   circuittypes.ModuleName,
				Config: appconfig.WrapAny(&circuitmodulev1.Module{}),
			},
			{
				Name: alerttypes.ModuleName,
				Config: appconfig.WrapAny(&alertmodulev1.Module{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				}),
			},
			{
				Name: marketmaptypes.ModuleName,
				Config: appconfig.WrapAny(&marketmapmodulev1.Module{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				}),
			},
		},
	}),
		depinject.Provide(alerttypes.ProvideMsgAlertGetSigners),
		depinject.Provide(ProvideIncentives),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
				govtypes.ModuleName: gov.NewAppModuleBasic(
					[]govclient.ProposalHandler{
						paramsclient.ProposalHandler,
					},
				),
			},

			// Supply the Incentive Handler for the Alerts module's ProvideModule Inputs
			strategies.DefaultHandleValidatorIncentive(),
		))
)
