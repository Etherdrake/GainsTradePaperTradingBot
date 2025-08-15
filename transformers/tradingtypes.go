package transformers

// DoneEventTransform represents the Done event
type DoneEventTransform struct {
	Done bool `json:"done" bson:"done"`
}

// PausedEventTransform represents the Paused event
type PausedEventTransform struct {
	Paused bool `json:"paused" bson:"paused"`
}

// NumberUpdatedEventTransform represents the NumberUpdated event
type NumberUpdatedEventTransform struct {
	Name  string `json:"name" bson:"name"`
	Value int64  `json:"value" bson:"value"`
}

// BypassTriggerLinkUpdatedEventTransform represents the BypassTriggerLinkUpdated event
type BypassTriggerLinkUpdatedEventTransform struct {
	User   string `json:"user" bson:"user"`
	Bypass bool   `json:"bypass" bson:"bypass"`
}

// MarketOrderInitiatedEventTransform represents the MarketOrderInitiated event
type MarketOrderInitiatedEventTransform struct {
	OrderID   string `json:"order_id"  bson:"_id"`
	Trader    string `json:"trader" bson:"trader"`
	PairIndex int64  `json:"pair_index" bson:"pair_index"`
	Open      bool   `json:"open" bson:"open"`
}

// OpenLimitPlacedEventTransform represents the OpenLimitPlaced event
type OpenLimitPlacedEventTransform struct {
	Trader    string `json:"trader"  bson:"_id"`
	PairIndex int64  `json:"pair_index" bson:"pair_index"`
	Index     int64  `json:"index" bson:"index"`
}

// OpenLimitUpdatedEventTransform represents the OpenLimitUpdated event
type OpenLimitUpdatedEventTransform struct {
	Trader       string  `json:"trader"  bson:"_id"`
	PairIndex    int64   `json:"pair_index" bson:"pair_index"`
	Index        int64   `json:"index" bson:"index"`
	NewPrice     float64 `json:"new_price" bson:"new_price"`
	NewTp        float64 `json:"new_tp" bson:"new_tp"`
	NewSl        float64 `json:"new_sl" bson:"new_sl"`
	MaxSlippageP float64 `json:"max_slippage_p" bson:"max_slippage_p"`
}

// OpenLimitCanceledEventTransform represents the OpenLimitCanceled event
type OpenLimitCanceledEventTransform struct {
	Trader    string `json:"trader"  bson:"_id"`
	PairIndex int64  `json:"pair_index" bson:"pair_index"`
	Index     int64  `json:"index" bson:"index"`
}

// TpUpdatedEventTransform represents the TpUpdated event
type TpUpdatedEventTransform struct {
	Trader    string  `json:"trader"  bson:"_id"`
	PairIndex int64   `json:"pair_index" bson:"pair_index"`
	Index     int64   `json:"index" bson:"index"`
	NewTp     float64 `json:"new_tp" bson:"new_tp"`
}

// SlUpdatedEventTransform represents the SlUpdated event
type SlUpdatedEventTransform struct {
	Trader    string  `json:"trader"  bson:"_id"`
	PairIndex int64   `json:"pair_index" bson:"pair_index"`
	Index     int64   `json:"index" bson:"index"`
	NewSl     float64 `json:"new_sl" bson:"new_sl"`
}

// NftOrderInitiatedEventTransform represents the NftOrderInitiated event
type NftOrderInitiatedEventTransform struct {
	OrderID          int64  `json:"order_id"  bson:"_id"`
	Trader           string `json:"trader" bson:"trader"`
	PairIndex        int64  `json:"pair_index" bson:"pair_index"`
	ByPassesLinkCost bool   `json:"bypasses_link_cost" bson:"bypasses_link_cost"`
}

// ChainlinkCallbackTimeoutEventTransform represents the ChainlinkCallbackTimeout event
type ChainlinkCallbackTimeoutEventTransform struct {
	OrderID int64                        `json:"order_id"  bson:"_id"`
	Order   *PendingMarketOrderTransform `json:"order" bson:"order"`
}

// CouldNotCloseTradeEventTransform represents the CouldNotCloseTrade event
type CouldNotCloseTradeEventTransform struct {
	Trader    string `json:"trader"  bson:"_id"`
	PairIndex int64  `json:"pair_index" bson:"pair_index"`
	Index     int64  `json:"index" bson:"index"`
}

// PendingMarketOrderTransform represents a pending market order
type PendingMarketOrderTransform struct {
	Trade            TradeTransform `json:"trade"  bson:"_id"`
	Block            int64          `json:"block" bson:"block"`
	WantedPrice      float64        `json:"wanted_price" bson:"wanted_price"` // PRECISION
	SlippageP        float64        `json:"slippage_p" bson:"slippage_p"`     // PRECISION (%)
	SpreadReductionP float64        `json:"spread_reduction_p" bson:"spread_reduction_p"`
	TokenID          int            `json:"token_id" bson:"token_id"` // index in supportedTokens
}
