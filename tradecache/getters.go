package tradecache

import "fmt"

func (tc *TradeCache) GetPairIndex(guid int64) (int64, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return userTrade.PairIndex, nil
}

func (tc *TradeCache) IsLong(guid int64) (bool, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return false, fmt.Errorf("user not found")
	}
	return userTrade.Buy, nil
}

func (tc *TradeCache) GetEntryPrice(guid int64) (float64, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return userTrade.OpenPrice, nil
}

func (tc *TradeCache) GetPositionSize(guid int64) (int64, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return userTrade.PositionSizeDai, nil
}

func (tc *TradeCache) GetLeverage(guid int64) (int64, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return userTrade.Leverage, nil
}

func (tc *TradeCache) GetTakeProfit(guid int64) (float64, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return userTrade.TP, nil
}

func (tc *TradeCache) GetStopLoss(guid int64) (float64, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return 0, fmt.Errorf("user not found")
	}
	return userTrade.SL, nil
}

func (tc *TradeCache) GetActiveTradeID(guid int64) (string, error) {
	tc.mx.Lock()
	defer tc.mx.Unlock()

	userTrade, exists := tc.Cache[guid]
	if !exists {
		return "", fmt.Errorf("user not found")
	}
	return userTrade.ActiveTradeID, nil
}
