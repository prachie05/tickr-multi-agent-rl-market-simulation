package main

import "testing"

// addOrder adds a valid order and fails the test if an error occurs.
func addOrder(t *testing.T, book *OrderBook, order Order) {
	t.Helper()

	if err := book.AddOrder(order); err != nil {
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

			err := book.AddOrder(tt.order)

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
