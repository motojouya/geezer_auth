package command

import (
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type VerifyEmailQuery interface {
	VerifyEmail(userEmail *transfer.UserEmail, now time.Time) (*transfer.UserEmail, error)
}

func VerifyEmail(executer gorp.SqlExecutor, userEmail *transfer.UserEmail, now time.Time) (*transfer.UserEmail, error) {
	var _, insertErr = executer.Update(userEmail)
	if insertErr != nil {
		return &transfer.UserEmail{}, insertErr
	}

	var expireErr = ExpireEmail(executer, userEmail.UserPersistKey, userEmail.PersistKey, false, now)
	if expireErr != nil {
		return &transfer.UserEmail{}, expireErr
	}

	return userEmail, nil
}
