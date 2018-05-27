package main

import (
	"container/heap"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/starius/status"
)

type OrderID [idl]byte

type order struct {
	pair   string
	amount float64
	price  float64
	date   time.Time
}

// make different maps and queues for each pair
var (
	mapOidOrder = make(map[OrderID]order)
	oq          = make(PriorityQueue, 1)
	OidMutex    sync.Mutex
)

func getOid() (OrderID, error) {
	for {
		randoms, err := getRandoms32()
		if err != nil {
			return OrderID{}, status.Format("getRandom32: %v", err)
		}
		hb := toHex(randoms)
		var id OrderID
		copy(id[:], hb[:])

		if _, has := mapOidOrder[id]; !has {
			return id, nil
		}
	}
}

func makeQueueItem(pa string, am, pr float64) {
	ord := &Item{
		value:    order{pair: pa, amount: am, price: pr},
		priority: pr,
	}
	heap.Push(&oq, ord)
	oq.update(ord, ord.value, float64(time.Now().UnixNano()))
}

func limitOrder(userInfo user, pair, amount, price string) error {
	amountfl, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return status.WithCode(http.StatusBadRequest, "Wrong amount format:ParseFloat: %v", err)
	}
	pricefl, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return status.WithCode(http.StatusBadRequest, "Wrong price fomat: ParseFloat: %v", err)
	}
	OidMutex.Lock()
	defer OidMutex.Unlock()

	oid, err := getOid()
	if err != nil {
		return err
	}
	mapOidOrder[oid] = order{pair, amountfl, pricefl, time.Now()}
	makeQueueItem(pair, amountfl, pricefl)
	return nil
}

func marketOrder(userInfo user, pair, amount string) error {
	// convert amount to float64
	// get unique order id
	// execute
	return nil
}

func cancelOrder(userInfo user, pair, oid string) error {
	// convert oid to OrderID
	// remove order from that pair queue
	return nil
}
