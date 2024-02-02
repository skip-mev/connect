package types

import (
	"fmt"
	"time"
)

// ResponseKey is a type restriction interface for the key of a GetResponse.
type ResponseKey interface {
	comparable
	fmt.Stringer
}

// ResponseValue is a type restriction interface for the value of a GetResponse.
type ResponseValue interface {
	fmt.Stringer
}

// GetResponse is the GET response from the API data handler.
type GetResponse[K ResponseKey, V ResponseValue] struct {
	Resolved   map[K]Result[V]
	UnResolved map[K]error
}

// Result is the result of a single requested Ticker.
type Result[V ResponseValue] struct {
	// Value is the value of the requested Ticker.
	Value V
	// Timestamp is the timestamp of the value.
	Timestamp time.Time
}

// NewGetResponse creates a new GetResponse.
func NewGetResponse[K ResponseKey, V ResponseValue](resolved map[K]Result[V], unresolved map[K]error) GetResponse[K, V] {
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
func NewGetResponseWithErr[K ResponseKey, V ResponseValue](ids []K, err error) GetResponse[K, V] {
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
func NewResult[V ResponseValue](value V, timestamp time.Time) Result[V] {
	return Result[V]{
		Value:     value,
		Timestamp: timestamp,
	}
}

// String returns a string representation of the Result. This is mostly used for logging
// and testing purposes.
func (r Result[V]) String() string {
	return fmt.Sprintf(
		"(value: %s, timestamp: %s)",
		r.Value.String(),
		r.Timestamp.String(),
	)
}
