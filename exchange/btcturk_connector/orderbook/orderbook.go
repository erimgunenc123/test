package orderbook

import (
	http_endpoints2 "genericAPI/exchange/binanceconnector/http_endpoints"
	"genericAPI/exchange/common"
	"log"
	"strconv"
	"sync"
)

type Orderbook struct {
	symbol           string // e.g. BTCUSDT
	bids             *PriceLevelList
	asks             *PriceLevelList
	dataStream       chan []byte
	snapshotID       uint64
	snapshotDataLock sync.Mutex
}

func NewOrderbook(symbol string, fill bool) (ob *Orderbook, err error) {
	ob = &Orderbook{
		symbol:     symbol,
		bids:       &PriceLevelList{side: common.Bid},
		asks:       &PriceLevelList{side: common.Ask},
		snapshotID: 99999999999999999, // need to change this in 100 years
	}
	if err == nil {
		log.Printf("%s Orderbook initialized", symbol)
	}
	return
}

func (ob *Orderbook) GetSnapshot() map[string]any {
	var bids []common.Order
	var asks []common.Order
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		bids = ob.bids.GetAll()
	}()
	go func() {
		defer wg.Done()
		asks = ob.asks.GetAll()
	}()
	wg.Wait()
	return map[string]any{"bids": bids, "asks": asks}
}

func (ob *Orderbook) processBid(bids [][]string) {
	for _, bid := range bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		qty, _ := strconv.ParseFloat(bid[1], 64)
		ob.bids.Update(price, uint64(qty))
	}
}

func (ob *Orderbook) processAsk(asks [][]string) {
	for _, ask := range asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		qty, _ := strconv.ParseFloat(ask[1], 64)
		ob.asks.Update(price, uint64(qty))
	}
}

func (ob *Orderbook) processSnapshot(snapshot *http_endpoints2.OrderbookSnapshot) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, bid := range snapshot.Bids {
			price, _ := strconv.ParseFloat(bid[0], 64)
			qty, _ := strconv.ParseFloat(bid[1], 64)
			ob.bids.Insert(price, uint64(qty))
		}
	}()
	go func() {
		defer wg.Done()
		for _, ask := range snapshot.Asks {
			price, _ := strconv.ParseFloat(ask[0], 64)
			qty, _ := strconv.ParseFloat(ask[1], 64)
			ob.asks.Insert(price, uint64(qty))
		}
	}()
	wg.Wait()
}
