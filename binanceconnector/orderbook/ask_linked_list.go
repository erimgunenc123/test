package orderbook

import (
	"genericAPI/binanceconnector/dto"
	"sync"
)

type askPriceLevelNode struct {
	order dto.Order
	next  *askPriceLevelNode
}

type AskPriceLevelList struct {
	head *askPriceLevelNode
	size uint32
	lock sync.Mutex
}

func (l *AskPriceLevelList) GetAllAsks() []dto.Order {
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

func (l *AskPriceLevelList) BestAsk() float64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.head.order.PriceLevel
}

func (l *AskPriceLevelList) Insert(priceLevel float64, qty uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.head == nil {
		l.head = &askPriceLevelNode{
			order: dto.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  nil,
		}
		l.size += 1
		return
	}

	prevNode := l.head
	tempNode := l.head
	for (tempNode.next != nil) && (tempNode.order.PriceLevel < priceLevel) {
		prevNode = tempNode
		tempNode = tempNode.next
	}

	if tempNode.order.PriceLevel == priceLevel {
		tempNode.order.Quantity += qty
	} else if tempNode.order.PriceLevel > priceLevel {
		prevNode.next = &askPriceLevelNode{
			order: dto.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  tempNode,
		}
	} else {
		tempNode.next = &askPriceLevelNode{
			order: dto.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  nil,
		}
	}
	l.size += 1
}

func (l *AskPriceLevelList) Update(priceLevel float64, qty uint64) {
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
