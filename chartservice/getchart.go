package chartservice

import (
	"HootTelegram/constants"
	"bytes"
	"encoding/json"
	"net/http"
)

func GetChart1M(pair string) (string, error) {
	// Create payload with default values
	payload := Payload{
		Pair:                pair,
		CandlesWidthSeconds: 60,
		NumCandles:          100,
	}

	// Perform API request with payload
	apiResponse, err := PerformApiRequest(constants.CHARTS_URI+"", payload)
	if err != nil {
		return "", err
	}

	// Return only the chart string
	return apiResponse.Chart, nil
}

func GetChart5M(pair string) (string, error) {
	// Create payload with default values
	payload := Payload{
		Pair:                pair,
		CandlesWidthSeconds: 300,
		NumCandles:          100,
	}

	// Perform API request with payload
	apiResponse, err := PerformApiRequest(constants.CHARTS_URI+"", payload)
	if err != nil {
		return "", err
	}

	// Return only the chart string
	return apiResponse.Chart, nil
}

func GetChart15M(pair string) (string, error) {
	// Create payload with default values
	payload := Payload{
		Pair:                pair,
		CandlesWidthSeconds: 900,
		NumCandles:          100,
	}

	// Perform API request with payload
	apiResponse, err := PerformApiRequest(constants.CHARTS_URI+"", payload)
	if err != nil {
		return "", err
	}

	// Return only the chart string
	return apiResponse.Chart, nil
}

func GetChart1H(pair string) (string, error) {
	// Create payload with default values
	payload := Payload{
		Pair:                pair,
		CandlesWidthSeconds: 3600,
		NumCandles:          100,
	}

	// Perform API request with payload
	apiResponse, err := PerformApiRequest(constants.CHARTS_URI+"", payload)
	if err != nil {
		return "", err
	}

	// Return only the chart string
	return apiResponse.Chart, nil
}

func GetChart4H(pair string) (string, error) {
	// Create payload with default values
	payload := Payload{
		Pair:                pair,
		CandlesWidthSeconds: 14400,
		NumCandles:          100,
	}

	// Perform API request with payload
	apiResponse, err := PerformApiRequest(constants.CHARTS_URI+"", payload)
	if err != nil {
		return "", err
	}

	// Return only the chart string
	return apiResponse.Chart, nil
}

func GetChart1D(pair string) (string, error) {
	// Create payload with default values
	payload := Payload{
		Pair:                pair,
		CandlesWidthSeconds: 86400,
		NumCandles:          100,
	}

	// Perform API request with payload
	apiResponse, err := PerformApiRequest(constants.CHARTS_URI+"", payload)
	if err != nil {
		return "", err
	}

	// Return only the chart string
	return apiResponse.Chart, nil
}

func PerformApiRequest(uri string, payload Payload) (ApiResponse, error) {
	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return ApiResponse{}, err
	}

	// Make HTTP POST request
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return ApiResponse{}, err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return ApiResponse{}, err
	}
	return apiResponse, nil
}
