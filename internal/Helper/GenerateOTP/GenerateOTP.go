package helper

import (
	"math/rand"
	"time"
)

// Exported function
func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}
