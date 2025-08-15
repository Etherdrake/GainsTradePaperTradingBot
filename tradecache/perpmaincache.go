package tradecache

import (
	"HootTelegram/pairmaps"
	"HootTelegram/utils"
	"fmt"
)
import "HootTelegram/hooterrors"

func (tc *TradeCache) GenerateOrderID(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}
	userTrade.TradeID = utils.GenerateTradeID()

	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetOrderTypeToMarket(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.OrderType = 0
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetOrderTypeToLimit(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.OrderType = 1
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetOrderTypeToStop(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.OrderType = 2
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementPage(guid int64, class string) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	if userTrade.PairPage == 0 {
		userTrade.PairPage = 1
	}

	userTrade.PairPage++

	if len(pairmaps.FilteredIndexToCryptoPage(int(userTrade.PairPage), 10)) == 0 {
		userTrade.PairPage = 1
	}

	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) ResetPage(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.PairPage = 0

	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementPage(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	if userTrade.PairPage == 0 {
		userTrade.PairPage = 1
	}

	userTrade.PairPage--

	if userTrade.PairPage == 0 {
		userTrade.PairPage = 4
	}

	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetPairIndex(pairIdx, guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.PairIndex = pairIdx

	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementPairIndex(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	if userTrade.PairIndex == 0 {
		return hooterrors.ErrIndexZero
	}

	userTrade.PairIndex--

	for {
		if !pairmaps.IndexToDelisted[int(userTrade.PairPage)] {
			break
		}
		userTrade.PairPage--
	}

	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementPairIndex(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.PairIndex++

	for {
		if !pairmaps.IndexToDelisted[int(userTrade.PairIndex)] {
			break
		}
		userTrade.PairIndex++
	}
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) ToggleLongShort(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.Buy = !userTrade.Buy
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) ToggleRealPaper(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.Paper = !userTrade.Paper
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetEntryPrice(guid int64, entryPrice float64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.OpenPrice = entryPrice
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementStopLoss(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.SL = userTrade.SL * 0.9975
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementStopLoss(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.SL = userTrade.SL * 1.0025
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementTakeProfit(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.SL = userTrade.SL * 0.9975
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementTakeProfit(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.TP = userTrade.TP * 1.0025
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementEntryPrice(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.OpenPrice = userTrade.OpenPrice * 0.9975
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementEntryPrice(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}
	userTrade.OpenPrice = userTrade.OpenPrice * 1.0025
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementPositionSize(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.PositionSizeDai = userTrade.PositionSizeDai - 100
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetPositionSize(guid int64, positionSize int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.PositionSizeDai = positionSize
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementPositionSize(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.PositionSizeDai = userTrade.PositionSizeDai + 100
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) DecrementLeverage(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.Leverage = userTrade.Leverage - 10
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetLeverage(guid int64, leverage int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.Leverage = leverage
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) IncrementLeverage(guid int64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.Leverage = userTrade.Leverage + 10
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetTakeProfit(guid int64, takeProfit float64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.TP = takeProfit
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetStopLoss(guid int64, stopLoss float64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.SL = stopLoss
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetLiquidation(guid int64, liq float64) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}
	userTrade.Liq = liq
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetChain(guid int64, chain string) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.Chain = chain
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetActiveWindow(guid int64, window uint8) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.ActiveWindow = window
	tc.Cache[guid] = userTrade
	return nil
}

func (tc *TradeCache) SetActiveTradeID(guid int64, tradeID string) error {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return fmt.Errorf("user not found")
	}

	userTrade.ActiveTradeID = tradeID
	tc.Cache[guid] = userTrade
	return nil
}
