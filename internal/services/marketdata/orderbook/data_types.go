package orderbook

import binance_orderbook "genericAPI/binanceconnector/orderbook"

type SymbolOrderbook struct {
	BinanceOrderbook *binance_orderbook.Orderbook
	// btcturk orderbook todo
}
