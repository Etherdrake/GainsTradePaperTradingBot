package stringbuilders

import (
	"HootTelegram/pairmaps"
	"HootTelegram/utils"
	"fmt"
	"strings"
)

func BuildCryptoPairString(page int) string {
	header := "📌 **All Crypto Pairs**\n\nSelect Asset To Trade (24hr)\n\n"
	var pairsList string

	pageArray := pairmaps.FilteredIndexToCryptoPage(page, 10) // -> THIS IS NOT UPDATED in COMMIT 100644

	for _, pair := range pageArray {
		pairReturn := strings.Replace(pair, "/", "", -1)
		chartURL := utils.GetChartURLGains(pair)
		name := pairmaps.SymbolToName[pairmaps.CryptoToIndex[pair]]
		pairsList += fmt.Sprintf("/%s: %s [📊 Charts](%s)\n", pairReturn, name, chartURL)
	}
	return header + pairsList
}

func BuildFxPairString(page int) string {
	header := "📌 **All Forex Pairs**\n\nSelect Asset To Trade (24hr)\n\n"
	var pairsList string

	pageArray := pairmaps.FilteredIndexToFxPage(page, 10) // -> THIS IS NOT UPDATED in COMMIT 100644

	for _, pair := range pageArray {
		pairReturn := strings.Replace(pair, "/", "", -1)
		chartURL := utils.GetChartURLGains(pair)
		name := pairmaps.FxToFlag[pair]
		pairsList += fmt.Sprintf("/%s: %s [📊 Charts](%s)\n", pairReturn, name, chartURL)
	}
	return header + pairsList
}

func BuildCommoditiesPairString() string {
	header := "📌 **All Commodity Pairs**\n\nSelect Asset To Trade (24hr)\n\n"
	var pairsList string

	// Get the Commodities array
	commoditiesArray := pairmaps.IndexToCommodity // FIX THIS -> WAS AN ARRAY

	// Construct the list
	for _, pair := range commoditiesArray {
		if pair != "" {
			pairReturn := strings.Replace(pair, "/", "", -1)
			chartURL := utils.GetChartURLGains(pair)
			pairsList += fmt.Sprintf("/%s: %s [📊 Charts](%s)\n", pairReturn, pair, chartURL)
		}
	}

	return header + pairsList
}
