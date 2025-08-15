package getters

import "fmt"

// GetStopLoss retrieves the "stop_loss" value for the given userID from the MongoDB collection.
func GetStopLoss(userID int64) (uint64, error) {
	tradeSetting, err := GetTradeSettingsByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get stop loss: %v", err)
	}
	if tradeSetting == nil {
		return 0, fmt.Errorf("no trade settings found for the given userID")
	}
	return tradeSetting.StopLoss, nil
}
