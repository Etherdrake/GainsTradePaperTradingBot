package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetOpenOrdersAPI(payload PayloadOpenOrders) ([]OpenOrdersString, error) {
	// Construct the URL with parameters
	url := "http://localhost:3030/" + "getopenlimitorders"

	// Convert struct to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return []OpenOrdersString{}, err
	}

	// Make the HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return []OpenOrdersString{}, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return []OpenOrdersString{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []OpenOrdersString{}, err
	}

	// Unmarshal the response body into OpenLimitOrderQuery
	var openOrders []OpenOrdersString
	err = json.Unmarshal(responseBody, &openOrders)
	if err != nil {
		return []OpenOrdersString{}, err
	}

	return openOrders, nil
}
