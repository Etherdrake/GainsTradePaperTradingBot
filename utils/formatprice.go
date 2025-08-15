package utils

import (
	"fmt"
	"math"
	"strconv"
)

func FormatPrice(price float64) string {
	// Check if the price is less than 1
	if price < 1 {
		return fmt.Sprintf("%.7f", price)
	}

	// Get the number of digits before the decimal point
	numDigits := int(math.Log10(price)) + 1

	// Format the number according to the rules you provided
	switch numDigits {
	case 1:
		return fmt.Sprintf("%.5f", price)
	case 2:
		return fmt.Sprintf("%.4f", price)
	case 3:
		return fmt.Sprintf("%.3f", price)
	case 4:
		return fmt.Sprintf("%.2f", price)
	case 5:
		return fmt.Sprintf("%.1f", price)
	default:
		// Return the price as-is or you can specify a default format
		return strconv.FormatFloat(price, 'f', -1, 64)
	}
}
