package erc4626sharepriceoracle

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/providers/evm"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func (p *Provider) getPriceForPair(pair oracletypes.CurrencyPair) (aggregator.QuotePrice, error) {
	metadata, ok := p.config.TokenNameToMetadata[pair.Quote]
	if !ok {
		return aggregator.QuotePrice{}, fmt.Errorf("token %s metadata not found", pair.Quote)
	}

	client, err := ethclient.Dial(p.rpcEndpoint)
	if err != nil {
		return aggregator.QuotePrice{}, err
	}

	contractAddress, found := p.getPairContractAddress(pair)
	if !found {
		return aggregator.QuotePrice{}, fmt.Errorf("contract address for pair %v not found", pair)
	}

	contract, err := NewERC4626SharePriceOracle(common.HexToAddress(contractAddress), client)
	if err != nil {
		return aggregator.QuotePrice{}, err
	}

	latest, err := contract.GetLatest(&bind.CallOpts{})
	if err != nil || latest.NotSafeToUse {
		return aggregator.QuotePrice{}, err
	}

	var _price *big.Int
	var price *uint256.Int
	if metadata.IsTWAP {
		_price = latest.TimeWeightedAverageAnswer
	} else {
		_price = latest.Ans
	}
	price, ok = uint256.FromBig(_price)
	if !ok {
		return aggregator.QuotePrice{}, fmt.Errorf("failed to convert price %v to uint256 for pair %v", _price, pair)
	}

	quote, err := aggregator.NewQuotePrice(price, time.Now())
	if err != nil {
		return aggregator.QuotePrice{}, err
	}

	return quote, nil
}

// getRPCEndpoint returns the endpoint to fetch prices from.
func getRPCEndpoint(config evm.Config) string {
	return fmt.Sprintf("%s/%s", config.RPCEndpoint, config.APIKey)
}
