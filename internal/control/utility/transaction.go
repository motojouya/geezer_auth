package utility

import (
	"github.com/motojouya/geezer_auth/internal/db"
	pkgUser "github.com/motojouya/geezer_auth/pkg/core/user"
)

func RollbackWithError(database db.Transactional, err error) error {
	if rollbackErr := database.Rollback(); rollbackErr != nil {
		return rollbackErr
	}
	return err
}

// control処理の頭から最後までトランザクションとする場合に有効な関数。例えば、DBアクセスもするし、APIアクセスもして、トランザクションの粒度を操作したい場合は、control処理内でbegin/commitすべき
func Transact[C db.TransactionalDatabase, E any, R any](callback func (C, E, *pkgUser.User) (R, error)) func (C, E, *pkgUser.User) (R, error) {
	return func(control C, entry E, authentic *pkgUser.User) (R, error) {

		var zero R
		if err := control.Begin(); err != nil {
			return zero, err
		}

		result, err := callback(control, entry, authentic)

		if err != nil {
			if err := control.Rollback(); err != nil {
				return zero, err
			}
		} else {
			if err := control.Commit(); err != nil {
				return zero, err
			}
		}

		return result, err
	}
}
