package orderbook

import (
	"encoding/json"
	"errors"
	"genericAPI/binanceconnector/connection_manager"
	"genericAPI/binanceconnector/dto"
	"genericAPI/binanceconnector/http_endpoints"
	"log"
	"strconv"
	"sync"
)

type Orderbook struct {
	symbol           string // e.g. BTCUSDT
	bids             *BidPriceLevelList
	asks             *AskPriceLevelList
	dataStream       chan []byte
	snapshotID       uint64
	snapshotDataLock sync.Mutex
}

func NewOrderbook(symbol string, fill bool) (ob *Orderbook, err error) {
	ob = &Orderbook{
		symbol:     symbol,
		bids:       &BidPriceLevelList{},
		asks:       &AskPriceLevelList{},
		snapshotID: 99999999999999999, // need to change this in 100 years
	}
	if fill {
		err = ob.initOrderbook(false)
	}
	if err == nil {
		log.Printf("%s Orderbook initialized", symbol)
	}
	return
}

func (ob *Orderbook) GetSnapshot() map[string]any {
	var bids []dto.Order
	var asks []dto.Order
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		bids = ob.bids.GetAllBids()
	}()
	go func() {
		defer wg.Done()
		asks = ob.asks.GetAllAsks()
	}()
	wg.Wait()
	return map[string]any{"bids": bids, "asks": asks}
}

func (ob *Orderbook) initOrderbook(restart bool) error {
	if !restart {
		ob.dataStream = make(chan []byte, 10000)
		err := connection_manager.BinanceConnectionManager.Listen(ob.symbol, dto.DepthStream, ob.dataStream)
		if err != nil {
			return err
		}
	}
	go ob.listen(true)
	ob.snapshotDataLock.Lock()
	defer ob.snapshotDataLock.Unlock()
	snapshot, err := http_endpoints.GetOrderbookSnapshot(ob.symbol)
	if err != nil {
		return err
	}
	ob.snapshotID = snapshot.LastUpdateId
	ob.processSnapshot(snapshot)
	return nil
}

func (ob *Orderbook) listen(initial bool) {
	defer func() {
		if initial {
			go ob.listen(false)
		} else {
			go ob.initOrderbook(true)
		}

	}()
	if initial {
		wg := sync.WaitGroup{}
		for {
			ob.snapshotDataLock.Lock()
			msg := <-ob.dataStream
			var depthStreamResp dto.DepthStreamWsResponse
			err := json.Unmarshal(msg, &depthStreamResp)
			if err != nil {
				log.Printf("Failed unmarshalling depth stream response: %s", msg)
			}
			if depthStreamResp.Data.FirstUpdateID <= ob.snapshotID && depthStreamResp.Data.FinalUpdateID >= ob.snapshotID {
				wg.Add(2)
				go func() {
					defer wg.Done()
					ob.processAsk(depthStreamResp.Data.Asks)
				}()
				go func() {
					defer wg.Done()
					ob.processBid(depthStreamResp.Data.Bids)
				}()
				wg.Wait()
				return // process first message and call itself again without "initial = true"
			}
			ob.snapshotDataLock.Unlock()
		}
	} else {
		msgProcessor := ob.processSingleMessage()
		for {
			msg := <-ob.dataStream
			var depthStreamResp dto.DepthStreamWsResponse
			err := json.Unmarshal(msg, &depthStreamResp)
			if err != nil {
				log.Printf("Failed unmarshalling depth stream response: %s", msg)
			}
			err = msgProcessor(&depthStreamResp)
			if err != nil {
				log.Printf("%s orderbook order is broken. Restarting...", ob.symbol)
				return
			}
		}
	}
}

func (ob *Orderbook) processSingleMessage() func(msg *dto.DepthStreamWsResponse) error {
	var lastUpdateID uint64
	wg := sync.WaitGroup{}
	return func(msg *dto.DepthStreamWsResponse) error {
		if lastUpdateID != 0 && msg.Data.Pu != lastUpdateID {
			lastUpdateID = 0
			return errors.New("broken orderbook")
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			ob.processAsk(msg.Data.Asks)
		}()
		go func() {
			defer wg.Done()
			ob.processBid(msg.Data.Bids)
		}()
		wg.Wait()
		return nil
	}
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

func (ob *Orderbook) processSnapshot(snapshot *http_endpoints.OrderbookSnapshot) {
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
