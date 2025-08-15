package chartservice

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"os"
	"time"
)

func SendChartWithKeyboard(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
	// Decode base64 string
	baseExtract, err := extractBase64(base64String)
	if err != nil {
		return fmt.Errorf("error decoding base64: %v", err)
	}
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

	// Create a PhotoConfig to send the photo
	photoConfig := tgbotapi.NewPhoto(chatID, fileBytes)

	// Create an inline keyboard
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5M", "button1"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "button2"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "button2"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "button2"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "button2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "/leverageback"),
		),
	)

	// Set the inline keyboard in the message
	photoConfig.ReplyMarkup = &inlineKeyboard

	// Send the photo with the inline keyboard
	_, err = bot.Send(photoConfig)
	if err != nil {
		return fmt.Errorf("error sending photo: %v", err)
	}

	return nil
}

func SendChartWithKeyboardReturn(bot *tgbotapi.BotAPI, chatID int64, base64String string) (*tgbotapi.PhotoConfig, tgbotapi.InlineKeyboardMarkup, error) {
	// Decode base64 string
	baseExtract, err := extractBase64(base64String)
	if err != nil {
		return nil, tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("error decoding base64: %v", err)
	}
	imageBytes, err := base64.StdEncoding.DecodeString(baseExtract)
	if err != nil {
		return nil, tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("error decoding base64: %v", err)
	}

	// Generate a unique file name based on the current timestamp
	fileName := fmt.Sprintf("photo_%d.jpg", time.Now().UnixNano())

	// Create a temporary file to store the decoded image
	tmpfile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return nil, tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("error creating temporary file: %v", err)
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
		return nil, tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("error writing to temporary file: %v", err)
	}

	// Seek to the beginning of the file for reading
	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return nil, tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("error seeking to the beginning of the file: %v", err)
	}

	// Create a FileBytes from the temporary file
	fileBytes := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: imageBytes,
	}

	// Create a PhotoConfig to send the initial photo
	photoConfig := tgbotapi.NewPhoto(chatID, fileBytes)

	// Create an inline keyboard
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5M", "button1"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "button2"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "button2"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "button2"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "button2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<- Back", "/newtrade"),
		),
	)

	// Set the inline keyboard in the message config
	photoConfig.ReplyMarkup = &inlineKeyboard

	return &photoConfig, inlineKeyboard, nil
}

func EditMessageWithChartAndKeyboard(bot *tgbotapi.BotAPI, chatID int64, messageID int, base64String string, updatedInlineKeyboard tgbotapi.InlineKeyboardMarkup) error {
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

	//// Create a PhotoConfig to edit the existing message
	//photoConfig := tgbotapi.NewPhoto(chatID, fileBytes)
	//photoConfig.ReplyMarkup = &updatedInlineKeyboard
	//photoConfig.ChatID = chatID

	// Create a FileBytes from the temporary file
	fileBytes := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: imageBytes,
	}

	baseInputMedia := tgbotapi.BaseInputMedia{
		Type:      "photo", // Set the desired media type
		Media:     fileBytes,
		ParseMode: "markdown", // Set the desired parse mode
	}

	// Create an EditMessageMediaConfig to update the message
	editMessageConfig := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: &updatedInlineKeyboard,
		},
		Media: tgbotapi.InputMediaPhoto{
			BaseInputMedia: baseInputMedia,
		},
	}

	// Edit the existing message with the updated photo and inline keyboard
	_, err = bot.Send(&editMessageConfig)
	if err != nil {
		return fmt.Errorf("error editing message: %v", err)
	}
	return nil
}
