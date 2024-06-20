// THIS IS GENERATED CODE. SPECIFICALLY, THIS IS ANCHOR BINDINGS FOR THE RAYDIUM V4 CONTRACTS
// GENERATED VIA solana-go.
package schema

import (
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

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

func (obj Fees) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
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

func (obj *Fees) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
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
	TokenCoin          ag_solanago.PublicKey
	TokenPc            ag_solanago.PublicKey
	CoinMint           ag_solanago.PublicKey
	PcMint             ag_solanago.PublicKey
	LpMint             ag_solanago.PublicKey
	OpenOrders         ag_solanago.PublicKey
	Market             ag_solanago.PublicKey
	SerumDex           ag_solanago.PublicKey
	TargetOrders       ag_solanago.PublicKey
	WithdrawQueue      ag_solanago.PublicKey
	TokenTempLp        ag_solanago.PublicKey
	AmmOwner           ag_solanago.PublicKey
	LpAmount           uint64
	ClientOrderId      uint64 //nolint:all
	Padding            [2]uint64
}

func (obj AmmInfo) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
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
	err = encoder.Encode(obj.ClientOrderId)
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

func (obj *AmmInfo) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Status`:
	err = decoder.Decode(&obj.Status)
	if err != nil {
		return err
	}
	// Deserialize `Nonce`:
	err = decoder.Decode(&obj.Nonce)
	if err != nil {
		return err
	}
	// Deserialize `OrderNum`:
	err = decoder.Decode(&obj.OrderNum)
	if err != nil {
		return err
	}
	// Deserialize `Depth`:
	err = decoder.Decode(&obj.Depth)
	if err != nil {
		return err
	}
	// Deserialize `CoinDecimals`:
	err = decoder.Decode(&obj.CoinDecimals)
	if err != nil {
		return err
	}
	// Deserialize `PcDecimals`:
	err = decoder.Decode(&obj.PcDecimals)
	if err != nil {
		return err
	}
	// Deserialize `State`:
	err = decoder.Decode(&obj.State)
	if err != nil {
		return err
	}
	// Deserialize `ResetFlag`:
	err = decoder.Decode(&obj.ResetFlag)
	if err != nil {
		return err
	}
	// Deserialize `MinSize`:
	err = decoder.Decode(&obj.MinSize)
	if err != nil {
		return err
	}
	// Deserialize `VolMaxCutRatio`:
	err = decoder.Decode(&obj.VolMaxCutRatio)
	if err != nil {
		return err
	}
	// Deserialize `AmountWave`:
	err = decoder.Decode(&obj.AmountWave)
	if err != nil {
		return err
	}
	// Deserialize `CoinLotSize`:
	err = decoder.Decode(&obj.CoinLotSize)
	if err != nil {
		return err
	}
	// Deserialize `PcLotSize`:
	err = decoder.Decode(&obj.PcLotSize)
	if err != nil {
		return err
	}
	// Deserialize `MinPriceMultiplier`:
	err = decoder.Decode(&obj.MinPriceMultiplier)
	if err != nil {
		return err
	}
	// Deserialize `MaxPriceMultiplier`:
	err = decoder.Decode(&obj.MaxPriceMultiplier)
	if err != nil {
		return err
	}
	// Deserialize `SysDecimalValue`:
	err = decoder.Decode(&obj.SysDecimalValue)
	if err != nil {
		return err
	}
	// Deserialize `Fees`:
	err = decoder.Decode(&obj.Fees)
	if err != nil {
		return err
	}
	// Deserialize `OutPut`:
	err = decoder.Decode(&obj.OutPut)
	if err != nil {
		return err
	}
	// Deserialize `TokenCoin`:
	err = decoder.Decode(&obj.TokenCoin)
	if err != nil {
		return err
	}
	// Deserialize `TokenPc`:
	err = decoder.Decode(&obj.TokenPc)
	if err != nil {
		return err
	}
	// Deserialize `CoinMint`:
	err = decoder.Decode(&obj.CoinMint)
	if err != nil {
		return err
	}
	// Deserialize `PcMint`:
	err = decoder.Decode(&obj.PcMint)
	if err != nil {
		return err
	}
	// Deserialize `LpMint`:
	err = decoder.Decode(&obj.LpMint)
	if err != nil {
		return err
	}
	// Deserialize `OpenOrders`:
	err = decoder.Decode(&obj.OpenOrders)
	if err != nil {
		return err
	}
	// Deserialize `Market`:
	err = decoder.Decode(&obj.Market)
	if err != nil {
		return err
	}
	// Deserialize `SerumDex`:
	err = decoder.Decode(&obj.SerumDex)
	if err != nil {
		return err
	}
	// Deserialize `TargetOrders`:
	err = decoder.Decode(&obj.TargetOrders)
	if err != nil {
		return err
	}
	// Deserialize `WithdrawQueue`:
	err = decoder.Decode(&obj.WithdrawQueue)
	if err != nil {
		return err
	}
	// Deserialize `TokenTempLp`:
	err = decoder.Decode(&obj.TokenTempLp)
	if err != nil {
		return err
	}
	// Deserialize `AmmOwner`:
	err = decoder.Decode(&obj.AmmOwner)
	if err != nil {
		return err
	}
	// Deserialize `LpAmount`:
	err = decoder.Decode(&obj.LpAmount)
	if err != nil {
		return err
	}
	// Deserialize `ClientOrderId`:
	err = decoder.Decode(&obj.ClientOrderId)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
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
	SwapCoinInAmount    ag_binary.Uint128
	SwapPcOutAmount     ag_binary.Uint128
	SwapTakePcFee       uint64
	SwapPcInAmount      ag_binary.Uint128
	SwapCoinOutAmount   ag_binary.Uint128
	SwapTakeCoinFee     uint64
}

func (obj OutPutData) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
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

func (obj *OutPutData) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
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
