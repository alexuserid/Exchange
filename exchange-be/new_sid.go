package main

import (
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type SessionID [idl]byte

type session struct {
	id UserID
}

var (
	mapSidUid = make(map[SessionID]session)
)

func EmailAndPassChecker(em, pass string) (UserID, bool) {
	uid, ok := mapEmailUid[em]
	if !ok {
		return UserID{}, false
	}
	err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pass))
	if err != nil {
		return UserID{}, false
	}
	return uid, true
}

func getSid() (SessionID, errorc) {
	for i:=0; ; i++ {
		randoms, err := getRandoms32()
		if err != errNo {
			return SessionID{}, err
		}
		hb := hexMakerb32(randoms)
		var id SessionID
		copy(id[:], hb[:])

		_, ok := mapSidUid[id]
		if !ok {
			return id, errNo
		}
		if i > 100 {
			return SessionID{}, errNoToken
		}
	}
}

func newSid(email []string, password []string) (string, errorc) {
	em := strings.Join(email, "")
	pass := strings.Join(password, "")

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	uid, ok := EmailAndPassChecker(em, pass)
	if !ok {
		return "", errWrongEmailPassword
	}

	sid, err := getSid()
	if err != errNo {
		return "", err
	}

	mapSidUid[sid] = session{id: uid}
	return string(sid[:]), errNo
}
