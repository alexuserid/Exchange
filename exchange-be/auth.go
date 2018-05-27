package main

import (
	"net/http"
	"sync"

	"github.com/starius/status"
	"golang.org/x/crypto/bcrypt"
)

/* This part creates user account after registration. */

type userID [idl]byte

type user struct {
	email    string
	password []byte
	money    map[string]float64
	history  []OrderID
}

var (
	mapEmailUid = make(map[string]userID)
	mapUidUser  = make(map[userID]user)
	authMutex   sync.Mutex
)

func newUid() (userID, error) {
	for {
		randoms, err := getRandoms32()
		if err != nil {
			return userID{}, status.Format("getRansom32: %v", err)
		}
		hb := toHex(randoms)
		var id userID
		copy(id[:], hb[:])

		if _, has := mapUidUser[id]; !has {
			return id, nil
		}
	}
}

func register(email string, password string) error {
	cryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	authMutex.Lock()
	defer authMutex.Unlock()

	if _, has := mapEmailUid[email]; has {
		return status.WithCode(http.StatusConflict, "The email already exist")
	}
	uid, err := newUid()
	if err != nil {
		return status.Format("newUid: %v", err)
	}

	mapEmailUid[email] = uid
	mapUidUser[uid] = user{email: email, password: cryptedPass, money: make(map[string]float64)}
	return nil
}

/* This part creates session for user after login. */

type sessionID [idl]byte

type session struct {
	id userID
}

var (
	mapSidSession = make(map[sessionID]session)
)

func newSid() (sessionID, error) {
	for {
		randoms, err := getRandoms32()
		if err != nil {
			return sessionID{}, status.Format("getRandom32: %v", err)
		}
		hb := toHex(randoms)
		var id sessionID
		copy(id[:], hb[:])

		if _, has := mapSidSession[id]; !has {
			return id, nil
		}
	}
}

func login(email string, password string) (string, error) {
	authMutex.Lock()
	defer authMutex.Unlock()

	emailAndPassCheck := func(email, pass string) (userID, bool) {
		uid, has := mapEmailUid[email]
		if !has {
			return userID{}, false
		}
		if err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pass)); err != nil {
			return userID{}, false
		}
		return uid, true
	}
	uid, ok := emailAndPassCheck(email, password)
	if !ok {
		return "", status.WithCode(http.StatusBadRequest, "Wrong email or password")
	}

	sid, err := newSid()
	if err != nil {
		return "", status.Format("newSid: %v", err)
	}

	mapSidSession[sid] = session{id: uid}
	return string(sid[:]), nil
}
