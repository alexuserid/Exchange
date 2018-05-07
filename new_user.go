package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	email    string
	password []byte
	wallet   []money
	history  []deals
}

type money struct {
	name   string
	amount float64
}

type deals struct {
	date      time.Time
	currency1 string
	currency2 string
	amountC1  float64
	amountC2  float64
	price     float64
}

var (
	mapEmailUid = make(map[string]b32)
	mapUidUser  = make(map[b32]user)
)

func newUser(email []string, password []string, w http.ResponseWriter) error {
	emj := strings.Join(email, "")
	pass := []byte(strings.Join(password, ""))

	bytePass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error. Can't generate hash from password. Please, contact support.", http.StatusInternalServerError)
		return err
	}

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()
	
	_, ok := mapEmailUid[emj]
	if ok {
		http.Error(w, "The email is already exist.", http.StatusBadRequest)
	}

	uid, err := getUniqueId(markerUid, w)
	if err != nil {
		return err
	}

	mapEmailUid[emj] = uid
	mapUidUser[uid] = user{email: emj, password: bytePass}
	return nil
}
