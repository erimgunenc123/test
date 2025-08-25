package orderbook

import (
	"genericAPI/exchange/common"
	"sync"
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
	//Bids   *BidsHeap
	//Asks   *AsksHeap

	bidMap sync.Map
	askMap sync.Map
}

func NewHeapOrderbook(symbol string) *HeapOrderbook {
	return &HeapOrderbook{
		At:     time.Now(),
		Symbol: symbol,
		//Bids:   &BidsHeap{},
		//Asks:   &AsksHeap{},
		bidMap: sync.Map{},
		askMap: sync.Map{},
	}
}
func (ob *HeapOrderbook) Insert(isBid bool, price, quantity float64) {
	if isBid {
		if qty, ok := ob.bidMap.Load(price); ok {
			if quantity == 0 {
				ob.bidMap.Delete(price)
				return
			}
			ob.bidMap.Store(price, qty.(float64)+quantity)
		} else {
			ob.bidMap.Store(price, quantity)
		}
	} else {
		if qty, ok := ob.askMap.Load(price); ok {
			if quantity == 0 {
				ob.askMap.Delete(price)
				return
			}
			ob.askMap.Store(price, qty.(float64)+quantity)
		} else {
			ob.askMap.Store(price, quantity)
		}
	}
}
func (o *HeapOrderbook) GetAllAsList() ([][2]float64, [][2]float64) {
	bids := [][2]float64{}
	asks := [][2]float64{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		o.bidMap.Range(func(k, v interface{}) bool {
			bids = append(bids, [2]float64{k.(float64), v.(float64)})
			return true
		})
	}()
	go func() {
		defer wg.Done()
		o.askMap.Range(func(k, v interface{}) bool {
			asks = append(asks, [2]float64{k.(float64), v.(float64)})
			return true
		})
	}()
	wg.Wait()
	return bids, asks
}

//
//func (o *HeapOrderbook) GetSnapshot() (bids [][2]float64, asks [][2]float64) {
//	bidCopy := make(BidsHeap, len(*o.Bids))
//	copy(bidCopy, *o.Bids)
//	heap.Init(&bidCopy)
//	for bidCopy.Len() > 0 {
//		bid := heap.Pop(&bidCopy).(*Order)
//		bids = append(bids, [2]float64{bid.PriceLevel, bid.Quantity})
//	}
//
//	askCopy := make(AsksHeap, len(*o.Asks))
//	copy(askCopy, *o.Asks)
//	heap.Init(&askCopy)
//	for askCopy.Len() > 0 {
//		ask := heap.Pop(&askCopy).(*Order)
//		asks = append(asks, [2]float64{ask.PriceLevel, ask.Quantity})
//	}
//
//	return bids, asks
//}
//
//func (o *HeapOrderbook) Insert(isBid bool, price, amount float64) {
//	order := &Order{
//		Order: common.Order{
//			PriceLevel: price,
//			Quantity:   amount,
//		},
//		Index: 0,
//	}
//	if isBid {
//		heap.Push(o.Bids, order)
//	} else {
//		heap.Push(o.Asks, order)
//	}
//}
//
//func (o *HeapOrderbook) Update(isBid bool, price, amount float64) {
//	var hHeap heap.Interface
//	var hMap map[float64]*Order
//
//	if isBid {
//		hHeap = o.Bids
//		hMap = o.bidMap
//	} else {
//		hHeap = o.Asks
//		hMap = o.askMap
//	}
//
//	if existing, ok := hMap[price]; ok {
//		if amount == 0 {
//			heap.Remove(hHeap, existing.Index)
//			delete(hMap, price)
//		} else {
//			existing.Quantity = amount
//		}
//	} else if amount > 0 {
//		order := &Order{
//			Order: common.Order{
//				PriceLevel: price,
//				Quantity:   amount,
//			},
//		}
//		heap.Push(hHeap, order)
//		hMap[price] = order
//	}
//}
//
//func (o *HeapOrderbook) GetAllAsList() (bids [][2]float64, asks [][2]float64) {
//	bidCopy := make(BidsHeap, len(*o.Bids))
//	copy(bidCopy, *o.Bids)
//	askCopy := make(AsksHeap, len(*o.Asks))
//	copy(askCopy, *o.Asks)
//
//	heap.Init(&bidCopy)
//	for bidCopy.Len() > 0 {
//		copiedBid := heap.Pop(&bidCopy).(*Order)
//		bids = append(bids, [2]float64{copiedBid.PriceLevel, copiedBid.Quantity})
//	}
//
//	heap.Init(&askCopy)
//	for askCopy.Len() > 0 {
//		copiedAsk := heap.Pop(&askCopy).(*Order)
//		asks = append(asks, [2]float64{copiedAsk.PriceLevel, copiedAsk.Quantity})
//	}
//
//	return bids, asks
//}
