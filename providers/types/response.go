package types

import (
	"fmt"
	"time"
)

// ResponseCode is an optional code that can be attached to responses to provide
// additional context.
type ResponseCode int

const (
	// ResponseCodeOK is a code that notifies the base provider that the response
	// is OK.
	ResponseCodeOK ResponseCode = 0
	// ResponseCodeUnchange is a code that notifies the base provider that the response
	// is unchanged for the given ID. This is useful when the provider has a cache
	// and the value has not changed.
	ResponseCodeUnchanged ResponseCode = 1
)

func (r ResponseCode) String() string {
	switch r {
	case ResponseCodeOK:
		return "ok"
	case ResponseCodeUnchanged:
		return "unchanged"
	default:
		return "unknown"
	}
}

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
	Resolved   map[K]ResolvedResult[V]
	UnResolved map[K]UnresolvedResult
}

// ResolvedResult is the result of a single requested ID.
type ResolvedResult[V ResponseValue] struct {
	// Value is the value of the requested ID.
	Value V
	// Timestamp is the timestamp of the value.
	Timestamp time.Time
	// ResponseCode is an optional code that can be attached to responses to provide
	// additional context.
	ResponseCode ResponseCode
}

// UnresolvedResult is an unresolved (failed) result of a single requested ID.
type UnresolvedResult struct {
	ErrorWithCode
}

// NewGetResponse creates a new GetResponse.
func NewGetResponse[K ResponseKey, V ResponseValue](resolved map[K]ResolvedResult[V], unresolved map[K]UnresolvedResult) GetResponse[K, V] {
	if resolved == nil {
		resolved = make(map[K]ResolvedResult[V])
	}

	if unresolved == nil {
		unresolved = make(map[K]UnresolvedResult)
	}

	return GetResponse[K, V]{
		Resolved:   resolved,
		UnResolved: unresolved,
	}
}

// NewGetResponseWithErr creates a new GetResponse with the given error. This populates
// the unresolved map with the given IDs and error.
func NewGetResponseWithErr[K ResponseKey, V ResponseValue](ids []K, err ErrorWithCode) GetResponse[K, V] {
	unresolved := make(map[K]UnresolvedResult, len(ids))
	for _, id := range ids {
		unresolved[id] = UnresolvedResult{
			err,
		}
	}

	return GetResponse[K, V]{
		Resolved:   make(map[K]ResolvedResult[V]),
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

// NewResult creates a new ResolvedResult.
func NewResult[V ResponseValue](value V, timestamp time.Time) ResolvedResult[V] {
	return ResolvedResult[V]{
		Value:     value,
		Timestamp: timestamp,
	}
}

// NewResultWithCode creates a new ResolvedResult with the given error code.
func NewResultWithCode[V ResponseValue](value V, timestamp time.Time, code ResponseCode) ResolvedResult[V] {
	return ResolvedResult[V]{
		Value:        value,
		Timestamp:    timestamp,
		ResponseCode: code,
	}
}

// String returns a string representation of the ResolvedResult. This is mostly used for logging
// and testing purposes.
func (r ResolvedResult[V]) String() string {
	return fmt.Sprintf(
		"(value: %s, timestamp: %s, response code: %s)",
		r.Value.String(),
		r.Timestamp.String(),
		r.ResponseCode.String(),
	)
}
