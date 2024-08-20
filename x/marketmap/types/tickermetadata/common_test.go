package tickermetadata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/marketmap/types/tickermetadata"
)

func Test_UnmarshalAggregatorID(t *testing.T) {
	t.Run("can marshal and unmarshal the same struct and values", func(t *testing.T) {
		elem := tickermetadata.NewAggregatorID("coingecko", "id")

		bz, err := tickermetadata.MarshalAggregatorID(elem)
		require.NoError(t, err)

		elem2, err := tickermetadata.AggregatorIDFromJSONBytes(bz)
		require.NoError(t, err)
		require.Equal(t, elem, elem2)
	})

	t.Run("can unmarshal a JSON string into a struct", func(t *testing.T) {
		elemJSON := `{"venue":"coingecko","ID":"id"}`
		elem, err := tickermetadata.AggregatorIDFromJSONString(elemJSON)
		require.NoError(t, err)

		require.Equal(t, tickermetadata.NewAggregatorID("coingecko", "id"), elem)
	})
}
