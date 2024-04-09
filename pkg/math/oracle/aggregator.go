package oracle

import (
	"maps"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/oracle/types"
)

// GetProviderData returns the provider data the aggregator has.
func (m *MedianAggregator) GetProviderData() types.AggregatedProviderPrices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cpy := make(types.AggregatedProviderPrices)
	maps.Copy(cpy, m.providerPrices)

	return cpy
}

// GetDataByProvider returns the data currently stored for a given provider.
func (m *MedianAggregator) GetDataByProvider(provider string) types.AggregatorPrices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cpy := make(types.AggregatorPrices)
	maps.Copy(cpy, m.providerPrices[provider])

	return cpy
}

// SetProviderData updates the data aggregator with the given provider and data.
func (m *MedianAggregator) SetProviderData(provider string, data types.AggregatorPrices) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if data == nil {
		data = make(types.AggregatorPrices)
	}

	m.providerPrices[provider] = data
}

// ResetProviderData resets the data aggregator for all providers.
func (m *MedianAggregator) ResetProviderData() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.providerPrices = make(types.AggregatedProviderPrices)
}

// GetAggregatedData returns the aggregated data the aggregator has.
func (m *MedianAggregator) GetAggregatedData() types.AggregatorPrices {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cpy := make(types.AggregatorPrices)
	maps.Copy(cpy, m.scaledPrices)

	return cpy
}

// SetAggregatedData updates the data aggregator with the given aggregated data.
func (m *MedianAggregator) SetAggregatedData(data types.AggregatorPrices) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if data == nil {
		data = make(types.AggregatorPrices)
	}

	m.scaledPrices = data
}

// AggregateDataFromContext is a no-op for the median aggregator.
func (m *MedianAggregator) AggregateDataFromContext(_ sdk.Context) {
	// no-op
}
