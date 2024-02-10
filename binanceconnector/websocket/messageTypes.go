package websocket

type SymbolListenRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

type SymbolListenResponse struct {
	Result interface{} `json:"result"`
	Id     int         `json:"id"`
}
