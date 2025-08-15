package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var LeverageKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "2x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "5x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "10x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "20x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "30x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "40x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "50x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "75x",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "MAX",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
)
