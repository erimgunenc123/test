package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"genericAPI/api/api_config"
	"genericAPI/btcturk_connector/constants"
	"genericAPI/internal/websocketclient"
	"log"
	"time"
)

type BtcturkWsConnection struct {
	wsClient *websocketclient.WebsocketClient
	respChan chan []byte
}

func (b *BtcturkWsConnection) Listen() {
	for {
		msg, err := b.wsClient.ReadMessage()
		if err != nil {
			continue
		}
		b.respChan <- msg
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
		log.Println("ERROR:", err)
		return
	}
}

func NewBtcturkWsConnection(respChan chan []byte) (*BtcturkWsConnection, error) {
	client := websocketclient.NewWebsocketClient("BTCTURK_WS_CLIENT", constants.BaseWsUrl)
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	return &BtcturkWsConnection{
		wsClient: client,
		respChan: respChan,
	}, nil
}
