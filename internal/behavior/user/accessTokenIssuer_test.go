package user_test

import (
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/behavior/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	//"errors"
)

type accessTokenIssuerDBMock struct {
	getUserAccessToken func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error)
	testUtility.SqlExecutorMock
}

func (mock accessTokenIssuerDBMock) GetUserAccessToken(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
	return mock.getUserAccessToken(identifier, now)
}

func getShelterUserAuthenticForAccToken() *shelterUser.UserAuthentic {
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

func getLocalerMockForAccToken(t *testing.T, expectUUID uuid.UUID, now time.Time) *testUtility.LocalerMock {
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

func getAccessTokenIssueDbMock(t *testing.T, expectId string, expectToken string, firstNow time.Time) accessTokenIssuerDBMock {
	var insert = func(userAccessTokens ...interface{}) error {
		assert.Equal(t, userAccessTokens[0].AccessToken, expectToken, "Expected token to match")
		return userRefreshToken, nil
	}
	var getUserAccessToken = func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
		assert.Equal(t, expectId, identifier, "Expected identifier to match")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return []dbUser.UserAccessTokenFull{}, nil
	}
	return accessTokenIssuerDBMock{
		SqlExecutorMock: testUtility.SqlExecutorMock{
			FakeInsert: insert,
		},
		getUserAccessToken: getUserAccessToken,
	}
}

func getJwtHandlerMock(t *testing.T, expectId string, expectToken string, expectUUID string, firstNow time.Time) *testUtility.JwtHandlerMock {
	var generate = func(user pkgUser.User, now time.Time, tokenId string) (pkgUser.Authentic, pkgText.JwtToken, error) {
		assert.Equal(t, expectId, user.Identifier, "Expected token ID to match")
		assert.Equal(t, expectUUID, tokenId, "Expected token ID to match")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return shelterUser.TokenData{
			ExpiresAt: shelterUser.NewExpiresAt(now.Add(time.Hour)),
		}, pkgText.JwtToken(expectToken), nil
	}
	return &testUtility.JwtHandlerMock{
		FakeGenerate: generate,
	}
}

func TestRefreshTokenIssuer(t *testing.T) {
	var expectUUID, _ = uuid.NewUUID()
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken()

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getRefreshTokenIssueDbMock(t, expectUUID.String(), firstNow)

	setter := user.NewRefreshTokenIssue(localerMock, dbMock)
	refreshToken, err := setter.Execute(userAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectUUID.String(), string(refreshToken), "Expected refresh token to match generated UUID")
}
