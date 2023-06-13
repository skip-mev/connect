package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
}

func NewKeeper(sk storetypes.StoreKey) Keeper {
	return Keeper{
		storeKey: sk,
	}
}

// Get the TickerPrice for a given Ticker
func (k Keeper) GetPriceForTicker(ctx sdk.Context, t types.Ticker) (types.TickerPrice, error) {
	// get TickerPrice for ticker (if any is stored)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(t.GetStoreKeyForTicker())
	
	if len(bz) == 0 {
		return types.TickerPrice{}, fmt.Errorf("no ticker price found for ticker")
	}

	tp := types.TickerPrice{}
	if err := tp.Unmarshal(bz); err != nil {
		return types.TickerPrice{}, err
	}
	return tp, nil
}

// Set the TickerPrice for a given Ticker
func (k Keeper) SetPriceForTicker(ctx sdk.Context, t types.Ticker, tp types.TickerPrice) error {
	store := ctx.KVStore(k.storeKey)
	// marshal TickerPrice
	bz, err := tp.Marshal()
	if err != nil {
		return fmt.Errorf("error marshalling TickerPrice: %v", err)
	}
	// set the marshalled TickerPrice to state under the Ticker's store-key
	store.Set(t.GetStoreKeyForTicker(), bz)
	return nil
}
