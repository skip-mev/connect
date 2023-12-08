package aggregator

// DataAggregatorOption is a function that is used to parametrize a DataAggregator
// instance.
type DataAggregatorOption[K comparable, V any] func(*DataAggregator[K, V])

// WithAggregateFn sets the aggregateFn of a DataAggregatorOptions.
func WithAggregateFn[K comparable, V any](fn AggregateFn[K, V]) DataAggregatorOption[K, V] {
	return func(opts *DataAggregator[K, V]) {
		opts.aggregateFn = fn
	}
}

// WithAggregateFnFromContext sets the aggregateFnFromContext of a DataAggregatorOptions.
func WithAggregateFnFromContext[K comparable, V any](fn AggregateFnFromContext[K, V]) DataAggregatorOption[K, V] {
	return func(opts *DataAggregator[K, V]) {
		opts.aggregateFnFromContext = fn
	}
}
