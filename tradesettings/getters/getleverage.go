package getters

import "fmt"

// GetLeverage retrieves the "leverage" value for the given userID from the MongoDB collection.
func GetLeverage(userID int64) (int, error) {
	tradeSetting, err := GetTradeSettingsByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get leverage: %v", err)
	}
	if tradeSetting == nil {
		return 0, fmt.Errorf("no trade settings found for the given userID")
	}
	return tradeSetting.Leverage, nil
}
