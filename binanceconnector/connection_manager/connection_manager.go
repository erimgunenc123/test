package connection_manager

import (
	"fmt"
	"genericAPI/binanceconnector/dto"
	"genericAPI/binanceconnector/websocket/binancewebsocket"
	"github.com/google/uuid"
	"sync"
	"time"
)

// WebSocket connections have a limit of 5 incoming messages per second. A message is considered:
// A PING frame
// A PONG frame
// A JSON controlled message (e.g. subscribe, unsubscribe)
// A connection that goes beyond the limit will be disconnected; IPs that are repeatedly disconnected may be banned.
// A single connection can listen to a maximum of 1024 streams.
// There is a limit of 300 connections per attempt every 5 minutes per IP.

type binanceConnectionManager struct {
	connections        map[string]*binancewebsocket.BinanceSocket // stream name as key. example:  btcusdt@depth -> conn
	connectionsLock    sync.Mutex
	connectionAttempts chan time.Time
}

var BinanceConnectionManager *binanceConnectionManager

func InitBinanceConnectionManager() {
	BinanceConnectionManager = &binanceConnectionManager{
		connections:        map[string]*binancewebsocket.BinanceSocket{},
		connectionsLock:    sync.Mutex{},
		connectionAttempts: make(chan time.Time, 299), // do not change 299
	}
	go BinanceConnectionManager.connectionAttemptListener()
}

func (bcm *binanceConnectionManager) Listen(symbol string, stream dto.BinanceStream, ch chan []byte) (err error) {
	if conn := bcm.GetConnection(fmt.Sprintf("%s@%s", symbol, stream)); conn != nil {
		conn.AddSubscriber(uuid.NewString(), ch)
		return nil
	}
	bcm.connectionAttempts <- time.Now() // this will block until the channel has at least 1 space
	conn, err := binancewebsocket.NewBinanceWebsocket(true, fmt.Sprintf("%s_%s_BINANCE_WS_CLIENT", symbol, stream))
	if err != nil {
		return err
	}
	conn.AddSubscriber(uuid.NewString(), ch)
	bcm.addConnection(fmt.Sprintf("%s@%s", symbol, stream), conn)
	return nil
}

func (bcm *binanceConnectionManager) GetConnection(stream string) *binancewebsocket.BinanceSocket {
	bcm.connectionsLock.Lock()
	defer bcm.connectionsLock.Unlock()
	if conn, ok := bcm.connections[stream]; ok {
		return conn
	}
	return nil
}

func (bcm *binanceConnectionManager) connectionAttemptListener() {
	for {
		lastAttempt := <-bcm.connectionAttempts // consume oldest element
		time.Sleep((5 * time.Minute) - time.Now().Sub(lastAttempt))
	}
}

func (bcm *binanceConnectionManager) addConnection(stream string, conn *binancewebsocket.BinanceSocket) {
	bcm.connectionsLock.Lock()
	defer bcm.connectionsLock.Unlock()
	bcm.connections[stream] = conn
}

func (bcm *binanceConnectionManager) removeConnection(stream string) {
	bcm.connectionsLock.Lock()
	defer bcm.connectionsLock.Unlock()
	delete(bcm.connections, stream)
}
