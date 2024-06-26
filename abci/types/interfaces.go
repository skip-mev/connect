package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	servertypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// OracleKeeper defines the interface that must be fulfilled by the oracle keeper. This
// interface is utilized by the PreBlock handler to write oracle data to state for the
// supported assets.
//
//go:generate mockery --name OracleKeeper --filename mock_oracle_keeper.go
type OracleKeeper interface { //golint:ignore
	GetAllCurrencyPairs(ctx sdk.Context) []slinkytypes.CurrencyPair
	SetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair, qp oracletypes.QuotePrice) error
}

// OracleClient defines the interface that must be fulfilled by the slinky client.
// This interface is utilized by the vote extension handler to fetch prices.
type OracleClient interface {
	Prices(ctx context.Context, in *servertypes.QueryPricesRequest, opts ...grpc.CallOption) (*servertypes.QueryPricesResponse, error)
}
