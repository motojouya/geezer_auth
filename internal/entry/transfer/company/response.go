package company

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	text "github.com/motojouya/geezer_auth/internal/core/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/internal/utility"
)

type CompanyGetResponse struct {
	Company common.Company `json:"company"`
}

type CompanyUserResponse struct {
	Users []*common.User `json:"users"`
}

func FromCoreCompany(coreCompany core.Company) CompanyGetResponse {
	var commonUser = common.FromCoreCompany(coreCompany)
	return CompanyGetResponse{
		User: commonUser,
	}
}

func FromCoreUserAuthentic(coreUsers []*core.UserAuthentic) *CompanyUserResponse {
	var users []*common.User = utility.Map(coreUsers, common.FromCoreUser)
	return &CompanyUserResponse{
		Users: users,
	}
}
