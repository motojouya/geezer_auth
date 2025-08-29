package company_test

import (
	//"errors"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type roleAssignerDBMock struct {
	getUserAuthentic func(identifier string, now time.Time) (*dbUser.UserAuthentic, error)
	dbUtility.SqlExecutorMock
}

func (mock roleAssignerDBMock) GetUserAuthentic(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthentic(identifier, now)
}

func getShelterCompanyForRoleAssign(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func getShelterRoleForRoleAssign(expectLabel string) shelterRole.Role {
	var label, _ = pkgText.NewLabel(expectLabel)
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var roleRegisteredDate = time.Now()

	return shelterRole.NewRole(roleName, label, description, roleRegisteredDate)
}

func getShelterUserAuthenticForRoleAssign(expectUserId string, expectCompanyId string, expectLabel string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectUserId)
	var emailId, _ = pkgText.NewEmail("test@exaple.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = shelterUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier(expectCompanyId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()
	var company = shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)

	var label, _ = pkgText.NewLabel(expectLabel)
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var roleRegisteredDate = time.Now()

	var roles = []shelterRole.Role{shelterRole.NewRole(roleName, label, description, roleRegisteredDate)}
	var companyRole = shelterUser.NewCompanyRole(company, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func getDbUserAuthenticForUserRoleAssigner(expectUserId string, expectCompanyId string, expectOldLabel string, expectNewLabel string) *dbUser.UserAuthentic {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)

	var userCompanyRole1 = &dbUser.UserCompanyRoleFull{
		UserCompanyRole: dbUser.UserCompanyRole{
			PersistKey:        1,
			UserPersistKey:    2,
			CompanyPersistKey: 3,
			RoleLabel:         expectOldLabel,
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
	var userCompanyRole2 = &dbUser.UserCompanyRoleFull{
		UserCompanyRole: dbUser.UserCompanyRole{
			PersistKey:        2,
			UserPersistKey:    2,
			CompanyPersistKey: 3,
			RoleLabel:         expectNewLabel,
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
		RoleName:              "TestRolaName",
		RoleDescription:       "TestRolaDescription",
		RoleRegisteredDate:    now.Add(5 * time.Hour),
	}
	var userCompanyRoles = []dbUser.UserCompanyRoleFull{*userCompanyRole1, *userCompanyRole2}

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

func getLocalerMockForRoleAssign(now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow:       getNow,
	}
}

func getRoleAssignDbMock(t *testing.T, expectUserId string, expectCompanyId string, expectOldLabel string, expectNewLabel string, firstNow time.Time) roleAssignerDBMock {
	var dbUserAuthentic = getDbUserAuthenticForUserRoleAssigner(expectUserId, expectCompanyId, expectOldLabel, expectNewLabel)
	var getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, expectUserId, identifier)
		return dbUserAuthentic, nil
	}
	var insert = func(args ...interface{}) error {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		userCompanyRole, ok := args[0].(*dbUser.UserCompanyRole)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbUser.User, got %T", args[0])
		}

		assert.NotNil(t, userCompanyRole)
		assert.Equal(t, expectNewLabel, userCompanyRole.RoleLabel)

		return nil
	}
	return roleAssignerDBMock{
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeInsert: insert,
		},
		getUserAuthentic: getUserAuthentic,
	}
}

func TestRoleAssigner(t *testing.T) {
	var expectUserId = "US-TESTES"
	var expectCompanyId = "CP-TESTES"
	var expectOldLabel = "TEST_ROLE"
	var expectNewLabel = "TEST_ROLA"
	var firstNow = time.Now()

	var userAuthentic = getShelterUserAuthenticForRoleAssign(expectUserId, expectCompanyId, expectOldLabel)
	var shelterCompany = getShelterCompanyForRoleAssign(expectCompanyId)
	var role = getShelterRoleForRoleAssign(expectOldLabel)

	var localerMock = getLocalerMockForRoleAssign(firstNow)
	var dbMock = getRoleAssignDbMock(t, expectUserId, expectCompanyId, expectOldLabel, expectNewLabel, firstNow)

	assigner := company.NewRoleAssign(localerMock, dbMock)
	result, err := assigner.Execute(shelterCompany, userAuthentic, role)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectUserId, string(result.Identifier))
	assert.Equal(t, expectCompanyId, string(result.CompanyRole.Company.Identifier))
	assert.Equal(t, expectOldLabel, string(result.CompanyRole.Roles[0].Label))
	assert.Equal(t, expectNewLabel, string(result.CompanyRole.Roles[1].Label))

	t.Logf("Result: %+v", result)
}
