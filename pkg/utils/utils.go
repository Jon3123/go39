package utils

import (
	"math/rand"
	"time"
)

var rInit = false

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	if rInit == false {
		rand.Seed(time.Now().UnixNano())
		rInit = true
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//GenerateConnectionID TODO return a random connection id
func GenerateConnectionID() string {
	return "conn-" + randStringBytes(10)
}
