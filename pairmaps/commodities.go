package pairmaps

var CommodityToIndex = map[string]int{
	"XAU/USD": 90,
	"XAG/USD": 91,
	// Add more commodities if needed.
}

var IndexToCommodity = map[int]string{
	90: "XAU/USD",
	91: "XAG/USD",
	// Add more commodities if needed.
}
