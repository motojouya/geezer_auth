package db

import (
	"github.com/motojouya/geezer_auth/internal/db/query/command"
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
	user.GetUserAuthenticQuery
	user.GetUserAuthenticOfCompanyQuery
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
func (orp ORPImpl) GetUserAccessToken(identifier string, now time.Time) ([]transferUser.UserAccessTokenFull, error) {
	return user.GetUserAccessToken(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserAccessToken(identifier string, now time.Time) ([]transferUser.UserAccessTokenFull, error) {
	return user.GetUserAccessToken(orp, identifier, now)
}

/** GetUserCompanyRole */
func (orp ORPImpl) GetUserCompanyRole(identifiers []string, now time.Time) ([]transferUser.UserCompanyRoleFull, error) {
	return user.GetUserCompanyRole(orp, identifiers, now)
}

func (orp ORPTransactionImpl) GetUserCompanyRole(identifiers []string, now time.Time) ([]transferUser.UserCompanyRoleFull, error) {
	return user.GetUserCompanyRole(orp, identifiers, now)
}

/** GetUserEmailOfToken */
func (orp ORPImpl) GetUserEmailOfToken(identifier string, email string, verifyToken string, now time.Time) (*transferUser.UserEmailFull, error) {
	return user.GetUserEmailOfToken(orp, identifier, email, verifyToken, now)
}

func (orp ORPTransactionImpl) GetUserEmailOfToken(identifier string, email string, verifyToken string, now time.Time) (*transferUser.UserEmailFull, error) {
	return user.GetUserEmailOfToken(orp, identifier, email, verifyToken, now)
}

/** GetUserEmail */
func (orp ORPImpl) GetUserEmail(email string) ([]transferUser.UserEmailFull, error) {
	return user.GetUserEmail(orp, email)
}

func (orp ORPTransactionImpl) GetUserEmail(email string) ([]transferUser.UserEmailFull, error) {
	return user.GetUserEmail(orp, email)
}

/** GetUserPassword */
func (orp ORPImpl) GetUserPassword(identifier string, now time.Time) (*transferUser.UserPasswordFull, error) {
	return user.GetUserPassword(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserPassword(identifier string, now time.Time) (*transferUser.UserPasswordFull, error) {
	return user.GetUserPassword(orp, identifier, now)
}

/** GetUserRefreshToken */
func (orp ORPImpl) GetUserRefreshToken(identifier string, now time.Time) (*transferUser.UserRefreshTokenFull, error) {
	return user.GetUserRefreshToken(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserRefreshToken(identifier string, now time.Time) (*transferUser.UserRefreshTokenFull, error) {
	return user.GetUserRefreshToken(orp, identifier, now)
}

/** GetUserAuthentic */
func (orp ORPImpl) GetUserAuthentic(identifier string, now time.Time) (*transferUser.UserAuthentic, error) {
	return user.GetUserAuthentic(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserAuthentic(identifier string, now time.Time) (*transferUser.UserAuthentic, error) {
	return user.GetUserAuthentic(orp, identifier, now)
}

/** GetUserAuthenticOfCompany */
func (orp ORPImpl) GetUserAuthenticOfCompany(identifier string, now time.Time) ([]transferUser.UserAuthentic, error) {
	return user.GetUserAuthenticOfCompany(orp, identifier, now)
}

func (orp ORPTransactionImpl) GetUserAuthenticOfCompany(identifier string, now time.Time) ([]transferUser.UserAuthentic, error) {
	return user.GetUserAuthenticOfCompany(orp, identifier, now)
}

type Command interface {
	command.AddEmailQuery
	command.VerifyEmailQuery
	command.AddPasswordQuery
	command.AddRefreshTokenQuery
}

func (orp ORPTransactionImpl) AddEmail(userEmail *transferUser.UserEmail, now time.Time) (*transferUser.UserEmail, error) {
	return command.AddEmail(orp, userEmail, now)
}

func (orp ORPTransactionImpl) VerifyEmail(userEmail *transferUser.UserEmail, now time.Time) (*transferUser.UserEmail, error) {
	return command.VerifyEmail(orp, userEmail, now)
}

func (orp ORPTransactionImpl) AddPassword(userPassword *transferUser.UserPassword, now time.Time) (*transferUser.UserPassword, error) {
	return command.AddPassword(orp, userPassword, now)
}

func (orp ORPTransactionImpl) AddRefreshToken(userRefreshToken *transferUser.UserRefreshToken, now time.Time) (*transferUser.UserRefreshToken, error) {
	return command.AddRefreshToken(orp, userRefreshToken, now)
}
