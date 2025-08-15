package constants

const (
	// Constants
	Chain = "polygon"

	// Chain
	Poly = "polygon"
	Arb  = "arbitrum"

	// MongoDB related constants
	ConnectionString = "mongodb://localhost:27017"
	DbNamePoly       = "gains_trades"
	DbNameArb        = "gains_trades_arb"
	CanceledOrders   = "canceled_orders"
	ClosedOrders     = "closed_orders"
	ClosedTrades     = "closed_trades"
	LimitTrades      = "limit_trades"
	MarketTrades     = "market_trades"
	CanceledActions  = "canceled_actions"
	OpenOrders       = "open_orders"

	ChainstackPolyHTTP                  = "https://nd-064-349-324.p2pify.com/2ca8de61471307d6c9a238bbfe23e52e"
	ChainstackPolyWSS                   = "wss://ws-nd-064-349-324.p2pify.com/2ca8de61471307d6c9a238bbfe23e52e"
	GainsTradingAddressPoly             = "0xb0901FEaD3112f6CaF9353ec5c36DC3DdE111F61"
	GainsTradingAddressArb              = "0x42069"
	GainsCallbackAddressPoly            = "0x82e59334da8c667797009bbe82473b55c7a6b311"
	GainsCallbackAddressArb             = "0x42069"
	GainsTradingAbiPath                 = "./gains_abi/GNSTradingV6_4_1_abi.json"
	GainsCallbackAbiPath                = "./gains_abi/GNSTradingCallbacksV6_4_1_abi.json"
	EventListenerTradingTelegramToken   = "6770906613:AAFCSQPy2jaZ_-SND15345D7tUIa8QRxL2I"
	EventListenerCallbacksTelegramToken = "6663124980:AAGBRx35m3YXnG43KLFQZssOTo-8v1qkYM4"
	EventlistenerAmalgamatedToken       = "6505441961:AAFIV05fIS5bCTrrQkL6GMFbb7mncDAKrUw"
	HootCommandControlID                = -1001805208006
	HootMarketExecutedID                = 44
	HootLimitExecuted                   = 48
	HootAnomalies                       = 57
	HootOpenLimit                       = 307
	HootOpenLimitUpdated                = 311
	HootTakeProfitStoploss              = 315
	HootNftOrderInitiated               = 319
	HootCouldNotClose                   = 323
	HootOpenLimitCanceled               = 327
	HootCallbackCancellations           = 331
	HootCallbackErrors                  = 498
	HootAmalgamatedHashmaps             = 5922
)
