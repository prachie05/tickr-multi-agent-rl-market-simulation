package main

import "testing"

func TestCreateOrder(t *testing.T) {
	order := Order{
		ID:        101,
		Price:     10250,
		Quantity:  10,
		Side:      Buy,
		Timestamp: 1728700000,
	}

	if order.ID != 101 {
		t.Errorf("Expected ID 101 go %d", order.ID)
	}

	if order.Price != 10250 {
		t.Errorf("Expected Price of 10250, got %d", order.Price)
	}
	if order.Quantity != 10 {
		t.Errorf("expected quantity 10, got %d", order.Quantity)
	}

	if order.Side != Buy {
		t.Errorf("expected side Buy, got %v", order.Side)
	}

	if order.Timestamp != 1728700000 {
		t.Errorf("expected timestamp 1728700000, got %d", order.Timestamp)
	}
}
