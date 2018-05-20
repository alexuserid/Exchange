package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/starius/status"
)

var mutexDepositAndWithdraw sync.Mutex

func getUserInfo(r *http.Request) (user, error) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return user{}, status.WithCode(StatusBadRequest, "You are not logged in: %v", err)
	}
	var sid SessionID
	copy(sid[:], []byte(cookie.Value))
	uid := mapSidSession[sid]
	return mapUidUser[uid.id], nil
}

func depositAndWithdraw(userInfo user, operation, currency, amount string) error {
	amf, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return status.WithCode(StatusBadRequest, "Wrong amount format: ParseFloat: %v", err)
	}

	mutexDepositAndWithdraw.Lock()
	defer mutexDepositAndWithdraw.Unlock()

	if userInfo.money == nil {
		return status.WithCode(StatusBadRequest, "You are not logged in: userInfo.money == nil")
	}

	switch operation {
	case "deposit":
		userInfo.money[currency] += amf
	case "withdraw":
		if res := userInfo.money[currency] - amf; res < 0 {
			return status.WithCode(StatusBadRequest, "You can't withdraw more money than you have")
		}
		userInfo.money[currency] -= amf
	}
	return nil
}
