package main

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type session struct {
	uid b32
}

var (
	mapSidUid = make(map[b32]session)
)

func checker(em, pass string) (b32, bool) {
	uid, ok := mapEmailUid[em]
	if !ok {
		return b32{0}, false
	}
	err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pass))
	if err != nil {
		return b32{0}, false
	}
	return uid, true
}

func newSid(email []string, password []string, w http.ResponseWriter) (string, error) {
	em := strings.Join(email, "")
	pass := strings.Join(password, "")

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	uid, ok := checker(em, pass)
	if !ok {
		http.Error(w, "Wrong email or password.", http.StatusBadRequest)
		return "", errors.New("Wrong email or password")
	}

	sid, err := getUniqueId(markerSid, w)
	if err != nil {
		return "", err
	}
	mapSidUid[sid] = session{uid: uid}

	uidBytes := make([]byte, len(uid))
	copy(uidBytes, uid[:])
	return string(uidBytes), nil
}
