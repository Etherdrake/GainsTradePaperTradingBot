package chartservice

import (
	"HootTelegram/boardbuilders"
	"HootTelegram/tradecache"
	"encoding/base64"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"os"
	"time"
)

func SendChartBoardTimeframe(bot *tgbotapi.BotAPI, chatID int64, base64String string, tradeCache *tradecache.TradeCache, rdbPrice *redis.Client, guid int64, timeframe int) error {
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

	var chartKeyboard tgbotapi.InlineKeyboardMarkup
	// Create an inline keyboard
	if timeframe == 1 {
		chartKeyboard = boardbuilders.BuildNewTradeBoardChart1M(tradeCache, rdbPrice, guid)
	}
	if timeframe == 5 {
		chartKeyboard = boardbuilders.BuildNewTradeBoardChart5M(tradeCache, rdbPrice, guid)
	}
	if timeframe == 15 {
		chartKeyboard = boardbuilders.BuildNewTradeBoardChart15M(tradeCache, rdbPrice, guid)
	}
	if timeframe == 60 {
		chartKeyboard = boardbuilders.BuildNewTradeBoardChart1H(tradeCache, rdbPrice, guid)
	}
	if timeframe == 240 {
		chartKeyboard = boardbuilders.BuildNewTradeBoardChart4H(tradeCache, rdbPrice, guid)
	}
	if timeframe == 960 {
		chartKeyboard = boardbuilders.BuildNewTradeBoardChart1D(tradeCache, rdbPrice, guid)
	}

	// Set the inline keyboard in the message
	photoConfig.ReplyMarkup = &chartKeyboard

	// Send the photo with the inline keyboard
	_, err = bot.Send(photoConfig)
	if err != nil {
		return fmt.Errorf("error sending photo: %v", err)
	}

	return nil
}

func SendChartLargeBoardTest(bot *tgbotapi.BotAPI, chatID int64, base64String string, tradeCache *tradecache.TradeCache, rdbPrice *redis.Client, guid int64, timeframe int) error {
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
	inlineKeyboard := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)

	// Set the inline keyboard in the message
	photoConfig.ReplyMarkup = &inlineKeyboard

	// Send the photo with the inline keyboard
	_, err = bot.Send(photoConfig)
	if err != nil {
		return fmt.Errorf("error sending photo: %v", err)
	}

	return nil
}
