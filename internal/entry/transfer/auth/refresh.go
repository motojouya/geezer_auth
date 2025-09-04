package auth

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
)

type RefreshTokenGetter interface {
	GetRefreshToken() (text.Token, error)
}

type AuthRefresh struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshRequest struct {
	AuthRefresh
}

func (a AuthRefreshRequest) GetRefreshToken() (text.Token, error) {
	return text.NewToken(a.RefreshToken)
}
