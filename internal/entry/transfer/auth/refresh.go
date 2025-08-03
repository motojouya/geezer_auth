package auth

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type AuthRefresh struct {
	AuthIdentifier
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshRequest struct {
	AuthRefresh AuthRefresh `http:"body"`
}

func (a AuthRefreshRequest) GetIdentifier() (*pkgText.Identifier, error) {
	return a.AuthRefresh.GetIdentifier()
}

func (a AuthRefreshRequest) GetEmailIdentifier() (*pkgText.Email, error) {
	return a.AuthRefresh.GetEmailIdentifier()
}

func (a AuthRefreshRequest) GetRefreshToken() (text.Token, error) {
	return text.NewToken(a.AuthRefresh.RefreshToken)
}
