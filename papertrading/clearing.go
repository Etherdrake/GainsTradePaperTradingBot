package papertrading

import (
	"HootTelegram/api"
	"HootTelegram/priceserver"
	"HootTelegram/redismanagers/ordercache"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// CentralClearingPaper takes non-market orders and then:
// 1. Puts the order in an array:
// 2. Checks the RedisPrice every second and compares it with the prices in the array:
// 3.
func CentralClearingPaper(rdbPrice, rdbPositionsPaper *redis.Client, pendingPaperChan chan api.OpenTradeJSON) {
	var (
		mu             sync.Mutex
		pendingOrders  []api.OpenTradeJSON
		pendingOrdersC = make(chan []api.OpenTradeJSON)
	)

	// Create a Goroutine for continuous reading from pendingPaperChan
	go func() {
		for trade := range pendingPaperChan {
			mu.Lock()
			pendingOrders = append(pendingOrders, trade)
			mu.Unlock()
		}
	}()

	// Create a Goroutine to continuously update the copy of pendingOrders
	go func() {
		var lastIndexProcessed int

		for {
			mu.Lock()
			// Only copy the trades that haven't been processed yet
			pendingOrdersCopy := pendingOrders[lastIndexProcessed:]
			lastIndexProcessed = len(pendingOrders) // Update the index to the last trade
			mu.Unlock()

			pendingOrdersC <- pendingOrdersCopy
			time.Sleep(time.Second)
		}
	}()

	// Create a ticker to trigger the iteration every t second
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentPendingOrders := <-pendingOrdersC
		// Iteration function
		for i := 0; i < len(currentPendingOrders); i++ {
			fmt.Println("PENDING RECEIVED: ", currentPendingOrders[i])
			trade := currentPendingOrders[i]
			fmt.Println("TRADE: ", trade)
			tradeIdxStr, err := strconv.Atoi(trade.PairIndex)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
			}
			openPrice, err := strconv.ParseFloat(trade.OpenPrice, 64)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
			}
			currentPrice, err := priceserver.GetPrice(rdbPrice, tradeIdxStr)
			if err != nil {
				log.Println("Error in CentralClearingPaper", err)
			}

			// LONG
			if trade.Buy {
				// LIMIT
				if trade.OrderType == 1 {
					if openPrice > currentPrice {
						mu.Lock()
						err := ordercache.ChangeOrderStatus(rdbPositionsPaper, strconv.FormatInt(trade.ID, 10), trade.TradeID, "3")
						mu.Unlock()
						if err != nil {
							log.Println("Error in clearing.go changing order status", err)
						}
					}
				}
				// STOP
				if trade.OrderType == 2 {
					if openPrice < currentPrice {
						mu.Lock()
						err := ordercache.ChangeOrderStatus(rdbPositionsPaper, strconv.FormatInt(trade.ID, 10), trade.TradeID, "3")
						mu.Unlock()
						if err != nil {
							log.Println("Error in clearing.go changing order status", err)
						}
					}
				}
			} else { // SHORT
				// LIMIT
				if trade.OrderType == 1 {
					if openPrice < currentPrice {
						mu.Lock()
						err := ordercache.ChangeOrderStatus(rdbPositionsPaper, strconv.FormatInt(trade.ID, 10), trade.TradeID, "3")
						mu.Unlock()
						if err != nil {
							log.Println("Error in clearing.go changing order status", err)
						}
					}
				}
				// STOP
				if trade.OrderType == 2 {
					if openPrice > currentPrice {
						mu.Lock()
						err := ordercache.ChangeOrderStatus(rdbPositionsPaper, strconv.FormatInt(trade.ID, 10), trade.TradeID, "3")
						mu.Unlock()
						if err != nil {
							log.Println("Error in clearing.go changing order status", err)
						}
					}
				}
			}
		}
	}
}
