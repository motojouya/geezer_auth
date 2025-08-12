package user_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
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

func getDbUserAuthenticForNameChanger(id string) *dbUser.UserAuthentic {
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
		UserIdentifier:     id,
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &email,
		UserCompanyRole:    userCompanyRoles,
	}
}

func getShelterUser(expectId string, expectName string) shelterUser.UnsavedUser {
	var identifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName(expectName)
	var botFlag = false
	var registeredDate = time.Now()

	return shelterUser.CreateUser(identifier, emailId, name, botFlag, registeredDate)
}

func getLocalerMockForNameChange(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow:               getNow,
	}
}

func getNameChangeDbMock(t *testing.T, expectId string, expectName string, firstNow time.Time) userCreatorDBMock {
	var getUser = func(identifier string) (*dbUser.User, error) {
		assert.Equal(t, "US-TESTES", identifier, "Expected identifier 'US-TESTES'")
		return &dbUser.User{}, nil
	}
	var insert = func(args ...interface{}) error {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		user, ok := args[0].(*dbUser.User)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbUser.User, got %T", args[0])
		}

		assert.NotNil(t, user, "Expected user to be not nil")
		assert.Equal(t, "US-TESTES", user.Identifier, "Expected user identifier 'US-TESTES'")

		return nil
	}
	var dbUserAuthentic = getDbUserAuthenticForCreator(expectId)
	var getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, "US-TESTES", identifier, "Expected identifier 'US-TESTES'")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return dbUserAuthentic, nil
	}
	return userCreatorDBMock{
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeInsert: insert,
		},
		getUserAuthentic: getUserAuthentic,
		getUser:          getUser,
	}
}

func getEntryNameChangeMock(t *testing.T, expectId string, expectName string, firstNow time.Time) userNameEntryMock {
	var shelterUserVal = getShelterUser(expectId, expectName)
	var applyShelterUser = func(userArg shelterUser.User, now time.Time) (shelterUser.User, error) {
		assert.Equal(t, expectId, string(userArg.Identifier), "Expected identifier 'US-TESTES'")
		assert.Equal(t, expectName, string(userArg.Name), "Expected user name 'TestName'")
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
	var localerMock = getLocalerMockForNameChange(t, firstNow)
	var dbMock = getNameChangeDbMock(t, expectId, expectName, firstNow)
	var entryMock = getEntryNameChangeMock(t, expectId, expectName, firstNow)

	changer := user.NewNameChange(localerMock, dbMock)
	result, err := changer.Execute(entryMock)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "US-TESTES", string(result.Identifier), "Expected user identifier 'US-TESTES'")
	assert.Equal(t, expectNewName, string(result.Name), "Expected user name to be updated to 'TestNewName'")

	t.Logf("User created: %+v", result)
}
