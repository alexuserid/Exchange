package main

import (
	"sync"

	"github.com/starius/status"
	"golang.org/x/crypto/bcrypt"
)

//Make user after registration.
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
	mutexAuth   sync.Mutex
)

func getUid() (UserID, error) {
	for {
		randoms, err := getRandoms32()
		if err != nil {
			return UserID{}, status.Format("getRansom32: %v", err)
		}
		hb := makeHex(randoms)
		var id UserID
		copy(id[:], hb[:])

		if _, has := mapUidUser[id]; !has {
			return id, nil
		}
	}
}

func newUser(email string, password string) error {
	cryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	mutexAuth.Lock()
	defer mutexAuth.Unlock()

	if _, has := mapEmailUid[email]; has {
		return status.WithCode(statusConflict, "The email is already exist")
	}
	uid, err := getUid()
	if err != nil {
		return status.Format("getUid: %v", err)
	}

	mapEmailUid[email] = uid
	mapUidUser[uid] = user{email: email, password: cryptedPass, money: make(map[string]float64)}
	return nil
}

//Make session after login.
type SessionID [idl]byte

type session struct {
	id UserID
}

var (
	mapSidSession = make(map[SessionID]session)
	mutexAuth     sync.Mutex
)

func emailAndPassChecker(em, pass string) (UserID, bool) {
	mutexAuth.Lock()
	defer mutexAuth.Unlock()

	uid, has := mapEmailUid[em]
	if !has {
		return UserID{}, false
	}
	if err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pass)); err != nil {
		return UserID{}, false
	}
	return uid, true
}

func getSid() (SessionID, error) {
	for {
		randoms, err := getRandoms32()
		if err != nil {
			return SessionID{}, status.Format("getRandom32: %v", err)
		}
		hb := makeHex(randoms)
		var id SessionID
		copy(id[:], hb[:])

		if _, has := mapSidSession[id]; !has {
			return id, nil
		}
	}
}

func newSid(email string, password string) (string, error) {
	mutexAuth.Lock()
	defer mutexAuth.Unlock()

	uid, ok := emailAndPassChecker(email, password)
	if !ok {
		return "", status.WithCode(statusBadRequest, "Wrong email or password")
	}

	sid, err := getSid()
	if err != nil {
		return "", status.Format("getSid: %v", err)
	}

	mapSidSession[sid] = session{id: uid}
	return string(sid[:]), nil
}
