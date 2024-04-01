package orderbook

import (
	binance_orderbook "genericAPI/binanceconnector/orderbook"
	"genericAPI/internal/services/marketdata/exchange_info"
	"log"
	"sync"
)

type orderbookService struct {
	liveOrderbooks     map[string]*binance_orderbook.Orderbook
	liveOrderbooksLock sync.Mutex
}

var OrderbookService *orderbookService

func InitOrderbookService() {
	allSymbols := exchange_info.ExchangeInfoService.GetSymbols()
	OrderbookService = &orderbookService{
		liveOrderbooks: make(map[string]*binance_orderbook.Orderbook, len(allSymbols)),
	}
	wg := sync.WaitGroup{}
	for _, symbol := range allSymbols {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			ob, err := binance_orderbook.NewOrderbook(s, true)
			if err != nil {
				log.Printf("Failed initializing orderbook for symbol: %s", s)
				return
			}
			OrderbookService.liveOrderbooksLock.Lock()
			defer OrderbookService.liveOrderbooksLock.Unlock()
			OrderbookService.liveOrderbooks[s] = ob
		}(symbol)
	}
	wg.Wait()
}
