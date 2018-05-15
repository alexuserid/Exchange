package main

import (
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type session struct {
	id b32
}

var (
	mapSidUid = make(map[b32]session)
)

func checker(em, pass string) (b32, bool) {
	uid, ok := mapEmailUid[em]
	if !ok {
		return b32{}, false
	}
	err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pass))
	if err != nil {
		return b32{}, false
	}
	return uid, true
}

func newSid(email []string, password []string) (string, errorc) {
	em := strings.Join(email, "")
	pass := strings.Join(password, "")

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	uid, ok := checker(em, pass)
	if !ok {
		return "", errWrongEmailPassword
	}

	sid, err := getUniqueId(markerSid)
	if err != errNo {
		return "", err
	}
	mapSidUid[sid] = session{id: uid}

	return b32ToString(sid), errNo
}
