package auth

import (
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	userBehavior "github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/control/utility"
	"github.com/motojouya/geezer_auth/internal/db"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type RefreshAccessTokenControl struct {
	db.TransactionalDatabase
	refreshTokenChecker userBehavior.RefreshTokenChecker
	accessTokenIssuer   userBehavior.AccessTokenIssuer
}

func NewRefreshAccessTokenControl(
	database db.TransactionalDatabase,
	refreshTokenChecker userBehavior.RefreshTokenChecker,
	accessTokenIssuer userBehavior.AccessTokenIssuer,
) *RefreshAccessTokenControl {
	return &RefreshAccessTokenControl{
		TransactionalDatabase: database,
		refreshTokenChecker:   refreshTokenChecker,
		accessTokenIssuer:     accessTokenIssuer,
	}
}

func CreateRefreshAccessTokenControl() (*RefreshAccessTokenControl, error) {
	var local = localPkg.CreateLocal()
	var env = localPkg.CreateEnvironment()

	database, err := configBehavior.NewDatabaseGet(env).GetDatabase()
	if err != nil {
		return nil, err
	}

	jwtHandler, err := configBehavior.NewJwtHandlerGet(env).GetJwtHandler()
	if err != nil {
		return nil, err
	}

	refreshTokenChecker := userBehavior.NewRefreshTokenCheck(local, database)
	accessTokenIssuer := userBehavior.NewAccessTokenIssue(local, database, jwtHandler)

	return NewRefreshAccessTokenControl(database, refreshTokenChecker, accessTokenIssuer), nil
}

var RefreshAccessTokenControlExecute = utility.Transact(func(control *RefreshAccessTokenControl, entry entryAuth.AuthRefreshRequest, _ *pkgUser.Authentic) (*entryAuth.AuthRefreshResponse, error) {

	userAuthentic, err := control.refreshTokenChecker.Execute(entry)
	if err != nil {
		return nil, err
	}

	accessToken, err := control.accessTokenIssuer.Execute(userAuthentic)
	if err != nil {
		return nil, err
	}

	return entryAuth.FromShelterUserAuthenticToRefreshResponse(userAuthentic, accessToken), nil
})
