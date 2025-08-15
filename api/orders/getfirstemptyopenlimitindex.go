package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetFirstEmptyOpenLimitIndex(payload PayloadFirstEmptyOpenLimitIndex) (int64, error) {
	// Construct the URL with parameters
	url := "http://localhost:3030/" + "firstlimitslot"

	// Convert struct to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return -1, err
	}

	// Make the HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	// Unmarshal the response body into OpenLimitOrderQuery
	var firstEmptyOpenLimitIndex int64
	err = json.Unmarshal(responseBody, &firstEmptyOpenLimitIndex)
	if err != nil {
		return -1, err
	}

	return firstEmptyOpenLimitIndex, nil
}
