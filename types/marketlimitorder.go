package types

// OrderType represents different order types
type OrderType int

const (
	MarketTrade OrderType = iota // 0
	LimitTrade                   // 1
	OpenOrder                    // 2
	// Add more order types if needed
)
