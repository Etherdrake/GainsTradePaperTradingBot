package types

type DexScrSearchJSON struct {
	SchemaVersion string `json:"schemaVersion"`
	Pairs         []struct {
		ChainID     string   `json:"chainId"`
		DexID       string   `json:"dexId"`
		URL         string   `json:"url"`
		PairAddress string   `json:"pairAddress"`
		Labels      []string `json:"labels"`
		BaseToken   struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"baseToken"`
		QuoteToken struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"quoteToken"`
		PriceNative string `json:"priceNative"`
		PriceUSD    string `json:"priceUsd"`
		Txns        struct {
			H1  TransactionData `json:"h1"`
			H24 TransactionData `json:"h24"`
			H6  TransactionData `json:"h6"`
			M5  TransactionData `json:"m5"`
		} `json:"txns"`
		Volume struct {
			H24 float64 `json:"h24"`
			H6  float64 `json:"h6"`
			H1  float64 `json:"h1"`
			M5  float64 `json:"m5"`
		} `json:"volume"`
		PriceChange struct {
			H1  float64 `json:"h1"`
			H24 float64 `json:"h24"`
			H6  float64 `json:"h6"`
			M5  float64 `json:"m5"`
		} `json:"priceChange"`
		Liquidity struct {
			USD   float64 `json:"usd"`
			Base  int64   `json:"base"`
			Quote float64 `json:"quote"`
		} `json:"liquidity"`
		FDV           int64 `json:"fdv"`
		PairCreatedAt int64 `json:"pairCreatedAt"`
	} `json:"pairs"`
}
