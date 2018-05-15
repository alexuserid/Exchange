package main

import (
	"time"
	"container/heap"
	"sync"
	"strconv"
)

type order struct {
	pair string
	amount float64
	price float64
	date time.Time
}

// make different maps and queues for each pair
var (
	mapOidOrder = make(map[b32]order)
	oq = make(PriorityQueue, 1)
)

func makeQueueItem(pa string, am, pr float64) {
	ord := &Item{
		value: order{pair: pa, amount: am, price: pr},
		priority: pr,
	}
	heap.Push(&oq, ord)
	oq.update(ord, ord.value, float64(time.Now().UnixNano()))
}

func limitOrder(userInfo user, pair, amount, price string) errorc {
	amountfl, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return errParseFloat
	}
	pricefl, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return errParseFloat
	}
	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	oid, errc := getUniqueId(markerOid)
	if errc != errNo {
		return errc
	}
	mapOidOrder[oid] = order{pair, amountfl, pricefl, time.Now()}
	makeQueueItem(pair, amountfl, pricefl)
	return errNo
}

func marketOrder(userInfo user, pair, amount string) errorc {
	// convert amount to float64
	// get unique order id
	// execute
	return errNo
}

func cancelOrder(userInfo user, pair, oid string) errorc {
	// convert oid to b32
	// remove order from that pair queue
	return errNo
}
