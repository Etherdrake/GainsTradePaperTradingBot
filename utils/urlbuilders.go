package utils

import (
	"fmt"
	"strings"
)

// getChartURLGains takes a pair string like "BTC/USD" and returns the corresponding URL.
func GetChartURLGains(pair string) string {
	// Replace the slash with a hyphen and concatenate with the base URL
	return fmt.Sprintf("https://gains.trade/trading#%s", strings.ReplaceAll(pair, "/", "-"))
}
