package chartservice

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"os"
	"time"
)

func SendChartWithKeyboard1M(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
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
			tgbotapi.NewInlineKeyboardButtonData("5M", "/sendchart5M"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "/sendchart15M"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "/sendchart1H"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "/sendchart4H"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "/sendchart1D"),
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

func SendChartWithKeyboard5M(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
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
			tgbotapi.NewInlineKeyboardButtonData("1M", "/sendchart"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "/sendchart15M"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "/sendchart1H"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "/sendchart4H"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "/sendchart1D"),
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

func SendChartWithKeyboard15M(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
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
			tgbotapi.NewInlineKeyboardButtonData("1M", "/sendchart"),
			tgbotapi.NewInlineKeyboardButtonData("5M", "/sendchart5M"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "/sendchart1H"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "/sendchart4H"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "/sendchart1D"),
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

func SendChartWithKeyboard1H(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
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
			tgbotapi.NewInlineKeyboardButtonData("1M", "/sendchart"),
			tgbotapi.NewInlineKeyboardButtonData("5M", "/sendchart5M"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "/sendchart15M"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "/sendchart4H"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "/sendchart1D"),
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

func SendChartWithKeyboard4H(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
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
			tgbotapi.NewInlineKeyboardButtonData("1M", "/sendchart"),
			tgbotapi.NewInlineKeyboardButtonData("5M", "/sendchart5M"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "/sendchart15M"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "/sendchart1H"),
			tgbotapi.NewInlineKeyboardButtonData("1D", "/sendchart1D"),
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

func SendChartWithKeyboard1D(bot *tgbotapi.BotAPI, chatID int64, base64String string) error {
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
			tgbotapi.NewInlineKeyboardButtonData("1M", "/sendchart"),
			tgbotapi.NewInlineKeyboardButtonData("5M", "/sendchart5M"),
			tgbotapi.NewInlineKeyboardButtonData("15M", "/sendchart15M"),
			tgbotapi.NewInlineKeyboardButtonData("1H", "/sendchart1H"),
			tgbotapi.NewInlineKeyboardButtonData("4H", "/sendchart4H"),
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
