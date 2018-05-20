package main

import (
	"sync"

	"github.com/starius/status"
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
	mutexGetUid sync.Mutex
)

func getUid() (UserID, error) {
	for i := 0; ; i++ {
		randoms, err := getRandoms32()
		if err != nil {
			return UserID{}, status.Format("getRansom32: %v", err)
		}
		hb := hexMakerb32(randoms)
		var id UserID
		copy(id[:], hb[:])

		if _, ok := mapUidUser[id]; !ok {
			return id, nil
		}
		if i > 100 {
			return UserID{}, status.WithCode(StatusInternalServerError, "no free token after 100 iteration")
		}
	}
}

func newUser(email string, password string) error {
	cryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	mutexGetUid.Lock()
	defer mutexGetUid.Unlock()

	if _, ok := mapEmailUid[email]; ok {
		return status.WithCode(StatusBadRequest, "The email is already exist")
	}
	uid, err := getUid()
	if err != nil {
		return status.Format("getUid: %v", err)
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
	mutexGetSid   sync.Mutex
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

func getSid() (SessionID, error) {
	for i := 0; ; i++ {
		randoms, err := getRandoms32()
		if err != nil {
			return SessionID{}, status.Format("getRandom32: %v", err)
		}
		hb := hexMakerb32(randoms)
		var id SessionID
		copy(id[:], hb[:])

		if _, ok := mapSidSession[id]; !ok {
			return id, nil
		}
		if i > 100 {
			return SessionID{}, status.WithCode(StatusInternalServerError, "no free token after 100 iteration")
		}
	}
}

func newSid(email string, password string) (string, error) {
	mutexGetSid.Lock()
	defer mutexGetSid.Unlock()

	uid, ok := EmailAndPassChecker(email, password)
	if !ok {
		return "", status.WithCode(StatusBadRequest, "Wrong email or password")
	}

	sid, err := getSid()
	if err != nil {
		return "", status.Format("getSid: %v", err)
	}

	mapSidSession[sid] = session{id: uid}
	return string(sid[:]), nil
}
