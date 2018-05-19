package main

import (
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
	mutexGetUid = &sync.RWMutex{}
)

func getUid() (UserID, *errorc) {
	for i := 0; ; i++ {
		randoms, errc := getRandoms32()
		if errc != nil {
			return UserID{}, errc
		}
		hb := hexMakerb32(randoms)
		var id UserID
		copy(id[:], hb[:])

		if _, ok := mapUidUser[id]; !ok {
			return id, nil
		}
		if i > 100 {
			return UserID{}, errNoToken
		}
	}
}

func newUser(email string, password string) *errorc {
	cryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fullError(errHashGen, err)
	}

	mutexGetUid.Lock()
	defer mutexGetUid.Unlock()

	if _, ok := mapEmailUid[email]; ok {
		return errExistingEmail
	}
	uid, errc := getUid()
	if errc != nil {
		return errc
	}

	mapEmailUid[email] = uid
	mapUidUser[uid] = user{email: email, password: cryptedPass, money: make(map[string]float64)}
	return nil
}
