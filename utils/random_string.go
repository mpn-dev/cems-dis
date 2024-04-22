package utils

import (
	"math/rand"
	"time"
)

const (
	RANDOM_CHAR_MAP = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for range length {
		cpos := rand.Intn(len(RANDOM_CHAR_MAP))
		s = s + RANDOM_CHAR_MAP[cpos:cpos + 1]
	}
	return s
}

func RandomString16() string {
	return GenerateRandomString(16)
}

func RandomString32() string {
	return GenerateRandomString(32)
}

func RandomString64() string {
	return GenerateRandomString(64)
}
