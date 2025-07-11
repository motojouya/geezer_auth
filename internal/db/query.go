package db

import (
	"github.com/motojouya/geezer_auth/internal/db/query/company"
	"github.com/motojouya/geezer_auth/internal/db/query/role"
	"github.com/motojouya/geezer_auth/internal/db/query/user"
	transferCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	transferRole "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	transferUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
)

type Query interface {
	company.GetCompanyQuery
	company.GetCompanyInviteQuery
	role.GetRoleQuery
	role.GetRolePermissionQuery
	user.GetUserQuery
}

/** GetCompany */
func (orp ORPImpl) GetCompany(identifier string) (*transferCompany.Company, error) {
	return company.GetCompany(orp, identifier)
}

func (orp ORPTransactionImpl) GetCompany(identifier string) (*transferCompany.Company, error) {
	return company.GetCompany(orp, identifier)
}

/** GetCompanyInvite */
func (orp ORPImpl) GetCompanyInvite(identifier string, verifyToken string) (*transferCompany.CompanyInviteFull, error) {
	return company.GetCompanyInvite(orp, identifier, verifyToken)
}

func (orp ORPTransactionImpl) GetCompanyInvite(identifier string, verifyToken string) (*transferCompany.CompanyInviteFull, error) {
	return company.GetCompanyInvite(orp, identifier, verifyToken)
}

/** GetRole */
func (orp ORPImpl) GetRole() ([]transferRole.Role, error) {
	return role.GetRole(orp)
}

func (orp ORPTransactionImpl) GetRole() ([]transferRole.Role, error) {
	return role.GetRole(orp)
}

/** GetRolePermission */
func (orp ORPImpl) GetRolePermission() ([]transferRole.RolePermission, error) {
	return role.GetRolePermission(orp)
}

func (orp ORPTransactionImpl) GetRolePermission() ([]transferRole.RolePermission, error) {
	return role.GetRolePermission(orp)
}

/** GetUser */
func (orp ORPImpl) GetUser(identifier string) (*transferUser.User, error) {
	return user.GetUser(orp, identifier)
}

func (orp ORPTransactionImpl) GetUser(identifier string) (*transferUser.User, error) {
	return user.GetUser(orp, identifier)
}
