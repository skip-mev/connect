package erc4626

import (
	"testing"

	"github.com/skip-mev/slinky/oracle/types"
	moduletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	pairs := []moduletypes.CurrencyPair{
		{
			Base:  "eth",
			Quote: "ryeth",
		},
		{
			Base:  "usdc",
			Quote: "weth",
		},
	}
	testTokenNameToMetadata := map[string]types.TokenMetadata{
		"ryeth": {
			Symbol:   "0x0000000000000000000000000000000000000000",
			Decimals: 18,
		},
		"usdc": {
			Symbol:   "0x0000000000000000000000000000000000000000",
			Decimals: 6,
		},
	}

	// create a new provider
	_, err := NewProvider(nil, pairs, "", testTokenNameToMetadata)
	require.NoError(t, err)
}
