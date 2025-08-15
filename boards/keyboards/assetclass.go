package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var AssetClasskeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "Crypto",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "Forex",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "Comm",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
)
