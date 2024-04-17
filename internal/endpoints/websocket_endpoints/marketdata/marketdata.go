package marketdata

import (
	"encoding/json"
	"genericAPI/internal/common/constants"
	"genericAPI/internal/websocketclient"
	"github.com/gin-gonic/gin"
	"log"
)

func MarketdataWsHandler(c *gin.Context) {
	conn_, ok := c.Get(constants.ContextWebsocketConnectionKey)
	if !ok {
		panic("Ws upgrade middleware is not working properly.")
	}
	conn := conn_.(*websocketclient.WebsocketClient)
	for {
		msg, err := conn.ReadMessage()
		if err != nil {
			continue
		}
		go marketdataRequestHandler(msg)
	}
}

func marketdataRequestHandler(msg []byte) {
	var requestBody marketDataRequest
	if err := json.Unmarshal(msg, &requestBody); err != nil {
		log.Printf("Failed unmarshalling frontend request: %s", msg)
		return
	}
	switch requestBody.Action {
	case subscriptionRequest:
		subscriptionHandler(requestBody.Parameters)
	case unsubscriptionRequest:
		unsubscriptionHandler(requestBody.Parameters)
	default:
		log.Printf("Unknown action: %s", requestBody.Action)
		return
	}
}

func subscriptionHandler(params map[string]any) {

}

func unsubscriptionHandler(params map[string]any) {

}
