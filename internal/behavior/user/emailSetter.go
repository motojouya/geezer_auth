package user

import (
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	"github.com/motojouya/geezer_auth/internal/local"
)

type EmailSetterDB interface {
	userQuery.GetUserEmailQuery
	commandQuery.AddEmailQuery
}

type EmailSetter interface {
	Execute(entry entryUser.EmailGetter, userAuthentic *shelterUser.UserAuthentic) error
}

type EmailSet struct {
	local local.Local
	db    EmailSetterDB
}

func NewEmailSet(local local.Local, database EmailSetterDB) *EmailSet {
	return &EmailSet{
		db:    database,
		local: local,
	}
}

func (setter EmailSet) Execute(entry entryUser.EmailGetter, userAuthentic *shelterUser.UserAuthentic) error {
	now := setter.local.GetNow()

	email, err := entry.GetEmail()
	if err != nil {
		return err
	}

	userEmails, err := setter.db.GetUserEmail(string(email))
	if err != nil {
		return err
	}

	if len(userEmails) > 0 {
		keys := map[string]string{"email": string(email)}
		return essence.NewDuplicateError("user_email", keys, "email already exists")
	}

	verifyTokenSource, err := setter.local.GenerateUUID()
	if err != nil {
		return err
	}

	verifyToken, err := shelterText.CreateToken(verifyTokenSource)
	if err != nil {
		return err
	}

	userEmail := shelterUser.CreateUserEmail(userAuthentic.GetUser(), email, verifyToken, now)

	dbUserEmail := dbUser.FromCoreUserEmail(userEmail)

	if _, err = setter.db.AddEmail(dbUserEmail, now); err != nil {
		return err
	}

	//FIXME 未実装 ここでverify tokenを当該メールアドレスに通知する処理が入る。
	return nil
}
