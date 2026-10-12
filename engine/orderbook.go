package main

import "fmt"

type OrderBook struct {
	Bids []Order
	Asks []Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make([]Order, 0),
		Asks: make([]Order, 0),
	}
}

func (ob *OrderBook) AddOrder(order Order) error {
	if order.Side != Buy && order.Side != Sell {
		return fmt.Errorf("invalid order side: %d", order.Side)
	}

	if order.Price <= 0 {
		return fmt.Errorf("Price must be positive!")
	}

	if order.Quantity <= 0 {
		return fmt.Errorf("Quantity should be greater than zero")
	}

	if order.Side == Buy {
		ob.Bids = insertByPrice(ob.Bids, order, true)
	} else {
		ob.Asks = insertByPrice(ob.Asks, order, false)
	}

	return nil
}

func insertByPrice(orders []Order, order Order, descending bool) []Order {
	index := 0

	for index < len(orders) {
		if descending && order.Price > orders[index].Price {
			break
		}

		if !descending && order.Price < orders[index].Price {
			break
		}

		index++
	}
	orders = append(orders, Order{})
	copy(orders[index+1:], orders[index:])
	orders[index] = order

	return orders

}

func (ob *OrderBook) GetBids() []Order {
	return append([]Order(nil), ob.Bids...)
}

func (ob *OrderBook) GetAsks() []Order {
	return append([]Order(nil), ob.Asks...)
}
