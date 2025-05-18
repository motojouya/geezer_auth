package here

import (
	"math/rand"
	"github.com/google/uuid"
	"time"
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

func GenerateUUID() (string, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func GetNow() time.Time {
	return time.Now()
}
