package orderbook

import (
	"genericAPI/exchange/common"
	"sync"
)

type priceLevelNode struct {
	order common.Order
	next  *priceLevelNode
}

type PriceLevelList struct {
	head *priceLevelNode
	size uint32
	lock sync.Mutex
	side common.Side
}

func (l *PriceLevelList) GetAll() []common.Order {
	l.lock.Lock()
	defer l.lock.Unlock()
	result := make([]common.Order, l.size)
	curNode := l.head
	for curNode.next != nil {
		result = append(result, curNode.order) // these aren't pointers so no need to copy values
		curNode = curNode.next
	}
	return result
}

func (l *PriceLevelList) Best() float64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.head.order.PriceLevel
}

func (l *PriceLevelList) Insert(priceLevel float64, qty uint64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.head == nil {
		l.head = &priceLevelNode{
			order: common.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  nil,
		}
		l.size += 1
		return
	}

	prevNode := l.head
	tempNode := l.head
	for tempNode.next != nil {
		switch l.side {
		case common.Ask:
			if tempNode.order.PriceLevel < priceLevel {
				prevNode = tempNode
				tempNode = tempNode.next
			}
		case common.Bid:
			if tempNode.order.PriceLevel > priceLevel {
				prevNode = tempNode
				tempNode = tempNode.next
			}
		}
	}

	if tempNode.order.PriceLevel == priceLevel {
		tempNode.order.Quantity += qty
	} else if (l.side == common.Ask && tempNode.order.PriceLevel > priceLevel) || (l.side == common.Bid && tempNode.order.PriceLevel < priceLevel) {
		prevNode.next = &priceLevelNode{
			order: common.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  tempNode,
		}
	} else {
		tempNode.next = &priceLevelNode{
			order: common.Order{PriceLevel: priceLevel, Quantity: qty},
			next:  nil,
		}
	}
	l.size += 1
}

func (l *PriceLevelList) Update(priceLevel float64, qty uint64) {
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
