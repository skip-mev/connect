package http

import (
	"bytes"
	"io"
	"net/http"
	urlpkg "net/url"
	"sync"
)

// RoundTripper is a wrapper over the standard RoundTripper that permits for simple mock generation
//
//go:generate mockery --name RoundTripper --output ./mocks --outpkg mocks --case underscore
type RoundTripper interface {
	http.RoundTripper
}

// Filter is a function that can be used to filter http.Response objects, it is intended that the only
// Responses given to the Filter are valid (i.e non-error http).
type Filter func([]*http.Response) (*http.Response, error)

type RoundTripperWithURL struct {
	// url is the URL that this round-tripper will mutate it's http.Request with
	url *urlpkg.URL

	// roundTripper is the RoundTripper that this round-tripper will delegate to
	roundTripper RoundTripper
}

// RoundTrip handles the round-trip in accordance with the underlying RoundTripper.
// However, before the request is sent, the URL of the request is mutated to the URL of the RoundTripperWithURL.
func (rt *RoundTripperWithURL) RoundTrip(req *http.Request) (*http.Response, error) {
	// mutate the URL of the request
	req.URL = rt.url

	// delegate to the underlying round-tripper
	return rt.roundTripper.RoundTrip(req)
}

// URL returns the URL of the RoundTripperWithURL.
func (rt *RoundTripperWithURL) URL() *urlpkg.URL {
	return rt.url
}

// NewRoundTripperWithURL returns a new RoundTripperWithURL.
func NewRoundTripperWithURL(urlString string, rt RoundTripper) (*RoundTripperWithURL, error) {
	url, err := urlpkg.Parse(urlString)
	if err != nil {
		return nil, err
	}

	return &RoundTripperWithURL{
		url:          url,
		roundTripper: rt,
	}, nil
}

// MultiRoundTripper is an RoundTripper that delegates the transport of a request to multiple RoundTrippers.
// and applies a filter over all non-error responses.
type MultiRoundTripper struct {
	// roundTrippers is the list of RoundTrippers that this MultiRoundTripper will delegate to
	roundTrippers []RoundTripperWithURL

	// filter is the filter that will be applied to the responses of the round-trippers
	filter Filter
}

// NewMultiRoundTripper returns a new MultiRoundTripper.
func NewMultiRoundTripper(roundTrippers []RoundTripperWithURL, filter Filter) *MultiRoundTripper {
	return &MultiRoundTripper{
		roundTrippers: roundTrippers,
		filter:        filter,
	}
}

// RoundTrip delegates the request to the underlying RoundTripperWithURLs, waits to collect all responses, and
// applies the filter to the responses. This method adheres to the request.Context's closure.
func (mrt *MultiRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	responsesCh := make(chan *http.Response, len(mrt.roundTrippers))

	// fan-out the requests to the round-trippers
	var wg sync.WaitGroup
	for i := range mrt.roundTrippers {
		wg.Add(1)

		// we have to clone the request for each round-tripper to
		// avoid race-conditions over the request body
		clonedRequest, err := http.NewRequestWithContext(
			req.Context(),
			req.Method,
			req.URL.String(),
			cloneReader(req.Body),
		)
		if err != nil {
			return nil, err
		}

		// make the outbound request concurrently
		go func(rt RoundTripperWithURL) {
			defer wg.Done()
			resp, err := rt.RoundTrip(clonedRequest)
			if err == nil {
				responsesCh <- resp
			}
		}(mrt.roundTrippers[i])
	}

	// wait for all responses to be collected, or for the context to be cancelled
	go func() {
		select {
		case <-doneChannelForWaitGroup(&wg):
		case <-req.Context().Done():
		}
		close(responsesCh)
	}()

	// collect responses
	var responses []*http.Response
	for resp := range responsesCh {
		// only collect valid responses
		if resp != nil && IsValidHTTPResponse(resp) {
			responses = append(responses, resp)
		}
	}

	// apply the filter to the responses
	return mrt.filter(responses)
}

// doneChannelForWaitGroup returns a channel that is closed when the wait group is done.
func doneChannelForWaitGroup(wg *sync.WaitGroup) chan struct{} {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}

// clone reader is a helper function that clones a reader.
func cloneReader(reader io.Reader) io.Reader {
	if reader == nil {
		return nil
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)
	return bytes.NewReader(buffer.Bytes())
}
