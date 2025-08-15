package transformers

type MarketExecutedTransform struct {
	OrderID         int64          `json:"order_id" bson:"_id"`
	Trade           TradeTransform `json:"trade" bson:"trade"`
	Open            bool           `json:"open" bson:"open"`
	Price           float64        `json:"price" bson:"price"`
	PriceImpactP    float64        `json:"price_impact_p" bson:"price_impact_p"`
	PositionSizeDai float64        `json:"position_size_dai" bson:"position_size_dai"`
	PercentProfit   float64        `json:"percent_profit" bson:"percent_profit"`
	DaiSentToTrader float64        `json:"dai_sent_to_trader" bson:"dai_sent_to_trader"`
	CollateralToken string         `json:"collateral_token" bson:"collateral_token"`
}

type LimitExecutedTransform struct {
	OrderID         int64          `json:"order_id"  bson:"_id"`
	LimitIndex      int64          `json:"limit_index" bson:"limit_index"`
	Trade           TradeTransform `json:"trade" bson:"trade"`
	NftHolder       string         `json:"nft_holder" bson:"nft_holder"`
	OrderType       uint8          `json:"order_type" bson:"order_type"`
	Price           float64        `json:"price" bson:"price"`
	PriceImpactP    float64        `json:"price_impact_p" bson:"price_impact_p"`
	PositionSizeDai float64        `json:"position_size_dai" bson:"position_size_dai"`
	PercentProfit   float64        `json:"percent_profit" bson:"percent_profit"`
	DaiSentToTrader float64        `json:"dai_sent_to_trader" bson:"dai_sent_to_trader"`
	ExactExecution  bool           `json:"exact_execution" bson:"exact_execution"`
	CollateralToken string         `json:"collateral_token" bson:"collateral_token"`
}

type TradeTransform struct {
	Trader          string  `json:"trader" bson:"trader"`
	PairIndex       int64   `json:"pairIndex" bson:"pair_index"`
	Index           int64   `json:"index" bson:"index"`
	InitialPosToken float64 `json:"initialPosToken" bson:"initial_pos_token"`
	PositionSizeDai float64 `json:"positionSizeDai" bson:"position_size_dai"`
	OpenPrice       float64 `json:"openPrice" bson:"open_price"`
	Buy             bool    `json:"buy" bson:"buy"`
	Leverage        int64   `json:"leverage" bson:"leverage"`
	Tp              float64 `json:"tp" bson:"tp"`
	Sl              float64 `json:"sl" bson:"sl"`
	CollateralToken string  `json:"collateral_token" bson:"collateral_token"`
}

type MarketOpenCanceledTransform struct {
	OrderID         int64  `json:"order_id"  bson:"_id"`
	Trader          string `json:"trader" bson:"trader"`
	PairIndex       uint64 `json:"pair_index" bson:"pair_index"`
	CancelReason    uint8  `json:"cancel_reason" bson:"cancel_reason"`
	CollateralToken string `json:"collateral_token" bson:"collateral_token"`
}

type MarketCloseCanceledTransform struct {
	OrderID         int64  `json:"order_id" bson:"_id"`
	Trader          string `json:"trader" bson:"trader"`
	PairIndex       uint64 `json:"pair_index" bson:"pair_index"`
	Index           uint64 `json:"index" bson:"index"`
	CancelReason    uint8  `json:"cancel_reason" bson:"cancel_reason"`
	CollateralToken string `json:"collateral_token" bson:"collateral_token"`
}

type NftOrderCanceledTransform struct {
	OrderID         int64  `json:"order_id" bson:"order_id"`
	NftHolder       string `json:"nft_holder" bson:"nft_holder"`
	OrderType       uint8  `json:"order_type" bson:"order_type"`
	CancelReason    uint8  `json:"cancel_reason" bson:"cancel_reason"`
	CollateralToken string `json:"collateral_token" bson:"collateral_token"`
}

type MixedExecutedTransform struct {
	OrderID         int64          `json:"order_id" bson:"_id"`
	Trade           TradeTransform `json:"trade" bson:"trade"`
	Price           float64        `json:"price" bson:"price"`
	PriceImpactP    float64        `json:"price_impact_p" bson:"price_impact_p"`
	PositionSizeDai float64        `json:"position_size_dai" bson:"position_size_dai"`
	PercentProfit   float64        `json:"percent_profit" bson:"percent_profit"`
	DaiSentToTrader float64        `json:"dai_sent_to_trader" bson:"dai_sent_to_trader"`
	CollateralToken string         `json:"collateral_token" bson:"collateral_token"`

	// Fields specific to MarketExecutedTransform
	Open bool `json:"open" bson:"open"`

	// Fields specific to LimitExecutedTransform
	LimitIndex     int64  `json:"limit_index" bson:"limit_index"`
	NftHolder      string `json:"nft_holder" bson:"nft_holder"`
	OrderType      uint8  `json:"order_type" bson:"order_type"`
	ExactExecution bool   `json:"exact_execution" bson:"exact_execution"`

	// Additional fields from LimitExecutedTransform
	MinPrice int64 `json:"min_price" bson:"min_price"`
	MaxPrice int64 `json:"max_price" bson:"max_price"`
	Block    int64 `json:"block" bson:"block"`
	TokenID  int64 `json:"token_id" bson:"token_id"`
}
