package boardbuilders

import (
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func BuildSlEditBoard(tradeCache *tradecache.TradeCache, guid int64) tgbotapi.InlineKeyboardMarkup {
	trade, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}
	var newTradeStr string
	var activeTradesStr string

	// Check on which Window we are:
	if trade.ActiveWindow == 0 { // We are on new trade
		newTradeStr = "⭐️ NEW TRADE"
		activeTradesStr = "ACTIVE TRADES"
	} else {
		if trade.ActiveWindow == 1 {
			newTradeStr = "NEW TRADE"
			activeTradesStr = "⭐️ ACTIVE TRADES"
		} else {
			log.Println("Error during BuildPerpMainBoard: ActiveWindow not equal to 0 or 1")
		}
	}

	//pairStr := pairmaps.IndexToPair[int(user.PairIndex)]

	openPrice := utils.FormatPrice(trade.SL)

	return tgbotapi.NewInlineKeyboardMarkup(
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("⬅️", "/previouspair"),
		//	tgbotapi.NewInlineKeyboardButtonData("🔄 "+pairStr, "/refresh"),
		//	tgbotapi.NewInlineKeyboardButtonData("➡️", "/nextpair"),
		//),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newTradeStr, "/newtrade"),
			tgbotapi.NewInlineKeyboardButtonData(activeTradesStr, "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrtpedit"),
			tgbotapi.NewInlineKeyboardButtonData(openPrice, "/customsledit"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrsledit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("None", "/zerosledit"),
			tgbotapi.NewInlineKeyboardButtonData("-10%", "/minus10edit"),
			tgbotapi.NewInlineKeyboardButtonData("-25%", "/minus25edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("-33%", "/minus33edit"),
			tgbotapi.NewInlineKeyboardButtonData("-50%", "/minus50edit"),
			tgbotapi.NewInlineKeyboardButtonData("-75%", "/minus75edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/activetrades"),
			tgbotapi.NewInlineKeyboardButtonData("💾 Save", "/savesledit"),
		),
	)
}
