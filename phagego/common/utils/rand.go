package utils

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytesLower = "abcdefghijklmnopqrstuvwxyz"
const letterBytesLowerWithNum = "0123456789abcdefghijklmnopqrstuvwxyz"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	return randStr(n, letterBytes, time.Now().UnixNano())
}

func RandStringLower(n int) string {
	return randStr(n, letterBytesLower, time.Now().UnixNano())
}

func RandStrLowerWithNum(n int) string {
	return randStr(n, letterBytesLowerWithNum, time.Now().UnixNano())
}

func RandStrLowerWithNumSeed(n int, seed int64) string {
	return randStr(n, letterBytesLowerWithNum, seed)
}

func randStr(n int, letters string, seed int64) string {
	var src = rand.NewSource(seed)
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandNum(min, max int) int {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	return randNum
}
