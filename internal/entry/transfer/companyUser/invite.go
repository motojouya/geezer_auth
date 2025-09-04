package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type RoleGetter interface {
	GetRoleLabel() (pkgText.Label, error)
}

type RoleInvite struct {
	RoleLabel string `json:"role_label"`
}

type CompanyUserInviteRequest struct {
	company.CompanyGetRequest
	RoleInvite
}

func (c CompanyUserInviteRequest) GetRoleLabel() (pkgText.Label, error) {
	return pkgText.NewLabel(c.RoleLabel)
}
