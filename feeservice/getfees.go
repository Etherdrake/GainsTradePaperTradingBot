package feeservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GET http://localhost:3001/api/v1/getFees/chain/arb/isLong/false/collat/1000/lev/30/pairIndex/3
func GetFees(chain string, isLong bool, collat int, lev int, pairIndex int) (FeeData, error) {
	url := fmt.Sprintf("http://localhost:3001/api/v1/getFees/chain/%s/isLong/%t/collat/%d/lev/%d/pairIndex/%d", chain, isLong, collat, lev, pairIndex)

	response, err := http.Get(url)
	if err != nil {
		return FeeData{}, err
	}
	defer response.Body.Close()

	var GainsFees FeeData
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return FeeData{}, err
	}

	err = json.Unmarshal(body, &GainsFees)
	if err != nil {
		return FeeData{}, err
	}

	return GainsFees, nil
}
