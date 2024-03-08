package connection_manager

import (
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
	connections        map[uint64]*BinanceConnection // userId -> conn
	connectionsLock    sync.Mutex
	connectionAttempts chan time.Time
}

var BinanceConnectionManager *binanceConnectionManager

func InitBinanceConnectionManager() {
	BinanceConnectionManager = &binanceConnectionManager{
		connections:        map[uint64]*BinanceConnection{},
		connectionsLock:    sync.Mutex{},
		connectionAttempts: make(chan time.Time, 300),
	}
	go BinanceConnectionManager.connectionAttemptListener()
}

func (bcm *binanceConnectionManager) NewConnection(userId uint64) (*BinanceConnection, error) {
	bcm.connectionAttempts <- time.Now() // this will block until the channel has at least 1 space
	conn, err := newBinanceConnection(userId)
	if err != nil {
		return nil, err
	}
	bcm.addConnection(userId, conn)
	return conn, nil
}

func (bcm *binanceConnectionManager) GetConnection(userId uint64) *BinanceConnection {
	bcm.connectionsLock.Lock()
	defer bcm.connectionsLock.Unlock()
	if conn, ok := bcm.connections[userId]; ok {
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

func (bcm *binanceConnectionManager) addConnection(userId uint64, conn *BinanceConnection) {
	bcm.connectionsLock.Lock()
	defer bcm.connectionsLock.Unlock()
	bcm.connections[userId] = conn
}

func (bcm *binanceConnectionManager) removeConnection(userId uint64) {
	bcm.connectionsLock.Lock()
	defer bcm.connectionsLock.Unlock()
	delete(bcm.connections, userId)
}
