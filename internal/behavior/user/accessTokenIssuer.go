package user

import (
	"github.com/go-gorp/gorp"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/motojouya/geezer_auth/pkg/shelter/jwt"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type AccessTokenIssuerDB interface {
	gorp.SqlExecutor
	db.Transactional
	userQuery.GetUserAccessTokenQuery
}

type AccessTokenIssuer interface {
	Execute(userAuthentic *shelterUser.UserAuthentic) (pkgText.JwtToken, error)
}

type AccessTokenIssue struct {
	local io.Local
	db    AccessTokenIssuerDB
	jwt   jwt.JwtHandler
}

func NewAccessTokenIssue(local io.Local, database AccessTokenIssuerDB, jwtHandler jwt.JwtHandler) *AccessTokenIssue {
	return &AccessTokenIssue{
		local: local,
		db:    database,
		jwt:   jwtHandler,
	}
}

func (issuer AccessTokenIssue) Execute(userAuthentic *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
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

	userAccessToken := shelterUser.CreateUserAccessToken(userAuthentic.GetUser(), accessToken, now, tokenData.ExpiresAt.Time)
	if err != nil {
		return pkgText.JwtToken(""), err
	}

	if err = issuer.db.Insert(userAccessToken); err != nil {
		return pkgText.JwtToken(""), err
	}

	return accessToken, nil
}
