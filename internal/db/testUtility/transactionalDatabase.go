package testUtility

type TransactionalDatabaseMock struct {
	FakeBegin    func() error
	FakeRollback func() error
	FakeCommit   func() error
	FakeClose    func() error
}

func (mock TransactionalDatabaseMock) Begin() error {
	return mock.FakeBegin()
}

func (mock TransactionalDatabaseMock) Rollback() error {
	return mock.FakeRollback()
}

func (mock TransactionalDatabaseMock) Commit() error {
	return mock.FakeCommit()
}

func (mock TransactionalDatabaseMock) Close() error {
	return mock.FakeClose()
}

type TransactionCalledCount struct {
	BeginCalled    int
	RollbackCalled int
	CommitCalled   int
	CloseCalled    int
}

func GetTransactionalDatabaseMock(calledCount *TransactionCalledCount) TransactionalDatabaseMock {
	return TransactionalDatabaseMock{
		FakeBegin: func() error {
			calledCount.BeginCalled++
			return nil
		},
		FakeRollback: func() error {
			calledCount.RollbackCalled++
			return nil
		},
		FakeCommit: func() error {
			calledCount.CommitCalled++
			return nil
		},
		FakeClose: func() error {
			calledCount.CloseCalled++
			return nil
		},
	}
}
