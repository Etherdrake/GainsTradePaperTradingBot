package approval

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAllowance(chain, pubKey, srcToken, dstContract string) (*AllowanceResponse, error) {
	// Define the URL
	url := fmt.Sprintf("http://localhost:4000/api/v1/allowance/chain/%s/pubKey/%s/srcToken/%s/dstContract/%s",
		chain, pubKey, srcToken, dstContract)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// Read and parse the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the AllowanceResponse struct
	var allowanceResponse AllowanceResponse
	if err := json.Unmarshal(responseBody, &allowanceResponse); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return &allowanceResponse, nil
}
