package binancewebsocket

import (
	"encoding/json"
	"fmt"
	"genericAPI/exchange/binanceconnector/dto"
	"genericAPI/internal/customErrors"
	"genericAPI/internal/websocketclient"
	"strings"
	"sync"
)

type BinanceSocket struct {
	client         *websocketclient.WebsocketClient
	subscribers    map[string]chan []byte
	subscriberLock sync.Mutex
	symbol         string
	stream         string
}

func NewBinanceWebsocket(start bool, symbol string, stream string) (socket *BinanceSocket, err error) {
	clientName := fmt.Sprintf("%s_%s_BINANCE_WS_CLIENT", symbol, stream)
	socket = &BinanceSocket{
		client:      websocketclient.NewWebsocketClient(clientName, BaseWsUrl),
		subscribers: make(map[string]chan []byte),
		symbol:      symbol,
		stream:      stream,
	}
	if start {
		err = socket.Start()
		if err != nil {
			return nil, err
		}
	}
	return socket, err
}

func (bs *BinanceSocket) Start() error {
	err := bs.client.Connect()
	if err != nil {
		return customErrors.ErrUnsuccessfulListenRequest
	}
	listenMsg, _ := json.Marshal(dto.SymbolListenRequest{
		Method: MethodSubscribe,
		Params: []string{
			fmt.Sprintf("%s@%s", strings.ToLower(bs.symbol), strings.ToLower(bs.stream)),
		},
		Id: 1,
	})
	bs.client.WriteMessage(listenMsg)
	_, _ = bs.client.ReadMessage()
	go bs.listen()
	return nil
}

func (bs *BinanceSocket) AddSubscriber(uuid string, sub chan []byte) {
	bs.subscriberLock.Lock()
	defer bs.subscriberLock.Unlock()
	bs.subscribers[uuid] = sub
}

func (bs *BinanceSocket) RemoveSubscriber(uuid string) {
	bs.subscriberLock.Lock()
	defer bs.subscriberLock.Unlock()
	delete(bs.subscribers, uuid)
}

func (bs *BinanceSocket) listen() {
	wg := sync.WaitGroup{}
	for {
		msg, err := bs.client.ReadMessage()
		if err != nil {
			continue
		}
		bs.subscriberLock.Lock()
		for _, subscriber := range bs.subscribers {
			wg.Add(1)
			go func(subChan chan []byte) {
				defer wg.Done()
				subChan <- msg
			}(subscriber)
		}
		wg.Wait()
		bs.subscriberLock.Unlock()
	}
}
