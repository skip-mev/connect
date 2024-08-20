package tickermetadata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/marketmap/types/tickermetadata"
)

func Test_UnmarshalCoreMetadata(t *testing.T) {
	t.Run("can marshal and unmarshal the same struct and values", func(t *testing.T) {
		elem := tickermetadata.NewCoreMetadata(
			[]tickermetadata.AggregatorID{
				tickermetadata.NewAggregatorID("coingecko", "id"),
				tickermetadata.NewAggregatorID("cmc", "id"),
			},
		)

		bz, err := tickermetadata.MarshalCoreMetadata(elem)
		require.NoError(t, err)

		elem2, err := tickermetadata.CoreMetadataFromJSONBytes(bz)
		require.NoError(t, err)
		require.Equal(t, elem, elem2)
	})

	t.Run("can marshal and unmarshal the same struct and values with empty AggregatorIDs", func(t *testing.T) {
		elem := tickermetadata.NewCoreMetadata(nil)

		bz, err := tickermetadata.MarshalCoreMetadata(elem)
		require.NoError(t, err)

		elem2, err := tickermetadata.CoreMetadataFromJSONBytes(bz)
		require.NoError(t, err)
		require.Equal(t, elem, elem2)
	})

	t.Run("can unmarshal a JSON string into a struct", func(t *testing.T) {
		elemJSON := `{"aggregate_ids":[{"venue":"coingecko","ID":"id"},{"venue":"cmc","ID":"id"}]}`
		elem, err := tickermetadata.CoreMetadataFromJSONString(elemJSON)
		require.NoError(t, err)

		require.Equal(t, tickermetadata.NewCoreMetadata(
			[]tickermetadata.AggregatorID{
				tickermetadata.NewAggregatorID("coingecko", "id"),
				tickermetadata.NewAggregatorID("cmc", "id"),
			},
		), elem)
	})
}
