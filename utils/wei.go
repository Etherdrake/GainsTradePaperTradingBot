package utils

import (
	"fmt"
	"math/big"
)

func WeiToEther(weiStr string) (string, error) {
	wei, ok := new(big.Int).SetString(weiStr, 10)
	if !ok {
		return "", fmt.Errorf("Failed to parse input string as big.Int")
	}

	etherValue := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetFloat64(1e18))
	return etherValue.Text('f', 2), nil // Set precision to 2 decimal places
}
