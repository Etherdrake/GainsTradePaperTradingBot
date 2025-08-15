package boards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var OpenPosition = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("🔄 getPosition", "/refreshposition"), // FIND ACTIVE PAIR FOR USER USING CACHE - DEFAULT = BTC/USD
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("closePosition(guid, activePosition)", "/closeposition"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ TP", "/custom_tp_open_position"),
		tgbotapi.NewInlineKeyboardButtonData("🔴 SL", "/custom_sl_open_position"),
		tgbotapi.NewInlineKeyboardButtonData("📈 Charts", "/charts"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("NEW TRADE", "/newtrade"),
		tgbotapi.NewInlineKeyboardButtonData("MAIN MENU", "/mainmenu"),
	),
)
