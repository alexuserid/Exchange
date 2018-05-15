package main

import (
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserID [idl]byte

type user struct {
	email    string
	password []byte
	money    map[string]float64
	history  []OrderID
}

var (
	mapEmailUid = make(map[string]UserID)
	mapUidUser  = make(map[UserID]user)
)

func getUid() (UserID, errorc) {
	for i:=0; ; i++ {
		randoms, err := getRandoms32()
		if err != errNo {
			return UserID{}, err
		}
		hb := hexMakerb32(randoms)
		var id UserID
		copy(id[:], hb[:])

		_, ok := mapUidUser[id]
		if !ok {
			return id, errNo
		}
		if i > 100 {
			return UserID{}, errNoToken
		}
	}
}

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
	uid, errc := getUid()
	if errc != errNo {
		return errc
	}

	mapEmailUid[em] = uid
	mapUidUser[uid] = user{email: em, password: cryptedPass, money: make(map[string]float64)}
	return errNo
}
