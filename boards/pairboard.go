package boards

import (
	"HootTelegram/pairmaps"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var PairsBoardOne = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Hoot Bot", "/"),
		tgbotapi.NewInlineKeyboardButtonData("Trading Pairs", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[0], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[1], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[2], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[3], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[4], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[5], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[6], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pairmaps.IndexToPair[7], "/"),
		tgbotapi.NewInlineKeyboardButtonData("$", "/"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("Page 1", "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	),
)
