package orderbook

import (
	"context"
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
	sequenceId       uint64
	snapshotDataLock sync.Mutex
	dbSink           *sink_service.QuestSinkService
}

func NewOrderbook(symbol string, fill bool, dbSink *sink_service.QuestSinkService) (ob *Orderbook, err error) {
	ob = &Orderbook{
		symbol:     symbol,
		orderbook:  NewHeapOrderbook(symbol),
		sequenceId: 99999999999999999, // need to change this in 100 years
		dbSink:     dbSink,
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
	time.Sleep(1 * time.Second) // shit solution but works
	ob.snapshotDataLock.Lock()
	defer ob.snapshotDataLock.Unlock()
	snapshot, err := http_endpoints.GetOrderbookSnapshot(ob.symbol)
	if err != nil {
		return err
	}
	err = ob.dbSink.InsertSingleRow(context.Background(),
		&models.OrderbookSnapshots{
			At:           time.Now(),
			Symbol:       ob.symbol,
			Bids:         snapshot.Bids,
			Asks:         snapshot.Asks,
			LastUpdateId: int64(snapshot.LastUpdateId),
		},
		"orderbook_snapshots",
	)
	if err != nil {
		slog.Info("Failed to insert snapshot", slog.AnyValue(err))
	}
	ob.sequenceId = snapshot.LastUpdateId
	ob.processSnapshot(snapshot)
	slog.Info(fmt.Sprintf("Snapshot initialized. Sequence id: %d", snapshot.LastUpdateId))
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
	slog.Info(fmt.Sprintf("Started listening with initial: %v", initial))
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
			if depthStreamResp.FinalUpdateID < ob.sequenceId {
				slog.Info(fmt.Sprintf("Snapshot not yet initialized. Discarding depth event with sequence id: %d", depthStreamResp.FinalUpdateID))
				ob.snapshotDataLock.Unlock()
				continue
			} else if depthStreamResp.FirstUpdateID <= ob.sequenceId && depthStreamResp.FinalUpdateID > ob.sequenceId {
				slog.Info("Received first message after snapshot initialization.")
				ob.sendOrderbookDeltasToSink(&depthStreamResp)
				ob.sequenceId = depthStreamResp.FinalUpdateID
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
				ob.snapshotDataLock.Unlock()
				slog.Info(fmt.Sprintf("Parsed message with Sequence id: %d", depthStreamResp.FirstUpdateID))
				return // process first message and call itself again without "initial = true"
			} else {
				panic(fmt.Sprintf("Invalid state. Last snapshot update id: %d, Received first update id: %d last update id: %d",
					ob.sequenceId, depthStreamResp.FirstUpdateID, depthStreamResp.FinalUpdateID))
			}
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
			ob.sendOrderbookDeltasToSink(&depthStreamResp)

		}
	}
}

func (ob *Orderbook) sendOrderbookDeltasToSink(deltas *binance_dtypes.DepthStreamWsResponse) {
	if ob.dbSink != nil {
		ob.dbSink.GetDataChan() <- &models.OrderbookDelta{
			At:            time.UnixMilli(deltas.EventTime),
			Symbol:        ob.symbol,
			Bids:          deltas.Bids,
			Asks:          deltas.Asks,
			EventType:     deltas.EventType,
			FirstUpdateId: int64(deltas.FirstUpdateID),
			FinalUpdateId: int64(deltas.FinalUpdateID),
		}
	}
}

func (ob *Orderbook) processSingleMessage() func(msg *binance_dtypes.DepthStreamWsResponse) error {
	wg := sync.WaitGroup{}
	return func(msg *binance_dtypes.DepthStreamWsResponse) error {
		if msg.FirstUpdateID != ob.sequenceId+1 {
			ob.sequenceId = 0
			return errors.New("broken orderbook")
		}
		ob.sequenceId = msg.FinalUpdateID
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
		ob.orderbook.Insert(true, price, qty)
	}
}

func (ob *Orderbook) processAsk(asks [][]string) {
	for _, ask := range asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		qty, _ := strconv.ParseFloat(ask[1], 64)
		ob.orderbook.Insert(false, price, qty)
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
