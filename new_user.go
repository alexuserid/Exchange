package main

import (
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

func getUniqueUserId(w http.ResponseWriter) (b32, error) {
	for i := 0; ; i++ {
		token, err := getToken(32)
		if err != nil {
			http.Error(w, "Internal server error. Can't create a new token.Please, contact support.", http.StatusInternalServerError)
			return b32{0}, err
		}
		id := hexMakerb32(token)
		_, ok := mapUidUser[id]
		if !ok {
			return id, nil
		}
		if i >= 100 {
			http.Error(w, "There is no free tokens. Please, try again, or contact support.", http.StatusInternalServerError)
			return b32{0}, errors.New("No free tokens")
		}
	}
}

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
		return errors.New("Existing email")
	}

	uid, err := getUniqueUserId(w)
	if err != nil {
		return err
	}

	mapEmailUid[emj] = uid
	mapUidUser[uid] = user{email: emj, password: bytePass}
	return nil
}
