package main

import (
	"HootTelegram/api"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/tradecache"
	"errors"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func handleSubmitPaper(bot *tgbotapi.BotAPI, rdbPrice *redis.Client, tradeCache *tradecache.TradeCache, chatID int64, guid int64, pendingPaper chan api.OpenTradeJSON) (string, api.OpenTradeJSON, tgbotapi.Message, error) {
	// Fetch the trade from the global tradeCache
	trade, exists := tradeCache.Get(guid)
	if !exists {
		return "Cache Error", api.OpenTradeJSON{}, tgbotapi.Message{}, errors.New("trade not found in cache")
	}

	priceReal, err := priceserver.GetPrice(rdbPrice, int(trade.Index))
	if err != nil {
		log.Println("Error retrieving price in handle submit paper")
	}

	if trade.Buy {
		if trade.SL > priceReal {
			return "Stop Loss", api.OpenTradeJSON{}, tgbotapi.Message{}, errors.New("stop loss can not be higher than the current price for a long")
		}
	} else {
		if trade.SL < priceReal {
			return "Stop Loss", api.OpenTradeJSON{}, tgbotapi.Message{}, errors.New("stop loss can not be lower than the current price for a short")
		}
	}
	buyOrSell := ""
	if trade.Buy {
		buyOrSell = "BUY"
	} else {
		buyOrSell = "SELL"
	}

	// We have a market order
	if trade.OrderType == 0 {
		trade.OpenPrice = priceReal
	}

	inProgress := "Executing " + "*" + buyOrSell + "*" + " for " + "*" + pairmaps.IndexToPair[int(trade.PairIndex)] + "*" + " @ " + "*" + strconv.FormatFloat(trade.OpenPrice, 'f', 2, 64) + "*" + "\n\n Kindly wait for confirmation."
	progressMsg := tgbotapi.NewMessage(chatID, inProgress)
	progressMsg.ParseMode = tgbotapi.ModeMarkdown
	inProgressDelMsg, err := bot.Send(progressMsg)
	if err != nil {
		log.Printf("Error sending progress message", err)
	}
	// if MARKET set to OPEN_TRADE
	if trade.OrderType == 0 {
		trade.OrderStatus = 3
	}
	// If LIMIT set to OPEN_ORDER
	if trade.OrderType == 1 {
		trade.OrderStatus = 2
	}
	// if STOP set to OPEN_ORDER
	if trade.OrderType == 2 {
		trade.OrderStatus = 2
	}

	// Sent trade to clearing if it's not a market order
	if trade.OrderType != 0 {
		tradeJSON := tradecache.ConvertToOpenTradeJSON(trade)
		pendingPaper <- tradeJSON
	}

	tradeJson := tradecache.ConvertToOpenTradeJSON(trade)

	txReceipt := "Your paper trade has been placed!"
	return txReceipt, tradeJson, inProgressDelMsg, nil
}
