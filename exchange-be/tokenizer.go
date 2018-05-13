package main

import (
	"crypto/rand"
	"encoding/json"
	"encoding/hex"
	"errors"
	"net/http"
)

type b32 [32]byte

const (
	markerUid = 1
	markerSid = 2
	markerOid = 3
)

func getRandoms(length int) ([]byte, error) {
	randoms := make([]byte, length)
	n, err := rand.Read(randoms)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, errors.New("Mismatched length.")
	}
	return randoms, nil
}

func hexMakerb32(b []byte) b32 {
	hexed := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(hexed, b)

	var result b32
	copy(result[:32], hexed)
	return result
}

func getUniqueId(w http.ResponseWriter, marker int) (b32, error) {
	for i := 0; ; i++ {
		randoms, err := getRandoms(32)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(jsons{Err: "Internal server error. Can't create a new token.Please, contact support."})
			return b32{0}, err
		}

		id := hexMakerb32(randoms)
		var ok bool
		switch marker {
		case 1:
			_, ok = mapUidUser[id]
		case 2:
			_, ok = mapSidUid[id]
		case 3:
			_, ok = mapOidOrder[id]
		}
		if !ok {
			return id, nil
		}

		if i >= 100 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(jsons{Err: "There is no free tokens. Please, try again, or contact support."})
			return b32{0}, errors.New("No free tokens.")
		}
	}
}

func stringToB32(s string) b32 {
	var token b32
	copy(token[:], []byte(s))
	return token
}

func b32ToString(token b32) string {
	tokenBytes := make([]byte, len(token))
	copy(tokenBytes, token[:])
	return string(tokenBytes)
}
