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
	mapSidSession = make(map[SessionID]session)
	mutexGetSid = &sync.RWMutex{}
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

func getSid() (SessionID, *errorc) {
	for i:=0; ; i++ {
		randoms, errc := getRandoms32()
		if errc != nil {
			return SessionID{}, errc
		}
		hb := hexMakerb32(randoms)
		var id SessionID
		copy(id[:], hb[:])

		_, ok := mapSidSession[id]
		if !ok {
			return id, nil
		}
		if i > 100 {
			return SessionID{}, errNoToken
		}
	}
}

func newSid(email []string, password []string) (string, *errorc) {
	em := strings.Join(email, "")
	pass := strings.Join(password, "")

	mutexGetSid.Lock()
	defer mutexGetSid.Unlock()

	uid, ok := EmailAndPassChecker(em, pass)
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
