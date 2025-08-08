package testUtility

import (
	"github.com/google/uuid"
	"time"
)

type LocalerMock struct {
	FakeGenerateRamdomString func(length int, source string) string
	FakeGenerateUUID         func() (uuid.UUID, error)
	FakeGetNow               func() time.Time
}

func (mock LocalerMock) GenerateRamdomString(length int, source string) string {
	return mock.FakeGenerateRamdomString(length, source)
}

func (mock LocalerMock) GenerateUUID() (uuid.UUID, error) {
	return mock.FakeGenerateUUID()
}

func (mock LocalerMock) GetNow() time.Time {
	return mock.FakeGetNow()
}

// func CreateLocalerMock() *LocalerMock {
// 	var generateRamdomString = func(length int, source string) string {
// 		return "mockedRandomString"
// 	}
// 	var generateUUID = func() (uuid.UUID, error) {
// 		return uuid.UUID{}, nil
// 	}
// 	var getNow = func() time.Time {
// 		return time.Now()
// 	}
// 	return &LocalerMock{
// 		generateRamdomString: generateRamdomString,
// 		generateUUID:         generateUUID,
// 		getNow:               getNow,
// 	}
// }
