package marketdata

import (
	"encoding/json"
	tickers2 "genericAPI/exchange/btcturk_connector/tickers"
	"genericAPI/internal/common/constants"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

func BestBidAskWsHandler(c *gin.Context) {
	conn_, ok := c.Get(constants.ContextWebsocketConnectionKey)
	if !ok {
		panic("Ws upgrade middleware is not working properly.")
	}
	log.Printf("New ws connection received!")
	conn := conn_.(*websocket.Conn)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			continue
		}
		go bestBidAskHandler(conn, msg)
	}
}

func bestBidAskHandler(c *websocket.Conn, msg []byte) {
	var msgMap map[string]string
	json.Unmarshal(msg, &msgMap)
	if symbol, ok := msgMap["ticker"]; ok {
		go provideTickerFeed(c, symbol)
	}
}

func provideTickerFeed(c *websocket.Conn, symbol string) {
	readChan := make(chan *tickers2.Tick, 100)
	err := tickers2.NewTicker(symbol, readChan)
	if err != nil {
		log.Printf("Error while initializing ticker: %s", err.Error())
		return
	}
	go listenTicker(c, readChan)
}

func listenTicker(c *websocket.Conn, readChan chan *tickers2.Tick) {
	for {
		tick := <-readChan
		msgBytes, _ := json.Marshal(tick)
		c.WriteMessage(websocket.TextMessage, msgBytes)
	}
}
