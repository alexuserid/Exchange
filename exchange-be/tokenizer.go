package main

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/starius/status"
)

const idl = 32

func getRandoms32() ([]byte, error) {
	randoms := make([]byte, idl)
	n, err := rand.Read(randoms)
	if err != nil {
		return nil, status.WithCode(statusInternalServerError, "rand.Read: %v", err)
	}
	if n != idl {
		return nil, status.WithCode(statusInternalServerError, "rand.Read: length error")
	}
	return randoms, nil
}

func makeHex(b []byte) [idl]byte {
	hexed := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(hexed, b)

	var result [idl]byte
	copy(result[:], hexed)
	return result
}
