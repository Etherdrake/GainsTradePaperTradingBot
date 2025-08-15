package gnsfees

type FeesResult struct {
	Borrowing struct {
		FeePerBlock float64 `json:"feePerBlock"`
		FeePerHour  float64 `json:"feePerHour"`
	} `json:"borrowing"`
	Fees      float64 `json:"fees"`
	AvgSpread float64 `json:"avgSpread"`
}
