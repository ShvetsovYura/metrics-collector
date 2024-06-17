package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// Contains, проверяет содержится ли указанная строка в слайсе строк.
func Contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}

	return false
}

// Hash, вычисляет хэш-строку для переданного значени
//
// BUG(ShvetsovYura): ключ не используется
func Hash(value []byte, _ string) string {
	h := sha256.New()
	h.Write(value)
	res := h.Sum(nil)

	return hex.EncodeToString(res)
}
