package boardbuilders

import (
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func BuildTpEditBoard(tradeCache *tradecache.TradeCache, guid int64) tgbotapi.InlineKeyboardMarkup {
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

	openPrice := utils.FormatPrice(trade.TP)

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
			tgbotapi.NewInlineKeyboardButtonData(openPrice, "/customtpedit"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrtpedit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("None", "/zerotpedit"),
			tgbotapi.NewInlineKeyboardButtonData("25%", "/plus25edit"),
			tgbotapi.NewInlineKeyboardButtonData("50%", "/plus50edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("100%", "/plus100edit"),
			tgbotapi.NewInlineKeyboardButtonData("200%", "/plus150edit"),
			tgbotapi.NewInlineKeyboardButtonData("900%", "/plus900edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/activetrades"),
			tgbotapi.NewInlineKeyboardButtonData("💾 Save", "/savetpedit"),
		),
	)
}
