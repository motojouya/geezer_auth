package command

import (
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type AddPasswordQuery interface {
	AddPassword(userPassword *transfer.UserPassword, now time.Time) (*transfer.UserPassword, error)
}

func AddPassword(executer gorp.SqlExecutor, userPassword *transfer.UserPassword, now time.Time) (*transfer.UserPassword, error) {
	var expireErr = ExpirePassword(executer, userPassword.UserPersistKey, now)
	if expireErr != nil {
		return &transfer.UserPassword{}, expireErr
	}

	var insertErr = executer.Insert(userPassword)
	if insertErr != nil {
		return &transfer.UserPassword{}, insertErr
	}

	return userPassword, nil
}
