package stringbuilders

import (
	"HootTelegram/api"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/transformers"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func BuildSinglePositionString(rdbPrice *redis.Client, trade api.OpenTradeJSON, tradeIndex int) (string, error) {
	//positionSizeDai, err := strconv.ParseFloat(trade.PositionSizeDai, 64)
	//if err != nil {
	//	return "", err
	//}

	positionSize, err := strconv.ParseFloat(trade.PositionSizeDai, 64)
	openPrice, err := strconv.ParseFloat(trade.OpenPrice, 64)
	pairIdx, _ := strconv.Atoi(trade.PairIndex)
	leverage, err := strconv.ParseFloat(trade.Leverage, 64)

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, pairIdx) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	getPnlUSD := utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, leverage, trade.Buy)
	getPnlPrcnt := utils.CalculatePercentageProfitOrLoss(positionSize, getPnlUSD, trade.Buy)

	tradeDirection := "📉"
	if trade.Buy {
		tradeDirection = "📈"
	}

	var plusOrMinus string
	var theDot string
	if getPnlUSD > 0 {
		plusOrMinus = "+"
		theDot = "🟢"
	} else {
		plusOrMinus = ""
		theDot = "🔴"
	}

	tickerInt, _ := strconv.Atoi(trade.PairIndex)

	tradeIndexIncr := tradeIndex + 1

	positionString := fmt.Sprintf("/t%d %s %s %s%.2f%% %s %s*$%.2f*",
		tradeIndexIncr,
		pairmaps.IndexToPair[tickerInt],
		tradeDirection,
		plusOrMinus,
		getPnlPrcnt,
		theDot,
		plusOrMinus,
		getPnlUSD,
	)
	return positionString, nil
}

func BuildSinglePositionStringGNS(rdbPrice *redis.Client, trade transformers.MixedExecutedTransform, tradeIndex int) (string, error) {
	positionSize := trade.PositionSizeDai
	openPrice := trade.Trade.OpenPrice
	pairIndex := trade.Trade.PairIndex
	leverage := trade.Trade.Leverage

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, int(pairIndex)) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	getPnlUSD := utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, float64(leverage), trade.Trade.Buy)
	getPnlPrcnt := utils.CalculatePercentageProfitOrLoss(positionSize, getPnlUSD, trade.Trade.Buy)

	tradeDirection := "📉"
	if trade.Trade.Buy {
		tradeDirection = "📈"
	}

	var plusOrMinus string
	var theDot string
	if getPnlUSD > 0 {
		plusOrMinus = "+"
		theDot = "🟢"
	} else {
		plusOrMinus = ""
		theDot = "🔴"
	}

	tickerInt := trade.Trade.PairIndex

	tradeIndexIncr := tradeIndex + 1

	positionString := fmt.Sprintf("/x%d %s %s %s%.2f%% %s %s*$%.2f*",
		tradeIndexIncr,
		pairmaps.IndexToPair[int(tickerInt)],
		tradeDirection,
		plusOrMinus,
		getPnlPrcnt,
		theDot,
		plusOrMinus,
		getPnlUSD,
	)
	return positionString, nil
}

func BuildSingleTradeStringGNS(rdbPrice *redis.Client, trade transformers.MixedExecutedTransform, tradeIndex int) (string, error) {
	convert := int(trade.Trade.PairIndex)
	//positionSizeDai, err := strconv.ParseFloat(trade.PositionSizeDai, 64)
	//if err != nil {
	//	return "", err
	//}

	positionSize := trade.PositionSizeDai

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, convert) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	getPnlUSD := utils.CalculateDollarProfitOrLossWithLeverage(trade.PositionSizeDai, trade.Trade.OpenPrice, currentPrice, float64(trade.Trade.Leverage), trade.Trade.Buy)
	getPnlPrcnt := utils.CalculatePercentageProfitOrLoss(positionSize, getPnlUSD, trade.Trade.Buy)

	tradeDirection := "📉"
	if trade.Trade.Buy {
		tradeDirection = "📈"
	}

	var plusOrMinus string
	var theDot string
	if getPnlUSD > 0 {
		plusOrMinus = "+"
		theDot = "🟢"
	} else {
		plusOrMinus = ""
		theDot = "🔴"
	}

	tickerInt, _ := strconv.Atoi(strconv.Itoa(int(trade.Trade.PairIndex)))

	tradeIndexIncr := tradeIndex + 1

	positionString := fmt.Sprintf("/t%d %s %s %s%.2f%% %s %s*$%.2f*",
		tradeIndexIncr,
		pairmaps.IndexToPair[tickerInt],
		tradeDirection,
		plusOrMinus,
		getPnlPrcnt,
		theDot,
		plusOrMinus,
		getPnlUSD,
	)
	return positionString, nil
}
