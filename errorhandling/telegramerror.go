package errorhandling

import (
	"HootTelegram/alertatron"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func HandleTelegramError(bot *tgbotapi.BotAPI, err error, chatID string) {
	// Check for the specific error
	if strings.Contains(err.Error(), "must be of type Array") {
		fmt.Println("Error: Invalid inline keyboard format. Must be an array.")
		alertatron.SendAlert(bot, chatID, "Please press /start again!")

		// Handle this specific error condition
	} else {
		// Handle other Telegram API errors
		fmt.Printf("Generic Error: %s\n", err.Error())
	}
}
