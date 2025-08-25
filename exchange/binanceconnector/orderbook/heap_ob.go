package orderbook

import (
	"container/heap"
	"genericAPI/exchange/common"
	"time"
)

type Order struct {
	common.Order
	Index int
}
type BidsHeap []*Order

func (h BidsHeap) Len() int            { return len(h) }
func (h BidsHeap) Less(i, j int) bool  { return h[i].PriceLevel > h[j].PriceLevel } // max-heap
func (h BidsHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].Index, h[j].Index = i, j }
func (h *BidsHeap) Push(x interface{}) { *h = append(*h, x.(*Order)) }
func (h *BidsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type AsksHeap []*Order

func (h AsksHeap) Len() int            { return len(h) }
func (h AsksHeap) Less(i, j int) bool  { return h[i].PriceLevel < h[j].PriceLevel } // min-heap
func (h AsksHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i]; h[i].Index, h[j].Index = i, j }
func (h *AsksHeap) Push(x interface{}) { *h = append(*h, x.(*Order)) }
func (h *AsksHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type HeapOrderbook struct {
	At     time.Time
	Symbol string
	Bids   *BidsHeap
	Asks   *AsksHeap
}

func NewHeapOrderbook(symbol string) *HeapOrderbook {
	return &HeapOrderbook{
		At:     time.Now(),
		Symbol: symbol,
		Bids:   &BidsHeap{},
		Asks:   &AsksHeap{},
	}
}

func (o *HeapOrderbook) Insert(isBid bool, price, amount float64) {
	order := &Order{
		Order: common.Order{
			PriceLevel: price,
			Quantity:   amount,
		},
		Index: 0,
	}
	if isBid {
		heap.Push(o.Bids, order)
	} else {
		heap.Push(o.Asks, order)
	}
}

func (o *HeapOrderbook) Update(isBid bool, price, amount float64) {
	if isBid {
		newHeap := BidsHeap{}
		for _, ord := range *o.Bids {
			if ord.PriceLevel == price {
				if amount > 0 {
					ord.Quantity = amount
					newHeap = append(newHeap, ord)
				}
			} else {
				newHeap = append(newHeap, ord)
			}
		}
		o.Bids = &newHeap
		heap.Init(o.Bids)
	} else {
		newHeap := AsksHeap{}
		for _, ord := range *o.Asks {
			if ord.PriceLevel == price {
				if amount > 0 {
					ord.Quantity = amount
					newHeap = append(newHeap, ord)
				}
			} else {
				newHeap = append(newHeap, ord)
			}
		}
		o.Asks = &newHeap
		heap.Init(o.Asks)
	}
}

func (o *HeapOrderbook) GetAllAsList() (bids [][2]float64, asks [][2]float64) {
	bidCopy := make(BidsHeap, len(*o.Bids))
	copy(bidCopy, *o.Bids)
	askCopy := make(AsksHeap, len(*o.Asks))
	copy(askCopy, *o.Asks)

	heap.Init(&bidCopy)
	for bidCopy.Len() > 0 {
		copiedBid := heap.Pop(&bidCopy).(*Order)
		bids = append(bids, [2]float64{copiedBid.PriceLevel, copiedBid.Quantity})
	}

	heap.Init(&askCopy)
	for askCopy.Len() > 0 {
		copiedAsk := heap.Pop(&askCopy).(*Order)
		asks = append(asks, [2]float64{copiedAsk.PriceLevel, copiedAsk.Quantity})
	}

	return bids, asks
}
