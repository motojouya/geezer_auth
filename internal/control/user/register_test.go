package user_test

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/behavior/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"errors"
)

type userCreatorMock struct {
	execute func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error)
}

func (mock userCreatorMock) Execute(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
	return mock.execute(entry)
}

type emailSetterMock struct {
	execute func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error
}

func (emailSetterMock) Execute(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
	return mock.execute(entry, user)
}

type passwordSetterMock struct {
	execute func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error
}

func (mock passwordSetterMock) Execute(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
	return mock.execute(entry, user)
}

type refreshTokenIssuerMock struct {
	execute func(user *shelterUser.UserAuthentic) (shelterText.Token, error)
}

func (mock refreshTokenIssuerMock) Execute(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
	return mock.execute(user)
}

type accessTokenIssuerMock struct {
	execute func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error)
}

func (mock accessTokenIssuerMock) Execute(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
	return mock.execute(user)
}

type TransactionalDatabaseMock struct {
	FakeBegin func() error
	FakeRollback func() error
	FakeCommit func() error
	FakeClose func() error
}

func (mock TransactionalDatabaseMock) Begin() error {
	return mock.FakeBegin()
}

func (mock TransactionalDatabaseMock) Rollback() error {
	return mock.FakeRollback()
}

func (mock TransactionalDatabaseMock) Commit() error {
	return mock.FakeCommit()
}

func (mock TransactionalDatabaseMock) Close() error {
	return mock.FakeClose()
}

func getBehavior() (*userCreatorMock, *emailSetterMock, *passwordSetterMock, *refreshTokenIssuerMock, *accessTokenIssuerMock) {
	var userCreator = &userCreatorMock{
		execute: func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
			return getShelterUserAuthenticForAccToken(string(entry.GetIdentifier())), nil
		},
	}

	var emailSetter = &emailSetterMock{
		execute: func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
			return nil
		},
	}

	var passwordSetter = &passwordSetterMock{
		execute: func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
			return nil
		},
	}

	var refreshTokenIssuer = &refreshTokenIssuerMock{
		execute: func(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
			return shelterText.Token("test-refresh-token"), nil
		},
	}

	var accessTokenIssuer = &accessTokenIssuerMock{
		execute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
			return pkgText.JwtToken("test-access-token"), nil
		},
	}

	return userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer
}

func getDB() TransactionalDatabaseMock {
	return TransactionalDatabaseMock{
		FakeBegin: func() error {
			return nil
		},
		FakeRollback: func() error {
			return nil
		},
		FakeCommit: func() error {
			return nil
		},
		FakeClose: func() error {
			return nil
		},
	}
}

func getShelterUserAuthenticForRegister(expectId string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = shelterUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()
	var company = shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var roleRegisteredDate = time.Now()

	var roles = []shelterRole.Role{shelterRole.NewRole(roleName, label, description, roleRegisteredDate)}
	var companyRole = shelterUser.NewCompanyRole(company, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func TestRegister(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	accessToken, err := issuer.Execute(userAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectToken, string(accessToken), "Expected refresh token to match generated UUID")
}
