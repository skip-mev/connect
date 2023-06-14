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

// Get the QuotePrice for a given CurrencyPair
func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (types.QuotePrice, error) {
	// check validity of cp
	if err := cp.ValidateBasic(); err != nil {
		return types.QuotePrice{}, err
	}

	// get QuotePrice for CurrencyPair (if any is stored)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(cp.GetStoreKeyForCurrencyPair())
	
	if len(bz) == 0 {
		return types.QuotePrice{}, fmt.Errorf("no CurrencyPair price found for CurrencyPair")
	}

	tp := types.QuotePrice{}
	if err := tp.Unmarshal(bz); err != nil {
		return types.QuotePrice{}, err
	}
	return tp, nil
}

// Set the QuotePrice for a given CurrencyPair
func (k Keeper) SetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair, qp types.QuotePrice) error {
	// check validity of currency pair
	if err := cp.ValidateBasic(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	// marshal QuotePrice
	bz, err := qp.Marshal()
	if err != nil {
		return fmt.Errorf("error marshalling QuotePrice: %v", err)
	}
	// set the marshalled QuotePrice to state under the CurrencyPair's store-key
	store.Set(cp.GetStoreKeyForCurrencyPair(), bz)
	return nil
}


func (k Keeper) GetAllTickers(ctx sdk.Context) ([]types.CurrencyPair, error) {
}