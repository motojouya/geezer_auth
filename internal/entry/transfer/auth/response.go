package auth

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	text "github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
)

type AuthRefreshResponse struct {
	user.UserUpdateResponse
}

type AuthLoginResponse struct {
	user.UserRegisterResponse
}

type UserUpdateResponse struct {
	UserGetResponse
	AccessToken string `json:"access_token"`
}

type UserRegisterResponse struct {
	UserUpdateResponse
	RefreshToken string `json:"refresh_token"`
}

func FromCoreUserAuthenticToRefreshResponse(coreUser *core.UserAuthentic, accessToken pkgText.JwtToken) *AuthRefreshResponse {
	var userGetResponse = user.FromCoreUserAuthenticToGetResponse(coreUser)
	return &AuthRefreshResponse{
		UserUpdateResponse: UserUpdateResponse{
			UserGetResponse: *userGetResponse,
			AccessToken: string(accessToken),
		},
	}
}

func FromCoreUserAuthenticToLoginResponse(coreUser *core.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *AuthLoginResponse {
	var userUpdateResponse = FromCoreUserAuthenticToRefreshResponse(coreUser, accessToken)
	return &AuthLoginResponse{
		UserRegisterResponse: UserRegisterResponse{
			UserUpdateResponse: *userUpdateResponse,
			RefreshToken: string(refreshToken),
		},
	}
}
