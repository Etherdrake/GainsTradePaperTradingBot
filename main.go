package main

import (
	"HootTelegram/api"
	"HootTelegram/concurrentmaps"
	"HootTelegram/papertrading"
	"HootTelegram/redislocker"
	"HootTelegram/tradecache"
	"HootTelegram/types"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	ownerID        = 1099978784 // Replace with the ID of the bot owner
	botToken       = "6559595319:AAHy3F8_2mfwa9FwGeQM0Fmgd_S6VIkuhlY"
	databaseName   = "hooterdb"
	collectionName = "users"
	mongoURI       = "mongodb://localhost:27017"
	apiURI         = "localhost:3030"
)

var GlobalPriceCache types.PriceCache

var ctx = context.Background()

func InitRedisPriceStream() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // change this to your Redis server address
		DB:   1,                // using logical database 1 for pricedata
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil
	}

	fmt.Println("Connected to Redis!")
	return rdb
}

func main() {
	// Initialize the standard Redis client
	rdbPaper := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdbPaper.Close()

	// Initialize the Price Stream Redis client
	rdbPrice := InitRedisPriceStream()
	defer rdbPrice.Close()

	rdbMinMaxPos := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       4,
	})
	defer rdbMinMaxPos.Close()

	// Initialize the Redis Tracker
	rdbLocker := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   5,
	})

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}
	ctxTimeSeries, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctxTimeSeries)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(ctxTimeSeries)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Existing logic for bot updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updatesChan := bot.GetUpdatesChan(u)

	// Create a channel to distribute updates among workers
	workerUpdateChan := make(chan tgbotapi.Update)

	tradeCache := tradecache.New()
	boardCache := concurrentmaps.New()

	// Papertrade order channel
	pendingPaperChan := make(chan api.OpenTradeJSON)

	errWipe := redislocker.WipeCache(ctx, rdbLocker)
	if errWipe != nil {
		log.Println("Fatal encountered during RedisLocker Flush: ", errWipe)
		return
	}

	go func() {
		for update := range updatesChan {
			// Check if we are sending this update ourselves
			if !update.SentFrom().IsBot {
				go func() {
					userID := update.SentFrom().ID
					exists, errCheck := redislocker.CheckUserID(ctx, rdbLocker, userID)
					if errCheck != nil {
						log.Println("Error checking UserID: ", userID)
						return
					}
					if exists {
						log.Println("User ")
						return
					}
					// Add userID to Redis
					errAdd := redislocker.AddUserID(ctx, rdbLocker, userID)
					if errAdd != nil {
						log.Printf("Error add UserID %d to Redis: %v", userID, errAdd)
					}
					// Distribute the update to worker goroutines
					go func() {
						handleUpdates(bot, client, update, tradeCache, boardCache, rdbPaper, rdbMinMaxPos, rdbPrice, workerUpdateChan, pendingPaperChan)
						// Remove UserID from Redis after handleUpdates is done
						if errRemove := redislocker.RemoveUserID(ctx, rdbLocker, userID); errRemove != nil {
							log.Printf("Error removing UserID %d from Redis: %v", userID, errRemove)
						}
					}()
				}()
			}
		}
	}()

	go papertrading.CentralClearingPaper(rdbPrice, rdbPaper, pendingPaperChan)

	go papertrading.CentralDispatchingPaper(rdbPrice, rdbPaper, pendingPaperChan)

	// Block the main function from exiting
	select {}
}
