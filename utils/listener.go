package utils

import (
	"fmt"
	"math/big"
	"strings"
)

// RemoveLeadingZeroes removes leading zeroes from a string
func RemoveLeadingZeroes(s string) string {
	return strings.TrimLeft(s, "0")
}

// RemoveLeadingZeroes removes leading zeroes from a hex-string while keeping the "0x" prefix
func RemoveLeadingZeroesHex(s string) string {
	// Check if the input has "0x" prefix
	hasPrefix := strings.HasPrefix(s, "0x")

	// Remove leading zeroes
	trimmed := strings.TrimLeft(s, "0")

	// Add "0x" prefix back if it was originally present
	if hasPrefix {
		return "0x" + trimmed
	}
	return trimmed
}

// BigToIntHexString converts a big.Int to a hex-string without leading zeros
func BigToIntHexString(value *big.Int) string {
	hexString := fmt.Sprintf("%x", value)
	return "0x" + hexString
}

func HexStringToInt64(hexStr string) (int64, error) {
	// Remove the "0x" prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Convert the hex string to a big.Int
	hexBigInt, success := new(big.Int).SetString(hexStr, 16)
	if !success {
		return 0, fmt.Errorf("invalid hex string: %s", hexStr)
	}

	// Extract the int64 value
	int64Value := hexBigInt.Int64()

	return int64Value, nil
}

func LowercaseString(input string) string {
	return strings.ToLower(input)
}
