package stringbuilders

import (
	"HootTelegram/api/gnsfees"
	"HootTelegram/database"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func BuildNewTradeString(rdbPrice *redis.Client, guid int64, trade tradecache.OpenTradeCache, pairIndex int64) string {
	balanceFloat, err := database.GetPaperBalance(guid)
	if err != nil {
		log.Println(err)
	}

	// Retrieve the price from Redis
	price, err := priceserver.GetPrice(rdbPrice, int(pairIndex)) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	var LongShortString string
	if trade.Buy {
		LongShortString = "Long"
	} else {
		LongShortString = "Short"
	}

	var orderType string
	if trade.OrderType == 0 {
		orderType = "Market"
	}
	if trade.OrderType == 1 {
		orderType = "Limit"
	}
	if trade.OrderType == 2 {
		orderType = "Stop"
	}

	chainStr := "Practice"

	directOrderString := orderType + " " + LongShortString

	formattedPrice := utils.FormatPrice(price)

	//totalPosSize := trade.Leverage * trade.PositionSizeDai

	realPaper := "🟡 Practice Mode [ Switch to REAL TRADING -> @HootTradeBot ]"

	activeCollateral := "Active Collateral: *USDC* [GUIDE -> HOW DO THIS WORK ON REAL GAINS?]"

	positionSize := trade.PositionSizeDai * trade.Leverage

	allFees, err := gnsfees.GetFeesAPI("arbitrum", trade.Buy, trade.PositionSizeDai, trade.Leverage, trade.PairIndex, 2, trade.PositionSizeDai)
	if err != nil {
		log.Println("Error calling GetFeesAPI")
	}
	msg := fmt.Sprintf(`
------------------------------------------------------
%s 
💸 Balance: *$%2.f* /wallet
⛽ Gas Balance: *%s %s* [WHAT IS GAS WORK?]
🎟 XP Points: *%d* /redeem
🪙 %s
------------------------------------------------------
⭐️ *NEW TRADE*

%s: %s @ *$%s*
with *$%d* collateral and *%dx* lever

💰 Total position size: *$%d*
🤸 Spread: *%.4f%%*
⛽ Fees:* $%2.f* | LIQ: *$%s*
TP: *$%s* SL: *$%s* 
🔗 %s | Switch to %s [WHAT IS ARBITRUM?]`,
		realPaper,
		balanceFloat,
		"♾️",
		"USDC",
		0,
		activeCollateral,
		directOrderString,
		utils.EscapeMarkdownV2(pairmaps.IndexToPair[int(pairIndex)]),
		formattedPrice,
		trade.PositionSizeDai,
		trade.Leverage,
		positionSize,      // Notice this change
		allFees.AvgSpread, //allFees.AvgSpread,
		allFees.Fees,      //allFees.Fees,
		utils.FormatPrice(trade.Liq),
		utils.FormatPrice(trade.TP),
		utils.FormatPrice(trade.SL),
		chainStr,
		"arbitrum",
	)
	return msg
}
