package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func Contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

func Hash(value []byte, key string) string {
	h := sha256.New()
	h.Write(value)
	res := h.Sum(nil)
	return hex.EncodeToString(res)
}
