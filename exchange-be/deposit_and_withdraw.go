package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/starius/status"
)

var depositAndWithdrawMutex sync.Mutex

func getUserInfo(r *http.Request) (user, error) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return user{}, status.WithCode(http.StatusBadRequest, "You are not logged in: %v", err)
	}
	var sid sessionID
	copy(sid[:], []byte(cookie.Value))
	uid := mapSidSession[sid]
	return mapUidUser[uid.id], nil
}

func depositAndWithdraw(userInfo user, operation, currency, amount string) error {
	amf, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return status.WithCode(http.StatusBadRequest, "Wrong amount format: ParseFloat: %v", err)
	}

	depositAndWithdrawMutex.Lock()
	defer depositAndWithdrawMutex.Unlock()

	if userInfo.money == nil {
		return status.WithCode(http.StatusBadRequest, "You are not logged in: userInfo.money == nil")
	}

	switch operation {
	case "deposit":
		userInfo.money[currency] += amf
	case "withdraw":
		if res := userInfo.money[currency] - amf; res < 0 {
			return status.WithCode(http.StatusBadRequest, "You can't withdraw more money than you have")
		}
		userInfo.money[currency] -= amf
	}
	return nil
}
