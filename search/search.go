package search

import "strings"

func SearchPairs(query string) []Pair {
	var matches []Pair

	query = strings.ToUpper(query)

	// Iterate over the pairs and find matches
	for _, pair := range GainsSearchablePairs {
		base := strings.ToUpper(pair.Base)

		// Check if the query is a prefix of the base
		if strings.HasPrefix(base, query) {
			matches = append(matches, pair)
		}

		// Break if we found 10 matches
		if len(matches) >= 10 {
			break
		}
	}
	return matches
}
