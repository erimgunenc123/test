package dto

import (
	"encoding/json"
	"fmt"
	"genericAPI/exchange/btcturk_connector/constants"
)

type WsSubscriptionMessage struct {
	Type    constants.WsChannel     `json:"type"`
	Channel constants.WsChannelName `json:"channel"`
	Event   string                  `json:"event"`
	Join    bool                    `json:"join"`
}

// example message: [151,{"type":151,"channel":"orderbook","event":"BTCUSDT","join":true}]
func (w *WsSubscriptionMessage) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(w)
	prefix := fmt.Sprintf("[%d,", w.Type)
	b = append([]byte(prefix), append(b, []byte(`]`)...)...)
	return b, err
}

type WsResponse interface {
	GetType() int
}

type WsSubscriptionResponse struct {
	Type    int    `json:"type"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (r WsSubscriptionResponse) GetType() int {
	return r.Type
}

type OrderBookFullResponse struct {
	CS int    `json:"CS"`
	PS string `json:"PS"`
	AO []struct {
		A string `json:"A"`
		P string `json:"P"`
	} `json:"AO"`
	BO []struct {
		A string `json:"A"`
		P string `json:"P"`
	} `json:"BO"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Type    int    `json:"type"`
}

func (r OrderBookFullResponse) GetType() int {
	return r.Type
}
