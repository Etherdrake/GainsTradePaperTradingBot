package getters

import (
	"fmt"
)

// GetActiveInstrument retrieves the "active_instrument" value for the given userID from the MongoDB collection.
func GetActiveInstrument(userID int64) (int, error) {
	tradeSetting, err := GetTradeSettingsByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get active instrument: %v", err)
	}
	if tradeSetting == nil {
		return 0, fmt.Errorf("no trade settings found for the given userID")
	}
	return tradeSetting.ActiveInstrument, nil
}
