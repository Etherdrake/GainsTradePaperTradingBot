package orders

import "math/big"

type PayloadOpenOrder struct {
	Trader    string `json:"trader"`
	PairIndex int    `json:"pairindex"`
	Index     int    `json:"index"`
	Chain     string `json:"chain"`
}

type PayloadOpenOrders struct {
	Chain string `json:"chain"`
}

type OpenOrderReturn struct {
	Trader           string   `json:"trader"`
	PairIndex        *big.Int `json:"pair_index"`
	Index            *big.Int `json:"index"`
	PositionSize     *big.Int `json:"position_size"`
	SpreadReductionP *big.Int `json:"spread_reduction_p"`
	Buy              bool     `json:"buy"`
	Leverage         *big.Int `json:"leverage"`
	Tp               *big.Int `json:"tp"`
	Sl               *big.Int `json:"sl"`
	MinPrice         *big.Int `json:"min_price"`
	MaxPrice         *big.Int `json:"max_price"`
	Block            *big.Int `json:"block"`
	TokenID          *big.Int `json:"token_id"`
}

type OpenOrdersString struct {
	Trader           string `json:"trader"`
	PairIndex        string `json:"pair_index"`
	Index            string `json:"index"`
	PositionSize     string `json:"position_size"`
	SpreadReductionP string `json:"spread_reduction_p"`
	Buy              bool   `json:"buy"`
	Leverage         string `json:"leverage"`
	Tp               string `json:"tp"`
	Sl               string `json:"sl"`
	MinPrice         string `json:"min_price"`
	MaxPrice         string `json:"max_price"`
	Block            string `json:"block"`
	TokenID          string `json:"token_id"`
}

type OpenOrdersBig struct {
	Trader           string   `json:"trader"`
	PairIndex        *big.Int `json:"pair_index"`
	Index            *big.Int `json:"index"`
	PositionSize     *big.Int `json:"position_size"`
	SpreadReductionP *big.Int `json:"spread_reduction_p"`
	Buy              bool     `json:"buy"`
	Leverage         *big.Int `json:"leverage"`
	Tp               *big.Int `json:"tp"`
	Sl               *big.Int `json:"sl"`
	MinPrice         *big.Int `json:"min_price"`
	MaxPrice         *big.Int `json:"max_price"`
	Block            *big.Int `json:"block"`
	TokenID          *big.Int `json:"token_id"`
}

type OpenOrdersHuman struct {
	Trader           string  `json:"trader" bson:"trader"`
	PairIndex        int64   `json:"pair_index" bson:"pair_index"`
	Index            int64   `json:"index" bson:"index"`
	PositionSize     float64 `json:"position_size" bson:"position_size"`
	SpreadReductionP float64 `json:"spread_reduction_p" bson:"spread_reduction_p"`
	Buy              bool    `json:"buy" bson:"buy"`
	Leverage         int64   `json:"leverage" bson:"leverage"`
	Tp               float64 `json:"tp" bson:"tp"`
	Sl               float64 `json:"sl" bson:"sl"`
	MinPrice         float64 `json:"min_price" bson:"min_price"`
	MaxPrice         float64 `json:"max_price" bson:"max_price"`
	Block            int64   `json:"block" bson:"block"`
	TokenID          int64   `json:"token_id" bson:"token_id"`
}

type PayloadFirstEmptyOpenLimitIndex struct {
	Trader    string `json:"trader"`
	PairIndex int    `json:"pair_index"`
	Chain     string `json:"chain"`
}
