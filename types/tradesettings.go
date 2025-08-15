package types

// TradeSetting represents the structure of a trade settings document in the collection.
type TradeSetting struct {
	UserID           int64   `bson:"userid"`
	EntryPrice       float64 `bson:"entryprice"`
	ActiveInstrument int     `bson:"active_instrument"`
	PositionSize     float64 `bson:"position_size"`
	Long             bool    `bson:"long"`
	Leverage         int     `bson:"leverage"`
	TakeProfit       uint64  `bson:"take_profit"`
	StopLoss         uint64  `bson:"stop_loss"`
}
