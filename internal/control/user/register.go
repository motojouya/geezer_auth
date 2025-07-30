package user

import (
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/core/essence"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/motojouya/geezer_auth/internal/service"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
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
	DB    UserRegisterDB
	JWT   service.JwtHandler
}

func NewRegisterControl(local io.Local, database UserRegisterDB, jwtHandler service.JwtHandler) *RegisterControl {
	return &RegisterControl{
		DB:    database,
		Local: local,
		JWT:   jwtHandler,
	}
}

func CreateRegisterControl() (*RegisterControl, error) {
	var local = io.CreateLocal()
	var env = io.CreateEnvironment()

	var loader = service.GetLoader()

	database, err := loader.LoadDatabase(env)
	if err != nil {
		return nil, err
	}

	jwtHandler, err := loader.LoadJwtHandler(env)
	if err != nil {
		return nil, err
	}

	return NewRegisterControl(local, database, jwtHandler), nil
}

func createUserIdentifier(local io.Local) func() (pkgText.Identifier, error) {
	return func() (pkgText.Identifier, error) {
		var ramdomString = local.GenerateRamdomString(pkgText.IdentifierLength, pkgText.IdentifierChar)
		var identifier, err = coreUser.CreateUserIdentifier(ramdomString)
		if err != nil {
			return pkgText.Identifier(""), err
		}
		return identifier, nil
	}
}

func checkUserIdentifier(userRegisterDB UserRegisterDB) func(pkgText.Identifier) (bool, error) {
	return func(identifier pkgText.Identifier) (bool, error) {
		var user, err = userRegisterDB.GetUser(string(identifier))
		if err != nil {
			return false, err
		}
		return user == nil, nil
	}
}

func RegisterExecute(control *RegisterControl, entry entryUser.UserRegisterRequest, user *coreUser.UserAuthentic) (*entryUser.UserRegisterResponse, error) {
	if err := control.DB.Begin(); err != nil {
		return nil, err
	}

	var now = control.Local.GetNow()

	userIdentifier, err := coreText.GetString(createUserIdentifier(control.Local), checkUserIdentifier(control.DB), 10)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	unsavedUser, err := entry.ToCoreUser(userIdentifier, now)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	var dbUserValue = dbUser.FromCoreUser(unsavedUser)

	if err = control.DB.Insert(&dbUserValue); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	savedUser, err := dbUserValue.ToCoreUser()
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	password, err := entry.GetPassword()
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	hashedPassword, err := coreText.HashPassword(password)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	userPassword := coreUser.CreateUserPassword(savedUser, hashedPassword, now)

	dbUserPassword := dbUser.FromCoreUserPassword(userPassword)

	if err = control.DB.Insert(dbUserPassword); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	verifyTokenSource, err := control.Local.GenerateUUID()
	verifyToken, err := coreText.CreateToken(verifyTokenSource)

	userEmail := coreUser.CreateUserEmail(savedUser, savedUser.ExposeEmailId, verifyToken, now)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	dbUserEmail := dbUser.FromCoreUserEmail(userEmail)

	if _, err = control.DB.AddEmail(dbUserEmail, now); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	dbUserAuthentic, err := control.DB.GetUserAuthentic(string(savedUser.Identifier), now)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": string(savedUser.Identifier)}
		err = essence.NewNotFoundError("user", keys, "user not found")
		return nil, db.RollbackWithError(control.DB, err)
	}

	userAuthentic, err := dbUserAuthentic.ToCoreUserAuthentic()

	pkgUser := userAuthentic.ToJwtUser()

	refreshTokenSource, err := control.Local.GenerateUUID()
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	refreshToken, err := coreText.CreateToken(refreshTokenSource)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	userRefreshToken := coreUser.CreateUserRefreshToken(savedUser, refreshToken, now)
	dbUserRefreshToken := dbUser.FromCoreUserRefreshToken(userRefreshToken)

	if err := control.DB.Insert(dbUserRefreshToken); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	tokenId, err := control.Local.GenerateUUID()
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	tokenData, accessToken, err := control.JWT.Generate(pkgUser, now, tokenId.String())
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	userAccessToken := coreUser.CreateUserAccessToken(savedUser, accessToken, now, tokenData.ExpiresAt.Time)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	if err = control.DB.Insert(userAccessToken); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	response := entryUser.FromCoreUserAuthenticToRegisterResponse(userAuthentic, refreshToken, accessToken)

	if err := control.DB.Commit(); err != nil {
		return nil, err
	}

	return response, nil
}
