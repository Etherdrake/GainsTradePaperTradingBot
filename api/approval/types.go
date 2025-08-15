package approval

type AllowanceResponse struct {
	Message        string `json:"message"`
	AllowanceWei   string `json:"allowance_wei"`
	AllowanceEther string `json:"allowance_ether"`
	PubKey         string `json:"pubKey"`
	SrcToken       string `json:"srcToken"`
	DstContract    string `json:"dstContract"`
	Chain          string `json:"chain"`
}
