package binancewebsocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"genericAPI/binanceconnector/dto"
	"genericAPI/internal/customErrors"
	"genericAPI/internal/websocketclient"
	"log"
	"strconv"
	"sync"
)

// BinanceSocket Use binance socket to listen specific market streams
// How it works: Creates a websocket client to the base url, sends a subscription request for the demanded market stream
// endpoint, creates a response channel for that subscription, listens to the socket and writes the responses to
// their corresponding response channels based on provided identifiers
type BinanceSocket struct {
	client                     *websocketclient.WebsocketClient
	streamResponseChannels     map[uint64]chan []byte // identifier -> response channel
	streamResponseChannelsLock sync.Mutex
}

func NewBinanceWebsocket(start bool, clientName string) (socket *BinanceSocket, err error) {
	socket = &BinanceSocket{
		client:                     websocketclient.NewWebsocketClient(clientName, BaseWsUrl),
		streamResponseChannelsLock: sync.Mutex{},
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
		return err
	}
	go bs.listen()
	return nil
}

func (bs *BinanceSocket) GetResponseChannel(identifier uint64) chan []byte {
	bs.streamResponseChannelsLock.Lock()
	defer bs.streamResponseChannelsLock.Unlock()
	if respChan, ok := bs.streamResponseChannels[identifier]; ok {
		return respChan
	}
	return nil
}

func (bs *BinanceSocket) Subscribe(symbol dto.BinanceSymbol, stream dto.BinanceStream, identifier uint64) (err error) {
	requestBytes, _ := json.Marshal(
		dto.SymbolListenRequest{
			Method: MethodSubscribe,
			Params: []string{fmt.Sprintf("%s@%s", symbol, stream)},
			Id:     identifier,
		})

	respChan := bs.addResponseChan(identifier)
	defer func() {
		if err != nil {
			bs.removeResponseChan(identifier)
		}
	}()

	err = bs.client.WriteMessage(requestBytes)
	if err != nil {
		return err
	}

	var response map[string]any
	responseBytes := <-respChan
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}
	if res, ok := response["result"]; ok {
		if res == nil {
			log.Printf("Listen request successful. Symbol (%s) Stream(%s)", symbol, stream)
			return
		}
	}
	err = customErrors.ErrUnsuccessfulListenRequest
	return
}

func (bs *BinanceSocket) listen() {
	for {
		msg := bs.client.ReadMessage() // need to handle error cases, maybe restart the service subscription
		if msg != nil {
			idIdx := bytes.Index(msg, []byte(`"id"`)) // todo start from the end
			if idIdx == -1 {
				log.Printf("ID field not found in message: %s", msg)
				continue
			}

			colonIdx := bytes.Index(msg[idIdx:], []byte(":"))
			if colonIdx == -1 {
				log.Printf("Invalid format. message: %s", msg)
				continue
			}

			bracketIdx := bytes.IndexAny(msg[colonIdx:], "}")
			if bracketIdx == -1 {
				log.Printf("Invalid format. message: %s", msg)
				continue
			}
			chanIdentifier, err := strconv.Atoi(string(msg[idIdx+colonIdx : idIdx+colonIdx+bracketIdx]))
			if err != nil {
				log.Printf("Invalid channel identifier. message: %s", msg)
				continue
			}
			bs.streamResponseChannelsLock.Lock()
			if respChan, ok := bs.streamResponseChannels[uint64(chanIdentifier)]; ok {
				respChan <- msg
			} else {
				log.Printf("Identifier channel %d not found. Unsubscribing...", chanIdentifier)
				// todo unsub
			}
			bs.streamResponseChannelsLock.Unlock()
		}
	}
}

func (bs *BinanceSocket) addResponseChan(identifier uint64) chan []byte {
	bs.streamResponseChannelsLock.Lock()
	defer bs.streamResponseChannelsLock.Unlock()
	respChan := make(chan []byte, 1000) // faster than unbuffered
	bs.streamResponseChannels[identifier] = respChan
	return respChan
}

func (bs *BinanceSocket) removeResponseChan(identifier uint64) {
	bs.streamResponseChannelsLock.Lock()
	defer bs.streamResponseChannelsLock.Unlock()
	delete(bs.streamResponseChannels, identifier)
}
