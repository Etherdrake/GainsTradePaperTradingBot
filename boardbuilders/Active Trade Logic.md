
1. User has no active trade -> no active trade in cache 
2. User opens board first -> trade[0] is active || order [0] is active
3. User selects trade -> trade[x] is active 
4. User selects order -> order[x] is active

**Build ActiveTradeStringV2 ALWAYS before ActiveTradeBoardV2**

We populate the cache on /start. 

For papertrading 

1. User has no active trade -> no active trade in cache
2. User opens board first -> trade[0].orderID is active || order[0].orderID is active 
3. User selects trade -> trade[x].orderID is active
4. User selects order -> order[x].orderID is active