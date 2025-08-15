package feeservice

type FeeData struct {
	Borrowing struct {
		FeePerBlock float64 `json:"feePerBlock"`
		FeePerHour  float64 `json:"feePerHour"`
	} `json:"borrowing"`
	Fees      int     `json:"fees"`
	AvgSpread float64 `json:"avgSpread"`
}
