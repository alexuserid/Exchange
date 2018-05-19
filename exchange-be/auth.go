package main

import (
	"sync"

	"golang.org/x/crypto/bcrypt"
)

//making user after registration
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



//making session after login
type SessionID [idl]byte

type session struct {
	id UserID
}

var (
	mapSidSession = make(map[SessionID]session)
	mutexGetSid   = &sync.RWMutex{}
)

func EmailAndPassChecker(em, pass string) (UserID, bool) {
	uid, ok := mapEmailUid[em]
	if !ok {
		return UserID{}, false
	}
	if err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pass)); err != nil {
		return UserID{}, false
	}
	return uid, true
}

func getSid() (SessionID, *errorc) {
	for i := 0; ; i++ {
		randoms, errc := getRandoms32()
		if errc != nil {
			return SessionID{}, errc
		}
		hb := hexMakerb32(randoms)
		var id SessionID
		copy(id[:], hb[:])

		if _, ok := mapSidSession[id]; !ok {
			return id, nil
		}
		if i > 100 {
			return SessionID{}, errNoToken
		}
	}
}

func newSid(email string, password string) (string, *errorc) {
	mutexGetSid.Lock()
	defer mutexGetSid.Unlock()

	uid, ok := EmailAndPassChecker(email, password)
	if !ok {
		return "", errWrongEmailPassword
	}

	sid, errc := getSid()
	if errc != nil {
		return "", errc
	}

	mapSidSession[sid] = session{id: uid}
	return string(sid[:]), nil
}
