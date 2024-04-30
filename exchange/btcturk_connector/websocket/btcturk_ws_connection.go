package websocket

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"genericAPI/api/api_config"
	"genericAPI/exchange/btcturk_connector/constants"
	"genericAPI/exchange/btcturk_connector/dto"
	"genericAPI/internal/websocketclient"
	"log"
	"time"
)

type BtcturkWsConnection struct {
	wsClient *websocketclient.WebsocketClient
	respChan chan dto.WsResponse
}

func (b *BtcturkWsConnection) listen() {
	for {
		msg, err := b.wsClient.ReadMessage()
		if err != nil {
			continue
		}
		var data dto.WsResponse
		if bytes.Equal(msg[1:3], orderBookFull) {
			data = data.(dto.OrderBookFullResponse)
		} else if bytes.Equal(msg[1:3], subscriptionResponse) {
			data = data.(dto.WsSubscriptionResponse)
		}

		err = json.Unmarshal(msg[3:len(msg)-1], &data)
		if err != nil {
			log.Printf("Error while unmarshalling btcturk ws response:%s", err.Error())
			continue
		}
		b.respChan <- data
	}
}

func (b *BtcturkWsConnection) Authenticate() {
	key, err := base64.StdEncoding.DecodeString(api_config.Config.Btcturk.PrivateKey)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	stamp := fmt.Sprint(time.Now().UTC().UnixMilli())
	nonce := 3000
	hash := hmac.New(sha256.New, key)
	hash.Write([]byte(fmt.Sprint(api_config.Config.Btcturk.PublicKey, nonce)))
	signature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	request := fmt.Sprintf(`[114,{"type":114, "publicKey":"%s", "timestamp":%s, "nonce":3000, "signature": "%s"}]`, api_config.Config.Btcturk.PublicKey, stamp, signature)
	message := []byte(request)
	err = b.wsClient.WriteMessage(message)
	if err != nil {
		log.Printf("error writing authentication message on btcturk ws:%s", err.Error())
		return
	}
}

func (b *BtcturkWsConnection) SubscribeToOrderbookStream(pairSymbol string) {
	msg, _ := json.Marshal(dto.WsSubscriptionMessage{
		Type:    constants.Subscription,
		Channel: constants.Orderbook,
		Event:   pairSymbol,
		Join:    true,
	})
	b.wsClient.WriteMessage(msg)
}

func NewBtcturkWsConnection(respChan chan dto.WsResponse) (*BtcturkWsConnection, error) {
	client := websocketclient.NewWebsocketClient("BTCTURK_WS_CLIENT", constants.BaseWsUrl)
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	wsConn := BtcturkWsConnection{
		wsClient: client,
		respChan: respChan,
	}
	go wsConn.listen()
	return &wsConn, nil
}
