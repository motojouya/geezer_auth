package testUtility

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

type LocalerMock struct {
	t *testing.T
	randomString string
	expectLength int
	expectSource string
	uuid         uuid.UUID
	uuidErr      error
	now          time.Time
}

func NewLocalerMock(t *testing.T, randomString string, expectLength int, expectSource string, uuid uuid.UUID, uuidErr error, now time.Time) LocalerMock {
	return LocalerMock{
		t:            t,
		randomString: randomString,
		expectLength: expectLength,
		expectSource: expectSource,
		uuid:         uuid,
		uuidErr:      uuidErr,
		now:          now,
	}
}

func (mock LocalerMock) GenerateRamdomString(length int, source string) string {
	if mock.expectLength != 0 && length != mock.expectLength {
		mock.t.Errorf("Expected length %d, got %d", mock.expectLength, length)
	}
	if mock.expectSource != "" && source != mock.expectSource {
		mock.t.Errorf("Expected source %s, got %s", mock.expectSource, source)
	}
	return mock.randomString
}

func (mock LocalerMock) GenerateUUID() (uuid.UUID, error) {
	if mock.uuidErr != nil {
		return uuid.UUID{}, mock.uuidErr
	} else {
		return mock.uuid, nil
	}
}

func (mock LocalerMock) GetNow() time.Time {
	return mock.now
}
