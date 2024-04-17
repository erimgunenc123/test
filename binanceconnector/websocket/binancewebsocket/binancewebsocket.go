package binancewebsocket

import (
	"genericAPI/internal/customErrors"
	"genericAPI/internal/websocketclient"
	"sync"
)

type BinanceSocket struct {
	client         *websocketclient.WebsocketClient
	subscribers    map[string]chan []byte
	subscriberLock sync.Mutex
}

func NewBinanceWebsocket(start bool, clientName string) (socket *BinanceSocket, err error) {
	socket = &BinanceSocket{
		client:      websocketclient.NewWebsocketClient(clientName, BaseWsUrl),
		subscribers: make(map[string]chan []byte),
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
