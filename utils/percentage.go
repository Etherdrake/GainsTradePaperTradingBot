package utils

import (
	"math/big"
)

func CalculateTakeProfit(price float64, leverage float64, percentage float64) float64 {
	rPrice := new(big.Float).SetFloat64(price)
	rLeverage := new(big.Float).SetFloat64(leverage)
	rPercentage := new(big.Float).SetFloat64(percentage)

	rHundred := big.NewFloat(100)
	rResult := new(big.Float).Mul(rPrice, rPercentage)
	rResult.Quo(rResult, rHundred)
	rResult.Quo(rResult, rLeverage)
	rResult.Add(rResult, rPrice)

	result, _ := rResult.Float64()
	return result
}

func CalculateStopLoss(price float64, leverage float64, percentage float64) float64 {
	rPrice := new(big.Float).SetFloat64(price)
	rLeverage := new(big.Float).SetFloat64(leverage)
	rPercentage := new(big.Float).SetFloat64(percentage)

	rHundred := big.NewFloat(100)
	rResult := new(big.Float).Mul(rPrice, rPercentage)
	rResult.Quo(rResult, rHundred)
	rResult.Quo(rResult, rLeverage)
	rResult.Sub(rPrice, rResult)

	result, _ := rResult.Float64()
	return result
}
