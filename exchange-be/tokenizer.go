package main

import (
	"crypto/rand"
	"encoding/hex"
)

const idl = 32

func getRandoms32() ([]byte, *errorc) {
	randoms := make([]byte, idl)
	n, err := rand.Read(randoms)
	if err != nil {
		return nil, fullError(errGetRandom, err)
	}
	if n != idl {
		return nil, errLength
	}
	return randoms, nil
}

func hexMakerb32(b []byte) [idl]byte {
	hexed := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(hexed, b)

	var result [idl]byte
	copy(result[:], hexed)
	return result
}
