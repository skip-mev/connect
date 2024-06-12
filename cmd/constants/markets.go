package constants

import (
	"encoding/json"
	"fmt"
	"os"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// RaydiumMarketMap is used to initialize the Raydium market map. This only includes
	// the markets that are supported by Raydium.
	RaydiumMarketMap mmtypes.MarketMap
	// RaydiumMarketMapJSON is the JSON representation of the Raydium MarketMap that can be used
	// to initialize for a genesis state or used by the sidecar as as static market map.
	RaydiumMarketMapJSON = `
	{
		"markets": {
		  "MOTHER/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "MOTHER",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
				{
					"name": "raydium_api",
					"off_chain_ticker": "SOL/MOTHER",
					"invert": true,
					"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"8uQwXPi1sWwUTVbDBnjznmf4mV44CETiNAh3UENvHejV\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"8ZcrNqaDbqy1H4R2DtmGnuZnJ6TKGSsaGmyRGQQeELLv\",\"token_decimals\":6},\"amm_info_address\":\"HcPgh6B2yHNvT6JsEmkrHYT8pVHu9Xiaoxm4Mmn2ibWw\",\"open_orders_address\":\"1z3rLR6AwR8gjVZ8ArqHp9kdaPrNwPvCdrZ6jGy6wwF\"}"
				}
			]
		  },
		  "$RETIRE/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "$RETIRE",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "$RETIRE/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2NJXbbLxfygPGusgTyxbFngLaodRRRCpXxeo1pv7M5XQ\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"8PcfqMPandh18bYZJKEvjuQRH5bcH4Y6TZzGUSWEPeYG\",\"token_decimals\":9},\"amm_info_address\":\"CQQDXt2M6Cx1J8N3cYsSmPiD7fcLdU5RpVtRbs9WaCXZ\",\"open_orders_address\":\"7uUa9ELGG1NQymmbNiUwYmkbrzhDzAdb1u9yHPNfnuZk\"}"
			  }
			]
		  },
		  "$WIF/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "$WIF",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "$WIF/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"7UYZ4vX13mmGiopayLZAduo8aie77yZ3o8FMzTeAX8uJ\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"7e9ExBAvDvuJP3GE6eKL5aSMi4RfXv3LkQaiNZBPmffR\",\"token_decimals\":9},\"amm_info_address\":\"D5MzuR2BVKhhLe5S2LiWHeLUv4QVr1mc2MC2PWdHZWtU\",\"open_orders_address\":\"CF3TGZTcBBv1ysxqTfRD2sP4Pu4TH6at62HU7Uc8mS3f\"}"
			  }
			]
		  },
		  "ANDY/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "ANDY",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "ANDY/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"pVCniSexvGFdDTjYuzoSKXDkoTqFjRhqJpmzzYMs7tY\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"2tt4w3C9hKjtzVgPqqa9Apbxz5qEKEAWodJEbbmfpwEm\",\"token_decimals\":9},\"amm_info_address\":\"72SvYQkXFwEUPtvRhhAEd7piardNzjzcSymHKrm9LPXA\",\"open_orders_address\":\"DW5dvuqw78jrmY8dgcPAtp7qJGEga9UgTm31Vea3q2oU\"}"
			  }
			]
		  },
		  "ANSEM/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "ANSEM",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "ANSEM/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5cM8YfCvkALkout2h9WxmYrc5e61YmvLgLLWsrZgumgK\",\"token_decimals\":1},\"quote_token_vault\":{\"token_vault_address\":\"B86KHBLhVVQQnsgbn6SDJR43NSKbqmpxZsrTai45yrMy\",\"token_decimals\":9},\"amm_info_address\":\"ArPbh31DT78N9Y1rQeABmcx5KYrUmVUcY9SrpfNjxUCB\",\"open_orders_address\":\"2gKn8iQ37nJXmQNU59Yb7WQiqjWGKECscGPRx1mCEK51\"}"
			  }
			]
		  },
		  "ATR/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "ATR",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "ATR/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4rm2qWHwGZGj9rWoRvU7m3FDdsZJV11wuHmczw27C3Wc\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Ahy4NhTyBaXZfsGjq4DDxGaMosBkAjaanGYdfeZjuDzP\",\"token_decimals\":6},\"amm_info_address\":\"2Ky6BskrcKNCJSrP4X6bgrPPe1erBArBAhyZi2C8nPwy\",\"open_orders_address\":\"DmY92VBBtKaNX2ZitPNZDX1G7VSrkpyeGwizXbCJc5ed\"}"
			  }
			]
		  },
		  "BAG/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "BAG",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BAG/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"7eLwyCqfhxKLsKeFwcN4JdfspKK22rSC4uQHNy3zWNPB\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Cr7Yo8Uf5f8pzMsY3ZwgDFNx85nb3UDvPfQxuWG4acxc\",\"token_decimals\":9},\"amm_info_address\":\"84n2mVWbHW2Gj4rkqQyhggEU548tdu4KtRmoB12eRqZp\",\"open_orders_address\":\"EupqWnZaYnYAACxVyf5yJ4dtwD6qJwLAsztaddN8RD5s\"}"
			  }
			]
		  },
		  "BODEN/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "BODEN",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BODEN/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"54zedUwxuSnmHHYg9oY1AfykeBDaCF6ZFZDW3ym2Nea4\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"DzpiXKsTUCacKyahLBUC5sfjj2fiWbwCpiCPEgyS3zDC\",\"token_decimals\":9},\"amm_info_address\":\"CiNw2KeeLT7oxj9e5H7eNmE9kQo74DasYUzuroNSY61i\",\"open_orders_address\":\"2vavMuz2sFMYwTKQ3U6EkC2uqAozwat9Bo7spisWky6H\"}"
			  }
			]
		  },
		  "BOME/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "BOME",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BOME/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FBba2XsQVhkoQDMfbNLVmo7dsvssdT39BMzVc2eFfE21\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"GuXKCb9ibwSeRSdSYqaCL3dcxBZ7jJcj6Y7rDwzmUBu9\",\"token_decimals\":9},\"amm_info_address\":\"AeBk8x4GVyWzLxmjWXA1WYkxkXaEsr8swFB3fExFy48n\",\"open_orders_address\":\"DyJKSuSFg6ZnNsmWUFWuCHswbF9muAVQRdKz3BkgSL5c\"}"
			  }
			]
		  },
		  "CHAT/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHAT",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CHAT/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FKCCPsYDgEoVpEhyE2XMFAXq5zWFrWHgpQjVEKQk1C54\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"9APPnbdEXbJpktfKSGbbdgYvj6K3ZFRDFwQUabFw6CHP\",\"token_decimals\":9},\"amm_info_address\":\"99tA8f5avoNLdNzWvMrKhTiXVSFsuCGMvw3btALS65fu\",\"open_orders_address\":\"BB95rL8vh9m922fMcVjBUCtHcWHUX3Jt9QmygNG4Naw5\"}"
			  }
			]
		  },
		  "CHEEMS/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHEEMS",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CHEEMS/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HRi4eJ8xWFG4hsv7FA9L7YnPLSxhQR4U7CCXKyZvcLXe\",\"token_decimals\":4},\"quote_token_vault\":{\"token_vault_address\":\"4gnEBvHQEx4nLcf9qWk1Wsxh9V1GFFDf4MfwEYmFo8hm\",\"token_decimals\":6},\"amm_info_address\":\"D9jaRQePbEXSDU17EDwcHuCfBUipqLRhQE1E9juYn2FY\",\"open_orders_address\":\"5G7co99pudpfd1Mp3xiGezJFTDvwMBhByaoZY4ttmtME\"}"
			  }
			]
		  },
		  "CHONKY/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHONKY",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CHONKY/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9rJqiGuRG971HCpapVNJtN4ho2fKMhkPiZRhQCAohonU\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"6Fbs4sCBH5jZp1v2Rr6zVdU44Fn4Vv9iPhY6eihjfdbz\",\"token_decimals\":9},\"amm_info_address\":\"E61pEDMEwf8iUHFhmGn3Wcj5P32DPjKDgo1UNjjaNrg1\",\"open_orders_address\":\"2PP6pkjc3QcQcB5qr6xg6gD4AWtjzDZzBDNebMQYFJqP\"}"
			  }
			]
		  },
		  "DUKO/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "DUKO",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "DUKO/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HD7aZ6YrqAxVbGNAMEKxozcW1ZDU7pbKfd7XMmZtxyzk\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"A9J2mXPXfRZ7Sh2ymUgCJM4p9iUjZBcyAfrz49PoBBN4\",\"token_decimals\":9},\"amm_info_address\":\"BGS69Ju7DRRVxw9b2B5TnrMLzVdJcscV8UtKywqNsgwx\",\"open_orders_address\":\"FoBQDGey332Ppv1KiTow8z9oZP8n6mEPLyhedPdG1nUG\"}"
			  }
			]
		  },
		  "FKETH/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "FKETH",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "FKETH/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"DNh9pRMiRX6zwWuzsXtmxmXLdbAGNuYg4dmmnzrpL871\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"47fCbQKnJYaMbPaPSyUrPXPUahizhYwAbwXnEcKN1vwD\",\"token_decimals\":9},\"amm_info_address\":\"HvAUaYpykFbUmzyGUCPbLR2nKA43cXspfxYNyYT2mw7j\",\"open_orders_address\":\"5VEfGLutckRLb3sFj9US8Dz4sFQ29xNMigpdxpTFj1bj\"}"
			  }
			]
		  },
		  "GMEOW/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "GMEOW",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "GMEOW/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9crtLiYfxYVYQ9sCfWix9vAMPJyBXCcMzCPXZ5isPFxB\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"ANLamVN6Df4HqC6YevQskovddsjhkqBqHsyLzhDibFEj\",\"token_decimals\":9},\"amm_info_address\":\"oykaqRr4gp6XRWxEraugTadtoMj2o1VMyvQCczf2GnQ\",\"open_orders_address\":\"9QWVztwM86spnAjAqGfVCCW8gycM69PMQ1hnq35Z2isR\"}"
			  }
			]
		  },
		  "HARAMBE/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "HARAMBE",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "HARAMBE/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5f9Fgcp2C9vdrp75GspNKBjzdaxq5uiqpLVkgtWKpDZZ\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Apgp3SzNB5VpVWbK5q2ucBvCJEsf1gqXL4iUAqvD9pgB\",\"token_decimals\":9},\"amm_info_address\":\"FT4S3M9MUPX7HKCGerTKN65rnHvwQpCrfpXFjS22Xbzd\",\"open_orders_address\":\"7ubVsaxLyMGyRQPmvHQ7uoj58jXmWFfpXgwPg7414eb3\"}"
			  }
			]
		  },
		  "KHAI/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "KHAI",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "KHAI/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"6g4rhxVTrN6SrtNvimq4QiU8yA5XScvwL6wxaMkegrtJ\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"7p2PrGMyeetNRqTKFraL7eYo2TbU3apWz6vfqrZFiPcG\",\"token_decimals\":9},\"amm_info_address\":\"ECbK6PSMZ5yQaUYBocsXaVrax2fWADw2ijTqLGPtt9sC\",\"open_orders_address\":\"2DaRg4UycKL9GSVfARBDrcensb89WD5PyyFX9NrMunLc\"}"
			  }
			]
		  },
		  "LIGMA/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "LIGMA",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "LIGMA/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"6aefaP7C6eKbW2taLqmyHinYH4ZMyY2G6MdqNu6PvfbL\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"G3kKDmmxwpjt4NVdQgdvgiuFxFsAsC1hSv4PVg63cKwM\",\"token_decimals\":9},\"amm_info_address\":\"GRz7v1v8AqpC1fyyjgAgnLk1PqjPcJ94gZC9Jur3f1fY\",\"open_orders_address\":\"E794XDGGZbnzrtCpU6RRUowknfUL619Q6SNi1oY24vST\"}"
			  }
			]
		  },
		  "LIKE/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "LIKE",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "LIKE/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"8LoHX6f6bMdQVs4mThoH2KwX2dQDSkqVFADi4ZjDQv9T\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"2Fwm8M8vuPXEXxvKz98VdawDxsK9W8uRuJyJhvtRdhid\",\"token_decimals\":6},\"amm_info_address\":\"GmaDNMWsTYWjaXVBjJTHNmCWAKU6cn5hhtWWYEZt4odo\",\"open_orders_address\":\"Crn5beRFeyj4Xw13E2wdJ9YkkLLEZzKYmtTV4LFDx3MN\"}"
			  }
			]
		  },
		  "MARVIN/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "MARVIN",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MARVIN/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"3vLTpZneGAzceAXKu2HuesT4rt6ksRJ3Q9WvjUmwksqA\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"3rWa7PstXZk4ZaEhLamfoqMVozwq7hfXEDqyNbHcL4uK\",\"token_decimals\":9},\"amm_info_address\":\"GR93XuWCWTqYpLCjqdJex8VdiSNyBhqKykNCjqfKdh3Q\",\"open_orders_address\":\"AMJqMphxaUXKsqy2Ymv2vamE4sByTfuNdBknDo1PRseq\"}"
			  }
			]
		  },
		  "MEW/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "MEW",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MEW/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4HqAjFKuQX7tnXEkxsjGQha3G4bcgER8qPXRahn9gj8F\",\"token_decimals\":5},\"quote_token_vault\":{\"token_vault_address\":\"BhNdEGJef9jSqT1iCEkFZ2bYZCdpC1vuiWtqDt87vBVp\",\"token_decimals\":9},\"amm_info_address\":\"7QQmXMuMDsrntXtzCc62twLLjDBmuQKKj9mUQVjiJ7te\",\"open_orders_address\":\"EGo9sULUczjtUKigCoykSZ3GqBQKL7d8hVAwur32VHq1\"}"
			  }
			]
		  },
		  "MONK/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "MONK",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MONK/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"CeLZLhC2nScSpsKqRL1eRr3L3LLfjDzakZLCUKcUHW1m\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"Be6DM12uTWtxHMSRJeah3J5PRP4CumR28Yy2qpLQFTE7\",\"token_decimals\":9},\"amm_info_address\":\"5mYCuaXmqW1McUe18Ry6gbWUQhtk1f4GxJ9j7vRj34HJ\",\"open_orders_address\":\"F1KFumMDuNonPprwUxarH6bbip9TR5wfsCZzRuKM8XBM\"}"
			  }
			]
		  },
		  "MOUTAI/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "MOUTAI",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MOUTAI/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4VkmzH14ETcNhSQLTK6AtL1ZP8UmvWpbNCgokDVfiCcD\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"HSQPdDCxtGo4fTHeZuBGWtQUqHgRsgdz4BVhTCCAtsTv\",\"token_decimals\":9},\"amm_info_address\":\"6AKb3QfbWTX57swhg39h71YjKb31LGLy3jFyQYNHDaCC\",\"open_orders_address\":\"HTDePdRKDZMASEZ1PYj6QBNKeBK4buAw7v5VVaGzmj7V\"}"
			  }
			]
		  },
		  "MPLX/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "MPLX",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MPLX/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5zdZza5N2TzV7cPtLeqCZQQRYCCHFVdXWLMeJo75DK24\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"7mwk7ZBiTtrKkKC5o34gpFBSCabEvLkp2fLjGNz43PyM\",\"token_decimals\":6},\"amm_info_address\":\"F3enDARqxRtTDuD8RxYLFXUEfPgph9VQLE6HXmWDqTSS\",\"open_orders_address\":\"ANiHc2R9qkDJxj9kLr6Bfwkw51niaK8q22CkuFTTUbTd\"}"
			  }
			]
		  },
		  "MUMU/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "MUMU",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MUMU/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2Re1H89emr8hNacyDTrm1NU8VEhuwaJX7JwcdDqy5Q6g\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"4VPXFMpndqZhME27vMqtkuGtBo7hVTA9kEvo87zbjXsA\",\"token_decimals\":9},\"amm_info_address\":\"9LhsF1JbMamXTaky4dzG1tiPFJ15nMtFXKPuRF87pqUR\",\"open_orders_address\":\"J4nb5Fj1EwB3q5uGNbJcb1bWxBhpN5QZbyP6num5VTYT\"}"
			  }
			]
		  },
		  "NICK/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "NICK",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NICK/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FDvQiwbJVHdSZE3ngZ6WCkJfciFTdg958W7bxyKU2PJ9\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"z6ZCZESyof3ZgCJ23hY31f1SSD33gQgyVRfMB8wP9iM\",\"token_decimals\":9},\"amm_info_address\":\"HrYiUQjzsZHF2t6ytm8rjUARZ3suCGi2Do2goZdJaL8Z\",\"open_orders_address\":\"4vj51zF3atwEr5PhWrT5zhU6jxqUapsusTk1HTTYSk3k\"}"
			  }
			]
		  },
		  "NINJA/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "NINJA",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NINJA/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5Mmie9Drh6RGMy8X8UQ3egyBi4Hvva1TR778bf77ViCV\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"HxVbv76N8EoRGEPJsKdtWCu3mz7ZXJi8dbZy8kM3QL3i\",\"token_decimals\":9},\"amm_info_address\":\"5oAhuuZMK4bCT8AGxxDQH7oHmLHGa9byDauHB3mvq6ie\",\"open_orders_address\":\"5TAZMKorJdh6ycqNTURsJaPWjeYNrWzNGkDPFAnnTRiM\"}"
			  }
			]
		  },
		  "NOS/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "NOS",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NOS/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9Gs4LvFZw18EBLrSmZbQBw4G2SpTu4bJRCWH1Dz33cUZ\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"FqKU4BxbabPd1tcZAVVv8JkdUWmdz32CocRM856gA3Lw\",\"token_decimals\":6},\"amm_info_address\":\"Cj7kD2VmzwSrwKBieuYYbjPEvr8gwhNi76KUESbGDNfF\",\"open_orders_address\":\"DEM1Zse8UWfKEk9dH1Jkjzepdb9DSaxMZ8uDe34rmbE5\"}"
			  }
			]
		  },
		  "NUB/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "NUB",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NUB/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9uNqUwneLXbQ6YKndciL5aBXTLJhwpyDXkZmaBbWfwWz\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"75DrZt3zmGSFfKaYDm7yHLKMrr35Wy8ffBNN1143PWbj\",\"token_decimals\":9},\"amm_info_address\":\"v4DyhFrb18Cjg8K8PoDR6Dr49A4Fu382KrTwKeYKXwX\",\"open_orders_address\":\"7nz62HHmxe7tmC32GJZPTCZDnG3Hui2fYd2qdUDyvRrk\"}"
			  }
			]
		  },
		  "PENG/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "PENG",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "PENG/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2g5q7fBGKZm2CXix8JjK4ZFdBTHQ1LerxkseBTqWuDdD\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"GmLJXUzjQAAU86a91hKesg5P9pKb6p9AZaGBEZLaDySD\",\"token_decimals\":9},\"amm_info_address\":\"AxBDdiMK9hRPLMPM7k6nCPC1gRARgXQHNejfP2LvNGr6\",\"open_orders_address\":\"9E5VWkY1UsbhkXW4Lk1YovkVouWMG57CuCNXUmecrGpC\"}"
			  }
			]
		  },
		  "PONKE/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "PONKE",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "PONKE/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"D7rw7fyEzo9EQcozjqAHJwbdbywGcSLw1at5MioZtMZ4\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"8DcvfWidQ53a3SCBrWxBWL2UU5zEBAKEypApiiCCFu2Y\",\"token_decimals\":9},\"amm_info_address\":\"6iNTc3N5f9FkDmchV89fizATcKcqP1UWukeREv1CcvPD\",\"open_orders_address\":\"2ohiiPwvoeZakCQ4gZ9dMYyGc84ZbGqMhbkgam7tCoxL\"}"
			  }
			]
		  },
		  "POPCAT/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "POPCAT",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "POPCAT/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4Vc6N76UBu26c3jJDKBAbvSD7zPLuQWStBk7QgVEoeoS\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"n6CwMY77wdEftf2VF6uPvbusYoraYUci3nYBPqH1DJ5\",\"token_decimals\":9},\"amm_info_address\":\"4qZaCqY3cq7t17Nabnx9aYjE45C7M4ZCKo9DYgXa4RPe\",\"open_orders_address\":\"EN95xVYrG6zHJFsHD2j19HCAD3JN8EvKCpcS6TqanUbn\"}"
			  }
			]
		  },
		  "POPCAT/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "POPCAT",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "POPCAT/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Ej1CZHUcHYcqAx3pJXUvqCTs3diVmEWSfozvQQLsQkyU\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"pTJo5c4ynoxxRAgDyWgQKasR8dgqQHP3CSeoXzDgZvZ\",\"token_decimals\":6},\"amm_info_address\":\"HBS7a3br8GMMWuqVa7VB3SMFa7xVi1tSFdoF5w4ZZ3kS\",\"open_orders_address\":\"9Q1E7B4w5Vhb5RjbmojpEuZbMZ944m9vDZZJoxBcGBRS\"}"
			  }
			]
		  },
		  "RAY/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RAY/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Em6rHi68trYgBFyJ5261A2nhwuQWfLcirgzZZYoRcrkX\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"3mEFzHsJyu2Cpjrz6zPmTzP7uoLFj9SbbecGVzzkL1mJ\",\"token_decimals\":9},\"amm_info_address\":\"E64ER3SCFZwqB2BWJFa6QvutLwAk35JfK2wFwZir46vz\",\"open_orders_address\":\"Bsdx1n9o449F45xnnGw9cu3w9E3WwRPGfHCHJ9KJ8iV2\"}"
			  }
			]
		  },
		  "RAY/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RAY/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FdmKUE4UMiJYFK5ogCngHzShuVKrFXBamPWcewDr31th\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"Eqrhxd7bDUCH3MepKmdVkgwazXRzY6iHhEoBpY7yAohk\",\"token_decimals\":6},\"amm_info_address\":\"EVzLJhqMtdC1nPmz8rNd6xGfVjDPxpLZgq7XJuNfMZ6\",\"open_orders_address\":\"3PJpMTYq82xmnUUERgsPsPKGk3jLEZS8wgNEHq91Yfhw\"}"
			  }
			]
		  },
		  "RAY/USDT": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "USDT"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RAY/USDT",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"3wqhzSB9avepM9xMteiZnbJw75zmTBDVmPFLTQAGcSMN\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"5GtSbKJEPaoumrDzNj4kGkgZtfDyUceKaHrPziazALC1\",\"token_decimals\":6},\"amm_info_address\":\"DVa7Qmb5ct9RCpaU7UTpSaf3GVMYz17vNVU67XpdCRut\",\"open_orders_address\":\"8MSPLj4c2hi1fZGDARvxLXQp1ooDQ8iGnWXbGdwvZxUQ\"}"
			  }
			]
		  },
		  "RETARDIO/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "RETARDIO",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RETARDIO/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HXzTvbuKKPyNMmLKJb8vaSUaRZsVS2J2AAsDuDm36rNC\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"HNcAAdLKHSRnwdmmWCYnP5Zcd11sfGpAoCuWFtugt2ma\",\"token_decimals\":9},\"amm_info_address\":\"3BWnpmiRG5nmzmsyVKNPcLeGKPMFWfuKLNMgSsLcLhnJ\",\"open_orders_address\":\"7LKjTUqBaYJ5qfaPEqn4ZDkoh39BPPUosBGaoeJkPXWq\"}"
			  }
			]
		  },
		  "SLERF/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "SLERF",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SLERF/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9A2ZsPz5Zg6jKN4o4KRMjTVPmkH51wYWFLmt4KBRy1Rq\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"5Zumc1SYPmQ89nqwXqzogeuhdJ85iEMpSk35A4P87pmD\",\"token_decimals\":9},\"amm_info_address\":\"AgFnRLUScRD2E4nWQxW73hdbSN7eKEUb2jHX7tx9YTYc\",\"open_orders_address\":\"FT5Ptk37g5r6D9BKt3hne8ovHZ1g56oJBvuZRwn3zS3j\"}"
			  }
			]
		  },
		  "SMOLE/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "SMOLE",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SMOLE/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2\",\"token_decimals\":9},\"amm_info_address\":\"J48ZPSiZ6G6TrZcC2PEMf3d2cYVoYRADMeUY3CTdsFqQ\",\"open_orders_address\":\"DURtPB9126j1TsTyv67o1Rd32E2HvZF4fYXNg1u5yE9P\"}"
			  }
			]
		  },
		  "SOL/EPIK": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "EPIK"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/EPIK",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9aGBQqKRyC5bbrZsnZJJtp59EqJj7vBkgV3HehgKEu5y\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"ANpMJb9ToMGNivLEdmBNBC2Qcf5ASaZkEdmUddV1FUZB\",\"token_decimals\":6},\"amm_info_address\":\"AZaaQaRhp1ys9VaJBRZYbmPz3JSBSp7m8cSSrLBn4BP9\",\"open_orders_address\":\"FjCKdnpN1t262QGGn6chWYRtoSaY6fuYxyKoqhinyGEK\"}"
			  }
			]
		  },
		  "SOL/HOBBES": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "HOBBES"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/HOBBES",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4J3cZy8G2qm6MBVGfeXhYETZvbRThv9TPPeY3p83QYLb\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"7DejqJN5iRcuUhR7C1Vif3SbjTXKCzkpyS3AxC28tiaF\",\"token_decimals\":6},\"amm_info_address\":\"9VffBiow5r5YQzgK56rirEWpu45gZGrDWzm9JUt6zL9G\",\"open_orders_address\":\"9q3x5stYdC6xuxdNjYQCRVktdHZiVrqWw2qcShrAQB2b\"}"
			  }
			]
		  },
		  "SOL/MAIL": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "MAIL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/MAIL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"BtJXP2iWPGW2x3EsscHtCuFLBECRCfrxsJ2SDi9jh96C\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"GmHpq7Wgri9TFitGj89quvwRc1ALhe7dePM6VBAiqxrC\",\"token_decimals\":6},\"amm_info_address\":\"9DYGj7g2b5fipk9wGsUhxdv5zpfTsoGzCiS29vH8Cfrs\",\"open_orders_address\":\"CLxCfhwK9SQAPLHu3KpMTLxvunG4JwRYxd5wY3tuNqQL\"}"
			  }
			]
		  },
		  "TOOKER/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "TOOKER",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TOOKER/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Cki9WdL3sCoNY3cLmfG4iqSbvB8g1Fr9tw8qa5tP1m3Y\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"2vTTh5pGbzc6guAJmt78XnTcXVBEZEWmGBkXkSNZwN59\",\"token_decimals\":9},\"amm_info_address\":\"F2JctNnqLJyf7mpUNadgZNpVLazX9GxjVeQT2FHv1QpM\",\"open_orders_address\":\"Gg4cXYqn6guMbWfuNdgfqM3AKeW3z5ZfV9NdxEwA7rYU\"}"
			  }
			]
		  },
		  "TREMP/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "TREMP",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TREMP/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"G2XNC6Rt2G7JZQWhqpJriYwZyxd2L52KSDbDNBCYCpvx\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"9DfnSR9h3hrmgy5pjqBP3SrVQRWPfjSqZZBrNNYGoyaN\",\"token_decimals\":9},\"amm_info_address\":\"EgF5PqgH3Y5dRnwQQ4eZhJ5JqktrYqNWSuJfjYpThkLb\",\"open_orders_address\":\"2XxdNxjjCyFsQSHSdoQUoSRhic9N4ppdgJpHk8JM2mYC\"}"
			  }
			]
		  },
		  "TRUMP/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRUMP",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TRUMP/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"GwUn3JYQ5PC6P9cB9HsatJN7y7aw3BXrNqHPpoHtWyKF\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"5DNPt6WYgj3x7EmU4Fyqe3jDYPk2HMAB21H5N4Ggbev9\",\"token_decimals\":9},\"amm_info_address\":\"CyZZM9ytS7ZREPqBXY2bTLUN95WXSyzVDqBzatgRzVyT\",\"open_orders_address\":\"5mUgTJLVhQ8McD7rpm4cvcjjgMqQUb3dJxkiVCn5xXMC\"}"
			  }
			]
		  },
		  "USDC/$MYRO": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDC",
				"Quote": "$MYRO"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "USDC/$MYRO",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"AMtPGYQS873njs35mD9MAAMKoospEuzNHPy7LQuuKo4A\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"BUvMbqP311JDU4ZGnf1rSZJLjutTU9VpNLEos393TYyW\",\"token_decimals\":9},\"amm_info_address\":\"5WGYajM1xtLy3QrLHGSX4YPwsso3jrjEsbU1VivUErzk\",\"open_orders_address\":\"2w1mZXi8PNqUz4gbezu4xvPzcGogDmVdLXGLhTpPbczd\"}"
			  }
			]
		  },
		  "VCAT/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "VCAT",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "VCAT/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"8gNjWm2wGubUiAXT9cXBeoY2NqAFkrnvKkh9J3gHZ7Wn\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"7G9HpLoYVhcBsg7ZEy928iUuzetJFK4AWBcfaCQTMp72\",\"token_decimals\":9},\"amm_info_address\":\"m9uiXqNMAxP7BNmtRf4NwkeKExjT5Hc6ftyjtNJB6E5\",\"open_orders_address\":\"9qnrgiBExCuugfUrK4Hgb3treG1YbBsNRqyk4cGDPzMF\"}"
			  }
			]
		  },
		  "VONSPEED/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "VONSPEED",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "VONSPEED/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"89CwpgTSsCc9u187kKvQQo6VAL5gKZViVub4eaNXfrtu\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"2fEchHP5r5wks9PPN1C2S39wPCe6Ni4247oMMTausc2a\",\"token_decimals\":9},\"amm_info_address\":\"34yiYY6kZmnVfmdQGtv2HugiNB5g1DcMDuc2VdckidB7\",\"open_orders_address\":\"CzfQYWMuLpnxCNMmQbwezkkivLD8Mn7W26MDNLgPtdS7\"}"
			  }
			]
		  },
		  "WHALES/SOL": {
			"ticker": {
			  "currency_pair": {
				"Base": "WHALES",
				"Quote": "SOL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "WHALES/SOL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"P9uSMnNEGHDP7Dhu7fKWfRViAGGHjEMv6urC8c2qG4k\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"CktEbT37HFRtwXVjwPEVfXHdcTAnqnmCvkgAw9SEN7zf\",\"token_decimals\":9},\"amm_info_address\":\"3dQ92wmjEnJEnZzxM9C3uSsd1MFAByrM5LbEDW2DCcxM\",\"open_orders_address\":\"B1FAXAPa7ZTUscypEYY5Xsb9AZyCGF8LDQHc7xC8K8Rn\"}"
			  }
			]
		  }
		}
	  }
	`

	// CoreMarketMap is used to initialize the Core market map.
	CoreMarketMap mmtypes.MarketMap
	// CoreMarketMapJSON is the JSON representation of the Core MarketMap that can be used
	// to initialize for a genesis state or used by the sidecar as as static market map.
	CoreMarketMapJSON = `
	{
		"markets": {
		  "AAVE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AAVE",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "AAVEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "aaveusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "AAVEUSD"
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AAVE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "AAVE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "binance_api",
				"off_chain_ticker": "AAVEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "AAVE-USD"
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "AAVE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ADA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ADA",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ADAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ADAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ADA-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ADA_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "adausdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ADAUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ADA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ADAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ADA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ADA_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "AEVO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AEVO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "AEVOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AEVOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "AEVO_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "AEVOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AEVO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "AGIX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AGIX",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "AGIXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AGIXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "AGIX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "AGIX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AGIX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "AGIXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ALGO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ALGO",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ALGOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ALGOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ALGO-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ALGOUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ALGO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ALGOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ALGO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "APE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "APE",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "APEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "APE-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "APE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "APEUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "APE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "APEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "APE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "APE_USD"
			  }
			]
		  },
		  "APT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "APT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "APTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "APTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "APT-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "APT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "aptusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "APT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "APTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "APT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "APT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ARB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ARB",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ARBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ARBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ARB-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ARB_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "arbusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ARB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ARBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ARB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ARB_USD"
			  }
			]
		  },
		  "ARKM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ARKM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ARKMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ARKMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ARKM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ARKM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ARKMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ASTR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ASTR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ASTRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ASTR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ASTRUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ASTR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ASTRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ASTR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ATOM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ATOM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ATOMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ATOMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ATOM-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ATOM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ATOMUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ATOM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ATOMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ATOM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ATOM_USD"
			  }
			]
		  },
		  "AVAX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AVAX",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "AVAXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AVAXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "AVAX-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "AVAX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "avaxusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "AVAXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "AVAX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AVAX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "AVAX_USD"
			  }
			]
		  },
		  "AXL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AXL",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "AXLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AXLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "AXL-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "WAXL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "WAXLUSD"
			  }
			]
		  },
		  "BCH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BCH",
				"Quote": "USD"
			  },
			  "decimals": 7,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "BCHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BCHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BCH-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "BCH_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "bchusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "BCHUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BCH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BCHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BCH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "BCH_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "BLUR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BLUR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BLUR-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "BLUR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "BLURUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BLUR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BLURUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BLUR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "BNB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BNB",
				"Quote": "USD"
			  },
			  "decimals": 7,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BNB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BNB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BNBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "binance_api",
				"off_chain_ticker": "BNBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BNBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "BNB_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "BONK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BONK",
				"Quote": "USD"
			  },
			  "decimals": 14,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "BONKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BONKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BONK-USD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BONK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BONK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BONKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "BONK_USD"
			  }
			]
		  },
		  "BTC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BTC",
				"Quote": "USD"
			  },
			  "decimals": 5,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "BTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BTC-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "btcusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XXBTZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "BTC_USD"
			  }
			]
		  },
		  "COMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COMP",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "COMPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "COMP-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "COMP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "COMPUSD"
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "COMPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "COMP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "CRV/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CRV",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "CRVUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "CRV-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "CRV_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "CRVUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "CRV-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "CRVUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "CRV-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "DOGE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DOGE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "DOGEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DOGEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "DOGE-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DOGE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "dogeusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XDGUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DOGE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DOGEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "DOGE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "DOGE_USD"
			  }
			]
		  },
		  "DOT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DOT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "DOTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DOTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "DOT-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DOT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "DOTUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DOT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DOTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "DOT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "DOT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "DYDX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DYDX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "DYDXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DYDXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DYDX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DYDX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DYDXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "DYDX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "DYM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DYM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "DYMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DYMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DYM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DYM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DYMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "EOS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "EOS",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "EOSUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "EOSUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "EOS-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "EOS_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "EOSUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "EOS-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "EOS-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "EOSUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ETC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETC",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ETCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ETC-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ETC_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "etcusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ETC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ETCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ETC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ETH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETH",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ETHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ETHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ETH-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ethusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XETHZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ETH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ETHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ETH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ETH_USD"
			  }
			]
		  },
		  "FET/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FET",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "FETUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "FET-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "FETUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "FET-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "FET-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "FETUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "FET_USD"
			  }
			]
		  },
		  "FIL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FIL",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "FILUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "FIL-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "FIL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "filusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "FILUSD"
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "FILUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "FIL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "GRT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GRT",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "GRTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "GRTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "GRT-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "GRT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "GRTUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "GRT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "GRTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "GRT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "GRT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "HBAR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HBAR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "HBARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bitstamp_ws",
				"off_chain_ticker": "hbarusd"
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "HBARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "HBAR-USD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "HBAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "HBARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "HBAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "HBAR_USD"
			  }
			]
		  },
		  "ICP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ICP",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ICPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ICPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ICP-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ICPUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ICP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ICP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ICPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ICP_USD"
			  }
			]
		  },
		  "IMX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "IMX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "IMXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "IMX-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "IMXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "IMX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "IMXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "IMX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "INJ/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "INJ",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "INJUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "INJUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "INJ-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "INJUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "INJ-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "INJ-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "INJUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "INJ_USD"
			  }
			]
		  },
		  "JTO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "JTO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "JTO-USD"
			  },
			  {
				"name": "binance_api",
				"off_chain_ticker": "JTOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "JTOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "JTOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "JTO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "JTO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "JUP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "JUP",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "JUP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "JUP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "binance_api",
				"off_chain_ticker": "JUPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "JUPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "JUP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "JUPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "JUP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "LDO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LDO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "LDOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "LDO-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "LDOUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "LDO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "LDOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "LDO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "LDO_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "LINK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LINK",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "LINKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "LINKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "LINK-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "LINKUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "LINK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "LINKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "LINK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "LINK_USD"
			  }
			]
		  },
		  "LTC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LTC",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "LTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "LTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "LTC-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ltcusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XLTCZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "LTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "LTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "LTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "LTC_USD"
			  }
			]
		  },
		  "MANA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MANA",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "MANAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "MANA-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "MANA_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "MANAUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "MANA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "MANAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "MANA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "MATIC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MATIC",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "MATICUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "MATICUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "MATIC-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "MATIC_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "maticusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "MATICUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "MATIC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "MATICUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "MATIC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "MATIC_USD"
			  }
			]
		  },
		  "MKR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MKR",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "MKRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "MKR-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "MKRUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "MKR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "MKRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "MKR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "NEAR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NEAR",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "NEARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "NEAR-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "NEAR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "nearusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "NEAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "NEARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "NEAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "NEAR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "NTRN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NTRN",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 2,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "NTRNUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "NTRN_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "NTRN-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "NTRN-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "OP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "OP",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "OPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "OP-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "OP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "OP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "OPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "OP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "OP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ORDI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ORDI",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "ORDIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ORDIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ORDI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ordiusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ORDI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ORDI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ORDIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "PEPE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PEPE",
				"Quote": "USD"
			  },
			  "decimals": 16,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "PEPEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "PEPEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "PEPE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "PEPEUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "PEPE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "PEPEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "PEPE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "PEPE_USD"
			  }
			]
		  },
		  "PYTH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PYTH",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "PYTHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "PYTHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "PYTH_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "PYTH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "PYTH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "PYTHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "RNDR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RNDR",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "RNDRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "RNDR-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "RNDRUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "RNDR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "RNDR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "RNDRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "RNDR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "RUNE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RUNE",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "RUNEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "RUNE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "RUNEUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "RUNE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "RUNEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SEI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SEI",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "SEIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SEIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SEI-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "SEI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "seiusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SEI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SEIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SHIB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SHIB",
				"Quote": "USD"
			  },
			  "decimals": 15,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "SHIBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SHIBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SHIB-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "SHIB_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SHIBUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SHIB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SHIBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SHIB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "SHIB_USD"
			  }
			]
		  },
		  "SNX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SNX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "SNXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SNXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SNX-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SNXUSD"
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SNXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SNX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SOL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "SOLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SOLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SOL-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "solusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SOLUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SOL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SOLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SOL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "SOL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "STRK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STRK",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "STRKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "STRKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "STRKUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "STRK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "STRK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "STRK_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "STX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "STXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "STXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "STX-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "STX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "STXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "STX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "STX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "STXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "STX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SUI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SUI",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "SUIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SUIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SUI-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "SUI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "suiusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SUI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SUIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SUI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "SUI_USD"
			  }
			]
		  },
		  "TIA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TIA",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "TIAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "TIAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "TIA-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "tiausdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "TIAUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "TIA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "TIAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "TIA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "TIA_USD"
			  }
			]
		  },
		  "TRX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRX",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "TRXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "TRXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "TRX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "trxusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "TRXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "TRX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "TRXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "TRX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "UNI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "UNI",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "UNIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "UNIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "UNI-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "UNI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "UNIUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "UNI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "UNI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "UNI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "USDT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "USDCUSDT",
				"invert": true
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "USDCUSDT",
				"invert": true
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "USDT-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ethusdt",
				"normalize_by_pair": {
				  "Base": "ETH",
				  "Quote": "USD"
				},
				"invert": true
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "USDTZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BTC-USDT",
				"normalize_by_pair": {
				  "Base": "BTC",
				  "Quote": "USD"
				},
				"invert": true
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "USDC-USDT",
				"invert": true
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "USDT_USD"
			  }
			]
		  },
		  "WLD/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WLD",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "WLDUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "WLDUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "WLD_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "wldusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "WLD-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "WLDUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "WLD-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "WLD_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "WOO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WOO",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "WOOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "WOO_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "WOO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "WOO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "WOOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "XLM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "XLM",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "XLMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "XLMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "XLM-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XXLMZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "XLM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "XLMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "XLM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "XLM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "XRP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "XRP",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_api",
				"off_chain_ticker": "XRPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "XRPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "XRP-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "XRP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "xrpusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XXRPZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "XRP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "XRPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "XRP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "XRP_USD"
			  }
			]
		  }
		}
	}
	  `

	// UniswapV3BaseMarketMap is used to initialize the Base market map. This only includes
	// the markets that are supported by uniswapv3 on Base.
	UniswapV3BaseMarketMap mmtypes.MarketMap

	// UniswapV3BaseMarketMapJSON is the JSON representation of UniswapV3BaseMarketMap.
	UniswapV3BaseMarketMapJSON = `
	{
		"markets": {
		  "BRETT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BRETT",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "BRETT/ETH",
				"metadata_JSON": "{\"address\":\"0xBA3F945812a83471d709BCe9C3CA699A19FB46f7\",\"base_decimals\":18,\"quote_decimals\":18,\"invert\":true}",
				"normalize_by_pair": {
					"Base": "ETH",
					"Quote": "USD"
				}
			  }
			]
		  },
		  "ETH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETH",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "ETH/USDT",
				"metadata_JSON": "{\"address\":\"0xd92E0767473D1E3FF11Ac036f2b1DB90aD0aE55F\",\"base_decimals\":18,\"quote_decimals\":6,\"invert\":false}",
				"normalize_by_pair": {
					"Base": "USDT",
					"Quote": "USD"
				}
			  }
			]
		  },
		  "USDT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDT",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "USDT/ETH",
				"metadata_JSON": "{\"address\":\"0xd92E0767473D1E3FF11Ac036f2b1DB90aD0aE55F\",\"base_decimals\":6,\"quote_decimals\":18,\"invert\":true}",
				"normalize_by_pair": {
				  "Base": "ETH",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "USDT/USDC",
				"metadata_JSON": "{\"address\":\"0xD56da2B74bA826f19015E6B7Dd9Dae1903E85DA1\",\"base_decimals\":6,\"quote_decimals\":6,\"invert\":true}"
			  }
			]
		  }
	    }
	}
	`
)

func init() {
	// Unmarshal the RaydiumMarketMapJSON into RaydiumMarketMap.
	if err := json.Unmarshal([]byte(RaydiumMarketMapJSON), &RaydiumMarketMap); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal RaydiumMarketMapJSON: %v\n", err)
		panic(err)
	}

	if err := RaydiumMarketMap.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to validate RaydiumMarketMap: %v\n", err)
		panic(err)
	}

	// Unmarshal the CoreMarketMapJSON into CoreMarketMap.
	if err := json.Unmarshal([]byte(CoreMarketMapJSON), &CoreMarketMap); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal CoreMarketMapJSON: %v\n", err)
		panic(err)
	}

	if err := CoreMarketMap.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to validate CoreMarketMap: %v\n", err)
		panic(err)
	}

	// Unmarshal the UniswapV3BaseMarketMapJSON into UniswapV3BaseMarketMap.
	if err := json.Unmarshal([]byte(UniswapV3BaseMarketMapJSON), &UniswapV3BaseMarketMap); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal UniswapV3BaseMarketMapJSON: %v\n", err)
		panic(err)
	}

	if err := UniswapV3BaseMarketMap.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to validate UniswapV3BaseMarketMap: %v\n", err)
		panic(err)
	}
}
