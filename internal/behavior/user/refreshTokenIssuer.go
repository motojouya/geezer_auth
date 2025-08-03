package user

import (
	"github.com/go-gorp/gorp"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	"github.com/motojouya/geezer_auth/internal/db"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
)

type RefreshTokenIssuerDB interface {
	gorp.SqlExecutor
	db.Transactional
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

func (issuer RefreshTokenIssue) Execute(userAuthentic *shelterUser.UserAuthentic) (shelterText.Token, error) {
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
	dbUserRefreshToken := dbUser.FromCoreUserRefreshToken(userRefreshToken)

	if err := issuer.db.Insert(dbUserRefreshToken); err != nil {
		return shelterText.Token(""), err
	}

	return refreshToken, nil
}
