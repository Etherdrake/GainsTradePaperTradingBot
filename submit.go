package main

import (
	"HootTelegram/api"
	"HootTelegram/pairmaps"
	"HootTelegram/tradecache"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func handleSubmit(bot *tgbotapi.BotAPI, tradeCache *tradecache.TradeCache, chatID int64, guid int64) (string, api.OpenTradeJSON, error) {
	// Fetch the trade from the global tradeCache
	trade, exists := tradeCache.Get(guid)
	if !exists {
		return "Cache Error", api.OpenTradeJSON{}, errors.New("trade not found in cache")
	}

	buyOrSell := ""
	if trade.Buy {
		buyOrSell = "BUY"
	} else {
		buyOrSell = "SELL"
	}

	inProgress := "Executing " + "*" + buyOrSell + "*" + " for " + "*" + pairmaps.IndexToPair[int(trade.PairIndex)] + "*" + " @ " + "*" + strconv.FormatFloat(trade.OpenPrice, 'f', 2, 64) + "*" + "\n\n Kindly wait for confirmation."
	progressMsg := tgbotapi.NewMessage(chatID, inProgress)
	progressMsg.ParseMode = tgbotapi.ModeMarkdown
	_, err := bot.Send(progressMsg)
	if err != nil {
		log.Println("Error sending progress message", err)
	}

	// We need to check here if market or not:
	tradeJson := tradecache.ConvertToOpenTradeJSON(trade)

	// Convert the trade data to JSON
	jsonData, err := json.Marshal(tradeJson)
	if err != nil {
		return "Marshalling Error", tradeJson, fmt.Errorf("failed to marshal trade data: %w", err)
	}

	// Make the POST request
	resp, err := http.Post("http://localhost:3030/opentradegns", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "POST not possible", api.OpenTradeJSON{}, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response if necessary
	// For now, just checking for a successful status code
	if resp.StatusCode != http.StatusOK {
		return "Wrong StatusCode", api.OpenTradeJSON{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "respBody Error", api.OpenTradeJSON{}, fmt.Errorf("failed to read response body: %w", err)
	}
	txReceipt := string(respBodyBytes)

	fmt.Println(txReceipt)

	return txReceipt, tradeJson, nil
}
