package db

import (
	"github.com/motojouya/geezer_auth/internal/db/query"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/transfer/role"
)

type Query interface {
	query.GetCompanyQuery
	query.GetCompanyInviteQuery
	query.GetRoleQuery
	query.GetRolePermissionQuery
}

/** GetCompany */
func (orp ORPImpl) GetCompany(identifier string) (*company.Company, error) {
	return query.GetCompany(orp, identifier)
}

func (orp ORPTransactionImpl) GetCompany(identifier string) (*company.Company, error) {
	return query.GetCompany(orp, identifier)
}

/** GetCompanyInvite */
func (orp ORPImpl) GetCompanyInvite(identifier string, verifyToken string) (*company.CompanyInviteFull, error) {
	return query.GetCompanyInvite(orp, identifier, verifyToken)
}

func (orp ORPTransactionImpl) GetCompanyInvite(identifier string, verifyToken string) (*company.CompanyInviteFull, error) {
	return query.GetCompanyInvite(orp, identifier, verifyToken)
}

/** GetRole */
func (orp ORPImpl) GetRole() ([]role.Role, error) {
	return query.GetRole(orp)
}

func (orp ORPTransactionImpl) GetRole() ([]role.Role, error) {
	return query.GetRole(orp)
}

/** GetRolePermission */
func (orp ORPImpl) GetRolePermission() ([]role.RolePermission, error) {
	return query.GetRolePermission(orp)
}

func (orp ORPTransactionImpl) GetRolePermission() ([]role.RolePermission, error) {
	return query.GetRolePermission(orp)
}
