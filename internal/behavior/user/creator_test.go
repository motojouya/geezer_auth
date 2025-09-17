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

type userCreatorDBMock struct {
	getUser          func(identifier string) (*dbUser.User, error)
	getUserAuthentic func(identifier string, now time.Time) (*dbUser.UserAuthentic, error)
	dbUtility.SqlExecutorMock
}

func (mock userCreatorDBMock) GetUser(identifier string) (*dbUser.User, error) {
	return mock.getUser(identifier)
}

func (mock userCreatorDBMock) GetUserAuthentic(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthentic(identifier, now)
}

type userEntryMock struct {
	toShelterUser func(identifier pkgText.Identifier, now time.Time) (shelterUser.UnsavedUser, error)
}

func (mock userEntryMock) ToShelterUser(identifier pkgText.Identifier, now time.Time) (shelterUser.UnsavedUser, error) {
	return mock.toShelterUser(identifier, now)
}

func getDbUserAuthenticForCreator(id string) *dbUser.UserAuthentic {
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

func getShelterUser(id string) shelterUser.UnsavedUser {
	var identifier, _ = pkgText.NewIdentifier(id)
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var name, _ = pkgText.NewName("TestName")
	var botFlag = false
	var registeredDate = time.Now()

	return shelterUser.CreateUser(identifier, emailId, name, botFlag, registeredDate)
}

func getLocalerMockForCreate(t *testing.T, idStr string, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	var generateRamdomString = func(length int, charSet string) string {
		return "TESTES"
	}
	return &localUtility.LocalerMock{
		FakeGenerateRamdomString: generateRamdomString,
		FakeGetNow:               getNow,
	}
}

func getCreateDbMock(t *testing.T, idStr string, firstNow time.Time) userCreatorDBMock {
	var getUser = func(identifier string) (*dbUser.User, error) {
		assert.Equal(t, "US-TESTES", identifier, "Expected identifier 'US-TESTES'")
		return nil, nil
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
	var dbUserAuthentic = getDbUserAuthenticForCreator(idStr)
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

func getEntryMock(t *testing.T, idStr string, firstNow time.Time) userEntryMock {
	var shelterUserVal = getShelterUser(idStr)
	var toShelterUser = func(identifier pkgText.Identifier, now time.Time) (shelterUser.UnsavedUser, error) {
		assert.Equal(t, "US-TESTES", string(identifier), "Expected identifier 'US-TESTES'")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return shelterUserVal, nil
	}
	return userEntryMock{
		toShelterUser: toShelterUser,
	}
}

func TestUserCreate(t *testing.T) {
	var idStr = "US-TESTES"
	var firstNow = time.Now()
	var localerMock = getLocalerMockForCreate(t, idStr, firstNow)
	var dbMock = getCreateDbMock(t, idStr, firstNow)
	var entryMock = getEntryMock(t, idStr, firstNow)

	creator := user.NewUserCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "US-TESTES", string(result.Identifier), "Expected user identifier 'US-TESTES'")

	t.Logf("User created: %+v", result)
}

func TestUserCreateErrGetUser(t *testing.T) {
	var idStr = "US-TESTES"
	var firstNow = time.Now()
	var localerMock = getLocalerMockForCreate(t, idStr, firstNow)
	var dbMock = getCreateDbMock(t, idStr, firstNow)
	var entryMock = getEntryMock(t, idStr, firstNow)

	dbMock.getUser = func(identifier string) (*dbUser.User, error) {
		return nil, errors.New("GetUser error")
	}

	creator := user.NewUserCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserCreateErrToShelter(t *testing.T) {
	var idStr = "US-TESTES"
	var firstNow = time.Now()
	var localerMock = getLocalerMockForCreate(t, idStr, firstNow)
	var dbMock = getCreateDbMock(t, idStr, firstNow)
	var entryMock = getEntryMock(t, idStr, firstNow)

	entryMock.toShelterUser = func(identifier pkgText.Identifier, now time.Time) (shelterUser.UnsavedUser, error) {
		return shelterUser.UnsavedUser{}, errors.New("ToShelterUser error")
	}

	creator := user.NewUserCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserCreateErrInsert(t *testing.T) {
	var idStr = "US-TESTES"
	var firstNow = time.Now()
	var localerMock = getLocalerMockForCreate(t, idStr, firstNow)
	var dbMock = getCreateDbMock(t, idStr, firstNow)
	var entryMock = getEntryMock(t, idStr, firstNow)

	dbMock.FakeInsert = func(args ...interface{}) error {
		return errors.New("Insert error")
	}

	creator := user.NewUserCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserCreateErrGetAuthentic(t *testing.T) {
	var idStr = "US-TESTES"
	var firstNow = time.Now()
	var localerMock = getLocalerMockForCreate(t, idStr, firstNow)
	var dbMock = getCreateDbMock(t, idStr, firstNow)
	var entryMock = getEntryMock(t, idStr, firstNow)

	dbMock.getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, errors.New("GetUserAuthentic error")
	}

	creator := user.NewUserCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserCreateErrGetAuthenticNil(t *testing.T) {
	var idStr = "US-TESTES"
	var firstNow = time.Now()
	var localerMock = getLocalerMockForCreate(t, idStr, firstNow)
	var dbMock = getCreateDbMock(t, idStr, firstNow)
	var entryMock = getEntryMock(t, idStr, firstNow)

	dbMock.getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, nil
	}

	creator := user.NewUserCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.Error(t, err)
	assert.Nil(t, result)
}
