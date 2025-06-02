package io

import (
	"math/rand"
	"github.com/google/uuid"
	"github.com/caarlos0/env/v11"
	"time"
)

type Local interface {
	GenerateRamdomString(length int, source string) string
	GenerateUUID() (UUID, error)
	GetNow() time.Time
	GetEnv(object *interface{}) error
}

type local interface {}

func CreateLocal() Local {
	return &local(interface{})
}

func (l local) GenerateRamdomString(length int, source string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func (l local) GenerateUUID() (UUID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return UUID(""), err
	}

	return uuid, nil
}

func (l local) GetNow() time.Time {
	return time.Now()
}

func (l local) GetEnv[T any]() (T, error) {
	return env.ParseAs[T]();
}
