package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
)

type RoleInvite struct {
	RoleLabel string `json:"role_label"`
}

type CompanyUserInviteRequest struct {
	company.CompanyGetRequest
	RoleInvite RoleInvite `http:"body"`
}

func (c CompanyUserInviteRequest) GetRoleLabel() (pkgText.Label, error) {
	return pkgText.NewLabel(c.RoleInvite.RoleLabel)
}
