package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// generate a random number between min and max number
func RanomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate a random string
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generate owner name
func RandomOwnerName() string {
	return RandomString(6)
}

// generate money
func RandomMoney() int64 {
	return RanomInt(0, 1000)
}

// generate currency
func RandomCurrrency() string {
	currencys := []string{USD, TK, EUR, CAD}
	n := len(currencys)
	return currencys[rand.Intn(n)]
}
