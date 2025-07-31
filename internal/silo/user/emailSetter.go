package user

import (
	"github.com/motojouya/geezer_auth/internal/core/essence"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	commandQuery "github.com/motojouya/geezer_auth/internal/db/query/command"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/io"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
)

type EmailSetterDB interface {
	userQuery.GetUserEmailQuery
	commandQuery.AddEmailQuery
}

type EmailSetter struct {
	local io.Local
	db    EmailSetterDB
}

func NewEmailSetter(local io.Local, database EmailSetterDB) *EmailSetter {
	return &EmailSetter{
		db:    database,
		local: local,
	}
}

type EmailGetter interface {
	GetEmail() (pkgText.Email, error)
}

func (setter EmailSetter) Execute(entry EmailGetter, userAuthentic *coreUser.UserAuthentic) error {
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

	verifyToken, err := coreText.CreateToken(verifyTokenSource)
	if err != nil {
		return err
	}

	userEmail := coreUser.CreateUserEmail(userAuthentic.GetUser(), email, verifyToken, now)

	dbUserEmail := dbUser.FromCoreUserEmail(userEmail)

	if _, err = setter.db.AddEmail(dbUserEmail, now); err != nil {
		return err
	}

	//FIXME 未実装 ここでverify tokenを当該メールアドレスに通知する処理が入る。
	return nil
}
