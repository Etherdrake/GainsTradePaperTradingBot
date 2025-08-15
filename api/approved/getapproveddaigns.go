package approved

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetApprovedDaiGNS(payload CheckAllowanceGNS) (string, error) {
	// Serialize the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Send the payload to the Rust API
	resp, err := http.Post("http://localhost:3030/getapprovedgns", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var txID string
	if err := json.NewDecoder(resp.Body).Decode(&txID); err != nil {
		return "", err
	}

	return txID, nil
}
