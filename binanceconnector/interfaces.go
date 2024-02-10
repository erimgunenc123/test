package binanceconnector

type ServiceInterface interface {
	GetSymbols() map[string]chan map[string]any
	Start() error
	AddSymbol(symbol string, sessionId string)
}
