package boardbuilders

import (
	"HootTelegram/pairmaps"
	"HootTelegram/tradecache"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildPerpMainBoard(tradeCache *tradecache.TradeCache, rdbPrice *redis.Client, guid int64) tgbotapi.InlineKeyboardMarkup {
	trade, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var newTradeStr string
	var activeTradesStr string

	// Check on which Window we are:
	if trade.ActiveWindow == 0 { // We are on new trade
		newTradeStr = "⭐️ NEW TRADE"
		activeTradesStr = "ACTIVE TRADES"
	} else {
		if trade.ActiveWindow == 1 {
			newTradeStr = "NEW TRADE"
			activeTradesStr = "⭐️ ACTIVE TRADES"
		} else {
			log.Println("Error during BuildPerpMainBoard: ActiveWindow not equal to 0 or 1")
		}
	}

	pairStr := pairmaps.IndexToPair[int(trade.PairIndex)]

	var longShortButton tgbotapi.InlineKeyboardButton
	if trade.Buy {
		longShortButton = tgbotapi.NewInlineKeyboardButtonData("🟢 LONG", "/toggletoshort")
	} else {
		longShortButton = tgbotapi.NewInlineKeyboardButtonData("🔴 SHORT", "/toggletolong")
	}

	size := strconv.FormatInt(trade.PositionSizeDai, 10)
	leverage := strconv.FormatInt(trade.Leverage, 10)

	var priceReal string
	if trade.OrderType != 0 {
		priceReal = utils.FormatPrice(trade.OpenPrice)
		prePendStr := "$"
		priceReal = prePendStr + priceReal
	} else {
		priceReal = "MARKET"
	}

	var orderTypeStr string
	if trade.OrderType == 0 {
		orderTypeStr = "MARKET"
	}
	if trade.OrderType == 1 {
		orderTypeStr = "LIMIT"
	}
	if trade.OrderType == 2 {
		orderTypeStr = "STOP"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newTradeStr, "/leverage"),
			tgbotapi.NewInlineKeyboardButtonData(activeTradesStr, "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("REFRESH", "/refresh"),
			tgbotapi.NewInlineKeyboardButtonData(pairStr+" 🔽", "/pairs"),
			tgbotapi.NewInlineKeyboardButtonURL("CHART", "https://gains.trade"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take Profit", "/takeprofit"),
			longShortButton,
			tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/stoploss"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrprice"),
			tgbotapi.NewInlineKeyboardButtonData(priceReal, "/customprice"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrprice"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrpossize"),
			tgbotapi.NewInlineKeyboardButtonData("$"+size, "/custompossize"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrpossize"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrleverage"),
			tgbotapi.NewInlineKeyboardButtonData(leverage+"X", "/customleverage"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrleverage"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("SUBMIT "+orderTypeStr+" ("+strings.Trim(longShortButton.Text, "🟢🔴 ")+")", "/submitconfirm"),
		),
	)
}

func BuildPerpMainBoardSubmitCancel(tradeCache *tradecache.TradeCache, rdbPrice *redis.Client, guid int64) tgbotapi.InlineKeyboardMarkup {
	trade, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

	pairStr := pairmaps.IndexToPair[int(trade.PairIndex)]

	var longShortButton tgbotapi.InlineKeyboardButton
	if trade.Buy {
		longShortButton = tgbotapi.NewInlineKeyboardButtonData("🟢 LONG", "/toggletoshort")
	} else {
		longShortButton = tgbotapi.NewInlineKeyboardButtonData("🔴 SHORT", "/toggletolong")
	}

	//price, err := priceserver.GetPrice(rdbPrice, int(user.PairIndex))
	//if err != nil {
	//	log.Printf("Error fetching price for pairIndex %d: %v", user.PairIndex, err)
	//	// Handle this error: perhaps return or set a default price value
	//}

	var priceReal string
	if trade.OrderType != 0 {
		priceReal = utils.FormatPrice(trade.OpenPrice)
		prePendStr := "$"
		priceReal = prePendStr + priceReal
	} else {
		priceReal = "MARKET"
	}

	size := strconv.FormatInt(trade.PositionSizeDai, 10)
	leverage := strconv.FormatInt(trade.Leverage, 10)

	var newTradeStr string
	var activeTradesStr string

	// Check on which Window we are:
	if trade.ActiveWindow == 0 { // We are on new trade
		newTradeStr = "⭐️ NEW TRADE"
		activeTradesStr = "ACTIVE TRADES"
	} else {
		if trade.ActiveWindow == 1 {
			newTradeStr = "NEW TRADE"
			activeTradesStr = "⭐️ ACTIVE TRADES"
		} else {
			log.Println("Error during BuildPerpMainBoard: ActiveWindow not equal to 0 or 1")
		}
	}

	hyperlinkStr := utils.ReplaceSlashWithDash(pairStr)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newTradeStr, "/newtrade"),
			tgbotapi.NewInlineKeyboardButtonData(activeTradesStr, "/activetrades"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("REFRESH", "/refresh"),
			tgbotapi.NewInlineKeyboardButtonData(pairStr+" 🔽", "/pairs"),
			tgbotapi.NewInlineKeyboardButtonURL("CHART", "https://gains.trade/trading#"+hyperlinkStr), // https://gains.trade/trading#ETH-USD
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take Profit", "/takeprofit"),
			longShortButton,
			tgbotapi.NewInlineKeyboardButtonData("Stop Loss", "/stoploss"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrprice"),
			tgbotapi.NewInlineKeyboardButtonData(priceReal, "/customprice"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrprice"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrpossize"),
			tgbotapi.NewInlineKeyboardButtonData("$"+size, "/custompossize"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrpossize"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➖", "/decrleverage"),
			tgbotapi.NewInlineKeyboardButtonData(leverage+"X", "/customleverage"),
			tgbotapi.NewInlineKeyboardButtonData("➕", "/incrleverage"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ CANCEL", "/newtrade"),
			tgbotapi.NewInlineKeyboardButtonData("✅ CONFIRM", "/submit"),
		),
	)
}
