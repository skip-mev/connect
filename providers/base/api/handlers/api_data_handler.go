package handlers

import (
	"net/http"

	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// APIDataHandler defines an interface that must be implemented by all providers that
// want to fetch data from an API using HTTP requests. This interface is meant to be
// paired with the APIQueryHandler. The APIQueryHandler will use the APIDataHandler to
// create the URL to be sent to the HTTP client and parse the response from the client.
//
//go:generate mockery --name APIDataHandler --output ./mocks/ --case underscore
type APIDataHandler[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
	// CreateURL is used to create the URL to be sent to the http client. The function
	// should utilize the IDs passed in as references to the data that needs to be fetched.
	CreateURL(ids []K) (string, error)

	// ParseResponse is used to parse the response from the client. The response should be
	// parsed into a map of IDs to results. If any IDs are not resolved, they should
	// be returned in the unresolved map. The timestamp associated with the result should
	// reflect either the time the data was fetched or the time the API last updated the data.
	ParseResponse(ids []K, response *http.Response) providertypes.GetResponse[K, V]
}
