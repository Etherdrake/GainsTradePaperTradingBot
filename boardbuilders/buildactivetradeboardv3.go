package boardbuilders

import (
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func BuildActiveTradeBoardGNSV2(
	rdbPrice *redis.Client,
	rdbPositionsPaper *redis.Client,
	guid int64, orderID string) tgbotapi.InlineKeyboardMarkup {

	newTradeStr := "NEW TRADE"
	activeTradesStr := "⭐️ ACTIVE TRADES"

	trade, err := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), orderID)
	if err != nil {
		log.Println("No trade found for: ", guid)
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("NEW TRADE", "/newtrade"),
				tgbotapi.NewInlineKeyboardButtonData("⭐️ ACTIVE TRADES", "/"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Refresh", "/"),
				tgbotapi.NewInlineKeyboardButtonData("No Trade or Order", "/"),
				tgbotapi.NewInlineKeyboardButtonData("Chart", "/"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Take Profit", "/"),
				tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Share on Telegram", "/"),
				tgbotapi.NewInlineKeyboardButtonData("Share on Twitter", "/"),
			),
			tgbotapi.NewInlineKeyboardRow(
				// The position in question
				tgbotapi.NewInlineKeyboardButtonData("No trade or order", "/")),
		)
	}

	var positionSize float64
	var openPrice float64
	var pairIdx int64
	var leverage int64

	positionSize, _ = strconv.ParseFloat(trade.PositionSizeDai, 64)
	openPrice, _ = strconv.ParseFloat(trade.OpenPrice, 64)
	pairIdx, _ = strconv.ParseInt(trade.PairIndex, 10, 64)
	leverage, _ = strconv.ParseInt(trade.Leverage, 10, 64)

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, int(pairIdx)) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	var getPnlUSD float64

	getPnlUSD = utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, float64(leverage), trade.Buy)

	var pnlStr string
	if getPnlUSD > 0 {
		plusSymbol := "+"
		pnlStr = strconv.FormatFloat(getPnlUSD, 'f', 2, 64)
		pnlStr = plusSymbol + pnlStr
	} else {
		pnlStr = strconv.FormatFloat(getPnlUSD, 'f', 2, 64)
	}

	var closeTradeOrCancelOrder string
	if trade.OrderStatus == 2 {
		closeTradeOrCancelOrder = "Cancel Order"
	} else {
		closeTradeOrCancelOrder = "Close Position: " + "💰" + pnlStr
	}

	// If trade or order /t or /o + ":" + "SYMBOL"
	symbol := pairmaps.IndexToCrypto[int(pairIdx)]

	closeButton := tgbotapi.NewInlineKeyboardButtonData(closeTradeOrCancelOrder, "/placeholder")

	// Change the callbacks on the TP / SL
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newTradeStr, "/newtrade"),
			tgbotapi.NewInlineKeyboardButtonData(activeTradesStr, "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Refresh", "/activetrades"),
			tgbotapi.NewInlineKeyboardButtonData(symbol, "/nextrade"),
			tgbotapi.NewInlineKeyboardButtonData("Chart", "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take Profit", "/tpedit"),
			tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/sledit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Share on Telegram", "/activetrades"),
			tgbotapi.NewInlineKeyboardButtonData("Share on Twitter", "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			// The position in question
			closeButton),
	)
}
