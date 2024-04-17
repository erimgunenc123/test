package orderbook

import (
	binance_orderbook "genericAPI/binanceconnector/orderbook"
	"genericAPI/internal/services/marketdata/exchange_info"
	"log"
	"sync"
	"time"
)

type orderbookService struct {
	liveOrderbooks     map[string]*SymbolOrderbook
	liveOrderbooksLock sync.Mutex
}

var OrderbookService *orderbookService

func InitOrderbookService() {
	binanceSymbols := exchange_info.BinanceExchangeInfo.GetSymbols()
	btcTurkSymbols := exchange_info.BtcTurkExchangeInfo.GetSymbols()

	OrderbookService = &orderbookService{
		liveOrderbooks: make(map[string]*SymbolOrderbook),
	}
	wg := sync.WaitGroup{}
	for _, symbol := range binanceSymbols {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			if !contains(s, btcTurkSymbols) {
				return
			}
			obChan := make(chan *binance_orderbook.Orderbook)
			go func() {
				ob, err := binance_orderbook.NewOrderbook(s, true)
				if err != nil {
					log.Printf("Failed initializing orderbook for symbol: %s", s)
					return
				}
				obChan <- ob
			}()
			select {
			case _ = <-time.After(10 * time.Second):
				log.Printf("Orderbook %s timed out", s)
				return
			case ob := <-obChan:
				OrderbookService.liveOrderbooksLock.Lock()
				defer OrderbookService.liveOrderbooksLock.Unlock()
				OrderbookService.liveOrderbooks[s] = &SymbolOrderbook{BinanceOrderbook: ob} // todo
			}
		}(symbol)
	}
	wg.Wait()
}

func contains(str string, list []string) bool {
	for _, compareStr := range list {
		if compareStr == str {
			return true
		}
	}
	return false
}
