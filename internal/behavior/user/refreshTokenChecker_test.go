package user_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type refreshTokenCheckerDBMock struct {
	getUserRefreshToken func(token string, now time.Time) (*dbUser.UserAuthentic, error)
}

func (mock refreshTokenCheckerDBMock) GetUserRefreshToken(token string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserRefreshToken(token, now)
}

func getDbUserAuthenticForRefTokenChecker() *dbUser.UserAuthentic {
	var companyId = "CP-TESTES"
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)

	var userCompanyRole1 = &dbUser.UserCompanyRoleFull{
		UserCompanyRole: dbUser.UserCompanyRole{
			PersistKey:        1,
			UserPersistKey:    2,
			CompanyPersistKey: 3,
			RoleLabel:         "TEST_ROLE",
			RegisterDate:      now,
			ExpireDate:        &expireDate,
		},
		UserIdentifier:        "US-TESTES",
		UserExposeEmailId:     "test02@example.com",
		UserName:              "TestUserName",
		UserBotFlag:           false,
		UserRegisteredDate:    now.Add(2 * time.Hour),
		UserUpdateDate:        now.Add(3 * time.Hour),
		CompanyIdentifier:     companyId,
		CompanyName:           "TestCompanyName",
		CompanyRegisteredDate: now.Add(4 * time.Hour),
		RoleName:              "TestRoleName",
		RoleDescription:       "TestRoleDescription",
		RoleRegisteredDate:    now.Add(5 * time.Hour),
	}
	var userCompanyRoles = []dbUser.UserCompanyRoleFull{*userCompanyRole1}

	var email = "test01@example.com"
	return &dbUser.UserAuthentic{
		UserPersistKey:     2,
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &email,
		UserCompanyRole:    userCompanyRoles,
	}
}

func getLocalerMockForRefTokenCheck(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow:       getNow,
	}
}

func getRefreshTokenCheckerDbMock(t *testing.T, expectToken string, firstNow time.Time) refreshTokenCheckerDBMock {
	var getUserRefreshToken = func(token string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, token, expectToken, "Expected token to match")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return getDbUserAuthenticForRefTokenChecker(), nil
	}
	return refreshTokenCheckerDBMock{
		getUserRefreshToken: getUserRefreshToken,
	}
}

type userRefreshTokenGetterMock struct {
	getRefreshToken func() (shelterText.Token, error)
}

func (getter userRefreshTokenGetterMock) GetRefreshToken() (shelterText.Token, error) {
	return getter.getRefreshToken()
}

func getUserRefreshTokenGetterMock(expectToken string) userRefreshTokenGetterMock {
	var getRefreshToken = func() (shelterText.Token, error) {
		return shelterText.NewToken(expectToken)
	}
	return userRefreshTokenGetterMock{
		getRefreshToken: getRefreshToken,
	}
}

func TestRefreshTokenChecker(t *testing.T) {
	var firstNow = time.Now()
	var expectToken = "refresh_token01"

	var localerMock = getLocalerMockForRefTokenCheck(t, firstNow)
	var dbMock = getRefreshTokenCheckerDbMock(t, expectToken, firstNow)
	var entryMock = getUserRefreshTokenGetterMock(expectToken)

	checker := user.NewRefreshTokenCheck(localerMock, dbMock)
	userAuthentic, err := checker.Execute(entryMock)

	assert.NoError(t, err)
	assert.NotNil(t, userAuthentic, "Expected user authentic to be returned")
}

func TestRefreshTokenCheckerErrEntry(t *testing.T) {
	var firstNow = time.Now()
	var expectToken = "refresh_token01"

	var localerMock = getLocalerMockForRefTokenCheck(t, firstNow)
	var dbMock = getRefreshTokenCheckerDbMock(t, expectToken, firstNow)
	var entryMock = getUserRefreshTokenGetterMock(expectToken)

	entryMock.getRefreshToken = func() (shelterText.Token, error) {
		return "", errors.New("error getting refresh token")
	}

	checker := user.NewRefreshTokenCheck(localerMock, dbMock)
	_, err := checker.Execute(entryMock)

	assert.Error(t, err)
}

func TestRefreshTokenCheckerErrDB(t *testing.T) {
	var firstNow = time.Now()
	var expectToken = "refresh_token01"

	var localerMock = getLocalerMockForRefTokenCheck(t, firstNow)
	var dbMock = getRefreshTokenCheckerDbMock(t, expectToken, firstNow)
	var entryMock = getUserRefreshTokenGetterMock(expectToken)

	dbMock.getUserRefreshToken = func(token string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, errors.New("error getting user refresh token from DB")
	}

	checker := user.NewRefreshTokenCheck(localerMock, dbMock)
	_, err := checker.Execute(entryMock)

	assert.Error(t, err)
}
