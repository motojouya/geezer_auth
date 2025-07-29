package user

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/service"
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/motojouya/geezer_auth/internal/db"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	entry "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
)

type UserRegisterDB interface {
	gorp.SqlExecutor
	db.Transactional
	userQuery.GetUserQuery
	userQuery.GetUserAuthenticQuery
	userQuery.GetUserEmailQuery
	commandQuery.AddPasswordQuery
	commandQuery.AddEmailQuery
	commandQuery.AddRefreshTokenQuery
}

type RegisterControl struct {
	Local io.Local
	DB UserRegisterDB
}

func NewRegisterControl(db UserRegisterDB, local io.Local) *RegisterControl {
	return &RegisterControl{
		DB:    db,
		Local: local,
	}
}

func CreateRegisterControl() (*RegisterControl, error) {
	var local = io.CreateLocal()
	var env = io.CreateEnvironment()

	var loader = service.GetLoader()
	var db, err = loader.LoadDatabase(env)
	if err != nil {
		return nil, err
	}

	return NewRegisterControl(db, local), nil
}

// TODO working
func Execute(control *RegisterControl, request entry.UserRegisterRequest, user *coreUser.UserAuthentic) (*entry.UserRegisterResponse, error) {
	return nil, nil
}
