package chartservice

type Price struct {
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Timestamp int64   `json:"timestamp"`
}

type ApiResponse struct {
	Prices []Price `json:"prices"`
	Chart  string  `json:"chart"`
}

type Payload struct {
	Pair                string `json:"pair"`
	CandlesWidthSeconds int    `json:"candlesWidthSeconds"`
	NumCandles          int    `json:"numCandles"`
}
