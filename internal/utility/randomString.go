package utility

import (
	"math/rand"
)

const (
	exposeIdDigit = 6
	exposeIdChar  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func NewExposeId(prefix string) string {
	b := make([]byte, exposeIdDigit)
	for i := range b {
		b[i] = exposeIdChar[rand.Intn(len(exposeIdChar))]
	}
	return prefix + string(b)
}
