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
	sid, err := getUniqueId(markerSid, w)
	if err != nil {
		return "", err
	}

	mapSidUid[sid] = session{uid: uid}

	uidBytes := make([]byte, len(uid))
	copy(uidBytes, uid[:])
	return string(uidBytes), nil
}
