package gnsfees

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetFeesAPI(chain string, isLong bool, posSize int64, lev, pairIndex int64, collat int, collatValue int64) (*FeesResult, error) {
	var collatName string
	if collat == 0 {
		collatName = "DAI"
	} else {
		if collat == 1 {
			collatName = "WETH"
		} else {
			if collat == 2 {
				collatName = "USDC"
			}
		}
	}

	// API endpoint
	apiURL := fmt.Sprintf("http://192.168.1.37:3002/api/v1/getFees/chain/%s/isLong/%t/collatName/%s/collatValue/%d/lev/%d/pairIndex/%d",
		chain, isLong, collatName, collatValue, lev, pairIndex)

	// Make the API call
	response, err := http.Get(apiURL)
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		// Parse the response JSON
		var result FeesResult
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		// Handle API error
		return nil, fmt.Errorf("Error: Unable to fetch data from the API. Status code: %d", response.StatusCode)
	}
}
