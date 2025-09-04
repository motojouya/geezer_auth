package user

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserGetter interface {
	ToShelterUser(pkgText.Identifier, time.Time) (shelter.UnsavedUser, error)
}

type UserRegister struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Bot      bool   `json:"bot"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	UserRegister
}

func (u UserRegisterRequest) ToShelterUser(identifier pkgText.Identifier, registerDate time.Time) (shelter.UnsavedUser, error) {
	var emailId, emailErr = pkgText.NewEmail(u.Email)
	if emailErr != nil {
		return shelter.UnsavedUser{}, emailErr
	}

	var name, nameErr = pkgText.NewName(u.Name)
	if nameErr != nil {
		return shelter.UnsavedUser{}, nameErr
	}

	return shelter.CreateUser(identifier, emailId, name, u.Bot, registerDate), nil
}

func (u UserRegisterRequest) GetPassword() (text.Password, error) {
	return text.NewPassword(u.Password)
}

func (u UserRegisterRequest) GetEmail() (pkgText.Email, error) {
	return pkgText.NewEmail(u.Email)
}
