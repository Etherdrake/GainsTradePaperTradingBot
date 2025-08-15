package keyboards

import (
	"HootTelegram/pairmaps"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func GeneratePairButton(index int) tgbotapi.KeyboardButton {
	return tgbotapi.KeyboardButton{
		Text:            pairmaps.IndexToCrypto[index],
		RequestContact:  false,
		RequestLocation: false,
		RequestPoll:     nil,
	}
}

func GeneratePairButtonExUSD(index int) tgbotapi.KeyboardButton {
	return tgbotapi.KeyboardButton{
		Text:            strings.Trim(pairmaps.IndexToCrypto[index], "/USD"),
		RequestContact:  false,
		RequestLocation: false,
		RequestPoll:     nil,
	}
}

func GeneratePairKeyboardCryptoExUSDThreeRowFiveCol(page int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup
	var rows []tgbotapi.KeyboardButton

	// Initialize top row with "CRYPTO", "FX", and "GOLD"
	topRow := []tgbotapi.KeyboardButton{
		{Text: "Crypto"},
		{Text: "Forex"},
		{Text: "Comms"},
	}
	keyboard.Keyboard = append(keyboard.Keyboard, topRow)

	// Calculate the number of pairs to generate
	// As there are 2 rows remaining with 5 buttons each, you need to generate 10 buttons.
	numPairs := 10
	start := page * numPairs

	for i := 0; i < numPairs; i++ {
		button := GeneratePairButtonExUSD(start + i)
		rows = append(rows, button)

		if len(rows) == 5 {
			// Add the row to the keyboard
			keyboard.Keyboard = append(keyboard.Keyboard, rows)

			// Clear the rows
			rows = []tgbotapi.KeyboardButton{}
		}
	}

	// Ensure that the third row and the last button exist
	if len(keyboard.Keyboard) >= 3 && len(keyboard.Keyboard[2]) >= 5 {
		keyboard.Keyboard[2][4].Text = "➡️"
	}

	return keyboard
}

func GeneratePairKeyboardCryptoExUSDThree(page int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup
	var rows []tgbotapi.KeyboardButton

	// Initialize top row with "CRYPTO", "FX", and "GOLD"
	topRow := []tgbotapi.KeyboardButton{
		{Text: "CRYPTO"},
		{Text: "FX"},
		{Text: "GOLD"},
	}
	keyboard.Keyboard = append(keyboard.Keyboard, topRow)

	// Calculate the number of pairs to generate
	// As there are 4 rows remaining with 3 buttons each, you need to generate 12 buttons.
	numPairs := 12
	start := page * numPairs

	for i := 0; i < numPairs; i++ {
		button := GeneratePairButtonExUSD(start + i)
		rows = append(rows, button)

		if len(rows) == 3 {
			// Add the row to the keyboard
			keyboard.Keyboard = append(keyboard.Keyboard, rows)

			// Clear the rows
			rows = []tgbotapi.KeyboardButton{}
		}
	}

	// Ensure that the last row and the last button exist
	if len(keyboard.Keyboard) >= 5 && len(keyboard.Keyboard[4]) >= 3 {
		keyboard.Keyboard[4][2].Text = "➡️"
	}

	return keyboard
}

func GeneratePairKeyboardCrypto(page int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup
	var rows []tgbotapi.KeyboardButton

	// Initialize top row with "CRYPTO", "FX", and "GOLD"
	topRow := []tgbotapi.KeyboardButton{
		{Text: "CRYPTO"},
		{Text: "FX"},
		{Text: "GOLD"},
		// Fill the top row with 5 buttons.
		// You can adjust these empty ones as needed.
		{Text: ""},
		{Text: ""},
	}
	keyboard.Keyboard = append(keyboard.Keyboard, topRow)

	// Calculate the number of pairs to generate
	// As there are 4 rows remaining with 5 buttons each, you need to generate 20 buttons.
	numPairs := 20
	start := page * numPairs

	for i := 0; i < numPairs; i++ {
		button := GeneratePairButton(start + i)
		rows = append(rows, button)

		if len(rows) == 5 {
			// Add the row to the keyboard
			keyboard.Keyboard = append(keyboard.Keyboard, rows)

			// Clear the rows
			rows = []tgbotapi.KeyboardButton{}
		}
	}

	// Ensure that the last row and the last button exist
	if len(keyboard.Keyboard) >= 5 && len(keyboard.Keyboard[4]) >= 5 {
		keyboard.Keyboard[4][4].Text = "➡️"
	}

	return keyboard
}

func GeneratePairKeyboardCryptoExUSD(page int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard tgbotapi.ReplyKeyboardMarkup
	var rows []tgbotapi.KeyboardButton

	// Initialize top row with "CRYPTO", "FX", and "GOLD"
	topRow := []tgbotapi.KeyboardButton{
		{Text: "CRYPTO"},
		{Text: "FX"},
		{Text: "GOLD"},
		// Fill the top row with 5 buttons.
		// You can adjust these empty ones as needed.
		{Text: ""},
		{Text: ""},
	}
	keyboard.Keyboard = append(keyboard.Keyboard, topRow)

	// Calculate the number of pairs to generate
	// As there are 4 rows remaining with 5 buttons each, you need to generate 20 buttons.
	numPairs := 20
	start := page * numPairs

	for i := 0; i < numPairs; i++ {
		button := GeneratePairButtonExUSD(start + i)
		rows = append(rows, button)

		if len(rows) == 5 {
			// Add the row to the keyboard
			keyboard.Keyboard = append(keyboard.Keyboard, rows)

			// Clear the rows
			rows = []tgbotapi.KeyboardButton{}
		}
	}

	// Ensure that the last row and the last button exist
	if len(keyboard.Keyboard) >= 5 && len(keyboard.Keyboard[4]) >= 5 {
		keyboard.Keyboard[4][4].Text = "➡️"
	}

	return keyboard
}

var PairsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.KeyboardButton{
			Text:            pairmaps.IndexToPair[0],
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
		tgbotapi.KeyboardButton{
			Text:            "🎁 Airdrops",
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
		tgbotapi.KeyboardButton{
			Text:            "🎁 Airdrops",
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
		tgbotapi.KeyboardButton{
			Text:            "🎁 Airdrops",
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
		},
	),
)
