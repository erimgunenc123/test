package websocketclient

import (
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
)

type WebsocketClient struct {
	clientName string
	conn       *websocket.Conn
	mutex      sync.Mutex
	url        string
}

func NewWebsocketClient(clientName string, url string) *WebsocketClient {
	return &WebsocketClient{
		clientName: clientName,
		mutex:      sync.Mutex{},
		url:        url,
	}

}

func (w *WebsocketClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(w.url, nil) // http.Header{"X-MBX-APIKEY": []string{config.Cfg.TRBinance.Key}}
	if err != nil {
		return err
	}
	w.conn = conn
	return nil
}

func (w *WebsocketClient) GetClientName() string {
	return w.clientName
}

func (w *WebsocketClient) GetConnection() *websocket.Conn {
	return w.conn
}

func (w *WebsocketClient) SendPong() {
	return
}

func (w *WebsocketClient) SendPing() {
	if err := w.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
		slog.Error("Error sending ping:", err)
		return
	}
}

func (w *WebsocketClient) ReadMessage() []byte {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		slog.Error("Error reading message:", err)
		return nil
	}
	return message
}

func (w *WebsocketClient) WriteMessage(message []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, message)
}
