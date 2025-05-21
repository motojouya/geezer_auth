package here

import (
	"math/rand"
	"github.com/google/uuid"
	"time"
)

func GenerateLargeCharactorString(length int, source string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func GenerateUUID() (UUID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return UUID(""), err
	}

	return uuid, nil
}

func GetNow() time.Time {
	return time.Now()
}
