package revertreason

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RevertResponse struct {
	ErrorMessage string `json:"ErrorMessage"`
}

func GetRevertReason(chain, txHash string) (RevertResponse, error) {
	url := fmt.Sprintf("http://192.168.1.37:3003/api/v1/txError/chain/%s/txHash/%s", chain, txHash)

	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		return RevertResponse{}, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return RevertResponse{}, err
	}

	// Unmarshal the response body into RevertResponse struct
	var resp RevertResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return RevertResponse{}, err
	}

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		return RevertResponse{}, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	return resp, nil
}
