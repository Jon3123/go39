package utils

import (
	"math/rand"
	"time"
)

var init bool = false

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	if init == false {
		rand.Seed(time.Now().UnixNano())
		init = true
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//GenerateConnectionID return a random connection id
func GenerateConnectionID() string {
	return "conn-" + randStringBytes(10)
}
