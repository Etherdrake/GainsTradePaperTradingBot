package search

import "HootTelegram/pairmaps"

func BuildSearchString(result []Pair) string {
	var searchResults string
	searchResults += "🔎 Your search results are: \n\n"
	for _, value := range result {
		base := value.Base
		quote := value.Quote
		mapString := base + "/" + quote
		idx := pairmaps.PairToIndex[mapString]
		fullName := pairmaps.SymbolToName[idx]
		resultString := "/" + base + ": " + fullName
		searchResults += resultString + "\n"
	}
	searchResults += "\n Press back to return to the Hoot trading board. "
	return searchResults
}
