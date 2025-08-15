package trends

import "HootTelegram/types"

func IndexToTrend(cache types.PriceCache) map[int]string {
	trendMap := make(map[int]string)
	for i := 0; i < len(cache.IndexToPriceDataClose); i++ {
		if cache.IndexToPriceDataOpen[i] < cache.IndexToPriceDataClose[i] {
			trendMap[i] = "🔼"
		} else {
			trendMap[i] = "🔽"
		}
	}

	return trendMap
}
