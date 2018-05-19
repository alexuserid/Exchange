package main

import (
	"net/http"
	"strconv"
	"sync"
)

func getUserInfo(r *http.Request) (user, *errorc) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return user{}, errNoCookie
	}
	var sid SessionID
	copy(sid[:], []byte(cookie.Value))
	uid := mapSidSession[sid]
	return mapUidUser[uid.id], nil
}

func depositAndWithdraw(userInfo user, operation, currency, amount string) *errorc {
	amf, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fullError(errParseFloat, err)
	}

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if userInfo.money == nil {
		return errLogin
	}

	switch operation {
	case "deposit":
		userInfo.money[currency] += amf
	case "withdraw":
		if res := userInfo.money[currency] - amf; res < 0 {
			return errWithdraw
		}
		userInfo.money[currency] -= amf
	}
	return nil
}
