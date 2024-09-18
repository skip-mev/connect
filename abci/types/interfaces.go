package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	servertypes "github.com/skip-mev/connect/v2/service/servers/oracle/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
)

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper. This
// interface is utilized by the PreBlock handler to write oracle data to state for the
// supported assets.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface { //golint:ignore
	GetAllCurrencyPairs(ctx sdk.Context) []connecttypes.CurrencyPair
	SetPriceForCurrencyPair(ctx sdk.Context, cp connecttypes.CurrencyPair, qp oracletypes.QuotePrice) error
}

// OracleClient defines the interface that must be fulfilled by the connect client.
// This interface is utilized by the vote extension handler to fetch prices.
type OracleClient interface {
	Prices(ctx context.Context, in *servertypes.QueryPricesRequest, opts ...grpc.CallOption) (*servertypes.QueryPricesResponse, error)
}
