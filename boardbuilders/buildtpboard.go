package boardbuilders

import (
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func BuildTpBoard(tradeCache *tradecache.TradeCache, guid int64) tgbotapi.InlineKeyboardMarkup {
	user, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

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

	var longShortButton tgbotapi.InlineKeyboardButton
	if user.Buy {
		longShortButton = tgbotapi.NewInlineKeyboardButtonData("🟢 LONG", "/toggletoshort")
	} else {
		longShortButton = tgbotapi.NewInlineKeyboardButtonData("🔴 SHORT", "/toggletolong")
	}

	takeProfit := utils.FormatPrice(user.TP)

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
			tgbotapi.NewInlineKeyboardButtonData("✅ Take Profit", "/perps"),
			longShortButton,
			tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/stoploss"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrtp"),
			tgbotapi.NewInlineKeyboardButtonData(takeProfit, "/customtp"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrtp"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("None", "/zerotp"),
			tgbotapi.NewInlineKeyboardButtonData("25%", "/plus25"),
			tgbotapi.NewInlineKeyboardButtonData("50%", "/plus50"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("100%", "/plus100"),
			tgbotapi.NewInlineKeyboardButtonData("200%", "/plus150"),
			tgbotapi.NewInlineKeyboardButtonData("900%", "/plus900"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Save", "/leverage"),
		),
	)
}
