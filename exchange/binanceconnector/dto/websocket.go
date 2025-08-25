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

type DepthStreamWsResponse struct {
	EventType       string     `json:"e"`
	EventTime       int64      `json:"E"`
	TransactionTime int64      `json:"T"`
	Symbol          string     `json:"s"`
	FirstUpdateID   uint64     `json:"U"`
	FinalUpdateID   uint64     `json:"u"`
	Pu              uint64     `json:"pu"` // Final update Id in last stream(ie `u` in last stream)
	Bids            [][]string `json:"b"`  // Bids to be updated
	Asks            [][]string `json:"a"`  // Asks to be updated
}
