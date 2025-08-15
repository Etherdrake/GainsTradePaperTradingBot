package main

import (
	"sync"
)

// Define a struct to represent the global cache
type GlobalCache struct {
	sync.RWMutex // Embed a RWMutex for thread-safety
	cache        map[int64]bool
}

// Function to check if a user exists in the global cache
func (gc *GlobalCache) exists(userID int64) bool {
	gc.RLock()
	defer gc.RUnlock()
	_, ok := gc.cache[userID]
	return ok
}

// Function to add a user to the global cache
func (gc *GlobalCache) addUser(userID int64) {
	gc.Lock()
	defer gc.Unlock()
	gc.cache[userID] = true
}

//func PriceServerHTTPS() {
//	ticker := time.NewTicker(5 * time.Second)
//	for _ = range ticker.C {
//		GlobalPriceCache = priceserver.GetHTTPSPriceCache()
//	}
//}
//
//type CurrentConfigPerps struct {
//	Pair string
//}
