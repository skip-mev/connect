package types

import (
	"fmt"
	"time"
)

type GetResult interface {
	fmt.Stringer
}

// GetResponse is the GET response from the API data handler.
type GetResponse[K comparable, V GetResult] struct {
	Resolved   map[K]Result[V]
	UnResolved map[K]error
}

// Result is the result of a single requested ID.
type Result[V GetResult] struct {
	// Value is the value of the requested ID.
	Value V
	// Timestamp is the timestamp of the value.
	Timestamp time.Time
}

// NewGetResponse creates a new GetResponse.
func NewGetResponse[K comparable, V GetResult](resolved map[K]Result[V], unresolved map[K]error) GetResponse[K, V] {
	if resolved == nil {
		resolved = make(map[K]Result[V])
	}

	if unresolved == nil {
		unresolved = make(map[K]error)
	}

	return GetResponse[K, V]{
		Resolved:   resolved,
		UnResolved: unresolved,
	}
}

// NewGetResponseWithErr creates a new GetResponse with the given error. This populates
// the unresolved map with the given IDs and error.
func NewGetResponseWithErr[K comparable, V GetResult](ids []K, err error) GetResponse[K, V] {
	unresolved := make(map[K]error, len(ids))
	for _, id := range ids {
		unresolved[id] = err
	}

	return GetResponse[K, V]{
		Resolved:   make(map[K]Result[V]),
		UnResolved: unresolved,
	}
}

// String returns a string representation of the GetResponse. This is mostly used for logging
// and testing purposes.
func (r GetResponse[K, V]) String() string {
	return fmt.Sprintf(
		"resolved: %v | unresolved: %v",
		r.Resolved,
		r.UnResolved,
	)
}

// NewResult creates a new Result.
func NewResult[V GetResult](value V, timestamp time.Time) Result[V] {
	return Result[V]{
		Value:     value,
		Timestamp: timestamp,
	}
}

// String returns a string representation of the Result. This is mostly used for logging
// and testing purposes.
func (r Result[V]) String() string {
	return fmt.Sprintf(
		"(value: %v, timestamp: %v)",
		r.Value,
		r.Timestamp,
	)
}
