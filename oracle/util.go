package oracle

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// tvwapCandlePeriod represents the time period we use for TVWAP in minutes
	tvwapCandlePeriod = 10 * time.Minute
)

var (
	minimumTimeWeight   = sdk.MustNewDecFromStr("0.2000")
	minimumCandleVolume = sdk.MustNewDecFromStr("0.0001")
)

// ComputeVWAPByProvider executes ComputeVWAP for each provider and returns the
// result.
func ComputeVWAPByProvider(prices types.AggregatedProviderPrices) map[string]map[string]sdk.Dec {
	providerVWAP := make(map[string]map[string]sdk.Dec)

	for pName, pPrices := range prices {
		singleProviderCandles := types.AggregatedProviderPrices{"_": pPrices}
		providerVWAP[pName] = ComputeVWAP(singleProviderCandles)
	}
	return providerVWAP
}

// ComputeVWAP computes the volume weighted average price for all price points
// for each ticker/exchange pair. The provided prices argument reflects a mapping
// of provider => {<base> => <TickerPrice>, ...}.
//
// Ref: https://en.wikipedia.org/wiki/Volume-weighted_average_price
func ComputeVWAP(prices types.AggregatedProviderPrices) map[string]sdk.Dec {
	var (
		weightedPrices = make(map[string]sdk.Dec)
		volumeSum      = make(map[string]sdk.Dec)
	)

	for _, pPrices := range prices {
		for base, tp := range pPrices {
			if _, ok := weightedPrices[base]; !ok {
				weightedPrices[base] = sdk.ZeroDec()
			}
			if _, ok := volumeSum[base]; !ok {
				volumeSum[base] = sdk.ZeroDec()
			}

			// weightedPrices[base] = Σ {P * V} for all TickerPrice
			weightedPrices[base] = weightedPrices[base].Add(tp.Price.Mul(tp.Volume))

			// track total volume for each base
			volumeSum[base] = volumeSum[base].Add(tp.Volume)
		}
	}

	return vwap(weightedPrices, volumeSum)
}

// ComputeTVWAPByProvider executes ComputeTVWAP for each provider and returns
// the result.
func ComputeTVWAPByProvider(providerCandles types.AggregatedProviderCandles) (map[string]map[string]sdk.Dec, error) {
	var (
		providerTVWAP = make(map[string]map[string]sdk.Dec)
		err           error
	)

	for pName, pCandles := range providerCandles {
		singleProviderCandles := types.AggregatedProviderCandles{"_": pCandles}
		providerTVWAP[pName], err = ComputeTVWAP(singleProviderCandles)
		if err != nil {
			return nil, err
		}
	}

	return providerTVWAP, nil
}

// ComputeTVWAP computes the time-volume-weighted average price for all prices
// for each exchange pair. Filters out any candles that did not occur within
// tvwapCandlePeriod. The provided prices argument reflects a mapping of
// provider => {<base> => []Candle}.
//
// Ref : https://en.wikipedia.org/wiki/Time-weighted_average_price
func ComputeTVWAP(providerCandles types.AggregatedProviderCandles) (map[string]sdk.Dec, error) {
	var (
		weightedPrices = make(map[string]sdk.Dec)
		volumeSum      = make(map[string]sdk.Dec)
		now            = time.Now().Unix()
		timePeriod     = time.Now().Add(tvwapCandlePeriod * -1).Unix()
	)

	for _, pCandles := range providerCandles {
		for base := range pCandles {
			candles := pCandles[base]
			if len(candles) == 0 {
				continue
			}

			if _, ok := weightedPrices[base]; !ok {
				weightedPrices[base] = sdk.ZeroDec()
			}
			if _, ok := volumeSum[base]; !ok {
				volumeSum[base] = sdk.ZeroDec()
			}

			// sort by timestamp old -> new
			sort.SliceStable(candles, func(i, j int) bool {
				return candles[i].Timestamp < candles[j].Timestamp
			})

			period := sdk.NewDec(now - candles[0].Timestamp)
			if period.IsZero() {
				return nil, fmt.Errorf("unable to divide by zero")
			}

			// weightUnit = (1 - minimumTimeWeight) / period
			weightUnit := sdk.OneDec().Sub(minimumTimeWeight).Quo(period)

			// get weighted prices, and sum of volumes
			for _, candle := range candles {
				// we only want candles within the last timePeriod
				if timePeriod < candle.Timestamp && candle.Timestamp <= now {
					// timeDiff = now - candle.TimeStamp
					timeDiff := sdk.NewDec(now - candle.Timestamp)

					// set minimum candle volume for low-trading assets
					if candle.Volume.Equal(sdk.ZeroDec()) {
						candle.Volume = minimumCandleVolume
					}

					// volume = candle.Volume * (weightUnit * (period - timeDiff) + minimumTimeWeight)
					volume := candle.Volume.Mul(
						weightUnit.Mul(period.Sub(timeDiff).Add(minimumTimeWeight)),
					)

					volumeSum[base] = volumeSum[base].Add(volume)
					weightedPrices[base] = weightedPrices[base].Add(candle.Price.Mul(volume))
				}
			}

		}
	}

	return vwap(weightedPrices, volumeSum), nil
}

// compute VWAP for each base by dividing the Σ {P * V} by Σ {V}
func vwap(weightedPrices, volumeSum map[string]sdk.Dec) map[string]sdk.Dec {
	vwap := make(map[string]sdk.Dec)

	for base, p := range weightedPrices {
		if !volumeSum[base].Equal(sdk.ZeroDec()) {
			if _, ok := vwap[base]; !ok {
				vwap[base] = sdk.ZeroDec()
			}

			vwap[base] = p.Quo(volumeSum[base])
		}
	}

	return vwap
}
