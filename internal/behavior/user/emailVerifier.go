package user

import (
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	pkgEssence "github.com/motojouya/geezer_auth/pkg/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type EmailVerifierDB interface {
	userQuery.GetUserEmailOfTokenQuery
	commandQuery.VerifyEmailQuery
	userQuery.GetUserAuthenticQuery
}

type EmailVerifier interface {
	Execute(entry entryUser.EmailGetter, userAuthentic *shelterUser.UserAuthentic) error
}

type EmailVerify struct {
	local localPkg.Localer
	db    EmailVerifierDB
}

func NewEmailVerify(local localPkg.Localer, database EmailVerifierDB) *EmailVerify {
	return &EmailVerify{
		db:    database,
		local: local,
	}
}

func (verifier EmailVerify) Execute(entry entryUser.EmailVerifier, userAuthentic *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
	if userAuthentic == nil {
		return nil, pkgEssence.NewNilError("userAuthentic", "user authentic is nil")
	}

	now := verifier.local.GetNow()

	email, err := entry.GetEmail()
	if err != nil {
		return nil, err
	}

	verifyToken, err := entry.GetVerifyToken()
	if err != nil {
		return nil, err
	}

	userIdentifier := string(userAuthentic.Identifier)
	dbUserEmailFull, err := verifier.db.GetUserEmailOfToken(userIdentifier, string(email))
	if err != nil {
		return nil, err
	}

	if dbUserEmailFull == nil {
		keys := map[string]string{"identifier": string(userAuthentic.Identifier), "email": string(email)}
		return nil, essence.NewNotFoundError("userEmail", keys, "user email not found")
	}

	userEmail, err := dbUserEmailFull.ToShelterUserEmail()
	if err != nil {
		return nil, err
	}

	if userEmail.VerifyToken != verifyToken {
		return nil, essence.NewInvalidArgumentError("verifyToken", string(verifyToken), "verify token is invalid")
	}

	dbUserEmail := dbUser.FromShelterUserEmail(userEmail)

	if _, err = verifier.db.VerifyEmail(dbUserEmail, now); err != nil {
		return nil, err
	}

	dbUserAuthentic, err := verifier.db.GetUserAuthentic(userIdentifier, now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": userIdentifier}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
