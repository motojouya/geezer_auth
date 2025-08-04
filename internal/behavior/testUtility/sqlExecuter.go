package testUtility

type SqlExecutorMock struct {}

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
	return nil, nil
}

func (mock SqlExecutorMock) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (mock SqlExecutorMock) SelectInt(query string, args ...interface{}) (int64, error) {
	return 0, nil
}

func (mock SqlExecutorMock) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	return sql.NullInt64{}, nil
}

func (mock SqlExecutorMock) SelectFloat(query string, args ...interface{}) (float64, error) {
	return 0, nil
}

func (mock SqlExecutorMock) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	return sql.NullFloat64{}, nil
}

func (mock SqlExecutorMock) SelectStr(query string, args ...interface{}) (string, error {
	return "", nil
}

func (mock SqlExecutorMock) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	return sql.NullString{}, nil
}

func (mock SqlExecutorMock) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return nil
}

func (mock SqlExecutorMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (mock SqlExecutorMock) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}
