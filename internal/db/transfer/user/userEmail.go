package user

import (
	text "github.com/motojouya/geezer_auth/internal/core/text"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
)

type UserEmail struct {
	PersistKey         uint
	UserPersistKey     uint
	Email              string
	VerifyToken        string
	RegisterDate       time.Time
	VerifyDate         *time.Time
	ExpireDate         *time.Time
}

type UserEmailFull struct {
	UserEmail
	UserIdentifier     string
	UserExposeEmailId  string
	UserName           string
	UserBotFlag        bool
	UserRegisteredDate time.Time
	UserUpdateDate     time.Time
}

func FromCoreUserEmail(coreUserEmail *core.UnsavedUserEmail) *UserEmail {
	return &UserEmail{
		UserPersistKey: coreUserEmail.User.PersistKey,
		Email:          string(coreUserEmail.Email),
		VerifyToken:    string(coreUserEmail.VerifyToken),
		RegisterDate:   coreUserEmail.RegisterDate,
		VerifyDate:     coreUserEmail.VerifyDate,
		ExpireDate:     coreUserEmail.ExpireDate,
	}
}

func (u UserEmailFull) ToCoreUserEmail() (*core.UserEmail, error) {
	var user, userErr = (User{
		PersistKey:     u.UserPersistKey,
		Identifier:     u.UserIdentifier,
		ExposeEmailId:  u.UserExposeEmailId,
		Name:           u.UserName,
		BotFlag:        u.UserBotFlag,
		RegisteredDate: u.UserRegisteredDate,
		UpdateDate:     u.UserUpdateDate,
	}).ToCoreUser()
	if userErr != nil {
		return &core.UserEmail{}, userErr
	}

	var email, emailErr = pkg.NewEmail(u.Email)
	if emailErr != nil {
		return &core.UserEmail{}, emailErr
	}

	var verifyToken, tokenErr = text.NewToken(u.VerifyToken)
	if tokenErr != nil {
		return &core.UserEmail{}, tokenErr
	}

	return core.NewUserEmail(
		u.PersistKey,
		user,
		email,
		verifyToken,
		u.RegisterDate,
		u.VerifyDate,
		u.ExpireDate,
	), nil
}
