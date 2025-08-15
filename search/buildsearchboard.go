package search

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BuildSearchBoard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 BACK", "/pairs"),
		),
	)
}
