package tradecache

import (
	"fmt"
	"sync"
)

type TradeCache struct {
	mx    sync.Mutex
	Cache map[int64]OpenTradeCache
}

var tradeCache = TradeCache{
	Cache: make(map[int64]OpenTradeCache),
}

func New() *TradeCache {
	return &TradeCache{
		Cache: make(map[int64]OpenTradeCache),
	}
}

func (tc *TradeCache) Get(key int64) (OpenTradeCache, bool) {
	tc.mx.Lock()
	defer tc.mx.Unlock()
	val, exists := tc.Cache[key]
	if !exists {
		fmt.Printf("Key %d not found in cache\n", key)
	}
	return val, exists
}

func (tc *TradeCache) Set(key int64, value OpenTradeCache) {
	tc.mx.Lock()
	defer tc.mx.Unlock()
	tc.Cache[key] = value
}

func (tc *TradeCache) Delete(key int64) {
	tc.mx.Lock()
	defer tc.mx.Unlock()
	delete(tc.Cache, key)
}
