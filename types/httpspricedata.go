package types

type PriceData struct {
	Time   int64     `json:"time"`
	Opens  []float64 `json:"opens"`
	Highs  []float64 `json:"highs"`
	Lows   []float64 `json:"lows"`
	Closes []float64 `json:"closes"`
}

type IndexToPriceDataHigh map[int]float64

type IndexToPriceDataOpen map[int]float64

type IndexToPriceDataClose map[int]float64

type IndexToPriceDataLow map[int]float64
