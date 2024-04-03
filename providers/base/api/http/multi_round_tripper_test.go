package http_test

import (
	"context"
	gohttp "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers/base/api/http"
	"github.com/skip-mev/slinky/providers/base/api/http/mocks"
)

func TestMultiRoundTripper(t *testing.T) {
	rt1 := mocks.NewRoundTripper(t)
	url1 := "http://one.com"
	rt2 := mocks.NewRoundTripper(t)
	url2 := "http://two.com"
	rt3 := mocks.NewRoundTripper(t)
	url3 := "http://three.com"
	rt1WithURL, err := http.NewRoundTripperWithURL(url1, rt1)
	require.NoError(t, err)
	rt2WithURL, err := http.NewRoundTripperWithURL(url2, rt2)
	require.NoError(t, err)
	rt3WithURL, err := http.NewRoundTripperWithURL(url3, rt3)
	require.NoError(t, err)

	// test that MultiRoundTripper adheres to the given context
	t.Run("test that MultiRoundTripper adheres to context", func(t *testing.T) {
		mrt := http.NewMultiRoundTripper([]http.RoundTripperWithURL{
			*rt1WithURL,
			*rt2WithURL,
			*rt3WithURL,
		}, func(req []*gohttp.Response) (*gohttp.Response, error) {
			require.Len(t, req, 0)
			return nil, nil
		})

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		req, err := gohttp.NewRequestWithContext(
			ctx,
			gohttp.MethodGet,
			"http://test.com", // this will be over-written
			nil,
		)
		require.NoError(t, err)

		// assert mocks
		req1 := req.Clone(ctx)
		req1.URL = rt1WithURL.URL()
		rt1.On("RoundTrip", req1).Return(nil, nil).After(10 * time.Second).Once()

		req2 := req.Clone(ctx)
		req2.URL = rt2WithURL.URL()
		rt2.On("RoundTrip", req2).Return(nil, nil).After(10 * time.Second).Once()

		req3 := req.Clone(ctx)
		req3.URL = rt3WithURL.URL()
		rt3.On("RoundTrip", req3).Return(nil, nil).After(10 * time.Second).Once()

		cancel()
		_, err = mrt.RoundTrip(req)
		require.NoError(t, err)
	})

	// test that the filter is applied to responses
	t.Run("test that the filter is applied to responses", func(t *testing.T) {
		ctx := context.Background()

		req, err := gohttp.NewRequestWithContext(
			ctx,
			gohttp.MethodGet,
			"http://test.com", // this will be over-written
			nil,
		)
		require.NoError(t, err)

		// assert mocks
		req1 := req.Clone(ctx)
		req1.URL = rt1WithURL.URL()
		res1 := &gohttp.Response{
			StatusCode: gohttp.StatusOK,
			Header: gohttp.Header{
				"round-tripper1": []string{},
			},
		}
		rt1.On("RoundTrip", req1).Return(res1, nil)

		req2 := req.Clone(ctx)
		req2.URL = rt2WithURL.URL()
		res2 := &gohttp.Response{
			StatusCode: gohttp.StatusBadGateway,
		}
		rt2.On("RoundTrip", req2).Return(res2, nil)

		req3 := req.Clone(ctx)
		req3.URL = rt3WithURL.URL()
		res3 := &gohttp.Response{
			StatusCode: gohttp.StatusOK,
			Header: gohttp.Header{
				"round-tripper3": []string{},
			},
		}
		rt3.On("RoundTrip", req3).Return(res3, nil)

		expectedResponses := map[string]struct{}{
			"round-tripper3": {},
			"round-tripper1": {},
		}

		mrt := http.NewMultiRoundTripper([]http.RoundTripperWithURL{
			*rt1WithURL,
			*rt2WithURL,
			*rt3WithURL,
		}, func(req []*gohttp.Response) (*gohttp.Response, error) {
			// expect only res 1 and 3
			require.Len(t, req, 2)
			for _, r := range req {
				// expect a single header-field
				require.Len(t, r.Header, 1)
				for k := range r.Header {
					_, ok := expectedResponses[k]
					require.True(t, ok)
					delete(expectedResponses, k)
				}
			}

			return nil, nil
		})
		_, err = mrt.RoundTrip(req)
		require.NoError(t, err)
	})
}
