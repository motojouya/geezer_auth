package user

import (
	"github.com/go-gorp/gorp"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	"github.com/motojouya/geezer_auth/internal/io"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/silo/config"
)

type AccessTokenIssuerDB interface {
	gorp.SqlExecutor
	db.Transactional
	userQuery.GetUserAccessTokenQuery
}

type AccessTokenIssuer struct {
	local io.Local
	db    AccessTokenIssuerDB
	jwt   config.JwtHandler
}

func NewAccessTokenIssuer(local io.Local, database AccessTokenIssuerDB, jwtHandler config.JwtHandler) *AccessTokenIssuer {
	return &AccessTokenIssuer{
		local: local,
		db:    database,
		jwt:   jwtHandler,
	}
}

func (issuer AccessTokenIssuer) Execute(userAuthentic *coreUser.UserAuthentic) (pkgText.JwtToken, error) {
	now := issuer.local.GetNow()

	dbAccessTokens, err := issuer.db.GetUserAccessToken(string(userAuthentic.Identifier), now)
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	// token期限切れ間近の場合に再発行したいという要件がありそうなので、2つまで発行する
	if len(dbAccessTokens) > 1 {
		userAccessToken, err := dbAccessTokens[0].ToCoreUserAccessToken()
		if err != nil {
			return pkgText.JwtToken(""), err
		}

		return userAccessToken.AccessToken, nil
	}

	tokenId, err := issuer.local.GenerateUUID()
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	pkgUser := userAuthentic.ToJwtUser()

	tokenData, accessToken, err := issuer.jwt.Generate(pkgUser, now, tokenId.String())
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	userAccessToken := coreUser.CreateUserAccessToken(userAuthentic.GetUser(), accessToken, now, tokenData.ExpiresAt.Time)
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	if err = issuer.db.Insert(userAccessToken); err != nil {
		return pkgText.JwtToken(""), err
	}

	return accessToken, nil
}
