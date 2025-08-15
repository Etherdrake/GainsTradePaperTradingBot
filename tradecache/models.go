package tradecache

// OpenTradeCache Changes to OpenTradeCache also need to be implemented in:
// - tradecache/initstruct.go
// - tradecache/convertopentradetojson.go
// - api/models.go
// - all redis functionality must be edited too!
type OpenTradeCache struct {
	ID              int64   `json:"_id"`
	TradeID         string  `json:"trade_id"`
	Trader          string  `json:"trader"`
	Paper           bool    `json:"paper"`
	PairIndex       int64   `json:"pair_index"`
	Index           int64   `json:"activeOrderIndex"`
	InitialPosToken int64   `json:"initial_pos_token"`
	PositionSizeDai int64   `json:"position_size_dai"`
	OpenPrice       float64 `json:"open_price"`
	ClosePrice      float64 `json:"close_price"`
	Buy             bool    `json:"buy"`
	Leverage        int64   `json:"leverage"`
	TP              float64 `json:"tp"`
	SL              float64 `json:"sl"`
	Liq             float64 `json:"liq"`
	PercentageTP    float64 `json:"percentagetp"`
	PercentageSL    float64 `json:"percentagesl"`
	OrderType       uint8   `json:"order_type"`
	OrderStatus     uint8   `json:"order_status"`
	PairPage        uint8   `json:"pair_page"`
	PnL             float64 `json:"pnl"`
	Chain           string  `json:"chain"`
	ActiveWindow    uint8   `json:"active_window"` // 0 is newtrade, 1 is active position window.
	ActiveTradeID   string  `json:"active_trade_id"`
}

var ordertypes = map[int]string{
	0: "MARKET",
	1: "LIMIT",
	2: "STOP",
}

var orderStatus = map[int]string{
	0: "NONE",
	1: "IN_PROGRESS",
	2: "OPEN_ORDER",
	3: "OPEN_TRADE",
	4: "CLOSED",
}
