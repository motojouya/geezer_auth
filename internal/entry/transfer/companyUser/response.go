package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	coreRole "github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/utility"
)

type CompanyUserResponse struct {
	user.UserGetResponse
}

type CompanyUserInviteResponse struct {
	Token string `json:"token"`
}

type RoleGetResponse struct {
	Roles []common.Role `json:"roles"`
}

func FromCoreRoles(coreRoles []coreRole.Role) RoleGetResponse {
	var roles = utility.Map(coreRoles, common.FromCoreRole)
	return RoleGetResponse{
		Roles: roles,
	}
}

func FromCoreUserAuthenticToGetResponse(coreUser *coreUser.UserAuthentic) *CompanyUserResponse {
	var commonUser = common.FromCoreUser(coreUser)
	return &CompanyUserResponse{
		UserGetResponse: user.UserGetResponse{
			User: *commonUser,
		},
	}
}

func FromToken(token text.Token) CompanyUserInviteResponse {
	return CompanyUserInviteResponse{
		Token: string(token),
	}
}
