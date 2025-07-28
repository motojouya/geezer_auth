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
	command.AddEmailQuery
	command.VerifyEmailQuery
	command.AddPasswordQuery
	command.AddRefreshTokenQuery
}

func (orp ORPImpl) GetCompany(identifier string) (*transferCompany.Company, error) {
	return company.GetCompany(orp, identifier)
}

func (orp ORPImpl) GetCompanyInvite(identifier string, verifyToken string) (*transferCompany.CompanyInviteFull, error) {
	return company.GetCompanyInvite(orp, identifier, verifyToken)
}

func (orp ORPImpl) GetRole() ([]transferRole.Role, error) {
	return role.GetRole(orp)
}

func (orp ORPImpl) GetRolePermission() ([]transferRole.RolePermission, error) {
	return role.GetRolePermission(orp)
}

func (orp ORPImpl) GetUser(identifier string) (*transferUser.User, error) {
	return user.GetUser(orp, identifier)
}

func (orp ORPImpl) GetUserAccessToken(identifier string, now time.Time) ([]transferUser.UserAccessTokenFull, error) {
	return user.GetUserAccessToken(orp, identifier, now)
}

func (orp ORPImpl) GetUserCompanyRole(identifiers []string, now time.Time) ([]transferUser.UserCompanyRoleFull, error) {
	return user.GetUserCompanyRole(orp, identifiers, now)
}

func (orp ORPImpl) GetUserEmailOfToken(identifier string, email string) (*transferUser.UserEmailFull, error) {
	return user.GetUserEmailOfToken(orp, identifier, email)
}

func (orp ORPImpl) GetUserEmail(email string) ([]transferUser.UserEmailFull, error) {
	return user.GetUserEmail(orp, email)
}

func (orp ORPImpl) GetUserPassword(identifier string) (*transferUser.UserPasswordFull, error) {
	return user.GetUserPassword(orp, identifier)
}

func (orp ORPImpl) GetUserRefreshToken(identifier string, now time.Time) (*transferUser.UserRefreshTokenFull, error) {
	return user.GetUserRefreshToken(orp, identifier, now)
}

func (orp ORPImpl) GetUserAuthentic(identifier string, now time.Time) (*transferUser.UserAuthentic, error) {
	return user.GetUserAuthentic(orp, identifier, now)
}

func (orp ORPImpl) GetUserAuthenticOfCompany(identifier string, now time.Time) ([]transferUser.UserAuthentic, error) {
	return user.GetUserAuthenticOfCompany(orp, identifier, now)
}

func (orp ORPImpl) AddEmail(userEmail *transferUser.UserEmail, now time.Time) (*transferUser.UserEmail, error) {
	var err = orp.checkTransaction()
	if err != nil {
		return nil, err
	}
	return command.AddEmail(orp, userEmail, now)
}

func (orp ORPImpl) VerifyEmail(userEmail *transferUser.UserEmail, now time.Time) (*transferUser.UserEmail, error) {
	var err = orp.checkTransaction()
	if err != nil {
		return nil, err
	}
	return command.VerifyEmail(orp, userEmail, now)
}

func (orp ORPImpl) AddPassword(userPassword *transferUser.UserPassword, now time.Time) (*transferUser.UserPassword, error) {
	var err = orp.checkTransaction()
	if err != nil {
		return nil, err
	}
	return command.AddPassword(orp, userPassword, now)
}

func (orp ORPImpl) AddRefreshToken(userRefreshToken transferUser.UserRefreshToken, now time.Time) (transferUser.UserRefreshToken, error) {
	var err = orp.checkTransaction()
	if err != nil {
		return transferUser.UserRefreshToken{}, err
	}
	return command.AddRefreshToken(orp, userRefreshToken, now)
}
