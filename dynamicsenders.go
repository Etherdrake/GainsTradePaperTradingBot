package main

import (
	"HootTelegram/boardbuilders"
	"HootTelegram/boardbuilders/stringbuilders"
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func UpdateNewTradeOrActiveTrades(
	bot *tgbotapi.BotAPI, client *mongo.Client, trade tradecache.OpenTradeCache,
	rdbPrice *redis.Client, rdbPositionsPaper *redis.Client, tradeCache *tradecache.TradeCache,
	boardMsgChatId int64, boardMsgMsgId int,
	idx int, guid int64) error {
	switch trade.ActiveWindow {
	case 0:
		board := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)
		updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
		finalMsg, err := utils.CreateEditMessage(board, updatedMessage, boardMsgChatId, boardMsgMsgId, "Markdown")
		if err != nil {
			log.Println("Error CreatingEditMessage message: ", err)
			return err
		}
		_, err = bot.Send(finalMsg)
		if err != nil {
			log.Println("Error sending message: ", err)
			return err
		}

	case 1:
		board := boardbuilders.BuildActiveTradeBoardGNSV2(rdbPrice, rdbPositionsPaper, guid, trade.ActiveTradeID)
		updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
		finalMsg, err := utils.CreateEditMessage(board, updatedMessage, boardMsgChatId, boardMsgMsgId, "Markdown")
		if err != nil {
			log.Println("Error in CreateEditMessage in case polygon")
		}
		_, err = bot.Send(finalMsg)
		if err != nil {
			log.Println("Error sending message")
			return err
		}
		return nil
	}
	return nil
}
