package boards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Define updated inline keyboard for Option 1
var TradeBoardInit = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("", "/"),
		tgbotapi.NewInlineKeyboardButtonData("Hoot Trade - Trading Bot", "/"),
		tgbotapi.NewInlineKeyboardButtonData("", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("BTC/USD", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Price: 29841", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("Charts", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("SHORT", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("LONG", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("-", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("DAI 1000", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("+", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("-", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("x10", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("+", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("TP: 100% @ 59682 ", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("SL: -10% @ 29542", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("LIQ: 16413", "/action2"),
	),
)

// Define updated inline keyboard for Option 2
var UpdatedOption2Board = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Action 3", "/action3"),
		tgbotapi.NewInlineKeyboardButtonData("Action 4", "/action4"),
	),
)

var TradeBoardSam = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("BTC/USD", "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Close Position 💰 +328.43", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Take Profit", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("🔄 Stoploss", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("📈 Charts", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("$28400", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("0%", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("-", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("+", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("-50%", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5M +0.21%", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("15M 1.4%", "/action2"),
		tgbotapi.NewInlineKeyboardButtonData("30M 1.84%", "/action2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ SAVE & RETURN", "/action2"),
	),
)
