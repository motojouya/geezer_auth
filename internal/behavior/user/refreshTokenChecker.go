package user

import (
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type RefreshTokenIssuerDB interface {
	commandQuery.AddRefreshTokenQuery
}

type RefreshTokenIssuer interface {
	Execute(userAuthentic *shelterUser.UserAuthentic) (shelterText.Token, error)
}

type RefreshTokenIssue struct {
	local localPkg.Localer
	db    RefreshTokenIssuerDB
}

func NewRefreshTokenIssue(local localPkg.Localer, database RefreshTokenIssuerDB) *RefreshTokenIssue {
	return &RefreshTokenIssue{
		db:    database,
		local: local,
	}
}

// TODO refresh tokenは単体で認証もできるので漏れるとだいぶまずい。passwordと同様に不可逆暗号として保存すべき
func (issuer RefreshTokenIssue) Execute(entry entryAuth.RefreshTokenGetter) (userAuthentic *shelterUser.UserAuthentic, error) {
	now := issuer.local.GetNow()

	refreshTokenSource, err := issuer.local.GenerateUUID()
	if err != nil {
		return shelterText.Token(""), err
	}

	refreshToken, err := shelterText.CreateToken(refreshTokenSource)
	if err != nil {
		return shelterText.Token(""), err
	}

	userRefreshToken := shelterUser.CreateUserRefreshToken(userAuthentic.GetUser(), refreshToken, now)
	dbUserRefreshToken := dbUser.FromShelterUserRefreshToken(userRefreshToken)

	_, err = issuer.db.AddRefreshToken(dbUserRefreshToken, now)
	if err != nil {
		return shelterText.Token(""), err
	}

	return refreshToken, nil
}
