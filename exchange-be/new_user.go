package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	email    string
	password []byte
	money    map[string]float64
	history  []deals
}

type deals struct {
	date      time.Time
	currency1 string
	currency2 string
	amountC1  float64
	amountC2  float64
	price     float64
}

type jsons struct {
	err string
}

var (
	mapEmailUid = make(map[string]b32)
	mapUidUser  = make(map[b32]user)
)

func newUser(email []string, password []string, w http.ResponseWriter) error {
	em := strings.Join(email, "")
	pass := []byte(strings.Join(password, ""))

	cryptedPass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsons{err: "Internal server error. Can't generate hash from password. Please, contact support."})
		return err
	}

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := mapEmailUid[em]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsons{err: "The email is already exist."})
		return errors.New("Existing email")
	}

	uid, err := getUniqueId(w, markerUid)
	if err != nil {
		return err
	}

	mapEmailUid[em] = uid
	mapUidUser[uid] = user{email: em, password: cryptedPass, money: make(map[string]float64)}
	return nil
}
