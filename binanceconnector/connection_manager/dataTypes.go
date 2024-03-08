package connection_manager

import (
	"fmt"
	"genericAPI/binanceconnector/binance_constants"
	"genericAPI/binanceconnector/dto"
	"genericAPI/binanceconnector/websocket/binancewebsocket"
	"genericAPI/internal/customErrors"
	"sync"
	"sync/atomic"
	"time"
)

// WebSocket connections have a limit of 5 incoming messages per second. A message is considered:
// A PING frame
// A PONG frame
// A JSON controlled message (e.g. subscribe, unsubscribe)
// A connection that goes beyond the limit will be disconnected; IPs that are repeatedly disconnected may be banned.
// A single connection can listen to a maximum of 1024 streams.
// There is a limit of 300 connections per attempt every 5 minutes per IP.
var globalIdentifier atomic.Uint64 // identifier for ws messages

type BinanceConnection struct {
	conn             *binancewebsocket.BinanceSocket
	openedAt         time.Time
	closedAt         time.Time
	owner            uint64     // userId
	listeningStreams *streamMap // depth, kline ... -> stream parameters
}

func (bc *BinanceConnection) subscribe(symbol dto.BinanceSymbol, stream dto.BinanceStream) error {
	return bc.conn.Subscribe(symbol, stream, globalIdentifier.Add(1))
}

func (bc *BinanceConnection) removeListeningStream(stream dto.BinanceStream) {
	bc.listeningStreams.remove(stream)
}

func (bc *BinanceConnection) addListeningStream(stream dto.BinanceStream, params dto.BinanceStreamParameters) error {
	if bc.listeningStreams.isFull() {
		return customErrors.ErrMaximumNumberOfConnections
	}
	bc.listeningStreams.add(stream, params)
	return nil
}

func (bc *BinanceConnection) getStartTime() time.Time {
	return bc.openedAt
}

func (bc *BinanceConnection) getOwner() uint64 {
	return bc.owner
}

func newBinanceConnection(userId uint64) (*BinanceConnection, error) {
	conn, err := binancewebsocket.NewBinanceWebsocket(true, fmt.Sprintf("%d_BINANCE_WS_CLIENT", userId))
	if err != nil {
		return nil, err
	}

	return &BinanceConnection{
		conn:     conn,
		openedAt: time.Now(),
		closedAt: time.Time{},
		owner:    userId,
		listeningStreams: &streamMap{
			streamMap: map[dto.BinanceStream]*dto.BinanceStreamParameters{},
			lock:      sync.Mutex{},
			size:      atomic.Uint32{},
		},
	}, nil
}

type streamMap struct {
	streamMap map[dto.BinanceStream]*dto.BinanceStreamParameters
	lock      sync.Mutex
	size      atomic.Uint32
}

func (s *streamMap) add(stream dto.BinanceStream, params dto.BinanceStreamParameters) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.streamMap[stream] = &params
	s.size.Add(1)
}

func (s *streamMap) remove(stream dto.BinanceStream) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.streamMap[stream]; ok {
		delete(s.streamMap, stream)
		s.size.Add(-1)
	}
}

func (s *streamMap) isFull() bool {
	return s.size.Load() == binance_constants.MaxBinanceWebsocketListeningSymbols
}

func (s *streamMap) getParameters(stream dto.BinanceStream) *dto.BinanceStreamParameters {
	s.lock.Lock()
	defer s.lock.Unlock()
	if params, ok := s.streamMap[stream]; ok {
		return params
	}
	return nil
}

func (s *streamMap) exists(stream dto.BinanceStream) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.streamMap[stream]
	return ok
}
