package company_test

import (
	//"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type allUserGetterDBMock struct {
	getUserAuthenticOfCompany func(identifier string, now time.Time) ([]dbUser.UserAuthentic, error)
}

func (mock allUserGetterDBMock) GetUserAuthenticOfCompany(identifier string, now time.Time) ([]dbUser.UserAuthentic, error) {
	return mock.getUserAuthenticOfCompany(identifier, now)
}

func getDbUserAuthenticForAllUserGetter(expectId string) *dbUser.UserAuthentic {
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
		CompanyIdentifier:     expectId,
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

func getLocalerMockForAllUserGet(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow: getNow,
	}
}

func getAllUserGetterDbMock(t *testing.T, expectId string, firstNow time.Time) allUserGetterDBMock {
	var dbUserAuthentic = getDbUserAuthenticForAllUserGetter(expectId)
	var getUserAuthenticOfCompany = func(identifier string, now time.Time) ([]dbUser.UserAuthentic, error) {
		assert.Equal(t, expectId, identifier)
		assert.WithinDuration(t, now, firstNow, time.Second)
		return []dbUser.UserAuthentic{*dbUserAuthentic}, nil
	}
	return allUserGetterDBMock{
		getUserAuthenticOfCompany: getUserAuthenticOfCompany,
	}
}

type companyGetterEntryMockForUser struct {
	getCompanyIdentifier func() (pkgText.Identifier, error)
}

func (mock companyGetterEntryMockForUser) GetCompanyIdentifier() (pkgText.Identifier, error) {
	return mock.getCompanyIdentifier()
}

func getCompanyGetterEntryMockForUser(expectId string) companyGetterEntryMockForUser {
	var getCompanyIdentifier = func() (pkgText.Identifier, error) {
		return pkgText.NewIdentifier(expectId)
	}
	return companyGetterEntryMockForUser{
		getCompanyIdentifier: getCompanyIdentifier,
	}
}

func TestAllUserGetter(t *testing.T) {
	var expectId = "CP-TESTES"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForAllUserGet(t, firstNow)
	var dbMock = getAllUserGetterDbMock(t, expectId, firstNow)
	var entryMock = getCompanyGetterEntryMockForUser(expectId)

	getter := company.NewAllUserGet(localerMock, dbMock)
	result, err := getter.Execute(entryMock)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, expectId, string(result[0].CompanyRole.Company.Identifier))

	t.Logf("User: %+v", result)
}
