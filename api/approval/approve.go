package approval

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// srcToken is the token you want to approve TO the dstContract. Ie. DAI to GNSTrading
func ApproveTokenToContract(guid, chain, amount, srcToken, dstContract string) (string, error) {
	// Define the URL
	url := fmt.Sprintf("http://localhost:4000/api/v1/approve/guid/%s/chain/%s/amount/%s/srctoken/%s/dstcontract/%s",
		guid, chain, amount, srcToken, dstContract)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// Read and parse the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Extract transaction hash from response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseMap); err != nil {
		return "", fmt.Errorf("error decoding response body: %v", err)
	}
	txHash := strings.ToLower(responseMap["txHash"].(string))

	return txHash, nil
}
