package priceserver

import (
	"HootTelegram/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	AllPricesURI = "https://backend-pricing.eu.gains.trade/charts"
)

func GetHTTPSPriceCache() types.PriceCache {
	resp, err := http.Get(AllPricesURI)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data types.PriceData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error encountered during price retrieval: ", err)
	}

	var PriceCacheStruct = types.PriceCache{
		IndexToPriceDataHigh:  IndexHTTPSPriceDataHigh(data),
		IndexToPriceDataLow:   IndexHTTPSPriceDataLow(data),
		IndexToPriceDataOpen:  IndexHTTPSPriceDataOpen(data),
		IndexToPriceDataClose: IndexHTTPSPriceDataClose(data),
	}
	return PriceCacheStruct
}

func IndexHTTPSPriceDataHigh(data types.PriceData) types.IndexToPriceDataHigh {
	indexData := make(types.IndexToPriceDataHigh)

	for index, value := range data.Highs {
		indexData[index] = value
	}

	return indexData
}

func IndexHTTPSPriceDataLow(data types.PriceData) types.IndexToPriceDataLow {
	indexData := make(types.IndexToPriceDataLow)

	for index, value := range data.Lows {
		indexData[index] = value
	}

	return indexData
}

func IndexHTTPSPriceDataOpen(data types.PriceData) types.IndexToPriceDataOpen {
	indexData := make(types.IndexToPriceDataOpen)

	for index, value := range data.Opens {
		indexData[index] = value
	}

	return indexData
}

func IndexHTTPSPriceDataClose(data types.PriceData) types.IndexToPriceDataClose {
	indexData := make(types.IndexToPriceDataClose)

	for index, value := range data.Closes {
		indexData[index] = value
	}

	return indexData
}
