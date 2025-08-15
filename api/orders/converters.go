package orders

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

func ConvertToOpenOrdersBig(openOrdersString OpenOrdersString) OpenOrdersBig {
	openOrdersBig := OpenOrdersBig{
		Trader:           openOrdersString.Trader,
		PairIndex:        ConvertToBigInt(openOrdersString.PairIndex),
		Index:            ConvertToBigInt(openOrdersString.Index),
		PositionSize:     ConvertToBigInt(openOrdersString.PositionSize),
		SpreadReductionP: ConvertToBigInt(openOrdersString.SpreadReductionP),
		Buy:              openOrdersString.Buy,
		Leverage:         ConvertToBigInt(openOrdersString.Leverage),
		Tp:               ConvertToBigInt(openOrdersString.Tp),
		Sl:               ConvertToBigInt(openOrdersString.Sl),
		MinPrice:         ConvertToBigInt(openOrdersString.MinPrice),
		MaxPrice:         ConvertToBigInt(openOrdersString.MaxPrice),
		Block:            ConvertToBigInt(openOrdersString.Block),
		TokenID:          ConvertToBigInt(openOrdersString.TokenID),
	}

	return openOrdersBig
}

func ConvertToBigInt(value string) *big.Int {
	value = strings.TrimPrefix(value, "0x")

	if len(value)%2 != 0 {
		value = "0" + value
	}

	decoded, err := hex.DecodeString(value)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return nil
	}

	result := new(big.Int).SetBytes(decoded)
	return result
}

func ConvertOpenOrdersBigToOpenOrdersHuman(openOrdersBig OpenOrdersBig) OpenOrdersHuman {
	eighteenFloat := big.NewFloat(1e18)
	tenFloat := big.NewFloat(1e10)

	positionSizeDai := new(big.Float).SetInt(openOrdersBig.PositionSize)
	positionSizeDai.Quo(positionSizeDai, eighteenFloat)
	positionSizeDaiF64, _ := positionSizeDai.Float64()
	spreadReductionP := new(big.Float).SetInt(openOrdersBig.SpreadReductionP)
	spreadReductionP.Quo(spreadReductionP, tenFloat)
	spreadReductionPFloat64, _ := spreadReductionP.Float64()
	minPrice := new(big.Float).SetInt(openOrdersBig.MinPrice)
	minPrice.Quo(minPrice, tenFloat)
	minPriceFloat64, _ := minPrice.Float64()
	maxPrice := new(big.Float).SetInt(openOrdersBig.MaxPrice)
	maxPrice.Quo(maxPrice, tenFloat)
	maxPriceFloat64, _ := maxPrice.Float64()
	tp := new(big.Float).SetInt(openOrdersBig.Tp)
	tp.Quo(tp, tenFloat)
	tpFloat64, _ := tp.Float64()
	sl := new(big.Float).SetInt(openOrdersBig.Sl)
	sl.Quo(sl, tenFloat)
	slFloat64, _ := sl.Float64()

	convert := OpenOrdersHuman{
		Trader:           openOrdersBig.Trader,
		PairIndex:        openOrdersBig.PairIndex.Int64(),
		Index:            openOrdersBig.Index.Int64(),
		PositionSize:     positionSizeDaiF64,
		SpreadReductionP: spreadReductionPFloat64,
		Buy:              false,
		Leverage:         openOrdersBig.Leverage.Int64(),
		Tp:               tpFloat64,
		Sl:               slFloat64,
		MinPrice:         minPriceFloat64,
		MaxPrice:         maxPriceFloat64,
		Block:            openOrdersBig.Block.Int64(),
		TokenID:          0,
	}
	return convert
}

func ConvertReceivedArrToTotalArr(openOrdersArr []OpenOrdersString) []OpenOrdersHuman {
	var convArr []OpenOrdersHuman
	convArr = make([]OpenOrdersHuman, len(openOrdersArr))
	for i := 0; i < len(openOrdersArr); i++ {
		convArr[i] = ConvertOpenOrdersBigToOpenOrdersHuman(ConvertToOpenOrdersBig(openOrdersArr[i]))
	}
	return convArr
}
