package raydium

// Option is a function that can be used to modify the APIPriceFetcher.
type Option func(*APIPriceFetcher)

// WithSolanaClient sets the SolanaJSONRPCClient used to query the API.
func WithSolanaClient(client SolanaJSONRPCClient) Option {
	return func(f *APIPriceFetcher) {
		f.client = client
	}
}
