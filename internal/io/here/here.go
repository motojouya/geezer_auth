package here

import (
	"math/rand"
	"github.com/google/uuid"
	"time"
)

type Here interface {
	GenerateRamdomString(length int, source string) string
	GenerateUUID() (UUID, error)
	GetNow() time.Time
}

type here interface {}

func CreateHere() Here {
	return &here(interface{})
}

func (h here) GenerateRamdomString(length int, source string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func (h here) GenerateUUID() (UUID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return UUID(""), err
	}

	return uuid, nil
}

func (h here) GetNow() time.Time {
	return time.Now()
}
