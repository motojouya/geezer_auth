package io

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Local interface {
	GenerateRamdomString(length int, source string) string
	GenerateUUID() (uuid.UUID, error)
	GetNow() time.Time
}

type local struct {}

func CreateLocal() Local {
	return &local{}
}

func (l local) GenerateRamdomString(length int, source string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func (l local) GenerateUUID() (uuid.UUID, error) {
	var uuidValue, err = uuid.NewUUID()
	if err != nil {
		var zero = uuid.UUID{}
		return zero, err
	}

	return uuidValue, nil
}

func (l local) GetNow() time.Time {
	return time.Now()
}
