package transformers

import "sort"

// CreateOrderedMixedTradeArray creates an array of MixedExecutedTransform ordered by OrderID in descending order
func CreateOrderedMixedTradeArray(marketExecutedArray []MarketExecutedTransform, limitExecutedArray []LimitExecutedTransform) []MixedExecutedTransform {
	// Combine both arrays into a single array of MixedExecutedTransform
	var mixedArray []MixedExecutedTransform

	// Add MarketExecutedTransform items to mixedArray
	for _, marketExecuted := range marketExecutedArray {
		mixedArray = append(mixedArray, MixedExecutedTransform{
			OrderID:         marketExecuted.OrderID,
			Trade:           marketExecuted.Trade,
			Price:           marketExecuted.Price,
			PriceImpactP:    marketExecuted.PriceImpactP,
			PositionSizeDai: marketExecuted.PositionSizeDai,
			PercentProfit:   marketExecuted.PercentProfit,
			DaiSentToTrader: marketExecuted.DaiSentToTrader,
			Open:            marketExecuted.Open,
		})
	}

	// Add LimitExecutedTransform items to mixedArray
	for _, limitExecuted := range limitExecutedArray {
		mixedArray = append(mixedArray, MixedExecutedTransform{
			OrderID:         limitExecuted.OrderID,
			Trade:           limitExecuted.Trade,
			Price:           limitExecuted.Price,
			PriceImpactP:    limitExecuted.PriceImpactP,
			PositionSizeDai: limitExecuted.PositionSizeDai,
			PercentProfit:   limitExecuted.PercentProfit,
			DaiSentToTrader: limitExecuted.DaiSentToTrader,
			LimitIndex:      limitExecuted.LimitIndex,
			NftHolder:       limitExecuted.NftHolder,
			OrderType:       limitExecuted.OrderType,
			ExactExecution:  limitExecuted.ExactExecution,
		})
	}

	// Sort mixedArray by OrderID in descending order
	sort.Slice(mixedArray, func(i, j int) bool {
		return mixedArray[i].OrderID > mixedArray[j].OrderID
	})

	return mixedArray
}
