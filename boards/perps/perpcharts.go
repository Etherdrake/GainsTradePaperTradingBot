package perps

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var perpChartsBoards = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("🔄 getPair(guid)", "/"), // FIND ACTIVE PAIR FOR USER USING CACHE - DEFAULT = BTC/USD
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🟢 getLongShort(guid)", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("getPrice(getPair(guid))", "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1M", "/"),
		tgbotapi.NewInlineKeyboardButtonData("5M", "/"),
		tgbotapi.NewInlineKeyboardButtonData("15M", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("30M️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("1H", "/"),
		tgbotapi.NewInlineKeyboardButtonData("4H", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("<- RETURN", "/"),
	),
)
