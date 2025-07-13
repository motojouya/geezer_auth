package command

import (
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type AddRefreshTokenQuery interface {
	AddRefreshToken(userRefreshToken *transfer.UserRefreshToken, now time.Time) (*transfer.UserRefreshToken, error)
}

func AddRefreshToken(executer gorp.SqlExecutor, userRefreshToken *transfer.UserRefreshToken, now time.Time) (*transfer.UserRefreshToken, error) {
	var expireErr = ExpireRefreshToken(executer, userRefreshToken.UserPersistKey, now)
	if expireErr != nil {
		return &transfer.UserRefreshToken{}, expireErr
	}

	var insertErr = executer.Insert(userRefreshToken)
	if insertErr != nil {
		return &transfer.UserRefreshToken{}, insertErr
	}

	return userRefreshToken, nil
}
