package aaaboards

import (
	"HootTelegram/api"
	"HootTelegram/database"
	"HootTelegram/mongolisten"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/tradecache"
	"HootTelegram/transformers"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
)

func BuildActiveTradeBoardGNS(client *mongo.Client, rdbPrice *redis.Client,
	orderID int, openTrade bool, openOrder bool, guid int64, isArb bool) tgbotapi.InlineKeyboardMarkup {
	publicKey, err := database.GetPublicKey(guid)
	if err != nil {
		fmt.Println("PublicKey Error: ", err)
	}

	var newTradeStr string
	var activeTradesStr string

	var activeOrder mongolisten.OpenLimitOrderConverted
	var activeMarketTrade transformers.MarketExecutedTransform
	var activeLimitTrade transformers.LimitExecutedTransform

	var marketActive bool

	// If openTrade, take the first Trade from the array
	if openTrade {
		marketTrades, limitTrades, err := mongolisten.GetAllTrades(client, publicKey, isArb)
		if err != nil {
			log.Println("GetAllTraders error: ", err)
		}

		if orderID == 0 {
			// If there is no active orderID, find whatver is the most recent trade and make it active
			if len(marketTrades) > 0 && len(limitTrades) > 0 {
				for x := 0; x < len(marketTrades); x++ {
					if marketTrades[x].OrderID > activeMarketTrade.OrderID {
						activeMarketTrade = marketTrades[x]
					}
				}
				for i := 0; i < len(limitTrades); i++ {
					if limitTrades[i].OrderID > activeLimitTrade.OrderID {
						activeLimitTrade = limitTrades[i]
					}
				}
			}
			if activeLimitTrade.OrderID > activeMarketTrade.OrderID {
				marketActive = false
			}
		}
	}

	// Build only orders at idx0
	if !openTrade && openOrder {
		activeOrders, err := mongolisten.GetOpenOrders(client, publicKey, isArb)
		if err != nil {
			log.Println("GetOpenOrders error: ", err)
		}
		var newestOrder mongolisten.OpenLimitOrderConverted
		var highestNumber int64
		// If we have more than one activeOrder, make the newest one active
		if len(activeOrders) > 1 {
			for p := 0; p < len(activeOrders); p++ {
				if activeOrders[p].Block > highestNumber {
					highestNumber = activeOrders[p].Block
					newestOrder = activeOrders[p]
				} else {
					continue
				}
			}
			activeOrder = newestOrder
		} else {
			activeOrder = activeOrders[0]
		}
	}

	var positionSize float64
	var openPrice float64
	var pairIdx int64
	var leverage int64

	if openTrade {
		if marketActive {
			positionSize = activeMarketTrade.PositionSizeDai
			openPrice = activeMarketTrade.Trade.OpenPrice
			pairIdx = activeMarketTrade.Trade.PairIndex
			leverage = activeMarketTrade.Trade.Leverage

		} else {
			positionSize = activeLimitTrade.PositionSizeDai
			openPrice = activeLimitTrade.Trade.OpenPrice
			pairIdx = activeLimitTrade.Trade.PairIndex
			leverage = activeLimitTrade.Trade.Leverage
		}
	} else {
		if openOrder {
			positionSize = activeOrder.PositionSize
			openPrice = activeOrder.MinPrice
			pairIdx = activeOrder.PairIndex
			leverage = activeOrder.Leverage
		}
	}

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, int(pairIdx)) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	var getPnlUSD float64
	if !openTrade && openOrder {
		getPnlUSD = 0
	} else {
		if marketActive {
			getPnlUSD = utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, float64(leverage), activeMarketTrade.Trade.Buy)
		} else {
			getPnlUSD = utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, float64(leverage), activeLimitTrade.Trade.Buy)
		}

	}

	var pnlStr string
	if getPnlUSD > 0 {
		plusSymbol := "+"
		pnlStr = strconv.FormatFloat(getPnlUSD, 'f', 2, 64)
		pnlStr = plusSymbol + pnlStr
	} else {
		pnlStr = strconv.FormatFloat(getPnlUSD, 'f', 2, 64)
	}

	newTradeStr = "NEW TRADE"
	activeTradesStr = "⭐️ ACTIVE TRADES"
	isPaperStr := "realtrade"

	var closeTradeOrCancelOrder string
	if !openTrade && openOrder {
		closeTradeOrCancelOrder = "Cancel Order"
	} else {
		closeTradeOrCancelOrder = "Close Position: " + "💰" + pnlStr + ""
	}

	// If trade or order /t or /o + ":" + "SYMBOL"
	symbol := pairmaps.IndexToCrypto[int(pairIdx)]

	closeButton := tgbotapi.NewInlineKeyboardButtonData(closeTradeOrCancelOrder, "/placeholder")

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newTradeStr, "/newtrade"),
			tgbotapi.NewInlineKeyboardButtonData(activeTradesStr, "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Refresh", "/activetrades"),
			tgbotapi.NewInlineKeyboardButtonData(": "+symbol, "/nextrade"),
			tgbotapi.NewInlineKeyboardButtonData("Chart", "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take Profit", "/edittakeprofit+"+"69"+"+"+isPaperStr+"index="),
			tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/editstoploss+"+"69"+"+"+isPaperStr+"index="),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Share on Telegram", "/sharetelegram+"+"X69"+"+"),
			tgbotapi.NewInlineKeyboardButtonData("Share on Twitter", "/sharetwitter+"+"X69"+"+"),
		),
		tgbotapi.NewInlineKeyboardRow(
			// The position in question
			closeButton),
	)
}

func BuildActiveTradeBoard(rdbPaper *redis.Client, rdbPrice *redis.Client,
	tradeCache *tradecache.TradeCache, idx int, guid int64, isPaper bool) tgbotapi.InlineKeyboardMarkup {
	var trade api.OpenTradeJSON

	var newTradeStr string
	var activeTradesStr string

	// Check if user has open positions in Redis
	if ordercache.HasTrades(rdbPaper, strconv.FormatInt(guid, 10)) {
		allOrders, err := ordercache.GetAllOrdersForUserArrayFromRedis(rdbPaper, strconv.FormatInt(guid, 10))
		if err != nil {
			log.Println("Error encountered get AllOrdersArray in BuildPositionsAndOrderString", err)
		}

		// Seperate open_trades and open_orders
		openTradesArray := ordercache.FilterOpenTrades(allOrders)
		openOrdersArray := ordercache.FilterOpenOrders(allOrders)

		// Sort both arrays
		openTradesArraySorted := ordercache.SortArrayByTimestamp(openTradesArray)
		openOrdersArraySorted := ordercache.SortArrayByTimestamp(openOrdersArray)

		if len(openTradesArraySorted) > 0 {
			// Get the position using the idx provided
			trade = openTradesArraySorted[idx]
		} else {
			if len(openOrdersArraySorted) > 0 {
				trade = openOrdersArray[idx]
			} else {
				fmt.Println("No OpenTrades or OpenOrders have been found.")
			}
		}
	}

	positionSize, err := strconv.ParseFloat(trade.PositionSizeDai, 64)
	openPrice, err := strconv.ParseFloat(trade.OpenPrice, 64)
	pairIdx, _ := strconv.Atoi(trade.PairIndex)
	leverage, err := strconv.ParseFloat(trade.Leverage, 64)

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, pairIdx) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	getPnlUSD := utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, leverage, trade.Buy)

	var pnlStr string
	if getPnlUSD > 0 {
		plusSymbol := "+"
		pnlStr = strconv.FormatFloat(getPnlUSD, 'f', 2, 64)
		pnlStr = plusSymbol + pnlStr
	} else {
		pnlStr = strconv.FormatFloat(getPnlUSD, 'f', 2, 64)
	}

	cacheTrade, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

	// Check on which Window we are:
	if cacheTrade.ActiveWindow == 0 { // We are on new trade
		newTradeStr = "⭐️ NEW TRADE"
		activeTradesStr = "ACTIVE TRADES"
	} else {
		if cacheTrade.ActiveWindow == 1 {
			newTradeStr = "NEW TRADE"
			activeTradesStr = "⭐️ ACTIVE TRADES"
		} else {
			log.Println("Error during BuildPerpMainBoard: ActiveWindow not equal to 0 or 1")
		}
	}

	var isPaperStr string
	if isPaper {
		isPaperStr = "papertrade"
	} else {
		isPaperStr = "realtrade"
	}

	idxStr := strconv.Itoa(idx)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newTradeStr, "/newtrade"),
			tgbotapi.NewInlineKeyboardButtonData(activeTradesStr, "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Take Profit", "/edittakeprofit+"+trade.TradeID+"+"+isPaperStr+"index="+idxStr),
			tgbotapi.NewInlineKeyboardButtonData("🛑 Stop Loss", "/editstoploss+"+trade.TradeID+"+"+isPaperStr+"index="+idxStr),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Share on Telegram", "/sharetelegram+"+trade.TradeID+"+"),
			tgbotapi.NewInlineKeyboardButtonData("Share on Twitter", "/sharetwitter+"+trade.TradeID+"+"),
		),
		tgbotapi.NewInlineKeyboardRow(
			// The position in question
			tgbotapi.NewInlineKeyboardButtonData("Close Position: "+"💰"+pnlStr+"", "/closeposition+"+trade.TradeID+"+"+isPaperStr),
		),
	)
}
