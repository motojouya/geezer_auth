package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type CompanyUserGetter interface {
	GetUserIdentifier() (pkgText.Identifier, error)
}

type RoleAssign struct {
	RoleInvite
	UserIdentifier string `json:"user_identifier"`
}

type CompanyUserAssignRequest struct {
	company.CompanyGetRequest
	RoleAssign
}

func (c CompanyUserAssignRequest) GetRoleLabel() (pkgText.Label, error) {
	return pkgText.NewLabel(c.RoleLabel)
}

func (c CompanyUserAssignRequest) GetUserIdentifier() (pkgText.Identifier, error) {
	return pkgText.NewIdentifier(c.UserIdentifier)
}
