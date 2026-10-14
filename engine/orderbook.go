package main

import "fmt"

type OrderBook struct {
	Bids        []Order
	Asks        []Order
	NextTradeID int64
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids:        make([]Order, 0),
		Asks:        make([]Order, 0),
		NextTradeID: 1,
	}
}

func (ob *OrderBook) AddOrder(order Order) ([]Trade, error) {
	if order.Side != Buy && order.Side != Sell {
		return nil, fmt.Errorf("invalid order side: %d", order.Side)
	}

	if order.Price <= 0 {
		return nil, fmt.Errorf("Price must be positive!")
	}

	if order.Quantity <= 0 {
		return nil, fmt.Errorf("Quantity should be greater than zero")
	}

	trades, remaining := ob.MatchOrder(order)

	if remaining.Quantity > 0 {
		if remaining.Side == Buy {
			ob.Bids = insertByPrice(ob.Bids, remaining, true)
		} else {
			ob.Asks = insertByPrice(ob.Asks, remaining, false)
		}
	}
	return trades, nil
}

func (ob *OrderBook) MatchOrder(order Order) ([]Trade, Order) {
	var trades []Trade

	switch order.Side {
	case Buy:
		for len(ob.Asks) > 0 && order.Quantity > 0 && ob.Asks[0].Price <= order.Price {
			bestAsk := &ob.Asks[0]

			tradeQty := order.Quantity
			if bestAsk.Quantity < tradeQty {
				tradeQty = bestAsk.Quantity
			}

			trade := Trade{
				TradeID:  ob.NextTradeID,
				BuyerID:  order.ID,
				SellerID: bestAsk.ID,
				Price:    bestAsk.Price,
				Quantity: tradeQty,
			}
			ob.NextTradeID++

			trades = append(trades, trade)

			bestAsk.Quantity -= tradeQty
			order.Quantity -= tradeQty

			if bestAsk.Quantity == 0 {
				ob.Asks = ob.Asks[1:]
			}
		}

	case Sell:
		for len(ob.Bids) > 0 && order.Quantity > 0 && order.Price <= ob.Bids[0].Price {
			bestBuy := &ob.Bids[0]

			tradeQty := order.Quantity
			if bestBuy.Quantity < tradeQty {
				tradeQty = bestBuy.Quantity
			}

			trade := Trade{
				TradeID:  ob.NextTradeID,
				BuyerID:  bestBuy.ID,
				SellerID: order.ID,
				Price:    bestBuy.Price,
				Quantity: tradeQty,
			}
			ob.NextTradeID++

			trades = append(trades, trade)

			bestBuy.Quantity -= tradeQty
			order.Quantity -= tradeQty

			if bestBuy.Quantity == 0 {
				ob.Bids = ob.Bids[1:]
			}
		}
	}

	return trades, order
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

func (ob *OrderBook) CancelOrder(id int64) error {

	for i, bid := range ob.Bids {
		if id == bid.ID {
			copy(ob.Bids[i:], ob.Bids[i+1:])
			ob.Bids = ob.Bids[:len(ob.Bids)-1]
			return nil
		}
	}

	for i, ask := range ob.Asks {
		if id == ask.ID {
			copy(ob.Asks[i:], ob.Asks[i+1:])
			ob.Asks = ob.Asks[:len(ob.Asks)-1]
			return nil
		}

	}
	return fmt.Errorf("order ID %d not found", id)

}
