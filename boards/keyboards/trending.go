package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var TrendingKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "Gainers",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "1M",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "Losers",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "5M",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "15M",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "30M",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            "1H",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "3H",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
		tgbotapi.KeyboardButton{
			Text:            "1D",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
)
