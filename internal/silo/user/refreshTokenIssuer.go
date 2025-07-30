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
	"time"
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

func createUser(control *RegisterControl, entry entryUser.UserRegisterRequest, now time.Time) (coreUser.User, error) {
	userIdentifier, err := coreText.GetString(createUserIdentifier(control.Local), checkUserIdentifier(control.DB), 10)
	if err != nil {
		return coreUser.User{}, err
	}

	unsavedUser, err := entry.ToCoreUser(userIdentifier, now)
	if err != nil {
		return coreUser.User{}, err
	}

	var dbUserValue = dbUser.FromCoreUser(unsavedUser)

	if err = control.DB.Insert(&dbUserValue); err != nil {
		return coreUser.User{}, err
	}

	return dbUserValue.ToCoreUser()
}

func setPassword(control *RegisterControl, entry entryUser.UserRegisterRequest, now time.Time, savedUser coreUser.User) error {
	password, err := entry.GetPassword()
	if err != nil {
		return err
	}

	hashedPassword, err := coreText.HashPassword(password)
	if err != nil {
		return err
	}

	userPassword := coreUser.CreateUserPassword(savedUser, hashedPassword, now)

	dbUserPassword := dbUser.FromCoreUserPassword(userPassword)

	if err = control.DB.Insert(dbUserPassword); err != nil {
		return err
	}

	return nil
}

func setEmail(control *RegisterControl, now time.Time, savedUser coreUser.User) error {
	verifyTokenSource, err := control.Local.GenerateUUID()
	if err != nil {
		return err
	}

	verifyToken, err := coreText.CreateToken(verifyTokenSource)
	if err != nil {
		return err
	}

	userEmail := coreUser.CreateUserEmail(savedUser, savedUser.ExposeEmailId, verifyToken, now)

	dbUserEmail := dbUser.FromCoreUserEmail(userEmail)

	if _, err = control.DB.AddEmail(dbUserEmail, now); err != nil {
		return err
	}

	//FIXME 未実装 ここでverify tokenを当該メールアドレスに通知する処理が入る。
	return nil
}

func getUserAuthentic(control *RegisterControl, now time.Time, savedUser coreUser.User) (*coreUser.UserAuthentic, error) {
	dbUserAuthentic, err := control.DB.GetUserAuthentic(string(savedUser.Identifier), now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": string(savedUser.Identifier)}
		err = essence.NewNotFoundError("user", keys, "user not found")
		return nil, err
	}

	return dbUserAuthentic.ToCoreUserAuthentic()
}

func createRefreshToken(control *RegisterControl, now time.Time, savedUser coreUser.User) (coreText.Token, error) {
	refreshTokenSource, err := control.Local.GenerateUUID()
	if err != nil {
		return coreText.Token(""), err
	}

	refreshToken, err := coreText.CreateToken(refreshTokenSource)
	if err != nil {
		return coreText.Token(""), err
	}

	userRefreshToken := coreUser.CreateUserRefreshToken(savedUser, refreshToken, now)
	dbUserRefreshToken := dbUser.FromCoreUserRefreshToken(userRefreshToken)

	if err := control.DB.Insert(dbUserRefreshToken); err != nil {
		return coreText.Token(""), err
	}

	return refreshToken, nil
}

func createAccessToken(control *RegisterControl, now time.Time, savedUser coreUser.User, userAuthentic *coreUser.UserAuthentic) (pkgText.JwtToken, error) {
	tokenId, err := control.Local.GenerateUUID()
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	pkgUser := userAuthentic.ToJwtUser()

	tokenData, accessToken, err := control.JWT.Generate(pkgUser, now, tokenId.String())
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	userAccessToken := coreUser.CreateUserAccessToken(savedUser, accessToken, now, tokenData.ExpiresAt.Time)
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	if err = control.DB.Insert(userAccessToken); err != nil {
		return pkgText.JwtToken(""), err
	}

	return accessToken, nil
}

func RegisterExecute(control *RegisterControl, entry entryUser.UserRegisterRequest, user *coreUser.UserAuthentic) (*entryUser.UserRegisterResponse, error) {
	if err := control.DB.Begin(); err != nil {
		return nil, err
	}

	now := control.Local.GetNow()

	savedUser, err := createUser(control, entry, now)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	if err = setPassword(control, entry, now, savedUser); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	if err = setEmail(control, now, savedUser); err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	userAuthentic, err := getUserAuthentic(control, now, savedUser)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	refreshToken, err := createRefreshToken(control, now, savedUser)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	accessToken, err := createAccessToken(control, now, savedUser, userAuthentic)
	if err != nil {
		return nil, db.RollbackWithError(control.DB, err)
	}

	response := entryUser.FromCoreUserAuthenticToRegisterResponse(userAuthentic, refreshToken, accessToken)

	if err := control.DB.Commit(); err != nil {
		return nil, err
	}

	return response, nil
}
