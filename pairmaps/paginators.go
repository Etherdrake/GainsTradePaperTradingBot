package pairmaps

import "sort"

type KeyValue struct {
	Key   int
	Value string
}

func FilteredIndexToCryptoPage(pageNumber, itemsPerPage int) []string {
	// Create a slice of KeyValue structs from the map
	s := make([]KeyValue, 0, len(IndexToCrypto))
	for k, v := range IndexToCrypto {
		s = append(s, KeyValue{k, v})
	}

	// Sort the slice by key
	sort.Slice(s, func(i, j int) bool {
		return s[i].Key < s[j].Key
	})

	// Adjust page number for zero-based index calculation (if the first page is 1, not 0)
	zeroBasedPageNumber := pageNumber - 1

	// Calculate the starting index for the page in the sorted slice
	startIdx := zeroBasedPageNumber * itemsPerPage
	if startIdx > len(s) {
		return []string{} // Start index is beyond the slice, return empty
	}

	// Calculate the ending index in the sorted slice
	endIdx := startIdx + itemsPerPage
	if endIdx > len(s) {
		endIdx = len(s) // Adjust end index if it goes beyond the slice
	}

	// Extract the page of values
	page := make([]string, 0, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		page = append(page, s[i].Value)
	}
	return page
}

func FilteredIndexToFxPage(pageNumber, itemsPerPage int) []string {
	// Create a slice of KeyValue structs from the map
	s := make([]KeyValue, 0, len(IndexToFx))
	for k, v := range IndexToFx {
		s = append(s, KeyValue{k, v})
	}

	// Sort the slice by key
	sort.Slice(s, func(i, j int) bool {
		return s[i].Key < s[j].Key
	})

	// Adjust page number for zero-based index calculation (if the first page is 1, not 0)
	zeroBasedPageNumber := pageNumber - 1

	// Calculate the starting index for the page in the sorted slice
	startIdx := zeroBasedPageNumber * itemsPerPage
	if startIdx > len(s) {
		return []string{} // Start index is beyond the slice, return empty
	}

	// Calculate the ending index in the sorted slice
	endIdx := startIdx + itemsPerPage
	if endIdx > len(s) {
		endIdx = len(s) // Adjust end index if it goes beyond the slice
	}

	// Extract the page of values
	page := make([]string, 0, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		page = append(page, s[i].Value)
	}
	return page
}
