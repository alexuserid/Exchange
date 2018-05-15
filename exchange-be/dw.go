package main

import (
	"net/http"
	"strconv"
	"sync"
)

func getUserInfo(r *http.Request) (user, errorc) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return user{}, errNoCookie
	}
	uid := mapSidUid[stringToB32(cookie.Value)]
	return mapUidUser[uid.id], errNo
}

func dw(userInfo user, op, cur, am string) errorc {
	amf, err := strconv.ParseFloat(am, 64)
	if err != nil {
		return errParseFloat
	}

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if userInfo.money == nil {
		return errLogin
	}

	switch op {
	case "deposit":
		userInfo.money[cur] += amf
	case "withdraw":
		if res := userInfo.money[cur] - amf; res < 0 {
			return errWithdraw
		}
		userInfo.money[cur] -= amf
	}
	return errNo
}
