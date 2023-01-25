package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func StartAsyncTimer(tickChan <-chan time.Time, doneChan <-chan struct{}, f func()) {
	for {
		select {
		case <-tickChan:
			go f()
		case <-doneChan:
			//log.Printf("ret")
			return
		}
	}
}
