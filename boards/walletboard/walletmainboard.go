package walletboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var WalletMainBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", "/backmain"),
	),
)
