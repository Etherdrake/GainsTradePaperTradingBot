# GainsTradePaperTradingBot
A completely simulated and working paper-trading bot for Telegram. This bot internally simulates the trading activity by processing orders aynchonously. 

This older iteration of the bot does not use Redis to store userdata, but instead uses an asyncronos-hashmap which does not scale properly and was later removed for a Redis-based solution.  
