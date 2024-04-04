package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	config "github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	url              = "https://api.raydium.io/v2/main/pairs"
	solanaNodeEnvVar = "SOLANA_NODE"
)

var nodes []string

func init() {
	nodes = strings.Split(os.Getenv(solanaNodeEnvVar), ",")
	if len(nodes) == 0 {
		panic("SOLANA_NODE environment variable not set")
	}
}

type TickerMetadata struct {
	cp             slinkytypes.CurrencyPair
	tickerMetaData raydium.TickerMetadata
}

func main() {
	tickers, err := getTickerMetadata()
	if err != nil {
		panic(err)
	}

	mm := makeMarketMap(tickers)

	cfg := config.OracleConfig{
		UpdateInterval: 500 * time.Millisecond,
		MaxPriceAge:    5 * time.Minute,
		Providers: []config.ProviderConfig{
			{
				Name: raydium.Name,
				API: config.APIConfig{
					Enabled:          true,
					Timeout:          5 * time.Second,
					Interval:         500 * time.Millisecond,
					ReconnectTimeout: 5 * time.Second,
					MaxQueries:       100,
					Name:             raydium.Name,
					Atomic:           true,
				},
				Type: "price_provider",
			},
		},
		Production: false,
		Host:       "localhost",
		Port:       "8080",
		Metrics: config.MetricsConfig{
			Enabled:                 true,
			PrometheusServerAddress: "localhost:8081",
		},
	}

	// update endpoints
	for _, node := range nodes {
		cfg.Providers[0].API.Endpoints = append(cfg.Providers[0].API.Endpoints, config.Endpoint{
			URL: node,
		})
	}

	// write the oracle config to oracle.json
	bz, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile("oracle.json", bz, 0o600); err != nil {
		panic(err)
	}

	// write the market map to marketmap.json
	bz, err = json.MarshalIndent(mm, "", " ")
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile("market.json", bz, 0o600); err != nil {
		panic(err)
	}
}

func getTickerMetadata() ([]TickerMetadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var data []map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	tickers := make([]TickerMetadata, 0)
	accounts := make([]solana.PublicKey, 0)
	for _, datum := range data {
		if len(tickers) >= 50 {
			break
		}

		pairs := strings.Split(strings.ToUpper(datum["name"].(string)), "-")

		if pairs[0] == "UNKNOWN" || pairs[1] == "UNKNOWN" {
			continue
		}

		cp := slinkytypes.CurrencyPair{
			Base:  pairs[0],
			Quote: pairs[1],
		}

		// get the market account
		amm := datum["ammId"].(string)
		tickers = append(tickers, TickerMetadata{
			cp: cp,
		})
		accounts = append(accounts, solana.MustPublicKeyFromBase58(amm))
	}

	client := rpc.New(nodes[0])

	// get amm infos
	infos, err := client.GetMultipleAccounts(context.Background(), accounts...)
	if err != nil {
		return nil, err
	}

	if len(infos.Value) != len(tickers) {
		return nil, fmt.Errorf("expected %d accounts, got %d", len(tickers), len(infos.Value))
	}

	for i, info := range infos.Value {
		// unmarshal into account info
		var accInfo AmmInfo
		if err = bin.NewBinDecoder(info.Data.GetBinary()).Decode(&accInfo); err != nil {
			return nil, err
		}

		tickers[i].tickerMetaData = raydium.TickerMetadata{
			BaseTokenVault: raydium.AMMTokenVaultMetadata{
				TokenVaultAddress: accInfo.TokenCoin.String(),
				TokenDecimals:     accInfo.CoinDecimals,
			},
			QuoteTokenVault: raydium.AMMTokenVaultMetadata{
				TokenVaultAddress: accInfo.TokenPc.String(),
				TokenDecimals:     accInfo.PcDecimals,
			},
		}
	}

	return tickers, nil
}

func makeMarketMap(tickers []TickerMetadata) mmtypes.MarketMap {
	mm := mmtypes.MarketMap{
		Tickers:   make(map[string]mmtypes.Ticker),
		Paths:     make(map[string]mmtypes.Paths),
		Providers: make(map[string]mmtypes.Providers),
	}

	for _, ticker := range tickers {
		mm.Tickers[ticker.cp.String()] = mmtypes.Ticker{
			CurrencyPair:     ticker.cp,
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
			Metadata_JSON:    marshalDataToJSON(ticker.tickerMetaData),
		}

		mm.Paths[ticker.cp.String()] = mmtypes.Paths{
			Paths: []mmtypes.Path{
				{
					Operations: []mmtypes.Operation{
						{
							Provider:     raydium.Name,
							CurrencyPair: ticker.cp,
						},
					},
				},
			},
		}

		mm.Providers[ticker.cp.String()] = mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				{
					Name:           raydium.Name,
					OffChainTicker: ticker.cp.String(),
				},
			},
		}
	}

	return mm
}

func marshalDataToJSON(data interface{}) string {
	bz, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}

	return string(bz)
}

type AmmInfo struct {
	Status             uint64
	Nonce              uint64
	OrderNum           uint64
	Depth              uint64
	CoinDecimals       uint64
	PcDecimals         uint64
	State              uint64
	ResetFlag          uint64
	MinSize            uint64
	VolMaxCutRatio     uint64
	AmountWave         uint64
	CoinLotSize        uint64
	PcLotSize          uint64
	MinPriceMultiplier uint64
	MaxPriceMultiplier uint64
	SysDecimalValue    uint64
	Fees               Fees
	OutPut             OutPutData
	TokenCoin          solana.PublicKey
	TokenPc            solana.PublicKey
	CoinMint           solana.PublicKey
	PcMint             solana.PublicKey
	LpMint             solana.PublicKey
	OpenOrders         solana.PublicKey
	Market             solana.PublicKey
	SerumDex           solana.PublicKey
	TargetOrders       solana.PublicKey
	WithdrawQueue      solana.PublicKey
	TokenTempLp        solana.PublicKey
	AmmOwner           solana.PublicKey
	LpAmount           uint64
	ClientOrderID      uint64
	Padding            [2]uint64
}

var AmmInfoDiscriminator = [8]byte{33, 217, 2, 203, 184, 83, 235, 91}

func (obj AmmInfo) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(AmmInfoDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Status` param:
	err = encoder.Encode(obj.Status)
	if err != nil {
		return err
	}
	// Serialize `Nonce` param:
	err = encoder.Encode(obj.Nonce)
	if err != nil {
		return err
	}
	// Serialize `OrderNum` param:
	err = encoder.Encode(obj.OrderNum)
	if err != nil {
		return err
	}
	// Serialize `Depth` param:
	err = encoder.Encode(obj.Depth)
	if err != nil {
		return err
	}
	// Serialize `CoinDecimals` param:
	err = encoder.Encode(obj.CoinDecimals)
	if err != nil {
		return err
	}
	// Serialize `PcDecimals` param:
	err = encoder.Encode(obj.PcDecimals)
	if err != nil {
		return err
	}
	// Serialize `State` param:
	err = encoder.Encode(obj.State)
	if err != nil {
		return err
	}
	// Serialize `ResetFlag` param:
	err = encoder.Encode(obj.ResetFlag)
	if err != nil {
		return err
	}
	// Serialize `MinSize` param:
	err = encoder.Encode(obj.MinSize)
	if err != nil {
		return err
	}
	// Serialize `VolMaxCutRatio` param:
	err = encoder.Encode(obj.VolMaxCutRatio)
	if err != nil {
		return err
	}
	// Serialize `AmountWave` param:
	err = encoder.Encode(obj.AmountWave)
	if err != nil {
		return err
	}
	// Serialize `CoinLotSize` param:
	err = encoder.Encode(obj.CoinLotSize)
	if err != nil {
		return err
	}
	// Serialize `PcLotSize` param:
	err = encoder.Encode(obj.PcLotSize)
	if err != nil {
		return err
	}
	// Serialize `MinPriceMultiplier` param:
	err = encoder.Encode(obj.MinPriceMultiplier)
	if err != nil {
		return err
	}
	// Serialize `MaxPriceMultiplier` param:
	err = encoder.Encode(obj.MaxPriceMultiplier)
	if err != nil {
		return err
	}
	// Serialize `SysDecimalValue` param:
	err = encoder.Encode(obj.SysDecimalValue)
	if err != nil {
		return err
	}
	// Serialize `Fees` param:
	err = encoder.Encode(obj.Fees)
	if err != nil {
		return err
	}
	// Serialize `OutPut` param:
	err = encoder.Encode(obj.OutPut)
	if err != nil {
		return err
	}
	// Serialize `TokenCoin` param:
	err = encoder.Encode(obj.TokenCoin)
	if err != nil {
		return err
	}
	// Serialize `TokenPc` param:
	err = encoder.Encode(obj.TokenPc)
	if err != nil {
		return err
	}
	// Serialize `CoinMint` param:
	err = encoder.Encode(obj.CoinMint)
	if err != nil {
		return err
	}
	// Serialize `PcMint` param:
	err = encoder.Encode(obj.PcMint)
	if err != nil {
		return err
	}
	// Serialize `LpMint` param:
	err = encoder.Encode(obj.LpMint)
	if err != nil {
		return err
	}
	// Serialize `OpenOrders` param:
	err = encoder.Encode(obj.OpenOrders)
	if err != nil {
		return err
	}
	// Serialize `Market` param:
	err = encoder.Encode(obj.Market)
	if err != nil {
		return err
	}
	// Serialize `SerumDex` param:
	err = encoder.Encode(obj.SerumDex)
	if err != nil {
		return err
	}
	// Serialize `TargetOrders` param:
	err = encoder.Encode(obj.TargetOrders)
	if err != nil {
		return err
	}
	// Serialize `WithdrawQueue` param:
	err = encoder.Encode(obj.WithdrawQueue)
	if err != nil {
		return err
	}
	// Serialize `TokenTempLp` param:
	err = encoder.Encode(obj.TokenTempLp)
	if err != nil {
		return err
	}
	// Serialize `AmmOwner` param:
	err = encoder.Encode(obj.AmmOwner)
	if err != nil {
		return err
	}
	// Serialize `LpAmount` param:
	err = encoder.Encode(obj.LpAmount)
	if err != nil {
		return err
	}
	// Serialize `ClientOrderId` param:
	err = encoder.Encode(obj.ClientOrderID)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

type Fees struct {
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
}

var FeesDiscriminator = [8]byte{151, 157, 50, 115, 130, 72, 179, 36}

func (obj Fees) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(FeesDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `MinSeparateNumerator` param:
	err = encoder.Encode(obj.MinSeparateNumerator)
	if err != nil {
		return err
	}
	// Serialize `MinSeparateDenominator` param:
	err = encoder.Encode(obj.MinSeparateDenominator)
	if err != nil {
		return err
	}
	// Serialize `TradeFeeNumerator` param:
	err = encoder.Encode(obj.TradeFeeNumerator)
	if err != nil {
		return err
	}
	// Serialize `TradeFeeDenominator` param:
	err = encoder.Encode(obj.TradeFeeDenominator)
	if err != nil {
		return err
	}
	// Serialize `PnlNumerator` param:
	err = encoder.Encode(obj.PnlNumerator)
	if err != nil {
		return err
	}
	// Serialize `PnlDenominator` param:
	err = encoder.Encode(obj.PnlDenominator)
	if err != nil {
		return err
	}
	// Serialize `SwapFeeNumerator` param:
	err = encoder.Encode(obj.SwapFeeNumerator)
	if err != nil {
		return err
	}
	// Serialize `SwapFeeDenominator` param:
	err = encoder.Encode(obj.SwapFeeDenominator)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Fees) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Read and check account discriminator:
	// {
	// 	discriminator, err := decoder.ReadTypeID()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !discriminator.Equal(FeesDiscriminator[:]) {
	// 		return fmt.Errorf(
	// 			"wrong discriminator: wanted %s, got %s",
	// 			"[151 157 50 115 130 72 179 36]",
	// 			fmt.Sprint(discriminator[:]))
	// 	}
	// }
	// Deserialize `MinSeparateNumerator`:
	err = decoder.Decode(&obj.MinSeparateNumerator)
	if err != nil {
		return err
	}
	// Deserialize `MinSeparateDenominator`:
	err = decoder.Decode(&obj.MinSeparateDenominator)
	if err != nil {
		return err
	}
	// Deserialize `TradeFeeNumerator`:
	err = decoder.Decode(&obj.TradeFeeNumerator)
	if err != nil {
		return err
	}
	// Deserialize `TradeFeeDenominator`:
	err = decoder.Decode(&obj.TradeFeeDenominator)
	if err != nil {
		return err
	}
	// Deserialize `PnlNumerator`:
	err = decoder.Decode(&obj.PnlNumerator)
	if err != nil {
		return err
	}
	// Deserialize `PnlDenominator`:
	err = decoder.Decode(&obj.PnlDenominator)
	if err != nil {
		return err
	}
	// Deserialize `SwapFeeNumerator`:
	err = decoder.Decode(&obj.SwapFeeNumerator)
	if err != nil {
		return err
	}
	// Deserialize `SwapFeeDenominator`:
	err = decoder.Decode(&obj.SwapFeeDenominator)
	if err != nil {
		return err
	}
	return nil
}

type OutPutData struct {
	NeedTakePnlCoin     uint64
	NeedTakePnlPc       uint64
	TotalPnlPc          uint64
	TotalPnlCoin        uint64
	PoolOpenTime        uint64
	PunishPcAmount      uint64
	PunishCoinAmount    uint64
	OrderbookToInitTime uint64
	SwapCoinInAmount    bin.Uint128
	SwapPcOutAmount     bin.Uint128
	SwapTakePcFee       uint64
	SwapPcInAmount      bin.Uint128
	SwapCoinOutAmount   bin.Uint128
	SwapTakeCoinFee     uint64
}

func (obj OutPutData) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Serialize `NeedTakePnlCoin` param:
	err = encoder.Encode(obj.NeedTakePnlCoin)
	if err != nil {
		return err
	}
	// Serialize `NeedTakePnlPc` param:
	err = encoder.Encode(obj.NeedTakePnlPc)
	if err != nil {
		return err
	}
	// Serialize `TotalPnlPc` param:
	err = encoder.Encode(obj.TotalPnlPc)
	if err != nil {
		return err
	}
	// Serialize `TotalPnlCoin` param:
	err = encoder.Encode(obj.TotalPnlCoin)
	if err != nil {
		return err
	}
	// Serialize `PoolOpenTime` param:
	err = encoder.Encode(obj.PoolOpenTime)
	if err != nil {
		return err
	}
	// Serialize `PunishPcAmount` param:
	err = encoder.Encode(obj.PunishPcAmount)
	if err != nil {
		return err
	}
	// Serialize `PunishCoinAmount` param:
	err = encoder.Encode(obj.PunishCoinAmount)
	if err != nil {
		return err
	}
	// Serialize `OrderbookToInitTime` param:
	err = encoder.Encode(obj.OrderbookToInitTime)
	if err != nil {
		return err
	}
	// Serialize `SwapCoinInAmount` param:
	err = encoder.Encode(obj.SwapCoinInAmount)
	if err != nil {
		return err
	}
	// Serialize `SwapPcOutAmount` param:
	err = encoder.Encode(obj.SwapPcOutAmount)
	if err != nil {
		return err
	}
	// Serialize `SwapTakePcFee` param:
	err = encoder.Encode(obj.SwapTakePcFee)
	if err != nil {
		return err
	}
	// Serialize `SwapPcInAmount` param:
	err = encoder.Encode(obj.SwapPcInAmount)
	if err != nil {
		return err
	}
	// Serialize `SwapCoinOutAmount` param:
	err = encoder.Encode(obj.SwapCoinOutAmount)
	if err != nil {
		return err
	}
	// Serialize `SwapTakeCoinFee` param:
	err = encoder.Encode(obj.SwapTakeCoinFee)
	if err != nil {
		return err
	}
	return nil
}

func (obj *OutPutData) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Deserialize `NeedTakePnlCoin`:
	err = decoder.Decode(&obj.NeedTakePnlCoin)
	if err != nil {
		return err
	}
	// Deserialize `NeedTakePnlPc`:
	err = decoder.Decode(&obj.NeedTakePnlPc)
	if err != nil {
		return err
	}
	// Deserialize `TotalPnlPc`:
	err = decoder.Decode(&obj.TotalPnlPc)
	if err != nil {
		return err
	}
	// Deserialize `TotalPnlCoin`:
	err = decoder.Decode(&obj.TotalPnlCoin)
	if err != nil {
		return err
	}
	// Deserialize `PoolOpenTime`:
	err = decoder.Decode(&obj.PoolOpenTime)
	if err != nil {
		return err
	}
	// Deserialize `PunishPcAmount`:
	err = decoder.Decode(&obj.PunishPcAmount)
	if err != nil {
		return err
	}
	// Deserialize `PunishCoinAmount`:
	err = decoder.Decode(&obj.PunishCoinAmount)
	if err != nil {
		return err
	}
	// Deserialize `OrderbookToInitTime`:
	err = decoder.Decode(&obj.OrderbookToInitTime)
	if err != nil {
		return err
	}
	// Deserialize `SwapCoinInAmount`:
	err = decoder.Decode(&obj.SwapCoinInAmount)
	if err != nil {
		return err
	}
	// Deserialize `SwapPcOutAmount`:
	err = decoder.Decode(&obj.SwapPcOutAmount)
	if err != nil {
		return err
	}
	// Deserialize `SwapTakePcFee`:
	err = decoder.Decode(&obj.SwapTakePcFee)
	if err != nil {
		return err
	}
	// Deserialize `SwapPcInAmount`:
	err = decoder.Decode(&obj.SwapPcInAmount)
	if err != nil {
		return err
	}
	// Deserialize `SwapCoinOutAmount`:
	err = decoder.Decode(&obj.SwapCoinOutAmount)
	if err != nil {
		return err
	}
	// Deserialize `SwapTakeCoinFee`:
	err = decoder.Decode(&obj.SwapTakeCoinFee)
	if err != nil {
		return err
	}
	return nil
}
