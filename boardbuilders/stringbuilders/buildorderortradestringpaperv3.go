package stringbuilders

import (
	"HootTelegram/api"
	"HootTelegram/priceserver"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// BuildOrderOrTradeStringPaperV3 is to be used to
func BuildOrderOrTradeStringPaperV3(
	rdbPrice *redis.Client,
	perp api.OpenTradeJSON,
) string {

	// Retrieve activeOrderOrTrade
	var currentPrice float64
	var closeFee int
	var pairSpreadStr string
	var positionSize float64
	//var liquidation float64
	var err error

	chainStr := "Practice"

	pairIdxStr, _ := strconv.Atoi(perp.PairIndex)

	// Retrieve the price from Redis
	currentPrice, err = priceserver.GetPrice(rdbPrice, pairIdxStr) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in BuildActiveTradeStringV2:", err)
	}

	closeFee = 420
	pairSpread := 69.420

	//pairSpreadFloat := float64(pairSpread)
	pairSpreadStr = fmt.Sprintf("%.2f%%", pairSpread)

	//liquidation = utils.CalculateLiquidationPrice(activePerpTrade.Trade.OpenPrice, float64(activePerpTrade.Trade.Leverage), activePerpTrade.Trade.Buy)

	msg := fmt.Sprintf(`
💲 Current Price: %.5f
💰 Total position size: $%2.5f
🤸 Spread: *%s* | Token: *%s*
TP: *$%s* SL: *$%s*
LIQ: *$%s* | 🎯 Gains Network
⛓️ *%s* | ⛽ Fees: *%d*`,
		currentPrice,
		positionSize, // Notice this change
		pairSpreadStr,
		"USDC",
		perp.TP,
		perp.SL,
		perp.Liq,
		chainStr,
		closeFee,
	)
	return msg
}
