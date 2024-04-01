package dto

type BinanceStream string

func (bs BinanceStream) ToString() string {
	return string(bs)
}

const (
	DepthStream BinanceStream = "depth"
)

type BinanceStreamParameters struct {
	Symbol     string
	Identifier uint64
	Parameters map[string]any // todo
}
