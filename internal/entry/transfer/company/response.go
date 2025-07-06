package company

import (
	coreCompany "github.com/motojouya/geezer_auth/internal/core/company"
	coreUser "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/internal/utility"
)

type CompanyGetResponse struct {
	Company common.Company `json:"company"`
}

type CompanyUserResponse struct {
	Users []*common.User `json:"users"`
}

func FromCoreCompany(coreCompany coreCompany.Company) CompanyGetResponse {
	var commonUser = common.FromCoreCompany(coreCompany)
	return CompanyGetResponse{
		Company: commonUser,
	}
}

func FromCoreUserAuthentic(coreUsers []*coreUser.UserAuthentic) *CompanyUserResponse {
	var users []*common.User = utility.Map(coreUsers, common.FromCoreUser)
	return &CompanyUserResponse{
		Users: users,
	}
}
