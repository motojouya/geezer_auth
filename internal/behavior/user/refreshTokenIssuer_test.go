package user_test

import (
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/behavior/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"errors"
)

type refreshTokenIssuerDBMock struct {
	addRefreshToken func(userRefreshToken dbUser.UserRefreshToken, now time.Time) (dbUser.UserRefreshToken, error)
}

func (mock refreshTokenIssuerDBMock) AddRefreshToken(userRefreshToken dbUser.UserRefreshToken, now time.Time) (dbUser.UserRefreshToken, error) {
	return mock.addRefreshToken(userRefreshToken, now)
}

func getShelterUserAuthenticForRefToken() *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
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

func getLocalerMockForRefToken(t *testing.T, expectUUID uuid.UUID, now time.Time) *testUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	var generateUUID = func() (uuid.UUID, error) {
		return expectUUID, nil
	}
	return &testUtility.LocalerMock{
		FakeGetNow: getNow,
		FakeGenerateUUID: generateUUID,
	}
}

func getRefreshTokenIssueDbMock(t *testing.T, expectToken string, firstNow time.Time) refreshTokenIssuerDBMock {
	var addRefreshToken = func(userRefreshToken dbUser.UserRefreshToken, now time.Time) (dbUser.UserRefreshToken, error) {
		assert.Equal(t, userRefreshToken.RefreshToken, expectToken, "Expected token to match")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return userRefreshToken, nil
	}
	return refreshTokenIssuerDBMock{
		addRefreshToken: addRefreshToken,
	}
}

func TestRefreshTokenIssuer(t *testing.T) {
	var expectUUID, _ = uuid.NewUUID()
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForRefToken()

	var localerMock = getLocalerMockForRefToken(t, expectUUID, firstNow)
	var dbMock = getRefreshTokenIssueDbMock(t, expectUUID.String(), firstNow)

	setter := user.NewRefreshTokenIssue(localerMock, dbMock)
	refreshToken, err := setter.Execute(userAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectUUID.String(), string(refreshToken), "Expected refresh token to match generated UUID")
}

func TestRefreshTokenIssuerErrGenerateUUID(t *testing.T) {
	var expectUUID, _ = uuid.NewUUID()
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForRefToken()

	var localerMock = getLocalerMockForRefToken(t, expectUUID, firstNow)
	var dbMock = getRefreshTokenIssueDbMock(t, expectUUID.String(), firstNow)

	localerMock.FakeGenerateUUID = func() (uuid.UUID, error) {
		return uuid.Nil, errors.New("UUID generation error")
	}

	setter := user.NewRefreshTokenIssue(localerMock, dbMock)
	_, err := setter.Execute(userAuthentic)

	assert.Error(t, err)
}

func TestRefreshTokenIssuerErrAddToken(t *testing.T) {
	var expectUUID, _ = uuid.NewUUID()
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForRefToken()

	var localerMock = getLocalerMockForRefToken(t, expectUUID, firstNow)
	var dbMock = getRefreshTokenIssueDbMock(t, expectUUID.String(), firstNow)

	dbMock.addRefreshToken = func(userRefreshToken dbUser.UserRefreshToken, now time.Time) (dbUser.UserRefreshToken, error) {
		return dbUser.UserRefreshToken{}, errors.New("Database error")
	}

	setter := user.NewRefreshTokenIssue(localerMock, dbMock)
	_, err := setter.Execute(userAuthentic)

	assert.Error(t, err)
}
