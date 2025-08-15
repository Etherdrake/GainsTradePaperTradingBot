# Database Model 

We have the following model for the database: 

# Goals 

# 1. Making sure the right pair is traded with the correct parameters for each user: 

MAIN CONTRACT: https://polygonscan.com/address/0x5c58a8acac721928b6f4495370c10762795d83da#writeContract

BSON Data for tradingconfig: 

userID_ leading key 
trader:"0x4Ff689D58fDda695fa2309ffB14DA170bc84cE74"
pairIndex:19
index: 0
inital: 0
positionSizeDai: 300000000000000000000 // 300 DAI 
openPrice: 6982580000 // 0.698258
buy: false 
leverage: 5
tp: 558606400
sl: 802996700

orderType: 0
spreadReductionId: 0
slippageP: 10131597140
referrer: 0x0000000000000000000000000000000000000000


# 2. Making sure open positions are shown correctly: 

How do we correctly query a contract to show us our open positions: 

MAIN CONTRACT: https://polygonscan.com/address/0xaee4d11a16B2bc65EDD6416Fb626EB404a6D65BD#readContract

1. Call function openTrades(trader address, pairIndex uint256, index uint256) 

openTrades returns: 

userID_



# 3.