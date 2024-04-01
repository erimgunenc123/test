package orderbook

import (
	"genericAPI/binanceconnector/dto"
	"sync"
)

type bidPriceLevelNode struct {
	order dto.Order
	next  *bidPriceLevelNode
}

type BidPriceLevelList struct {
	head *bidPriceLevelNode
	size uint32
	lock sync.Mutex
}

func (l *BidPriceLevelList) GetAllBids() []dto.Order {
	l.lock.Lock()
	defer l.lock.Unlock()
	result := make([]dto.Order, l.size)
	curNode := l.head
	for curNode.next != nil {
		result = append(result, curNode.order) // these aren't pointers so no need to copy values
		curNode = curNode.next
	}
	return result
}

func (l *BidPriceLevelList) BestBid() float64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.head.order.PriceLevel
}

func (l *BidPriceLevelList) Insert(priceLevel float64, qty uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.head == nil {
		l.head = &bidPriceLevelNode{
			order: dto.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  nil,
		}
		l.size = 1
		return
	}

	prevNode := l.head
	tempNode := l.head
	for (tempNode.next != nil) && (tempNode.order.PriceLevel > priceLevel) {
		prevNode = tempNode
		tempNode = tempNode.next
	}

	if tempNode.order.PriceLevel == priceLevel {
		tempNode.order.Quantity += qty
	} else if tempNode.order.PriceLevel < priceLevel {
		prevNode.next = &bidPriceLevelNode{
			order: dto.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  tempNode,
		}
	} else {
		tempNode.next = &bidPriceLevelNode{
			order: dto.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  nil,
		}
	}
	l.size += 1
}

func (l *BidPriceLevelList) Update(priceLevel float64, qty uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	prevNode := l.head
	tempNode := l.head
	for tempNode.order.PriceLevel != priceLevel {
		prevNode = tempNode
		tempNode = tempNode.next
	}
	if qty == 0 {
		prevNode.next = tempNode.next
		l.size -= 1
	} else {
		tempNode.order.Quantity = qty
	}
}
