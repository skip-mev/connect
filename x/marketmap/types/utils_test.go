package types_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestReadMarketMapFromFile(t *testing.T) {
	t.Run("invalid file does not exist", func(t *testing.T) {
		path := filepath.Join("testdata", "invalid.txt")
		_, err := types.ReadMarketMapFromFile(path)
		require.Error(t, err)
	})

	t.Run("invalid datatype", func(t *testing.T) {
		path := filepath.Join("testdata", "markets.json")
		_, err := types.ReadMarketMapFromFile(path)
		require.Error(t, err)
	})

	t.Run("valid file exists", func(t *testing.T) {
		path := filepath.Join("testdata", "marketmap.json")
		_, err := types.ReadMarketMapFromFile(path)
		require.NoError(t, err)
	})
}

func TestReadMarketFromFile(t *testing.T) {
	t.Run("invalid file does not exist", func(t *testing.T) {
		path := filepath.Join("testdata", "invalid.txt")
		_, err := types.ReadMarketsFromFile(path)
		require.Error(t, err)
	})

	t.Run("invalid datatype", func(t *testing.T) {
		path := filepath.Join("testdata", "marketmap.json")
		_, err := types.ReadMarketsFromFile(path)
		require.Error(t, err)
	})

	t.Run("valid file exists", func(t *testing.T) {
		path := filepath.Join("testdata", "markets.json")
		_, err := types.ReadMarketsFromFile(path)
		require.NoError(t, err)
	})
}

func TestToMarketMap(t *testing.T) {
	path := filepath.Join("testdata", "markets.json")
	ms, err := types.ReadMarketsFromFile(path)
	require.NoError(t, err)

	path = filepath.Join("testdata", "marketmap.json")
	mMap, err := types.ReadMarketMapFromFile(path)
	require.NoError(t, err)

	require.Equal(t, mMap, ms.ToMarketMap())
}
