package main

import "testing"

// addOrder adds a valid order and fails the test if an error occurs.
func addOrder(t *testing.T, book *OrderBook, order Order) {
	t.Helper()

	if _, err := book.AddOrder(order); err != nil {
		t.Fatalf("unexpected error adding order: %v", err)
	}
}

func TestBidOrderPriority(t *testing.T) {
	book := NewOrderBook()

	addOrder(t, book, Order{ID: 1, Price: 10000, Quantity: 10, Side: Buy})
	addOrder(t, book, Order{ID: 2, Price: 10500, Quantity: 10, Side: Buy})
	addOrder(t, book, Order{ID: 3, Price: 10200, Quantity: 10, Side: Buy})

	bids := book.GetBids()
	expectedIDs := []int64{2, 3, 1}

	if len(bids) != len(expectedIDs) {
		t.Fatalf("expected %d bids, got %d", len(expectedIDs), len(bids))
	}

	for i, expectedID := range expectedIDs {
		if bids[i].ID != expectedID {
			t.Errorf("position %d: expected order ID %d, got %d",
				i, expectedID, bids[i].ID)
		}
	}
}

func TestAskOrderPriority(t *testing.T) {
	book := NewOrderBook()

	addOrder(t, book, Order{ID: 1, Price: 10000, Quantity: 10, Side: Sell})
	addOrder(t, book, Order{ID: 2, Price: 10500, Quantity: 10, Side: Sell})
	addOrder(t, book, Order{ID: 3, Price: 10200, Quantity: 10, Side: Sell})

	asks := book.GetAsks()
	expectedIDs := []int64{1, 3, 2}

	if len(asks) != len(expectedIDs) {
		t.Fatalf("expected %d asks, got %d", len(expectedIDs), len(asks))
	}

	for i, expectedID := range expectedIDs {
		if asks[i].ID != expectedID {
			t.Errorf("position %d: expected order ID %d, got %d",
				i, expectedID, asks[i].ID)
		}
	}
}

func TestFIFOAtEqualBidPrice(t *testing.T) {
	book := NewOrderBook()

	addOrder(t, book, Order{ID: 1, Price: 10500, Quantity: 10, Side: Buy})
	addOrder(t, book, Order{ID: 2, Price: 10500, Quantity: 20, Side: Buy})

	bids := book.GetBids()

	if len(bids) != 2 {
		t.Fatalf("expected 2 bids, got %d", len(bids))
	}

	if bids[0].ID != 1 {
		t.Errorf("expected first order ID 1, got %d", bids[0].ID)
	}

	if bids[1].ID != 2 {
		t.Errorf("expected second order ID 2, got %d", bids[1].ID)
	}
}

func TestFIFOAtEqualAskPrice(t *testing.T) {
	book := NewOrderBook()

	addOrder(t, book, Order{ID: 1, Price: 10500, Quantity: 10, Side: Sell})
	addOrder(t, book, Order{ID: 2, Price: 10500, Quantity: 20, Side: Sell})

	asks := book.GetAsks()

	if len(asks) != 2 {
		t.Fatalf("expected 2 asks, got %d", len(asks))
	}

	if asks[0].ID != 1 {
		t.Errorf("expected first order ID 1, got %d", asks[0].ID)
	}

	if asks[1].ID != 2 {
		t.Errorf("expected second order ID 2, got %d", asks[1].ID)
	}
}

func TestEmptyOrderBook(t *testing.T) {
	book := NewOrderBook()

	if len(book.GetBids()) != 0 {
		t.Errorf("expected empty bids, got %d orders", len(book.GetBids()))
	}

	if len(book.GetAsks()) != 0 {
		t.Errorf("expected empty asks, got %d orders", len(book.GetAsks()))
	}
}

func TestBidAndAskSeparation(t *testing.T) {
	book := NewOrderBook()

	addOrder(t, book, Order{ID: 1, Price: 10000, Quantity: 10, Side: Buy})
	addOrder(t, book, Order{ID: 2, Price: 11000, Quantity: 10, Side: Sell})

	bids := book.GetBids()
	asks := book.GetAsks()

	if len(bids) != 1 {
		t.Fatalf("expected 1 bid, got %d", len(bids))
	}

	if len(asks) != 1 {
		t.Fatalf("expected 1 ask, got %d", len(asks))
	}

	if bids[0].ID != 1 {
		t.Errorf("expected bid ID 1, got %d", bids[0].ID)
	}

	if asks[0].ID != 2 {
		t.Errorf("expected ask ID 2, got %d", asks[0].ID)
	}
}

func TestAddOrderRejectsInvalidOrders(t *testing.T) {
	tests := []struct {
		name  string
		order Order
	}{
		{
			name:  "invalid side",
			order: Order{ID: 1, Price: 10000, Quantity: 10, Side: Side(99)},
		},
		{
			name:  "zero price",
			order: Order{ID: 1, Price: 0, Quantity: 10, Side: Buy},
		},
		{
			name:  "negative price",
			order: Order{ID: 1, Price: -100, Quantity: 10, Side: Buy},
		},
		{
			name:  "zero quantity",
			order: Order{ID: 1, Price: 10000, Quantity: 0, Side: Buy},
		},
		{
			name:  "negative quantity",
			order: Order{ID: 1, Price: 10000, Quantity: -5, Side: Buy},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := NewOrderBook()

			_, err := book.AddOrder(tt.order)

			if err == nil {
				t.Fatal("expected an error for an invalid order, got nil")
			}

			if len(book.GetBids()) != 0 {
				t.Errorf("expected no bids after rejection, got %d",
					len(book.GetBids()))
			}

			if len(book.GetAsks()) != 0 {
				t.Errorf("expected no asks after rejection, got %d",
					len(book.GetAsks()))
			}
		})
	}
}

func TestGettersReturnCopies(t *testing.T) {
	book := NewOrderBook()

	addOrder(t, book, Order{
		ID:       1,
		Price:    10000,
		Quantity: 10,
		Side:     Buy,
	})

	addOrder(t, book, Order{
		ID:       2,
		Price:    11000,
		Quantity: 10,
		Side:     Sell,
	})

	bids := book.GetBids()
	asks := book.GetAsks()

	// Modify the returned slices.
	bids[0].Price = 1
	asks[0].Price = 1

	// Verify the internal order book remains unchanged.
	if book.GetBids()[0].Price != 10000 {
		t.Error("modifying returned bids changed the internal order book")
	}

	if book.GetAsks()[0].Price != 11000 {
		t.Error("modifying returned asks changed the internal order book")
	}
}

func TestBuyMatchesRestingSell(t *testing.T) {
	book := NewOrderBook()

	sellOrder := Order{
		ID:        1,
		Price:     105,
		Quantity:  5,
		Side:      Sell,
		Timestamp: 1,
	}

	addOrder(t, book, sellOrder)

	buyOrder := Order{
		ID:        2,
		Price:     105,
		Quantity:  3,
		Side:      Buy,
		Timestamp: 2,
	}

	trades, remainingBuy := book.MatchOrder(buyOrder)

	// Exactly one trade should occur.
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}

	// Verify trade details.
	if trades[0].Quantity != 3 {
		t.Errorf("expected trade quantity 3, got %d", trades[0].Quantity)
	}

	if trades[0].Price != 105 {
		t.Errorf("expected trade price 105, got %d", trades[0].Price)
	}

	if trades[0].BuyerID != 2 {
		t.Errorf("expected buyer ID 2, got %d", trades[0].BuyerID)
	}

	if trades[0].SellerID != 1 {
		t.Errorf("expected seller ID 1, got %d", trades[0].SellerID)
	}

	if trades[0].TradeID != 1 {
		t.Errorf("expected trade ID 1, got %d", trades[0].TradeID)
	}

	// Verify the resting sell order has 2 units remaining.
	asks := book.GetAsks()
	if len(asks) != 1 {
		t.Fatalf("expected 1 remaining ask, got %d", len(asks))
	}

	if asks[0].Quantity != 2 {
		t.Errorf("expected remaining ask quantity 2, got %d", asks[0].Quantity)
	}

	// The incoming buy order should be fully filled.
	if remainingBuy.Quantity != 0 {
		t.Errorf("expected remaining buy quantity 0, got %d", remainingBuy.Quantity)
	}
}

func TestSellMatchesRestingBuy(t *testing.T) {
	book := NewOrderBook()

	buyOrder := Order{
		ID:        1,
		Price:     110,
		Quantity:  5,
		Side:      Buy,
		Timestamp: 1,
	}

	addOrder(t, book, buyOrder)

	sellOrder := Order{
		ID:        2,
		Price:     105,
		Quantity:  3,
		Side:      Sell,
		Timestamp: 2,
	}

	trades, remainingSell := book.MatchOrder(sellOrder)

	// Exactly one trade should occur.
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}

	// Verify trade details.
	if trades[0].Quantity != 3 {
		t.Errorf("expected trade quantity 3, got %d", trades[0].Quantity)
	}

	if trades[0].Price != 110 {
		t.Errorf("expected trade price 110, got %d", trades[0].Price)
	}

	if trades[0].BuyerID != 1 {
		t.Errorf("expected buyer ID 1, got %d", trades[0].BuyerID)
	}

	if trades[0].SellerID != 2 {
		t.Errorf("expected seller ID 2, got %d", trades[0].SellerID)
	}

	if trades[0].TradeID != 1 {
		t.Errorf("expected trade ID 1, got %d", trades[0].TradeID)
	}

	// Verify the resting buy order has 2 units remaining.
	bids := book.GetBids()
	if len(bids) != 1 {
		t.Fatalf("expected 1 remaining bid, got %d", len(bids))
	}

	if bids[0].Quantity != 2 {
		t.Errorf("expected remaining bid quantity 2, got %d", bids[0].Quantity)
	}

	// The incoming sell order should be fully filled.
	if remainingSell.Quantity != 0 {
		t.Errorf("expected remaining sell quantity 0, got %d", remainingSell.Quantity)
	}
}

func TestMatchesMultipleRestingOrders(t *testing.T) {
	book := NewOrderBook()

	sellerA := Order{ID: 1, Price: 100, Quantity: 3, Side: Sell, Timestamp: 1}
	sellerB := Order{ID: 2, Price: 105, Quantity: 4, Side: Sell, Timestamp: 2}
	buyer := Order{ID: 3, Price: 110, Quantity: 5, Side: Buy, Timestamp: 3}

	addOrder(t, book, sellerA)
	addOrder(t, book, sellerB)

	trades, remainingBuy := book.MatchOrder(buyer)
	asks := book.GetAsks()

	if len(trades) != 2 {
		t.Fatalf("Expected 2 trades, got %d", len(trades))
	}

	if remainingBuy.Quantity != 0 {
		t.Errorf("Expected remaining buyer quantity to be 0, got %d", remainingBuy.Quantity)
	}

	if trades[0].Quantity != 3 || trades[0].Price != 100 {
		t.Errorf("Expected first trade: quantity 3 at price 100, got quantity %d at price %d",
			trades[0].Quantity, trades[0].Price)
	}

	if trades[1].Quantity != 2 || trades[1].Price != 105 {
		t.Errorf("Expected second trade: quantity 2 at price 105, got quantity %d at price %d",
			trades[1].Quantity, trades[1].Price)
	}

	if trades[0].SellerID != 1 {
		t.Errorf("Expected first trade's seller ID to be 1, got %d", trades[0].SellerID)
	}

	if trades[1].SellerID != 2 {
		t.Errorf("Expected second trade's seller ID to be 2, got %d", trades[1].SellerID)
	}

	if len(asks) != 1 {
		t.Fatalf("Expected 1 remaining ask, got %d", len(asks))
	}

	if asks[0].ID != 2 || asks[0].Quantity != 2 {
		t.Errorf("Expected seller B (ID 2) to have quantity 2 remaining, got ID %d with quantity %d",
			asks[0].ID, asks[0].Quantity)
	}
}

func TestNoPriceMatch(t *testing.T) {
	book := NewOrderBook()

	seller := Order{ID: 1, Price: 110, Quantity: 5, Side: Sell, Timestamp: 1}
	buyer := Order{ID: 2, Price: 100, Quantity: 3, Side: Buy, Timestamp: 2}

	addOrder(t, book, seller)

	trades, remainingBuy := book.MatchOrder(buyer)

	if len(trades) != 0 {
		t.Errorf("Expeceted trades 0, got %d", len(trades))
	}

	if remainingBuy.Quantity != 3 {
		t.Errorf("Expected Remaining Buy 3, got %d", remainingBuy.Quantity)
	}

	asks := book.GetAsks()
	bids := book.GetBids()

	if len(asks) != 1 {
		t.Fatalf("Expected 1 ask, got %d", len(asks))
	}

	if asks[0].Quantity != 5 {
		t.Errorf("Expected ask quantity 5, got %d", asks[0].Quantity)
	}

	if len(bids) != 0 {
		t.Errorf("Expected 0, got %d", len(bids))
	}

}

func TestFIFOWhenMatching(t *testing.T) {
	book := NewOrderBook()

	buyerA := Order{ID: 1, Price: 100, Quantity: 5, Side: Buy, Timestamp: 1}
	buyerB := Order{ID: 2, Price: 100, Quantity: 5, Side: Buy, Timestamp: 2}
	seller := Order{ID: 3, Price: 100, Quantity: 6, Side: Sell, Timestamp: 3}

	addOrder(t, book, buyerA)
	addOrder(t, book, buyerB)

	trades, remainingSell := book.MatchOrder(seller)

	if len(trades) != 2 {
		t.Fatalf("Expected 2 trades, got %d", len(trades))
	}

	if remainingSell.Quantity != 0 {
		t.Errorf("Expected remaining sell quantity 0, got %d", remainingSell.Quantity)
	}

	if trades[0].BuyerID != 1 || trades[0].Quantity != 5 {
		t.Errorf("Expected first trade: buyer ID 1, quantity 5; got buyer ID %d, quantity %d",
			trades[0].BuyerID, trades[0].Quantity)
	}

	if trades[1].BuyerID != 2 || trades[1].Quantity != 1 {
		t.Errorf("Expected second trade: buyer ID 2, quantity 1; got buyer ID %d, quantity %d",
			trades[1].BuyerID, trades[1].Quantity)
	}

	bids := book.GetBids()

	if len(bids) != 1 {
		t.Fatalf("Expected 1 remaining bid, got %d", len(bids))
	}

	if bids[0].ID != 2 || bids[0].Quantity != 4 {
		t.Errorf("Expected buyer B (ID 2) to remain with quantity 4, got ID %d, quantity %d",
			bids[0].ID, bids[0].Quantity)
	}
}

func TestAddOrderFullFill(t *testing.T) {
	book := NewOrderBook()

	seller := Order{ID: 1, Price: 100, Quantity: 5, Side: Sell, Timestamp: 1}
	buyer := Order{ID: 2, Price: 100, Quantity: 5, Side: Buy, Timestamp: 2}

	addOrder(t, book, seller)

	trades, err := book.AddOrder(buyer)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades))
	}

	if trades[0].Quantity != 5 {
		t.Errorf("Expected trade quantity 5, got %d", trades[0].Quantity)
	}

	asks := book.GetAsks()
	bids := book.GetBids()

	if len(asks) != 0 {
		t.Errorf("Expected 0 asks, got %d", len(asks))
	}

	if len(bids) != 0 {
		t.Errorf("Expected 0 bids, got %d", len(bids))
	}

	if trades[0].BuyerID != 2 || trades[0].SellerID != 1 {
		t.Errorf("Expected buyer ID 2 and seller ID 1, got buyer %d and seller %d",
			trades[0].BuyerID, trades[0].SellerID)
	}
}

func TestAddOrderPartialFill(t *testing.T) {
	book := NewOrderBook()

	seller := Order{ID: 1, Price: 100, Quantity: 3, Side: Sell, Timestamp: 1}
	buyer := Order{ID: 2, Price: 100, Quantity: 5, Side: Buy, Timestamp: 2}

	addOrder(t, book, seller)

	trades, err := book.AddOrder(buyer)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades))
	}

	if trades[0].Quantity != 3 {
		t.Errorf("Expected trade quantity 3, got %d", trades[0].Quantity)
	}

	asks := book.GetAsks()
	bids := book.GetBids()

	if len(asks) != 0 {
		t.Errorf("Expected 0 asks, got %d", len(asks))
	}

	if len(bids) != 1 {
		t.Fatalf("Expected 1 bid, got %d", len(bids))
	}

	if bids[0].Quantity != 2 || bids[0].ID != 2 {
		t.Errorf("Expected buyer ID 2 with quantity 2, got ID %d with quantity %d",
			bids[0].ID, bids[0].Quantity)
	}
}

func TestAddOrderNoMatch(t *testing.T) {
	book := NewOrderBook()

	seller := Order{ID: 1, Price: 110, Quantity: 5, Side: Sell, Timestamp: 1}
	buyer := Order{ID: 2, Price: 100, Quantity: 3, Side: Buy, Timestamp: 2}

	addOrder(t, book, seller)

	trades, err := book.AddOrder(buyer)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(trades) != 0 {
		t.Fatalf("Expected 0 trades, got %d", len(trades))
	}

	asks := book.GetAsks()
	bids := book.GetBids()

	if len(asks) != 1 {
		t.Fatalf("Expected 1 ask, got %d", len(asks))
	}

	if len(bids) != 1 {
		t.Fatalf("Expected 1 bid, got %d", len(bids))
	}

	if asks[0].ID != 1 || asks[0].Quantity != 5 {
		t.Errorf("Expected seller ID 1 with quantity 5, got ID %d with quantity %d",
			asks[0].ID, asks[0].Quantity)
	}

	if bids[0].ID != 2 || bids[0].Quantity != 3 {
		t.Errorf("Expected buyer ID 2 with quantity 3, got ID %d with quantity %d",
			bids[0].ID, bids[0].Quantity)
	}
}

func TestCancelBid(t *testing.T) {
	book := NewOrderBook()

	bid := Order{ID: 1, Price: 100, Quantity: 5, Side: Buy, Timestamp: 1}
	addOrder(t, book, bid)

	err := book.CancelOrder(1)
	if err != nil {
		t.Fatalf("Unexpected error cancelling bid: %v", err)
	}

	bids := book.GetBids()
	if len(bids) != 0 {
		t.Errorf("Expected 0 bids, got %d", len(bids))
	}
}

func TestCancelAsk(t *testing.T) {
	book := NewOrderBook()

	ask := Order{ID: 1, Price: 110, Quantity: 5, Side: Sell, Timestamp: 1}
	addOrder(t, book, ask)

	err := book.CancelOrder(1)
	if err != nil {
		t.Fatalf("Unexpected error cancelling ask: %v", err)
	}

	asks := book.GetAsks()
	if len(asks) != 0 {
		t.Errorf("Expected 0 asks, got %d", len(asks))
	}
}

func TestCancelNonexistentOrder(t *testing.T) {
	book := NewOrderBook()

	err := book.CancelOrder(999)
	if err == nil {
		t.Fatal("Expected error when cancelling nonexistent order")
	}
}

func TestCancelPreservesOtherOrders(t *testing.T) {
	book := NewOrderBook()

	bidA := Order{ID: 1, Price: 110, Quantity: 5, Side: Buy, Timestamp: 1}
	bidB := Order{ID: 2, Price: 100, Quantity: 3, Side: Buy, Timestamp: 2}
	bidC := Order{ID: 3, Price: 90, Quantity: 7, Side: Buy, Timestamp: 3}

	addOrder(t, book, bidA)
	addOrder(t, book, bidB)
	addOrder(t, book, bidC)

	err := book.CancelOrder(2)
	if err != nil {
		t.Fatalf("Unexpected error cancelling bid: %v", err)
	}

	bids := book.GetBids()

	if len(bids) != 2 {
		t.Fatalf("Expected 2 bids, got %d", len(bids))
	}

	if bids[0].ID != 1 || bids[0].Quantity != 5 {
		t.Errorf("Expected first bid ID 1 with quantity 5, got ID %d with quantity %d",
			bids[0].ID, bids[0].Quantity)
	}

	if bids[1].ID != 3 || bids[1].Quantity != 7 {
		t.Errorf("Expected second bid ID 3 with quantity 7, got ID %d with quantity %d",
			bids[1].ID, bids[1].Quantity)
	}
}
