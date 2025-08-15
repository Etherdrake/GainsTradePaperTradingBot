package priceserver

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

var ctx = context.Background()

type PriceData struct {
	PairIndex int
	Price     float64
}

func InitPriceStream(pricedata chan PriceData) error {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("wss://backend-pricing.eu.gains.trade", nil)
	if err != nil {
		log.Println("Error connecting to WebSocket:", err)
		return err // Return the error instead of terminating the program
	}
	defer conn.Close()

	go SendPingMessages(conn)

	for {
		var rawData []float64
		err := conn.ReadJSON(&rawData)
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			return err // Return the error instead of breaking out of the loop
		}

		// Process rawData two items at a time (index, price)
		for i := 0; i < len(rawData)-1; i += 2 {
			index := int(rawData[i])
			price := rawData[i+1]
			pricedata <- PriceData{PairIndex: index, Price: price}
		}
	}
}

type TradeData struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func InitTradeStream() {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("wss://backend-polygon.gains.trade/", nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	for {
		var tradeData TradeData
		err := conn.ReadJSON(&tradeData)
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}
		fmt.Printf("Received Trade Data: %+v\n", tradeData)
	}
}

func SendPingMessages(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second) // Send a ping every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))
			if err != nil {
				log.Println("Error sending ping:", err)
				return
			}
		}
	}
}

func ReconnectLoop(pricedata chan PriceData) {
	for {
		err := InitPriceStream(pricedata)
		if err != nil {
			log.Println("Error encountered in reconnect loop.")
			return
		}
		log.Println("Connection lost. Attempting to reconnect...")
		time.Sleep(5 * time.Second) // Wait before attempting to reconnect
	}
}

func GetPrice(rdb *redis.Client, pairIndex int) (float64, error) {
	// Generate the key used in Redis based on pairIndex
	key := fmt.Sprintf("price:%d", pairIndex)

	// Fetch the value from Redis
	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("REDIS ERR", err)
		if err == redis.Nil {
			log.Printf("No price found for pairIndex %d", pairIndex)
			return 0, fmt.Errorf("no price found for pairIndex %d", pairIndex)
		}
		return 0, err
	}

	// Convert the string value to float64
	price, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing price value for pairIndex %d: %s", pairIndex, value)
		return 0, err
	}

	return price, nil
}

