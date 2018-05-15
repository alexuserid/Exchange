package main

import (
	"crypto/rand"
	"encoding/hex"
)

type b32 [32]byte

const (
	markerUid = 1
	markerSid = 2
	markerOid = 3
)

func getRandoms(length int) ([]byte, errorc) {
	randoms := make([]byte, length)
	n, err := rand.Read(randoms)
	if err != nil {
		return nil, errGetRandom
	}
	if n != length {
		return nil, errLength
	}
	return randoms, errNo
}

func hexMakerb32(b []byte) b32 {
	hexed := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(hexed, b)

	var result b32
	copy(result[:32], hexed)
	return result
}

func getUniqueId(marker int) (b32, errorc) {
	for i := 0; ; i++ {
		randoms, err := getRandoms(32)
		if err != errNo {
			return b32{}, err
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
			return id, errNo
		}

		if i >= 100 {
			return b32{}, errNoToken
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
