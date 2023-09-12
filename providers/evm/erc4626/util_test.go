package erc4626

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitValueFromDecimals(t *testing.T) {
	result := getUnitValueFromDecimals(18)
	require.Equal(t, big.NewInt(1000000000000000000), result)

	result = getUnitValueFromDecimals(6)
	require.Equal(t, big.NewInt(1000000), result)

	result = getUnitValueFromDecimals(0)
	require.Equal(t, big.NewInt(1), result)
}
