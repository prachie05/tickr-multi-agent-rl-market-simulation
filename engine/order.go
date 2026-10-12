package main

type Side int

const (
	Buy Side = iota
	Sell
)

type Order struct {
	ID        int64
	Price     int64
	Quantity  int64
	Side      Side
	Timestamp int64
}
