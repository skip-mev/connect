#!/usr/bin/env bash
set -e

NUM_MARKETS=$(echo "$MARKETS" | jq '.markets | length + 1')

jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
NUM_MARKETS=$NUM_MARKETS; jq --arg num "$NUM_MARKETS" '.app_state["oracle"]["next_id"] = $num' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["marketmap"]["market_map"] = ($markets | fromjson)' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] += [$markets | fromjson | .markets | values | .[].ticker.currency_pair | {"currency_pair": {"Base": .Base, "Quote": .Quote}, "currency_pair_price": null, "nonce": 0} ]' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] |= (to_entries | map(.value += {id: (.key + 1)} | .value))' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
