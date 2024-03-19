package aggregator

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/maps"
)

type (
	// AggregatedProviderData defines a type alias for a map
	// of provider -> data (i.e. a set of prices).
	AggregatedProviderData[K comparable, V any] map[K]V

	// AggregateFn is the function used to aggregate data from each provider. Given a
	// map of provider -> values, the aggregate function should return a final
	// value.
	AggregateFn[K comparable, V any] func(providers AggregatedProviderData[K, V]) V

	// AggregateFnFromContext is a function that is used to parametrize an aggregateFn
	// by an sdk.Context. This is used to allow the aggregateFn to access the latest state
	// of an application i.e computing a stake weighted median based on the latest validator set.
	AggregateFnFromContext[K comparable, V any] func(ctx sdk.Context) AggregateFn[K, V]
)

// Aggregator defines the expected interface that must be implemented by any custom data aggregator.
type Aggregator[K comparable, V any] interface {
	GetProviderData() AggregatedProviderData[K, V]
	GetDataByProvider(provider K) V
	SetProviderData(provider K, data V)
	ResetProviderData()
	AggregateData()
	AggregateDataFromContext(ctx sdk.Context)
	GetAggregatedData() V
	SetAggregatedData(aggregatedData V)
}

// DataAggregator is a simple aggregator for provider data. It is thread-safe since
// it is assumed to be called concurrently in data fetching goroutines. The DataAggregator
// requires one of either an aggregateFn or aggregateFnFromContext to be set.
type DataAggregator[K comparable, V any] struct {
	sync.Mutex

	// aggregateFn is the function used to aggregate data from each provider.
	aggregateFn AggregateFn[K, V]

	// aggregateFnFromContext is a function that is used to parametrize an aggregateFn
	// by an sdk.Context.
	aggregateFnFromContext AggregateFnFromContext[K, V]

	// providerData is a map of provider -> value (i.e. prices).
	providerData AggregatedProviderData[K, V]

	// aggregatedData is the current set of aggregated data across the providers.
	aggregatedData V
}

// NewDataAggregator returns a DataAggregator. The DataAggregator
// is responsible for aggregating data (such as prices) from each provider
// and computing the final aggregated data (final price). The DataAggregator is
// thread-safe since it is assumed to be called concurrently in price
// fetching goroutines.
func NewDataAggregator[K comparable, V any](opts ...DataAggregatorOption[K, V]) *DataAggregator[K, V] {
	agg := &DataAggregator[K, V]{
		providerData:   make(AggregatedProviderData[K, V]),
		aggregatedData: *new(V),
	}

	for _, opt := range opts {
		opt(agg)
	}

	return agg
}

// GetProviderData returns a copy of the aggregated provider data.
func (p *DataAggregator[K, V]) GetProviderData() AggregatedProviderData[K, V] {
	p.Lock()
	defer p.Unlock()

	cpy := make(AggregatedProviderData[K, V])
	maps.Copy(cpy, p.providerData)

	return cpy
}

// GetDataByProvider returns the data currently stored for a given provider.
func (p *DataAggregator[K, V]) GetDataByProvider(provider K) V {
	p.Lock()
	defer p.Unlock()

	cpy := make(AggregatedProviderData[K, V])
	maps.Copy(cpy, p.providerData)

	return cpy[provider]
}

// SetProviderData updates the data aggregator with the given provider
// and data.
func (p *DataAggregator[K, V]) SetProviderData(provider K, data V) {
	p.Lock()
	defer p.Unlock()

	p.providerData[provider] = data
}

// ResetProviderData resets the data aggregator for all providers.
func (p *DataAggregator[K, V]) ResetProviderData() {
	p.Lock()
	defer p.Unlock()

	p.providerData = make(AggregatedProviderData[K, V])
}

// AggregateData aggregates the current set of data by using the aggregate function.
func (p *DataAggregator[K, V]) AggregateData() {
	if p.aggregateFn == nil {
		panic("aggregateFn cannot be nil")
	}

	providerData := p.GetProviderData()
	p.SetAggregatedData(p.aggregateFn(providerData))
}

// AggregateDataFromContext aggregates the current set of data by using the aggregate function
// parametrized by the given context.
func (p *DataAggregator[K, V]) AggregateDataFromContext(ctx sdk.Context) {
	if p.aggregateFnFromContext == nil {
		panic("aggregateFnFromContext cannot be nil")
	}

	aggregateFn := p.aggregateFnFromContext(ctx)
	providerData := p.GetProviderData()
	p.SetAggregatedData(aggregateFn(providerData))
}

// GetAggregatedData returns the aggregated data based on the provided data.
func (p *DataAggregator[K, V]) GetAggregatedData() V {
	p.Lock()
	defer p.Unlock()

	return p.aggregatedData
}

// SetAggregatedData sets the current set of aggregated data.
func (p *DataAggregator[K, V]) SetAggregatedData(aggregatedData V) {
	p.Lock()
	defer p.Unlock()

	p.aggregatedData = aggregatedData
}
