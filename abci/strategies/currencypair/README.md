# Currency Pair Strategies

## Overview

This document overviews the different price information strategies that are available for applications to use. These strategies are primarily utilized to optimize how much data is transmitted over the wire, with two primary implementations:

1. **DefaultCurrencyPairStrategy**: This strategy utilizes raw prices.
2. **DeltaCurrencyPairStrategy**: This strategy utilizes the delta between the current price and the previous price.

## DefaultCurrencyPairStrategy

The default strategy is the simplest strategy, but is not the most efficient. This strategy simply transmits the raw price information for each currency pair. As a result, a single price update may take up to 32 bytes of data.

## DeltaCurrencyPairStrategy

The delta strategy is a more efficient strategy, but is more complex. This strategy transmits the delta between the current price and the previous price. As a result, the worst case scenario remains the same as the default strategy, but the average case scenario is much more efficient. This strategy is most efficient when the price changes are small.

## Usage

To implement a custom strategy, simply implement the `CurrencyPairStrategy` interface. The `CurrencyPairStrategy` interface is defined as follows:

```go
// CurrencyPairStrategy is a strategy for generating a unique ID and price representation for a given currency pair.
type CurrencyPairStrategy interface {
	// ID returns the on-chain ID of the given currency pair. This method returns an error if the given currency
	// pair is not found in the x/oracle state.
	ID(ctx sdk.Context, cp oracletypes.CurrencyPair) (uint64, error)

	// FromID returns the currency pair with the given ID. This method returns an error if the given ID is not
	// currently present for an existing currency pair.
	FromID(ctx sdk.Context, id uint64) (oracletypes.CurrencyPair, error)

	// GetEncodedPrice returns the encoded price for the given currency pair. This method returns an error if the
	// given currency pair is not found in the x/oracle state or if the price cannot be encoded.
	GetEncodedPrice(
		ctx sdk.Context,
		cp oracletypes.CurrencyPair,
		price *big.Int,
	) ([]byte, error)

	// GetDecodedPrice returns the decoded price for the given currency pair. This method returns an error if the
	// given currency pair is not found in the x/oracle state or if the price cannot be decoded.
	GetDecodedPrice(
		ctx sdk.Context,
		cp oracletypes.CurrencyPair,
		priceBytes []byte,
	) (*big.Int, error)
} 
```
