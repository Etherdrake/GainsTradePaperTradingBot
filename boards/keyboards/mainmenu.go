package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	MainMenuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{
				Text:            "⚖️ Leverage",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
			tgbotapi.KeyboardButton{
				Text:            "🔀 Swap",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
			tgbotapi.KeyboardButton{
				Text:            "🔫 Sniper",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{
				Text:            "🎁 Airdrops",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
			tgbotapi.KeyboardButton{
				Text:            "📋 Copytrade",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
			tgbotapi.KeyboardButton{
				Text:            "👜 Wallet",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{
				Text:            "⚡ Premium",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
			tgbotapi.KeyboardButton{
				Text:            "⚙️ Settings",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
			tgbotapi.KeyboardButton{
				Text:            "ℹ️ Help",
				RequestContact:  false,
				RequestLocation: false,
				RequestPoll:     nil,
			},
		),
	)
	MainMenuKeyboard.OneTimeKeyboard = false
	MainMenuKeyboard.InputFieldPlaceholder = "Hoot Trade - Executor Class Trading"
	return MainMenuKeyboard
}

var MainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "⚖️ Perps",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "🔀 Swap",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "🔫 Sniper",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "🎁 Airdrops",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "📋 Copytrade",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "👜 Wallet",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "⚡ Premium",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "⚙️ Settings",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "ℹ️ Help",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
)
