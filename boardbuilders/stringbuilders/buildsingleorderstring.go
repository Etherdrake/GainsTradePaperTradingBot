package stringbuilders

import (
	"HootTelegram/api"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func BuildSingleOrderString(rdbPrice *redis.Client, trade api.OpenTradeJSON, tradeIndex int) (string, error) {
	openPrice, err := strconv.ParseFloat(trade.OpenPrice, 64)
	if err != nil {
		return "", err
	}

	//positionSizeDai, err := strconv.ParseFloat(trade.PositionSizeDai, 64)
	//if err != nil {
	//	return "", err
	//}

	//positionSize, err := strconv.ParseFloat(trade.PositionSizeDai, 64)
	openPrice, err = strconv.ParseFloat(trade.OpenPrice, 64)
	convert, _ := strconv.Atoi(trade.PairIndex)

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, convert) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	//leverage, err := strconv.ParseFloat(trade.Leverage, 64)

	tradeDirection := "📉"
	if trade.Buy {
		tradeDirection = "📈"
	}

	tradeDirectionStr := "Short"
	if trade.Buy {
		tradeDirectionStr = "Long"
	}

	tickerInt, _ := strconv.Atoi(trade.PairIndex)

	tradeIndexIncr := tradeIndex + 1

	positionString := fmt.Sprintf("/o*%d* %s %s %s @ *$%s*, now: $%s",
		tradeIndexIncr,
		pairmaps.IndexToPair[tickerInt],
		tradeDirection,
		tradeDirectionStr,
		utils.FormatPrice(openPrice),
		utils.FormatPrice(currentPrice),
	)
	return positionString, nil
}
