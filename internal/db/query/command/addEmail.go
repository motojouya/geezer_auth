package command

import (
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type AddEmailQuery interface {
	AddEmail(userEmail *transfer.UserEmail, now time.Time) (*transfer.UserEmail, error)
}

func AddEmail(executer gorp.SqlExecutor, userEmail *transfer.UserEmail, now time.Time) (*transfer.UserEmail, error) {
	var insertErr = executer.Insert(userEmail)
	if insertErr != nil {
		return &transfer.UserEmail{}, insertErr
	}

	var expireErr = ExpireEmail(executer, userEmail.UserPersistKey, userEmail.PersistKey, true, now)
	if expireErr != nil {
		return &transfer.UserEmail{}, expireErr
	}

	return userEmail, nil
}
