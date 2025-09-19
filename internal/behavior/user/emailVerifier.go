package user

import (
	"github.com/go-gorp/gorp"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgEssence "github.com/motojouya/geezer_auth/pkg/shelter/essence"
)

type EmailVerifierDB interface {
	gorp.SqlExecutor
	userQuery.GetUserEmailOfTokenQuery
	commandQuery.VerifyEmailQuery
	userQuery.GetUserAuthenticQuery
}

type EmailVerifier interface {
	Execute(entry entryUser.EmailVerifier, userAuthentic *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error)
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

	dbUserEmailFull, err := verifier.db.GetUserEmailOfToken(string(userAuthentic.Identifier), string(email))
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

	userEmail.VerifyDate = &now

	dbUserEmail := dbUser.FromShelterUserEmail(userEmail)

	if _, err = verifier.db.VerifyEmail(dbUserEmail, now); err != nil {
		return nil, err
	}

	updatedUser := userAuthentic.GetUser().Update(now)
	var dbUserValue = dbUser.FromShelterUser(updatedUser)

	_, err = verifier.db.Update(&dbUserValue)
	if err != nil {
		return nil, err
	}

	dbUserAuthentic, err := verifier.db.GetUserAuthentic(string(userAuthentic.Identifier), now)
	if err != nil {
		return nil, err
	}

	if dbUserAuthentic == nil {
		keys := map[string]string{"identifier": string(userAuthentic.Identifier)}
		return nil, essence.NewNotFoundError("user", keys, "user not found")
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
