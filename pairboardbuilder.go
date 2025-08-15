package main

import (
	"HootTelegram/pairmaps"
	"HootTelegram/trends"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildPairsBoard(pageNumber int) tgbotapi.InlineKeyboardMarkup {
	trendIndex := trends.IndexToTrend(GlobalPriceCache)

	const itemsPerPage = 8
	startIdx := (pageNumber - 1) * itemsPerPage

	// Top row
	topRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("FX", "/"),
		tgbotapi.NewInlineKeyboardButtonData("COMMODITIES", "/"),
	)

	catRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("PAIR:", "/"),
		tgbotapi.NewInlineKeyboardButtonData("PRICE:", "/"),
		tgbotapi.NewInlineKeyboardButtonData("TREND:", "/"),
	)

	// Rows with data
	var dataRows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsPerPage; i++ {
		idx := startIdx + i
		pair, _ := pairmaps.IndexToPair[idx]
		price := GlobalPriceCache.IndexToPriceDataClose[idx]
		trend := trendIndex[idx]
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(pair, "/"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("$%.2f", price), "/"),
			tgbotapi.NewInlineKeyboardButtonData(trend, "/"),
		)
		dataRows = append(dataRows, row)
	}

	// Bottom row
	bottomRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", "/"),
		tgbotapi.NewInlineKeyboardButtonData("Page "+strconv.Itoa(pageNumber), "/"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", "/"),
	)

	// Add all rows to the markup
	rows := append([][]tgbotapi.InlineKeyboardButton{topRow, catRow}, dataRows...)
	rows = append(rows, bottomRow)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
