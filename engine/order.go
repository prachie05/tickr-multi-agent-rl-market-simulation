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

type Trade struct {
	TradeID  int64
	BuyerID  int64
	SellerID int64
	Price    int64
	Quantity int64
}
