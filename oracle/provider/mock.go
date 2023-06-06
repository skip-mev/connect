package provider

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

const ProviderNameMock = "mock"

var _ Provider = (*MockProvider)(nil)

type (
	// MockProvider defines a mocked exchange rate provider using fixed exchange
	// rates.
	MockProvider struct {
		exchangeRates map[string]types.TickerPrice
	}
)

func NewMockProvider() *MockProvider {
	return &MockProvider{
		exchangeRates: map[string]types.TickerPrice{
			"ATOM/USDC": {Price: sdk.MustNewDecFromStr("11.34"), Volume: sdk.MustNewDecFromStr("1827884.77")},
			"ATOM/USDT": {Price: sdk.MustNewDecFromStr("11.36"), Volume: sdk.MustNewDecFromStr("1827884.77")},
			"ATOM/USD":  {Price: sdk.MustNewDecFromStr("11.35"), Volume: sdk.MustNewDecFromStr("1827884.77")},
			"OSMO/USDC": {Price: sdk.MustNewDecFromStr("1.34"), Volume: sdk.MustNewDecFromStr("834120.21")},
			"OSMO/USDT": {Price: sdk.MustNewDecFromStr("1.36"), Volume: sdk.MustNewDecFromStr("834120.21")},
			"OSMO/USD":  {Price: sdk.MustNewDecFromStr("1.35"), Volume: sdk.MustNewDecFromStr("834120.21")},
			"WETH/USDC": {Price: sdk.MustNewDecFromStr("1560.34"), Volume: sdk.MustNewDecFromStr("51342578.34")},
			"WETH/USDT": {Price: sdk.MustNewDecFromStr("1560.36"), Volume: sdk.MustNewDecFromStr("51342578.34")},
			"WETH/USD":  {Price: sdk.MustNewDecFromStr("1560.35"), Volume: sdk.MustNewDecFromStr("51342578.34")},
		},
	}
}

func (p MockProvider) Name() string {
	return ProviderNameMock
}

func (p MockProvider) GetTickerPrices(pairs ...types.CurrencyPair) (map[string]types.TickerPrice, error) {
	tickerMap := make(map[string]struct{})
	for _, cp := range pairs {
		tickerMap[strings.ToUpper(cp.String())] = struct{}{}
	}

	tickerPrices := make(map[string]types.TickerPrice, len(pairs))
	for ticker, er := range p.exchangeRates {
		if _, ok := tickerMap[ticker]; !ok {
			// skip records that are not requested
			continue
		}

		if _, ok := tickerPrices[ticker]; ok {
			return nil, fmt.Errorf("duplicate ticker: %s", ticker)
		}

		tickerPrices[ticker] = types.TickerPrice{Price: er.Price, Volume: er.Volume}
	}

	for t := range tickerMap {
		if _, ok := tickerPrices[t]; !ok {
			return nil, fmt.Errorf("%s: %w", t, ErrMissingExchangeRate)
		}
	}

	return tickerPrices, nil
}

func (p MockProvider) GetCandlePrices(pairs ...types.CurrencyPair) (map[string][]types.Candle, error) {
	prices, err := p.GetTickerPrices(pairs...)
	if err != nil {
		return nil, err
	}

	ts := time.Now().Add(time.Minute * -1).Unix() // 1 minute ago

	candles := make(map[string][]types.Candle, len(prices))
	for ticker, price := range prices {
		candles[ticker] = []types.Candle{
			{
				Price:     price.Price,
				Volume:    price.Volume,
				Timestamp: ts,
			},
		}
	}
	return candles, nil
}

func (p MockProvider) GetAvailablePairs() (map[string]struct{}, error) {
	availablePairs := make(map[string]struct{}, len(p.exchangeRates))
	for ticker := range p.exchangeRates {
		availablePairs[ticker] = struct{}{}
	}

	return availablePairs, nil
}
