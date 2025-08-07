package user_test

import (
	"github.com/motojouya/geezer_auth/internal/behavior/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"errors"
)

type emailSetterDBMock struct {
	getUserEmail func(email string) ([]dbUser.UserEmailFull, error)
	addEmail func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error)
}

func (mock emailSetterDBMock) GetUserEmail(email string) ([]dbUser.UserEmailFull, error) {
	return mock.getUserEmail(email)
}

func (mock emailSetterDBMock) AddEmail(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
	return mock.addEmail(userEmail, now)
}

type emailGetEntryMock struct {
	getEmail func() (pkgText.Email, error)
}

func (mock emailGetEntryMock) GetEmail() (pkgText.Email, error) {
	return mock.getEmail()
}

func getShelterUserAuthenticForEmail(id string) *shelterUser.UserAuthentic {
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

func getLocalerMockForEmail(t *testing.T, now time.Time) *testUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	var generateUUID = func() (uuid.UUID, error) {
		return uuid.NewUUID()
	}
	return &testUtility.LocalerMock{
		FakeGenerateUUID: generateUUID,
		FakeGetNow:       getNow,
	}
}

func getEmailSetDbMock(t *testing.T, expectEmail string, firstNow time.Time) emailSetterDBMock {
	var getUserEmail = func(email string) ([]dbUser.UserEmailFull, error) {
		assert.Equal(t, expectEmail, email)
		return []dbUser.UserEmailFull{}, nil
	}
	var addEmail = func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
		assert.Equal(t, expectEmail, userEmail.Email)
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return userEmail, nil
	}
	return emailSetterDBMock{
		getUserEmail: getUserEmail,
		addEmail:     addEmail,
	}
}

func getGetEmailEntryMock(t *testing.T, expectEmail string, firstNow time.Time) emailGetEntryMock {
	var email, _ = pkgText.NewEmail(expectEmail)
	var getEmail = func() (pkgText.Email, error) {
		return email, nil
	}
	return emailGetEntryMock{
		getEmail: getEmail,
	}
}

func TestEmailSetter(t *testing.T) {
	var idStr = "US-TESTES"
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(idStr)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.NoError(t, err)
}
