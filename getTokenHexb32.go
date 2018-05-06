package main

import (
	"crypto/rand"
	"encoding/hex"
)

type b32 [32]byte

func getToken(length int) ([]byte, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func hexMakerb32(b []byte) b32 {
	hexed := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(hexed, b)

	var result b32
	copy(result[:32], hexed)
	return result
}
