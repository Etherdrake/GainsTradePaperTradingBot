package utils

import (
	"math"
)

func CalculateDollarProfitOrLossWithLeverage(initialInvestment, entryPrice, currentPrice, leverage float64, isLong bool) float64 {
	percentageDifference := (currentPrice - entryPrice) / entryPrice
	if isLong {
		return initialInvestment * percentageDifference * leverage
	}
	return (initialInvestment * percentageDifference * leverage) * -1
}

func CalculatePercentageProfitOrLossWithLeverage(entryPrice, currentPrice, leverage float64, isLong bool) float64 {
	percentageDifference := (currentPrice - entryPrice) / entryPrice
	if isLong {
		return percentageDifference * leverage * 100
	}
	return (percentageDifference * leverage * 100) * -1
}

func CalculatePercentageProfitOrLoss(initialInvestment, profitOrLoss float64, isLong bool) float64 {
	if isLong {
		return (profitOrLoss / initialInvestment) * 100
	}
	return ((profitOrLoss / initialInvestment) * 100) * -1
}

func CalculateLiquidationPrice(openPrice float64, leverage float64, isLong bool) float64 {
	directionMultiplier := 1.0
	if !isLong {
		directionMultiplier = -1.0
	}

	liquidationPrice := openPrice * (1 - directionMultiplier/float64(leverage))
	return liquidationPrice
}

func CalculatePayout(entryPrice, takeProfit, positionSize, leverage float64, isLong bool) float64 {
	directionMultiplier := 1.0
	if !isLong {
		directionMultiplier = -1.0
	}

	priceDifferencePercent := math.Abs((entryPrice-takeProfit)/entryPrice) * 100
	factor := priceDifferencePercent * directionMultiplier
	payout := (factor / 100) * positionSize
	return payout * leverage
}
