package orderbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"genericAPI/exchange/binanceconnector/connection_manager"
	binance_dtypes "genericAPI/exchange/binanceconnector/dto"
	"genericAPI/exchange/binanceconnector/http_endpoints"
	"genericAPI/internal/qdb/quest/models"
	"genericAPI/internal/qdb/quest/sink_service"
	"log"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

type Orderbook struct {
	symbol           string // e.g. BTCUSDT
	orderbook        *HeapOrderbook
	dataStream       chan []byte
	snapshotID       uint64
	snapshotDataLock sync.Mutex
	dbSinkChan       chan sink_service.TableData
}

func NewOrderbook(symbol string, fill bool, dbSink *sink_service.QuestSinkService) (ob *Orderbook, err error) {
	ob = &Orderbook{
		symbol:     symbol,
		orderbook:  NewHeapOrderbook(symbol),
		snapshotID: 99999999999999999, // need to change this in 100 years
		dbSinkChan: dbSink.GetDataChan(),
	}
	if fill {
		err = ob.initOrderbook(false)
	}
	if err == nil {
		log.Printf("%s Orderbook initialized", symbol)
	}
	return
}

func (ob *Orderbook) initOrderbook(restart bool) error {
	if !restart {
		ob.dataStream = make(chan []byte, 10000)
		err := connection_manager.BinanceConnectionManager.Listen(ob.symbol, binance_dtypes.DepthStream, ob.dataStream)
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
	slog.Info("Snapshot initialized.")
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
	wg := sync.WaitGroup{}
	if initial {
		for {
			ob.snapshotDataLock.Lock()
			msg := <-ob.dataStream
			var depthStreamResp binance_dtypes.DepthStreamWsResponse
			err := json.Unmarshal(msg, &depthStreamResp)
			if err != nil {
				log.Printf("Failed unmarshalling depth stream response: %s", msg)
			}
			if depthStreamResp.FirstUpdateID <= ob.snapshotID && depthStreamResp.FinalUpdateID >= ob.snapshotID {
				wg.Add(2)
				go func() {
					defer wg.Done()
					ob.processAsk(depthStreamResp.Asks)
				}()
				go func() {
					defer wg.Done()
					ob.processBid(depthStreamResp.Bids)
				}()
				wg.Wait()
				ob.sendOrderbookToSink()
			} else {
				slog.Info("Snapshot not yet initialized. Discarding depth event...")
			}
			ob.snapshotDataLock.Unlock()
			return // process first message and call itself again without "initial = true"
		}
	} else {
		msgProcessor := ob.processSingleMessage()
		for {
			msg := <-ob.dataStream
			var depthStreamResp binance_dtypes.DepthStreamWsResponse
			err := json.Unmarshal(msg, &depthStreamResp)
			if err != nil {
				log.Printf("Failed unmarshalling depth stream response: %s", msg)
			}
			err = msgProcessor(&depthStreamResp)
			if err != nil {
				log.Printf("%s orderbook order is broken. Restarting...", ob.symbol)
				return
			}
			ob.sendOrderbookToSink()

		}
	}
}

func (ob *Orderbook) sendOrderbookToSink() {
	if ob.dbSinkChan != nil {
		bids, asks := ob.orderbook.GetAllAsList()
		slog.Info("Sent bids and asks to questdb!")
		ob.dbSinkChan <- &models.Orderbook{
			At:     time.Now(),
			Symbol: ob.symbol,
			Bids:   bids,
			Asks:   asks,
		}
	}
}
func (ob *Orderbook) processSingleMessage() func(msg *binance_dtypes.DepthStreamWsResponse) error {
	var lastUpdateID uint64
	wg := sync.WaitGroup{}
	return func(msg *binance_dtypes.DepthStreamWsResponse) error {
		if lastUpdateID != 0 && msg.Pu != lastUpdateID {
			lastUpdateID = 0
			return errors.New("broken orderbook")
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			ob.processAsk(msg.Asks)
		}()
		go func() {
			defer wg.Done()
			ob.processBid(msg.Bids)
		}()
		wg.Wait()
		return nil
	}
}

func (ob *Orderbook) processBid(bids [][]string) {
	for _, bid := range bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		qty, _ := strconv.ParseFloat(bid[1], 64)
		ob.orderbook.Update(true, price, qty)
	}
}

func (ob *Orderbook) processAsk(asks [][]string) {
	for _, ask := range asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		qty, _ := strconv.ParseFloat(ask[1], 64)
		ob.orderbook.Update(false, price, qty)
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
			ob.orderbook.Insert(true, price, qty)
		}
	}()
	go func() {
		defer wg.Done()
		slog.Info(fmt.Sprintf("Inserting asks from snapshot. Ask count: %d", len(snapshot.Asks)))
		for _, ask := range snapshot.Asks {
			price, _ := strconv.ParseFloat(ask[0], 64)
			qty, _ := strconv.ParseFloat(ask[1], 64)
			ob.orderbook.Insert(false, price, qty)
		}
	}()
	wg.Wait()
}
