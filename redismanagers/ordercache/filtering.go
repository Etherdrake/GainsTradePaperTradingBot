package ordercache

import "HootTelegram/api"

func FilterOpenOrdersAndTrades(trades []api.OpenTradeJSON) []api.OpenTradeJSON {
	filteredTrades := make([]api.OpenTradeJSON, 0)

	for _, trade := range trades {
		if trade.OrderStatus == 2 || trade.OrderStatus == 3 {
			filteredTrades = append(filteredTrades, trade)
		}
	}

	return filteredTrades
}

func FilterOpenOrders(trades []api.OpenTradeJSON) []api.OpenTradeJSON {
	filteredOrders := make([]api.OpenTradeJSON, 0)

	for _, trade := range trades {
		if trade.OrderStatus == 2 {
			filteredOrders = append(filteredOrders, trade)
		}
	}

	return filteredOrders
}

func FilterOpenTrades(trades []api.OpenTradeJSON) []api.OpenTradeJSON {
	filteredTrades := make([]api.OpenTradeJSON, 0)

	for _, trade := range trades {
		if trade.OrderStatus == 3 {
			filteredTrades = append(filteredTrades, trade)
		}
	}

	return filteredTrades
}
