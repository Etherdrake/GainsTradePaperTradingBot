package tg_utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// DeleteMessage attempts to delete a message and returns any error encountered
func DeleteMessage(bot *tgbotapi.BotAPI, msgToDelete tgbotapi.DeleteMessageConfig) error {
	_, err := bot.Request(msgToDelete)
	return err
}
