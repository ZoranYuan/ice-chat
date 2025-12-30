package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenJoinCode(length int) string {
	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"
	randString := make([]byte, length)

	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	for i := range randString {
		randString[i] = charSet[r.Intn(len(charSet))]
	}

	return strings.ToUpper(string(randString))
}
