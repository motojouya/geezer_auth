package user_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
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

type nameChangerDBMock struct {
	getUserAuthentic func(identifier string, now time.Time) (*dbUser.UserAuthentic, error)
	dbUtility.SqlExecutorMock
}

func (mock nameChangerDBMock) GetUserAuthentic(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthentic(identifier, now)
}

type userNameEntryMock struct {
	applyShelterUser func(user shelterUser.User, now time.Time) (shelterUser.User, error)
}

func (mock userNameEntryMock) ApplyShelterUser(user shelterUser.User, now time.Time) (shelterUser.User, error) {
	return mock.applyShelterUser(user, now)
}

func getDbUserAuthenticForNameChanger(id string, name string) *dbUser.UserAuthentic {
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
		UserIdentifier:        id,
		UserExposeEmailId:     "test02@example.com",
		UserName:              name,
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
		UserIdentifier:     id,
		UserExposeEmailId:  "test02@example.com",
		UserName:           name,
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &email,
		UserCompanyRole:    userCompanyRoles,
	}
}

func getShelterUserAuthenticForNameChange(expectId string, expectName string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName(expectName)
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

func getShelterUserForNameChanger(expectId string, expectName string) shelterUser.User {
	var identifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName(expectName)
	var botFlag = false
	var registeredDate = time.Now()

	return shelterUser.NewUser(1, identifier, name, emailId, botFlag, registeredDate, registeredDate)
}

func getLocalerMockForNameChange(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow: getNow,
	}
}

func getNameChangeDbMock(t *testing.T, expectId string, expectName string, firstNow time.Time) nameChangerDBMock {
	var update = func(args ...interface{}) (int64, error) {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		user, ok := args[0].(*dbUser.User)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbUser.User, got %T", args[0])
		}

		assert.NotNil(t, user)
		assert.Equal(t, expectId, user.Identifier)
		assert.Equal(t, expectName, user.Name)

		return 0, nil
	}
	var dbUserAuthentic = getDbUserAuthenticForNameChanger(expectId, expectName)
	var getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, "US-TESTES", identifier, "Expected identifier 'US-TESTES'")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return dbUserAuthentic, nil
	}
	return nameChangerDBMock{
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeUpdate: update,
		},
		getUserAuthentic: getUserAuthentic,
	}
}

func getEntryNameChangeMock(t *testing.T, expectId string, expectOldName string, expectNewName string, firstNow time.Time) userNameEntryMock {
	var shelterUserVal = getShelterUserForNameChanger(expectId, expectNewName)
	var applyShelterUser = func(userArg shelterUser.User, now time.Time) (shelterUser.User, error) {
		assert.Equal(t, expectId, string(userArg.Identifier), "Expected identifier 'US-TESTES'")
		assert.Equal(t, expectOldName, string(userArg.Name), "Expected user name 'TestName'")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return shelterUserVal, nil
	}
	return userNameEntryMock{
		applyShelterUser: applyShelterUser,
	}
}

func TestNameChanger(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldName = "TestOldName"
	var expectNewName = "TestNewName"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForNameChange(expectId, expectOldName)
	var localerMock = getLocalerMockForNameChange(t, firstNow)
	var dbMock = getNameChangeDbMock(t, expectId, expectNewName, firstNow)
	var entryMock = getEntryNameChangeMock(t, expectId, expectOldName, expectNewName, firstNow)

	changer := user.NewNameChange(localerMock, dbMock)
	result, err := changer.Execute(entryMock, userAuthentic)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectId, string(result.Identifier), "Expected user identifier 'US-TESTES'")
	assert.Equal(t, expectNewName, string(result.Name), "Expected user name to be updated to 'TestNewName'")

	t.Logf("User created: %+v", result)
}

func TestNameChangerErrApply(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldName = "TestOldName"
	var expectNewName = "TestNewName"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForNameChange(expectId, expectOldName)

	var localerMock = getLocalerMockForNameChange(t, firstNow)
	var dbMock = getNameChangeDbMock(t, expectId, expectNewName, firstNow)
	var entryMock = getEntryNameChangeMock(t, expectId, expectOldName, expectNewName, firstNow)

	entryMock.applyShelterUser = func(userArg shelterUser.User, now time.Time) (shelterUser.User, error) {
		return shelterUser.User{}, errors.New("apply error")
	}

	changer := user.NewNameChange(localerMock, dbMock)
	_, err := changer.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestNameChangerErrUpdate(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldName = "TestOldName"
	var expectNewName = "TestNewName"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForNameChange(expectId, expectOldName)

	var localerMock = getLocalerMockForNameChange(t, firstNow)
	var dbMock = getNameChangeDbMock(t, expectId, expectNewName, firstNow)
	var entryMock = getEntryNameChangeMock(t, expectId, expectOldName, expectNewName, firstNow)

	dbMock.SqlExecutorMock.FakeUpdate = func(args ...interface{}) (int64, error) {
		return 0, errors.New("update error")
	}

	changer := user.NewNameChange(localerMock, dbMock)
	_, err := changer.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestNameChangerErrGetAuthentic(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldName = "TestOldName"
	var expectNewName = "TestNewName"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForNameChange(expectId, expectOldName)

	var localerMock = getLocalerMockForNameChange(t, firstNow)
	var dbMock = getNameChangeDbMock(t, expectId, expectNewName, firstNow)
	var entryMock = getEntryNameChangeMock(t, expectId, expectOldName, expectNewName, firstNow)

	dbMock.getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, errors.New("get authentic error")
	}

	changer := user.NewNameChange(localerMock, dbMock)
	_, err := changer.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestNameChangerErrGetAuthenticNil(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldName = "TestOldName"
	var expectNewName = "TestNewName"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForNameChange(expectId, expectOldName)

	var localerMock = getLocalerMockForNameChange(t, firstNow)
	var dbMock = getNameChangeDbMock(t, expectId, expectNewName, firstNow)
	var entryMock = getEntryNameChangeMock(t, expectId, expectOldName, expectNewName, firstNow)

	dbMock.getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, nil
	}

	changer := user.NewNameChange(localerMock, dbMock)
	_, err := changer.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}
