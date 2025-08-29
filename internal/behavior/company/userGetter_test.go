package company_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type userGetterDBMock struct {
	getUserAuthenticOfCompanyUser func(companyIdentifier string, userIdentifier string, now time.Time) (*dbUser.UserAuthentic, error)
}

func (mock userGetterDBMock) GetUserAuthenticOfCompanyUser(companyIdentifier string, userIdentifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthenticOfCompanyUser(companyIdentifier, userIdentifier, now)
}

func getDbUserAuthenticForUserGetter(expectCompanyId string, expectUserId string) *dbUser.UserAuthentic {
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
		UserIdentifier:        expectUserId,
		UserExposeEmailId:     "test02@example.com",
		UserName:              "TestUserName",
		UserBotFlag:           false,
		UserRegisteredDate:    now.Add(2 * time.Hour),
		UserUpdateDate:        now.Add(3 * time.Hour),
		CompanyIdentifier:     expectCompanyId,
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
		UserIdentifier:     expectUserId,
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &email,
		UserCompanyRole:    userCompanyRoles,
	}
}

func getLocalerMockForUserGet(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow: getNow,
	}
}

func getUserGetterDbMock(t *testing.T, expectCompanyId string, expectUserId string, firstNow time.Time) *userGetterDBMock {
	var dbUserAuthentic = getDbUserAuthenticForUserGetter(expectCompanyId, expectUserId)
	var getUserAuthenticOfCompanyUser = func(companyIdentifier string, userIdentifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, expectUserId, userIdentifier, "Expected identifier 'US-TESTES'")
		assert.Equal(t, expectCompanyId, companyIdentifier, "Expected identifier 'US-TESTES'")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return dbUserAuthentic, nil
	}
	return &userGetterDBMock{
		getUserAuthenticOfCompanyUser: getUserAuthenticOfCompanyUser,
	}
}

type companyUserGetterMock struct {
	getUserIdentifier func() (pkgText.Identifier, error)
}

func (mock companyUserGetterMock) GetUserIdentifier() (pkgText.Identifier, error) {
	return mock.getUserIdentifier()
}

func getCompanyUserGetterMock(t *testing.T, expectId string) companyUserGetterMock {
	var getUserIdentifier = func() (pkgText.Identifier, error) {
		return pkgText.NewIdentifier(expectId)
	}
	return companyUserGetterMock{
		getUserIdentifier: getUserIdentifier,
	}
}

func getShelterCompanyForGetter(expectId string) shelterCompany.Company {
	var companyId uint = 1
	var identifier, _ = pkgText.NewIdentifier(expectId)
	var name, _ = pkgText.NewName("TestRole")
	var registeredDate = time.Now()

	return shelterCompany.NewCompany(companyId, identifier, name, registeredDate)
}

func TestUserGetter(t *testing.T) {
	var expectCompanyId = "CP-TESTES"
	var expectUserId = "US-TESTES"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForUserGet(t, firstNow)
	var dbMock = getUserGetterDbMock(t, expectCompanyId, expectUserId, firstNow)
	var entryMock = getCompanyUserGetterMock(t, expectUserId)
	var shelterCompany = getShelterCompanyForGetter(expectCompanyId)

	getter := company.NewUserGet(localerMock, dbMock)
	result, err := getter.Execute(entryMock, shelterCompany)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectUserId, string(result.Identifier))
	assert.Equal(t, expectCompanyId, string(result.CompanyRole.Company.Identifier))

	t.Logf("User created: %+v", result)
}

func TestUserGetterNil(t *testing.T) {
	var expectCompanyId = "CP-TESTES"
	var expectUserId = "US-TESTES"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForUserGet(t, firstNow)
	var dbMock = getUserGetterDbMock(t, expectCompanyId, expectUserId, firstNow)
	var entryMock = getCompanyUserGetterMock(t, expectUserId)
	dbMock.getUserAuthenticOfCompanyUser = func(companyIdentifier string, userIdentifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, nil
	}

	var shelterCompany = getShelterCompanyForGetter(expectCompanyId)

	getter := company.NewUserGet(localerMock, dbMock)
	result, err := getter.Execute(entryMock, shelterCompany)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestUserGetterErrEntry(t *testing.T) {
	var expectCompanyId = "CP-TESTES"
	var expectUserId = "US-TESTES"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForUserGet(t, firstNow)
	var dbMock = getUserGetterDbMock(t, expectCompanyId, expectUserId, firstNow)
	var entryMock = getCompanyUserGetterMock(t, expectUserId)
	entryMock.getUserIdentifier = func() (pkgText.Identifier, error) {
		return pkgText.Identifier(""), errors.New("test error")
	}

	var shelterCompany = getShelterCompanyForGetter(expectCompanyId)

	getter := company.NewUserGet(localerMock, dbMock)
	_, err := getter.Execute(entryMock, shelterCompany)

	assert.Error(t, err)
}

func TestUserGetterErrGet(t *testing.T) {
	var expectCompanyId = "CP-TESTES"
	var expectUserId = "US-TESTES"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForUserGet(t, firstNow)
	var dbMock = getUserGetterDbMock(t, expectCompanyId, expectUserId, firstNow)
	var entryMock = getCompanyUserGetterMock(t, expectUserId)
	dbMock.getUserAuthenticOfCompanyUser = func(companyIdentifier string, userIdentifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, errors.New("test error")
	}

	var shelterCompany = getShelterCompanyForGetter(expectCompanyId)

	getter := company.NewUserGet(localerMock, dbMock)
	_, err := getter.Execute(entryMock, shelterCompany)

	assert.Error(t, err)
}

func TestUserGetterErrTrans(t *testing.T) {
	var expectCompanyId = "CP-TESTES"
	var expectUserId = "US-TESTES"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForUserGet(t, firstNow)
	var dbMock = getUserGetterDbMock(t, expectCompanyId, expectUserId, firstNow)
	var entryMock = getCompanyUserGetterMock(t, expectUserId)
	dbMock.getUserAuthenticOfCompanyUser = func(companyIdentifier string, userIdentifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return &dbUser.UserAuthentic{}, errors.New("test error")
	}

	var shelterCompany = getShelterCompanyForGetter(expectCompanyId)

	getter := company.NewUserGet(localerMock, dbMock)
	_, err := getter.Execute(entryMock, shelterCompany)

	assert.Error(t, err)
}
