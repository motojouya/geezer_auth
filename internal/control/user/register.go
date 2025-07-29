package user

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/service"
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/pkg/core/text"
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

func NewRegisterControl(database UserRegisterDB, local io.Local) *RegisterControl {
	return &RegisterControl{
		DB:    database,
		Local: local,
	}
}

func CreateRegisterControl() (*RegisterControl, error) {
	var local = io.CreateLocal()
	var env = io.CreateEnvironment()

	var loader = service.GetLoader()
	var database, err = loader.LoadDatabase(env)
	if err != nil {
		return nil, err
	}

	return NewRegisterControl(database, local), nil
}

// TODO working
func RegisterExecute(control *RegisterControl, entryUser entry.UserRegisterRequest, user *coreUser.UserAuthentic) (*entry.UserRegisterResponse, error) {
	if err := control.DB.Begin(); err != nil {
		return nil, err
	}
	// db.RollbackWithError(control.DB, err)

	var now = control.Local.GetNow()

	// TODO db checkを何回かまでretry
	var ramdomString = control.Local.GenerateRamdomString(text.IdentifierLength, text.IdentifierChar)
	var identifier, identifierErr = coreUser.CreateUserIdentifier(ramdomString)
	if identifierErr != nil {
		return nil, db.RollbackWithError(control.DB, identifierErr)
	}

	var user, getUserErr = control.DB.GetUser(identifier)
	if getUserErr != nil {
		return nil, db.RollbackWithError(control.DB, getUserErr)
	}

	var unsavedUser, userErr = entryUser.ToCoreUser(user.Identifier, now)

	var password, passwordErr = entryUser.GetPassword()

















	if err := control.DB.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
