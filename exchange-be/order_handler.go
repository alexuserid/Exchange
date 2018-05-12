package main

import (
	"time"
)

type deal struct {
	pair string
	amoutn float64
	price float64
	date time.Time
}

var mapDeals = make(map[b32]deal)

func limitOrder(userInfo user, pair, amount, price string) error {
	// convert amount to float64
	// get unique order id
	// push order to queue
	return nil
}

func nowOrder(userInfo user, pair, amount string) error {
	// convert amount to float64
	// get unique order if
	// execute
	return nil
}

func cancelOrder(userInfo user, pair, oid string) error {
	// convert oid to b32
	// pop order from that pair queue
	return nil
}
