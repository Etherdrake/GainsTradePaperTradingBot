package api

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TradeResponse struct {
	IsOpen      bool                 `json:"is_open"`
	TradeResult OpenTradeQueryReturn `json:"trade_result"`
}

type OpenTradeQueryReturn struct {
	TraderAddress   common.Address
	PairIndex       *big.Int
	Index           *big.Int
	InitialPosToken *big.Int
	PositionSizeDai *big.Int
	OpenPrice       *big.Int
	Buy             bool
	Leverage        *big.Int
	TP              *big.Int
	SL              *big.Int
}
