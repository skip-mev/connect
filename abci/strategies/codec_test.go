package strategies_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/abci/strategies"
	vetypes "github.com/skip-mev/slinky/abci/ve/types"
)

func TestDefaultVoteExtensionCodec(t *testing.T) {
	t.Run("test encoding / decoding", func(t *testing.T) {
		// create a sample vote extension
		ve := vetypes.OracleVoteExtension{
			Prices: map[uint64][]byte{
				1: []byte("1"),
				2: []byte("2"),
			},
		}
		// encode it
		codec := strategies.NewDefaultVoteExtensionCodec()
		bz, err := codec.Encode(ve)
		require.NoError(t, err)

		// decode it
		decodedVe, err := codec.Decode(bz)
		require.NoError(t, err)

		// make sure it's the same
		require.Equal(t, ve.Prices, decodedVe.Prices)
	})

	t.Run("test decoding empty byte array", func(t *testing.T) {
		codec := strategies.NewDefaultVoteExtensionCodec()
		_, err := codec.Decode([]byte{})
		require.Nil(t, err)
	})
}

func TestCompressionVoteExtensionCodec(t *testing.T) {
	t.Run("test encoding / decoding", func(t *testing.T) {
		// create a sample vote extension
		samplePrice := []byte("nocapongodskiptoonicewititshiiiiiiiii")
		ve := vetypes.OracleVoteExtension{
			Prices: make(map[uint64][]byte),
		}

		// add 200 prices
		for i := uint64(0); i < 200; i++ {
			ve.Prices[i] = samplePrice
		}

		// create a codec
		defaultCodec := strategies.NewDefaultVoteExtensionCodec()
		codec := strategies.NewCompressionVoteExtensionCodec(defaultCodec, strategies.NewZLibCompressor())

		// encode it
		bz, err := codec.Encode(ve)
		require.NoError(t, err)

		defaultBz, err := defaultCodec.Encode(ve)
		require.NoError(t, err)

		// make sure it's smaller
		require.True(t, len(bz) < len(defaultBz))

		// decode it
		decodedVe, err := codec.Decode(bz)
		require.NoError(t, err)

		// make sure it's the same
		require.Equal(t, ve.Prices, decodedVe.Prices)
	})

	t.Run("test decoding empty byte array", func(t *testing.T) {
		codec := strategies.NewCompressionVoteExtensionCodec(strategies.NewDefaultVoteExtensionCodec(), strategies.NewZLibCompressor())
		_, err := codec.Decode([]byte{})
		require.Nil(t, err)
	})
}
