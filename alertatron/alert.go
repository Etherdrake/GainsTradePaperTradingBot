package alertatron

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func SendAlert(bot *tgbotapi.BotAPI, chatID string, answerText string) {
	// Create an answer to the callback query
	answer := tgbotapi.NewCallbackWithAlert(chatID, answerText)

	// Set show_alert to t	rue to display it as an alert
	answer.ShowAlert = true

	// Send the answer to the callback query
	_, err := bot.Send(answer)
	if err != nil {
		log.Println(err)
	}
}
