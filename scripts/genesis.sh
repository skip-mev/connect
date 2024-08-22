#!/usr/bin/env bash
set -eux

go run $SCRIPT_DIR/genesis.go --use-core=$USE_CORE_MARKETS --use-raydium=$USE_RAYDIUM_MARKETS \
    --use-uniswapv3-base=$USE_UNISWAPV3_BASE_MARKETS --use-coingecko=$USE_COINGECKO_MARKETS \
    --use-polymarket=$USE_POLYMARKET_MARKETS --use-coinmarketcap=$USE_COINMARKETCAP_MARKETS \
    --use-osmosis=$USE_OSMOSIS_MARKETS --temp-file=markets.json
MARKETS=$(cat markets.json)

echo "MARKETS content: $MARKETS"

NUM_MARKETS=$(echo "$MARKETS" | jq '.markets | length + 1')

./build/connectd init validator --chain-id skip-1 --home "$HOMEDIR"
./build/connectd keys add validator --home "$HOMEDIR" --keyring-backend test
./build/connectd genesis add-genesis-account validator 10000000000000000000000000stake --home "$HOMEDIR" --keyring-backend test
./build/connectd genesis add-genesis-account cosmos1see0htr47uapjvcvh0hu6385rp8lw3em24hysg 10000000000000000000000000stake --home "$HOMEDIR" --keyring-backend test
./build/connectd genesis gentx validator 1000000000stake --chain-id skip-1 --home "$HOMEDIR" --keyring-backend test
./build/connectd genesis collect-gentxs --home "$HOMEDIR"

jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
NUM_MARKETS=$NUM_MARKETS; jq --arg num "$NUM_MARKETS" '.app_state["oracle"]["next_id"] = $num' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["marketmap"]["market_map"] = ($markets | fromjson)' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] += [$markets | fromjson | .markets | values | .[].ticker.currency_pair | {"currency_pair": {"Base": .Base, "Quote": .Quote}, "currency_pair_price": null, "nonce": 0} ]' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] |= (to_entries | map(.value += {id: (.key + 1)} | .value))' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"


# Adding all genesis accounts to the genesis file, look for all genesis accounts to be in env variables
# GENESIS_ACCOUNT_1, GENESIS_ACCOUNT_2, GENESIS_ACCOUNT_3, etc, and their respective balances
# GENESIS_ACCOUNT_1_BALANCE, GENESIS_ACCOUNT_2_BALANCE, GENESIS_ACCOUNT_3_BALANCE, etc
for i in $(env | grep -v "_BALANCE" | grep -o "GENESIS_ACCOUNT_[0-9]*" | sort -n -t _ -k 3); do
    # get the account address
    ACCOUNT=$(printenv $i)

    # get the account balance
    BALANCE=$(printenv ${i}_BALANCE)

    # add the account to the genesis file
    ./build/connectd genesis add-genesis-account $ACCOUNT $BALANCE --home "$HOMEDIR" --keyring-backend test
done

# Check if MARKET_MAP_AUTHORITY environment variable exists and is not empty
if [ -n "${MARKET_MAP_AUTHORITY+x}" ] && [ -n "$MARKET_MAP_AUTHORITY" ]; then
    MARKET_MAP_AUTHORITY=$(printenv MARKET_MAP_AUTHORITY)
    jq --arg authority "$MARKET_MAP_AUTHORITY" \
    '.app_state["marketmap"]["params"]["market_authorities"] += [$authority]' \
    "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
fi

rm markets.json
