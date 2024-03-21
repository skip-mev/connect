package constants

// Chain corresponds to the chain we are running the oracle for.
type Chain struct {
	// Name is the name of the chain.
	Name string `json:"name"`

	// ID is the chain id.
	ID string `json:"id"`
}

var (
	// Ref: https://v4-teacher.vercel.app/network/network_constants
	//
	// DYDXMainnet is the chain id for the mainnet.
	DYDXMainnet = Chain{
		Name: "dYdX Mainnet",
		ID:   "dydx-mainnet-1",
	}
	// DYDXTestnet is the chain id for the testnet.
	DYDXTestnet = Chain{
		Name: "dYdX Testnet",
		ID:   "dydx-testnet-4",
	}
)
