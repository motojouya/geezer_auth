package auth

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type IdentifierGetter interface {
	GetIdentifier() (*pkgText.Identifier, error)
	GetEmailIdentifier() (*pkgText.Email, error)
}

type AuthIdentifier struct {
	Identifier      *string `json:"identifier"`
	EmailIdentifier *string `json:"email_identifier"`
}

func (a AuthIdentifier) GetIdentifier() (*pkgText.Identifier, error) {
	if a.Identifier == nil {
		return nil, nil
	}
	var identifier, err = pkgText.NewIdentifier(*a.Identifier)
	if err != nil {
		return nil, err
	}

	return &identifier, nil
}

func (a AuthIdentifier) GetEmailIdentifier() (*pkgText.Email, error) {
	if a.EmailIdentifier == nil {
		return nil, nil
	}
	var email, err = pkgText.NewEmail(*a.EmailIdentifier)
	if err != nil {
		return nil, err
	}

	return &email, nil
}

type AuthLoginner interface {
	IdentifierGetter
	GetPassword() (text.Password, error)
}

type AuthLogin struct {
	AuthIdentifier
	Password string `json:"password"`
}

type AuthLoginRequest struct {
	AuthLogin
}

func (a AuthLoginRequest) GetPassword() (text.Password, error) {
	return text.NewPassword(a.Password)
}
