package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetOpenorderAPI(payload PayloadOpenOrder) (OpenOrdersString, error) {
	// Construct the URL with parameters
	url := "http://localhost:3030/" + "getopenorder"

	// Convert struct to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return OpenOrdersString{}, err
	}

	// Make the HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return OpenOrdersString{}, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return OpenOrdersString{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OpenOrdersString{}, err
	}

	// Unmarshal the response body into OpenLimitOrderQuery
	var openOrder OpenOrdersString
	err = json.Unmarshal(responseBody, &openOrder)
	if err != nil {
		return OpenOrdersString{}, err
	}

	return openOrder, nil
}
