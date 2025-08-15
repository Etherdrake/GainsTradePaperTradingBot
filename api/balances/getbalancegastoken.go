package balances

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// GetBalanceGasToken Function is working multichain -> February 7 / 2023
func GetBalanceGasToken(payload BalanceGasTokenPayload) (string, error) {
	// Serialize the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Send the payload to the Rust API
	resp, err := http.Post("http://localhost:3030/balancegastoken", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var balance string
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return "", err
	}

	return balance, nil
}
