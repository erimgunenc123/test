package common

type Side uint32

const (
	Ask Side = iota
	Bid
)

type Order struct { // todo might need side?
	Quantity   uint64  `json:"qty"`
	PriceLevel float64 `json:"price"`
}
