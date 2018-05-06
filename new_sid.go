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

func getUniqueSessionId(w http.ResponseWriter) (b32, error) {
	for i := 0; ; i++ {
		token, err := getToken(32)
		if err != nil {
			http.Error(w, "Internal server error. Can't create a new token.Please, contact support.", http.StatusInternalServerError)
			return b32{0}, err
		}
		id := hexMakerb32(token)
		_, ok := mapSidUid[id]
		if !ok {
			return id, nil
		}
		if i >= 100 {
			http.Error(w, "There is no free tokens. Please, try again, or contact support.", http.StatusInternalServerError)
			return b32{0}, errors.New("No free tokens")
		}
	}
}

func newSid(email []string, password []string, w http.ResponseWriter) (string, error) {
	emj := strings.Join(email, "")
	pasj := strings.Join(password, "")

	mutex := &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	uid, ok := mapEmailUid[emj]
	if !ok {
		http.Error(w, "Wrong email.", http.StatusBadRequest)
		return "", errors.New("Existing email")
	}
	err := bcrypt.CompareHashAndPassword(mapUidUser[uid].password, []byte(pasj))
	if err != nil {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return "", errors.New("Wrong password.")
	}
	sid, err := getUniqueSessionId(w)
	if err != nil {
		return "", err
	}

	mapSidUid[sid] = session{uid: uid}

	uidBytes := make([]byte, len(uid))
	copy(uidBytes, uid[:])
	return string(uidBytes), nil
}
