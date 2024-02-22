package constants

import (
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// DISCLAIMER: The conversion paths are defined based on the available markets
// on the supported exchanges. Additionally, whether the conversion paths are
// utilized is dependent on the aggregation function that the main oracle
// utilizes. Please consult the oracle documentation for more information.

var (
	// ATOM_USD_PATHS defines the conversion paths for ATOM to USD.
	//
	// The first path is a direct conversion from ATOM to USD.
	// The second path is a conversion from ATOM to USDT to USD.
	// The third path is a conversion from ATOM to USDC to USD.
	ATOM_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from ATOM to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: ATOM_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
			{
				// Path from ATOM to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: ATOM_USDT.CurrencyPair,
						Invert:       false,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
			{
				// Path from ATOM to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: ATOM_USDC.CurrencyPair,
						Invert:       false,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
		},
	}

	// AVAX_USD_PATHS defines the conversion paths for AVAX to USD.
	//
	// The first path is a direct conversion from AVAX to USD.
	// The second path is a conversion from AVAX to USDT to USD.
	// The third path is a conversion from AVAX to USDC to USD.
	AVAX_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from AVAX to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: AVAX_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
			{
				// Path from AVAX to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: AVAX_USDT.CurrencyPair,
						Invert:       false,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
			{
				// Path from AVAX to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: AVAX_USDC.CurrencyPair,
						Invert:       false,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
		},
	}

	// BITCOIN_USD_PATHS defines the conversion paths for BITCOIN to USD.
	//
	// The first path is a direct conversion from BITCOIN to USD.
	// The second path is a conversion from BITCOIN to USDT to USD.
	// The third path is a conversion from BITCOIN to USDC to USD.
	BITCOIN_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from BITCOIN to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: BITCOIN_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
			{
				// Path from BITCOIN to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: BITCOIN_USDT.CurrencyPair,
						Invert:       false,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
			{
				// Path from BITCOIN to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: BITCOIN_USDC.CurrencyPair,
						Invert:       false,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
						Invert:       false,
					},
				},
			},
		},
	}

	// CELESTIA_USD_PATHS defines the conversion paths for CELESTIA to USD.
	//
	// The first path is a direct conversion from CELESTIA to USD.
	// The second path is a conversion from CELESTIA to USDT to USD.
	// The third path is a conversion from CELESTIA to USDC to USD.
	CELESTIA_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from Celestia to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: CELESTIA_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from Celestia to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: CELESTIA_USDT.CurrencyPair,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from Celestia to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: CELESTIA_USDC.CurrencyPair,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
					},
				},
			},
		},
	}

	// DYDX_USD_PATHS defines the conversion paths for DYDX to USD.
	//
	// The first path is a direct conversion from DYDX to USD.
	// The second path is a conversion from DYDX to USDT to USD.
	// The third path is a conversion from DYDX to USDC to USD.
	DYDX_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from DYDX to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: DYDX_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from DYDX to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: DYDX_USDT.CurrencyPair,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from DYDX to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: DYDX_USDC.CurrencyPair,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
					},
				},
			},
		},
	}

	// ETHEREUM_USD_PATHS defines the conversion paths for ETHEREUM to USD.
	//
	// The first path is a direct conversion from ETHEREUM to USD.
	// The second path is a conversion from ETHEREUM to USDT to USD.
	// The third path is a conversion from ETHEREUM to USDC to USD.
	ETHEREUM_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from ETHEREUM to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: ETHEREUM_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from ETHEREUM to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: ETHEREUM_USDT.CurrencyPair,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from ETHEREUM to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: ETHEREUM_USDC.CurrencyPair,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
					},
				},
			},
		},
	}

	// OSMOSIS_USD_PATHS defines the conversion paths for OSMOSIS to USD.
	//
	// The first path is a direct conversion from OSMOSIS to USD.
	// The second path is a conversion from OSMOSIS to USDT to USD.
	// The third path is a conversion from OSMOSIS to USDC to USD.
	OSMOSIS_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from OSMOSIS to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: OSMOSIS_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from OSMOSIS to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: OSMOSIS_USDT.CurrencyPair,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from OSMOSIS to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: OSMOSIS_USDC.CurrencyPair,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
					},
				},
			},
		},
	}

	// SOLANA_USD_PATHS defines the conversion paths for SOLANA to USD.
	//
	// The first path is a direct conversion from SOLANA to USD.
	// The second path is a conversion from SOLANA to USDT to USD.
	// The third path is a conversion from SOLANA to USDC to USD.
	SOLANA_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from SOLANA to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: SOLANA_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from SOLANA to USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: SOLANA_USDT.CurrencyPair,
					},
					{
						CurrencyPair: USDT_USD.CurrencyPair,
					},
				},
			},
			{
				// Path from SOLANA to USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: SOLANA_USDC.CurrencyPair,
					},
					{
						CurrencyPair: USDC_USD.CurrencyPair,
					},
				},
			},
		},
	}

	// USDC_USD_PATHS defines the conversion paths for USDC to USD.
	//
	// The first path is a direct conversion from USDC to USD.
	USDC_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from USDC to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: USDC_USD.CurrencyPair,
					},
				},
			},
		},
	}

	// USDT_USD_PATHS defines the conversion paths for USDT to USD.
	//
	// The first path is a direct conversion from USDT to USD.
	USDT_USD_PATHS = mmtypes.Paths{
		Paths: []mmtypes.Path{
			{
				// Direct path from USDT to USD
				Operations: []mmtypes.Operation{
					{
						CurrencyPair: USDT_USD.CurrencyPair,
					},
				},
			},
		},
	}
)
