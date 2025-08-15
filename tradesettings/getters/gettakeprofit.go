package getters

import "fmt"

// GetTakeProfit retrieves the "take_profit" value for the given userID from the MongoDB collection.
func GetTakeProfit(userID int64) (uint64, error) {
	tradeSetting, err := GetTradeSettingsByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get take profit: %v", err)
	}
	if tradeSetting == nil {
		return 0, fmt.Errorf("no trade settings found for the given userID")
	}
	return tradeSetting.TakeProfit, nil
}
