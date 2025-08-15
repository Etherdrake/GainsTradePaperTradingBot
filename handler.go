package main

import (
	"HootTelegram/alertatron"
	"HootTelegram/api"
	"HootTelegram/boardbuilders"
	"HootTelegram/boardbuilders/stringbuilders"
	"HootTelegram/boards/keyboards"
	"HootTelegram/boards/walletboard"
	"HootTelegram/chartservice"
	"HootTelegram/concurrentmaps"
	"HootTelegram/database"
	"HootTelegram/errorhandling"
	"HootTelegram/minmaxpos"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/redismanagers/openpositionedit"
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/search"
	"HootTelegram/tradecache"
	"HootTelegram/tradesettings"
	"HootTelegram/utils"
	"HootTelegram/wallet"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUpdates(bot *tgbotapi.BotAPI, client *mongo.Client, update tgbotapi.Update, tradeCache *tradecache.TradeCache, boardCache *concurrentmaps.BoardMessagesCache, rdbPositionsPaper *redis.Client, rdbMinMaxPos *redis.Client, rdbPrice *redis.Client, updates <-chan tgbotapi.Update, pendingPaperChan chan api.OpenTradeJSON) {
	if update.CallbackQuery != nil {
		// Handle callback queries here
		callbackData := update.CallbackQuery.Data
		fmt.Println("CALLBACKDATA: ", callbackData)
		//boardMsg, ok := boardMessages[update.CallbackQuery.From.ID]
		boardCache.Set(update.CallbackQuery.From.ID, update.CallbackQuery.Message)

		var msg tgbotapi.Chattable
		guid := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.Chat.ID
		//guidStr := strconv.FormatInt(guid, 10)

		boardMsg, ok := boardCache.Get(update.CallbackQuery.From.ID)
		if !ok {
			log.Println("Error retrieving BoardMessageNew")
		}
		switch callbackData {
		case "/search":
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please reply to this message with your query: ")
			msg.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply:            true,
				InputFieldPlaceholder: "Type the symbol or part of the symbol for which you want to search.",
				Selective:             false,
			}
			_, _ = bot.Send(msg)
			return
		case "/leveragefromchart":
			utils.DeleteMessage(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

			board := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)
			// Update the message content with the desired message
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			perpMenuString := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			// Create a new editable message with the updated text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, perpMenuString)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msg.ReplyMarkup = board
			msg.ParseMode = tgbotapi.ModeMarkdown

			spoofMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Loading...")
			spoofMsg.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: false,
				Selective:  true,
			}
			spoofDel, _ := bot.Send(spoofMsg)
			utils.DeleteMessage(bot, spoofDel.Chat.ID, spoofDel.MessageID)
			return
		case "/leverage":
			board := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)
			// Update the message content with the desired message
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// Update the message content with the desired message
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			// Create a new editable message with the updated text
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/leverageback":
			utils.DeleteMessage(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			board := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)
			// Update the message content with the desired message
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			perpMenuString := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			// Create a new editable message with the updated text
			msgNew := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, perpMenuString)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgNew.ReplyMarkup = board
			msgNew.ParseMode = tgbotapi.ModeMarkdown
			_, err := bot.Send(msgNew)
			if err != nil {
				log.Println(err)
			}
			return
		case "/sendchart":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			b64, err := chartservice.GetChart1M(pairmaps.IndexToCrypto[int(trade.PairIndex)])
			if err != nil {
				log.Println(err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, boardMsg.MessageID)

			err = chartservice.SendChartWithKeyboard1M(bot, update.CallbackQuery.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
			return
		case "/sendchart5M":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			b64, err := chartservice.GetChart5M(pairmaps.IndexToCrypto[int(trade.PairIndex)])
			if err != nil {
				log.Println(err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, boardMsg.MessageID)

			err = chartservice.SendChartWithKeyboard5M(bot, update.CallbackQuery.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
			return
		case "/sendchart15M":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			b64, err := chartservice.GetChart15M(pairmaps.IndexToCrypto[int(trade.PairIndex)])
			if err != nil {
				log.Println(err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, boardMsg.MessageID)

			err = chartservice.SendChartWithKeyboard15M(bot, update.CallbackQuery.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
			return
		case "/sendchart1H":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			b64, err := chartservice.GetChart1H(pairmaps.IndexToCrypto[int(trade.PairIndex)])
			if err != nil {
				log.Println(err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, boardMsg.MessageID)

			err = chartservice.SendChartWithKeyboard1H(bot, update.CallbackQuery.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
			return
		case "/sendchart4H":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			b64, err := chartservice.GetChart4H(pairmaps.IndexToCrypto[int(trade.PairIndex)])
			if err != nil {
				log.Println(err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, boardMsg.MessageID)

			err = chartservice.SendChartWithKeyboard4H(bot, update.CallbackQuery.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
			return
		case "/sendchart1D":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			b64, err := chartservice.GetChart1D(pairmaps.IndexToCrypto[int(trade.PairIndex)])
			if err != nil {
				log.Println(err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, boardMsg.MessageID)

			err = chartservice.SendChartWithKeyboard1D(bot, update.CallbackQuery.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
			return
		case "/newtrade":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			err := tradeCache.SetActiveWindow(guid, 0)
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/activetrades":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			hasTrades := ordercache.HasTrades(rdbPositionsPaper, strconv.FormatInt(guid, 10))
			if hasTrades {
				err := tradeCache.SetActiveWindow(guid, 1)
				if err != nil {
					log.Println("Error setting active window in active trades.")
				}
				trade, exists = tradeCache.Get(guid)
				if !exists {
					log.Println("User does not exist.")
					return
				}
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			} else {
				alertatron.SendAlert(bot, update.CallbackQuery.ID, "You do not have any open trades or orders.")
				return
			}
		case "/submitconfirm":
			board := boardbuilders.BuildPerpMainBoardSubmitCancel(tradeCache, rdbPrice, guid)
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// LONG
			if trade.Buy {
				if trade.SL > trade.OpenPrice {
					alertatron.SendAlert(bot, update.CallbackQuery.ID, "Your stoploss is higher than your open price, this is not possible on a long as it would cause an instant close of your trade.")
					return
				}
				// SHORT
			} else {
				if trade.SL < trade.OpenPrice {
					alertatron.SendAlert(bot, update.CallbackQuery.ID, "Your stoploss is lower than your open price, this is not possible on a short as it would cause an instant close of your trade.")
					return
				}
			}
			// Update the message content with the desired message
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			// Create a new editable message with the updated text
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/submit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			err := tradeCache.GenerateOrderID(guid)
			if err != nil {
				log.Println("Generation of new OrderID failed: ", err)
			}
			txReceipt, tradeData, inProgressMsgDel, err := handleSubmitPaper(bot, rdbPrice, tradeCache, chatID, guid, pendingPaperChan)
			if err != nil {
				if txReceipt == "Stop Loss" {
					errTxt := fmt.Sprint(err)
					errMsg := tgbotapi.NewMessage(chatID, errTxt)
					errDelMsg, err := bot.Send(errMsg)
					if err != nil {
						log.Println("Error sending errDelMsg")
					}
					time.Sleep(2000 * time.Millisecond)
					utils.DeleteMessage(bot, errDelMsg.Chat.ID, errDelMsg.MessageID)
					return
				}
				err = ordercache.StoreOpenTradeInRedis(rdbPositionsPaper, strconv.FormatInt(guid, 10), tradeData)
				if err != nil {
					log.Printf("Error with Redis: %v", err)
				}
				errMsg := tgbotapi.NewMessage(chatID, err.Error())
				_, err := bot.Send(errMsg)
				if err != nil {
					log.Println("Error during storage of new trade to Redis: ", err)
					return
				}
				// NO ERROR BECAUSE ERROR == NIL
			} else {
				var msgTx tgbotapi.MessageConfig
				msgTx = tgbotapi.NewMessage(chatID, "Trade execution successful: [TxReceipt](https://hootscan.com/tx/"+strings.Trim(txReceipt, "\"")+")")
				msgTx.ParseMode = tgbotapi.ModeMarkdown
				err = ordercache.StoreOpenTradeInRedis(rdbPositionsPaper, strconv.FormatInt(guid, 10), tradeData)
				if err != nil {
					log.Println("Error during storage of new trade to Redis: ", err)
				}
				msgTxReceiptDel, err := bot.Send(msgTx)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				trade, exists = tradeCache.Get(guid)
				if !exists {
					log.Println("User does not exist.")
				}
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				// Now delete the old message
				utils.DeleteMessage(bot, inProgressMsgDel.Chat.ID, inProgressMsgDel.MessageID)
				utils.DeleteMessage(bot, msgTxReceiptDel.Chat.ID, msgTxReceiptDel.MessageID)
				return
			}
			return
		case "/wallet":
			isSet, err := database.CheckIfSet(guid)
			if err != nil {
				fmt.Println(err)
				return
			}
			if isSet {
				pubkey, err := database.GetPublicKey(guid)
				if err != nil {
					log.Println("Error retrieving public key in wallet", err)
					return
				}
				updatedMessage, err := stringbuilders.BuildWalletMainString(guid)
				if err != nil {
					fmt.Println("Error building string for users: " + strconv.Itoa(int(guid)))
					fmt.Println(err)
				}
				// Create a new editable message with the updated text
				msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
				qrFileName, err := wallet.CreateQR(pubkey)
				if err != nil {
					fmt.Println(err)
				}
				// Send the QR code as a photo
				photoMedia := wallet.NewInputMediaPhotoFromFile(qrFileName)

				photoMsg := tgbotapi.NewPhoto(guid, photoMedia.BaseInputMedia.Media)
				_, err = bot.Send(photoMsg)
				if err != nil {
					fmt.Println("Photo error")
				}

				_, err = bot.Send(msg)
				if err != nil {
					log.Println("Failed to send QR code photo:", err)
					return
				}
				// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
				msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
				if !ok {
					log.Println("Cannot convert msg to type EditMessageTextConfig")
					return
				}
				// Assign the new inline keyboard to the message
				msgEdit.ReplyMarkup = &walletboard.WalletMainBoard
				msg = msgEdit

			} else {
				pubkey, pk := wallet.GenerateWallet()
				database.AddPublicKey(guid, pubkey)
				database.AddPrivateKey(guid, pk)

				updatedMessage, err := stringbuilders.BuildWalletMainString(guid)
				if err != nil {
					fmt.Println("Error building string for users: " + strconv.Itoa(int(guid)))
					fmt.Println(err)
				}
				// Create a new editable message with the updated text
				msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
				// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
				msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
				if !ok {
					log.Println("Cannot convert msg to type EditMessageTextConfig")
					return
				}
				// Assign the new inline keyboard to the message
				msgEdit.ReplyMarkup = &walletboard.WalletMainBoard
				msg = msgEdit

			}
		// Perpmainboard Interactions
		case "/pairs":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildPairsBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildCryptoPairString(int(trade.PairPage))
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/pairkeyboard":
			msgKeyboard := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Price Refreshed!")
			msgKeyboard.ReplyMarkup = keyboards.GeneratePairKeyboardCryptoExUSDThreeRowFiveCol(0)
			_, err := bot.Send(msgKeyboard)
			if err != nil {
				log.Println("Error in /openpositionselectkeyboard sending message: ", err)
			}
		case "/refresh":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			price, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
			if err != nil {
				log.Printf("Error fetching price for pairIndex %d: %v", trade.PairIndex, err)
				// Handle this error: perhaps return or set a default price value
			}
			err = tradeCache.SetEntryPrice(guid, price)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/stoploss":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/takeprofit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/toggletolong":
			err := tradeCache.ToggleLongShort(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			fmt.Println("This is the trade in the cache: ", trade)
			price, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
			if err != nil {
				log.Println("Error getting price")
			}
			err = tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 75))
			if err != nil {
				fmt.Println(err)
			}
			err = tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 900))
			if err != nil {
				fmt.Println(err)
			}

			if trade.OpenPrice < price && trade.OrderType != 1 {
				err := tradeCache.SetOrderTypeToLimit(guid)
				if err != nil {
					log.Println("User does not exist.")
				}
			}
			if trade.OpenPrice > price && trade.OrderType != 2 {
				err := tradeCache.SetOrderTypeToStop(guid)
				if err != nil {
					log.Println("User does not exist.")
				}
			}

			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/toggletoshort":
			fmt.Println("TOGGLE", tradeCache)
			err := tradeCache.ToggleLongShort(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			price, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
			if err != nil {
				log.Println("Error getting price")
			}
			err = tradeCache.SetTakeProfit(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 75))
			if err != nil {
				fmt.Println(err)
			}
			err = tradeCache.SetStopLoss(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 900))
			if err != nil {
				fmt.Println(err)
			}
			if trade.OrderType != 0 {
				if !trade.Buy {
					if trade.OpenPrice > price && trade.OrderType != 1 {
						err := tradeCache.SetOrderTypeToLimit(guid)
						if err != nil {
							log.Println("User does not exist.")
						}
					}
					if trade.OpenPrice < price && trade.OrderType != 2 {
						err := tradeCache.SetOrderTypeToStop(guid)
						if err != nil {
							log.Println("User does not exist.")
						}
					}
				}
			}

			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/decrprice":
			err := tradeCache.DecrementEntryPrice(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			if trade.OpenPrice <= 0 {
				alertatron.SendAlert(bot, update.CallbackQuery.ID, "Entryprice is lower than or equal to zero! Open price has been reset to current live-price. Happy trading sir :-)")
				price, _ := priceserver.GetPrice(rdbPrice, int(trade.Index))
				err := tradeCache.SetEntryPrice(guid, price)
				if err != nil {
					log.Println("Error encountered during SetEntryPrice")
				}
			}
			price, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
			if err != nil {
				log.Printf("Error fetching price for pairIndex %d: %v", trade.PairIndex, err)
				// Handle this error: perhaps return or set a default price value
			}

			if trade.Buy {
				if trade.OpenPrice < price && trade.OrderType != 1 {
					err := tradeCache.SetOrderTypeToLimit(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
				if trade.OpenPrice > price && trade.OrderType != 2 {
					err := tradeCache.SetOrderTypeToStop(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
			} else {
				if trade.OpenPrice > price && trade.OrderType != 1 {
					err := tradeCache.SetOrderTypeToLimit(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
				if trade.OpenPrice < price && trade.OrderType != 2 {
					err := tradeCache.SetOrderTypeToStop(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/incrprice":
			err := tradeCache.IncrementEntryPrice(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			price, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
			if err != nil {
				log.Printf("Error fetching price for pairIndex %d: %v", trade.PairIndex, err)
				// Handle this error: perhaps return or set a default price value
			}
			if trade.Buy {
				if trade.OpenPrice < price && trade.OrderType != 1 {
					err := tradeCache.SetOrderTypeToLimit(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
				if trade.OpenPrice > price && trade.OrderType != 2 {
					err := tradeCache.SetOrderTypeToStop(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
			} else {
				if trade.OpenPrice > price && trade.OrderType != 1 {
					err := tradeCache.SetOrderTypeToLimit(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
				if trade.OpenPrice < price && trade.OrderType != 2 {
					err := tradeCache.SetOrderTypeToStop(guid)
					if err != nil {
						log.Println("User does not exist.")
					}
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/customprice":
			msgConfig := tgbotapi.NewMessage(boardMsg.Chat.ID, "Please provide your order entry price, type \"market\" for market orders.")
			msgConfig.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			msg = msgConfig
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("This is the customprice: ", err)
			}
			return
		case "/decrpossize":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if !exists {
				log.Println("Does not exist decrpossize")
			}

			err := tradeCache.DecrementPositionSize(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/custompossize":
			msgConfig := tgbotapi.NewMessage(boardMsg.Chat.ID, "Please provide position size")
			msgConfig.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			msg = msgConfig
			// Optional: use ForceReply to show reply interface to the user
			//msg = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		case "/incrpossize":
			err := tradeCache.IncrementPositionSize(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/decrleverage":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			if 2 == trade.Leverage-1 {
				answerTxt := fmt.Sprintf("Leverage for this pair can not be set lower than: %d", 2)
				alertatron.SendAlert(bot, update.CallbackQuery.ID, answerTxt)
				return
			}

			err := tradeCache.DecrementLeverage(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/customleverage":
			msgConfig := tgbotapi.NewMessage(boardMsg.Chat.ID, "Please provide leverage")
			msgConfig.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			msg = msgConfig
		case "/incrleverage":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			maxLev, err := minmaxpos.GetMaxLeverage(int(trade.Index), rdbMinMaxPos)

			if maxLev == trade.Leverage+1 {
				answerTxt := fmt.Sprintf("Leverage for this pair can not be set higher than: %d", maxLev)
				alertatron.SendAlert(bot, update.CallbackQuery.ID, answerTxt)
				return
			}
			err = tradeCache.IncrementLeverage(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		case "/openchart":
			return
			// Take Profit Queries
		case "/zerotp":
			err := tradeCache.SetTakeProfit(guid, 0)
			if err != nil {
				log.Println("Take profit failed to set at /plus25", err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/decrtp":
			err := tradeCache.DecrementTakeProfit(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/customtp":
			msgConfig := tgbotapi.NewMessage(boardMsg.Chat.ID, "Please provide take profit")
			msgConfig.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			msg = msgConfig
			// Take Profit Queries
		case "/incrtp":
			err := tradeCache.IncrementTakeProfit(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/decrpairpage":
			err := tradeCache.DecrementPage(guid)
			if err != nil {
				log.Println("Error decrementing pair page for user: ", guid)
			}
			user, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildPairsBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildCryptoPairString(int(user.PairPage))

			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)

			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/incrpairpage":
			err := tradeCache.IncrementPage(guid, "crypto")
			if err != nil {
				log.Println("Error incrementing pair page for user: ", guid)
			}
			user, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildPairsBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildCryptoPairString(int(user.PairPage))

			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)

			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/decrpairpagefx":
			err := tradeCache.DecrementPage(guid)
			if err != nil {
				log.Println("Error decrementing pair page for user: ", guid)
			}
			user, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildForexBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildFxPairString(int(user.PairPage))

			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)

			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/incrpairpagefx":
			err := tradeCache.IncrementPage(guid, "forex")
			if err != nil {
				log.Println("Error incrementing pair page for user: ", guid)
			}
			user, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildForexBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildFxPairString(int(user.PairPage))

			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)

			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/plus25":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 25))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			} else {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 25))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/plus50":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 50))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			} else {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 50))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/plus100":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 100))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			} else {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 100))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/plus150":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 200))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			} else {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 200))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/plus900":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 900))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			} else {
				err := tradeCache.SetTakeProfit(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 900))
				if err != nil {
					log.Println("Take profit failed to set at /plus25", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/zerotpedit":
			err := tradeCache.SetTakeProfit(guid, 0)
			if err != nil {
				log.Println("Take profit failed to set at /plus25", err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/decrtpedit":
			err := tradeCache.DecrementTakeProfit(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/customtpedit":
			msgConfig := tgbotapi.NewMessage(boardMsg.Chat.ID, "Please provide take profit")
			msgConfig.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			msg = msgConfig
			// Take Profit Queries
		case "/incrtpedit":
			err := tradeCache.IncrementTakeProfit(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildTpBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/plus25edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newTakeProfit := utils.CalculateTakeProfit(openPrice, leverage, 25)
			newTpStr := strconv.FormatFloat(newTakeProfit, 'E', -1, 64)
			err := openpositionedit.UpdateTPInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newTpStr)
			if err != nil {
				log.Println("plus25edit: openpositionedit.UpdateTPInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/plus50edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newTakeProfit := utils.CalculateTakeProfit(openPrice, leverage, 50)
			newTpStr := strconv.FormatFloat(newTakeProfit, 'E', -1, 64)
			err := openpositionedit.UpdateTPInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newTpStr)
			if err != nil {
				log.Println("plus25edit: openpositionedit.UpdateTPInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/plus100edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newTakeProfit := utils.CalculateTakeProfit(openPrice, leverage, 100)
			newTpStr := strconv.FormatFloat(newTakeProfit, 'E', -1, 64)
			err := openpositionedit.UpdateTPInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newTpStr)
			if err != nil {
				log.Println("plus25edit: openpositionedit.UpdateTPInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/plus150edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newTakeProfit := utils.CalculateTakeProfit(openPrice, leverage, 150)
			newTpStr := strconv.FormatFloat(newTakeProfit, 'E', -1, 64)
			err := openpositionedit.UpdateTPInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newTpStr)
			if err != nil {
				log.Println("plus25edit: openpositionedit.UpdateTPInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/plus900edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newTakeProfit := utils.CalculateTakeProfit(openPrice, leverage, 900)
			newTpStr := strconv.FormatFloat(newTakeProfit, 'E', -1, 64)
			err := openpositionedit.UpdateTPInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newTpStr)
			if err != nil {
				log.Println("plus25edit: openpositionedit.UpdateTPInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/decrsl":
			err := tradeCache.DecrementStopLoss(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/customsl":
			msgConfig := tgbotapi.NewMessage(boardMsg.Chat.ID, "Please provide stop loss")
			msgConfig.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: true,
				Selective:  true,
			}
			msg = msgConfig
		case "/incrsl":
			err := tradeCache.IncrementStopLoss(guid)
			if err != nil {
				fmt.Println(err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/zerosl":
			err := tradeCache.SetStopLoss(guid, 0)
			if err != nil {
				log.Println("Stop loss failed to set at /plus300", err)
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/minus10":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 10))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
				fmt.Println("BUY")
			} else {
				err := tradeCache.SetStopLoss(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 10))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
				fmt.Println("SELL")
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/minus25":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 25))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
			} else {
				err := tradeCache.SetStopLoss(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 25))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/minus33":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 33))
				if err != nil {
					log.Println("Stop loss failed to set at /minus33", err)
				}
			} else {
				err := tradeCache.SetStopLoss(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 33))
				if err != nil {
					log.Println("Stop loss failed to set at /minus33", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/minus50":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 50))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
			} else {
				err := tradeCache.SetStopLoss(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 50))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/minus75":
			trade, exists := tradeCache.Get(guid)
			if exists {
				log.Println("User does not exist.")
			}
			if trade.Buy {
				err := tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 75))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
			} else {
				err := tradeCache.SetStopLoss(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 75))
				if err != nil {
					log.Println("Stop loss failed to set at /plus300", err)
				}
			}
			trade, exists = tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			board := boardbuilders.BuildSlBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/zerosledit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			err := openpositionedit.UpdateSLInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, "0.00")
			if err != nil {
				log.Println("plus25edit: openpositionedit.UpdateTPInCache", err)
			}
			board := boardbuilders.BuildSlEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/minus10edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newSl := utils.CalculateStopLoss(openPrice, leverage, 10)
			newSlStr := strconv.FormatFloat(newSl, 'E', -1, 64)
			err := openpositionedit.UpdateSLInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newSlStr)
			if err != nil {
				log.Println("minus25edit: openpositionedit.UpdateSLInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/minus25edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newSl := utils.CalculateStopLoss(openPrice, leverage, 25)
			newSlStr := strconv.FormatFloat(newSl, 'E', -1, 64)
			err := openpositionedit.UpdateSLInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newSlStr)
			if err != nil {
				log.Println("minus25edit: openpositionedit.UpdateSLInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/minus33edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newSl := utils.CalculateStopLoss(openPrice, leverage, 33)
			newSlStr := strconv.FormatFloat(newSl, 'E', -1, 64)
			err := openpositionedit.UpdateSLInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newSlStr)
			if err != nil {
				log.Println("minus25edit: openpositionedit.UpdateSLInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/minus50edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newSl := utils.CalculateStopLoss(openPrice, leverage, 50)
			newSlStr := strconv.FormatFloat(newSl, 'E', -1, 64)
			err := openpositionedit.UpdateSLInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newSlStr)
			if err != nil {
				log.Println("minus25edit: openpositionedit.UpdateSLInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		case "/minus75edit":
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			// TODO
			perp, _ := ordercache.GetTradeSplit(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID)

			openPrice, _ := strconv.ParseFloat(perp.OpenPrice, 64)
			leverage, _ := strconv.ParseFloat(perp.Leverage, 64)

			newSl := utils.CalculateStopLoss(openPrice, leverage, 75)
			newSlStr := strconv.FormatFloat(newSl, 'E', -1, 64)
			err := openpositionedit.UpdateSLInCache(rdbPositionsPaper, strconv.FormatInt(guid, 10), trade.ActiveTradeID, newSlStr)
			if err != nil {
				log.Println("minus25edit: openpositionedit.UpdateSLInCache", err)
			}
			board := boardbuilders.BuildTpEditBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(tradeCache, rdbPrice, rdbPositionsPaper, trade.ActiveTradeID, guid)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		//case "/savesledit":
		//	trade, exists := tradeCache.Get(guid)
		//	if !exists {
		//		log.Println("User does not exist.")
		//	}
		//	newSl := utils.FloatToString(trade.OpenTradeSl)
		//
		//	err := api.UpdateSlGns(guid, strconv.Itoa(trade.OpenTradePairIndex), strconv.Itoa(trade.OpenTradeIndex), newSl, trade.Chain, trade.ActiveCollateral)
		//	if err != nil {
		//		log.Println("Error in update SL GNS /savesledit", err)
		//	}
		//	return
		//case "/savetpedit":
		//	trade, exists := tradeCache.Get(guid)
		//	if !exists {
		//		log.Println("User does not exist.")
		//	}
		//	newTp := utils.FloatToString(trade.OpenTradeTp)
		//
		//	err := api.UpdateSlGns(guid, strconv.Itoa(trade.OpenTradePairIndex), strconv.Itoa(trade.OpenTradeIndex), newTp, trade.Chain, trade.ActiveCollateral)
		//	if err != nil {
		//		log.Println("Error in update SL GNS /savetpedit", err)
		//	}
		//	return
		case "/forexpairs":
			//trade, exists := tradeCache.Get(guid)
			//if !exists {
			//	log.Println("User does not exist.")
			//}
			//user, exists := tradeCache.Get(guid)
			//if !exists {
			//	log.Println("User does not exist: ", guid)
			//}

			err := tradeCache.ResetPage(guid)
			if err != nil {
				log.Println("Error resetting page.")
			}

			board := boardbuilders.BuildForexBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildFxPairString(1)
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/cryptopairs":
			err := tradeCache.ResetPage(guid)
			if err != nil {
				log.Println("Error resetting page.")
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}

			board := boardbuilders.BuildPairsBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildCryptoPairString(int(trade.PairPage))
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}

			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		case "/commoditypairs":
			err := tradeCache.ResetPage(guid)
			if err != nil {
				log.Println("Error resetting page.")
			}
			//trade, exists := tradeCache.Get(guid)
			//if !exists {
			//	log.Println("User does not exist.")
			//}
			board := boardbuilders.BuildCommoditiesBoard(tradeCache, guid)
			updatedMessage := stringbuilders.BuildCommoditiesPairString()
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.DisableWebPagePreview = true
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			msg = msgEdit
		}
		if strings.Contains(callbackData, "/edittakeprofit+") {
			// Extract the tradeID from the Callback
			tradeID := utils.ExtractTradeIDPlusMethod(callbackData)

			var correctRedisInstance *redis.Client
			var isPaperStr string
			// Check if paper or real
			if strings.Contains(callbackData, "papertrade") {
				correctRedisInstance = rdbPositionsPaper
				isPaperStr = "papertrade"
			}
			if strings.Contains(callbackData, "realtrade") {
				correctRedisInstance = rdbPositionsPaper
				isPaperStr = "realtrade"

			}

			// Extract the Index
			//rdbTrade, _ := ordercache.GetTradeSplit(correctRedisInstance, strconv.FormatInt(guid, 10), tradeID)
			idxStr := utils.ExtractIndex(callbackData)
			idx, _ := strconv.Atoi(idxStr)

			board := boardbuilders.BuildTpBoardEdit(rdbPositionsPaper, rdbPositionsPaper, tradeCache, guid, tradeID, idxStr, isPaperStr)
			// Update the message content with the desired message
			updatedMessage := stringbuilders.BuildOpenPositionBoardString(correctRedisInstance, rdbPrice, idx, guid)
			// Create a new editable message with the updated text
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit

		}
		if strings.Contains(callbackData, "/editstoploss+") {
			// Extract the tradeID from the Callback
			tradeID := utils.ExtractTradeIDPlusMethod(callbackData)

			var correctRedisInstance *redis.Client
			var isPaperStr string
			// Check if paper or real
			if strings.Contains(callbackData, "papertrade") {
				correctRedisInstance = rdbPositionsPaper
				isPaperStr = "papertrade"
			}
			if strings.Contains(callbackData, "realtrade") {
				correctRedisInstance = rdbPositionsPaper
				isPaperStr = "papertrade"

			}
			//rdbTrade, _ := ordercache.GetTradeSplit(correctRedisInstance, strconv.FormatInt(guid, 10), tradeID)
			idxStr := utils.ExtractIndex(callbackData)
			idx, _ := strconv.Atoi(idxStr)

			board := boardbuilders.BuildSlBoardEdit(rdbPositionsPaper, rdbPositionsPaper, tradeCache, guid, tradeID, idxStr, isPaperStr)
			// Update the message content with the desired message
			updatedMessage := stringbuilders.BuildOpenPositionBoardString(correctRedisInstance, rdbPrice, idx, guid)
			// Create a new editable message with the updated text
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		}
		if strings.Contains(callbackData, "/refreshtpeditboard+") {
			// Extract the tradeID from the Callback
			tradeID := utils.ExtractTradeIDPlusMethod(callbackData)
			idx := utils.ExtractIndex(callbackData)
			idxInt, _ := strconv.Atoi(idx)

			var correctRedisInstance *redis.Client
			var isPaperStr string
			// Check if paper or real
			if strings.Contains(callbackData, "papertrade") {
				correctRedisInstance = rdbPositionsPaper
				isPaperStr = "papertrade"
			}
			if strings.Contains(callbackData, "realtrade") {
				correctRedisInstance = rdbPositionsPaper
				isPaperStr = "papertrade"

			}
			board := boardbuilders.BuildTpBoardEdit(rdbPositionsPaper, rdbPositionsPaper, tradeCache, guid, tradeID, idx, isPaperStr)
			// Update the message content with the desired message
			updatedMessage := stringbuilders.BuildOpenPositionBoardString(correctRedisInstance, rdbPrice, idxInt, guid)
			// Create a new editable message with the updated text
			msg = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msgEdit, ok := msg.(tgbotapi.EditMessageTextConfig)
			if !ok {
				log.Println("Cannot convert msg to type EditMessageTextConfig")
				return
			}
			msgEdit.ParseMode = tgbotapi.ModeMarkdown
			// Assign the new inline keyboard to the message
			msgEdit.ReplyMarkup = &board
			msg = msgEdit
		}
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
			errorhandling.HandleTelegramError(bot, err, update.CallbackQuery.ID)
		}
	} else if update.Message != nil && update.Message.ReplyToMessage == nil {
		if update.Message.Text != "/start" {
			userCache, exists := tradeCache.Get(update.Message.From.ID)
			if !exists {
				log.Println(userCache)
				alertMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please send the /start command again to re-initialize your session.")
				_, err := bot.Send(alertMsg)
				if err != nil {
					log.Println("Error sending alert message: ", err)
				}
				return
			}
		}

		var msg tgbotapi.MessageConfig
		boardMsg, ok := boardCache.Get(update.Message.From.ID)
		if !ok {
			log.Println("Message does not exist")
		}
		// Handle regular messages and initial board
		//var msg tgbotapi.MessageConfig
		guid := update.Message.From.ID
		//
		//// Prehandling to stop clogging this file
		//// Naive way of checking for opentrade or openorder [ real / paper ]
		//if len(update.Message.Text) >= 2 && len(update.Message.Text) <= 5 {
		//	textChars := update.Message.Text[:2]
		//	// Handle Open Positions
		//	if textChars == "/o" {
		//		trade, exists := tradeCache.Get(guid)
		//		if !exists {
		//			log.Println("User does not exist.")
		//		}
		//		var isArb bool
		//		if trade.Chain == "arbitrum" {
		//			isArb = true
		//		} else {
		//			isArb = false
		//		}
		//		publicKey, err := database.GetPublicKey(guid)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//		hasOrders, err := mongolisten.HasOpenOrders(client, publicKey, isArb)
		//		if err != nil {
		//			log.Println("Error Checking OpenTradesOrOrders: ", err)
		//			return
		//		}
		//		if hasOrders {
		//			//trade, _ := tradeCache.Get(int64(guid))
		//			// Delete the message containing the textChars
		//			utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
		//			idx, err := utils.ExtractNumber(update.Message.Text)
		//			if err != nil {
		//				log.Println("ExtractNumber error in /x handling", err)
		//			}
		//
		//			// We decrease this because the index is incremented by one because users do not read 0.
		//			idxDecr := idx - 1
		//			board := boardbuilders.BuildManageOrderOrTradeBoardGNS(client, rdbPrice, publicKey, idxDecr, false, true, guid, "polygon", isArb)
		//			updatedMessage := stringbuilders.BuildOpenPositionBoardStringGNS(tradeCache, client, rdbPrice, guid, publicKey, idxDecr, false, true, "polygon", false, isArb)
		//			finalMsg, err := utils.CreateEditMessage(board, updatedMessage, boardMsg.Chat.ID, boardMsg.MessageID, "Markdown")
		//			if err != nil {
		//				log.Println("Error in CreateEditMessage in case polygon")
		//			}
		//			_, err = bot.Send(finalMsg)
		//			if err != nil {
		//				log.Println("Error sending message:", err)
		//				return
		//			}
		//			return
		//		}
		//	}
		//	if textChars == "/t" {
		//		trade, exists := tradeCache.Get(guid)
		//		if !exists {
		//			log.Println("User does not exist.")
		//		}
		//		var isArb bool
		//		if trade.Chain == "arbitrum" {
		//			isArb = true
		//		} else {
		//			isArb = false
		//		}
		//		publicKey, err := database.GetPublicKey(guid)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//		hasTrades, err := mongolisten.HasOpenTrades(client, publicKey, isArb)
		//		if err != nil {
		//			log.Println("Error Checking OpenTradesOrOrders: ", err)
		//			return
		//		}
		//		var msgChattable tgbotapi.Chattable
		//		utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
		//
		//		if hasTrades {
		//			// Check if the active trade right now is
		//			isTrade, cacheErr := tradeCache.GetIsTrade(guid)
		//			if cacheErr != nil {
		//				log.Println("Could not retrievve IsTrade from tradeCache in /t: ", err)
		//			}
		//			if isTrade {
		//				// We decrease because we count from 1
		//				idx, err := utils.ExtractNumber(update.Message.Text)
		//				err = tradeCache.SetActiveOpenTrade(guid, idx-1)
		//				if err != nil {
		//					fmt.Println("Unable to SetActiveOpenTrade in /t")
		//				}
		//				if err != nil {
		//					log.Println("ExtractNumber error in /t handling", err)
		//				}
		//
		//				//idxDecr := idx - 1
		//				updatedMessage := stringbuilders.BuildActiveGainsTradeStringV2(client, tradeCache, rdbPositions, rdbPrice, guid, 0, isArb)
		//				board := boardbuilders.BuildActiveTradeBoardGNSV2(client, rdbPrice, 0, tradeCache, guid, isArb)
		//
		//				msgChattable = tgbotapi.NewEditMessageText(boardMsg.Chat.ID, boardMsg.MessageID, updatedMessage)
		//				// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
		//				msgEdit, ok := msgChattable.(tgbotapi.EditMessageTextConfig)
		//				if !ok {
		//					log.Println("Cannot convert msg to type EditMessageTextConfig")
		//					return
		//				}
		//				// Assign the new inline keyboard to the message
		//				msgEdit.ReplyMarkup = &board
		//				msgEdit.ParseMode = tgbotapi.ModeMarkdown
		//				msgChattable = msgEdit
		//
		//				_, err = bot.Send(msgChattable)
		//				if err != nil {
		//					log.Println("Error sending message: ", err)
		//				}
		//				return
		//			} else {
		//				log.Println("/t was selected but no trade was found in the cache")
		//				return
		//			}
		//		}
		//	}
		//
		//}
		switch update.Message.Text {
		case "/debug":
			//ordercache.GetTradeSplit(rdbPositionsPaper)
		case "/start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to Hoot Trade! \n\nClick on menu icon to access all of the bot’s features. \n\nYou can try it out in PRACTICE MODE or LIVE MODE\n\nWe auto-generated a fresh WALLET for you. \n\nBeta LIVE: Deposit on Arbitrum or Polygon only.\n\nYou can deposit DAI on Polygon or Arbitrum to start leverage trading\n\nSniping, Swapping, Bridging, Copytrading and Airdrop Farming is still under development but is being worked on in parallel.")
			msg.ParseMode = tgbotapi.ModeMarkdown
			//msg.ReplyMarkup = keyboards.GetMainMenuKeyboard()
			userCache, exists := tradeCache.Get(guid) // Use tradeCache's Get function to check for existence
			if !exists {
				tradeCache.InitUser(guid)
				err := database.InitUser(guid, update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
				if err != nil {
					fmt.Println("Error: InitUser Database ", guid)
				}
				err = tradesettings.InitTradeSettings(guid)
				if err != nil {
					fmt.Println("Error: InitUser Tradesettings ", guid)
				}
				cacheErr := SetTradeId(tradeCache, rdbPositionsPaper, guid)
				if cacheErr != nil {
					log.Println("Could not SetStartTradeOrOrder for guid: ", guid)
				}
			} else {
				log.Println("User already exists: ", userCache)
				//tradeCache.InitUser(guid)
			}
		case "/chart":
			b64, err := chartservice.GetChart1M("BTC/USD")
			if err != nil {
				log.Println(err)
			}
			err = chartservice.SendChart(bot, update.Message.Chat.ID, b64)
			if err != nil {
				log.Println(err)
			}
		case "/switch":
			err := tradeCache.ToggleRealPaper(guid)
			if err != nil {
				log.Println("ToggleRealPaper Failed!")
				errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error encountered switching to Paper")
				_, err := bot.Send(errMsg)
				if err != nil {
					log.Println("Error sending message for /switch")
				}
			}
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}

			utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
			return
		case "/leverage":
			board := boardbuilders.BuildNewTradeBoard(tradeCache, rdbPrice, guid)
			// Update the message content with the desired message
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist /leverage.")
				return
			}
			perpMenuString := stringbuilders.BuildNewTradeString(rdbPrice, guid, trade, trade.PairIndex)
			// Create a new editable message with the updated text
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, perpMenuString)
			// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			msg.ReplyMarkup = board
			msg.ParseMode = tgbotapi.ModeMarkdown

			// Delete the perp message
			msgToDelete := tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}
			_, err := bot.Request(msgToDelete)
			if err != nil {
				errorhandling.HandleTelegramError(bot, err, update.CallbackQuery.ID)
			}

			spoofMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Loading...")
			spoofMsg.ReplyMarkup = tgbotapi.ForceReply{
				ForceReply: false,
				Selective:  true,
			}
			spoofDel, err := bot.Send(spoofMsg)
			if err != nil {
				log.Println("/leverage error: ", err)
				errorhandling.HandleTelegramError(bot, err, update.CallbackQuery.ID)
			}

			utils.DeleteMessage(bot, spoofDel.Chat.ID, spoofDel.MessageID)
		case "/sniper":
			log.Println("sniper")
		case "/wallet":
			isSet, err := database.CheckIfSet(guid)
			if err != nil {
				fmt.Println(err)
				return
			}
			if isSet {
				pubkey, err := database.GetPublicKey(guid)
				if err != nil {
					log.Println("Error retrieving public key in wallet", err)
					return
				}
				updatedMessage, err := stringbuilders.BuildWalletMainString(guid)
				if err != nil {
					fmt.Println("Error building string for users: " + strconv.Itoa(int(guid)))
					fmt.Println(err)
				}
				// Create a new editable message with the updated text
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, updatedMessage)
				msg.ParseMode = tgbotapi.ModeMarkdown
				qrFileName, err := wallet.CreateQRWithLogo(pubkey)
				if err != nil {
					fmt.Println(err)
				}
				// Send the QR code as a photo
				photoMedia := wallet.NewInputMediaPhotoFromFile(qrFileName)

				photoMsg := tgbotapi.NewPhoto(guid, photoMedia.BaseInputMedia.Media)
				_, err = bot.Send(photoMsg)
				if err != nil {
					log.Println("Failed to send QR code photo:", err)
					return
				}
			} else {
				pubkey, pk := wallet.GenerateWallet()
				database.AddPublicKey(guid, pubkey)
				database.AddPrivateKey(guid, pk)

				updatedMessage, err := stringbuilders.BuildWalletMainString(guid)
				if err != nil {
					fmt.Println("Error building string for users: " + strconv.Itoa(int(guid)))
					fmt.Println(err)
				}
				// Create a new editable message with the updated text
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, updatedMessage)
				// This type assertion is needed because tgbotapi.NewEditMessageText returns a tgbotapi.MessageConfig, not a tgbotapi.EditMessageTextConfig.
			}
			// REALKEYBOARD HANDLING
		case "➡️":
			// Delete the received right arrow message
			msgToDelete := tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}
			_, err := bot.Request(msgToDelete)
			if err != nil {
				fmt.Println(err)
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyMarkup = keyboards.GeneratePairKeyboardCryptoExUSDThreeRowFiveCol(1)
		}
		// PAIRCATCHER : This is where we catch the pairs
		if utils.IsFirstThreeCaps(update.Message.Text) {
			pair := utils.ExtractPair(update.Message.Text)
			fmt.Println("Catching pairs: ")
			index := pairmaps.NakedFilteredPairToIndex[pair]
			guid = int64(guid)

			// Set the new pair in the cache
			err := tradeCache.SetPairIndex(int64(index), guid)
			if err != nil {
				log.Println("Error setting PairIndex: ", err)
			}

			// Retrieve the price of the new asset from Redis
			price, err := priceserver.GetPrice(rdbPrice, index)
			if err != nil {
				log.Println("Error retrieving the new price of the asset from Redis: ", err)
			}

			leverage, err := tradeCache.GetLeverage(guid)
			if err != nil {
				log.Println("Error retrieving leverage from tradeCache: ", err)
			}

			isLong, err := tradeCache.IsLong(guid)
			if err != nil {
				log.Println("Error retrieving isLong from tradeCache: ", err)
			}

			// Update the open-price in the cache
			err = tradeCache.SetEntryPrice(guid, price)
			if err != nil {
				log.Println("Error setting the entry price: ", err)
			}
			// Update the TP
			tp := utils.CalculateTakeProfit(price, float64(leverage), 900)
			err = tradeCache.SetTakeProfit(guid, tp)
			if err != nil {
				log.Println("Error setting take-profit", err)
			}

			// Update the SL
			sl := 0.00
			err = tradeCache.SetStopLoss(guid, sl)
			if err != nil {
				log.Println("Error setting stop-loss", err)
			}

			// Update the liq
			liq := utils.CalculateLiquidationPrice(price, float64(leverage), isLong)
			err = tradeCache.SetLiquidation(guid, liq)
			if err != nil {
				log.Println("Error setting liquidation", err)
			}

			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist: ", guid)
			}
			idx := 0
			err = UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
			return
		}
		if msg.Text != "" {
			boardMsg, err := bot.Send(msg)
			if err != nil {
				errorhandling.HandleTelegramError(bot, err, strconv.FormatInt(update.Message.Chat.ID, 10))
			} else {
				boardCache.Set(update.Message.From.ID, &boardMsg)
			}
		}

	} else if update.Message.ReplyToMessage != nil {
		guid := update.Message.From.ID

		boardMsg, ok := boardCache.Get(update.Message.From.ID)
		if !ok {
			log.Println("Message does not exist")
		}

		var responseText string // Declare a string to store the response message
		switch update.Message.ReplyToMessage.Text {
		case "Please reply to this message with your query:":
			query := update.Message.Text
			resultArr := search.SearchPairs(query)
			board := search.BuildSearchBoard()
			msgTxt := search.BuildSearchString(resultArr)

			finalMsg, err := utils.CreateEditMessage(board, msgTxt, boardMsg.Chat.ID, boardMsg.MessageID, "Markdown")
			if err != nil {
				log.Println("Error CreatingEditMessage message in query reply search: ", err)
			}

			utils.DeleteMessage(bot, boardMsg.Chat.ID, update.Message.ReplyToMessage.MessageID)
			utils.DeleteMessage(bot, boardMsg.Chat.ID, update.Message.MessageID)

			_, err = bot.Send(finalMsg)
			if err != nil {
				log.Println("Error with message Hoot! Hoot no comprendo")
			}

			return
		case "Please provide your order entry price, type \"market\" for market orders.":
			// Refresh the board
			trade, exists := tradeCache.Get(guid)
			if !exists {
				fmt.Println(guid)
				log.Println("User does not exist.")
			}
			marketStr := "market"
			if update.Message.Text == marketStr || strings.ToLower(update.Message.Text) == marketStr {
				priceReal, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
				if err != nil {
					log.Println("Market to entry price error", err)
				}
				err = tradeCache.SetEntryPrice(guid, priceReal)
				if err != nil {
					fmt.Println(err)
				}
				err = tradeCache.SetOrderTypeToMarket(guid)
				err = tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 10))
				if err != nil {
					fmt.Println("SetStopLoss error inside provide position size: ", err)
				}
				err = tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 900))
				if err != nil {
					fmt.Println("SetTakeProfit error inside provide position size: ", err)
				}
				responseText = "Entry price set to MARKET"
				msgSet := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				sleepMsg, err := bot.Send(msgSet)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				time.Sleep(1500 * time.Millisecond)

				// Delete the initial ForceReply message
				utils.DeleteMessage(bot, update.Message.ReplyToMessage.Chat.ID, update.Message.ReplyToMessage.MessageID)
				// Delete the user entry
				utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
				// Delete ENTRY PRICE SET TO MARKET
				utils.DeleteMessage(bot, sleepMsg.Chat.ID, sleepMsg.MessageID)

				// Refresh the board
				_, exist := tradeCache.Get(guid)
				if exist {
					log.Println("User does not exist.")
				}
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			}

			entryPrice := update.Message.Text
			entryPriceInt, err := strconv.Atoi(entryPrice)
			entryPriceFloat, err := strconv.ParseFloat(entryPrice, 64)
			if err != nil {
				// Parsing as float failed, try parsing as int
				entryPriceInt, err = strconv.Atoi(entryPrice)
				if err != nil {
					fmt.Println("entryPrice Custom Error: ", err)
					invalidMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hoot! Hoot no comprendo.")
					send, err := bot.Send(invalidMsg)
					if err != nil {
						log.Println("Error with message Hoot! Hoot no comprendo")
					}
					time.Sleep(3000 * time.Millisecond)
					// Delete the no comprendo
					utils.DeleteMessage(bot, update.Message.ReplyToMessage.Chat.ID, send.MessageID)
					// Delete the initial ForceReply message
					utils.DeleteMessage(bot, update.Message.ReplyToMessage.Chat.ID, update.Message.ReplyToMessage.MessageID)
					// Delete the user entry
					utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
					return
				} else {
					// Entry price is an integer
					err = tradeCache.SetEntryPrice(guid, float64(entryPriceInt))
					if err != nil {
						fmt.Println(err)
					}

					fmt.Printf("Entry price as integer: %d\n", entryPriceInt)
					err = tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(float64(entryPriceInt), float64(trade.Leverage), 10))
					if err != nil {
						fmt.Println("SetStopLoss error inside provide position size: ", err)
					}
					err = tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(float64(entryPriceInt), float64(trade.Leverage), 900))
					if err != nil {
						fmt.Println("SetTakeProfit error inside provide position size: ", err)
					}

					err = tradeCache.SetOrderTypeToLimit(guid)
					// Delete the old message "Please provide entry price
					msgToDelete := tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.ReplyToMessage.MessageID,
					}
					_, err := bot.Request(msgToDelete)
					if err != nil {
						fmt.Println(err)
					}

					// Delete the old reply to "Please provide entry price
					replyToDelete := tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.MessageID,
					}
					_, err = bot.Request(replyToDelete)
					if err != nil {
						fmt.Println(err)
					}

					responseText = "Entry price set"
					msgSet := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
					sleepMsg, err := bot.Send(msgSet)
					if err != nil {
						log.Printf("Error sending message: %v", err)
					}
					time.Sleep(2 * time.Second)
					// Delete the old message
					msgToDelete = tgbotapi.DeleteMessageConfig{
						ChatID:    sleepMsg.Chat.ID,
						MessageID: sleepMsg.MessageID,
					}
					_, err = bot.Request(msgToDelete)
					if err != nil {
						fmt.Println(err)
					}
					// Delete the Set message "Please provide entry price
					msgToDeleteTwo := tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.ReplyToMessage.MessageID,
					}
					_, err = bot.Request(msgToDeleteTwo)
					if err != nil {
						fmt.Println(err)
					}
					tradeUpdated, exist := tradeCache.Get(guid)
					if !exist {
						log.Println("User does not exist: ", guid)
					}
					idx := 0
					err = UpdateNewTradeOrActiveTrades(
						bot, client, tradeUpdated,
						rdbPrice, rdbPositionsPaper,
						tradeCache,
						boardMsg.Chat.ID, boardMsg.MessageID,
						idx, guid)
					if err != nil {
						log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
					}
					return
				}
			} else {
				// Entry price is a float
				fmt.Printf("Entry price as float: %.2f\n", entryPriceFloat)
				err = tradeCache.SetEntryPrice(guid, float64(entryPriceInt))
				err = tradeCache.SetOrderTypeToLimit(guid)
				err = tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(float64(entryPriceInt), float64(trade.Leverage), 10))
				if err != nil {
					fmt.Println("SetStopLoss error inside provide position size: ", err)
				}
				err = tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(float64(entryPriceInt), float64(trade.Leverage), 900))
				if err != nil {
					fmt.Println("SetTakeProfit error inside provide position size: ", err)
				}
				if err != nil {
					fmt.Println(err)
				}
				// Delete the old message "Please provide entry price
				msgToDelete := tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.ReplyToMessage.MessageID,
				}
				_, err := bot.Request(msgToDelete)
				if err != nil {
					fmt.Println(err)
				}

				// Delete the old reply to "Please provide entry price
				replyToDelete := tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				}
				_, err = bot.Request(replyToDelete)
				if err != nil {
					fmt.Println(err)
				}

				responseText = "Entry price set"
				msgSet := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				sleepMsg, err := bot.Send(msgSet)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				time.Sleep(2 * time.Second)
				// Delete the old message
				msgToDelete = tgbotapi.DeleteMessageConfig{
					ChatID:    sleepMsg.Chat.ID,
					MessageID: sleepMsg.MessageID,
				}
				_, err = bot.Request(msgToDelete)
				if err != nil {
					fmt.Println(err)
				}

				tradeUpdated, exist := tradeCache.Get(guid)
				if !exist {
					log.Println("User does not exist: ", guid)
				}
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, tradeUpdated,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			}

		case "Please provide position size":
			positionSize := update.Message.Text
			positionSizeInt, err := strconv.Atoi(positionSize)
			if err != nil {
				fmt.Println("positionSize Custom Error: ", err)
				responseText = "Invalid position size, please check your input! Your input was: " + update.Message.Text
			} else {
				err = tradeCache.SetPositionSize(guid, int64(positionSizeInt))

				trade, exists := tradeCache.Get(guid)
				if !exists {
					fmt.Println("User not found in cache")
				}
				err = tradeCache.SetStopLoss(guid, utils.CalculateStopLoss(trade.OpenPrice, float64(trade.Leverage), 10))
				if err != nil {
					fmt.Println("SetStopLoss error inside provide position size: ", err)
				}
				err = tradeCache.SetTakeProfit(guid, utils.CalculateTakeProfit(trade.OpenPrice, float64(trade.Leverage), 900))
				if err != nil {
					fmt.Println("SetTakeProfit error inside provide position size: ", err)
				}
				if err != nil {
					fmt.Println(err)
				}
				responseText = "Position size set!"

				msgBoard := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				respMsgDel, err := bot.Send(msgBoard)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				time.Sleep(500 * time.Millisecond)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, update.Message.ReplyToMessage.MessageID)
				utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, respMsgDel.MessageID)
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			}
		case "Please provide leverage":
			leverageFactor := update.Message.Text
			leverageFactorInt, err := strconv.Atoi(leverageFactor)
			if err != nil {
				fmt.Println("positionSize Custom Error: ", err)
				responseText = "Invalid leverage size, please check your input! Your input was: " + update.Message.Text
				utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
				msgError := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				msgErrDel, err := bot.Send(msgError)
				if err != nil {
					log.Println(msgError)
				}
				time.Sleep(2000 * time.Millisecond)
				utils.DeleteMessage(bot, msgErrDel.Chat.ID, msgErrDel.MessageID)
			} else {
				trade, exists := tradeCache.Get(guid)
				if !exists {
					log.Println("Trade does not exist.")
				}
				// Update leverage
				err = tradeCache.SetLeverage(guid, int64(leverageFactorInt))
				if err != nil {
					fmt.Println(err)
				}

				trade, ok := tradeCache.Get(update.Message.From.ID)
				if !ok {
					log.Println("Error getting trade.")
				}

				price, err := priceserver.GetPrice(rdbPrice, int(trade.PairIndex))
				if err != nil {
					log.Println("Error getting price.")
				}

				// Calculate take profit
				takeProfit := utils.CalculateTakeProfit(price, float64(leverageFactorInt), 900)
				// Update take profit
				err = tradeCache.SetTakeProfit(guid, takeProfit)
				if err != nil {
					fmt.Println(err)
				}

				// Calculate stop loss
				stopLoss := utils.CalculateStopLoss(price, float64(leverageFactorInt), 75)
				// Update stop loss
				err = tradeCache.SetStopLoss(guid, stopLoss)
				if err != nil {
					fmt.Println(err)
				}

				// Calculate liquidation
				liquidation := utils.CalculateLiquidationPrice(price, float64(leverageFactorInt), trade.Buy)
				// Update liquidation
				err = tradeCache.SetLiquidation(guid, liquidation)
				if err != nil {
					fmt.Println(err)
				}

				responseText = "Leverage set!"
				setMsg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				respMsgDel, err := bot.Send(setMsg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				time.Sleep(500 * time.Millisecond)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, update.Message.ReplyToMessage.MessageID)
				utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, respMsgDel.MessageID)
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			}
		case "Please provide stop loss":
			trade, ok := tradeCache.Get(update.Message.From.ID)
			if !ok {
				log.Println("Error getting trade.")
			}
			positionSize := update.Message.Text
			positionSizeInt, err := strconv.Atoi(positionSize)
			if err != nil {
				fmt.Println("stopLoss Custom Error: ", err)
				responseText = "Invalid stop loss, please check your input! Your input was: " + update.Message.Text
			} else {
				err = tradeCache.SetStopLoss(guid, float64(positionSizeInt))
				if err != nil {
					fmt.Println(err)
				}
				responseText = "Stop loss set!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				respMsgDel, err := bot.Send(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				time.Sleep(500 * time.Millisecond)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, update.Message.ReplyToMessage.MessageID)
				utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, respMsgDel.MessageID)
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			}
		case "Please provide take profit":
			trade, ok := tradeCache.Get(update.Message.From.ID)
			if !ok {
				log.Println("Error getting trade.")
			}
			positionSize := update.Message.Text
			positionSizeInt, err := strconv.Atoi(positionSize)
			if err != nil {
				fmt.Println("stopLoss Custom Error: ", err)
				responseText = "Invalid take profit, please check your input! Your input was: " + update.Message.Text
			} else {
				err = tradeCache.SetTakeProfit(guid, float64(positionSizeInt))
				if err != nil {
					fmt.Println(err)
				}
				responseText = "Take profit set!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				respMsgDel, err := bot.Send(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				time.Sleep(500 * time.Millisecond)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, update.Message.ReplyToMessage.MessageID)
				utils.DeleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
				utils.DeleteMessage(bot, respMsgDel.Chat.ID, respMsgDel.MessageID)
				idx := 0
				err = UpdateNewTradeOrActiveTrades(
					bot, client, trade,
					rdbPrice, rdbPositionsPaper,
					tradeCache,
					boardMsg.Chat.ID, boardMsg.MessageID,
					idx, guid)
				if err != nil {
					log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
				}
				return
			}
		}
		// Check if responseText has been set, then create the message and send it
		if responseText != "" {
			trade, exists := tradeCache.Get(guid)
			if !exists {
				log.Println("User does not exist.")
			}
			idx := 0
			err := UpdateNewTradeOrActiveTrades(
				bot, client, trade,
				rdbPrice, rdbPositionsPaper,
				tradeCache,
				boardMsg.Chat.ID, boardMsg.MessageID,
				idx, guid)
			if err != nil {
				log.Println("Error in UpdateNewTradeOrActiveTrades: ", err)
			}
			return
		}
	}
}
