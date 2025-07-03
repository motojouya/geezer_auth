package auth

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
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

func FromCoreUserAuthenticToRefreshResponse(coreUser *core.UserAuthentic, accessToken pkgText.JwtToken) *AuthRefreshResponse {
	var userGetResponse = user.FromCoreUserAuthenticToGetResponse(coreUser)
	return &AuthRefreshResponse{
		UserUpdateResponse: user.UserUpdateResponse{
			UserGetResponse: *userGetResponse,
			AccessToken: string(accessToken),
		},
	}
}

func FromCoreUserAuthenticToLoginResponse(coreUser *core.UserAuthentic, refreshToken text.Token, accessToken pkgText.JwtToken) *AuthLoginResponse {
	var userUpdateResponse = user.FromCoreUserAuthenticToUpdateResponse(coreUser, accessToken)
	return &AuthLoginResponse{
		UserRegisterResponse: user.UserRegisterResponse{
			UserUpdateResponse: *userUpdateResponse,
			RefreshToken: string(refreshToken),
		},
	}
}
