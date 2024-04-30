package constants

type WsChannel int

const (
	Result              WsChannel = 100
	Request             WsChannel = 101
	UserLoginResult     WsChannel = 114
	Subscription        WsChannel = 151
	TickerAll           WsChannel = 401
	TickerPair          WsChannel = 402
	TradeSingle         WsChannel = 422
	UserTrade           WsChannel = 423
	OrderBookFull       WsChannel = 431
	OrderBookDifference WsChannel = 432
)

type WsChannelName string

const (
	Orderbook WsChannelName = "orderbook"
	Trade     WsChannelName = "trade "
	Ticker    WsChannelName = "ticker"
)
