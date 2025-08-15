package tradecache

import (
	"HootTelegram/database"
	"HootTelegram/priceserver"
	"HootTelegram/utils"
	"fmt"
)

func (tc *TradeCache) InitUser(guid int64) {
	// Lock the mutex before accessing the cache
	tc.mx.Lock()

	userPubKey, err := database.GetPublicKey(guid)
	if err != nil {
		fmt.Println(err)
	}

	OpenPrice := priceserver.GetHTTPSPriceCache().IndexToPriceDataOpen[0]

	percTP := 900.0
	percSL := 25.0

	TP := utils.CalculateTakeProfit(OpenPrice, 100.0, percTP)
	//SL := utils.CalculateStopLoss(OpenPrice, 100.0, percSL)

	Liq := utils.CalculateLiquidationPrice(OpenPrice, 100.0, true)

	tradeID := utils.GenerateTradeID()

	// Replace with your actual default values
	defaultOpenTradeCache := OpenTradeCache{
		ID:              guid,
		TradeID:         tradeID,
		Trader:          userPubKey,
		Paper:           true,
		PairIndex:       0,
		Index:           0,
		InitialPosToken: 0,
		PositionSizeDai: 1000,
		OpenPrice:       OpenPrice,
		Buy:             true,
		Leverage:        100,
		TP:              TP,
		SL:              0,
		Liq:             Liq,
		PercentageTP:    percTP,
		PercentageSL:    percSL,
		OrderType:       0,
		OrderStatus:     0,
		PairPage:        1,
		PnL:             0,
		Chain:           "arbitrum", // Change
		ActiveWindow:    0,
		ActiveTradeID:   "NOT_SET",
	}
	tc.mx.Unlock()
	tc.Set(guid, defaultOpenTradeCache)
}
