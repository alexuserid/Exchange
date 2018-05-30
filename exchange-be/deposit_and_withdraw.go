package main

import (
	"net/http"
	"strconv"

	"github.com/starius/status"
)

func depositAndWithdraw(userInfo *user, operation, currency, amount string) error {
	amf, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return status.WithCode(http.StatusBadRequest, "Wrong amount format: ParseFloat: %v", err)
	}

	mutex.Lock()
	defer mutex.Unlock()

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
