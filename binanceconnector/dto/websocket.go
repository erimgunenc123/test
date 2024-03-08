package dto

type BinanceWsMethod string

const (
	MethodSubscribe   BinanceWsMethod = "SUBSCRIBE"
	MethodUnsubscribe BinanceWsMethod = "UNSUBSCRIBE"
)

type SymbolListenRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     uint64   `json:"id"`
}

type SymbolListenResponse struct {
	Result interface{} `json:"result"`
	Id     int         `json:"id"`
}
