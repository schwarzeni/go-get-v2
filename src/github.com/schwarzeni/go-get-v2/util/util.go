package util

import (
	"log"
	"math/rand"
	"time"
)

// 只是模拟http延迟
func SleepAtRandomTime() {
	r := rand.Intn(2000)
	rateLimit := time.Tick(time.Duration(r) * time.Millisecond)
	<-rateLimit
}

func LogP(v interface{}) {
	log.Println(">> ", v)
}

func LogE(v interface{}) {
	log.Println(">>", v)
}
