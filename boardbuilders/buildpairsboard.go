package boardbuilders

import (
	"HootTelegram/tradecache"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func BuildPairsBoard(tradeCache *tradecache.TradeCache, guid int64) tgbotapi.InlineKeyboardMarkup {
	user, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

	pairPage := strconv.Itoa(int(user.PairPage))

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "/decrpairpage"),
			tgbotapi.NewInlineKeyboardButtonData("PAGE "+pairPage, "/perps"),
			tgbotapi.NewInlineKeyboardButtonData("➡️", "/incrpairpage"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Forex", "/forexpairs"),
			tgbotapi.NewInlineKeyboardButtonData("Commodities", "/commoditypairs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔎 Search", "/search"),
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/leverage"),
		),
	)
}

func BuildForexBoard(tradeCache *tradecache.TradeCache, guid int64) tgbotapi.InlineKeyboardMarkup {
	user, exists := tradeCache.Get(guid)
	if !exists {
		fmt.Println("User not found in cache")
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var pairPageInt int
	if user.PairPage == 0 {
		pairPageInt = 1
	}

	pairPage := strconv.Itoa(pairPageInt)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "/decrpairpagefx"),
			tgbotapi.NewInlineKeyboardButtonData("PAGE "+pairPage, "/perps"),
			tgbotapi.NewInlineKeyboardButtonData("➡️", "/incrpairpagefx"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Crypto", "/pairs"),
			tgbotapi.NewInlineKeyboardButtonData("Commodities", "/commoditypairs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/perps"),
		),
	)
}

func BuildCommoditiesBoard(tradeCache *tradecache.TradeCache, guid int64) tgbotapi.InlineKeyboardMarkup {
	//user, exists := tradeCache.Get(guid)
	//if !exists {
	//	fmt.Println("User not found in cache")
	//	return tgbotapi.InlineKeyboardMarkup{}
	//}

	//pairPage := strconv.Itoa(int(user.PairPage))

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "/commoditypairs"),
			tgbotapi.NewInlineKeyboardButtonData("PAGE 1", "/perps"),
			tgbotapi.NewInlineKeyboardButtonData("➡️", "/commoditypairs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Forex", "/forexpairs"),
			tgbotapi.NewInlineKeyboardButtonData("Crypto", "/pairs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/perps"),
		),
	)
}
