package papertrading

import (
	"HootTelegram/api"
	"HootTelegram/database"
	"HootTelegram/priceserver"
	"HootTelegram/redismanagers/ordercache"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// CentralDispatchingPaper takes non-market orders and then:
// 1. Retrieve ALL trades that are OPEN
// 2. Check for every single trade, whether this trade should be closed every single second
// 3. If the trades should be closed, change the ORDER_STATUS to 4.
func CentralDispatchingPaper(rdbPrice, rdbPositionsPaper *redis.Client, pendingPaperChan chan api.OpenTradeJSON) {
	// Create a ticker to trigger every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Get all trades into an array
		allOrders, err := ordercache.GetAllPositionsRdb(rdbPositionsPaper)
		if err != nil {
			log.Println("Error in Dispatching GetAllPositionsRdb", err)
			continue // Skip this iteration if there's an error
		}

		for idx := 0; idx < len(allOrders); idx++ {
			trade := allOrders[idx]
			if trade.PairIndex == "" {
				continue
			}
			tradeIdxStr, err := strconv.Atoi(trade.PairIndex)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
				continue // Skip this trade if there's an error
			}
			TP, err := strconv.ParseFloat(trade.TP, 64)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
				continue // Skip this trade if there's an error
			}
			SL, err := strconv.ParseFloat(trade.SL, 64)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
				continue // Skip this trade if there's an error
			}

			currentPrice, err := priceserver.GetPrice(rdbPrice, tradeIdxStr)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
				continue // Skip this trade if there's an error
			}

			if trade.Buy {
				// TP is hit, close trade
				if (currentPrice > TP && TP != 0) || (currentPrice < SL && SL != 0) {
					//posSize, _ := strconv.ParseFloat(trade.PositionSizeDai, 64)
					//entryPrice, _ := strconv.ParseFloat(trade.OpenPrice, 64)
					//leverage, _ := strconv.ParseFloat(trade.Leverage, 64)
					//PnL := utils.CalculateDollarProfitOrLossWithLeverage(posSize, entryPrice, currentPrice, leverage, trade.Buy)
					err := ordercache.ChangeOrderStatus(rdbPositionsPaper, strconv.FormatInt(trade.ID, 10), trade.TradeID, "4")
					if err != nil {
						log.Println("Error in dispatching ChangeOrderStatus CentralDispatchingPaper")
					}
					PnL := ConvertOpenTradeToPnL(trade, rdbPrice)
					if PnL > 0 {
						database.AddPaperBalance(trade.ID, trade.PnL)
					} else {
						database.DecrPaperBalance(trade.ID, trade.PnL)
					}
				}
			} else {
				if (currentPrice < TP && TP != 0) || (currentPrice > SL && SL != 0) {
					//posSize, _ := strconv.ParseFloat(trade.PositionSizeDai, 64)
					//entryPrice, _ := strconv.ParseFloat(trade.OpenPrice, 64)
					//leverage, _ := strconv.ParseFloat(trade.Leverage, 64)
					//PnL := utils.CalculateDollarProfitOrLossWithLeverage(posSize, entryPrice, currentPrice, leverage, trade.Buy)
					ordercache.ChangeOrderStatus(rdbPositionsPaper, strconv.FormatInt(trade.ID, 10), trade.TradeID, "4")
					if trade.PnL > 0 {
						database.AddPaperBalance(trade.ID, trade.PnL)
					} else {
						database.DecrPaperBalance(trade.ID, trade.PnL)
					}
				}
			}
		}
	}
}
