package main

import (
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	email    string
	password []byte
	money    map[string]float64
	history  []b32
}

var (
	mapEmailUid = make(map[string]b32)
	mapUidUser  = make(map[b32]user)
)

func newUser(email []string, password []string) errorc {
	em := strings.Join(email, "")
	pass := []byte(strings.Join(password, ""))

	cryptedPass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return errHashGen
	}

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := mapEmailUid[em]
	if ok {
		return errExistingEmail
	}

	uid, errc := getUniqueId(markerUid)
	if errc != errNo {
		return errc
	}

	mapEmailUid[em] = uid
	mapUidUser[uid] = user{email: em, password: cryptedPass, money: make(map[string]float64)}
	return errNo
}
