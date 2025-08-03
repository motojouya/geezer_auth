package testUtility

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

type SqlExecutor interface {
	WithContext(ctx context.Context) SqlExecutor
	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectInt(query string, args ...interface{}) (int64, error)
	SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)
	SelectFloat(query string, args ...interface{}) (float64, error)
	SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error)
	SelectStr(query string, args ...interface{}) (string, error)
	SelectNullStr(query string, args ...interface{}) (sql.NullString, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}


type SqlExecutorMock struct {
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

func (mock SqlExecutorMock) WithContext(ctx context.Context) SqlExecutor {
	return mock
}

func (mock SqlExecutorMock) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return nil, nil
}

func (mock SqlExecutorMock) Insert(list ...interface{}) error {
	return nil
}

func (mock SqlExecutorMock) Update(list ...interface{}) (int64, error) {
	return 0, nil
}

func (mock SqlExecutorMock) Delete(list ...interface{}) (int64, error) {
	return 0, nil
}

func (mock SqlExecutorMock) Exec(query string, args ...interface{}) (sql.Result, error) {
}

func (mock SqlExecutorMock) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (mock SqlExecutorMock) SelectInt(query string, args ...interface{}) (int64, error) {
	return 0, nil
}

func (mock SqlExecutorMock) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	return nil, nil
}

func (mock SqlExecutorMock) SelectFloat(query string, args ...interface{}) (float64, error) {
	return 0, nil
}

func (mock SqlExecutorMock) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	return nil, nil
}

func (mock SqlExecutorMock) SelectStr(query string, args ...interface{}) (string, error {
	return "", nil
}
)
func (mock SqlExecutorMock) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	return nil, nil
}

func (mock SqlExecutorMock) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return nil
}

func (mock SqlExecutorMock) Query(query string, args ...interface{}) (*sql.Rows, error) {

}

func (mock SqlExecutorMock) QueryRow(query string, args ...interface{}) *sql.Row {
}
