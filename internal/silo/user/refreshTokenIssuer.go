package user

import (
	"github.com/go-gorp/gorp"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/db"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
)

type RefreshTokenIssuerDB interface {
	gorp.SqlExecutor
	db.Transactional
	commandQuery.AddRefreshTokenQuery
}

type RefreshTokenIssuer interface {
	Execute(userAuthentic *coreUser.UserAuthentic) (coreText.Token, error)
}

type RefreshTokenIssue struct {
	local io.Local
	db    RefreshTokenIssuerDB
}

func NewRefreshTokenIssue(local io.Local, database RefreshTokenIssuerDB) *RefreshTokenIssue {
	return &RefreshTokenIssue{
		db:    database,
		local: local,
	}
}

func (issuer RefreshTokenIssue) Execute(userAuthentic *coreUser.UserAuthentic) (coreText.Token, error) {
	now := issuer.local.GetNow()

	refreshTokenSource, err := issuer.local.GenerateUUID()
	if err != nil {
		return coreText.Token(""), err
	}

	refreshToken, err := coreText.CreateToken(refreshTokenSource)
	if err != nil {
		return coreText.Token(""), err
	}

	userRefreshToken := coreUser.CreateUserRefreshToken(userAuthentic.GetUser(), refreshToken, now)
	dbUserRefreshToken := dbUser.FromCoreUserRefreshToken(userRefreshToken)

	if err := issuer.db.Insert(dbUserRefreshToken); err != nil {
		return coreText.Token(""), err
	}

	return refreshToken, nil
}
