package perps

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var PerpMainBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("🔄 getPair", "/"), // FIND ACTIVE PAIR FOR USER USING CACHE - DEFAULT = BTC/USD
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🟢 getLongShort(guid)", "/"),
		tgbotapi.NewInlineKeyboardButtonData("🔴 getLongShort(guid)", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("getPrice(getPair(guid))", "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("getPositionSize(guid)", "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("getLeverage(guid)", "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back", "/backmain"),
	),
)
