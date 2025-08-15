package types

import "HootTelegram/api"

type OpenPositions struct {
	ID        int64 `json:"_id"`
	Positions []api.OpenTradeJSON
}
