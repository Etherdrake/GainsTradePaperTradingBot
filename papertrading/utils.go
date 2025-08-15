package papertrading

import (
	"HootTelegram/api"
	"HootTelegram/priceserver"
	"HootTelegram/utils"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func ConvertOpenTradeToPnL(openTrade api.OpenTradeJSON, rdbPrice *redis.Client) float64 {
	pairIndex, _ := strconv.Atoi(openTrade.PairIndex)
	posSize, _ := strconv.ParseFloat(openTrade.PositionSizeDai, 64)
	openPrice, _ := strconv.ParseFloat(openTrade.OpenPrice, 64)
	currentPrice, _ := priceserver.GetPrice(rdbPrice, pairIndex)
	leverage, _ := strconv.ParseFloat(openTrade.Leverage, 64)
	openTradePnL := utils.CalculateDollarProfitOrLossWithLeverage(posSize, openPrice, currentPrice, leverage, openTrade.Buy)
	return openTradePnL
}
