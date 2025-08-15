package boardbuilders

import (
	"HootTelegram/api"
	"HootTelegram/pairmaps"
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func BuildTpBoardEdit(rdbPositions, rdbPositionsPaper *redis.Client, tradeCache *tradecache.TradeCache, guid int64, tradeID string, idx string, isPaperStr string) tgbotapi.InlineKeyboardMarkup {
	guidStr := strconv.Itoa(int(guid))

	//var paperStr string
	//if user.Paper {
	//	paperStr = "papertrade"
	//} else {
	//	paperStr = "realtrade"
	//}
	idxIncr, _ := strconv.Atoi(idx)
	idxIncr += 1

	var trade api.OpenTradeJSON
	if isPaperStr == "papertrade" {
		trade, _ = ordercache.GetOpenTradeFromRedis(rdbPositionsPaper, guidStr, tradeID)
		//if err != nil {
		//	log.Println("Error in GetOpenTradeFromRedis during HandleStopLossEdit", err)
		//}
		//isPaper = true
	} else {
		trade, _ = ordercache.GetOpenTradeFromRedis(rdbPositions, guidStr, tradeID)
		//if err != nil {
		//	log.Println("Error in GetOpenTradeFromRedis during HandleStopLossEdit", err)
		//}
		//isPaper = false
	}

	pairIdxStr, _ := strconv.Atoi(trade.PairIndex)
	pairStr := pairmaps.IndexToPair[pairIdxStr]

	redisTpStr, err := strconv.ParseFloat(trade.TP, 64)
	if err != nil {
		log.Println("Error parsing redisTpStr in editpositionboards.go")
	}
	tpRdbPos := utils.FormatPrice(redisTpStr)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 "+pairStr, "/refreshtpeditboard+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Take Profit", "/refreshtpeditboard+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/refreshsleditboard+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/editdecrtp+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData(tpRdbPos, "/editcustomtp+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/editincrtp+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("None", "/editzerotp+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("25%", "/editplus25+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("50%", "/editplus50+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("100%", "/editplus100+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("200%", "/editplus150+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("900%", "/editplus900+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Save", "/backpos+"+idx),
		),
	)
}

func BuildSlBoardEdit(rdbPositions, rdbPositionsPaper *redis.Client, tradeCache *tradecache.TradeCache, guid int64, tradeID string, idx string, isPaperStr string) tgbotapi.InlineKeyboardMarkup {
	guidStr := strconv.Itoa(int(guid))

	user, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

	//var paperStr string
	//if user.Paper {
	//	paperStr = "papertrade"
	//} else {
	//	paperStr = "realtrade"
	//}

	pairStr := pairmaps.IndexToPair[int(user.PairIndex)]

	idxIncr, _ := strconv.Atoi(idx)
	idxIncr += 1
	//idxInrStr := strconv.Itoa(idxIncr)

	//var backSave string
	var trade api.OpenTradeJSON
	if isPaperStr == "papertrade" {
		trade, _ = ordercache.GetOpenTradeFromRedis(rdbPositionsPaper, guidStr, tradeID)
		//if err != nil {
		//	log.Println("Error in GetOpenTradeFromRedis during HandleStopLossEdit", err)
		//}
		//isPaper = true
		//backSave = "/x" + idxInrStr
	} else {
		trade, _ = ordercache.GetOpenTradeFromRedis(rdbPositions, guidStr, tradeID)
		//if err != nil {
		//	log.Println("Error in GetOpenTradeFromRedis during HandleStopLossEdit", err)
		//}
		//isPaper = false
		//backSave = "/t" + idxInrStr
	}

	redisSlStr, err := strconv.ParseFloat(trade.SL, 64)
	if err != nil {
		log.Println("Error parsing redisTpStr in editpositionboards.go")
	}
	slRdbPos := utils.FormatPrice(redisSlStr)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 "+pairStr, "/refreshtpeditboard+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take Profit", "/refreshtpeditboard+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("🛑 Stop Loss", "/refreshsleditboard+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/editdecrsl+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData(slRdbPos, "/editcustomsl+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/editincrsl+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("None", "/editzerosl+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("-10%", "/editminus10+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("-25%", "/editminus25+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("-33%", "/editminus33+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("-50%", "/editminus50+"+tradeID+"+"+isPaperStr+"index="+idx),
			tgbotapi.NewInlineKeyboardButtonData("-75%", "/editminus75+"+tradeID+"+"+isPaperStr+"index="+idx),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Save", "/backpos+"+idx),
		),
	)
}
