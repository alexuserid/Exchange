package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
)

func getUserInfo(r *http.Request) (user, error) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return user{}, err
	}
	uid := mapSidUid[stringToB32(cookie.Value)]
	return mapUidUser[uid.id], nil
}

func dw(w http.ResponseWriter, userInfo user, op, cur, am string) error {
	amf, err := strconv.ParseFloat(am, 64)
	if err != nil {
		return err
	}

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if userInfo.money == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsons{Err: "You are not logged in."})
		return errors.New("Nil map")
	}

	switch op {
	case "deposit":
		userInfo.money[cur] += amf
	case "withdraw":
		if res := userInfo.money[cur] - amf; res < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(jsons{Err: "Can't withdraw more money than you have."})
			return errors.New("Withdraw excess")
		}
		userInfo.money[cur] -= amf
	}
	return nil
}
