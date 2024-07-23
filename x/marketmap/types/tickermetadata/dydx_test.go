package tickermetadata_test

import (
	"testing"

	"github.com/skip-mev/slinky/x/marketmap/types/tickermetadata"
	"github.com/stretchr/testify/require"
)

func Test_Unmarshal(t *testing.T) {
	t.Run("can marshal and unmarshal the same struct and values", func(t *testing.T) {
		elem := tickermetadata.NewDyDx(
			100,
			1000,
			[]tickermetadata.AggregatorID{
				tickermetadata.NewAggregatorID("coingecko", "id"),
				tickermetadata.NewAggregatorID("cmc", "id"),
			},
		)

		bz, err := tickermetadata.MarshalDyDx(elem)
		require.NoError(t, err)

		elem2, err := tickermetadata.DyDxFromJSONBytes(bz)
		require.NoError(t, err)
		require.Equal(t, elem, elem2)
	})

	t.Run("can marshal and unmarshal the same struct and values with empty AggregatorIDs", func(t *testing.T) {
		elem := tickermetadata.NewDyDx(100, 1000, nil)

		bz, err := tickermetadata.MarshalDyDx(elem)
		require.NoError(t, err)

		elem2, err := tickermetadata.DyDxFromJSONBytes(bz)
		require.NoError(t, err)
		require.Equal(t, elem, elem2)
	})

	t.Run("can unmarshal a JSON string into a struct", func(t *testing.T) {
		elemJSON := `{"reference_price":100,"liquidity":1000,"aggregate_ids":[{"venue":"coingecko","ID":"id"},{"venue":"cmc","ID":"id"}]}`
		elem, err := tickermetadata.DyDxFromJSONString(elemJSON)
		require.NoError(t, err)

		require.Equal(t, tickermetadata.NewDyDx(
			100,
			1000,
			[]tickermetadata.AggregatorID{
				tickermetadata.NewAggregatorID("coingecko", "id"),
				tickermetadata.NewAggregatorID("cmc", "id"),
			},
		), elem)
	})
}
