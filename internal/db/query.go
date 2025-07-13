package db

import (
	"github.com/motojouya/geezer_auth/internal/db/query/company"
	"github.com/motojouya/geezer_auth/internal/db/query/role"
	"github.com/motojouya/geezer_auth/internal/db/query/user"
	transferCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	transferRole "github.com/motojouya/geezer_auth/internal/db/transfer/role"
	transferUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"time"
)

type Query interface {
	company.GetCompanyQuery
	company.GetCompanyInviteQuery
	role.GetRoleQuery
	role.GetRolePermissionQuery
	user.GetUserQuery
	user.GetUserAccessTokenQuery
	user.GetUserCompanyRoleQuery
	user.GetUserEmailOfTokenQuery
	user.GetUserEmailQuery
	user.GetUserPasswordQuery
	user.GetUserRefreshTokenQuery
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

/** GetUserAccessToken */
func (orp ORPImpl) GetUserAccessToken(identifier string, now time.Time) ([]transferUser.UserAccessToken, error) {
	return user.GetUserAccessToken(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserAccessToken(identifier string, now time.Time) ([]transferUser.UserAccessToken, error) {
	return user.GetUserAccessToken(orp, identifier, now)
}

/** GetUserCompanyRole */
func (orp ORPImpl) GetUserCompanyRole(identifier string, now time.Time) ([]transferUser.UserCompanyRole, error) {
	return user.GetUserCompanyRole(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserCompanyRole(identifier string, now time.Time) ([]transferUser.UserCompanyRole, error) {
	return user.GetUserCompanyRole(orp, identifier, now)
}

/** GetUserEmailOfToken */
func (orp ORPImpl) GetUserEmailOfToken(identifier string, email string, verifyToken string, now time.Time) (*transferUser.UserEmail, error) {
	return user.GetUserEmailOfToken(orp, identifier, email, verifyToken, now)
}

func (orp ORPTransactionImpl) GetUserEmailOfToken(identifier string, email string, verifyToken string, now time.Time) (*transferUser.UserEmail, error) {
	return user.GetUserEmailOfToken(orp, identifier, email, verifyToken, now)
}

/** GetUserEmail */
func (orp ORPImpl) GetUserEmail(email string) ([]transferUser.UserEmail, error) {
	return user.GetUserEmail(orp, email)
}

func (orp ORPTransactionImpl) GetUserEmail(email string) ([]transferUser.UserEmail, error) {
	return user.GetUserEmail(orp, email)
}

/** GetUserPassword */
func (orp ORPImpl) GetUserPassword(identifier string, now time.Time) (*transferUser.UserPassword, error) {
	return user.GetUserPassword(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserPassword(identifier string, now time.Time) (*transferUser.UserPassword, error) {
	return user.GetUserPassword(orp, identifier, now)
}

/** GetUserRefreshToken */
func (orp ORPImpl) GetUserRefreshToken(identifier string, now time.Time) (*transferUser.UserRefreshToken, error) {
	return user.GetUserRefreshToken(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserRefreshToken(identifier string, now time.Time) (*transferUser.UserRefreshToken, error) {
	return user.GetUserRefreshToken(orp, identifier, now)
}
