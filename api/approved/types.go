package approved

type CheckAllowanceGNS struct {
	TraderAddress string `json:"trader_address"`
	Chain         string `json:"chain"`
}

type ApproveDAI struct {
	Guid          int64  `json:"guid"`
	TraderAddress string `json:"trader_address"`
	Amount        string `json:"amount"`
	Chain         string `json:"chain"`
}
