package main

import (
	"errors"
	"encoding/json"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsons{err: "Wrong email or password."})
		return "", errors.New("Wrong email or password")
	}

	sid, err := getUniqueId(w, markerSid)
	if err != nil {
		return "", err
	}
	mapSidUid[sid] = session{uid: uid}

	return b32ToString(uid), nil
}
