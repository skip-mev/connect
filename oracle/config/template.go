package config

const (
	DefaultConfigTemplate = `

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
## Update Interval (in seconds) is the time between each time the oracle triggers providers to update price-data
update_interval = "1s"

## Timeout is the time that the vote-extension handler will wait for a response from the oracle (either running in / out-of process), generally this parameter should be 
## less than the timeout_prevote parameter in the consensus config
timeout = "3s"

## InProcess specifies whether the oracle configured, is currently running as a remote grpc-server, or will be run in process
in_process = true

## RemoveAddress is the address of the remote oracle grpc-server, only used if in_process is set to false
remote_address="localhost:8080"

# Providers
[[oracle.providers]]
## Name of provier to query price-data from
name = "coingecko"

[[oracle.providers]]
name = "coinmarketcap"
[oracle.providers.token_name_to_symbol]
BITCOIN = "BTC"
ETHEREUM = "ETH"
COSMOS = "ATOM"
USD = "USD"
OSMOSIS = "OSMO"


# Currency Pairs
[[oracle.currency_pairs]]
base = "BITCOIN"
quote = "USD"
quote_decimals = 8
`
)
