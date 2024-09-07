package ratelimit

import "math/rand"

func GetLlmResult() int {
	return rand.Intn(5)
}
