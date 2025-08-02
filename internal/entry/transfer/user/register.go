package user

import (
	text "github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserGetter interface {
	ToCoreUser(pkgText.Identifier, time.Time) (core.UnsavedUser, error)
}

type EmailGetter interface {
	GetEmail() (pkgText.Email, error)
}

type UserRegister struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Bot      bool   `json:"bot"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	UserRegister UserRegister `http:"body"`
}

func (u UserRegisterRequest) ToCoreUser(identifier pkgText.Identifier, registerDate time.Time) (core.UnsavedUser, error) {
	var emailId, emailErr = pkgText.NewEmail(u.UserRegister.Email)
	if emailErr != nil {
		return core.UnsavedUser{}, emailErr
	}

	var name, nameErr = pkgText.NewName(u.UserRegister.Name)
	if nameErr != nil {
		return core.UnsavedUser{}, nameErr
	}

	return core.CreateUser(identifier, emailId, name, u.UserRegister.Bot, registerDate), nil
}

func (u UserRegisterRequest) GetPassword() (text.Password, error) {
	return text.NewPassword(u.UserRegister.Password)
}

func (u UserRegisterRequest) GetEmail() (pkgText.Email, error) {
	return pkgText.NewEmail(u.UserRegister.Email)
}
