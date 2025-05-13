package utility

import (
	"math/rand"
)

const (
	digits     = 6
	rs2Letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func NewExposeId(prefix string) string {
	b := make([]byte, digits)
	for i := range b {
		b[i] = rs2Letters[rand.Intn(len(rs2Letters))]
	}
	return prefix + string(b)
}
