#!/usr/bin/env bash
set -eux

go run $SCRIPT_DIR/genesis.go --use-core=$USE_CORE_MARKETS --use-raydium=$USE_RAYDIUM_MARKETS \
--use-uniswapv3-base=$USE_UNISWAPV3_BASE_MARKETS --use-coingecko=$USE_COINGECKO_MARKETS \
--use-coinmarketcap=$USE_COINMARKETCAP_MARKETS --use-osmosis=$USE_OSMOSIS_MARKETS --temp-file=markets.json
MARKETS=$(cat markets.json)

echo "MARKETS content: $MARKETS"

NUM_MARKETS=$(echo "$MARKETS" | jq '.markets | length + 1')

./build/slinkyd init validator --chain-id skip-1 --home "$HOMEDIR"
./build/slinkyd keys add validator --home "$HOMEDIR" --keyring-backend test
./build/slinkyd genesis add-genesis-account validator 10000000000000000000000000stake --home "$HOMEDIR" --keyring-backend test
./build/slinkyd genesis add-genesis-account cosmos1see0htr47uapjvcvh0hu6385rp8lw3em24hysg 10000000000000000000000000stake --home "$HOMEDIR" --keyring-backend test
./build/slinkyd genesis gentx validator 1000000000stake --chain-id skip-1 --home "$HOMEDIR" --keyring-backend test
./build/slinkyd genesis collect-gentxs --home "$HOMEDIR"

jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
NUM_MARKETS=$NUM_MARKETS; jq --arg num "$NUM_MARKETS" '.app_state["oracle"]["next_id"] = $num' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["marketmap"]["market_map"] = ($markets | fromjson)' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] += [$markets | fromjson | .markets | values | .[].ticker.currency_pair | {"currency_pair": {"Base": .Base, "Quote": .Quote}, "currency_pair_price": null, "nonce": 0} ]' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] |= (to_entries | map(.value += {id: (.key + 1)} | .value))' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"

rm markets.json