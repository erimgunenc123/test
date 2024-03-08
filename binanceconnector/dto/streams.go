package dto

type BinanceStream string

func (bs BinanceStream) ToString() string {
	return string(bs)
}

const (
	DepthStream BinanceStream = "depth"
)

type BinanceStreamParameters struct {
	firstSymbol  BinanceSymbol
	secondSymbol BinanceSymbol
	parameters   map[string]any // todo
}
