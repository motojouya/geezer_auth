package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type RoleAssign struct {
	RoleInvite
	UserIdentifier string `json:"user_identifier"`
}

type CompanyUserAssignRequest struct {
	company.CompanyGetRequest
	RoleAssign RoleAssign `http:"body"`
}

func (c CompanyUserAssignRequest) GetRoleLabel() (pkgText.Label, error) {
	return pkgText.NewLabel(c.RoleAssign.RoleLabel)
}

func (c CompanyUserAssignRequest) GetUserIdentifier() (pkgText.Identifier, error) {
	return pkgText.NewIdentifier(c.RoleAssign.UserIdentifier)
}
