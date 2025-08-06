package user

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserVerifyEmail struct {
	Email       string `json:"email"`
	VerifyToken string `json:"verify_token"`
}

type UserVerifyEmailRequest struct {
	UserVerifyEmail UserVerifyEmail `http:"body"`
}

func (u UserVerifyEmailRequest) ToShelterUserEmail(user shelter.User, registerDate time.Time) (*shelter.UnsavedUserEmail, error) {
	var email, emailErr = pkgText.NewEmail(u.UserVerifyEmail.Email)
	if emailErr != nil {
		return &shelter.UnsavedUserEmail{}, emailErr
	}

	var verifyToken, tokenErr = text.NewToken(u.UserVerifyEmail.VerifyToken)
	if tokenErr != nil {
		return &shelter.UnsavedUserEmail{}, tokenErr
	}

	return shelter.CreateUserEmail(user, email, verifyToken, registerDate), nil
}
