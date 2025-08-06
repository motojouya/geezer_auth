package testUtility

import (
	"context"
	"database/sql"
	"github.com/go-gorp/gorp"
)

type SqlExecutorMock struct {
	FakeWithContext     func(ctx context.Context) gorp.SqlExecutor
	FakeGet             func(i interface{}, keys ...interface{}) (interface{}, error)
	FakeInsert          func(list ...interface{}) error
	FakeUpdate          func(list ...interface{}) (int64, error)
	FakeDelete          func(list ...interface{}) (int64, error)
	FakeExec            func(query string, args ...interface{}) (sql.Result, error)
	FakeSelect          func(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	FakeSelectInt       func(query string, args ...interface{}) (int64, error)
	FakeSelectNullInt   func(query string, args ...interface{}) (sql.NullInt64, error)
	FakeSelectFloat     func(query string, args ...interface{}) (float64, error)
	FakeSelectNullFloat func(query string, args ...interface{}) (sql.NullFloat64, error)
	FakeSelectStr       func(query string, args ...interface{}) (string, error)
	FakeSelectNullStr   func(query string, args ...interface{}) (sql.NullString, error)
	FakeSelectOne       func(holder interface{}, query string, args ...interface{}) error
	FakeQuery           func(query string, args ...interface{}) (*sql.Rows, error)
	FakeQueryRow        func(query string, args ...interface{}) *sql.Row
}

func (mock SqlExecutorMock) WithContext(ctx context.Context) gorp.SqlExecutor {
	return mock.FakeWithContext(ctx)
}

func (mock SqlExecutorMock) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return mock.FakeGet(i, keys...)
}

func (mock SqlExecutorMock) Insert(list ...interface{}) error {
	return mock.FakeInsert(list...)
}

func (mock SqlExecutorMock) Update(list ...interface{}) (int64, error) {
	return mock.FakeUpdate(list...)
}

func (mock SqlExecutorMock) Delete(list ...interface{}) (int64, error) {
	return mock.FakeDelete(list...)
}

func (mock SqlExecutorMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	return mock.FakeExec(query, args...)
}

func (mock SqlExecutorMock) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return mock.FakeSelect(i, query, args...)
}

func (mock SqlExecutorMock) SelectInt(query string, args ...interface{}) (int64, error) {
	return mock.FakeSelectInt(query, args...)
}

func (mock SqlExecutorMock) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	return mock.FakeSelectNullInt(query, args...)
}

func (mock SqlExecutorMock) SelectFloat(query string, args ...interface{}) (float64, error) {
	return mock.FakeSelectFloat(query, args...)
}

func (mock SqlExecutorMock) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	return mock.FakeSelectNullFloat(query, args...)
}

func (mock SqlExecutorMock) SelectStr(query string, args ...interface{}) (string, error) {
	return mock.FakeSelectStr(query, args...)
}

func (mock SqlExecutorMock) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	return mock.FakeSelectNullStr(query, args...)
}

func (mock SqlExecutorMock) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return mock.FakeSelectOne(holder, query, args...)
}

func (mock SqlExecutorMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return mock.FakeQuery(query, args...)
}

func (mock SqlExecutorMock) QueryRow(query string, args ...interface{}) *sql.Row {
	return mock.FakeQueryRow(query, args...)
}
