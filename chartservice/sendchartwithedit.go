package chartservice

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"os"
	"time"
)

func EditMessageWithChartAndKeyboardDelete(bot *tgbotapi.BotAPI, chatID int64, messageID int, base64String string, updatedInlineKeyboard tgbotapi.InlineKeyboardMarkup) error {
	// Decode base64 string
	baseExtract, err := extractBase64(base64String)
	if err != nil {
		return fmt.Errorf("error decoding base64: %v", err)
	}

	// Decode base64 to image bytes
	imageBytes, err := base64.StdEncoding.DecodeString(baseExtract)
	if err != nil {
		return fmt.Errorf("error decoding base64: %v", err)
	}

	// Generate a unique file name based on the current timestamp
	fileName := fmt.Sprintf("photo_%d.jpg", time.Now().UnixNano())

	// Create a temporary file to store the decoded image
	tmpfile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return fmt.Errorf("error creating temporary file: %v", err)
	}
	defer func() {
		// Close the file before removing it
		tmpfile.Close()
		// Remove the temporary file
		os.Remove(tmpfile.Name())
	}()

	// Write the image data to the temporary file
	_, err = tmpfile.Write(imageBytes)
	if err != nil {
		return fmt.Errorf("error writing to temporary file: %v", err)
	}

	// Seek to the beginning of the file for reading
	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking to the beginning of the file: %v", err)
	}

	// Create a FileBytes from the temporary file
	fileBytes := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: imageBytes,
	}

	// Send a new media message
	mediaMessage := tgbotapi.NewPhoto(chatID, fileBytes)
	_, err = bot.Send(mediaMessage)
	if err != nil {
		return fmt.Errorf("error sending media message: %v", err)
	}

	// Update the existing text message to reference the media message
	editMessageTextConfig := tgbotapi.NewEditMessageText(chatID, messageID, "Updated message text")
	editMessageTextConfig.ReplyMarkup = &updatedInlineKeyboard

	_, err = bot.Send(editMessageTextConfig)
	if err != nil {
		return fmt.Errorf("error editing text message: %v", err)
	}

	return nil
}
